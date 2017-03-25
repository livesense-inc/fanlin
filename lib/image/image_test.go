package imageprocessor

import (
	"image"
	"image/color"
	"io/ioutil"
	"testing"
)

var (
	jpgPath  = "../test/img/Lenna.jpg"
	bmpPath  = "../test/img/Lenna.bmp"
	pngPath  = "../test/img/Lenna.png"
	gifPath  = "../test/img/Lenna.gif"
	confPath = "../test/test_conf.json"
)

var (
	jpegBin, _ = ioutil.ReadFile(jpgPath)
	bmpBin, _  = ioutil.ReadFile(bmpPath)
	pngBin, _  = ioutil.ReadFile(pngPath)
	gifBin, _  = ioutil.ReadFile(gifPath)
	confBin, _ = ioutil.ReadFile(confPath)
)

var testRect = image.Rect(0, 0, 100, 100)
var ResizeImage = resizeImage
var ResizeAndFillImage = resizeAndFillImage
var Crop = crop

func TestEncodeJpeg(t *testing.T) {
	img, _ := DecodeImage(jpegBin)
	bin, err := EncodeJpeg(img.GetImg(), -1)
	if err != nil {
		t.Fatalf("err is %v.", err)
	}
	if bin == nil {
		t.Fatalf("bin is nil.")
	}

	img, _ = DecodeImage(confBin)
	_, err = EncodeJpeg(img.GetImg(), 50)
	if err == nil {
		t.Fatalf("err is %v.", err)
	}
}

func TestDecodeImage(t *testing.T) {
	img, err := DecodeImage(nil)
	if err == nil {
		t.Fatalf("err is nil")
	}

	img, err = DecodeImage(jpegBin)
	if err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	if img == nil {
		t.Fatalf("can not decode.")
	}

	img, err = DecodeImage(bmpBin)
	if err != nil {
		t.Fatalf("err is not nil. : %v", err)
	}
	if img == nil {
		t.Fatalf("img.%v", img)
	}

	img, err = DecodeImage(pngBin)
	if err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	if img == nil {
		t.Fatalf("can not decode.")
	}

	img, err = DecodeImage(gifBin)
	if err != nil {
		t.Log(err)
		t.Fatalf("err is not nil.")
	}
	if img == nil {
		t.Fatalf("can not decode.")
	}

	img, err = DecodeImage(confBin)
	if err == nil {
		t.Log(err)
		t.Fatalf("err is nil")
	}

	if img == nil {
		t.Fatalf("can not decode.")
	}
}

func TestResizeImage(t *testing.T) {
	img, _ := DecodeImage(confBin)

	resizeImg := ResizeImage(*img.GetImg(), 100, 100, 10000, 10000)
	if resizeImg != nil {
		t.Fatalf("value is not nil.")
	}

	img, _ = DecodeImage(jpegBin)
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
	img, _ := DecodeImage(confBin)

	fillImg := ResizeAndFillImage(*img.GetImg(), 100, 100, c, 10000, 10000)
	if fillImg != nil {
		t.Fatalf("value is not nil.")
	}

	img, _ = DecodeImage(pngBin)

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

func TestCrop(t *testing.T) {
	img, _ := DecodeImage(confBin)

	cropImg := Crop(*img.GetImg(), 100, 100)
	if cropImg != nil {
		t.Fatalf("value is not nil.")
	}

	img, _ = DecodeImage(pngBin)

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
