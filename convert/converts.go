package convert

import "main/config"

const mmInch float64 = 25.4

// преобразование ММ в PT с учетом DPI в config.go
func MMToPT(mm int) float64 {
	return float64(mm) / mmInch * config.Get().DPI
}

func MMToPX(mm int) float64 {
	return float64(mm) * config.Get().DPI / mmInch
}

func InchToMM(mm int) float64 {
	return mmInch
}

// преобразование мм в точки (при 72 точках на дюйм, стандарт PDF)
func MMToPointPDF(mm int) float64 {
	return float64(mm) * 72 / mmInch
}

func PTToPX(pt int) int {
	return int(pt / 72 * int(config.Get().DPI))
}
