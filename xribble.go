package xribble

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"
)

const (
	baseDirectory string = "./xribbled"
	pathSuffix    string = ".xml"
)

var ErrDatabaseMiss error = errors.New(
	`xribble: Could not fetch the specified item as 
		it does not exist in the database`)

type XribbleDriver struct {
	mu      sync.RWMutex
	baseDir string
	e       Encrypter
	fs      FileSystem
}

type Item struct {
	Key       string    `xml:"info>key"`
	Value     []byte    `xml:"info>value"`
	ExpiresAt time.Time `xml:"info>expiration"`
}

type Encrypter interface {
	Encrypt(val []byte) ([]byte, error)
	Decrypt(cipherText []byte) ([]byte, error)
}

type Provider interface {
	Add(i *Item) error
	Get(key string) (*Item, error)
	Delete(key string) error
	Drop() error
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

func (x *XribbleDriver) Get(key string) (*Item, error) {
	x.mu.RLock()
	defer x.mu.RUnlock()

	path := x.path(key)

	if ok := x.fs.IsFile(path); !ok {
		return nil, ErrDatabaseMiss
	}

	data, err := x.fs.Read(path)

	if err != nil {
		return nil, err
	}

	i := new(Item)

	if err = xml.Unmarshal(data, i); err != nil {
		return i, err
	}

	return i, nil
}

func (x *XribbleDriver) Delete(key string) error {
	x.mu.Lock()
	defer x.mu.Unlock()

	path := x.path(key)

	if ok := x.fs.IsFile(path); ok {
		return x.fs.Remove(path)
	}

	return fmt.Errorf(
		`xribble: An error occurred while trying to delete the key %s`,
		key)
}

func (x *XribbleDriver) Drop() error {
	x.mu.Lock()
	defer x.mu.Unlock()

	return x.fs.Flush(x.baseDir)
}

func (x *XribbleDriver) path(key string) string {

	hash := md5.Sum([]byte(key))

	hashSumAsString := hex.EncodeToString(hash[:])

	return filepath.Join(x.baseDir,
		string(hashSumAsString[0:2]),
		string(hashSumAsString[2:4]), hashSumAsString) + pathSuffix
}
