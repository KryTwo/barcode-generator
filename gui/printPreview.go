package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

func MakePrintPreview() *canvas.Image {
	previewImage := canvas.NewImageFromImage(nil)
	previewImage.SetMinSize(fyne.NewSize(600, 600))
	previewImage.FillMode = canvas.ImageFillContain
	previewImage.ScaleMode = canvas.ImageScaleSmooth
	return previewImage
}
