package main

import (
	"fmt"
	"github.com/adelolmo/delauncher/config"
	"github.com/adelolmo/delauncher/crypt"
	"github.com/adelolmo/delauncher/magnet"
	"github.com/adelolmo/delauncher/notifications"
	"github.com/adelolmo/delugeclient"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
	"strings"
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
	if err != nil { // TODO show error in popup and the close the app
		errorMessage := fmt.Sprintf("Unable to read configuration. Error %s", err.Error())
		fmt.Println(errorMessage)
		notifications.Message(errorMessage)
		os.Exit(1)
	}

	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}
	win.SetTitle("Delauncher")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	win.SetDefaultSize(500, 195)
	win.SetPosition(gtk.WIN_POS_CENTER)
	win.SetIconName("delauncher")

	serverUrlLbl, err := gtk.LabelNew("Server URL")
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	serverUrlEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	textBuffer, err := serverUrlEntry.GetBuffer()
	textBuffer.SetText(configProperties.ServerUrl)
	serverUrlBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		log.Fatal("Unable to create serverUrlBox:", err)
	}
	serverUrlBox.PackStart(serverUrlLbl, false, false, 6)
	serverUrlBox.PackStart(serverUrlEntry, true, true, 6)

	passwordUrlLbl, err := gtk.LabelNew("Password")
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	passwordEntry, err := gtk.EntryNew()
	if err != nil {
		log.Fatal("Unable to create TextView:", err)
	}
	passwordEntry.SetInputPurpose(gtk.INPUT_PURPOSE_PASSWORD)
	passwordEntry.SetVisibility(false)
	passwordBuffer, err := passwordEntry.GetBuffer()
	passwordBuffer.SetText(configProperties.Password)
	passwordUrlBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	if err != nil {
		log.Fatal("Unable to create passwordUrlBox:", err)
	}
	passwordUrlBox.PackStart(passwordUrlLbl, false, false, 6)
	passwordUrlBox.PackStart(passwordEntry, true, true, 6)

	connectionErrorImage, _ := gtk.ImageNew()
	connectionErrorImage.SetFromIconName("delauncher-error", 0)
	connectionSuccessImage, _ := gtk.ImageNew()
	connectionSuccessImage.SetFromIconName("delauncher-success", 0)

	testBtn, err := gtk.ButtonNewWithLabel("Test")
	if err != nil {
		log.Fatal("Unable to create testBtn:", err)
	}
	testBtn.SetSizeRequest(90, 0)

	buttonFirstRowBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	buttonFirstRowBox.SetHAlign(gtk.ALIGN_END)
	buttonFirstRowBox.PackStart(connectionSuccessImage, false, false, 0)
	buttonFirstRowBox.PackStart(connectionErrorImage, false, false, 0)
	buttonFirstRowBox.PackStart(testBtn, false, false, 6)

	exitBtn, err := gtk.ButtonNewWithLabel("Exit")
	if err != nil {
		log.Fatal("Unable to create exitBtn:", err)
	}
	exitBtn.SetSizeRequest(90, 0)
	saveBtn, err := gtk.ButtonNewWithLabel("Save")
	if err != nil {
		log.Fatal("Unable to create saveBtn:", err)
	}
	saveBtn.SetSizeRequest(90, 0)

	buttonSecondRowBox, _ := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 6)
	buttonSecondRowBox.SetHAlign(gtk.ALIGN_END)
	buttonSecondRowBox.SetMarginTop(0)
	buttonSecondRowBox.SetMarginBottom(12)
	buttonSecondRowBox.SetMarginStart(12)
	buttonSecondRowBox.SetMarginEnd(12)
	buttonSecondRowBox.PackStart(exitBtn, false, false, 0)
	buttonSecondRowBox.PackStart(saveBtn, false, false, 0)

	pane, _ := gtk.FrameNew("Deluge configuration")
	pane.SetMarginTop(12)
	pane.SetMarginBottom(6)
	pane.SetMarginStart(12)
	pane.SetMarginEnd(12)

	configFormBox, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	configFormBox.PackStart(serverUrlBox, false, false, 0)
	configFormBox.PackStart(passwordUrlBox, false, false, 0)
	configFormBox.PackStart(buttonFirstRowBox, false, false, 6)

	pane.Add(configFormBox)

	box, _ := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 6)
	box.PackStart(pane, false, false, 0)
	box.PackStart(buttonSecondRowBox, false, false, 0)

	win.Add(box)

	testBtn.Connect("clicked", func() {
		connectionSuccessImage.Hide()
		connectionErrorImage.Hide()

		serverUrl, _ := serverUrlEntry.GetText()
		password, _ := passwordEntry.GetText()
		if len(password) == 0 {
			fmt.Println("ERROR")
			connectionErrorImage.Show()
			return
		}

		deluge := delugeclient.NewDeluge(serverUrl, password)
		if err := deluge.Connect(); err == nil {
			fmt.Println("OK")
			connectionSuccessImage.Show()
			return
		}
		fmt.Println("ERROR")
		connectionErrorImage.Show()
	})
	exitBtn.Connect("clicked", func() {
		gtk.MainQuit()
	})
	saveBtn.Connect("clicked", func() {
		serverUrl, _ := serverUrlEntry.GetText()
		serverUrl = strings.TrimSuffix(serverUrl, "/")
		password, _ := passwordEntry.GetText()
		conf.Save(serverUrl, password)
	})

	win.ShowAll()
	connectionSuccessImage.Hide()
	connectionErrorImage.Hide()
	gtk.Main()
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
