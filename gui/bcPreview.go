package gui

import (
	"main/config"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

type Barcode struct {
	BCImage *canvas.Image
	BCSize  fyne.Size
}

// возвращает контейнер с предпросмотром ШК в заданных характеристиках
func (b *Barcode) MakeBarcodePreviewContainer() fyne.CanvasObject {
	b.BCImage = canvas.NewImageFromImage(nil)
	b.BCSize = fyne.Size{
		Width:  float32(config.Get().Width) / 72 * float32(config.Get().DPI),
		Height: float32(config.Get().Higth) / 72 * float32(config.Get().DPI),
	}
	b.BCImage.SetMinSize(b.BCSize)
	b.BCImage.FillMode = canvas.ImageFillContain
	b.BCImage.ScaleMode = canvas.ImageScaleSmooth
	BCContainer := container.NewStack(b.BCImage)
	BCContainer.Resize(b.BCSize)

	return BCContainer
}
