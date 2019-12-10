package filestorage

import (
	"github.com/pkg/errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

func safeWrite(targetFile string, content []byte) error {
	filename := filepath.Base(targetFile)
	workDir := filepath.Dir(targetFile)
	err := os.MkdirAll(workDir, filePermission)
	if err != nil {
		return errors.Wrap(err, "create base dir")
	}

	file, err := ioutil.TempFile(workDir, filename+".tmp-*")
	if err != nil {
		return errors.Wrap(err, "create temp file")
	}
	err = ioutil.WriteFile(file.Name(), content, filePermission)
	_ = file.Close()
	if err != nil {

		return errors.Wrap(err, "write temp file")
	}

	return os.Rename(file.Name(), targetFile)
}
