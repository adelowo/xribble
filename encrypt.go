package xribble

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"errors"
	"io"
	"os"
)

type XribbleEncrypter struct {
	key []byte
}

//Sticking with AES-256 encryption here
//https://golang.org/pkg/crypto/aes/#NewCipher
const KEY_BYTES_LENGTH = 32

func NewXribbleEncrypter() *XribbleEncrypter {

	k := []byte(os.Getenv("XRIBBLE_KEY"))

	if binary.Size(k) != KEY_BYTES_LENGTH {
		panic("The key for encryption should be 32 bytes in length")
	}

	return &XribbleEncrypter{k}
}

func (x *XribbleEncrypter) Encrypt(text []byte) ([]byte, error) {

	var b []byte

	ciph, err := aes.NewCipher(x.key)

	if err != nil {
		return b, err
	}

	gcm, err := cipher.NewGCM(ciph)

	if err != nil {
		return b, err
	}

	nonce := make([]byte, gcm.NonceSize())

	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, text, nil), nil
}

func (x *XribbleEncrypter) Decrypt(cipherText []byte) ([]byte, error) {
	c, err := aes.NewCipher(x.key)

	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)

	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	if binary.Size(cipherText) < nonceSize {
		return nil, errors.New("encrypter: ciphertext too short")
	}

	nonce, ciphertext := cipherText[:nonceSize], cipherText[nonceSize:]

	return gcm.Open(nil, nonce, ciphertext, nil)
}
