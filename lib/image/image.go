package imageprocessor

import (
	"bytes"
	"errors"
	"image"
	"image/color"
	"image/draw"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"

	"github.com/jobtalk/fanlin/lib/error"
	"github.com/nfnt/resize"
	_ "golang.org/x/image/bmp"
)

type Image struct {
	img image.Image
}

func max(v uint, max uint) uint {
	if v > max {
		return max
	}
	return v
}

func EncodeJpeg(img *image.Image) ([]byte, error) {
	if *img == nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, errors.New("img is nil."))
	}

	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, *img, nil)
	return buf.Bytes(), imgproxyerr.New(imgproxyerr.WARNING, err)
}

//DecodeImage is return image.Image
func DecodeImage(bin []byte) (*Image, error) {
	img, _, err := image.Decode(bytes.NewReader(bin))
	return &Image{img}, imgproxyerr.New(imgproxyerr.WARNING, err)
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

func (i *Image) GetImg() *image.Image {
	return &i.img
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
	return EncodeJpeg(img.GetImg())
}
