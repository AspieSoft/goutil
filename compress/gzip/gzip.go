package gzip

import (
	"bytes"
	"compress/gzip"
	"io"
)

// gzip.Zip is Gzip compression to a utf8 []byte
//
// @quality: 1-9 (1 = fastest) (9 = best)
func Zip(msg []byte, quality ...int) ([]byte, error) {
	q := 6
	if len(quality) != 0 {
		q := quality[0]
		if q < 1 {
			q = 1
		}else if q > 9 {
			q = 9
		}
	}

	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, q)
	if err != nil {
		w = gzip.NewWriter(&b)
	}

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

// gzip.UnZip is Gzip decompression from a utf8 []byte
func UnZip(b []byte) ([]byte, error) {
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
