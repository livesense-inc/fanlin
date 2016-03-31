package imgproxyerr

import (
	"errors"
	"testing"
)

var GetLevel = getLevel
var Level = level
var TestStr = "John Doe"
var testWarn = &Err{WARNING, errors.New("test warning")}
var testError = &Err{ERROR, errors.New("test error")}

func TestNew(t *testing.T) {
	t.Log("test func New()")
	errNil := New("error", nil)
	if errNil != nil {
		t.Error(errNil, errNil.Error())
	}
	errNNil := New("error", errors.New("test error"))
	if errNNil == nil {
		t.Error(errNNil, errNNil.Error())
	}

	func() {
		t.Log("warning -> errorへ上書きが発生するときのテスト")
		warningErr := New(WARNING, errors.New("warning"))
		errorErr := New(ERROR, warningErr)
		if e, ok := errorErr.(*Err); ok {
			if e.Type != ERROR {
				t.Fatal(e)
			}
		} else {
			t.Fatal("can not cast")
		}
	}()

	func() {
		t.Log("warning -> errorへ上書きが発生するときのテスト")
		errorErr := New(ERROR, errors.New("error"))
		warningErr := New(WARNING, errorErr)
		if e, ok := warningErr.(*Err); ok {
			if e.Type != ERROR {
				t.Fatal(e)
			}
		} else {
			t.Fatal("can not cast")
		}
	}()
}

func TestGetLevel(t *testing.T) {
	t.Log("test func getLevel()")
	t.Log("level変数に存在している時")

	if 0 != GetLevel(ERROR) {
		t.Fatal("level:", GetLevel(ERROR), ", status:", ERROR)
	}
	if 1 != GetLevel(WARNING) {
		t.Fatal("level:", GetLevel(WARNING), ", status:", WARNING)
	}
	if -1 != GetLevel("unknown") {
		t.Fatal("level:", GetLevel("unknown"), ", status:", "unknown")
	}
	if func(status string) bool {
		for i := range Level {
			if i == GetLevel(status) {
				return true
			}
		}
		return false
	}(TestStr) {
		t.Fatal("level:", GetLevel(TestStr), ", status", TestStr)
	}
}

func TestError(t *testing.T) {
	t.Log("test Error()")

	if "test warning" != testWarn.Error() {
		t.Fatal(testWarn.Type, testWarn.Error())
	}
	if "test error" != testError.Error() {
		t.Fatal(testError.Type, testError.Error())
	}
}

func TestCmp(t *testing.T) {
	t.Log("test cmp")

	if !testError.cmp(ERROR) {
		t.Fatal(testError.cmp(ERROR))
	}
	if !testWarn.cmp(WARNING) {
		t.Fatal(testWarn.cmp(WARNING))
	}
	if !testError.cmp(WARNING) {
		t.Fatal(testError.cmp(WARNING))
	}
	if testWarn.cmp(ERROR) {
		t.Fatal(testWarn.cmp(ERROR))
	}
}
