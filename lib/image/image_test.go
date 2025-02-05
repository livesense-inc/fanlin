package imageprocessor

import (
	"bytes"
	"image/color"
	"log"
	"os"
	"testing"
)

var (
	jpgPath          = "../test/img/Lenna.jpg"
	bmpPath          = "../test/img/Lenna.bmp"
	pngPath          = "../test/img/Lenna.png"
	gifPath          = "../test/img/Lenna.gif"
	webpLosslessPath = "../test/img/Lenna_lossless.webp"
	webpLossyPath    = "../test/img/Lenna_lossy.webp"
	confPath         = "../test/test_conf.json"
)

var (
	jpegBin, _         = os.Open(jpgPath)
	bmpBin, _          = os.Open(bmpPath)
	pngBin, _          = os.Open(pngPath)
	gifBin, _          = os.Open(gifPath)
	webpLosslessBin, _ = os.Open(webpLosslessPath)
	webpLossyBin, _    = os.Open(webpLossyPath)
	confBin, _         = os.Open(confPath)
)

var ResizeImage = resizeImage
var ResizeAndFillImage = resizeAndFillImage
var Crop = crop

func TestEncodeJpeg(t *testing.T) {
	img, _ := DecodeImage(jpegBin, new(bytes.Buffer))
	jpegBin.Seek(0, 0)
	if format := img.GetFormat(); format != "jpeg" {
		t.Fatalf("format is %v, expected jpeg", format)
	}

	var b bytes.Buffer
	err := EncodeJpeg(&b, img.GetImg(), -1)
	if err != nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()

	img, _ = DecodeImage(confBin, new(bytes.Buffer))
	confBin.Seek(0, 0)
	err = EncodeJpeg(&b, img.GetImg(), 50)
	if err == nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()
}

func BenchmarkEncodeJpeg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		img, _ := DecodeImage(jpegBin, new(bytes.Buffer))
		jpegBin.Seek(0, 0)
		if format := img.GetFormat(); format != "jpeg" {
			log.Fatalf("format is %v, expected jpeg", format)
		}

		var b bytes.Buffer
		err := EncodeJpeg(&b, img.GetImg(), -1)
		if err != nil {
			log.Fatalf("err is %v.", err)
		}
		b.Reset()

		img, _ = DecodeImage(confBin, new(bytes.Buffer))
		confBin.Seek(0, 0)
		err = EncodeJpeg(&b, img.GetImg(), 50)
		if err == nil {
			log.Fatalf("err is %v.", err)
		}
		b.Reset()
	}
}

func TestEncodePNG(t *testing.T) {
	img, _ := DecodeImage(pngBin, new(bytes.Buffer))
	pngBin.Seek(0, 0)
	if format := img.GetFormat(); format != "png" {
		t.Fatalf("format is %v, expected png", format)
	}

	var b bytes.Buffer
	err := EncodePNG(&b, img.GetImg(), -1)
	if err != nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()

	img, _ = DecodeImage(confBin, new(bytes.Buffer))
	confBin.Seek(0, 0)
	err = EncodePNG(&b, img.GetImg(), 50)
	if err == nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()
}

func TestEncodeGIF(t *testing.T) {
	img, _ := DecodeImage(gifBin, new(bytes.Buffer))
	gifBin.Seek(0, 0)
	if format := img.GetFormat(); format != "gif" {
		t.Fatalf("format is %v, expected png", format)
	}

	var b bytes.Buffer
	err := EncodeGIF(&b, img.GetImg(), -1)
	if err != nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()

	img, _ = DecodeImage(confBin, new(bytes.Buffer))
	confBin.Seek(0, 0)
	err = EncodeGIF(&b, img.GetImg(), 50)
	if err == nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()
}

func TestEncodeWebP(t *testing.T) {
	// Lossless
	img, _ := DecodeImage(webpLosslessBin, new(bytes.Buffer))
	webpLosslessBin.Seek(0, 0)
	if format := img.GetFormat(); format != "webp" {
		t.Fatalf("format is %v, expected webp", format)
	}

	var b bytes.Buffer
	err := EncodeWebP(&b, img.GetImg(), -1, true)
	if err != nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()

	// Lossy
	img, _ = DecodeImage(webpLossyBin, new(bytes.Buffer))
	webpLossyBin.Seek(0, 0)
	if format := img.GetFormat(); format != "webp" {
		t.Fatalf("format is %v, expected webp", format)
	}

	err = EncodeWebP(&b, img.GetImg(), -1, false)
	if err != nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()

	// error
	img, _ = DecodeImage(confBin, new(bytes.Buffer))
	confBin.Seek(0, 0)
	err = EncodeWebP(&b, img.GetImg(), 50, true)
	if err == nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()
}

func TestEncodeAVIF(t *testing.T) {
	var b bytes.Buffer

	img, _ := DecodeImage(pngBin, new(bytes.Buffer))
	pngBin.Seek(0, 0)
	if err := EncodeAVIF(&b, img.GetImg(), 50); err != nil {
		t.Fatal(err)
	}
	b.Reset()

	img, _ = DecodeImage(confBin, new(bytes.Buffer))
	pngBin.Seek(0, 0)
	if err := EncodeAVIF(&b, img.GetImg(), 50); err == nil {
		t.Fatal("err is nil")
	}
	b.Reset()
}

func TestDecodeImage(t *testing.T) {
	img, err := DecodeImage(jpegBin, new(bytes.Buffer))
	jpegBin.Seek(0, 0)
	if err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	if img == nil {
		t.Fatalf("can not decode.")
	}

	img, err = DecodeImage(bmpBin, new(bytes.Buffer))
	bmpBin.Seek(0, 0)
	if err != nil {
		t.Fatalf("err is not nil. : %v", err)
	}
	if img == nil {
		t.Fatalf("img.%v", img)
	}

	img, err = DecodeImage(pngBin, new(bytes.Buffer))
	pngBin.Seek(0, 0)
	if err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	if img == nil {
		t.Fatalf("can not decode.")
	}

	img, err = DecodeImage(gifBin, new(bytes.Buffer))
	gifBin.Seek(0, 0)
	if err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	if img == nil {
		t.Fatalf("can not decode.")
	}

	img, err = DecodeImage(webpLosslessBin, new(bytes.Buffer))
	webpLosslessBin.Seek(0, 0)
	if err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	if img == nil {
		t.Fatalf("can not decode.")
	}

	img, err = DecodeImage(webpLossyBin, new(bytes.Buffer))
	webpLossyBin.Seek(0, 0)
	if err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	if img == nil {
		t.Fatalf("can not decode.")
	}

	img, err = DecodeImage(confBin, new(bytes.Buffer))
	confBin.Seek(0, 0)
	if err == nil {
		t.Log(err)
		t.Fatalf("err is nil")
	}

	if img == nil {
		t.Fatalf("can not decode.")
	}
}

func TestResizeImage(t *testing.T) {
	img, _ := DecodeImage(confBin, new(bytes.Buffer))
	confBin.Seek(0, 0)

	resizeImg := ResizeImage(*img.GetImg(), 100, 100, 10000, 10000)
	if resizeImg != nil {
		t.Fatalf("value is not nil.")
	}

	img, _ = DecodeImage(jpegBin, new(bytes.Buffer))
	jpegBin.Seek(0, 0)
	ii := *img.GetImg()
	jpegRect := ii.Bounds()

	resizeImg = ResizeImage(*img.GetImg(), uint(jpegRect.Max.X*2), uint(jpegRect.Max.Y*2), 10000, 10000)
	if resizeImg == nil {
		t.Fatalf("value is nil.")
	}

	resizeImg = ResizeImage(*img.GetImg(), uint(jpegRect.Max.X), uint(jpegRect.Max.Y), 10000, 10000)
	if resizeImg == nil {
		t.Fatalf("value is nil.")
	}

	resizeImg = ResizeImage(*img.GetImg(), uint(jpegRect.Max.X*100000), uint(jpegRect.Max.Y*100000), 10000, 10000)
	if resizeImg == nil {
		t.Fatalf("value is nil.")
	}

	resizeImg = ResizeImage(*img.GetImg(), uint(jpegRect.Max.X), uint(jpegRect.Max.Y*100000), 10000, 10000)
	if resizeImg == nil {
		t.Fatalf("value is nil.")
	}

	resizeImg = ResizeImage(*img.GetImg(), uint(jpegRect.Max.X*100000), uint(jpegRect.Max.Y), 10000, 10000)
	if resizeImg == nil {
		t.Fatalf("value is nil.")
	}

	resizeImg = ResizeImage(*img.GetImg(), uint(jpegRect.Max.X), uint(jpegRect.Max.Y+1000), 10000, 10000)
	if resizeImg == nil {
		t.Fatalf("value is nil.")
	}

	resizeImg = ResizeImage(*img.GetImg(), uint(jpegRect.Max.X+1000), uint(jpegRect.Max.Y), 10000, 10000)
	if resizeImg == nil {
		t.Fatalf("value is nil.")
	}

	resizeImg = ResizeImage(*img.GetImg(), uint(jpegRect.Max.X), uint(jpegRect.Max.Y-100), 10000, 10000)
	if resizeImg == nil {
		t.Fatalf("value is nil.")
	}

	resizeImg = ResizeImage(*img.GetImg(), uint(jpegRect.Max.X-100), uint(jpegRect.Max.Y), 10000, 10000)
	if resizeImg == nil {
		t.Fatalf("value is nil.")
	}

	resizeImg = ResizeImage(*img.GetImg(), 0, 0, 10000, 10000)
	if resizeImg == nil {
		t.Fatalf("value is nil.")
	}
}

func TestResizeAndFillImage(t *testing.T) {
	c := color.RGBA{
		R: 0xff,
		G: 0xff,
		B: 0xff,
		A: 0xff,
	}
	img, _ := DecodeImage(confBin, new(bytes.Buffer))
	confBin.Seek(0, 0)

	fillImg := ResizeAndFillImage(*img.GetImg(), 100, 100, c, 10000, 10000)
	if fillImg != nil {
		t.Fatalf("value is not nil.")
	}

	img, _ = DecodeImage(pngBin, new(bytes.Buffer))
	pngBin.Seek(0, 0)

	fillImg = ResizeAndFillImage(*img.GetImg(), 100, 100, c, 10000, 10000)
	if fillImg == nil {
		t.Fatalf("value is nil.")
	}
	if fillImg.Bounds().Max.X != 100 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}
	if fillImg.Bounds().Max.Y != 100 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}

	fillImg = ResizeAndFillImage(*img.GetImg(), 1000, 100, c, 10000, 10000)
	if fillImg == nil {
		t.Fatalf("value is not nil.")
	}
	if fillImg.Bounds().Max.X != 1000 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}
	if fillImg.Bounds().Max.Y != 100 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}

	fillImg = ResizeAndFillImage(*img.GetImg(), 100, 1000, c, 10000, 10000)
	if fillImg == nil {
		t.Fatalf("value is not nil.")
	}
	if fillImg.Bounds().Max.X != 100 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}
	if fillImg.Bounds().Max.Y != 1000 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}

	fillImg = ResizeAndFillImage(*img.GetImg(), 5000, 5000, c, 10000, 10000)
	if fillImg == nil {
		t.Fatalf("value is not nil.")
	}
	if fillImg.Bounds().Max.X != 5000 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}
	if fillImg.Bounds().Max.Y != 5000 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}

	fillImg = ResizeAndFillImage(*img.GetImg(), 0, 0, c, 10000, 10000)
	if fillImg == nil {
		t.Fatalf("value is not nil.")
	}
	if fillImg.Bounds().Max.X != 512 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}
	if fillImg.Bounds().Max.Y != 512 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}

	fillImg = ResizeAndFillImage(*img.GetImg(), 1000000, 1000000, c, 10000, 10000)
	if fillImg == nil {
		t.Fatalf("value is not nil.")
	}
	if fillImg.Bounds().Max.X != 512 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}
	if fillImg.Bounds().Max.Y != 512 {
		t.Fatalf("x is %v.", fillImg.Bounds().Max.X)
	}
}

func TestCrop(t *testing.T) {
	img, _ := DecodeImage(confBin, new(bytes.Buffer))
	confBin.Seek(0, 0)

	cropImg := Crop(*img.GetImg(), 100, 100)
	if cropImg != nil {
		t.Fatalf("value is not nil.")
	}

	img, _ = DecodeImage(pngBin, new(bytes.Buffer))
	pngBin.Seek(0, 0)

	cropImg = Crop(*img.GetImg(), 100, 100)
	if cropImg == nil {
		t.Fatalf("value is nil.")
	}
	if cropImg.Bounds().Max.X != 512 {
		t.Fatalf("x is %v.", cropImg.Bounds().Max.X)
	}
	if cropImg.Bounds().Max.Y != 512 {
		t.Fatalf("y is %v.", cropImg.Bounds().Max.Y)
	}

	cropImg = Crop(*img.GetImg(), 1000, 100)
	if cropImg == nil {
		t.Fatalf("value is not nil.")
	}
	if cropImg.Bounds().Max.X != 512 {
		t.Fatalf("x is %v.", cropImg.Bounds().Max.X)
	}
	if cropImg.Bounds().Max.Y != 51 {
		t.Fatalf("y is %v.", cropImg.Bounds().Max.Y)
	}

	cropImg = Crop(*img.GetImg(), 100, 1000)
	if cropImg == nil {
		t.Fatalf("value is not nil.")
	}
	if cropImg.Bounds().Max.X != 51 {
		t.Fatalf("x is %v.", cropImg.Bounds().Max.X)
	}
	if cropImg.Bounds().Max.Y != 512 {
		t.Fatalf("y is %v.", cropImg.Bounds().Max.Y)
	}

	cropImg = Crop(*img.GetImg(), 5000, 5000)
	if cropImg == nil {
		t.Fatalf("value is not nil.")
	}
	if cropImg.Bounds().Max.X != 512 {
		t.Fatalf("x is %v.", cropImg.Bounds().Max.X)
	}
	if cropImg.Bounds().Max.Y != 512 {
		t.Fatalf("y is %v.", cropImg.Bounds().Max.Y)
	}

	cropImg = Crop(*img.GetImg(), 0, 0)
	if cropImg == nil {
		t.Fatalf("value is not nil.")
	}
	if cropImg.Bounds().Max.X != 512 {
		t.Fatalf("x is %v.", cropImg.Bounds().Max.X)
	}
	if cropImg.Bounds().Max.Y != 512 {
		t.Fatalf("y is %v.", cropImg.Bounds().Max.Y)
	}

	cropImg = Crop(*img.GetImg(), 1000000, 1000000)
	if cropImg == nil {
		t.Fatalf("value is not nil.")
	}
	if cropImg.Bounds().Max.X != 512 {
		t.Fatalf("x is %v.", cropImg.Bounds().Max.X)
	}
	if cropImg.Bounds().Max.Y != 512 {
		t.Fatalf("y is %v.", cropImg.Bounds().Max.Y)
	}
}
