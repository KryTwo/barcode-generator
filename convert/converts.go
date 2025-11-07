package convert

import "main/config"

const mmInch = 25.4

func MMToPT(mm int) float64 {
	return float64(mm) / mmInch * config.Get().DPI
}

func MMToPX(mm int) float64 {
	return float64(mm) * config.Get().DPI / mmInch
}

func InchToMM(mm int) float64 {
	return mmInch
}
