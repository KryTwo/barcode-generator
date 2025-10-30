package label

import (
	"fmt"
	"main/structs"
	"os"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/math/fixed"
)

func LoadFontFromFile(font string) *truetype.Font {

	//читаем шрифт
	fontBytes, err := os.ReadFile(font)
	if err != nil {
		fmt.Printf("cant read fontfile with err: %v\n", err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		fmt.Printf("cant parse fontfile with err: %v\n", err)
	}

	return f
}

func getTextMeasuresInPixels(f *truetype.Font, text string, cfg structs.Config) (float64, int) {
	//определение ширины текста в пикселях
	face := truetype.NewFace(f, &truetype.Options{
		Size: cfg.Size,
		DPI:  cfg.DPI,
	})

	var (
		x        fixed.Int26_6 // текущая позиция курсора
		prevRune rune          // предыдущий символ (для кернинга)
	)

	dot := fixed.Point26_6{
		X: x,
		Y: 0,
	}
	var h int
	for i, r := range text {
		// Кернинг между предыдущим и текущим символом
		if i > 0 {
			x += face.Kern(prevRune, r)
		}

		// Получаем метрики символа
		dr, _, _, advance, _ := face.Glyph(dot, r)
		// fmt.Printf("rune measure from face.glyph %v\n", advance)
		// Двигаем курсор
		x += advance

		prevRune = r
		// fmt.Println(dr.Dy())
		if dr.Dy() > h {
			h = dr.Dy()
		}
	}

	widthPixels := float64(x) / 64.0
	heightPixels := h
	// fmt.Printf("ширина текста в пикселях: %v\n", widthPixels)
	// fmt.Printf(" текста в пикселях: %v\n", heightPixels)

	return widthPixels, heightPixels
}
