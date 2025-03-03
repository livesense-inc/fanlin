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

func TestEncodeJpeg(t *testing.T) {
	img, err := DecodeImage(jpegBin)
	jpegBin.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if format := img.GetFormat(); format != "jpeg" {
		t.Fatalf("format is %v, expected jpeg", format)
	}

	var b bytes.Buffer
	if err := EncodeJpeg(&b, img.GetImg(), -1); err != nil {
		t.Fatal(err)
	}
	b.Reset()

	defer confBin.Seek(0, 0)
	if _, err := DecodeImage(confBin); err == nil {
		t.Error("no error")
	}
	b.Reset()
}

func BenchmarkEncodeJpeg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		img, err := DecodeImage(jpegBin)
		jpegBin.Seek(0, 0)
		if err != nil {
			b.Fatal(err)
		}
		if format := img.GetFormat(); format != "jpeg" {
			log.Fatalf("format is %v, expected jpeg", format)
		}

		var buf bytes.Buffer
		if err := EncodeJpeg(&buf, img.GetImg(), -1); err != nil {
			log.Fatal(err)
		}
		buf.Reset()
	}
}

func TestEncodePNG(t *testing.T) {
	img, err := DecodeImage(pngBin)
	pngBin.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if format := img.GetFormat(); format != "png" {
		t.Fatalf("format is %v, expected png", format)
	}

	var b bytes.Buffer
	if err := EncodePNG(&b, img.GetImg(), -1); err != nil {
		t.Fatal(err)
	}
	b.Reset()

	defer confBin.Seek(0, 0)
	if _, err := DecodeImage(confBin); err == nil {
		t.Error("no error")
	}
	b.Reset()
}

func TestEncodeGIF(t *testing.T) {
	img, err := DecodeImage(gifBin)
	gifBin.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if format := img.GetFormat(); format != "gif" {
		t.Fatalf("format is %v, expected png", format)
	}

	var b bytes.Buffer
	if err := EncodeGIF(&b, img.GetImg(), -1); err != nil {
		t.Fatal(err)
	}
	b.Reset()

	defer confBin.Seek(0, 0)
	if _, err = DecodeImage(confBin); err == nil {
		t.Error("no error")
	}
	b.Reset()
}

func TestEncodeWebP(t *testing.T) {
	// Lossless
	img, err := DecodeImage(webpLosslessBin)
	webpLosslessBin.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if format := img.GetFormat(); format != "webp" {
		t.Fatalf("format is %v, expected webp", format)
	}

	var b bytes.Buffer
	if err := EncodeWebP(&b, img.GetImg(), -1, true); err != nil {
		t.Fatal(err)
	}
	b.Reset()

	// Lossy
	img, err = DecodeImage(webpLossyBin)
	webpLossyBin.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if format := img.GetFormat(); format != "webp" {
		t.Fatalf("format is %v, expected webp", format)
	}

	if err := EncodeWebP(&b, img.GetImg(), -1, false); err != nil {
		t.Fatal(err)
	}
	b.Reset()

	// error
	defer confBin.Seek(0, 0)
	if _, err := DecodeImage(confBin); err == nil {
		t.Error("no error")
	}
	b.Reset()
}

func TestEncodeAVIF(t *testing.T) {
	var b bytes.Buffer

	img, err := DecodeImage(pngBin)
	pngBin.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}
	if err := EncodeAVIF(&b, img.GetImg(), 50); err != nil {
		t.Fatal(err)
	}
	b.Reset()

	defer confBin.Seek(0, 0)
	if _, err := DecodeImage(confBin); err == nil {
		t.Error("no error")
	}
	b.Reset()
}

func TestDecodeImage(t *testing.T) {
	if _, err := DecodeImage(jpegBin); err != nil {
		t.Error(err)
	}
	defer jpegBin.Seek(0, 0)

	if _, err := DecodeImage(bmpBin); err != nil {
		t.Error(err)
	}
	defer bmpBin.Seek(0, 0)

	if _, err := DecodeImage(pngBin); err != nil {
		t.Error(err)
	}
	defer pngBin.Seek(0, 0)

	if _, err := DecodeImage(gifBin); err != nil {
		t.Error(err)
	}
	defer gifBin.Seek(0, 0)

	if _, err := DecodeImage(webpLosslessBin); err != nil {
		t.Error(err)
	}
	defer webpLosslessBin.Seek(0, 0)

	if _, err := DecodeImage(webpLossyBin); err != nil {
		t.Error(err)
	}
	defer webpLossyBin.Seek(0, 0)

	if _, err := DecodeImage(confBin); err == nil {
		t.Error("err is nil")
	}
	defer confBin.Seek(0, 0)
}

func TestImageProcess(t *testing.T) {
	img, err := DecodeImage(jpegBin)
	jpegBin.Seek(0, 0)
	if err != nil {
		t.Fatal(err)
	}

	img.ApplyOrientation()
	img.ResizeAndFill(1618, 1000, color.RGBA{uint8(32), uint8(32), uint8(32), 0xff})
	img.Process()

	innerP := img.GetImg()
	inner := *innerP
	if inner.Bounds().Max.X != 1618 {
		t.Errorf("want=%d, got=%d", 1618, inner.Bounds().Max.X)
	}
	if inner.Bounds().Max.Y != 1000 {
		t.Errorf("want=%d, got=%d", 1000, inner.Bounds().Max.Y)
	}
}
