package xribble

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	defaultFilePerm          os.FileMode = 0666
	defaultDirectoryFilePerm             = 0755
)

type FileSystem interface {
	IsDirectory(dir string) bool
	IsFile(file string) bool
	Write(path string, data []byte) error
	Read(path string) ([]byte, error)
	CreateDirectory(dir string) error
	Flush(dir string) error
	Remove(path string) error
}

type XribbleIO struct {
	filePermission os.FileMode
	dirPermission  os.FileMode
}

func NewXribbleIO() *XribbleIO {
	return &XribbleIO{defaultFilePerm, defaultDirectoryFilePerm}
}

func (x *XribbleIO) Read(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

func (x *XribbleIO) IsDirectory(dir string) bool {
	_, err := os.Stat(dir)

	if err != nil {
		return false
	}

	return true

}

func (x *XribbleIO) IsFile(file string) bool {
	_, err := os.Lstat(file)

	if err != nil {
		return false
	}

	return true
}

func (x *XribbleIO) Write(path string, data []byte) error {

	if err := x.CreateDirectory(filepath.Dir(path)); err != nil {
		return err
	}

	return ioutil.WriteFile(path, data, x.filePermission)
}

func (x *XribbleIO) CreateDirectory(dir string) error {
	return os.MkdirAll(dir, x.dirPermission)
}

func (x *XribbleIO) Flush(dir string) error {
	return os.RemoveAll(dir)
}

func (x *XribbleIO) Remove(path string) error {
	return os.Remove(path)
}
