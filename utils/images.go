package utils

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// Bilinear resize
func resizeBilinear(src image.Image, maxWidth, maxHeight int) image.Image {
	srcBounds := src.Bounds()
	srcW := srcBounds.Dx()
	srcH := srcBounds.Dy()

	ratioW := float64(maxWidth) / float64(srcW)
	ratioH := float64(maxHeight) / float64(srcH)
	ratio := ratioW
	if ratioH < ratioW {
		ratio = ratioH
	}

	if ratio >= 1 {
		return src
	}

	newW := int(float64(srcW) * ratio)
	newH := int(float64(srcH) * ratio)
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	for y := 0; y < newH; y++ {
		for x := 0; x < newW; x++ {
			fx := float64(x) / ratio
			fy := float64(y) / ratio
			x1 := int(fx)
			y1 := int(fy)
			x2 := x1 + 1
			y2 := y1 + 1

			if x2 >= srcW {
				x2 = srcW - 1
			}
			if y2 >= srcH {
				y2 = srcH - 1
			}

			c11 := src.At(x1, y1)
			c12 := src.At(x1, y2)
			c21 := src.At(x2, y1)
			c22 := src.At(x2, y2)

			r := bilinearInterpolate(getR(c11), getR(c21), getR(c12), getR(c22), fx-float64(x1), fy-float64(y1))
			g := bilinearInterpolate(getG(c11), getG(c21), getG(c12), getG(c22), fx-float64(x1), fy-float64(y1))
			b := bilinearInterpolate(getB(c11), getB(c21), getB(c12), getB(c22), fx-float64(x1), fy-float64(y1))
			a := bilinearInterpolate(getA(c11), getA(c21), getA(c12), getA(c22), fx-float64(x1), fy-float64(y1))

			dst.Set(x, y, color.RGBA{r, g, b, a})
		}
	}

	return dst
}

func bilinearInterpolate(c00, c10, c01, c11 uint8, dx, dy float64) uint8 {
	top := float64(c00)*(1-dx) + float64(c10)*dx
	bottom := float64(c01)*(1-dx) + float64(c11)*dx
	return uint8(top*(1-dy) + bottom*dy)
}

func getR(c color.Color) uint8 { r, _, _, _ := c.RGBA(); return uint8(r >> 8) }
func getG(c color.Color) uint8 { _, g, _, _ := c.RGBA(); return uint8(g >> 8) }
func getB(c color.Color) uint8 { _, _, b, _ := c.RGBA(); return uint8(b >> 8) }
func getA(c color.Color) uint8 { _, _, _, a := c.RGBA(); return uint8(a >> 8) }

// Watermark teks besar, semi-transparan putih
func addWatermarkold(img *image.RGBA, text string) {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	face := basicfont.Face7x13

	// Hitung ukuran teks
	textWidth := font.MeasureString(face, text).Round()
	textHeight := face.Metrics().Height.Ceil()

	// Posisi center
	x := (w - textWidth) / 2
	y := (h + textHeight) / 2

	// Warna putih semi transparan
	col := image.NewUniform(color.RGBA{255, 255, 255, 255})

	d := &font.Drawer{
		Dst:  img,
		Src:  col,
		Face: face,
		Dot: fixed.Point26_6{
			X: fixed.I(x),
			Y: fixed.I(y),
		},
	}

	d.DrawString(text)
}

func addWatermarkTTF(img *image.RGBA, text string) error {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	// Load font file
	fontBytes, err := os.ReadFile("assets/arialbd.ttf")
	if err != nil {
		return err
	}

	ft, err := opentype.Parse(fontBytes)
	if err != nil {
		return err
	}

	// Ukuran font (dinamis berdasarkan lebar gambar)
	fontSize := float64(w) / 30

	face, err := opentype.NewFace(ft, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return err
	}

	// Warna putih semi transparan
	col := image.NewUniform(color.RGBA{255, 255, 255, 255})

	d := &font.Drawer{
		Dst:  img,
		Src:  col,
		Face: face,
	}

	// Hitung posisi center
	textWidth := d.MeasureString(text).Round()
	textHeight := face.Metrics().Height.Ceil()

	x := (w - textWidth) / 2
	y := (h + textHeight) / 2

	d.Dot = fixed.Point26_6{
		X: fixed.I(x),
		Y: fixed.I(y),
	}

	d.DrawString(text)

	return nil
}

// Simpan base64 -> resize -> watermark
func SaveResizeBase64ToFile(base64str string, fileName string) error {
	if strings.Contains(base64str, ",") {
		parts := strings.Split(base64str, ",")
		base64str = parts[1]
	}

	data, err := base64.StdEncoding.DecodeString(base64str)
	if err != nil {
		return err
	}

	img, format, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return err
	}

	resized := resizeBilinear(img, 800, 600)
	rgba := image.NewRGBA(resized.Bounds())
	draw.Draw(rgba, resized.Bounds(), resized, image.Point{}, draw.Src)

	// watermark: yyyymmddhis
	now := time.Now().Format("2006-01-02 15:04:05")
	addWatermarkTTF(rgba, now)

	f, err := os.Create("./uploads/" + fileName)
	if err != nil {
		return err
	}
	defer f.Close()

	if format == "png" {
		return png.Encode(f, rgba)
	} else {
		opts := jpeg.Options{Quality: 90}
		return jpeg.Encode(f, rgba, &opts)
	}
}
