package convert

import "main/config"

const mmInch float64 = 25.4

// преобразование ММ в PT с учетом DPI в config.go
func MMToPT(mm int) float64 {
	return float64(config.Get().DPI) / mmInch * float64(mm)
}
