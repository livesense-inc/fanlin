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
	b                      Bounds
	preliminaryImageSource string
	fillColor              color.Color
	quality                int
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

	q.b.H = uint(h)
	q.b.W = uint(w)
	q.preliminaryImageSource = params.Get("psrc")
	q.fillColor = c
	q.quality = quality
	return &q
}

func (q *Query) Bounds() *Bounds {
	return &q.b
}

func (q *Query) PreliminaryImageSource() string {
	return q.preliminaryImageSource
}

func (q *Query) FillColor() *color.Color {
	return &q.fillColor
}

func (q *Query) Quality() int {
	return q.quality
}
