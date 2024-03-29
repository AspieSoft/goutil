package goutil

import (
	"bytes"
	"reflect"
	"strconv"
)

// Type Convertions
type typeConv struct {}
var Conv *typeConv = &typeConv{}

type NullType[T any] struct {
	Null T
}

type ToInterface struct {
	Val interface{}
}

var VarType map[string]reflect.Type

func init(){
	VarType = map[string]reflect.Type{}

	VarType["[]interface{}"] = reflect.TypeOf([]interface{}{})
	VarType["array"] = VarType["[]interface{}"]
	VarType["[][]byte"] = reflect.TypeOf([][]byte{})
	VarType["map[string]interface{}"] = reflect.TypeOf(map[string]interface{}{})
	VarType["map"] = VarType["map[string]interface{}"]

	VarType["int"] = reflect.TypeOf(int(0))
	VarType["int64"] = reflect.TypeOf(int64(0))
	VarType["float64"] = reflect.TypeOf(float64(0))
	VarType["float32"] = reflect.TypeOf(float32(0))

	VarType["string"] = reflect.TypeOf("")
	VarType["[]byte"] = reflect.TypeOf([]byte{})
	VarType["byteArray"] = VarType["[]byte"]
	VarType["byte"] = reflect.TypeOf([]byte{0}[0])

	// ' ' returned int32 instead of byte
	VarType["int32"] = reflect.TypeOf(int32(0))
	VarType["rune"] = reflect.TypeOf(rune(0))

	VarType["func"] = reflect.TypeOf(func(){})

	VarType["bool"] = reflect.TypeOf(false)

	VarType["int8"] = reflect.TypeOf(int8(0))
	VarType["int16"] = reflect.TypeOf(int16(0))
	
	VarType["uint"] = reflect.TypeOf(uint(0))
	VarType["uint8"] = reflect.TypeOf(uint8(0))
	VarType["uint16"] = reflect.TypeOf(uint16(0))
	VarType["uint32"] = reflect.TypeOf(uint32(0))
	VarType["uint64"] = reflect.TypeOf(uint64(0))
	VarType["uintptr"] = reflect.TypeOf(uintptr(0))

	VarType["complex128"] = reflect.TypeOf(complex128(0))
	VarType["complex64"] = reflect.TypeOf(complex64(0))

	VarType["map[byte]interface{}"] = reflect.TypeOf(map[byte]interface{}{})
	VarType["map[rune]interface{}"] = reflect.TypeOf(map[byte]interface{}{})
	VarType["map[int]interface{}"] = reflect.TypeOf(map[int]interface{}{})
	VarType["map[int64]interface{}"] = reflect.TypeOf(map[int64]interface{}{})
	VarType["map[int32]interface{}"] = reflect.TypeOf(map[int32]interface{}{})
	VarType["map[float64]interface{}"] = reflect.TypeOf(map[float64]interface{}{})
	VarType["map[float32]interface{}"] = reflect.TypeOf(map[float32]interface{}{})

	VarType["map[int8]interface{}"] = reflect.TypeOf(map[int8]interface{}{})
	VarType["map[int16]interface{}"] = reflect.TypeOf(map[int16]interface{}{})

	VarType["map[uint]interface{}"] = reflect.TypeOf(map[uint]interface{}{})
	VarType["map[uint8]interface{}"] = reflect.TypeOf(map[uint8]interface{}{})
	VarType["map[uint16]interface{}"] = reflect.TypeOf(map[uint16]interface{}{})
	VarType["map[uint32]interface{}"] = reflect.TypeOf(map[uint32]interface{}{})
	VarType["map[uint64]interface{}"] = reflect.TypeOf(map[uint64]interface{}{})
	VarType["map[uintptr]interface{}"] = reflect.TypeOf(map[uintptr]interface{}{})

	VarType["map[complex128]interface{}"] = reflect.TypeOf(map[complex128]interface{}{})
	VarType["map[complex64]interface{}"] = reflect.TypeOf(map[complex64]interface{}{})

	VarType["[]string"] = reflect.TypeOf([]string{})
	VarType["[]bool"] = reflect.TypeOf([]bool{})
	VarType["[]rune"] = reflect.TypeOf([]bool{})
	VarType["[]int"] = reflect.TypeOf([]int{})
	VarType["[]int64"] = reflect.TypeOf([]int64{})
	VarType["[]int32"] = reflect.TypeOf([]int32{})
	VarType["[]float64"] = reflect.TypeOf([]float64{})
	VarType["[]float32"] = reflect.TypeOf([]float32{})

	VarType["[]int8"] = reflect.TypeOf([]int8{})
	VarType["[]int16"] = reflect.TypeOf([]int16{})

	VarType["[]uint"] = reflect.TypeOf([]uint{})
	VarType["[]uint8"] = reflect.TypeOf([]uint8{})
	VarType["[]uint16"] = reflect.TypeOf([]uint16{})
	VarType["[]uint32"] = reflect.TypeOf([]uint32{})
	VarType["[]uint64"] = reflect.TypeOf([]uint64{})
	VarType["[]uintptr"] = reflect.TypeOf([]uintptr{})

	VarType["[]complex128"] = reflect.TypeOf([]complex128{})
	VarType["[]complex64"] = reflect.TypeOf([]complex64{})
}

// toString converts multiple types to a string|[]byte
//
// accepts: string, []byte, byte, int (and variants), [][]byte, []interface{}
func toString[T interface{string | []byte}](val interface{}) T {
	switch reflect.TypeOf(val) {
		case VarType["string"]:
			return T(val.(string))
		case VarType["[]byte"]:
			return T(val.([]byte))
		case VarType["byte"]:
			return T([]byte{val.(byte)})
		case VarType["int"]:
			return T(strconv.Itoa(val.(int)))
		case VarType["int64"]:
			return T(strconv.Itoa(int(val.(int64))))
		case VarType["int32"]:
			return T([]byte{byte(val.(int32))})
		case VarType["int16"]:
			return T([]byte{byte(val.(int16))})
		case VarType["int8"]:
			return T([]byte{byte(val.(int8))})
		case VarType["uintptr"]:
			return T(strconv.FormatUint(uint64(val.(uintptr)), 10))
		case VarType["uint"]:
			return T(strconv.FormatUint(uint64(val.(uint)), 10))
		case VarType["uint64"]:
			return T(strconv.FormatUint(val.(uint64), 10))
		case VarType["uint32"]:
			return T(strconv.FormatUint(uint64(val.(uint32)), 10))
		case VarType["uint16"]:
			return T(strconv.FormatUint(uint64(val.(uint16)), 10))
		case VarType["uint8"]:
			return T(strconv.FormatUint(uint64(val.(uint8)), 10))
		case VarType["float64"]:
			return T(strconv.FormatFloat(val.(float64), 'f', -1, 64))
		case VarType["float32"]:
			return T(strconv.FormatFloat(float64(val.(float32)), 'f', -1, 32))
		case VarType["rune"]:
			return T([]byte{byte(val.(rune))})
		case VarType["[]interface{}"]:
			b := make([]byte, len(val.([]interface{})))
			for i, v := range val.([]interface{}) {
				b[i] = byte(toNumber[int32](v))
			}
			return T(b)
		case VarType["[]int"]:
			b := make([]byte, len(val.([]int)))
			for i, v := range val.([]int) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]int64"]:
			b := make([]byte, len(val.([]int64)))
			for i, v := range val.([]int64) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]int32"]:
			b := make([]byte, len(val.([]int32)))
			for i, v := range val.([]int32) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]int16"]:
			b := make([]byte, len(val.([]int16)))
			for i, v := range val.([]int16) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]int8"]:
			b := make([]byte, len(val.([]int8)))
			for i, v := range val.([]int8) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]uint"]:
			b := make([]byte, len(val.([]uint)))
			for i, v := range val.([]uint) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]uint8"]:
			b := make([]byte, len(val.([]uint8)))
			for i, v := range val.([]uint8) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]uint16"]:
			b := make([]byte, len(val.([]uint16)))
			for i, v := range val.([]uint16) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]uint32"]:
			b := make([]byte, len(val.([]uint32)))
			for i, v := range val.([]uint32) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]uint64"]:
			b := make([]byte, len(val.([]uint64)))
			for i, v := range val.([]uint64) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]uintptr"]:
			b := make([]byte, len(val.([]uintptr)))
			for i, v := range val.([]uintptr) {
				b[i] = byte(v)
			}
			return T(b)
		case VarType["[]string"]:
			b := []byte{}
			for _, v := range val.([]string) {
				b = append(b, []byte(v)...)
			}
			return T(b)
		case VarType["[][]byte"]:
			b := []byte{}
			for _, v := range val.([][]byte) {
				b = append(b, v...)
			}
			return T(b)
		case VarType["[]rune"]:
			b := []byte{}
			for _, v := range val.([]rune) {
				b = append(b, byte(v))
			}
			return T(b)
		default:
			return T("")
	}
}

// ToString converts multiple types to a string
//
// accepts: string, []byte, byte, int (and variants), [][]byte, []interface{}
func (conv *typeConv) ToString(val interface{}) string {
	return toString[string](val)
}

// ToBytes converts multiple types to a []byte
//
// accepts: string, []byte, byte, int (and variants), [][]byte, []interface{}
func (conv *typeConv) ToBytes(val interface{}) []byte {
	return toString[[]byte](val)
}

// toNumber converts multiple types to a number
//
// accepts: int (and variants), string, []byte, byte, bool
func toNumber[T interface{int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | uintptr | float64 | float32}](val interface{}) T {
	switch reflect.TypeOf(val) {
		case VarType["int"]:
			return T(val.(int))
		case VarType["int32"]:
			return T(val.(int32))
		case VarType["int64"]:
			return T(val.(int64))
		case VarType["float64"]:
			return T(val.(float64))
		case VarType["float32"]:
			return T(val.(float32))
		case VarType["string"]:
			var varT interface{} = T(0)
			if _, ok := varT.(float64); ok {
				if f, err := strconv.ParseFloat(val.(string), 64); err == nil {
					return T(f)
				}
			}else if _, ok := varT.(float32); ok {
				if f, err := strconv.ParseFloat(val.(string), 32); err == nil {
					return T(f)
				}
			}else if i, err := strconv.Atoi(val.(string)); err == nil {
				return T(i)
			}
			return 0
		case VarType["[]byte"]:
			if i, err := strconv.Atoi(string(val.([]byte))); err == nil {
				return T(i)
			}
			return 0
		case VarType["byte"]:
			if i, err := strconv.Atoi(string(val.(byte))); err == nil {
				return T(i)
			}
			return 0
		case VarType["bool"]:
			if val.(bool) == true {
				return 1
			}
			return 0
		case VarType["int8"]:
			return T(val.(int8))
		case VarType["int16"]:
			return T(val.(int16))
		case VarType["uint"]:
			return T(val.(uint))
		case VarType["uint8"]:
			return T(val.(uint8))
		case VarType["uint16"]:
			return T(val.(uint16))
		case VarType["uint32"]:
			return T(val.(uint32))
		case VarType["uint64"]:
			return T(val.(uint64))
		case VarType["uintptr"]:
			return T(val.(uintptr))
		case VarType["rune"]:
			if i, err := strconv.Atoi(string(val.(rune))); err == nil {
				return T(i)
			}
			return 0
		default:
			return 0
	}
}

// ToInt converts multiple types to an int
//
// accepts: int (and variants), string, []byte, byte, bool
func (conv *typeConv) ToInt(val interface{}) int {
	return toNumber[int](val)
}

// ToUint converts multiple types to a uint
//
// accepts: int (and variants), string, []byte, byte, bool
func (conv *typeConv) ToUint(val interface{}) uint {
	return toNumber[uint](val)
}

// ToUintptr converts multiple types to a uintptr
//
// accepts: int (and variants), string, []byte, byte, bool
func (conv *typeConv) ToUintptr(val interface{}) uintptr {
	return toNumber[uintptr](val)
}

// ToFloat converts multiple types to a float64
//
// accepts: int (and variants), string, []byte, byte, bool
func (conv *typeConv) ToFloat(val interface{}) float64 {
	return toNumber[float64](val)
}

// ToMap converts multiple types to a map[string]interface{}
func (conv *typeConv) ToMap(val interface{}) map[string]interface{} {
	switch reflect.TypeOf(val) {
		case VarType["map[string]interface{}"]:
			return val.(map[string]interface{})
		case VarType["map[byte]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[byte]interface{})))
			for k, v := range val.(map[byte]interface{}) {
				m[string(k)] = v
			}
			return m
		case VarType["map[int]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[int]interface{})))
			for k, v := range val.(map[int]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[int64]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[int64]interface{})))
			for k, v := range val.(map[int64]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[int32]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[int32]interface{})))
			for k, v := range val.(map[int32]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[int16]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[int16]interface{})))
			for k, v := range val.(map[int16]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[int8]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[int8]interface{})))
			for k, v := range val.(map[int8]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[uintptr]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[uintptr]interface{})))
			for k, v := range val.(map[uintptr]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[uint]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[uint]interface{})))
			for k, v := range val.(map[uint]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[uint64]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[uint64]interface{})))
			for k, v := range val.(map[uint64]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[uint32]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[uint32]interface{})))
			for k, v := range val.(map[uint32]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[uint16]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[uint16]interface{})))
			for k, v := range val.(map[uint16]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[uint8]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[uint8]interface{})))
			for k, v := range val.(map[uint8]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[float64]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[float64]interface{})))
			for k, v := range val.(map[float64]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[float32]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[float32]interface{})))
			for k, v := range val.(map[float32]interface{}) {
				m[toString[string](k)] = v
			}
			return m
		case VarType["map[rune]interface{}"]:
			m := make(map[string]interface{}, len(val.(map[rune]interface{})))
			for k, v := range val.(map[rune]interface{}) {
				m[string(k)] = v
			}
			return m
		case VarType["[]interface{}"]:
			m := make(map[string]interface{}, len(val.([]interface{})))
			for k, v := range val.([]interface{}) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]byte"]:
			m := make(map[string]interface{}, len(val.([]byte)))
			for k, v := range val.([]byte) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]string"]:
			m := make(map[string]interface{}, len(val.([]string)))
			for k, v := range val.([]string) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]bool"]:
			m := make(map[string]interface{}, len(val.([]bool)))
			for k, v := range val.([]bool) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]int"]:
			m := make(map[string]interface{}, len(val.([]int)))
			for k, v := range val.([]int) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]int64"]:
			m := make(map[string]interface{}, len(val.([]int64)))
			for k, v := range val.([]int64) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]int32"]:
			m := make(map[string]interface{}, len(val.([]int32)))
			for k, v := range val.([]int32) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]int16"]:
			m := make(map[string]interface{}, len(val.([]int16)))
			for k, v := range val.([]int16) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]int8"]:
			m := make(map[string]interface{}, len(val.([]int8)))
			for k, v := range val.([]int8) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]uintptr"]:
			m := make(map[string]interface{}, len(val.([]uintptr)))
			for k, v := range val.([]uintptr) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]uint"]:
			m := make(map[string]interface{}, len(val.([]uint)))
			for k, v := range val.([]uint) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]uint64"]:
			m := make(map[string]interface{}, len(val.([]uint64)))
			for k, v := range val.([]uint64) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]uint32"]:
			m := make(map[string]interface{}, len(val.([]uint32)))
			for k, v := range val.([]uint32) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]uint16"]:
			m := make(map[string]interface{}, len(val.([]uint16)))
			for k, v := range val.([]uint16) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]uint8"]:
			m := make(map[string]interface{}, len(val.([]uint8)))
			for k, v := range val.([]uint8) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]float64"]:
			m := make(map[string]interface{}, len(val.([]float64)))
			for k, v := range val.([]float64) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]float32"]:
			m := make(map[string]interface{}, len(val.([]float32)))
			for k, v := range val.([]float32) {
				m[strconv.Itoa(k)] = v
			}
			return m
		case VarType["[]rune"]:
			m := make(map[string]interface{}, len(val.([]rune)))
			for k, v := range val.([]rune) {
				m[strconv.Itoa(k)] = v
			}
			return m
		default:
			return map[string]interface{}{}
	}
}

// ToArray converts multiple types to an []interface{}
func (conv *typeConv) ToArray(val interface{}) []interface{} {
	switch reflect.TypeOf(val) {
		case VarType["[]interface{}"]:
			return val.([]interface{})
		case VarType["[]byte"]:
			a := make([]interface{}, len(val.([]byte)))
			for i, v := range val.([]byte) {
				a[i] = v
			}
			return a
		case VarType["[][]byte"]:
			a := make([]interface{}, len(val.([][]byte)))
			for i, v := range val.([][]byte) {
				a[i] = v
			}
			return a
		case VarType["[]string"]:
			a := make([]interface{}, len(val.([]string)))
			for i, v := range val.([]string) {
				a[i] = v
			}
			return a
		case VarType["[]bool"]:
			a := make([]interface{}, len(val.([]bool)))
			for i, v := range val.([]bool) {
				a[i] = v
			}
			return a
		case VarType["[]int"]:
			a := make([]interface{}, len(val.([]int)))
			for i, v := range val.([]int) {
				a[i] = v
			}
			return a
		case VarType["[]int64"]:
			a := make([]interface{}, len(val.([]int64)))
			for i, v := range val.([]int64) {
				a[i] = v
			}
			return a
		case VarType["[]int32"]:
			a := make([]interface{}, len(val.([]int32)))
			for i, v := range val.([]int32) {
				a[i] = v
			}
			return a
		case VarType["[]int16"]:
			a := make([]interface{}, len(val.([]int16)))
			for i, v := range val.([]int16) {
				a[i] = v
			}
			return a
		case VarType["[]int8"]:
			a := make([]interface{}, len(val.([]int8)))
			for i, v := range val.([]int8) {
				a[i] = v
			}
			return a
		case VarType["[]uint"]:
			a := make([]interface{}, len(val.([]uint)))
			for i, v := range val.([]uint) {
				a[i] = v
			}
			return a
		case VarType["[]uint64"]:
			a := make([]interface{}, len(val.([]uint64)))
			for i, v := range val.([]uint64) {
				a[i] = v
			}
			return a
		case VarType["[]uint32"]:
			a := make([]interface{}, len(val.([]uint32)))
			for i, v := range val.([]uint32) {
				a[i] = v
			}
			return a
		case VarType["[]uint16"]:
			a := make([]interface{}, len(val.([]uint16)))
			for i, v := range val.([]uint16) {
				a[i] = v
			}
			return a
		case VarType["[]uint8"]:
			a := make([]interface{}, len(val.([]uint8)))
			for i, v := range val.([]uint8) {
				a[i] = v
			}
			return a
		case VarType["[]float64"]:
			a := make([]interface{}, len(val.([]float64)))
			for i, v := range val.([]float64) {
				a[i] = v
			}
			return a
		case VarType["[]float32"]:
			a := make([]interface{}, len(val.([]float32)))
			for i, v := range val.([]float32) {
				a[i] = v
			}
			return a
		case VarType["[]rune"]:
			a := make([]interface{}, len(val.([]rune)))
			for i, v := range val.([]rune) {
				a[i] = v
			}
			return a
		case VarType["map[string]interface{}"]:
			a := make([]interface{}, len(val.(map[string]interface{})))
			for i, v := range val.(map[string]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[byte]interface{}"]:
			a := make([]interface{}, len(val.(map[byte]interface{})))
			for i, v := range val.(map[byte]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[int]interface{}"]:
			a := make([]interface{}, len(val.(map[int]interface{})))
			for i, v := range val.(map[int]interface{}) {
				a[i] = v
			}
			return a
		case VarType["map[int64]interface{}"]:
			a := make([]interface{}, len(val.(map[int64]interface{})))
			for i, v := range val.(map[int64]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[int32]interface{}"]:
			a := make([]interface{}, len(val.(map[int32]interface{})))
			for i, v := range val.(map[int32]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[int16]interface{}"]:
			a := make([]interface{}, len(val.(map[int16]interface{})))
			for i, v := range val.(map[int16]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[int8]interface{}"]:
			a := make([]interface{}, len(val.(map[int8]interface{})))
			for i, v := range val.(map[int8]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[uint]interface{}"]:
			a := make([]interface{}, len(val.(map[uint]interface{})))
			for i, v := range val.(map[uint]interface{}) {
				a[i] = v
			}
			return a
		case VarType["map[uint64]interface{}"]:
			a := make([]interface{}, len(val.(map[uint64]interface{})))
			for i, v := range val.(map[uint64]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[uint32]interface{}"]:
			a := make([]interface{}, len(val.(map[uint32]interface{})))
			for i, v := range val.(map[uint32]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[uint16]interface{}"]:
			a := make([]interface{}, len(val.(map[uint16]interface{})))
			for i, v := range val.(map[uint16]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[uint8]interface{}"]:
			a := make([]interface{}, len(val.(map[uint8]interface{})))
			for i, v := range val.(map[uint8]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[float64]interface{}"]:
			a := make([]interface{}, len(val.(map[float64]interface{})))
			for i, v := range val.(map[float64]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[float32]interface{}"]:
			a := make([]interface{}, len(val.(map[float32]interface{})))
			for i, v := range val.(map[float32]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		case VarType["map[rune]interface{}"]:
			a := make([]interface{}, len(val.(map[rune]interface{})))
			for i, v := range val.(map[rune]interface{}) {
				a[toNumber[int](i)] = v
			}
			return a
		default:
			return []interface{}{}
	}
}

// SupportedType is an interface containing the types which are supported by the ToType method
type SupportedType interface {
	string | []byte | byte | bool |
	int | int64 | int32 | int16 | int8 |
	uint | uint64 | uint32 | uint16 | /* uint8 | */ uintptr |
	float64 | float32 |
	[]interface{} | []string | [][]byte | []bool |
	[]int | []int64 | []int32 | []int16 | []int8 |
	[]uint | []uint64 | []uint32 | []uint16 | /* []uint8 | */ []uintptr |
	[]float64 | []float32 |
	map[string]interface{} | map[byte]interface{} |
	map[int]interface{} | map[int64]interface{} | map[int32]interface{} | map[int16]interface{} | map[int8]interface{} |
	map[uint]interface{} | map[uint64]interface{} | map[uint32]interface{} | map[uint16]interface{} | /* map[uint8]interface{} | */ map[uintptr]interface{} |
	map[float64]interface{} | map[float32]interface{}
}

// ToType attempts to converts an interface{} from the many possible types in golang, to a specific type of your choice
//
// if it fails to convert, it will return a nil/zero value for the appropriate type
func ToType[T SupportedType](val interface{}) T {
	// basic
	var varT interface{} = ""
	if _, ok := varT.(T); ok {
		return ToInterface{toString[string](val)}.Val.(T)
	}

	varT = []byte{}
	if _, ok := varT.(T); ok {
		return ToInterface{toString[[]byte](val)}.Val.(T)
	}

	varT = byte(0)
	if _, ok := varT.(T); ok {
		if b := toString[[]byte](val); len(b) != 0 {
			return ToInterface{b[0]}.Val.(T)
		}
		return ToInterface{byte(0)}.Val.(T)
	}

	varT = rune(0)
	if _, ok := varT.(T); ok {
		if b := toString[[]byte](val); len(b) != 0 {
			return ToInterface{rune(b[0])}.Val.(T)
		}
		return ToInterface{rune(0)}.Val.(T)
	}

	varT = false
	if _, ok := varT.(T); ok {
		return ToInterface{!IsZeroOfUnderlyingType(val)}.Val.(T)
	}

	// int
	varT = int(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[int](val)}.Val.(T)
	}

	varT = int64(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[int64](val)}.Val.(T)
	}

	varT = int32(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[int32](val)}.Val.(T)
	}

	varT = int16(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[int16](val)}.Val.(T)
	}

	varT = int8(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[int8](val)}.Val.(T)
	}

	// uint
	varT = uintptr(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uintptr](val)}.Val.(T)
	}

	varT = uint(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uint](val)}.Val.(T)
	}

	varT = uint64(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uint64](val)}.Val.(T)
	}

	varT = uint32(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uint32](val)}.Val.(T)
	}

	varT = uint16(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uint16](val)}.Val.(T)
	}

	varT = uint8(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[uint8](val)}.Val.(T)
	}

	// float
	varT = float64(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[float64](val)}.Val.(T)
	}

	varT = float32(0)
	if _, ok := varT.(T); ok {
		return ToInterface{toNumber[float32](val)}.Val.(T)
	}

	// array - basic
	varT = []interface{}{}
	if _, ok := varT.(T); ok {
		return ToInterface{Conv.ToArray(val)}.Val.(T)
	}

	varT = []string{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]string, len(r))
		for i, v := range r {
			a[i] = toString[string](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = [][]byte{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([][]byte, len(r))
		for i, v := range r {
			a[i] = toString[[]byte](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []rune{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]rune, len(r))
		for i, v := range r {
			if b := toString[[]byte](v); len(b) != 0 {
				a[i] = rune(b[0])
			}else{
				a[i] = rune(toNumber[int32](v))
			}
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []bool{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]bool, len(r))
		for i, v := range r {
			a[i] = !IsZeroOfUnderlyingType(v)
		}
		return ToInterface{a}.Val.(T)
	}

	// array - int
	varT = []int{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]int, len(r))
		for i, v := range r {
			a[i] = toNumber[int](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []int64{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]int64, len(r))
		for i, v := range r {
			a[i] = toNumber[int64](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []int32{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]int32, len(r))
		for i, v := range r {
			a[i] = toNumber[int32](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []int16{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]int16, len(r))
		for i, v := range r {
			a[i] = toNumber[int16](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []int8{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]int8, len(r))
		for i, v := range r {
			a[i] = toNumber[int8](v)
		}
		return ToInterface{a}.Val.(T)
	}

	// array - uint
	varT = []uintptr{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]uintptr, len(r))
		for i, v := range r {
			a[i] = toNumber[uintptr](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []uint{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]uint, len(r))
		for i, v := range r {
			a[i] = toNumber[uint](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []uint64{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]uint64, len(r))
		for i, v := range r {
			a[i] = toNumber[uint64](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []uint32{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]uint32, len(r))
		for i, v := range r {
			a[i] = toNumber[uint32](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []uint16{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]uint16, len(r))
		for i, v := range r {
			a[i] = toNumber[uint16](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []uint8{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]uint8, len(r))
		for i, v := range r {
			a[i] = toNumber[uint8](v)
		}
		return ToInterface{a}.Val.(T)
	}

	// array - float
	varT = []float64{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]float64, len(r))
		for i, v := range r {
			a[i] = toNumber[float64](v)
		}
		return ToInterface{a}.Val.(T)
	}

	varT = []float32{}
	if _, ok := varT.(T); ok {
		r := Conv.ToArray(val)
		a := make([]float32, len(r))
		for i, v := range r {
			a[i] = toNumber[float32](v)
		}
		return ToInterface{a}.Val.(T)
	}

	// map - basic
	varT = map[string]interface{}{}
	if _, ok := varT.(T); ok {
		return ToInterface{Conv.ToMap(val)}.Val.(T)
	}

	varT = map[byte]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[byte]interface{}, len(r))
		for i, v := range r {
			if b := toString[[]byte](i); len(b) != 0 {
				m[b[0]] = v
			}else{
				m[byte(toNumber[int32](i))] = v
			}
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[rune]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[rune]interface{}, len(r))
		for i, v := range r {
			if b := toString[[]byte](i); len(b) != 0 {
				m[rune(b[0])] = v
			}else{
				m[rune(byte(toNumber[int32](i)))] = v
			}
		}
		return ToInterface{m}.Val.(T)
	}

	// map - int
	varT = map[int]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[int]interface{}, len(r))
		for i, v := range r {
			m[toNumber[int](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[int64]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[int64]interface{}, len(r))
		for i, v := range r {
			m[toNumber[int64](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[int32]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[int32]interface{}, len(r))
		for i, v := range r {
			m[toNumber[int32](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[int16]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[int16]interface{}, len(r))
		for i, v := range r {
			m[toNumber[int16](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[int8]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[int8]interface{}, len(r))
		for i, v := range r {
			m[toNumber[int8](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	// map - uint
	varT = map[uintptr]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[uintptr]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uintptr](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[uint]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[uint]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uint](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[uint64]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[uint64]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uint64](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[uint32]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[uint32]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uint32](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[uint16]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[uint16]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uint16](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[uint8]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[uint8]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uint8](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	// map - float
	varT = map[float64]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[float64]interface{}, len(r))
		for i, v := range r {
			m[toNumber[float64](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	varT = map[float32]interface{}{}
	if _, ok := varT.(T); ok {
		r := Conv.ToMap(val)
		m := make(map[float32]interface{}, len(r))
		for i, v := range r {
			m[toNumber[float32](i)] = v
		}
		return ToInterface{m}.Val.(T)
	}

	return NullType[T]{}.Null
}

// ToVarTypeInterface attempts to convert an interface to match the unknown type of another interface
//
// this method is similar to the ToType method, but it simply returns nil if it cannot find the proper var type
func ToVarTypeInterface(val interface{}, ref interface{}) interface{} {
	refT := reflect.TypeOf(ref)

	// basic
	if refT == VarType["string"] {
		return ToInterface{toString[string](val)}.Val
	}

	if refT == VarType["[]byte"] {
		return ToInterface{toString[[]byte](val)}.Val
	}

	if refT == VarType["byte"] {
		if b := toString[[]byte](val); len(b) != 0 {
			return ToInterface{b[0]}.Val
		}
		return ToInterface{byte(0)}.Val
	}

	if refT == VarType["rune"] {
		if b := toString[[]byte](val); len(b) != 0 {
			return ToInterface{rune(b[0])}.Val
		}
		return ToInterface{rune(0)}.Val
	}

	if refT == VarType["bool"] {
		return ToInterface{!IsZeroOfUnderlyingType(val)}.Val
	}

	// int
	if refT == VarType["int"] {
		return ToInterface{toNumber[int](val)}.Val
	}

	if refT == VarType["int64"] {
		return ToInterface{toNumber[int64](val)}.Val
	}

	if refT == VarType["int32"] {
		return ToInterface{toNumber[int32](val)}.Val
	}

	if refT == VarType["int16"] {
		return ToInterface{toNumber[int16](val)}.Val
	}

	if refT == VarType["int8"] {
		return ToInterface{toNumber[int8](val)}.Val
	}

	// uint
	if refT == VarType["uintptr"] {
		return ToInterface{toNumber[uintptr](val)}.Val
	}

	if refT == VarType["uint"] {
		return ToInterface{toNumber[uint](val)}.Val
	}

	if refT == VarType["uint64"] {
		return ToInterface{toNumber[uint64](val)}.Val
	}

	if refT == VarType["uint32"] {
		return ToInterface{toNumber[uint32](val)}.Val
	}

	if refT == VarType["uint16"] {
		return ToInterface{toNumber[uint16](val)}.Val
	}

	if refT == VarType["uint8"] {
		return ToInterface{toNumber[uint8](val)}.Val
	}

	// float
	if refT == VarType["float64"] {
		return ToInterface{toNumber[float64](val)}.Val
	}

	if refT == VarType["float32"] {
		return ToInterface{toNumber[float32](val)}.Val
	}

	// array - basic
	if refT == VarType["interface{}"] {
		return ToInterface{Conv.ToArray(val)}.Val
	}

	if refT == VarType["[]string"] {
		r := Conv.ToArray(val)
		a := make([]string, len(r))
		for i, v := range r {
			a[i] = toString[string](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[][]byte"] {
		r := Conv.ToArray(val)
		a := make([][]byte, len(r))
		for i, v := range r {
			a[i] = toString[[]byte](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]rune"] {
		r := Conv.ToArray(val)
		a := make([]rune, len(r))
		for i, v := range r {
			if b := toString[[]byte](v); len(b) != 0 {
				a[i] = rune(b[0])
			}else{
				a[i] = rune(toNumber[int32](v))
			}
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]bool"] {
		r := Conv.ToArray(val)
		a := make([]bool, len(r))
		for i, v := range r {
			a[i] = !IsZeroOfUnderlyingType(v)
		}
		return ToInterface{a}.Val
	}

	// array - int
	if refT == VarType["[]int"] {
		r := Conv.ToArray(val)
		a := make([]int, len(r))
		for i, v := range r {
			a[i] = toNumber[int](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]int64"] {
		r := Conv.ToArray(val)
		a := make([]int64, len(r))
		for i, v := range r {
			a[i] = toNumber[int64](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]int32"] {
		r := Conv.ToArray(val)
		a := make([]int32, len(r))
		for i, v := range r {
			a[i] = toNumber[int32](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]int16"] {
		r := Conv.ToArray(val)
		a := make([]int16, len(r))
		for i, v := range r {
			a[i] = toNumber[int16](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]int8"] {
		r := Conv.ToArray(val)
		a := make([]int8, len(r))
		for i, v := range r {
			a[i] = toNumber[int8](v)
		}
		return ToInterface{a}.Val
	}

	// array - uint
	if refT == VarType["[]uintptr"] {
		r := Conv.ToArray(val)
		a := make([]uintptr, len(r))
		for i, v := range r {
			a[i] = toNumber[uintptr](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]uint"] {
		r := Conv.ToArray(val)
		a := make([]uint, len(r))
		for i, v := range r {
			a[i] = toNumber[uint](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]uint64"] {
		r := Conv.ToArray(val)
		a := make([]uint64, len(r))
		for i, v := range r {
			a[i] = toNumber[uint64](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]uint32"] {
		r := Conv.ToArray(val)
		a := make([]uint32, len(r))
		for i, v := range r {
			a[i] = toNumber[uint32](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]uint16"] {
		r := Conv.ToArray(val)
		a := make([]uint16, len(r))
		for i, v := range r {
			a[i] = toNumber[uint16](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]uint8"] {
		r := Conv.ToArray(val)
		a := make([]uint8, len(r))
		for i, v := range r {
			a[i] = toNumber[uint8](v)
		}
		return ToInterface{a}.Val
	}

	// array - float
	if refT == VarType["[]float64"] {
		r := Conv.ToArray(val)
		a := make([]float64, len(r))
		for i, v := range r {
			a[i] = toNumber[float64](v)
		}
		return ToInterface{a}.Val
	}

	if refT == VarType["[]float32"] {
		r := Conv.ToArray(val)
		a := make([]float32, len(r))
		for i, v := range r {
			a[i] = toNumber[float32](v)
		}
		return ToInterface{a}.Val
	}

	// map - basic
	if refT == VarType["map[string]interface{}"] {
		return ToInterface{Conv.ToMap(val)}.Val
	}

	if refT == VarType["map[byte]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[byte]interface{}, len(r))
		for i, v := range r {
			if b := toString[[]byte](i); len(b) != 0 {
				m[b[0]] = v
			}else{
				m[byte(toNumber[int32](i))] = v
			}
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[rune]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[rune]interface{}, len(r))
		for i, v := range r {
			if b := toString[[]byte](i); len(b) != 0 {
				m[rune(b[0])] = v
			}else{
				m[rune(byte(toNumber[int32](i)))] = v
			}
		}
		return ToInterface{m}.Val
	}

	// map - int
	if refT == VarType["map[int]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[int]interface{}, len(r))
		for i, v := range r {
			m[toNumber[int](i)] = v
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[int64]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[int64]interface{}, len(r))
		for i, v := range r {
			m[toNumber[int64](i)] = v
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[int32]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[int32]interface{}, len(r))
		for i, v := range r {
			m[toNumber[int32](i)] = v
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[int16]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[int16]interface{}, len(r))
		for i, v := range r {
			m[toNumber[int16](i)] = v
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[int8]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[int8]interface{}, len(r))
		for i, v := range r {
			m[toNumber[int8](i)] = v
		}
		return ToInterface{m}.Val
	}

	// map - uint
	if refT == VarType["map[uintptr]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[uintptr]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uintptr](i)] = v
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[uint]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[uint]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uint](i)] = v
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[uint64]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[uint64]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uint64](i)] = v
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[uint32]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[uint32]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uint32](i)] = v
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[uint16]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[uint16]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uint16](i)] = v
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[uint8]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[uint8]interface{}, len(r))
		for i, v := range r {
			m[toNumber[uint8](i)] = v
		}
		return ToInterface{m}.Val
	}

	// map - float
	if refT == VarType["map[float64]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[float64]interface{}, len(r))
		for i, v := range r {
			m[toNumber[float64](i)] = v
		}
		return ToInterface{m}.Val
	}

	if refT == VarType["map[float32]interface{}"] {
		r := Conv.ToMap(val)
		m := make(map[float32]interface{}, len(r))
		for i, v := range r {
			m[toNumber[float32](i)] = v
		}
		return ToInterface{m}.Val
	}

	return nil
}

// ToVarTypeInterface attempts to convert an interface to match the unknown type of another interface
//
// this method is similar to the ToType method, but it simply returns nil if it cannot find the proper var type
func TypeEqual(val1 interface{}, val2 interface{}) bool {
	val2 = ToVarTypeInterface(val2, val1)
	refT := reflect.TypeOf(val1)

	if val1 == nil {
		return val2 == nil
	}else if val2 == nil {
		return val1 == nil
	}

	// basic
	if refT == VarType["string"] {
		return val1.(string) == val2.(string)
	}

	if refT == VarType["[]byte"] {
		return bytes.Equal(val1.([]byte), val2.([]byte))
	}

	if refT == VarType["byte"] {
		return val1.(byte) == val2.(byte)
	}

	if refT == VarType["rune"] {
		return val1.(rune) == val2.(rune)
	}

	if refT == VarType["bool"] {
		return val1.(bool) == val2.(bool)
	}

	// int
	if refT == VarType["int"] {
		return val1.(int) == val2.(int)
	}

	if refT == VarType["int64"] {
		return val1.(int64) == val2.(int64)
	}

	if refT == VarType["int32"] {
		return val1.(int32) == val2.(int32)
	}

	if refT == VarType["int16"] {
		return val1.(int16) == val2.(int16)
	}

	if refT == VarType["int8"] {
		return val1.(int8) == val2.(int8)
	}

	// uint
	if refT == VarType["uintptr"] {
		return val1.(uintptr) == val2.(uintptr)
	}

	if refT == VarType["uint"] {
		return val1.(uint) == val2.(uint)
	}

	if refT == VarType["uint64"] {
		return val1.(uint64) == val2.(uint64)
	}

	if refT == VarType["uint32"] {
		return val1.(uint32) == val2.(uint32)
	}

	if refT == VarType["uint16"] {
		return val1.(uint16) == val2.(uint16)
	}

	if refT == VarType["uint8"] {
		return val1.(uint8) == val2.(uint8)
	}

	// float
	if refT == VarType["float64"] {
		return val1.(float64) == val2.(float64)
	}

	if refT == VarType["float32"] {
		return val1.(float32) == val2.(float32)
	}

	// array - basic
	if refT == VarType["interface{}"] {
		return ArrayEqual(val1.([]interface{}), val2.([]interface{}))
	}

	if refT == VarType["[]string"] {
		return ArrayEqual(val1.([]string), val2.([]string))
	}

	if refT == VarType["[][]byte"] {
		return ArrayEqual(val1.([][]byte), val2.([][]byte))
	}

	if refT == VarType["[]rune"] {
		return ArrayEqual(val1.([]rune), val2.([]rune))
	}

	if refT == VarType["[]bool"] {
		return ArrayEqual(val1.([]bool), val2.([]bool))
	}

	// array - int
	if refT == VarType["[]int"] {
		return ArrayEqual(val1.([]int), val2.([]int))
	}

	if refT == VarType["[]int64"] {
		return ArrayEqual(val1.([]int64), val2.([]int64))
	}

	if refT == VarType["[]int32"] {
		return ArrayEqual(val1.([]int32), val2.([]int32))
	}

	if refT == VarType["[]int16"] {
		return ArrayEqual(val1.([]int16), val2.([]int16))
	}

	if refT == VarType["[]int8"] {
		return ArrayEqual(val1.([]int8), val2.([]int8))
	}

	// array - uint
	if refT == VarType["[]uintptr"] {
		return ArrayEqual(val1.([]uintptr), val2.([]uintptr))
	}

	if refT == VarType["[]uint"] {
		return ArrayEqual(val1.([]uint), val2.([]uint))
	}

	if refT == VarType["[]uint64"] {
		return ArrayEqual(val1.([]uint64), val2.([]uint64))
	}

	if refT == VarType["[]uint32"] {
		return ArrayEqual(val1.([]uint32), val2.([]uint32))
	}

	if refT == VarType["[]uint16"] {
		return ArrayEqual(val1.([]uint16), val2.([]uint16))
	}

	if refT == VarType["[]uint8"] {
		return ArrayEqual(val1.([]uint8), val2.([]uint8))
	}

	// array - float
	if refT == VarType["[]float64"] {
		return ArrayEqual(val1.([]float64), val2.([]float64))
	}

	if refT == VarType["[]float32"] {
		return ArrayEqual(val1.([]float32), val2.([]float32))
	}

	// map - basic
	if refT == VarType["map[string]interface{}"] {
		return MapEqual(val1.(map[string]interface{}), val2.(map[string]interface{}))
	}

	if refT == VarType["map[byte]interface{}"] {
		return MapEqual(val1.(map[byte]interface{}), val2.(map[byte]interface{}))
	}

	if refT == VarType["map[rune]interface{}"] {
		return MapEqual(val1.(map[rune]interface{}), val2.(map[rune]interface{}))
	}

	// map - int
	if refT == VarType["map[int]interface{}"] {
		return MapEqual(val1.(map[int]interface{}), val2.(map[int]interface{}))
	}

	if refT == VarType["map[int64]interface{}"] {
		return MapEqual(val1.(map[int64]interface{}), val2.(map[int64]interface{}))
	}

	if refT == VarType["map[int32]interface{}"] {
		return MapEqual(val1.(map[int32]interface{}), val2.(map[int32]interface{}))
	}

	if refT == VarType["map[int16]interface{}"] {
		return MapEqual(val1.(map[int16]interface{}), val2.(map[int16]interface{}))
	}

	if refT == VarType["map[int8]interface{}"] {
		return MapEqual(val1.(map[int8]interface{}), val2.(map[int8]interface{}))
	}

	// map - uint
	if refT == VarType["map[uintptr]interface{}"] {
		return MapEqual(val1.(map[uintptr]interface{}), val2.(map[uintptr]interface{}))
	}

	if refT == VarType["map[uint]interface{}"] {
		return MapEqual(val1.(map[uint]interface{}), val2.(map[uint]interface{}))
	}

	if refT == VarType["map[uint64]interface{}"] {
		return MapEqual(val1.(map[uint64]interface{}), val2.(map[uint64]interface{}))
	}

	if refT == VarType["map[uint32]interface{}"] {
		return MapEqual(val1.(map[uint32]interface{}), val2.(map[uint32]interface{}))
	}

	if refT == VarType["map[uint16]interface{}"] {
		return MapEqual(val1.(map[uint16]interface{}), val2.(map[uint16]interface{}))
	}

	if refT == VarType["map[uint8]interface{}"] {
		return MapEqual(val1.(map[uint8]interface{}), val2.(map[uint8]interface{}))
	}

	// map - float
	if refT == VarType["map[float64]interface{}"] {
		return MapEqual(val1.(map[float64]interface{}), val2.(map[float64]interface{}))
	}

	if refT == VarType["map[float32]interface{}"] {
		return MapEqual(val1.(map[float32]interface{}), val2.(map[float32]interface{}))
	}

	return false
}

// ToVarType grabs the type from another var as a reference, and runs the `ToType` with the ref type
func ToVarType[T SupportedType](val interface{}, ref T) T {
	return ToType[T](val)
}

// IsZeroOfUnderlyingType can be used to determine if an interface{} in null or empty
func IsZeroOfUnderlyingType(x interface{}) bool {
	// return x == nil || x == reflect.Zero(reflect.TypeOf(x)).Interface()
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
