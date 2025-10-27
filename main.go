package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
	"strings"
	"time"
)

/*

	func main() {
		a := app.New()
		w := a.NewWindow("Hello World")
		w2 := a.NewWindow("window 2")

		message := widget.NewLabel("Welcome")
		message2 := widget.NewLabel("Welcome Twice")
		button := widget.NewButton("Update", func() {
			formatted := time.Now().Format("Time: 03:04:05")
			message.SetText(formatted)
		})

		w.SetContent(container.NewVBox(message, button))
		w.SetMaster()
		w.Resize(fyne.NewSize(350, 350))
		w.Show()

		buttonNewWindow := widget.NewButton("new Window", func() {
			w3 := a.NewWindow("new window")
			w3.SetContent(message2)
			w3.Show()
		})
		w2.SetContent(container.NewHBox(buttonNewWindow))
		w2.Resize(fyne.NewSize(350, 350))
		w2.Show()

		a.Run()
	}
*/

func main() {
	//карта соответстий возможных названий столбцов к искомому столбцу
	headerSynonyms := map[string][]string{
		"id":   {"id", "barcode", "шк", "штрихкод", "баркод"},
		"name": {"name", "имя", "название", "наименование"},
	}
	commaList := []string{",", ";", " ", ".", ":"}
	file, err := os.Open("code.csv")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer file.Close()

	//находим первую строку и ищем наиболее часто встречаемый разделитель
	scanner := bufio.NewReader(file)
	firstLine, e := scanner.ReadString('\n')
	if err != nil {
		panic(e)
	}
	fmt.Println(firstLine)
	var commaSymbol string
	for _, v := range commaList {
		out := 0
		including := strings.Count(firstLine, v)

		if including > out {
			out = including
			commaSymbol = v
		}

	}
	fmt.Println(commaSymbol)
	//возвращаемся в начало
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(file)

	//читаем заголовки
	var newHeader []string
	header, e := reader.Read()
	if err != nil {
		fmt.Println(e)
	}
	//выявляем соответствия по карте
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
	fmt.Println(newHeader)

	for {
		record, e := reader.Read()
		if e != nil && strings.Contains(e.Error(), "wrong number of fields") {
			fmt.Println(e)
			continue
		}
		if e != nil {
			fmt.Println(e)
			break
		}

		fmt.Println(record)
		time.Sleep(time.Millisecond * 100)
	}

}
