package configure

import (
	"fmt"
	"testing"
)

var testConfig = "../test/test_conf.json"

func TestReadConfigure(t *testing.T) {
	conf := NewConfigure(testConfig)
	testExterls := func() map[string]string {
		return map[string]string{"example": "http://example.com/"}
	}()
	func() {
		fmt.Println("test conf struct all.")
		if conf == nil {
			t.Fatalf("conf is nil.")
		}
	}()

	func() {
		fmt.Println("port setting test.")
		if conf.Port() != 8080 {
			t.Fatalf("port is not equal 8080, value is \"%v\"", conf.Port())
		}
	}()

	func() {
		fmt.Println("local image path test.")
		if conf.LocalImagePath() != "../img/" {
			t.Fatalf("value is %v", conf.LocalImagePath())
		}
	}()

	func() {
		fmt.Println("max size test")
		w, h := conf.MaxSize()
		if w != 5000 {
			t.Fatalf("value is %v", w)
		}
		if h != 5000 {
			t.Fatalf("value is %v", h)
		}
	}()

	func() {
		fmt.Println("externals test.")
		externals := conf.Externals()
		if externals == nil {
			t.Fatalf("value is %v", nil)
		}
		for k, v := range testExterls {
			if v != externals[k] {
				t.Fatalf("k: %v, v: %v", k, externals[k])
			}
		}
	}()

}
