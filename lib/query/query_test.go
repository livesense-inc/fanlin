package query

import (
	"net/http"
	"testing"
)

var getRquest, _ = http.NewRequest("GET", "http://exsample.com/?w=100&h=100&psrc=http://google.co.jp/&rgb=255,255,255&crop=true&quality=75", nil)

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

func TestCrop(t *testing.T) {
	q := NewQueryFromGet(getRquest)

	if !q.Crop() {
		t.Fatalf("crop is false.")
	}
}

func TestQuality(t *testing.T) {
	q := NewQueryFromGet(getRquest)

	if q.Quality() != 75 {
		t.Fatalf("quality is %d.", q.Quality())
	}
}

func TestAVIF(t *testing.T) {
	q := NewQueryFromGet(getRquest)

	if q.UseAVIF() {
		t.Fatalf("avif is true.")
	}
}
