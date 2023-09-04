package brotli

import (
	"errors"
	"testing"
)

func TestCompress(t *testing.T){
	msg := "This is a test"

	comp, err := Zip([]byte(msg), 11)
	if err != nil {
		t.Error(err)
	}
	dec, err := UnZip(comp)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("Brotli did not return the correct output"))
	}
}
