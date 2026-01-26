package gui

import (
	"barcode-app/app"
	"barcode-app/logger"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type openFile struct {
	openFileLabel  *widget.Label
	openFileButton *widget.Button
}

func makeOpenFile(w fyne.Window, с *app.Controller) func() {
	logger.Log.Info("start makeOpenFile")

	tappedFunc := func() {
		logger.Log.Info("попытка открытия файла")

		dlg := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				logger.LogError(err, "ошибка открытия файла")
				dialog.ShowError(err, w)
				return
			}

			logger.Log.Info("проверка на нажатие отмены")
			if reader == nil {
				logger.Log.Info("Выбор файла отменен")
				return
			}
			defer reader.Close()

			logger.Log.Info("попытка чтения файла")
			err = с.HandleFileSelection(reader)
			if err != nil {
				logger.LogError(err, "ошибка чтения файла")
				dialog.ShowError(fmt.Errorf("File reading error%v\n", err), w)
				return
			}
		}, w)
		dlg.SetFilter(storage.NewExtensionFileFilter([]string{".csv"}))
		dlg.Show()
	}

	return tappedFunc
}
