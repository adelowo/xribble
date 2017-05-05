package xribble

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	defaultFilePerm          os.FileMode = 0666
	defaultDirectoryFilePerm             = 0755
	baseDirectory            string      = "./xribbled"
	pathSuffix               string      = ".xml"
)

type XribbleDriver struct {
	mu      sync.RWMutex
	baseDir string
	e       Encrypter
	fs      FileSystem
}

type FileSystem interface {
	IsDirectory(dir string) bool
	Write(path string, data []byte) error
	CreateDirectory(dir string) error
}

type XribbleIO struct {
	filePermission os.FileMode
	dirPermission  os.FileMode
}

func NewXribbleIO() *XribbleIO {
	return &XribbleIO{defaultFilePerm, defaultDirectoryFilePerm}
}

func (x *XribbleIO) IsDirectory(dir string) bool {
	_, err := os.Stat(dir)

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

type Item struct {
	Key       string    `xml:"info>key"`
	Value     []byte    `xml:"info>value"`
	ExpiresAt time.Time `xml:"info>expiration"`
}

type Encrypter interface {
	Encrypt(val string) (string, error)
	Decrypt(val string) (string, error)
}

type Provider interface {
	Add(i *Item) error
}

type Option func(*XribbleDriver)

func BaseDir(dir string) Option {
	return func(x *XribbleDriver) {
		x.baseDir = dir
	}
}

func NewXribble(opts ...Option) *XribbleDriver {

	d := &XribbleDriver{}

	for _, opt := range opts {
		if opt != nil {
			opt(d)
		}
	}

	if d.baseDir == "" {
		d.baseDir = baseDirectory
	}

	if d.fs == nil {
		d.fs = NewXribbleIO()
	}

	if !d.fs.IsDirectory(d.baseDir) {
		if err := d.fs.CreateDirectory(d.baseDir); err != nil {
			panic("Could not create the storage path for xribble")
		}
	}

	return d
}

func (x *XribbleDriver) Add(i *Item) error {
	x.mu.Lock()
	defer x.mu.Unlock()

	output, err := xml.MarshalIndent(i, " ", "  ")

	if err != nil {
		return err
	}

	return x.fs.Write(x.path(i.Key), output)
}

func (x *XribbleDriver) path(key string) string {

	hash := md5.Sum([]byte(key))

	hashSumAsString := hex.EncodeToString(hash[:])

	return filepath.Join(x.baseDir,
		string(hashSumAsString[0:2]),
		string(hashSumAsString[2:4]), hashSumAsString) + pathSuffix
}
