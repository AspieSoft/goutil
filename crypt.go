package goutil

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"

	"github.com/AspieSoft/go-regex-re2/v2"
)

type crypt struct {
	// AES-CFB
	CFB cryptCFB

	// HMAC - SHA256
	Hash cryptHash
}

type cryptCFB struct {}
type cryptHash struct {}

// Encryption
var Crypt *crypt = &crypt{
	cryptCFB{},
	cryptHash{},
}

// Encrypt runs AES-CFB Encryption
//
// the key is also hashed with SHA256
func (crypt *cryptCFB) Encrypt(text []byte, key []byte) ([]byte, error) {
	keyHash := sha256.Sum256(key)

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return []byte{}, err
	}

	ciphertext := make([]byte, aes.BlockSize+len(text))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], text)

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// Decrypt runs AES-CFB Decryption
//
// the key is also hashed with SHA256
func (crypt *cryptCFB) Decrypt(text []byte, key []byte) ([]byte, error) {
	keyHash := sha256.Sum256(key)

	ciphertext, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return []byte{}, err
	}

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return []byte{}, err
	}

	if len(ciphertext) < aes.BlockSize {
		return []byte{}, errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)

	stream.XORKeyStream(ciphertext, ciphertext)
	return ciphertext, nil
}

// NewHash runs a key based HMAC hash using SHA256
//
// the key is also hashed with SHA256
func (crypt *cryptHash) New(text []byte, key []byte) ([]byte, error) {
	keyHash := sha256.Sum256(key)

	mac := hmac.New(sha256.New, keyHash[:])

	_, err := mac.Write(text)
	if err != nil {
		return []byte{}, err
	}

	return mac.Sum(nil), nil
}

// CompareHash compares a key based hash created by the `NewHash` func to another text, to safely check for equality
//
// uses HMAC with SHA256
//
// the key is also hashed with SHA256
//
// @compare should be a valid hash
func (crypt *cryptHash) Compare(text []byte, compare []byte, key []byte) bool {
	keyHash := sha256.Sum256(key)

	mac := hmac.New(sha256.New, keyHash[:])

	_, err := mac.Write(text)
	if err != nil {
		return false
	}

	return hmac.Equal(compare, mac.Sum(nil))
}


// RandBytes generates random bytes using crypto/rand
//
// @exclude[0] allows you can to pass an optional []byte to ensure that set of chars will not be included in the output string
//
// @exclude[1] provides a replacement string to put in place of the unwanted chars
//
// @exclude[2:] is currently ignored
func (crypt *crypt) RandBytes(size int, exclude ...[]byte) []byte {
	b := make([]byte, size)
	rand.Read(b)
	b = []byte(base64.URLEncoding.EncodeToString(b))

	if len(exclude) >= 2 {
		if exclude[0] == nil || len(exclude[0]) == 0 {
			b = regex.Comp(`[^\w_-]`).RepStr(b, exclude[1])
		}else{
			b = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStr(b, exclude[1])
		}
	}else if len(exclude) >= 1 {
		if exclude[0] == nil || len(exclude[0]) == 0 {
			b = regex.Comp(`[^\w_-]`).RepStr(b, []byte{})
		}else{
			b = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStr(b, []byte{})
		}
	}

	for len(b) < size {
		a := make([]byte, size)
		rand.Read(a)
		a = []byte(base64.URLEncoding.EncodeToString(a))
	
		if len(exclude) >= 2 {
			if exclude[0] == nil || len(exclude[0]) == 0 {
				a = regex.Comp(`[^\w_-]`).RepStr(a, exclude[1])
			}else{
				a = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStr(a, exclude[1])
			}
		}else if len(exclude) >= 1 {
			if exclude[0] == nil || len(exclude[0]) == 0 {
				a = regex.Comp(`[^\w_-]`).RepStr(a, []byte{})
			}else{
				a = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStr(a, []byte{})
			}
		}

		b = append(b, a...)
	}

	return b[:size]
}
