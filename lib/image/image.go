package imageprocessor

import (
	"bytes"
	_ "embed"
	"errors"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"runtime"
	"sync"

	"github.com/Kagami/go-avif"
	"github.com/chai2010/webp"
	"github.com/disintegration/gift"
	"github.com/ieee0824/libcmyk"
	imgproxyerr "github.com/livesense-inc/fanlin/lib/error"
	"github.com/livesense-inc/go-lcms/lcms"
	"github.com/rwcarlsen/goexif/exif"
	_ "github.com/strukturag/libheif-go"
	_ "golang.org/x/image/bmp"
)

var affines = map[int]gift.Filter{
	2: gift.FlipHorizontal(),
	3: gift.Rotate180(),
	4: gift.FlipVertical(),
	5: gift.Transpose(),
	6: gift.Rotate270(),
	7: gift.Transverse(),
	8: gift.Rotate90(),
}

var mlConverterCache = &sync.Map{}

//go:embed default.icc
var defaultICCProfile []byte

var cmykToRGBTransformer *lcms.Transform

type Image struct {
	img         image.Image
	format      string
	orientation int
	fillColor   color.Color
	outerBounds image.Rectangle
	filter      *gift.GIFT
}

func (i *Image) ConvertColor(networkPath string) error {
	sc := i.img.At(0, 0)
	_, ok := sc.(color.CMYK)
	if !ok {
		return nil
	}

	rect := i.img.Bounds()
	ret := image.NewRGBA(rect)

	var converter *libcmyk.Converter
	iface, ok := mlConverterCache.Load(networkPath)
	if !ok {
		cr, err := libcmyk.New(networkPath)
		if err != nil {
			return err
		}
		mlConverterCache.Store(networkPath, cr)
		converter = cr
	} else {
		converter = iface.(*libcmyk.Converter)
	}

	w := rect.Max.X
	h := rect.Max.Y

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			cmyk := i.img.At(x, y).(color.CMYK)
			rgba, err := converter.CMYK2RGBA(&cmyk)
			if err != nil {
				return err
			}
			ret.Set(x, y, rgba)
		}
	}
	i.img = ret
	return nil
}

func (i *Image) ConvertColorWithICCProfile() {
	switch src := i.img.(type) {
	case *image.CMYK:
		if cmykToRGBTransformer == nil {
			return
		}
		dst := image.NewRGBA(i.img.Bounds())
		cmykToRGBTransformer.DoTransform(src.Pix, dst.Pix, len(src.Pix)/4)
		for i := range dst.Pix {
			if (i+1)%4 == 0 {
				dst.Pix[i] = 255 // Alpha
			}
		}
		i.img = dst
	}
}

func EncodeJpeg(buf io.Writer, img *image.Image, q int) error {
	if *img == nil {
		return imgproxyerr.New(imgproxyerr.WARNING, errors.New("img is nil"))
	}

	if !(0 <= q && q <= 100) {
		q = jpeg.DefaultQuality
	}

	err := jpeg.Encode(buf, *img, &jpeg.Options{Quality: q})
	return imgproxyerr.New(imgproxyerr.WARNING, err)
}

func EncodePNG(buf io.Writer, img *image.Image, q int) error {
	if *img == nil {
		return imgproxyerr.New(imgproxyerr.WARNING, errors.New("img is nil"))
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

	err := e.Encode(buf, *img)
	return imgproxyerr.New(imgproxyerr.WARNING, err)
}

func EncodeGIF(buf io.Writer, img *image.Image, q int) error {
	if *img == nil {
		return imgproxyerr.New(imgproxyerr.WARNING, errors.New("img is nil"))
	}

	// GIF is not support quality

	err := gif.Encode(buf, *img, &gif.Options{})
	return imgproxyerr.New(imgproxyerr.WARNING, err)
}

func EncodeWebP(buf io.Writer, img *image.Image, q int, lossless bool) error {
	if *img == nil {
		return imgproxyerr.New(imgproxyerr.WARNING, errors.New("img is nil"))
	}
	if !(0 <= q && q < 100) {
		// webp.DefaulQuality = 90 is large, adjust to JPEG
		q = jpeg.DefaultQuality
	}

	var option webp.Options
	if lossless {
		option.Lossless = true
	} else {
		option.Lossless = false
		option.Quality = float32(q)
	}

	err := webp.Encode(buf, *img, &option)
	return imgproxyerr.New(imgproxyerr.WARNING, err)
}

func EncodeAVIF(buf io.Writer, img *image.Image, q int) error {
	if *img == nil {
		return imgproxyerr.New(imgproxyerr.WARNING, errors.New("img is nil"))
	}

	// https://pkg.go.dev/github.com/Kagami/go-avif
	if q < 0 {
		// not specified
		q = avif.MinQuality + 20
	} else if q < avif.MinQuality {
		q = avif.MinQuality
	} else if q > avif.MaxQuality {
		q = avif.MaxQuality
	}
	q = avif.MaxQuality - q // lower is better, invert

	opts := avif.Options{
		Threads:        0,             // all available cores
		Speed:          avif.MaxSpeed, // bigger is faster, but lower compress ratio
		Quality:        q,             // lower is better, zero is lossless
		SubsampleRatio: nil,           // 4:2:0
	}
	if err := avif.Encode(buf, *img, &opts); err != nil {
		return imgproxyerr.New(imgproxyerr.WARNING, err)
	}

	return nil
}

func DecodeImage(r io.Reader) (*Image, error) {
	img, format, orientation, err := decode(r)
	if err != nil {
		return nil, imgproxyerr.New(imgproxyerr.WARNING, err)
	}
	return &Image{
		img:         img,
		format:      format,
		orientation: orientation,
		outerBounds: img.Bounds(),
		filter:      gift.New(),
	}, nil
}

func (i *Image) Process() {
	bounds := i.filter.Bounds(i.img.Bounds())
	dest := image.NewRGBA(bounds)
	i.filter.Draw(dest, i.img)

	if dest.Bounds() == i.outerBounds || i.fillColor == nil {
		i.img = dest
		return
	}

	bg := image.NewRGBA(i.outerBounds)
	draw.Draw(bg, bg.Bounds(), &image.Uniform{i.fillColor}, image.Point{}, draw.Src)
	centerH := math.Abs(float64(i.outerBounds.Max.Y-bounds.Max.Y)) / 2.0
	centerW := math.Abs(float64(i.outerBounds.Max.X-bounds.Max.X)) / 2.0
	center := bounds.Min.Sub(image.Pt(int(centerW), int(centerH)))
	draw.Draw(bg, bg.Bounds(), dest, center, draw.Over)
	i.img = bg
}

func (i *Image) ResizeAndFill(w, h uint, c color.Color) {
	if w == 0 || h == 0 {
		return
	}
	innerW := w
	innerH := h
	r := i.img.Bounds()
	if int(innerW)*r.Max.Y < int(innerH)*r.Max.X {
		innerH = 0
	} else {
		innerW = 0
	}
	i.filter.Add(gift.Resize(int(innerW), int(innerH), gift.LanczosResampling))
	i.outerBounds = image.Rect(0, 0, int(w), int(h))
	i.fillColor = c
}

func (i *Image) Crop(w, h uint) {
	if w == 0 || h == 0 {
		return
	}
	i.filter.Add(gift.ResizeToFill(int(w), int(h), gift.LanczosResampling, gift.CenterAnchor))
	i.outerBounds = image.Rect(0, 0, int(w), int(h))
}

func (i *Image) Invert() {
	i.filter.Add(gift.Invert())
}

func (i *Image) Grayscale() {
	i.filter.Add(gift.Grayscale())
}

func (i *Image) ApplyOrientation() {
	if affine, ok := affines[i.orientation]; ok {
		i.filter.Add(affine)
	}
}

func (i *Image) GetImg() *image.Image {
	if i.img == nil {
		return nil
	}
	return &i.img
}

func (i *Image) GetFormat() string {
	return i.format
}

func Set404Image(buf io.Writer, data io.Reader, w uint, h uint, c color.Color) error {
	img, err := DecodeImage(data)
	if err != nil {
		return imgproxyerr.New(imgproxyerr.ERROR, err)
	}
	img.ResizeAndFill(w, h, c)
	img.Process()
	return EncodeJpeg(buf, img.GetImg(), jpeg.DefaultQuality)
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

func decode(r io.Reader) (d image.Image, format string, o int, err error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	d, format, err = image.Decode(tee)
	if err != nil {
		return
	}

	raw := buf.Bytes()
	o, _ = readOrientation(bytes.NewReader(raw))
	return
}

func SetUpColorConverter() error {
	srcProf, err := lcms.OpenProfileFromMem(defaultICCProfile)
	if err != nil {
		return err
	}
	defer srcProf.CloseProfile()

	dstProf, err := lcms.CreateSRGBProfile()
	if err != nil {
		return err
	}
	defer dstProf.CloseProfile()

	t, err := lcms.CreateTransform(srcProf, lcms.TYPE_CMYK_8, dstProf, lcms.TYPE_RGBA_8)
	if err != nil {
		return err
	}
	cmykToRGBTransformer = t
	runtime.SetFinalizer(cmykToRGBTransformer, func(t *lcms.Transform) {
		t.DeleteTransform()
	})

	return nil
}
