package query

import (
	"net/http"
	"testing"
)

var getRquest, _ = http.NewRequest("GET", "https://example.com/?w=300&h=200&rgb=255,255,255&crop=true&quality=75", nil)

func TestNewGet(t *testing.T) {
	q := NewQueryFromGet(getRquest)

	if q == nil {
		t.Fatalf("query is not allocated.")
	}
}

func TestBounds(t *testing.T) {
	q := NewQueryFromGet(getRquest)
	b := Bounds{300, 200}

	if *q.Bounds() != b {
		t.Fatalf("Bounds: %v.", *q.Bounds())
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

func TestWebP(t *testing.T) {
	q := NewQueryFromGet(getRquest)

	if q.UseWebP() {
		t.Fatalf("webp is true.")
	}
}

func TestAVIF(t *testing.T) {
	q := NewQueryFromGet(getRquest)

	if q.UseAVIF() {
		t.Fatalf("avif is true.")
	}
}
