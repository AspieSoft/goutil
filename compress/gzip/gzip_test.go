package gzip

import (
	"errors"
	"testing"
)

func TestCompress(t *testing.T){
	msg := "This is a test"
	comp, err := Zip([]byte(msg))
	if err != nil {
		t.Error(err)
	}
	dec, err := UnZip(comp)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("Gzip did not return the correct output"))
	}
}
