package brotli

import (
	"bytes"
	"io"

	"github.com/andybalholm/brotli"
)

// brotli.Zip Compresses with brotli to a utf8 []byte
//
// @quality: 0-11 (0 = fastest) (11 = best)
func Zip(msg []byte, quality ...int) ([]byte, error) {
	q := 6
	if len(quality) != 0 {
		q := quality[0]
		if q < 0 {
			q = 0
		}else if q > 11 {
			q = 11
		}
	}

	var b bytes.Buffer
	w := brotli.NewWriterLevel(&b, q)
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

// brotli.UnZip Decompresses with brotli from a utf8 []byte
func UnZip(b []byte) ([]byte, error) {
	rdata := bytes.NewReader(b)
	r := brotli.NewReader(rdata)
	s, err := io.ReadAll(r)
	if err != nil {
		return []byte{}, err
	}
	return s, nil
}
