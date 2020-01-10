package utils

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

// see https://go-review.googlesource.com/c/go/+/1591/9/src/io/ioutil/ioutil.go#77
func CopyFile(dst, src string, perm os.FileMode) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), "")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		return err
	}
	tmp.Close()
	if err = os.Chmod(tmp.Name(), perm); err != nil {
		if err := tmp.Close(); err != nil {
			os.Remove(tmp.Name())
			return err
		}
		return err
	}
	return os.Rename(tmp.Name(), dst)
}
