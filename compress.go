package goutil

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"io"

	"github.com/andybalholm/brotli"
	"github.com/cespare/go-smaz"
)

type compGzip struct {}
var GZIP *compGzip = &compGzip{}

type compBrotli struct {}
var BROTLI *compBrotli = &compBrotli{}

type compSmaz struct {}
var SMAZ *compSmaz = &compSmaz{}

// GZIP.Zip is Gzip compression to a utf8 []byte
func (comp *compGzip) Zip(msg []byte) ([]byte, error) {
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
	return b.Bytes(), nil
}

// GZIP.UnZip is Gzip decompression from a utf8 []byte
func (comp *compGzip) UnZip(b []byte) ([]byte, error) {
	rdata := bytes.NewReader(b)
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

// BROTLI.Zip Compresses with brotli to a utf8 []byte
//
// @quality: 0-11 (0 = fastest) (11 = best)
func (comp *compBrotli) Zip(msg []byte, quality ...int) ([]byte, error) {
	q := 11
	if len(quality) != 0 {
		q := quality[0]
		if q < 0 {
			q = 0
		}else if q > 11 {
			q = 11
		}
	}

	var b bytes.Buffer
	w := brotli.NewWriter(&b)
	brotli.NewWriterLevel(w, q)
	if _, err := w.Write([]byte(msg)); err != nil {
		return []byte{}, err
	}
	if err := w.Flush(); err != nil {
		return []byte{}, err
	}
	if err := w.Close(); err != nil {
		return []byte{}, err
	}

	return b.Bytes(), nil
}

// BROTLI.UnZip Decompresses with brotli from a utf8 []byte
func (comp *compBrotli) UnZip(b []byte) ([]byte, error) {
	rdata := bytes.NewReader(b)
	r := brotli.NewReader(rdata)
	s, err := io.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}
	return s, nil
}

// SMAZ.Zip Compresses with SMAZ from a utf8 []byte
//
// @encode: true = encode to base64
func (comp *compSmaz) Zip(b []byte, encode ...bool) ([]byte) {
	b = smaz.Compress(b)
	if len(encode) != 0 && encode[0] == true {
		return []byte(base64.StdEncoding.EncodeToString(b))
	}
	return b
}

// SMAZ.UnZip Decompresses with SMAZ from a utf8 []byte
//
// this method will try to decode from base64, or will skip that step if it fails
// (this means you do not have to know if you encoded your string to base64 on compression)
func (comp *compSmaz) UnZip(b []byte) ([]byte, error) {
	if dec, err := base64.StdEncoding.DecodeString(string(b)); err == nil {
		if dec, err = smaz.Decompress(dec); err == nil {
			return dec, nil
		}
	}
	return smaz.Decompress(b)
}