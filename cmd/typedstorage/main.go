package main

import (
	"bytes"
	"flag"
	"github.com/dave/jennifer/jen"
	"github.com/knq/snaker"
	"io/ioutil"
	"os"
	"strings"
)

var (
	tpName        = flag.String("type", "", "Type name to wrap")
	pkName        = flag.String("package", "", "Output package (default: same as in input file)")
	ouName        = flag.String("out", "", "Output file (default: <type name>_storage.go)")
	codec         = flag.String("codec", "json", "Encoder/Decoder for the type: json, msgp")
	keyPrefix     = flag.String("prefix", "", "Custom key prefix")
	resultName    = flag.String("name", "Storage", "Name of result structure")
	importType    = flag.String("import", "", "Custom import for type name")
	interfaceName = flag.String("interface", "", "Make interface for type")
	interfaceKV   = flag.String("kv-interface", "", "Add KV-only subset for type (should be used together with --interface)")
)

func main() {
	flag.Parse()
	if *tpName == "" {
		panic("type name should be specified")
	}

	if *ouName == "" {
		*ouName = snaker.CamelToSnake(*tpName) + "_storage.go"
	}
	var key Keyer = &simpleKey{}
	if *keyPrefix != "" {
		key = &nsKey{*keyPrefix}
	}
	f, err := generate(".", *tpName, *pkName, *resultName, getCodec(*codec), key, *importType, *interfaceName, *interfaceKV)
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

type Keyer interface {
	ForStor() jen.Code
	ForView() jen.Code
	Filter() jen.Code
}

type simpleKey struct{}

func (sk *simpleKey) ForStor() jen.Code { return jen.Id("key") }
func (sk *simpleKey) ForView() jen.Code { return jen.String().Parens(jen.Id("key")) }
func (sk *simpleKey) Filter() jen.Code  { return jen.Null() }

type nsKey struct{ prefix string }

func (sk *nsKey) ForStor() jen.Code { return jen.Lit(*keyPrefix).Op("+").Id("key") }
func (sk *nsKey) ForView() jen.Code {
	return jen.String().Parens(jen.Id("key")).Index(jen.Lit(len(sk.prefix)), jen.Empty())
}
func (sk *nsKey) Filter() jen.Code {
	return jen.If(jen.Op("!").Qual("strings", "HasPrefix").Call(jen.String().Parens(jen.Id("key")), jen.Lit(sk.prefix))).BlockFunc(func(group *jen.Group) {
		group.Return(jen.Nil())
	})
}

type Codec interface {
	Encode() jen.Code
	Decode() jen.Code
	Header(typeName string) jen.Code
}

type jsonCodec struct{}

func (jc *jsonCodec) Encode() jen.Code {
	return jen.Qual("encoding/json", "Marshal").Call(jen.Id("item"))
}
func (jc *jsonCodec) Decode() jen.Code {
	return jen.Err().Op("=").Qual("encoding/json", "Unmarshal").Call(jen.Id("data"), jen.Op("&").Id("item"))
}
func (jc *jsonCodec) Header(typeName string) jen.Code {
	return jen.Null()
}

type msgCodec struct{}

func (jc *msgCodec) Encode() jen.Code {
	return jen.Id("item").Dot("MarshalMsg").Call(jen.Nil())
}
func (jc *msgCodec) Decode() jen.Code {
	return jen.List(jen.Id("_"), jen.Err()).Op("=").Id("item").Dot("UnmarshalMsg").Call(jen.Id("data"))
}
func (jc *msgCodec) Header(typeName string) jen.Code {
	return jen.Null()
}

func getCodec(name string) Codec {
	switch name {
	case "json":
		return &jsonCodec{}
	case "msgp":
		return &msgCodec{}
	default:
		panic("unknown codec " + name)
	}
}

func generate(dirName, typeName string, pkgName, resultName string, codec Codec, key Keyer, typeImport string, interfaceName string, kvSubset string) (*jen.File, error) {
	var file *jen.File

	file = jen.NewFile(pkgName)

	file.HeaderComment("Code generated by typedstorage. DO NOT EDIT.")
	file.Id("//go:generate " + strings.Join(os.Args, " ")).Line()
	var symQual *jen.Statement
	if typeImport != "" {
		typeName = strings.Split(typeName, ".")[1]
		symQual = jen.Qual(typeImport, typeName)
	} else {
		symQual = jen.Id(typeName)
	}

	stName := resultName
	var publicType = jen.Op("*").Id(stName)
	if interfaceName != "" {
		appendInterface(file, typeName, interfaceName, symQual)
		if kvSubset != "" {
			appendKVInterface(file, typeName, kvSubset, symQual)
		}
		stName = strings.ToLower(stName[:1]) + stName[1:]
		publicType = jen.Id(interfaceName)
	}

	file.Comment("Implementation of typed storage for " + typeName)
	file.Type().Id(stName).StructFunc(func(struc *jen.Group) {
		struc.Id("cold").Qual("github.com/reddec/storages", "Storage").Comment("persist storage")
	})

	file.Line()
	file.Comment("Creates new storage for " + typeName)
	file.Func().Id("New" + stName).Params(jen.Id("cold").Qual("github.com/reddec/storages", "Storage")).Add(publicType).BlockFunc(func(fn *jen.Group) {
		fn.Return(jen.Op("&").Id(stName).Values(jen.Id("cold").Op(":").Id("cold")))
	})

	file.Line()
	file.Comment("Put single " + typeName + " encoded in JSON into storage")
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("Put").Params(jen.Id("key").String(), jen.Id("item").Op("*").Add(symQual)).Params(jen.Error()).BlockFunc(func(fn *jen.Group) {
		fn.List(jen.Id("data"), jen.Err()).Op(":=").Add(codec.Encode())
		fn.If(jen.Err().Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
			group.Return(jen.Err())
		})
		fn.Err().Op("=").Id("cs").Dot("cold").Dot("Put").Call(jen.Index().Byte().Parens(key.ForStor()), jen.Id("data"))
		fn.Return(jen.Err())
	})

	file.Line()
	file.Comment("Get single " + typeName + " from storage and decode data as JSON")
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("Get").Params(jen.Id("key").String()).Parens(jen.List(jen.Op("*").Add(symQual), jen.Error())).BlockFunc(func(fn *jen.Group) {
		fn.List(jen.Id("data"), jen.Err()).Op(":=").Id("cs").Dot("cold").Dot("Get").Call(jen.Index().Byte().Parens(key.ForStor()))
		fn.If(jen.Err().Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
			group.Return(jen.Nil(), jen.Err())
		})
		fn.Var().Id("item").Add(symQual)
		fn.Add(codec.Decode())
		fn.If(jen.Err().Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
			group.Return(jen.Nil(), jen.Err())
		})
		fn.Return(jen.Op("&").Id("item"), jen.Nil())
	})

	file.Line()
	file.Comment("Del key from hot and cold storage")
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("Del").Params(jen.Id("key").String()).Error().BlockFunc(func(fn *jen.Group) {
		fn.Err().Op(":=").Id("cs").Dot("cold").Dot("Del").Call(jen.Index().Byte().Parens(key.ForStor()))
		fn.Return(jen.Err())
	})

	file.Line()
	file.Comment("Keys copied slice that cached in hot storage")
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("Keys").Params().Parens(jen.List(jen.Index().String(), jen.Error())).BlockFunc(func(fn *jen.Group) {
		fn.Var().Id("ans").Index().String()
		fn.Return(jen.Id("ans"), jen.Id("cs").Dot("cold").Dot("Keys").Call(jen.Func().Params(jen.Id("key").Index().Byte()).Error().BlockFunc(func(group *jen.Group) {
			group.Add(key.Filter())
			group.Id("ans").Op("=").Append(jen.Id("ans"), key.ForView())
			group.Return(jen.Nil())
		})))
	})

	file.Line()
	file.Comment("Iterate over all items")
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("Iterate").Params(jen.Id("handler").Func().Params(jen.String(), jen.Op("*").Add(symQual)).Error()).Error().BlockFunc(func(fn *jen.Group) {
		fn.Return(jen.Id("cs").Dot("cold").Dot("Keys").Call(jen.Func().Params(jen.Id("key").Index().Byte()).Error().BlockFunc(func(group *jen.Group) {
			group.Add(key.Filter())
			group.List(jen.Id("item"), jen.Err()).Op(":=").Id("cs").Dot("Get").Call(key.ForView())
			group.If(jen.Err().Op("!=").Nil()).BlockFunc(func(group *jen.Group) {
				group.Return(jen.Err())
			})
			group.Return(jen.Id("handler").Call(jen.String().Parens(key.ForView()), jen.Id("item")))
		})))
	})

	file.Line()
	file.Comment("View makes a map of all items")
	file.Func().Parens(jen.Id("cs").Op("*").Id(stName)).Id("View").Params().Parens(jen.List(jen.Map(jen.String()).Op("*").Add(symQual), jen.Error())).BlockFunc(func(fn *jen.Group) {
		fn.Var().Id("ans").Op("=").Make(jen.List(jen.Map(jen.String()).Op("*").Add(symQual)))
		fn.Return(jen.Id("ans"), jen.Id("cs").Dot("Iterate").Call(jen.Func().Params(jen.Id("key").String(), jen.Id("item").Op("*").Add(symQual)).Error().BlockFunc(func(group *jen.Group) {
			group.Id("ans").Index(jen.Id("key")).Op("=").Id("item")
			group.Return(jen.Nil())
		})))
	})

	file.Line()
	file.Add(codec.Header(stName))

	return file, nil
}

func appendInterface(file *jen.File, typeName string, interfaceName string, symQual *jen.Statement) {
	file.Comment("Typed storage for " + typeName)
	file.Type().Id(interfaceName).InterfaceFunc(func(decl *jen.Group) {
		decl.Id("View").Params().Parens(jen.List(jen.Map(jen.String()).Op("*").Add(symQual), jen.Error()))
		decl.Id("Put").Params(jen.Id("key").String(), jen.Id("item").Op("*").Add(symQual)).Params(jen.Error())
		decl.Id("Get").Params(jen.Id("key").String()).Parens(jen.List(jen.Op("*").Add(symQual), jen.Error()))
		decl.Id("Del").Params(jen.Id("key").String()).Error()
		decl.Id("Keys").Params().Parens(jen.List(jen.Index().String(), jen.Error()))
		decl.Id("Iterate").Params(jen.Id("handler").Func().Params(jen.String(), jen.Op("*").Add(symQual)).Error()).Error()
		decl.Id("View").Params().Parens(jen.List(jen.Map(jen.String()).Op("*").Add(symQual), jen.Error()))
	})
	file.Line()
}
func appendKVInterface(file *jen.File, typeName string, interfaceName string, symQual *jen.Statement) {
	file.Comment("KV-only typed storage for " + typeName)
	file.Type().Id(interfaceName).InterfaceFunc(func(decl *jen.Group) {
		decl.Id("Put").Params(jen.Id("key").String(), jen.Id("item").Op("*").Add(symQual)).Params(jen.Error())
		decl.Id("Get").Params(jen.Id("key").String()).Parens(jen.List(jen.Op("*").Add(symQual), jen.Error()))
		decl.Id("Del").Params(jen.Id("key").String()).Error()
		decl.Id("Keys").Params().Parens(jen.List(jen.Index().String(), jen.Error()))
	})
	file.Line()
}
