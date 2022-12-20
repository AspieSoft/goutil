package goutil

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/AspieSoft/go-regex/v4"
	"github.com/alphadose/haxmap"
	"github.com/fsnotify/fsnotify"
)

type Hashable interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | string | complex64 | complex128
}

var VarType map[string]reflect.Type

func init(){
	VarType = map[string]reflect.Type{}

	VarType["array"] = reflect.TypeOf([]interface{}{})
	VarType["arrayByte"] = reflect.TypeOf([][]byte{})
	VarType["map"] = reflect.TypeOf(map[string]interface{}{})

	VarType["int"] = reflect.TypeOf(int(0))
	VarType["int64"] = reflect.TypeOf(int64(0))
	VarType["float64"] = reflect.TypeOf(float64(0))
	VarType["float32"] = reflect.TypeOf(float32(0))

	VarType["string"] = reflect.TypeOf("")
	VarType["byteArray"] = reflect.TypeOf([]byte{})
	VarType["byte"] = reflect.TypeOf([]byte{0}[0])

	// int 32 returned instead of byte
	VarType["int32"] = reflect.TypeOf(' ')

	VarType["func"] = reflect.TypeOf(func(){})

	VarType["bool"] = reflect.TypeOf(true)
}

// JoinPath joins multiple file types with safety from backtracking
func JoinPath(path ...string) (string, error) {
	resPath, err := filepath.Abs(string(path[0]))
	if err != nil {
		return "", err
	}
	for i := 1; i < len(path); i++ {
		p := filepath.Join(resPath, string(path[i]))
		if p == resPath || !strings.HasPrefix(p, resPath) {
			return "", errors.New("path leaked outside of root")
		}
		resPath = p
	}
	return resPath, nil
}

// Contains returns true if an array contains a value
func Contains[T any](search []T, value T) bool {
	val := ToString(value)
	for _, v := range search {
		if ToString(v) == val {
			return true
		}
	}
	return false
}

// IndexOf returns the index of a value in an array
//
// returns -1 and an error if the value is not found
func IndexOf[T any](search []T, value T) (int, error) {
	val := ToString(value)
	for i, v := range search {
		if ToString(v) == val {
			return i, nil
		}
	}
	return -1, errors.New("array does not contain value: " + ToString(value))
}

// ContainsMap returns true if a map contains a value
func ContainsMap[T Hashable, J any](search map[T]J, value J) bool {
	val := ToString(value)
	for _, v := range search {
		if ToString(v) == val {
			return true
		}
	}
	return false
}

// IndexOfMap returns the index of a value in a map
//
// returns an error if the value is not found
func IndexOfMap[T Hashable, J any](search map[T]J, value J) (T, error) {
	val := ToString(value)
	for i, v := range search {
		if ToString(v) == val {
			return i, nil
		}
	}
	var blk T
	return blk, errors.New("map does not contain value: " + ToString(value))
}

// ContainsMapKey returns true if a map contains a key
func ContainsMapKey[T Hashable, J any](search map[T]J, key T) bool {
	/* for i := range search {
		if i == key {
			return true
		}
	}
	return false */

	_, ok := search[key]
	return ok
}


// TrimRepeats trims repeating adjacent characters and reduces them to one character
//
// @b: byte array to trim
//
// @chars: list of bytes to trim repeats of
func TrimRepeats(b []byte, chars []byte) []byte {
	r := []byte{}
	for i := 0; i < len(b); i++ {
		r = append(r, b[i])
		if Contains(chars, b[i]) {
			for i+1 < len(b) && b[i+1] == b[i] {
				i++
			}
		}
	}
	return r
}


// ToString converts multiple types to a string
//
// accepts: string, []byte, byte, int32, int, int64, float64, float32
func ToString(res interface{}) string {
	switch reflect.TypeOf(res) {
		case VarType["string"]:
			return res.(string)
		case VarType["byteArray"]:
			return string(res.([]byte))
		case VarType["byte"]:
			return string(res.(byte))
		case VarType["int32"]:
			return string(res.(int32))
		case VarType["int"]:
			return strconv.Itoa(res.(int))
		case VarType["int64"]:
			return strconv.Itoa(int(res.(int64)))
		case VarType["float64"]:
			return strconv.FormatFloat(res.(float64), 'f', -1, 64)
		case VarType["float32"]:
			return strconv.FormatFloat(float64(res.(float32)), 'f', -1, 32)
		default:
			return ""
	}
}

// ToByteArray converts multiple types to a []byte
//
// accepts: string, []byte, byte, int32, int, int64, float64, float32
func ToByteArray(res interface{}) []byte {
	switch reflect.TypeOf(res) {
		case VarType["string"]:
			return []byte(res.(string))
		case VarType["byteArray"]:
			return res.([]byte)
		case VarType["byte"]:
			return []byte{res.(byte)}
		case VarType["int32"]:
			return []byte{byte(res.(int32))}
		case VarType["int"]:
			return []byte(strconv.Itoa(res.(int)))
		case VarType["int64"]:
			return []byte(strconv.Itoa(int(res.(int64))))
		case VarType["float64"]:
			return []byte(strconv.FormatFloat(res.(float64), 'f', -1, 64))
		case VarType["float32"]:
			return []byte(strconv.FormatFloat(float64(res.(float32)), 'f', -1, 32))
		default:
			return []byte{}
	}
}

// ToInt converts multiple types to an int
//
// accepts: int, int32, int64, float64, float32, string, []byte, byte
func ToInt(res interface{}) int {
	switch reflect.TypeOf(res) {
		case VarType["int"]:
			return res.(int)
		case VarType["int32"]:
			return int(res.(int32))
		case VarType["int64"]:
			return int(res.(int64))
		case VarType["float64"]:
			return int(res.(float64))
		case VarType["float32"]:
			return int(res.(float32))
		case VarType["string"]:
			if i, err := strconv.Atoi(res.(string)); err == nil {
				return i
			}
			return 0
		case VarType["byteArray"]:
			if i, err := strconv.Atoi(string(res.([]byte))); err == nil {
				return i
			}
			return 0
		case VarType["byte"]:
			if i, err := strconv.Atoi(string(res.(byte))); err == nil {
				return i
			}
			return 0
		default:
			return 0
	}
}

// IsZeroOfUnderlyingType can be used to determine if an interface{} in null or empty
func IsZeroOfUnderlyingType(x interface{}) bool {
	// return x == nil || x == reflect.Zero(reflect.TypeOf(x)).Interface()
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// FormatMemoryUsage converts bytes to megabytes
func FormatMemoryUsage(b uint64) float64 {
	return math.Round(float64(b) / 1024 / 1024 * 100) / 100
}


var regIsAlphaNumeric *regex.Regexp = regex.Compile(`^[A-Za-z0-9]+$`)


// MapArgs will convert a bash argument array ([]string) into a map (map[string]string)
//
// When @args is left blank with no values, it will default to os.Args[1:]
//
// -- Arg Convertions:
//
// "--Key=value" will convert to "key:value"
//
// "--boolKey" will convert to "boolKey:true"
//
// "-flags" will convert to "f:true, l:true, a:true, g:true, s:true" (only if its alphanumeric [A-Za-z0-9])
// if -flags is not alphanumeric (example: "-test.paniconexit0" "-test.timeout=10m0s") it will be treated as a --flag (--key=value --boolKey)
//
// keys that match a number ("--1" or "-1") will start with a "-" ("--1=value" -> "-1:value", "-1" -> -1:true)
// this prevents a number key from conflicting with an index key
//
// everything else is given a number value index starting with 0
//
// this method will not allow --args to have their values modified after they have already been set
func MapArgs(args ...[]string) map[string]string {
	if len(args) == 0 {
		args = append(args, os.Args[1:])
	}

	argMap := map[string]string{}
	i := 0

	for _, argList := range args {
		for _, arg := range argList {
			if strings.HasPrefix(arg, "--") {
				arg = arg[2:]
				if strings.ContainsRune(arg, '=') {
					data := strings.SplitN(arg, "=", 2)
					if _, err := strconv.Atoi(data[0]); err == nil {
						if argMap["-"+data[0]] == "" {
							argMap["-"+data[0]] = data[1]
						}
					}else{
						if argMap[data[0]] == "" {
							argMap[data[0]] = data[1]
						}
					}
				}else{
					if _, err := strconv.Atoi(arg); err == nil {
						if argMap["-"+arg] == "" {
							argMap["-"+arg] = "true"
						}
					}else{
						if argMap[arg] == "" {
							argMap[arg] = "true"
						}
					}
				}
			}else if strings.HasPrefix(arg, "-") {
				arg = arg[1:]
				if regIsAlphaNumeric.Match([]byte(arg)) {
					flags := strings.Split(arg, "")
					for _, flag := range flags {
						if _, err := strconv.Atoi(flag); err == nil {
							if argMap["-"+flag] == "" {
								argMap["-"+flag] = "true"
							}
						}else{
							if argMap[flag] == "" {
								argMap[flag] = "true"
							}
						}
					}
				}else{
					if strings.ContainsRune(arg, '=') {
						data := strings.SplitN(arg, "=", 2)
						if _, err := strconv.Atoi(data[0]); err == nil {
							if argMap["-"+data[0]] == "" {
								argMap["-"+data[0]] = data[1]
							}
						}else{
							if argMap[data[0]] == "" {
								argMap[data[0]] = data[1]
							}
						}
					}else{
						if _, err := strconv.Atoi(arg); err == nil {
							if argMap["-"+arg] == "" {
								argMap["-"+arg] = "true"
							}
						}else{
							if argMap[arg] == "" {
								argMap[arg] = "true"
							}
						}
					}
				}
			}else{
				argMap[strconv.Itoa(i)] = arg
				i++
			}
		}
	}

	return argMap
}

// MapArgs is just like MapArgs, but it excepts and outputs using []byte instead of string
func MapArgsByte(args ...[][]byte) map[string][]byte {
	if len(args) == 0 {
		args = append(args, [][]byte{})
		for _, arg := range os.Args[1:] {
			args[0] = append(args[0], []byte(arg))
		}
	}

	argMap := map[string][]byte{}
	i := 0

	for _, argList := range args {
		for _, arg := range argList {
			if bytes.HasPrefix(arg, []byte("--")) {
				arg = arg[2:]
				if bytes.ContainsRune(arg, '=') {
					data := bytes.SplitN(arg, []byte{'='}, 2)
					if _, err := strconv.Atoi(string(data[0])); err == nil {
						if argMap["-"+string(data[0])] == nil {
							argMap["-"+string(data[0])] = data[1]
						}
					}else{
						if argMap[string(data[0])] == nil {
							argMap[string(data[0])] = data[1]
						}
					}
				}else{
					if _, err := strconv.Atoi(string(arg)); err == nil {
						if argMap["-"+string(arg)] == nil {
							argMap["-"+string(arg)] = []byte("true")
						}
					}else{
						if argMap[string(arg)] == nil {
							argMap[string(arg)] = []byte("true")
						}
					}
				}
			}else if bytes.HasPrefix(arg, []byte{'-'}) {
				arg = arg[1:]
				if regIsAlphaNumeric.Match(arg) {
					flags := bytes.Split(arg, []byte{})
					for _, flag := range flags {
						if _, err := strconv.Atoi(string(flag)); err == nil {
							if argMap["-"+string(flag)] == nil {
								argMap["-"+string(flag)] = []byte("true")
							}
						}else{
							if argMap[string(flag)] == nil {
								argMap[string(flag)] = []byte("true")
							}
						}
					}
				}else{
					if bytes.ContainsRune(arg, '=') {
						data := bytes.SplitN(arg, []byte{'='}, 2)
						if _, err := strconv.Atoi(string(data[0])); err == nil {
							if argMap["-"+string(data[0])] == nil {
								argMap["-"+string(data[0])] = data[1]
							}
						}else{
							if argMap[string(data[0])] == nil {
								argMap[string(data[0])] = data[1]
							}
						}
					}else{
						if _, err := strconv.Atoi(string(arg)); err == nil {
							if argMap["-"+string(arg)] == nil {
								argMap["-"+string(arg)] = []byte("true")
							}
						}else{
							if argMap[string(arg)] == nil {
								argMap[string(arg)] = []byte("true")
							}
						}
					}
				}
			}else{
				argMap[strconv.Itoa(i)] = arg
				i++
			}
		}
	}

	return argMap
}

var regEscHTML *regex.Regexp = regex.Compile(`[<>&]`)
var regEscFixAmp *regex.Regexp = regex.Compile(`&amp;(amp;)*`)

// EscapeHTML replaces HTML characters with html entities
//
// Also prevents and removes &amp;amp; from results
func EscapeHTML(html []byte) []byte {
	html = regEscHTML.RepFuncRef(&html, func(data func(int) []byte) []byte {
		if bytes.Equal(data(0), []byte("<")) {
			return []byte("&lt;")
		} else if bytes.Equal(data(0), []byte(">")) {
			return []byte("&gt;")
		}
		return []byte("&amp;")
	})
	return regEscFixAmp.RepStrRef(&html, []byte("&amp;"))
}

var regEscHTMLArgs *regex.Regexp = regex.Compile(`([\\]*)([\\"'\'])`)

// EscapeHTMLArgs escapes quotes and backslashes for use within HTML quotes
// @quote can be used to only escape specific quotes or chars
func EscapeHTMLArgs(html []byte, quote ...byte) []byte {
	if len(quote) == 0 {
		quote = []byte("\"'`")
	}

	return regEscHTMLArgs.RepFuncRef(&html, func(data func(int) []byte) []byte {
		if len(data(1)) % 2 == 0 && bytes.ContainsRune(quote, rune(data(2)[0])) {
			// return append([]byte("\\"), data(2)...)
			return regex.JoinBytes(data(1), '\\', data(2))
		}
		if bytes.ContainsRune(quote, rune(data(2)[0])) {
			return append(data(1), data(2)...)
		}
		return data(0)
	})
}

// StringifyJSON converts a map or array to a JSON string
func StringifyJSON(data interface{}, ind ...int) ([]byte, error) {
	var res []byte
	var err error
	if len(ind) != 0 {
		sp := "  "
		if len(ind) > 2 {
			sp = strings.Repeat(" ", ind[1])
		}
		res, err = json.MarshalIndent(data, strings.Repeat(" ", ind[0]), sp)
	}else{
		res, err = json.Marshal(data)
	}

	if err != nil {
		return []byte{}, err
	}
	res = bytes.ReplaceAll(res, []byte("\\u003c"), []byte("<"))
	res = bytes.ReplaceAll(res, []byte("\\u003e"), []byte(">"))

	return res, nil
}

// ParseJson converts a json string into a map of strings
func ParseJson(b []byte) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return res, nil
}

// DecodeJSON is useful for decoding a JSON output from the body of an http request
//
// example: goutil.DecodeJSON(r.Body)
func DecodeJSON(data io.Reader) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := json.NewDecoder(data).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeepCopyJson will stringify and parse json to create a deep copy and escape pointers
func DeepCopyJson(data map[string]interface{}) (map[string]interface{}, error) {
	b, err := StringifyJSON(data)
	if err != nil {
		return nil, err
	}
	return ParseJson(b)
}

// Compress is Gzip compression for a string
func Compress(msg []byte) ([]byte, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(msg)); err != nil {
		return []byte{}, err
	}
	if err := gz.Flush(); err != nil {
		return []byte{}, err
	}
	if err := gz.Close(); err != nil {
		return []byte{}, err
	}
	return []byte(base64.StdEncoding.EncodeToString(b.Bytes())), nil
}

// Decompress is Gzip decompression for a string
func Decompress(str []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(str))
	if err != nil {
		return []byte{}, err
	}
	rdata := bytes.NewReader(data)
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return []byte{}, err
	}
	s, err := io.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}
	return s, nil
}


// Encrypt runs AES-CFB Encryption
//
// the key is also hashed with SHA256
func Encrypt(text []byte, key []byte) ([]byte, error) {
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
func Decrypt(text []byte, key []byte) ([]byte, error) {
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
func NewHash(text []byte, key []byte) ([]byte, error) {
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
func CompareHash(text []byte, compare []byte, key []byte) bool {
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
			b = regex.Compile(`[^\w_-]`).RepStr(b, exclude[1])
		}else{
			b = regex.Compile(`[%1]`, string(exclude[0])).RepStr(b, exclude[1])
		}
	}else if len(exclude) >= 1 {
		if exclude[0] == nil || len(exclude[0]) == 0 {
			b = regex.Compile(`[^\w_-]`).RepStr(b, []byte{})
		}else{
			b = regex.Compile(`[%1]`, string(exclude[0])).RepStr(b, []byte{})
		}
	}

	for len(b) < size {
		a := make([]byte, size)
		rand.Read(a)
		a = []byte(base64.URLEncoding.EncodeToString(a))
	
		if len(exclude) >= 2 {
			if exclude[0] == nil || len(exclude[0]) == 0 {
				a = regex.Compile(`[^\w_-]`).RepStr(a, exclude[1])
			}else{
				a = regex.Compile(`[%1]`, string(exclude[0])).RepStr(a, exclude[1])
			}
		}else if len(exclude) >= 1 {
			if exclude[0] == nil || len(exclude[0]) == 0 {
				a = regex.Compile(`[^\w_-]`).RepStr(a, []byte{})
			}else{
				a = regex.Compile(`[%1]`, string(exclude[0])).RepStr(a, []byte{})
			}
		}

		b = append(b, a...)
	}

	return b[:size]
}


const localEncKeyAdd string = "txavzc5CMtpmqERcdTQCbs6cBKAyYc/9hP/s3wLREZBfoiEB8Vc00//i27FQ3twTmW0jAWNiTjXkx1iDAklqCXT1lvyGbSjb2iftyQRLFgM="

// EncryptLocal is a Non Standard AES-CFB Encryption method
//
// Notice This Feature Is Experimental
//
// purposely incompatible with other libraries and programing languages
//
// this was made by accident, and this bug is now a feature
func EncryptLocal(text []byte, key []byte) ([]byte, error) {
	addKey, err := Decrypt([]byte(localEncKeyAdd), []byte("9EruID5odGcw9hSWSG19xPrx4hbG8ggbjYNdQmibfHAsnm2Y3oYUtGXbXfgXKmtx"))
	if err != nil {
		return []byte{}, err
	}

	bKey := append(key, addKey...)

	l := len(bKey)
	for l % 32 != 0 {
		l--
	}
	bKey = bKey[:l]

	var eKey []byte = nil
	if len(bKey) > 32 {
		eKey = bKey[32:]
		bKey = bKey[:32]
	}

	keyHash := sha256.Sum256(bKey)

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return []byte{}, err
	}
	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	if eKey != nil {
		return encryptLocal(ciphertext, eKey)
	}

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

func encryptLocal(text []byte, bKey []byte) ([]byte, error) {
	var eKey []byte = nil
	if len(bKey) > 32 {
		eKey = bKey[32:]
		bKey = bKey[:32]
	}

	keyHash := sha256.Sum256(bKey)

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return []byte{}, err
	}

	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return []byte{}, err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	if eKey != nil {
		return encryptLocal(ciphertext, eKey)
	}

	return []byte(base64.StdEncoding.EncodeToString(ciphertext)), nil
}

// DecryptLocal is a Non Standard AES-CFB Decryption method
//
// Notice This Feature Is Experimental
//
// purposely incompatible with other libraries and programing languages
//
// this was made by accident, and this bug is now a feature
func DecryptLocal(ciphertext []byte, key []byte) ([]byte, error) {
	addKey, err := Decrypt([]byte(localEncKeyAdd), []byte("9EruID5odGcw9hSWSG19xPrx4hbG8ggbjYNdQmibfHAsnm2Y3oYUtGXbXfgXKmtx"))
	if err != nil {
		return nil, err
	}

	bKey := append(key, addKey...)

	l := len(bKey)
	for l % 32 != 0 {
		l--
	}
	bKey = bKey[:l]

	var eKey []byte = nil
	if len(bKey) > 32 {
		l := len(bKey) - 32
		eKey = bKey[:l]
		bKey = bKey[l:]
	}

	keyHash := sha256.Sum256(bKey)

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return nil, err
	}

	text, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return nil, err
	}

	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)

	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}

	if eKey != nil {
		return decryptLocal(data, eKey)
	}

	return data, nil
}

func decryptLocal(text []byte, bKey []byte) ([]byte, error) {
	var eKey []byte = nil
	if len(bKey) > 32 {
		l := len(bKey) - 32
		eKey = bKey[:l]
		bKey = bKey[l:]
	}

	keyHash := sha256.Sum256(bKey)

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return nil, err
	}

	if len(text) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}
	iv := text[:aes.BlockSize]
	text = text[aes.BlockSize:]
	cfb := cipher.NewCFBDecrypter(block, iv)
	cfb.XORKeyStream(text, text)

	data, err := base64.StdEncoding.DecodeString(string(text))
	if err != nil {
		return nil, err
	}

	if eKey != nil {
		return decryptLocal(data, eKey)
	}

	return data, nil
}


// CleanStr will sanitizes a string to valid UTF-8
func CleanStr(str string) string {
	//todo: sanitize inputs
	str = strings.ToValidUTF8(str, "")
	return str
}

// CleanByte will sanitizes a []byte to valid UTF-8
func CleanByte(b []byte) []byte {
	//todo: sanitize inputs
	b = bytes.ToValidUTF8(b, []byte{})
	return b
}

// CleanArray runs `CleanStr` on an []interface{}
//
// CleanStr sanitizes a string to valid UTF-8
func CleanArray(data []interface{}) []interface{} {
	cData := []interface{}{}
	for key, val := range data {
		t := reflect.TypeOf(val)
		if t == VarType["string"] {
			cData[key] = CleanStr(val.(string))
		}else if t == VarType["int"] || t == VarType["float64"] || t == VarType["float32"] || t == VarType["bool"] {
			cData[key] = val
		}else if t == VarType["byteArray"] {
			cData[key] = CleanStr(string(val.([]byte)))
		}else if t == VarType["byte"] {
			cData[key] = CleanStr(string(val.(byte)))
		}else if t == VarType["int32"] {
			cData[key] = CleanStr(string(val.(int32)))
		}else if t == VarType["array"] {
			cData[key] = CleanArray(val.([]interface{}))
		}else if t == VarType["map"] {
			cData[key] = CleanMap(val.(map[string]interface{}))
		}
	}
	return cData
}

// CleanMap runs `CleanStr` on a map[string]interface{}
//
// CleanStr sanitizes a string to valid UTF-8
func CleanMap(data map[string]interface{}) map[string]interface{} {
	cData := map[string]interface{}{}
	for key, val := range data {
		key = CleanStr(key)

		t := reflect.TypeOf(val)
		if t == VarType["string"] {
			cData[key] = CleanStr(val.(string))
		}else if t == VarType["int"] || t == VarType["float64"] || t == VarType["float32"] || t == VarType["bool"] {
			cData[key] = val
		}else if t == VarType["byteArray"] {
			cData[key] = CleanStr(string(val.([]byte)))
		}else if t == VarType["byte"] {
			cData[key] = CleanStr(string(val.(byte)))
		}else if t == VarType["int32"] {
			cData[key] = CleanStr(string(val.(int32)))
		}else if t == VarType["array"] {
			cData[key] = CleanArray(val.([]interface{}))
		}else if t == VarType["map"] {
			cData[key] = CleanMap(val.(map[string]interface{}))
		}
	}

	return cData
}

// CleanJSON runs `CleanStr` on a complex json object recursively
//
// CleanStr sanitizes a string to valid UTF-8
func CleanJSON(val interface{}) interface{} {
	t := reflect.TypeOf(val)
	if t == VarType["string"] {
		return CleanStr(val.(string))
	}else if t == VarType["int"] || t == VarType["float64"] || t == VarType["float32"] || t == VarType["bool"] {
		return val
	}else if t == VarType["byteArray"] {
		return CleanByte(val.([]byte))
	}else if t == VarType["byte"] {
		return CleanByte([]byte{val.(byte)})
	}else if t == VarType["int32"] {
		return CleanStr(string(val.(int32)))
	}else if t == VarType["array"] {
		return CleanArray(val.([]interface{}))
	}else if t == VarType["map"] {
		return CleanMap(val.(map[string]interface{}))
	}
	return nil
}

var regDirEndSlash *regex.Regexp = regex.Compile(`[\\/][^\\/]*$`)

// GetFileFromParent checks if the parent (or sub parent) directory of a file contains a specific file or folder
//
// @root is the highest grandparent to check before quitting
//
// @start is the lowest level to start searching from (if a directory is passed, it will not be included in your search)
//
// @search is what file you want to search fro
func GetFileFromParent(root string, start string, search string) (string, bool) {
	dir := string(regDirEndSlash.RepStr([]byte(start), []byte{}))
	if len(dir) == 0 || dir == root || !strings.HasPrefix(dir, root) {
		return "", false
	}

	if dirList, err := os.ReadDir(string(dir)); err == nil {
		for _, file := range dirList {
			name := file.Name()
			if name == search {
				if path, err := JoinPath(string(dir), name); err == nil {
					return path, true
				}
				return "", false
			}
		}
	}

	return GetFileFromParent(root, dir, search)
}


// A watcher instance for the `WatchDir` func
type Watcher struct {
	FileChange func(path string, op string)
	DirAdd func(path string, op string) (addWatcher bool)
	Remove func(path string, op string) (removeWatcher bool)
	Any func(path string, op string)
}

type watcherObj struct {
	root string
	watcher *fsnotify.Watcher
}

var watcherList *haxmap.Map[string, watcherObj] = haxmap.New[string, watcherObj]()

// WatchDir watches the files in a directory and its subdirectories for changes
func WatchDir(root string, cb *Watcher) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

	rand := root + "." + string(RandBytes(64))
	if  _, ok := watcherList.Get(rand); ok {
		loops := 10000
		for loops > 0 {
			loops--
			if  _, ok := watcherList.Get(rand); !ok {
				break
			}
			rand = root + "." + string(RandBytes(64))
		}
	}
	if  _, ok := watcherList.Get(rand); !ok {
		watcherList.Set(rand, watcherObj{root, watcher})
	}

	done := make(chan bool)
	go func() {
		defer close(done)
		for {
			if event, ok := <-watcher.Events; ok {
				filePath := event.Name

				stat, err := os.Stat(filePath)
				if err != nil {
					if cb.Remove == nil || cb.Remove(filePath, event.Op.String()){
						watcher.Remove(filePath)
					}
				}else if stat.IsDir() {
					if cb.DirAdd == nil || cb.DirAdd(filePath, event.Op.String()){
						watcher.Add(filePath)
					}
				}else{
					if cb.FileChange != nil {
						cb.FileChange(filePath, event.Op.String())
					}
				}

				if cb.Any != nil {
					cb.Any(filePath, event.Op.String())
				}
			}
		}
	}()

	err = watcher.Add(root)
	if err != nil {
		return
	}

	watchDirSub(watcher, root)

	<-done
}

func watchDirSub(watcher *fsnotify.Watcher, dir string){
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if file.IsDir() {
			if path, err := JoinPath(dir, file.Name()); err == nil {
				watcher.Add(path)
				watchDirSub(watcher, path)
			}
		}
	}
}

// CloseWatchers will close all the watchers with the given root if they were created by goutil
//
// @root pass a file path for a specific watcher or "*" for all watchers that exist
//
// note: this method may include other modules that are using goutil as a dependency
func CloseWatchers(root string) {
	watcherList.ForEach(func(id string, cache watcherObj) bool {
		if root == "*" || cache.root == root {
			cache.watcher.Close()
		}
		return true
	})
}


// InstallLinuxPkg attempts to install a linux package
//
// this method will also resolve the sudo command and ask for a user password if needed
//
// this method will not attempt to run an install, if it finds the package is already installed
func InstallLinuxPkg(pkg []string, man ...string){
	if !HasLinuxPkg(pkg) {
		var pkgMan string
		if len(man) != 0 {
			pkgMan = man[0]
		}else{
			pkgMan = GetLinuxInstaller([]string{`apt-get`, `apt`, `yum`})
		}

		cmd := exec.Command(`sudo`, append([]string{pkgMan, `install`, `-y`}, pkg...)...)

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return
		}

		go (func() {
			out := bufio.NewReader(stdout)
			for {
				s, err := out.ReadString('\n')
				if err == nil {
					fmt.Println(s)
				}
			}
		})()

		stderr, err := cmd.StderrPipe()
		if err != nil {
			return
		}

		go (func() {
			out := bufio.NewReader(stderr)
			for {
				s, err := out.ReadString('\n')
				if err == nil {
					fmt.Println(s)
				}
			}
		})()

		cmd.Run()
	}
}

// HasLinuxPkg attempt to check if a linux package is installed
func HasLinuxPkg(pkg []string) bool {
	for _, name := range pkg {
		hasPackage := false
		cmd := exec.Command(`dpkg`, `-s`, name)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			return true
		}
		go (func() {
			out := bufio.NewReader(stdout)
			for {
				_, err := out.ReadString('\n')
				if err == nil {
					hasPackage = true
				}
			}
		})()
		for i := 0; i < 3; i++ {
			cmd.Run()
			if hasPackage {
				break
			}
		}
		if !hasPackage {
			return false
		}
	}

	return true
}

// GetLinuxInstaller attempt to find out what package manager a linux distro is using or has available
func GetLinuxInstaller(man []string) string {
	hasInstaller := ""

	for _, m := range man {
		cmd := exec.Command(`dpkg`, `-s`, m)
		stdout, err := cmd.StdoutPipe()
		if err != nil {
			continue
		}
		go (func() {
			out := bufio.NewReader(stdout)
			for {
				_, err := out.Peek(1)
				if err == nil {
					hasInstaller = m
				}
			}
		})()

		for i := 0; i < 3; i++ {
			cmd.Run()
			if hasInstaller != "" {
				break
			}
		}

		if hasInstaller != "" {
			break
		}
	}

	return hasInstaller
}
