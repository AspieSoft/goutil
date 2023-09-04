package smaz

import (
	"errors"
	"testing"
)

func TestCompress(t *testing.T){
	msg := "This is a test"
	comp := Zip([]byte(msg), true)
	dec, err := UnZip(comp)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("SMAZ did not return the correct output"))
	}
}
