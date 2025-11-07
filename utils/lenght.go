package utils

import (
	"main/config"
)

type Millimeter int

const MMPerInch = 25.4

func (mm Millimeter) Inch() float64 {
	return float64(mm) / MMPerInch
}

func (mm Millimeter) Px() int {
	return int(mm * (Millimeter(config.Get().DPI * MMPerInch)))
}
