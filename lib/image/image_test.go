package imageprocessor

import (
	"bytes"
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
	if err != nil {
		t.Fatal(err)
	}
	jpegBin.Seek(0, 0)
	if format := img.GetFormat(); format != "jpeg" {
		t.Fatalf("format is %v, expected jpeg", format)
	}

	var b bytes.Buffer
	if err := EncodeJpeg(&b, img.GetImg(), -1); err != nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()

	if _, err := DecodeImage(confBin); err == nil {
		t.Fatal("no error")
	}
	b.Reset()
}

func BenchmarkEncodeJpeg(b *testing.B) {
	for i := 0; i < b.N; i++ {
		img, err := DecodeImage(jpegBin)
		if err != nil {
			b.Fatal(err)
		}
		jpegBin.Seek(0, 0)
		if format := img.GetFormat(); format != "jpeg" {
			log.Fatalf("format is %v, expected jpeg", format)
		}

		var buf bytes.Buffer
		if err := EncodeJpeg(&buf, img.GetImg(), -1); err != nil {
			log.Fatalf("err is %v.", err)
		}
		buf.Reset()
	}
}

func TestEncodePNG(t *testing.T) {
	img, err := DecodeImage(pngBin)
	if err != nil {
		t.Fatal(err)
	}
	pngBin.Seek(0, 0)
	if format := img.GetFormat(); format != "png" {
		t.Fatalf("format is %v, expected png", format)
	}

	var b bytes.Buffer
	if err := EncodePNG(&b, img.GetImg(), -1); err != nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()

	if _, err := DecodeImage(confBin); err == nil {
		t.Fatal("no error")
	}
	b.Reset()
}

func TestEncodeGIF(t *testing.T) {
	img, err := DecodeImage(gifBin)
	if err != nil {
		t.Fatal(err)
	}
	gifBin.Seek(0, 0)
	if format := img.GetFormat(); format != "gif" {
		t.Fatalf("format is %v, expected png", format)
	}

	var b bytes.Buffer
	if err := EncodeGIF(&b, img.GetImg(), -1); err != nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()

	if _, err = DecodeImage(confBin); err == nil {
		t.Fatal("no error")
	}
	b.Reset()
}

func TestEncodeWebP(t *testing.T) {
	// Lossless
	img, err := DecodeImage(webpLosslessBin)
	if err != nil {
		t.Fatal(err)
	}
	webpLosslessBin.Seek(0, 0)
	if format := img.GetFormat(); format != "webp" {
		t.Fatalf("format is %v, expected webp", format)
	}

	var b bytes.Buffer
	if err := EncodeWebP(&b, img.GetImg(), -1, true); err != nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()

	// Lossy
	img, err = DecodeImage(webpLossyBin)
	if err != nil {
		t.Fatal(err)
	}
	webpLossyBin.Seek(0, 0)
	if format := img.GetFormat(); format != "webp" {
		t.Fatalf("format is %v, expected webp", format)
	}

	if err := EncodeWebP(&b, img.GetImg(), -1, false); err != nil {
		t.Fatalf("err is %v.", err)
	}
	b.Reset()

	// error
	if _, err := DecodeImage(confBin); err == nil {
		t.Fatal("no error")
	}
	b.Reset()
}

func TestEncodeAVIF(t *testing.T) {
	var b bytes.Buffer

	img, err := DecodeImage(pngBin)
	if err != nil {
		t.Fatal(err)
	}
	pngBin.Seek(0, 0)
	if err := EncodeAVIF(&b, img.GetImg(), 50); err != nil {
		t.Fatal(err)
	}
	b.Reset()

	if _, err := DecodeImage(confBin); err == nil {
		t.Fatal("no error")
	}
	b.Reset()
}

func TestDecodeImage(t *testing.T) {
	if _, err := DecodeImage(jpegBin); err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	jpegBin.Seek(0, 0)

	if _, err := DecodeImage(bmpBin); err != nil {
		t.Fatalf("err is not nil. : %v", err)
	}
	bmpBin.Seek(0, 0)

	if _, err := DecodeImage(pngBin); err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	pngBin.Seek(0, 0)

	if _, err := DecodeImage(gifBin); err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	gifBin.Seek(0, 0)

	if _, err := DecodeImage(webpLosslessBin); err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	webpLosslessBin.Seek(0, 0)

	if _, err := DecodeImage(webpLossyBin); err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	webpLossyBin.Seek(0, 0)

	if _, err := DecodeImage(confBin); err == nil {
		t.Log(err)
		t.Fatalf("err is nil")
	}
	confBin.Seek(0, 0)
}
