package goutil

import (
	"bytes"
	"encoding/json"
	"io"
	"strings"

	"github.com/AspieSoft/go-regex/v5/re2-opt"
)

type encodeHtml struct {}
var HTML *encodeHtml = &encodeHtml{}

type encodeJson struct {}
var JSON *encodeJson = &encodeJson{}

var regEscHTML *regex.Regexp = regex.Comp(`[<>&]`)
var regEscFixAmp *regex.Regexp = regex.Comp(`&amp;(amp;)*`)

// EscapeHTML replaces HTML characters with html entities
//
// Also prevents and removes &amp;amp; from results
func (encHtml *encodeHtml) Escape(html []byte) []byte {
	html = regEscHTML.RepFunc(html, func(data func(int) []byte) []byte {
		if bytes.Equal(data(0), []byte("<")) {
			return []byte("&lt;")
		} else if bytes.Equal(data(0), []byte(">")) {
			return []byte("&gt;")
		}
		return []byte("&amp;")
	})
	return regEscFixAmp.RepStr(html, []byte("&amp;"))
}

var regEscHTMLArgs *regex.Regexp = regex.Comp(`([\\]*)([\\"'\'])`)

// EscapeHTMLArgs escapes quotes and backslashes for use within HTML quotes
// @quote can be used to only escape specific quotes or chars
func (encHtml *encodeHtml) EscapeArgs(html []byte, quote ...byte) []byte {
	if len(quote) == 0 {
		quote = []byte("\"'`")
	}

	return regEscHTMLArgs.RepFunc(html, func(data func(int) []byte) []byte {
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
func (encJson *encodeJson) Stringify(data interface{}, ind ...int) ([]byte, error) {
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
func (encJson *encodeJson) Parse(b []byte) (map[string]interface{}, error) {
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
func (encJson *encodeJson) Decode(data io.Reader) (map[string]interface{}, error) {
	var res map[string]interface{}
	err := json.NewDecoder(data).Decode(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// DeepCopyJson will stringify and parse json to create a deep copy and escape pointers
func (encJson *encodeJson) DeepCopy(data map[string]interface{}) (map[string]interface{}, error) {
	b, err := encJson.Stringify(data)
	if err != nil {
		return nil, err
	}
	return encJson.Parse(b)
}
