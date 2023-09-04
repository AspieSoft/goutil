package goutil

import (
	"errors"
	"testing"

	"github.com/AspieSoft/go-regex-re2/v2"
)

func Test(t *testing.T){
	// this test function is for debugging
}

func TestBasic(t *testing.T){
	if val := toString[string](0); val != "0" {
		t.Error("[", val, "]\n", errors.New("toString Method Failed"))
	}

	if val := toNumber[int]("1"); val != 1 {
		t.Error("[", val, "]\n", errors.New("toInt Method Failed"))
	}

	if val := ToType[float64]("5.0"); val != 5.0 {
		t.Error("[", val, "]\n", errors.New("ToType[float64] Method Failed"))
	}

	if val, err := FS.JoinPath("test", "1"); err != nil {
		t.Error("[", val, "]\n", errors.New("JoinPath Method Failed"))
	}

	if val, err := FS.JoinPath("test", "../out/of/root"); err == nil {
		t.Error("[", val, "]\n", errors.New("JoinPath Method Leaked Outsite The Root"))
	}

	if args := MapArgs([]string{"arg1", "--key=value", "--bool", "-flags"}); args["0"] != "arg1" || args["bool"] != "true" || args["key"] != "value" || args["f"] != "true" || args["l"] != "true" || args["s"] != "true" {
		t.Error(args, "\n", errors.New("MapArgs Produced The Wrong Output"))
	}
}

func TestEncrypt(t *testing.T){
	msg := "This is a test"
	enc, err := Crypt.CFB.Encrypt([]byte(msg), []byte("MyKey123"))
	if err != nil {
		t.Error(err)
	}
	dec, err := Crypt.CFB.Decrypt(enc, []byte("MyKey123"))
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("Decrypt did not return the correct output"))
	}

	hash, err := Crypt.Hash.New([]byte(msg), []byte("MyKey123"))
	if err != nil {
		t.Error(err)
	}
	if !Crypt.Hash.Compare([]byte(msg), hash, []byte("MyKey123")) {
		t.Error("[", msg, "]\n", errors.New("CompareHash did not return true"))
	}
}

func TestCompress(t *testing.T){
	msg := "This is a test"
	comp, err := GZIP.Zip([]byte(msg))
	if err != nil {
		t.Error(err)
	}
	dec, err := GZIP.UnZip(comp)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("Gzip did not return the correct output"))
	}

	comp, err = BROTLI.Zip([]byte(msg), 11)
	if err != nil {
		t.Error(err)
	}
	dec, err = BROTLI.UnZip(comp)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("Brotli did not return the correct output"))
	}

	comp = SMAZ.Zip([]byte(msg), true)
	dec, err = SMAZ.UnZip(comp)
	if err != nil {
		t.Error(err)
	}
	if string(dec) != msg {
		t.Error("[", msg, "]\n", errors.New("SMAZ did not return the correct output"))
	}
}

func TestHtmlEscape(t *testing.T) {
	html := regex.JoinBytes([]byte("<a href=\""), HTML.EscapeArgs([]byte(`test 1\\" attr="hack" null="`), '"'), '"', []byte("js=\""), HTML.EscapeArgs([]byte("this.media='all' `line \\n break`"), '"'), []byte("\">"), HTML.Escape([]byte(`<html>element & &amp; &amp;amp; &bad; test</html>`)), []byte("</a>"))
	html = regex.Comp(`href="(?:\\[\"]|.)*?"`).RepStrLit(html, []byte{})

	if regex.Comp(`(hack|<html>|&bad;|&amp;amp;|\\\\n|\\')`).Match(html) {
		t.Error(errors.New("'EscapeHTML' and/or 'EscapeHTMLArgs' method failed to prevent a test html hack properly"))
	}
}
