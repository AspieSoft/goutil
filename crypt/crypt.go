package crypt

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"io"
	"math"
	"strconv"
	"time"

	"github.com/AspieSoft/go-regex-re2/v2"
)

type cryptCFB struct {}
type cryptHash struct {}

// Encryption: AES-CFB
var CFB cryptCFB

// Hashing: HMAC - SHA256
var Hash cryptHash

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
func RandBytes(size int, exclude ...[]byte) []byte {
	b := make([]byte, size)
	rand.Read(b)
	b = []byte(base64.URLEncoding.EncodeToString(b))

	if len(exclude) >= 2 {
		if exclude[0] == nil || len(exclude[0]) == 0 {
			b = regex.Comp(`[^\w_-]`).RepStrLit(b, exclude[1])
		}else{
			b = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStrLit(b, exclude[1])
		}
	}else if len(exclude) >= 1 {
		if exclude[0] == nil || len(exclude[0]) == 0 {
			b = regex.Comp(`[^\w_-]`).RepStrLit(b, []byte{})
		}else{
			b = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStrLit(b, []byte{})
		}
	}

	for len(b) < size {
		a := make([]byte, size)
		rand.Read(a)
		a = []byte(base64.URLEncoding.EncodeToString(a))
	
		if len(exclude) >= 2 {
			if exclude[0] == nil || len(exclude[0]) == 0 {
				a = regex.Comp(`[^\w_-]`).RepStrLit(a, exclude[1])
			}else{
				a = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStrLit(a, exclude[1])
			}
		}else if len(exclude) >= 1 {
			if exclude[0] == nil || len(exclude[0]) == 0 {
				a = regex.Comp(`[^\w_-]`).RepStrLit(a, []byte{})
			}else{
				a = regex.Comp(`[`+regex.Escape(string(exclude[0]))+`]`).RepStrLit(a, []byte{})
			}
		}

		b = append(b, a...)
	}

	return b[:size]
}

var uuidGenLastTime int64

// GenUUID generates a Unique Identifier using a custom build method
//
// Notice: This feature is currently in beta
//
// @size: (minimum: 8) the bit size for the last part of the uuid
// (note: other parts may vary)
//
// @timezone: optionally add a timezone string to the uuid
// (note: you could also pass random info into here for a more complex algorithm)
//
// This method uses the following data:
//  - A hash of the current year and day of year
//  - A hash of the current timezone
//  - A hash of the current unix time (in seconds)
//  - A hash of the current unix time in nanoseconds and a random number
//
// The returned value is url encoded and will look something like this: xxxx-xxxx-xxxx-xxxxxxxx
func GenUUID(size int, timezone ...string) string {
	for time.Now().UnixNano() <= uuidGenLastTime {
		time.Sleep(1 * time.Millisecond)
	}
	uuidGenLastTime = time.Now().UnixNano()

	if size < 8 {
		size = 8
	}

	uuid := [][]byte{{}, {}, {}, {}}

	// year
	{
		s := int(math.Min(float64(size/4), 8))
		if s < 4 {
			s = 4
		}

		sm := s/2
		if s % 2 != 0 {
			sm++
		}

		b := sha1.Sum([]byte(strconv.Itoa(time.Now().Year())))
		uuid[0] = []byte(base64.URLEncoding.EncodeToString(b[:]))[:sm]
		b = sha1.Sum([]byte(strconv.Itoa(time.Now().YearDay())))
		uuid[0] = append(uuid[0], []byte(base64.URLEncoding.EncodeToString(b[:]))[:sm]...)
		uuid[0] = uuid[0][:s]
	}

	// time zone
	{
		s := int(math.Min(float64(size/8), 8))
		if s < 4 {
			s = 4
		}

		if len(timezone) != 0 {
			sm := s/len(timezone)
			if s % 2 != 0 {
				sm++
			}

			for _, zone := range timezone {
				b := sha1.Sum([]byte(zone))
				uuid[1] = append(uuid[1], []byte(base64.URLEncoding.EncodeToString(b[:]))[:sm]...)
			}
			uuid[1] = uuid[1][:s]
		}else{
			z, _ := time.Now().Zone()
			b := sha1.Sum([]byte(z))
			uuid[1] = []byte(base64.URLEncoding.EncodeToString(b[:]))[:s]
		}
	}

	// unix time
	{
		s := int(math.Min(float64(size/2), 16))
		if s < 4 {
			s = 4
		}

		b := sha1.Sum([]byte(strconv.Itoa(int(time.Now().Unix()))))
		uuid[2] = []byte(base64.URLEncoding.EncodeToString(b[:]))[:s]
	}

	// random
	{
		s := int(math.Min(float64(size/4), 64))
		if s < 4 {
			s = 4
		}

		b := sha512.Sum512([]byte(strconv.Itoa(int(time.Now().UnixNano()))))
		uuid[3] = []byte(base64.URLEncoding.EncodeToString(b[:]))[:s]
		uuid[3] = append(uuid[3], []byte(base64.URLEncoding.EncodeToString(RandBytes(size)))[:size-s]...)
	}

	if len(uuid[1]) == 0 {
		uuid = append(uuid[:1], uuid[2:]...)
	}

	for i := range uuid {
		uuid[i] = bytes.ReplaceAll(uuid[i], []byte{'-'}, []byte{'0'})
		uuid[i] = bytes.ReplaceAll(uuid[i], []byte{'_'}, []byte{'1'})
	}

	return string(bytes.Join(uuid, []byte{'-'}))
}
