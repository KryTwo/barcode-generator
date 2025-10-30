package label

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"main/structs"
	"math"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func DrawText(s string, cfg structs.Config, img image.Image) {
	//загружаем пользовательский шрифт
	f := LoadFontFromFile(cfg.FontFile)

	//получаем ширину строки в пикселях
	lenWidth, lenHeight := getTextMeasuresInPixels(f, s, cfg)

	//Создаем холст
	const imgW, imgH = 300, 100
	fg, _ := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	//заливаем фон под текст
	xLeft := (imgW - lenWidth) / 2
	xRight := (xLeft + lenWidth)

	for y := range lenHeight + lenHeight/4 {
		for x := int(xLeft) - (int(lenWidth) / 20); x < int(xRight+(lenWidth/20)); x++ {
			// fmt.Printf("x: %v, y: %v\n", x, y)
			// rgba.Set(x, y, color.CMYK{20, 0, 100, 41})
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
			Size:    cfg.Size,
			DPI:     cfg.DPI,
			Hinting: h,
		}),
		Dot: fixed.Point26_6{
			X: fixed.Int26_6(int(math.Ceil(xLeft)) << 6),
			Y: fixed.Int26_6(lenHeight << 6),
		},
	}
	fmt.Println("d.Dot.Y", d.Dot.Y)
	fmt.Printf("dot: %v %v", d.Dot.X, d.Dot.Y)
	d.DrawString(s)

	//выравнивание по y
	//y := 10 + int(math.Ceil(cfg.Size*cfg.DPI/72))
	// dy := int(math.Ceil(cfg.Size * cfg.Spacing * cfg.DPI / 72))
	// d.Dot = fixed.Point26_6{
	// 	X: (fixed.I(imgW) / 2),
	// 	Y: fixed.I(y),
	// }
	// d.DrawString("some text")
	// y += dy

	//создаем файл и записываем в него из буфера
	outFile, err := os.Create("out.png")
	if err != nil {
		fmt.Printf("cant create file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		fmt.Printf("cant encode file: %v\n", err)
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		fmt.Printf("cant flush file: %v\n", err)
		os.Exit(1)
	}

}
