package main

import (
	"fmt"
	"github.com/adelolmo/delauncher/config"
	"github.com/adelolmo/delauncher/crypt"
	"github.com/adelolmo/delauncher/magnet"
	"github.com/adelolmo/delauncher/notifications"
	"github.com/andlabs/ui"
	"os"
)

var secretKey = []byte{11, 22, 33, 44, 55, 66, 77, 88, 99, 00, 11, 22, 33, 44, 55, 66}
var key = crypt.Key{
	Value: secretKey,
}

var conf = config.NewConfig(key)

func main() {
	switch len(os.Args) {
	case 1:
		configure()
	case 2:
		link, err := magnet.NewLink(os.Args[1])
		if err != nil {
			notifications.Message(err.Error())
		}
		addMagnet(link)
	default:
		fmt.Print("usage: delauncher [MAGNET_LINK]")
		os.Exit(1)
	}
}

func configure() {
	configProperties, err := conf.Get()
	if err != nil {
		notifications.Message(fmt.Sprintf("Unable to read configuration. Error %s", err.Error()))
		os.Exit(1)
	}
	err = ui.Main(func() {
		serverUrlField := ui.NewEntry()
		serverUrlField.SetText(configProperties.ServerUrl)
		passwordField := ui.NewPasswordEntry()
		passwordField.SetText(configProperties.Password)
		saveButton := ui.NewButton("Save")
		quitButton := ui.NewButton("Quit")

		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel("Deluge server URL:"), false)
		box.Append(serverUrlField, false)
		box.Append(ui.NewLabel("Password:"), false)
		box.Append(passwordField, false)
		buttonsBox := ui.NewHorizontalBox()
		buttonsBox.Append(saveButton, true)
		buttonsBox.Append(quitButton, true)
		box.Append(buttonsBox, false)
		window := ui.NewWindow("Delauncher", 400, 150, false)
		window.SetMargined(true)
		window.SetChild(box)
		saveButton.OnClicked(func(*ui.Button) {
			conf.Save(serverUrlField.Text(), passwordField.Text())
		})
		quitButton.OnClicked(func(*ui.Button) {
			ui.Quit()
		})
		window.OnClosing(func(*ui.Window) bool {
			ui.Quit()
			return true
		})
		window.Show()
	})
	if err != nil {
		panic(err)
	}
}

func addMagnet(magnetLink magnet.Link) {
	configValues, err := conf.Get()
	if err != nil {
		notifications.Message(fmt.Sprintf("Unable to read configuration. Error %s", err.Error()))
		os.Exit(1)
	}
	var server = magnet.Server{Url: configValues.ServerUrl, Password: configValues.Password}
	if err := server.Add(magnetLink); err != nil {
		fmt.Printf(err.Error())
		notifications.Message(err.Error())
		os.Exit(2)
	}
	notifications.Message(fmt.Sprintf("Link added:\n%s", magnetLink.Name))
}
