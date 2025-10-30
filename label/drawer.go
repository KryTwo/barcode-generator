package label

import (
	"bufio"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"main/structs"
	"math"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

func DrawText(s string, cfg structs.Config) {
	//загружаем пользовательский шрифт
	f := LoadFontFromFile(cfg.FontFile)

	//Создаем холст
	const imgW, imgH = 300, 100
	fg, bg := image.Black, image.White
	rgba := image.NewRGBA(image.Rect(0, 0, imgW, imgH))
	draw.Draw(rgba, rgba.Bounds(), bg, image.Point{0, 0}, draw.Src)

	//пишем текст на холсте
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
		Dot: fixed.P(30, 80),
	}
	d.DrawString("hello world")

	//выравнивание по y
	y := 10 + int(math.Ceil(cfg.Size*cfg.DPI/72))
	dy := int(math.Ceil(cfg.Size * cfg.Spacing * cfg.DPI / 72))
	d.Dot = fixed.Point26_6{
		X: (fixed.I(imgW) / 2),
		Y: fixed.I(y),
	}
	d.DrawString("some text")
	y += dy

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
