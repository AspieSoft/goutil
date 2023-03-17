package goutil

import (
	"bytes"
	"reflect"
	"strings"
)

type clean struct {}
var Clean *clean = &clean{}

// Clean.Str will sanitizes a string to valid UTF-8
func (sanitize *clean) Str(str string) string {
	//todo: sanitize inputs
	str = strings.ToValidUTF8(str, "")
	return str
}

// Clean.Byte will sanitizes a []byte to valid UTF-8
func (sanitize *clean) Byte(b []byte) []byte {
	//todo: sanitize inputs
	b = bytes.ToValidUTF8(b, []byte{})
	return b
}

// Clean.Array runs `Clean.Str` on an []interface{}
//
// Clean.Str sanitizes a string to valid UTF-8
func (sanitize *clean) Array(data []interface{}) []interface{} {
	cData := []interface{}{}
	for key, val := range data {
		t := reflect.TypeOf(val)
		if t == VarType["string"] {
			cData[key] = sanitize.Str(val.(string))
		}else if t == VarType["int"] || t == VarType["float64"] || t == VarType["float32"] || t == VarType["bool"] {
			cData[key] = val
		}else if t == VarType["[]byte"] {
			cData[key] = sanitize.Str(string(val.([]byte)))
		}else if t == VarType["byte"] {
			cData[key] = sanitize.Str(string(val.(byte)))
		}else if t == VarType["int32"] {
			cData[key] = sanitize.Str(string(val.(int32)))
		}else if t == VarType["[]interface{}"] {
			cData[key] = sanitize.Array(val.([]interface{}))
		}else if t == VarType["map[string]interface{}"] {
			cData[key] = sanitize.Map(val.(map[string]interface{}))
		}
	}
	return cData
}

// Clean.Map runs `Clean.Str` on a map[string]interface{}
//
// Clean.Str sanitizes a string to valid UTF-8
func (sanitize *clean) Map(data map[string]interface{}) map[string]interface{} {
	cData := map[string]interface{}{}
	for key, val := range data {
		key = sanitize.Str(key)

		t := reflect.TypeOf(val)
		if t == VarType["string"] {
			cData[key] = sanitize.Str(val.(string))
		}else if t == VarType["int"] || t == VarType["float64"] || t == VarType["float32"] || t == VarType["bool"] {
			cData[key] = val
		}else if t == VarType["[]byte"] {
			cData[key] = sanitize.Str(string(val.([]byte)))
		}else if t == VarType["byte"] {
			cData[key] = sanitize.Str(string(val.(byte)))
		}else if t == VarType["int32"] {
			cData[key] = sanitize.Str(string(val.(int32)))
		}else if t == VarType["[]interface{}"] {
			cData[key] = sanitize.Array(val.([]interface{}))
		}else if t == VarType["map[string]interface{}"] {
			cData[key] = sanitize.Map(val.(map[string]interface{}))
		}
	}

	return cData
}

// CleanJSON runs `Clean.Str` on a complex json object recursively
//
// Clean.Str sanitizes a string to valid UTF-8
func (sanitize *clean) JSON(val interface{}) interface{} {
	t := reflect.TypeOf(val)
	if t == VarType["string"] {
		return sanitize.Str(val.(string))
	}else if t == VarType["int"] || t == VarType["float64"] || t == VarType["float32"] || t == VarType["bool"] {
		return val
	}else if t == VarType["[]byte"] {
		return sanitize.Byte(val.([]byte))
	}else if t == VarType["byte"] {
		return sanitize.Byte([]byte{val.(byte)})
	}else if t == VarType["int32"] {
		return sanitize.Str(string(val.(int32)))
	}else if t == VarType["[]interface{}"] {
		return sanitize.Array(val.([]interface{}))
	}else if t == VarType["map[string]interface{}"] {
		return sanitize.Map(val.(map[string]interface{}))
	}
	return nil
}
