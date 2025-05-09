package query

import (
	"image/color"
	"net/http"
	"strconv"
	"strings"
)

type Bounds struct {
	W uint
	H uint
}

type Query struct {
	b         Bounds
	fillColor color.Color
	crop      bool
	quality   int
	grayscale bool
	inverse   bool
	useWebp   bool
	useAvif   bool
}

func NewQueryFromGet(r *http.Request) *Query {
	params := r.URL.Query()
	q := Query{}
	w, _ := strconv.Atoi(params.Get("w"))
	h, _ := strconv.Atoi(params.Get("h"))
	rgb := strings.Split(strings.Trim(params.Get("rgb"), "\""), ",")

	var c color.Color
	if len(rgb) == 3 {
		c = func() color.Color {
			r, _ := strconv.Atoi(rgb[0])
			g, _ := strconv.Atoi(rgb[1])
			b, _ := strconv.Atoi(rgb[2])
			return color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
		}()
	} else {
		c = nil
	}

	quality, err := strconv.Atoi(params.Get("quality"))
	if err != nil {
		quality = -1
	}

	crop, err := strconv.ParseBool(params.Get("crop"))
	if err != nil {
		crop = false
	}

	grayscale, err := strconv.ParseBool(params.Get("grayscale"))
	if err != nil {
		grayscale = false
	}

	inverse, err := strconv.ParseBool(params.Get("inverse"))
	if err != nil {
		inverse = false
	}

	webp, err := strconv.ParseBool(params.Get("webp"))
	if err != nil {
		webp = false
	}

	avif, err := strconv.ParseBool(params.Get("avif"))
	if err != nil {
		avif = false
	}

	q.b.H = uint(h)
	q.b.W = uint(w)
	q.fillColor = c
	q.crop = crop
	q.quality = quality
	q.grayscale = grayscale
	q.inverse = inverse
	q.useWebp = webp
	q.useAvif = avif
	return &q
}

func (q *Query) Bounds() *Bounds {
	return &q.b
}

func (q *Query) FillColor() *color.Color {
	return &q.fillColor
}

func (q *Query) Crop() bool {
	return q.crop
}

func (q *Query) Quality() int {
	return q.quality
}

func (q *Query) Grayscale() bool {
	return q.grayscale
}

func (q *Query) Inverse() bool {
	return q.inverse
}

func (q *Query) UseWebP() bool {
	return q.useWebp
}

func (q *Query) UseAVIF() bool {
	return q.useAvif
}
