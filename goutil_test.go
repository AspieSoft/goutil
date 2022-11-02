package goutil

import (
	"errors"
	"testing"
)

func Test(t *testing.T){
	if val := ToString(0); val != "0" {
		t.Error("[", val, "]\n", errors.New("ToString Method Failed"))
	}

	if val := ToInt("1"); val != 1 {
		t.Error("[", val, "]\n", errors.New("ToInt Method Failed"))
	}

	if val, err := JoinPath("test", "1"); err != nil {
		t.Error("[", val, "]\n", errors.New("JoinPath Method Failed"))
	}

	if val, err := JoinPath("test", "../out/of/root"); err == nil {
		t.Error("[", val, "]\n", errors.New("JoinPath Method Leaked Outsite The Root"))
	}

	if args := MapArgs([]string{"arg1", "--key=value", "--bool", "-flags"}); args["0"] != "arg1" || args["bool"] != "true" || args["key"] != "value" || args["f"] != "true" || args["l"] != "true" || args["s"] != "true" {
		t.Error(args, "\n", errors.New("MapArgs Produced The Wrong Output"))
	}
}

func TestEncrypt(t *testing.T){
	msg := "This is a test"
	enc, err := Encrypt([]byte(msg), []byte("MyKey123"))
	if err != nil {
		t.Error(err)
	}
	dec, err := Decrypt(enc, []byte("MyKey123"))
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("Decrypt did not return the correct output"))
	}

	hash, err := NewHash([]byte(msg), []byte("MyKey123"))
	if err != nil {
		t.Error(err)
	}
	if !CompareHash([]byte(msg), hash, []byte("MyKey123")) {
		t.Error("[", msg, "]\n", errors.New("CompareHash did not return true"))
	}
}

func TestEncryptLocal(t *testing.T) {
	msg := "This is a test"
	key := []byte("ncu3l89yn298hidh8nxoiauj932oijqkhd8o3nq7i")
	enc, err := EncryptLocal([]byte(msg), key)
	if err != nil {
		t.Error(err)
	}
	dec, err := DecryptLocal(enc, key)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("DecryptLocal did not return the correct output"))
	}
}
