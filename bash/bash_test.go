package bash

import (
	"bytes"
	"testing"
)

func TestBash(t *testing.T){
	out, err := Run([]string{`echo`, `test`}, "", nil)
	if err != nil {
		t.Error(err)
	}
	if !bytes.Equal(out, []byte("test")) && !bytes.Equal(out, []byte("test\n")) && !bytes.Equal(out, []byte("test\r\n")) {
		t.Error("incorrect output [test]:", string(out))
	}
}
