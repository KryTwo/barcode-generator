package gui

import (
	"fmt"
	"io"
	"log"
	"main/app"
	"main/logger"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type openFile struct {
	openFileLabel  *widget.Label
	openFileButton *widget.Button
}

func makeOpenFile(w fyne.Window, controller *app.Controller) openFile {
	logger.Log.Info("start makeOpenFile")
	openFileLabel := widget.NewLabel("Выберите файл")
	openFileButton := widget.NewButton("file", func() {
		logger.Log.Info("попытка открытия файла")

		dlg := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				logger.LogError(err, "ошибка открытия файла")
				dialog.ShowError(err, w)
				return
			}

			logger.Log.Info("проверка на нажатие отмены")
			if reader == nil {
				log.Println("cancelled")
				return
			}
			defer reader.Close()

			logger.Log.Info("попытка чтения файла")
			data, err := io.ReadAll(reader)
			if err != nil {
				logger.LogError(err, "ошибка чтения файла")
				dialog.ShowError(fmt.Errorf("File reading error%v\n", err), w)
				return
			}

			logger.Log.Info("try ProcessFile")
			result := controller.ProcessFile(data)
			if result.Success == false {
				logger.Log.Info("failed ProcessFile")
				log.Fatalln("result.ProcessFile error")
			}

			logger.Log.Info("try RegeneradePreview")
			controller.RegeneratePreview()
			logger.Log.Info("success RegeneradePreview")

			if err != nil {
				dialog.ShowError(err, w)
				return
			}
		}, w)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		dlg.Show()
	})
	return openFile{
		openFileLabel:  openFileLabel,
		openFileButton: openFileButton,
	}
}
