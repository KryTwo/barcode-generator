package config

import (
	"flag"
	"main/structs"
)

var instance *structs.Config

func Init() {
	if instance != nil {
		return
	}

	dpi := flag.Float64("dpi", 48, "screen resolution in Dots Per Inch")
	fontfile := flag.String("fontfile", "./fonts/RobotoforLearning-Black_0.ttf", "filename of the ttf font")
	hinting := flag.String("hinting", "none", "none | full")
	size := flag.Float64("size", 32, "font size in points")
	spacing := flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb := flag.Bool("whiteonblack", false, "white text on a black background")
	height := flag.Int("height", 30, "set barcode height in mm")
	width := flag.Int("width", 70, "set barcode width in mm")

	flag.Parse()

	instance = &structs.Config{
		DPI:      *dpi,
		FontFile: *fontfile,
		Hinting:  *hinting,
		Size:     *size,
		Spacing:  *spacing,
		WONB:     *wonb,
		Height:   *height,
		Width:    *width,
	}
}

func Get() *structs.Config {
	if instance == nil {
		panic("config not initialized. Call config.Init() first.")
	}
	return instance
}
