package label

import (
	"bufio"
	"fmt"
	"image"
	"image/png"
	"os"
)

// data - название файла
func MakeFile(rgba image.Image, data string) {
	// создаем файл и записываем в него из буфера
	outFile, err := os.Create(data)
	if err != nil {
		fmt.Printf("cant create file: %v\n", err)
		os.Exit(1)
	}
	defer outFile.Close()

	b := bufio.NewWriter(outFile)
	err = png.Encode(b, rgba)
	if err != nil {
		fmt.Printf("cant encode file: %v\n", err)
		os.Exit(1)
	}

	err = b.Flush()
	if err != nil {
		fmt.Printf("cant flush file: %v\n", err)
		os.Exit(1)
	}
}
