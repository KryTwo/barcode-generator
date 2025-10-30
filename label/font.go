package label

import (
	"fmt"
	"os"

	"github.com/golang/freetype/truetype"
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
