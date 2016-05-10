package contentinfo

import "testing"

func TestLess(t *testing.T) {
	a := Providers{
		"a",
		nil,
	}
	b := Providers{
		"b",
		nil,
	}
	aa := Providers{
		"aa",
		nil,
	}

	if !a.Less(b) {
		t.Error(a, b)
	}
	if b.Less(a) {
		t.Error(b, a)
	}
	if a.Less(a) {
		t.Error(a, a)
	}

	if !a.Less(aa) {
		t.Error(a, aa)
	}
	if b.Less(aa) {
		t.Error(b, aa)
	}
}
