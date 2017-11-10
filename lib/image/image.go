package imageprocessor

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math"

	"github.com/BurntSushi/graphics-go/graphics"
	"github.com/BurntSushi/graphics-go/graphics/interp"
	"github.com/livesense-inc/fanlin/lib/error"
	"github.com/nfnt/resize"
	"github.com/rwcarlsen/goexif/exif"
	_ "golang.org/x/image/bmp"
)

var affines map[int]graphics.Affine = map[int]graphics.Affine{
	1: graphics.I,
	2: graphics.I.Scale(-1, 1),
	3: graphics.I.Scale(-1, -1),
	4: graphics.I.Scale(1, -1),
	5: graphics.I.Rotate(toRadian(90)).Scale(-1, 1),
	6: graphics.I.Rotate(toRadian(90)),
	7: graphics.I.Rotate(toRadian(-90)).Scale(-1, 1),
	8: graphics.I.Rotate(toRadian(-90)),
}

type Image struct {
	img    image.Image
	format string
}

func max(v uint, max uint) uint {
	if v > max {
		return max
	}
	return v
}

func EncodeJpeg(img *image.Image, q int) ([]byte, error) {
	if *img == nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, errors.New("img is nil."))
	}

	if !(0 <= q && q <= 100) {
		q = jpeg.DefaultQuality
	}

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, *img, &jpeg.Options{Quality: q})
	return buf.Bytes(), imgproxyerr.New(imgproxyerr.WARNING, err)
}

func EncodePNG(img *image.Image, q int) ([]byte, error) {
	if *img == nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, errors.New("img is nil."))
	}

	// Split quality from 0 to 100 in 4 CompressionLevel
	// https://golang.org/pkg/image/png/#CompressionLevel
	var e png.Encoder
	switch {
	case 0 <= q && q <= 25:
		e.CompressionLevel = png.BestCompression
	case 25 < q && q <= 50:
		e.CompressionLevel = png.DefaultCompression
	case 50 < q && q <= 75:
		e.CompressionLevel = png.BestSpeed
	case 75 < q && q <= 100:
		e.CompressionLevel = png.NoCompression
	default:
		e.CompressionLevel = png.DefaultCompression
	}

	buf := new(bytes.Buffer)
	err := e.Encode(buf, *img)
	return buf.Bytes(), imgproxyerr.New(imgproxyerr.WARNING, err)
}

func EncodeGIF(img *image.Image, q int) ([]byte, error) {
	if *img == nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, errors.New("img is nil."))
	}

	// GIF is not support quality

	buf := new(bytes.Buffer)
	err := gif.Encode(buf, *img, &gif.Options{})
	return buf.Bytes(), imgproxyerr.New(imgproxyerr.WARNING, err)
}

//DecodeImage is return image.Image
func DecodeImage(bin []byte) (*Image, error) {
	img, format, err := Decode(bin)
	return &Image{img: img, format: format}, imgproxyerr.New(imgproxyerr.WARNING, err)
}

//アス比を維持した時の長さを取得する
func keepAspect(img image.Image, w uint, h uint) (uint, uint) {
	r := img.Bounds()
	if int(w)*r.Max.Y < int(h)*r.Max.X {
		return w, 0
	} else {
		return 0, h
	}
}

func resizeImage(img image.Image, w uint, h uint, maxWidth uint, maxHeight uint) image.Image {
	if img == nil {
		return nil
	}
	//大きすぎる値はサポートしない
	w = max(w, maxWidth)
	h = max(h, maxHeight)
	w, h = keepAspect(img, w, h)
	// 速度・負荷的な問題出た時はアルゴリズム変更
	return resize.Resize(w, h, img, resize.Lanczos3)
}

func resizeAndFillImage(img image.Image, w uint, h uint, c color.Color, maxWidth uint, maxHeight uint) image.Image {
	if img == nil {
		return nil
	}
	if maxWidth < w || maxHeight < h {
		return img
	}
	ch0 := make(chan image.Image)
	ch1 := make(chan *image.RGBA)

	// ココらへんの並列化はベンチマーク次第で変更する
	go func() {
		ch0 <- resizeImage(img, w, h, maxWidth, maxHeight)
	}()
	go func() {
		ch1 <- image.NewRGBA(image.Rect(0, 0, int(w), int(h)))
	}()
	resizedImage := <-ch0
	m := <-ch1

	draw.Draw(m, m.Bounds(), &image.Uniform{c}, image.ZP, draw.Src)

	//画像の中心座標を計算
	centerH := int(h)/2 - (resizedImage.Bounds().Max.Y / 2)
	centerW := int(w)/2 - (resizedImage.Bounds().Max.X / 2)

	if resizedImage.Bounds().Max.X == int(w) {
		draw.Draw(m, m.Bounds(), resizedImage, resizedImage.Bounds().Min.Sub(image.Pt(0, centerH)), draw.Over)
	} else if resizedImage.Bounds().Max.Y == int(h) {
		draw.Draw(m, m.Bounds(), resizedImage, resizedImage.Bounds().Min.Sub(image.Pt(centerW, 0)), draw.Over)
	} else {
		return resizedImage
	}
	return m
}

func (i *Image) ResizeAndFill(w uint, h uint, c color.Color, maxW uint, maxH uint) {
	if maxW < w || maxH < h {
		return
	}
	if h == 0 || w == 0 {
		return
	}
	if c == nil {
		i.img = resizeImage(i.img, w, h, maxW, maxH)
		return
	}
	i.img = resizeAndFillImage(i.img, w, h, c, maxW, maxH)
}

func crop(img image.Image, w uint, h uint) image.Image {
	if img == nil {
		return nil
	}
	if h == 0 || w == 0 {
		return img
	}

	orgW := img.Bounds().Max.X
	orgH := img.Bounds().Max.Y

	r := float64(orgW) / float64(w)
	if (float64(orgW) / float64(orgH)) > (float64(w) / float64(h)) {
		r = float64(orgH) / float64(h)
	}

	startW := orgW/2 - int(float64(w)*r/2)
	startH := orgH/2 - int(float64(h)*r/2)

	result := image.NewRGBA(image.Rect(0, 0, int(float64(w)*r), int(float64(h)*r)))

	for y := 0; y < int(float64(h)*r); y++ {
		for x := 0; x < int(float64(w)*r); x++ {
			c := img.At(x+startW, y+startH)
			result.Set(x, y, c)
		}
	}

	return result
}

func (i *Image) Crop(w uint, h uint) {
	i.img = crop(i.img, w, h)
}

func (i *Image) GetImg() *image.Image {
	return &i.img
}

func (i *Image) GetFormat() string {
	return i.format
}

func Set404Image(path string, w uint, h uint, c color.Color, maxW uint, maxH uint) ([]byte, error) {
	bin, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
	}
	img, err := DecodeImage(bin)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.ERROR, err)
	}
	img.ResizeAndFill(w, h, c, maxW, maxH)
	return EncodeJpeg(img.GetImg(), jpeg.DefaultQuality)
}

func toRadian(n int) float64 {
	return float64(n) * math.Pi / 180.0
}

func applyOrientation(s image.Image, o int) (d draw.Image, e error) {
	bounds := s.Bounds()
	if o == 0 {
		o = 1
	}
	if o >= 5 && o <= 8 {
		bounds = rotateRect(bounds)
	}
	d = image.NewRGBA64(bounds)
	affine := affines[o]
	e = affine.TransformCenter(d, s, interp.Bilinear)
	return
}

func rotateRect(r image.Rectangle) image.Rectangle {
	s := r.Size()
	return image.Rectangle{r.Min, image.Point{s.Y, s.X}}
}

func readOrientation(r io.Reader) (o int, err error) {
	e, err := exif.Decode(r)
	if err != nil {
		return
	}
	tag, err := e.Get(exif.Orientation)
	if err != nil {
		return
	}
	o, err = tag.Int(0)
	if err != nil {
		return
	}
	return
}

func Decode(b []byte) (d image.Image, format string, err error) {
	s, format, err := image.Decode(bytes.NewReader(b))
	if err != nil {
		return
	}
	o, err := readOrientation(bytes.NewReader(b))
	if err != nil {
		return s, format, nil
	}
	d, err = applyOrientation(s, o)
	return
}
