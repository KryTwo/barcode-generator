package gui

import (
	"barcode-app/app"
	"barcode-app/updater"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func makeRecentFilesChildMenu(c *app.Controller) []*fyne.MenuItem {
	var res []*fyne.MenuItem

	c.GetRecentFiles()
	for _, v := range c.GetRecentFiles() {
		new := fyne.NewMenuItem(v, func() {
			c.OpenFileByPath(v)
		})
		res = append(res, new)
	}
	return res
}

func setupMenu(a fyne.App, w fyne.Window, c *app.Controller) {
	renderMenu := func() {
		mainMenu := makeMenu(a, w, c)
		w.SetMainMenu(mainMenu)
		w.MainMenu().Refresh()
	}

	c.OnRecentFilesChanged = func() {
		renderMenu()
	}
	renderMenu()
}

func makeMenu(a fyne.App, w fyne.Window, c *app.Controller) *fyne.MainMenu {
	c.CheckFileExist()
	fileOpenFile := fyne.NewMenuItem("Open File", makeOpenFile(w, c))

	fileOpenRecentFile := fyne.NewMenuItem("Open Recent", nil)
	//fyne.NewMenuItem("last opened", nil), fyne.NewMenuItem("first opened", nil)
	fileOpenRecentFile.ChildMenu = fyne.NewMenu("", makeRecentFilesChildMenu(c)...)

	// settingsItem := fyne.NewMenuItem("TODO settings", func() {
	// 	fmt.Println("settings print")
	// })

	showAbout := func() {
		w := a.NewWindow("About")
		w.CenterOnScreen()
		w.SetContent(widget.NewLabelWithStyle(aboutContent, fyne.TextAlignCenter, fyne.TextStyle{}))
		w.Show()
	}

	//проверка новой версии. Одной кнопкой? Вывести текущую версию.
	checkUpdates := func() {
		w := a.NewWindow("Version")
		w.CenterOnScreen()
		w.Resize(fyne.Size{Width: 200, Height: 100})

		currentVersion := a.Metadata().Version
		releaseInfo, _ := updater.GetReleaseInfo()

		content := container.NewVBox(
			widget.NewLabel("Current vesion: "+currentVersion),
			widget.NewLabel("Available version: "+releaseInfo.TagName),
			widget.NewButton("Скачать", func() {
				u, err := url.Parse(releaseInfo.HTMLURL)
				if err != nil {
					return
				}

				fyne.CurrentApp().OpenURL(u)
			}),
		)

		w.SetContent(container.NewCenter(content))
		w.Show()
	}

	makeFeedback := func() {
		w := a.NewWindow("Feedback")
		w.CenterOnScreen()
		w.Resize(fyne.Size{Width: 500, Height: 400})

		entryMsg := widget.NewMultiLineEntry()
		entryMsg.SetPlaceHolder("Опишите проблему здесь")
		sendButton := widget.NewButton("Отправить", func() {
			msg := entryMsg.Text
			if msg == "" {
				return
			}

			entryMsg.Disable()

			go func() {
				err := updater.SendFeedback(msg)

				if err != nil {
					dialog.ShowError(err, w)
					entryMsg.Enable()
				} else {
					d := dialog.NewInformation("Успех", "Сообщение отправлено!", w)
					d.SetOnClosed(func() {
						w.Close()
					})
					d.Show()
				}
			}()
		})

		sendButton.Importance = widget.HighImportance

		content := container.NewBorder(
			widget.NewLabel("Ваш отзыв помогает нам стать лучше"),
			sendButton,
			nil,
			nil,
			entryMsg,
		)
		w.SetContent(container.NewPadded(content))
		w.Show()
	}
	helpAbout := fyne.NewMenuItem("About", showAbout)
	helpCheckUpdates := fyne.NewMenuItem("Check updates", checkUpdates)
	feedback := fyne.NewMenuItem("Send Feedback", makeFeedback)
	//help_help := fyne.NewMenuItem("About", showAbout)

	File := fyne.NewMenu("File", fileOpenFile, fileOpenRecentFile)
	//Settings := fyne.NewMenu("Settings", settingsItem)
	help := fyne.NewMenu("Help", helpAbout, helpCheckUpdates, feedback)

	main_menu := fyne.NewMainMenu(File, help)
	return main_menu

}

var aboutContent string = "This program is designed to simplify the barcode printing process\n" +
	"and reduce paper consumption.\n\n" +
	"Specially for Vi.ru\n\n" +
	"Distributed free of charge."
