package app

import (
	"barcode-app/barcode"
	"barcode-app/config"
	"barcode-app/convert"
	"barcode-app/csvreader"
	"barcode-app/layout"
	"barcode-app/logger"
	"barcode-app/structs"
	"encoding/json"
	"errors"
	"fmt"
	"image"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
)

type Controller struct {
	config             *structs.Config
	CurrentPresetName  string
	CurrentRecords     [][]string
	OnPreviewUpdated   func(*image.RGBA)
	OnPresetChanged    func() //Обновление GUI актуальным списком пресетов
	OnValidationUpdate func()
}

type ProcessResult struct {
	PreviewPNG *image.RGBA
	PreviewBC  *image.RGBA
	Success    bool
	Error      error
}

func NewController(config *structs.Config) *Controller {
	return &Controller{config: config}
}

func (c *Controller) HandleFileSelection(reader fyne.URIReadCloser) error {
	if reader == nil {
		return nil
	}

	defer reader.Close()

	data, err := io.ReadAll(reader)
	if err != nil {
		logger.Log.Error("Не удалось прочитать файл")
		return err
	}

	result := c.ProcessFile(data)
	if result.Success {
		c.RegeneratePreview()
	}
	return nil
}

func (c *Controller) ProcessFile(data []byte) ProcessResult {
	logger.Log.Info("try csvreader.Read")
	records, _, err := csvreader.Read(data)
	logger.Log.Info("success csvreader.Read")

	if err != nil {
		logger.LogError(err, "ошибка processFile")
		return ProcessResult{Error: err}
	}

	c.CurrentRecords = records

	c.RegeneratePreview()

	return ProcessResult{Success: true}
}

func (c *Controller) CropBC(img *image.RGBA) *image.Image {
	// fmt.Printf("img.Bounds(): %v\n", img.Bounds())
	// fmt.Printf("c.config.Margin: %v\n", c.config.Margin)
	// x1 = float64(c.config.Margin)/72*float64(c.config.DPI) - float64(c.config.MarginToCrop)/72*float64(c.config.DPI)
	// y1 = float64(c.config.Margin)/72*float64(c.config.DPI) - 20
	// x2 = x1 + float64(convert.MMToPT(c.config.Width))
	// y2 = y1 + float64(convert.MMToPT(c.config.Higth)) + 40
	var x1, x2, y1, y2 float64

	x1 = convert.MMToPT(c.config.Margin - c.config.MarginToCrop)
	y1 = convert.MMToPT(c.config.Margin - 3)

	x2 = x1 + convert.MMToPT(c.config.Width)
	y2 = y1 + convert.MMToPT(c.config.Height+6)
	croppRect := image.Rect(int(x1), int(y1), int(x2), int(y2))
	croppImg := img.SubImage(croppRect)

	return &croppImg
}

func (c *Controller) SetBCWidth(data string) {
	d, err := strconv.Atoi(data)
	if err != nil {
		log.Fatalf("Failed convert ATOI in SetBCWidth: %v\n", err)
	}
	config.SetWidth(d)
	c.RegeneratePreview()
}

func (c *Controller) SetBCHeight(data string) {
	d, err := strconv.Atoi(data)
	if err != nil {
		log.Fatalf("Failed convert ATOI in SetBCHight: %v\n", err)
	}
	config.SetHeight(d)
	c.RegeneratePreview()
}

func (c *Controller) SetFontSize(data string) {
	d, err := strconv.Atoi(data)
	if err != nil {
		log.Fatalf("Failed convert ATOI in SetFontSize: %v\n", err)
	}
	config.SetFontSize(d)
	c.RegeneratePreview()
}

func (c *Controller) SetTextWrapping(data bool) {
	config.SetTextWrapping(data)
	c.RegeneratePreview()
}

func (c *Controller) SetMargin(data string) {
	d, err := strconv.Atoi(data)
	if err != nil {
		log.Fatalf("Failed convert ATOI in SetMargin: %v\n", err)
	}
	config.SetMargin(d)
	c.RegeneratePreview()
}

func (c *Controller) SetMarginToCrop(data string) {
	d, err := strconv.Atoi(data)
	if err != nil {
		log.Fatalf("Failed convert ATOI in SetMarginToCrop: %v\n", err)
	}
	config.SetMarginToCrop(d)
	c.RegeneratePreview()
}

func (c *Controller) SetYSpacing(data string) {
	d, err := strconv.ParseFloat(data, 64)
	if err != nil {
		log.Fatalf("Failed convert ATOI in SetSpacing: %v\n", err)
	}
	config.SetYSpacing(d)
	c.RegeneratePreview()
}

func (c *Controller) SetXSpacing(data string) {
	d, err := strconv.ParseFloat(data, 64)
	if err != nil {
		log.Fatalf("Failed convert ATOI in SetSpacing: %v\n", err)
	}
	config.SetXSpacing(d)
	c.RegeneratePreview()
}

func (c *Controller) RegeneratePreview() {
	logger.Log.Info("try RegeneratePreview")

	if len(c.CurrentRecords) == 0 {
		return
	}

	logger.Log.Info("try GenerateCode128")

	imgs, err := barcode.GenerateCode128(c.CurrentRecords)

	logger.Log.Info("done GenerateCode128")

	if err != nil {
		log.Fatalf("err: %v\n", err)
		return
	}

	logger.Log.Info("try MakePDF")
	PDFBytes := layout.MakePDF(imgs, c.CurrentRecords, false)
	if len(PDFBytes) == 0 {
		logger.LogError(errors.New("empty pdfbytes"), "cannot get PDFBytes")
		return
	}
	logger.Log.Info("done MakePDF")

	logger.Log.Info("try BytesPdfToPNGConvert")
	img := layout.BytesPdfToPNGConvert(PDFBytes)

	logger.Log.Info("done BytesPdfToPNGConvert")

	if c.OnPreviewUpdated != nil {
		c.OnPreviewUpdated(img)
	}
	logger.Log.Info("done RegeneratePreview")

}

func (c *Controller) SavingFile() {
	if len(c.CurrentRecords) == 0 {
		return
	}

	imgs, err := barcode.GenerateCode128(c.CurrentRecords)
	if err != nil {
		log.Fatalf("err: %v\n", err)
		return
	}

	layout.MakePDF(imgs, c.CurrentRecords, true)
}

func findIndex(s string) int {
	for i := range config.ConfigJSON.Presets {
		if config.ConfigJSON.Presets[i].Name == s {
			return i
		}
	}
	return -1
}

func (c *Controller) SetPreset(s string) {
	idx := findIndex(s)
	if idx == -1 {
		logger.LogError(errors.New("index not find in ConfigJSON"), "cant find index")
		return
	}
	p := config.ConfigJSON.Presets[idx].Setting

	c.SetBCWidth(strconv.Itoa(p.Width))
	c.SetBCHeight(strconv.Itoa(p.Height))
	c.SetFontSize(strconv.Itoa(p.FontSize))
	c.SetMargin(strconv.Itoa(p.Margin))
	c.SetMarginToCrop(strconv.Itoa(p.MarginToCrop))
	c.SetYSpacing(strconv.FormatFloat(p.YSpacing, 'f', -1, 64))
	c.SetXSpacing(strconv.FormatFloat(p.XSpacing, 'f', -1, 64))
	c.SetTextWrapping(p.TextWrapping)
}

func (c *Controller) CreatePreset(s string, onSuccess func()) {
	if strings.ContainsAny(s, `/\:*?"><|`) {
		fmt.Println("некорректные символыы")
		return
	}

	//получение json
	data, err := os.ReadFile("settings.json")
	if err != nil {
		logger.LogError(err, "falied to onen settings.json")

		//окно с ошибкой?
		return
	}
	var jsonConf config.JSONSettings
	json.Unmarshal(data, &jsonConf)

	//Проверка на неизменяемый пресет
	if s == "Стандарт" {
		fmt.Println("Нельзя изменять стандартный пресет")
		//ошибка? или мягкое изменение?
		onSuccess() //заменить на ?onFail?
		return
	}

	var presetNames []string

	for _, v := range jsonConf.Presets {
		presetNames = append(presetNames, v.Name)
	}

	var newPreset config.Preset
	newPreset.Name = s
	newPreset.Setting = *c.config
	//Конфликт имен
	for _, prsn := range presetNames {
		if s == prsn {
			fmt.Println("Такой пресет уже существует, удалите его перед сохранением")
			onSuccess() //заменить на ?onFail?
			return
			//ошибка? или перезапись?
		}
	}

	jsonConf.Presets = append(jsonConf.Presets, newPreset)

	file, err := os.Create("settings.json")

	if err != nil {
		logger.LogError(err, "settings.json creating failed")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonConf); err != nil {
		logger.LogError(err, "fail to encode JSONSettings to settings.json")
	}

	//запись в json
	onSuccess()
}

func (c *Controller) SavePreset(s string) {
	//получение json
	data, err := os.ReadFile("settings.json")
	if err != nil {
		logger.LogError(err, "falied to onen settings.json")
		//окно с ошибкой?
		return
	}
	var jsonConf config.JSONSettings
	json.Unmarshal(data, &jsonConf)

	//Проверка на неизменяемый пресет
	if s == "Стандарт" {
		fmt.Println("Нельзя изменять стандартный пресет")
		//ошибка? или мягкое изменение?
		return
	}

	for i, pn := range jsonConf.Presets {
		if pn.Name == s {
			jsonConf.Presets[i].Name = s
			jsonConf.Presets[i].Setting = *c.config

		}
	}

	file, err := os.Create("settings.json")

	if err != nil {
		logger.LogError(err, "settings.json creating failed")
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(jsonConf); err != nil {
		logger.LogError(err, "fail to encode JSONSettings to settings.json")
	}

	//запись в json
}

func (c *Controller) ReadJSON() error {
	data, err := os.ReadFile("settings.json")
	if err != nil {
		logger.LogError(err, "ошибка чтения json")
		fmt.Printf("ReadJSON err: %v\n", err)
		return nil
	}

	return json.Unmarshal(data, &config.ConfigJSON)
}
