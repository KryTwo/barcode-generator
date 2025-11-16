package gui

import (
	"fmt"
	"io"
	"log"
	"main/app"

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
	openFileLabel := widget.NewLabel("Выберите файл")
	openFileButton := widget.NewButton("file", func() {
		dlg := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if reader == nil {
				log.Println("cancelled")
				return
			}
			defer reader.Close()

			data, err := io.ReadAll(reader)
			if err != nil {
				dialog.ShowError(fmt.Errorf("File reading error%v\n", err), w)
				return
			}

			result := controller.ProcessFile(data)
			if result.Success == false {
				log.Fatalln("result.ProcessFile error")
			}

			controller.RegeneratePreview()

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
