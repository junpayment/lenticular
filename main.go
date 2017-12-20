package main

import (
	"flag"
	"gopkg.in/gographics/imagick.v3/imagick"
	"strconv"
	"fmt"
)

type Lenticular struct {
	img1, img2 *Image // 元画像
	L uint    // lenses per inch
	S uint    // piches per inch
	X float64 // length per pich
	P uint    // ppi
	R uint    // pixels per pich
}

func NewLenticular() (lenticular *Lenticular) {
	lenticular = &Lenticular{}

	lenticular.img1 = NewImage()
	lenticular.img2 = NewImage()

	return
}

func (lenticular *Lenticular) Destroy() {
	lenticular.img1.mw.Destroy()
	lenticular.img2.mw.Destroy()
}

func (lenticular *Lenticular) SetFiles(file1, file2 string) {
	lenticular.img1.ReadImage(file1)
	lenticular.img2.ReadImage(file2)
}

func (lenticular *Lenticular) Set(l, p uint) {
	lenticular.L = l
	lenticular.P = p
	lenticular.S = 2 * lenticular.L
	lenticular.X = 1.0 / float64(lenticular.S)
	lenticular.R = uint(lenticular.X * float64(lenticular.P))

	lenticular.img1.mw.SetImageUnits(imagick.RESOLUTION_PIXELS_PER_INCH)
	lenticular.img1.mw.SetImageResolution(float64(lenticular.P), float64(lenticular.P))
	lenticular.img2.mw.SetImageUnits(imagick.RESOLUTION_PIXELS_PER_INCH)
	lenticular.img2.mw.SetImageResolution(float64(lenticular.P), float64(lenticular.P))
}

type Image struct {
	filename string
	mw *imagick.MagickWand
}

func NewImage() *Image {
	image := &Image{}
	image.mw = imagick.NewMagickWand()

	return image
}

func (img *Image) ReadImage(filename string) {
	img.filename = filename
	img.mw.ReadImage(filename)
}

func (lenticular *Lenticular) Create(dst string) {
	dstMw := lenticular.img1.mw.Clone()

	iterator1 := dstMw.NewPixelIterator()
	iterator2 := lenticular.img2.mw.NewPixelIterator()
	for h := uint(0); h < lenticular.img1.mw.GetImageHeight(); h++ {
		fmt.Printf("%d / %d\n", h + 1, lenticular.img1.mw.GetImageHeight())
		pixels := iterator1.GetNextIteratorRow()
		active := true
		for w := range pixels {
			if active {
				iterator2.SetIteratorRow(int(h))
				pixels[w].SetColor(
					iterator2.GetCurrentIteratorRow()[w].GetColorAsString())
			}
			if 0 == w % int(lenticular.R) {
				active = !active
			}
		}
		iterator1.SyncIterator()
	}

	dstMw.SetImageUnits(imagick.RESOLUTION_PIXELS_PER_INCH)
	dstMw.SetImageResolution(float64(lenticular.P), float64(lenticular.P))

	dstMw.WriteImage(dst)
	dstMw.Destroy()
}

func main() {
	flag.Parse()
	file1, file2, dst :=
		flag.Arg(0), flag.Arg(1), flag.Arg(2)
	l, _ := strconv.Atoi(flag.Arg(3))
	P, _ := strconv.Atoi(flag.Arg(4))

	imagick.Initialize(); defer imagick.Terminate()

	lenticular := NewLenticular()
	lenticular.SetFiles(file1, file2)
	lenticular.Set(uint(l), uint(P))
	lenticular.Create(dst)
}
