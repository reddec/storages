package main

import (
	"bytes"
	"flag"
	"github.com/dave/jennifer/jen"
	"github.com/knq/snaker"
	"github.com/pkg/errors"
	"github.com/reddec/symbols"
	"io/ioutil"
)

var tpName = flag.String("type", "", "Type name to wrap")
var pkName = flag.String("package", "", "Output package (default: same as in input file)")
var ouName = flag.String("out", "", "Output file (default: <type name>_cached_storage.go)")

func main() {
	flag.Parse()
	if *tpName == "" {
		panic("type name should be specified")
	}

	if *ouName == "" {
		*ouName = snaker.CamelToSnake(*tpName) + "_cached_storage.go"
	}

	f, err := generate(".", *tpName, *pkName)
	if err != nil {
		panic(err)
	}
	data := &bytes.Buffer{}
	err = f.Render(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(*ouName, data.Bytes(), 0755)
	if err != nil {
		panic(err)
	}
}

func generate(dirName, typeName string, pkgName string) (*jen.File, error) {
	project, err := symbols.ProjectByDir(dirName)
	if err != nil {
		return nil, err
	}
	var file *jen.File
	if pkgName == "" {
		pkgName = project.Package.Package
		file = jen.NewFilePathName(project.Package.Import, pkgName)
	} else {
		file = jen.NewFile(pkgName)
	}

	sym := project.Package.FindSymbol(typeName)
	if sym == nil {
		return nil, errors.Errorf("symbol %v is not found", typeName)
	}

	symQual := jen.Qual(sym.Import.Import, typeName)
	stName := "Cached" + sym.Name + "Storage"
	file.Comment("Two level storage for " + typeName)
	file.Type().Id(stName).StructFunc(func(struc *jen.Group) {
		struc.Id("cold").Qual("github.com/reddec/storages", "Storage").Comment("persist storage")
		struc.Id("hot").Qual("github.com/reddec/storages", "Storage").Comment("cache storage")
		struc.Id("lock").Qual("sync", "RWMutex")
	})

	file.Line()
	file.Comment("Creates new storage for " + typeName + " with custom cache")
	file.Func().Id("New"+stName).Params(jen.Id("cold").Qual("github.com/reddec/storages", "Storage"), jen.Id("hot").Qual("github.com/reddec/storages", "Storage")).Op("*").Id(stName).BlockFunc(func(fn *jen.Group) {
		fn.Return(jen.Op("&").Id(stName).Values(jen.Id("cold").Op(":").Id("cold"), jen.Id("hot").Op(":").Id("hot")))
	})

	file.Line()
	file.Comment("Creates new storage for " + typeName + " with in-memory cache")
	file.Func().Id("NewMem" + stName).Params(jen.Id("cold").Qual("github.com/reddec/storages", "Storage")).Op("*").Id(stName).BlockFunc(func(fn *jen.Group) {
		fn.Return(jen.Id("New"+stName).Call(jen.Id("cold"), jen.Qual("github.com/reddec/storages/memstorage", "New").Call()))
	})

	file.Line()
	file.Comment("Put single " + typeName + " encoded in JSON into cold and hot storage")
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("Put").Params(jen.Id("key").String(), jen.Id("item").Op("*").Add(symQual)).Params(jen.Error()).BlockFunc(func(fn *jen.Group) {
		fn.List(jen.Id("data"), jen.Err()).Op(":=").Qual("encoding/json", "Marshal").Call(jen.Id("item"))
		fn.If(jen.Err().Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
			group.Return(jen.Err())
		})
		fn.Id("cs").Dot("lock").Dot("Lock").Call()
		fn.Defer().Id("cs").Dot("lock").Dot("Unlock").Call()
		fn.Err().Op("=").Id("cs").Dot("cold").Dot("Put").Call(jen.Index().Byte().Parens(jen.Id("key")), jen.Id("data"))
		fn.If(jen.Err().Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
			group.Return(jen.Err())
		})
		fn.Return(jen.Id("cs").Dot("hot").Dot("Put").Call(jen.Index().Byte().Parens(jen.Id("key")), jen.Id("data")))
	})

	file.Line()
	file.Comment("Get single " + typeName + " from hot storage and decode data as JSON. \nIf key is not in hot storage, the cold storage is used and obtained data is put to the hot storage for future cache")
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("Get").Params(jen.Id("key").String()).Params(jen.List(jen.Op("*").Add(symQual), jen.Error())).BlockFunc(func(fn *jen.Group) {
		fn.Id("cs").Dot("lock").Dot("RLock").Call()
		fn.List(jen.Id("data"), jen.Err()).Op(":=").Id("cs").Dot("hot").Dot("Get").Call(jen.Index().Byte().Parens(jen.Id("key")))
		fn.Id("cs").Dot("lock").Dot("RUnlock").Call()
		fn.If(jen.Err().Op("==").Qual("os", "ErrNotExist")).BlockFunc(func(group *jen.Group) {
			group.List(jen.Id("data"), jen.Err()).Op("=").Id("cs").Dot("getMissed").Call(jen.Id("key"))
		})
		fn.If(jen.Err().Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
			group.Return(jen.Nil(), jen.Err())
		})
		fn.Var().Id("item").Add(symQual)
		fn.Return(jen.Op("&").Id("item"), jen.Qual("encoding/json", "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("item")))
	})

	file.Line()
	file.Comment("Fetch all data from cold storage to the hot storage (warm cache)")
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("Fetch").Params().Error().BlockFunc(func(fn *jen.Group) {
		fn.Id("cs").Dot("lock").Dot("Lock").Call()
		fn.Defer().Id("cs").Dot("lock").Dot("Unlock").Call()
		fn.Return(jen.Id("cs").Dot("cold").Dot("Keys").CallFunc(func(group *jen.Group) {
			group.Func().Params(jen.Id("key").Index().Byte()).Error().BlockFunc(func(iterF *jen.Group) {
				iterF.List(jen.Id("data"), jen.Err()).Op(":=").Id("cs").Dot("cold").Dot("Get").Call(jen.Id("key"))
				iterF.If(jen.Err().Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
					group.Return(jen.Err())
				})
				iterF.Return(jen.Id("cs").Dot("hot").Dot("Put").Call(jen.Id("key"), jen.Id("data")))
			})
		}))
	})

	file.Line()
	file.Comment("Keys copied slice that cached in hot storage")
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("Keys").Params().Parens(jen.List(jen.Index().String(), jen.Error())).BlockFunc(func(fn *jen.Group) {
		fn.Id("cs").Dot("lock").Dot("RLock").Call()
		fn.Defer().Id("cs").Dot("lock").Dot("RUnlock").Call()
		fn.Var().Id("ans").Index().String()
		fn.Return(jen.Id("ans"), jen.Id("cs").Dot("hot").Dot("Keys").CallFunc(func(group *jen.Group) {
			group.Func().Params(jen.Id("key").Index().Byte()).Error().BlockFunc(func(iterF *jen.Group) {
				iterF.Id("ans").Op("=").Append(jen.Id("ans"), jen.String().Parens(jen.Id("key")))
				iterF.Return(jen.Nil())
			})
		}))
	})

	file.Line()
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("getMissed").Params(jen.Id("key").String()).Params(jen.List(jen.Index().Byte(), jen.Error())).BlockFunc(func(fn *jen.Group) {
		fn.Id("cs").Dot("lock").Dot("Lock").Call()
		fn.Defer().Id("cs").Dot("lock").Dot("Unlock").Call()
		fn.List(jen.Id("data"), jen.Err()).Op(":=").Id("cs").Dot("hot").Dot("Get").Call(jen.Index().Byte().Parens(jen.Id("key")))
		fn.If(jen.Err().Op("==").Nil()).BlockFunc(func(group *jen.Group) {
			group.Return(jen.Id("data"), jen.Nil())
		}).Else().If(jen.Err().Op("!=").Qual("os", "ErrNotExist")).BlockFunc(func(group *jen.Group) {
			group.Return(jen.Nil(), jen.Err())
		})
		fn.List(jen.Id("data"), jen.Err()).Op("=").Id("cs").Dot("cold").Dot("Get").Call(jen.Index().Byte().Parens(jen.Id("key")))
		fn.If(jen.Err().Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
			group.Return(jen.Nil(), jen.Err())
		})
		fn.Return(jen.Id("data"), jen.Id("cs").Dot("hot").Dot("Put").Call(jen.Index().Byte().Parens(jen.Id("key")), jen.Id("data")))
	})

	return file, nil
}
