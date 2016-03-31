package query

import (
	"net/http"
	"testing"
)

var getRquest, _ = http.NewRequest("GET", "http://exsample.com/?w=100&h=100&psrc=http://google.co.jp/&rgb=255,255,255", nil)

func TestNewGet(t *testing.T) {
	q := NewQueryFromGet(getRquest)

	if q == nil {
		t.Fatalf("query is not allocated.")
	}
}

func TestBounds(t *testing.T) {
	q := NewQueryFromGet(getRquest)
	b := Bounds{100, 100}

	if *q.Bounds() != b {
		t.Fatalf("Bounds: %v.", *q.Bounds())
	}
}

func TestPreliminaryImageSource(t *testing.T) {
	q := NewQueryFromGet(getRquest)

	if q.PreliminaryImageSource() != "http://google.co.jp/" {
		t.Fatalf("psrc is %v.", q.PreliminaryImageSource())
	}
}

func TestGetFillColor(t *testing.T) {
	q := NewQueryFromGet(getRquest)

	if q.FillColor() == nil {
		t.Fatalf("fillcolor is nil.")
	}
}
