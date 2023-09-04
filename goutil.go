package goutil

import (
	"bytes"
	"errors"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/AspieSoft/go-regex-re2/v2"
)

type Hashable interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float32 | float64 | string | complex64 | complex128
}

// Contains returns true if an array contains a value
func Contains[T any](search []T, value T) bool {
	val := toString[string](value)
	for _, v := range search {
		if toString[string](v) == val {
			return true
		}
	}
	return false
}

// IndexOf returns the index of a value in an array
//
// returns -1 and an error if the value is not found
func IndexOf[T any](search []T, value T) (int, error) {
	val := toString[string](value)
	for i, v := range search {
		if toString[string](v) == val {
			return i, nil
		}
	}
	return -1, errors.New("array does not contain value: " + toString[string](value))
}

// ContainsMap returns true if a map contains a value
func ContainsMap[T Hashable, J any](search map[T]J, value J) bool {
	val := toString[string](value)
	for _, v := range search {
		if toString[string](v) == val {
			return true
		}
	}
	return false
}

// IndexOfMap returns the index of a value in a map
//
// returns an error if the value is not found
func IndexOfMap[T Hashable, J any](search map[T]J, value J) (T, error) {
	val := toString[string](value)
	for i, v := range search {
		if toString[string](v) == val {
			return i, nil
		}
	}
	var blk T
	return blk, errors.New("map does not contain value: " + toString[string](value))
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

// ArrayEqual returns true if 2 arrays are equal and of the same length (even if they are in a different order)
func ArrayEqual[T any](arr1 []T, arr2 []T, ignoreLength ...bool) bool {
	if !(len(ignoreLength) != 0 && ignoreLength[0] == true) && len(arr1) != len(arr2) {
		return false
	}

	for _, val1 := range arr1 {
		hasMatch := false
		for _, val2 := range arr2 {
			if TypeEqual(val1, val2) {
				hasMatch = true
				break
			}
		}
		if !hasMatch {
			return false
		}
	}
	return true
}

// MapEqual returns true if 2 maps are equal and of the same length (even if they are in a different order)
func MapEqual[T Hashable, J any](map1 map[T]J, map2 map[T]J, ignoreLength ...bool) bool {
	if !(len(ignoreLength) != 0 && ignoreLength[0] == true) && len(map1) != len(map2) {
		return false
	}

	for key1, val1 := range map1 {
		if !TypeEqual(val1, map2[key1]) {
			return false
		}
	}
	return true
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

// FormatMemoryUsage converts bytes to megabytes
func FormatMemoryUsage(b uint64) float64 {
	return math.Round(float64(b) / 1024 / 1024 * 100) / 100
}


var regIsAlphaNumeric *regex.Regexp = regex.Comp(`^[A-Za-z0-9]+$`)

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
