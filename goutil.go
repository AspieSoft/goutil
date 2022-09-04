package goutil

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"github.com/AspieSoft/go-regex"
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

// Joins multiple file types with safety from backtracking
func JoinPath[T interface{string|[]byte}](path ...T) (string, error) {
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

// Returns true if an array contains a value
func Contains[T any](search []T, value T) bool {
	val := ToString(value)
	for _, v := range search {
		if ToString(v) == val {
			return true
		}
	}
	return false
}

// Returns true if a map contains a value
func ContainsMap[T Hashable, J any](search map[T]J, value J) bool {
	val := ToString(value)
	for _, v := range search {
		if ToString(v) == val {
			return true
		}
	}
	return false
}

// Returns true if a map contains a key
func ContainsMapKey[T Hashable, J any](search map[T]J, key T) bool {
	for i := range search {
		if i == key {
			return true
		}
	}
	return false
}

// Converts multiple types to a string
// accepts: string, byteArray, byte, int32, int, int64, float64, float32
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

// Converts multiple types to an int
// accepts: int, int32, int64, float64, float32, string, byteArray, byte
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

func IsZeroOfUnderlyingType(x interface{}) bool {
	// return x == nil || x == reflect.Zero(reflect.TypeOf(x)).Interface()
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}

// Converts bytes to megabytes
func FormatMemoryUsage(b uint64) float64 {
	return math.Round(float64(b) / 1024 / 1024 * 100) / 100
}

// Replaces HTML characters with html entities
// Also prevents &amp;amp; from results
func EscapeHTML(html []byte) []byte {
	html = regex.RepFunc(html, `[<>&]`, func(data func(int) []byte) []byte {
		if bytes.Equal(data(0), []byte("<")) {
			return []byte("&lt;")
		} else if bytes.Equal(data(0), []byte(">")) {
			return []byte("&gt;")
		}
		return []byte("&amp;")
	})
	return regex.RepStr(html, `&amp;(amp;)*`, []byte("&amp;"))
}

// Escapes quotes and backslashes for use within HTML quotes 
func EscapeHTMLArgs(html []byte) []byte {
	return regex.RepFunc(html, `[\\"'\']`, func(data func(int) []byte) []byte {
		return append([]byte("\\"), data(0)...)
	})
}

// Converts a map or array to a JSON string
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

// Converts a json string into a map of strings
func ParseJson(b []byte) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	err := json.Unmarshal(b, &res)
	if err != nil {
		return map[string]interface{}{}, err
	}
	return res, nil
}

// Useful for decoding a JSON output from the body of an http request
// goutil.DecodeJSON(r.Body)
func DecodeJSON(data io.Reader) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := json.NewDecoder(data).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Gzip compression for a string
func Compress[T interface{string|[]byte}](msg T) (string, error) {
	var b bytes.Buffer
	gz := gzip.NewWriter(&b)
	if _, err := gz.Write([]byte(msg)); err != nil {
		return "", err
	}
	if err := gz.Flush(); err != nil {
		return "", err
	}
	if err := gz.Close(); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b.Bytes()), nil
}

// Gzip decompression for a string
func Decompress[T interface{string|[]byte}](str T) (string, error) {
	data, err := base64.StdEncoding.DecodeString(string(str))
	if err != nil {
		return "", err
	}
	rdata := bytes.NewReader(data)
	r, err := gzip.NewReader(rdata)
	if err != nil {
		return "", err
	}
	s, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	return string(s), nil
}


// AES-CFB Encryption
func Encrypt[T interface{string|[]byte}, J interface{string|[]byte}](text T, key J) (string, error) {
	plaintext := []byte(text)

	keyHash := sha256.Sum256([]byte(key))

	block, err := aes.NewCipher(keyHash[:])
	if err != nil {
		return "", err
	}

	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", nil
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// AES-CFB Decryption
func Decrypt[T interface{string|[]byte}, J interface{string|[]byte}](text T, key J) ([]byte, error) {
	keyHash := sha256.Sum256([]byte(key))

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


/* func EncryptFast(msg []byte, key string) (string, error) {
	bKey := []byte(key)
	var eKey []byte = nil
	if len(bKey) > 32 {
		eKey = bKey[32:]
		bKey = bKey[:32]
	}

	block, err := aes.NewCipher(bKey)
	if err != nil {
		return "", err
	}
	b := base64.StdEncoding.EncodeToString(msg)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	if eKey != nil {
		return encryptByteFast(ciphertext, eKey)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func encryptByteFast(text []byte, key []byte) (string, error) {
	var eKey []byte = nil
	if len(key) > 32 {
		eKey = key[32:]
		key = key[:32]
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	b := base64.StdEncoding.EncodeToString(text)
	ciphertext := make([]byte, aes.BlockSize+len(b))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], []byte(b))

	if eKey != nil {
		return encryptByteFast(ciphertext, eKey)
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func DecryptFast(ciphertext string, key string) ([]byte, error) {
	bKey := []byte(key)
	var eKey []byte = nil
	if len(bKey) > 32 {
		l := len(bKey) - 32
		eKey = bKey[:l]
		bKey = bKey[l:]
	}

	block, err := aes.NewCipher(bKey)
	if err != nil {
		return nil, err
	}

	text, err := base64.StdEncoding.DecodeString(ciphertext)
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
		return decryptByteFast(data, eKey)
	}

	return data, nil
}

func decryptByteFast(text []byte, key []byte) ([]byte, error) {
	var eKey []byte = nil
	if len(key) > 32 {
		l := len(key) - 32
		eKey = key[:l]
		key = key[l:]
	}

	block, err := aes.NewCipher(key)
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
		return decryptByteFast(data, eKey)
	}

	return data, nil
} */


// Sanitizes a string to valid UTF-8
func CleanStr(str string) string {
	//todo: sanitize inputs
	str = strings.ToValidUTF8(str, "")
	return str
}

// Runs CleanStr on an array
// CleanStr: Sanitizes a string to valid UTF-8
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

// Runs CleanStr on a map[string]
// CleanStr: Sanitizes a string to valid UTF-8
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

// Runs CleanStr on a complex json value
// CleanStr: Sanitizes a string to valid UTF-8
func CleanJSON(val interface{}) interface{} {
	t := reflect.TypeOf(val)
	if t == VarType["string"] {
		return CleanStr(val.(string))
	}else if t == VarType["int"] || t == VarType["float64"] || t == VarType["float32"] || t == VarType["bool"] {
		return val
	}else if t == VarType["byteArray"] {
		return CleanStr(string(val.([]byte)))
	}else if t == VarType["byte"] {
		return CleanStr(string(val.(byte)))
	}else if t == VarType["int32"] {
		return CleanStr(string(val.(int32)))
	}else if t == VarType["array"] {
		return CleanArray(val.([]interface{}))
	}else if t == VarType["map"] {
		return CleanMap(val.(map[string]interface{}))
	}
	return nil
}


// Checks if the parent (or sub parent) directory of a file contains a specific file or folder
// root: the highest grandparent to check before quitting
// start: the lowest level to start searching from (if a directory is passed, it will not be included in your search)
// search: what file you want to search fro
func GetFileFromParent[T interface{string|[]byte}](root T, start T, search string) (string, bool) {
	rootB := []byte(root)
	dir := regex.RepStr([]byte(start), `[\\/][^\\/]*$`, []byte{})
	if len(dir) == 0 || bytes.Equal(dir, rootB) || !bytes.HasPrefix(dir, rootB) {
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

	return GetFileFromParent(rootB, dir, search)
}


// A watcher instance for the WatchDir function
type Watcher struct {
	FileChange func(path string, op string)
	DirAdd func(path string, op string) (addWatcher bool)
	Remove func(path string, op string) (removeWatcher bool)
	Any func(path string, op string)
}

// Watch the files in a directory and its subdirectories for changes
func WatchDir(root string, cb *Watcher) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()

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
