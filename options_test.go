package xribble

import (
	"reflect"
	"testing"
)

func Test_BaseDir(t *testing.T) {
	db := NewXribble(BaseDir("./skrr"))

	defer db.Drop()

	if dir := "./skrr"; db.baseDir != dir {
		t.Fatalf(
			`Base directory differs.. Expected %s \n Got %s`,
			dir, db.baseDir)
	}

}

type MockIO struct {
}

func (x *MockIO) Read(path string) ([]byte, error) {
	return nil, nil
}

func (x *MockIO) IsDirectory(dir string) bool {
	return true
}

func (x *MockIO) IsFile(file string) bool {
	return true
}

func (x *MockIO) Write(path string, data []byte) error {
	return nil
}

func (x *MockIO) CreateDirectory(dir string) error {
	return nil
}

func (x *MockIO) Flush(dir string) error {
	return nil
}

func (x *MockIO) Remove(path string) error {
	return nil
}

func Test_FS(t *testing.T) {
	m := &MockIO{}
	db := NewXribble(FS(m))

	defer db.Drop()

	if !reflect.DeepEqual(m, db.fs) {
		t.Fatalf(
			`FileSystem implementation differs. \n.. Expected %v.. \n 
			GOt %v`, m, db.fs)
	}
}
