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
}

func TestEncryptLocal(t *testing.T) {
	msg := "This is a test"
	key := "ncu3l89yn298hidh8nxoiauj932oijqkhd8o3nq7i"
	enc, err := EncryptLocal([]byte(msg), key)
	if err != nil {
		t.Error(err)
	}
	dec, err := DecryptLocal(enc, key)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("EncryptLocal Did Not Decrypt Correct Output"))
	}
}
