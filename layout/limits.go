package layout

import (
	"fmt"
	"main/config"
	"main/logger"
	"strconv"
)

// размеры листа а4
const (
	a4Width = 210
	a4Higth = 297

	minSpacingFromPageBorder = 10
	maxSpacingFromPageBorder = 100
)

const (
	//ограничение высота баркода
	bcHigthMax = a4Higth - minSpacingFromPageBorder*2
	bcHigthMin = 10

	//ограничение ширины баркода
	bcWidthMax = a4Width - minSpacingFromPageBorder*2
	bcWidthMin = 10

	//размер шрифта
	fontSizeMin = 6
	fontSizeMax = 100

	//размер отступов между баркодами
	xSpacingMin = 20.0
	ySpacingMin = 5.0

	//размер отступов для резки
	marginToCropMin = 0
	marginToCropMax = bcWidthMax
)

func ValidateBCHight(s string) bool {
	p, err := strconv.Atoi(s)
	if err != nil {
		logger.LogError(err, "failed validating bcHight in layout/limits.go")
		return false
	}
	if p > bcHigthMax || p < bcHigthMin {
		return false
	}
	return true
}

func ValidateBCWidth(s string) bool {
	p, err := strconv.Atoi(s)
	if err != nil {
		logger.LogError(err, "failed validating bcWidth in layout/limits.go")
		return false
	}
	if p > bcWidthMax || p < bcWidthMin {
		return false
	}
	return true
}

func ValidateFontSize(s string) bool {
	p, err := strconv.Atoi(s)
	if err != nil {
		logger.LogError(err, "failed validating fontSize in layout/limits.go")
		return false
	}
	if p > fontSizeMax || p < fontSizeMin {
		return false
	}
	return true
}

func ValidateMargin(s string) bool {
	p, err := strconv.Atoi(s)
	if err != nil {
		logger.LogError(err, "failed validating Margin in layout/limits.go")
		return false
	}
	if p < minSpacingFromPageBorder {
		return false
	}
	if p > a4Width-config.Get().Width-config.Get().MarginToCrop*2 {
		return false
	}
	return true
}

func ValidateMarginToCrop(s string) bool {
	p, err := strconv.Atoi(s)
	if err != nil {
		logger.LogError(err, "failed validating MarginToCrop in layout/limits.go")
		return false
	}
	if p < marginToCropMin {
		return false
	}
	if p > config.Get().Margin {
		return false
	}
	if p*2 > a4Width-minSpacingFromPageBorder*2-config.Get().Width {
		fmt.Printf("p margin: %v\n", p)
		fmt.Printf("a4Width: %v\n", a4Width)
		fmt.Printf("config.Get().Width: %v\n", config.Get().Width)
		fmt.Printf("minSpacingFromPageBorder: %v\n", minSpacingFromPageBorder)
		return false
	}
	return true
}

func ValidateXSpacing(s string) bool {
	p, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.LogError(err, "failed validating XSpacing in layout/limits.go")
		return false
	}
	if p < xSpacingMin {
		return false
	}
	if p > float64(config.Get().Width) {
		return false
	}
	return true
}

func ValidateYSpacing(s string) bool {
	p, err := strconv.ParseFloat(s, 64)
	if err != nil {
		logger.LogError(err, "failed validating YSpacing in layout/limits.go")
		return false
	}
	if p < ySpacingMin {
		return false
	}
	if p > float64(config.Get().Higth)*2 {
		return false
	}
	return true
}
