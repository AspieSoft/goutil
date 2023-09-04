package crypt

import (
	"errors"
	"testing"
)

func TestEncrypt(t *testing.T){
	msg := "This is a test"
	enc, err := CFB.Encrypt([]byte(msg), []byte("MyKey123"))
	if err != nil {
		t.Error(err)
	}
	dec, err := CFB.Decrypt(enc, []byte("MyKey123"))
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("Decrypt did not return the correct output"))
	}

	hash, err := Hash.New([]byte(msg), []byte("MyKey123"))
	if err != nil {
		t.Error(err)
	}
	if !Hash.Compare([]byte(msg), hash, []byte("MyKey123")) {
		t.Error("[", msg, "]\n", errors.New("CompareHash did not return true"))
	}
}
