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

	dpi := flag.Int("dpi", 300, "screen resolution in Dots Per Inch")
	fontfile := flag.String("fontfile", "./fonts/RobotoforLearning-Black_0.ttf", "filename of the ttf font")
	hinting := flag.String("hinting", "none", "none | full")
	fontSize := flag.Int("size", 20, "font size in points")
	ySpacing := flag.Float64("ySpacing", 30, "spacing btw bc (pt)")
	xSpacing := flag.Float64("xSpacing", 50, "spacing btw bc (pt)")
	wonb := flag.Bool("whiteonblack", false, "white text on a black background")
	hight := flag.Int("height", 30, "set barcode height in mm")
	width := flag.Int("width", 70, "set barcode width in mm")
	margin := flag.Int("margin", 50.0, "(pt) set margin from border list")
	marginToCrop := flag.Int("marginToCrop", 10, "(mm) set margin from left and right borders to crop")

	flag.Parse()

	instance = &structs.Config{
		DPI:          *dpi,
		FontFile:     *fontfile,
		Hinting:      *hinting,
		FontSize:     *fontSize,
		YSpacing:     *ySpacing,
		XSpacing:     *xSpacing,
		WONB:         *wonb,
		Higth:        *hight,
		Width:        *width,
		Margin:       *margin,
		MarginToCrop: *marginToCrop,
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
	cfg.Higth = mm
}

func SetFontSize(size int) {
	cfg := Get()
	cfg.FontSize = size
}

func SetMargin(margin int) {
	cfg := Get()
	cfg.Margin = margin
}

func SetMarginToCrop(marginToCrop int) {
	cfg := Get()
	cfg.MarginToCrop = marginToCrop
}

func SetYSpacing(spacing float64) {
	cfg := Get()
	cfg.YSpacing = float64(spacing)
}

func SetXSpacing(spacing float64) {
	cfg := Get()
	cfg.XSpacing = float64(spacing)
}

// func SetSize()
// func SetDPI()
