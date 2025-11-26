package label

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"main/config"
	"math"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func DrawText(s string, img image.Image, bcLenX int) *image.RGBA {
	cfg := config.Get()
	//загружаем пользовательский шрифт
	f := LoadFontFromFile(cfg.FontFile)

	//получаем ширину строки в пикселях
	lenWidth, lenHeight := getTextMeasuresInPixels(f, s)
	fmt.Printf("ширина ШК: %v Px\n", img.Bounds().Max.X)
	fmt.Printf("ширина текста: %v Px, Высота: %v Px\n", math.Round(lenWidth), lenHeight)
	fmt.Println()

	//Создаем холст
	fg, _ := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, bcLenX, cfg.Higth))
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	//заливаем фон под текст
	xLeft := (float64(bcLenX) - lenWidth) / 2
	xRight := (xLeft + lenWidth)

	for y := range lenHeight + lenHeight/4 {
		for x := int(xLeft) - (int(lenWidth) / 20); x < int(xRight+(lenWidth/20)); x++ {
			rgba.Set(x, y, color.CMYK{0, 0, 0, 0})
		}
	}

	// пишем текст на холсте
	h := font.HintingNone
	switch cfg.Hinting {
	case "full":
		h = font.HintingFull
	}

	d := &font.Drawer{
		Dst: rgba,
		Src: fg,
		Face: truetype.NewFace(f, &truetype.Options{
			Size:    float64(cfg.FontSize),
			DPI:     float64(cfg.DPI),
			Hinting: h,
		}),
		Dot: fixed.Point26_6{
			X: fixed.Int26_6(int(math.Ceil(xLeft)) << 6),
			Y: fixed.Int26_6(lenHeight << 6),
		},
	}
	d.DrawString(s)
	MakeFile(rgba, s)

	//выравнивание по y
	/* 	y := 10 + int(math.Ceil(cfg.Size*cfg.DPI/72))
	   	dy := int(math.Ceil(cfg.Size * cfg.Spacing * cfg.DPI / 72))
	   	d.Dot = fixed.Point26_6{
	   		X: (fixed.I(imgW) / 2),
	   		Y: fixed.I(y),
	   	}
	   	d.DrawString("some text")
	   	y += dy */

	return rgba
}
