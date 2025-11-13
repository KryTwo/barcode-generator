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

	dpi := flag.Float64("dpi", 300, "screen resolution in Dots Per Inch")
	fontfile := flag.String("fontfile", "./fonts/RobotoforLearning-Black_0.ttf", "filename of the ttf font")
	hinting := flag.String("hinting", "none", "none | full")
	fontSize := flag.Int("size", 22, "font size in points")
	spacing := flag.Float64("spacing", 1.5, "line spacing (e.g. 2 means double spaced)")
	wonb := flag.Bool("whiteonblack", false, "white text on a black background")
	hight := flag.Int("height", 30, "set barcode height in mm")
	width := flag.Int("width", 70, "set barcode width in mm")
	margin := flag.Float64("margin", 50.0, "(pt) set margin from border list")

	flag.Parse()

	instance = &structs.Config{
		DPI:      *dpi,
		FontFile: *fontfile,
		Hinting:  *hinting,
		FontSize: *fontSize,
		Spacing:  *spacing,
		WONB:     *wonb,
		Hight:    *hight,
		Width:    *width,
		Margin:   *margin,
	}
}

func Get() *structs.Config {
	if instance == nil {
		panic("config not initialized. Call config.Init() first.")
	}
	return instance
}

func SetWidth(mm int) {
	cfg := Get()
	cfg.Width = mm
}

func SetHight(mm int) {
	cfg := Get()
	cfg.Hight = mm
}

func SetFontSize(size int) {
	cfg := Get()
	cfg.FontSize = size
}

// func SetMargin()
// func SetSize()
// func SetDPI()
