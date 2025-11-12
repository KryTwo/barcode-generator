package csvreader

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"slices"
	"strings"
)

// карта соответствий названий колонок
var headerSynonyms = map[string][]string{
	"id":   {"id", "barcode", "шк", "штрихкод", "баркод"},
	"name": {"name", "имя", "название", "наименование"},
}

// возможные разделители
var commaList = []string{",", ";", " ", ".", ":"}

// читаем файл и возвращаем список пар (шк, название)
func Read(data []byte) ([][]string, []string, error) {
	// file, err := os.Open(filePath)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// defer file.Close()

	reader := bytes.NewReader(data)

	//находим первую строку и ищем наиболее часто встречаемый разделитель
	scanner := bufio.NewReader(reader)
	firstLine, err := scanner.ReadString('\n')
	if err != nil {
		panic(err)
	}
	// fmt.Println(firstLine)
	var commaSymbol string
	for _, v := range commaList {
		out := 0
		including := strings.Count(firstLine, v)

		if including > out {
			out = including
			commaSymbol = v
		}

	}

	// //возвращаемся в начало
	// _, err = file.Seek(0, io.SeekStart)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	allData := bytes.NewReader(data)
	csvReader := csv.NewReader(allData)
	csvReader.Comma = rune(commaSymbol[0])

	//читаем заголовок
	header, err := csvReader.Read()
	if err != nil {
		fmt.Println(err)
	}

	//выявляем соответствия по карте
	var newHeader []string
	for _, field := range header {
		lowerField := strings.ToLower(strings.TrimSpace(field))
		found := false
		for key, synonyms := range headerSynonyms {
			if slices.Contains(synonyms, lowerField) {
				newHeader = append(newHeader, key)
			}
		}
		if found {
			break
		}
	}

	// for _, v := range newHeader {
	// 	fmt.Printf("newHeader v: %v\n", v)
	// }

	var records [][]string

	for {
		record, e := csvReader.Read()
		if e != nil && strings.Contains(e.Error(), "wrong number of fields") {
			fmt.Println(e)
			continue
		}
		if e != nil {
			// fmt.Println(e)
			break
		}
		// fmt.Printf("record: %v\n", record)
		records = append(records, record)

	}
	//fmt.Println(records)
	return records, newHeader, nil
}
