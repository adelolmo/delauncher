package main

import (
	"fmt"
	"github.com/adelolmo/delauncher/config"
	"github.com/adelolmo/delauncher/crypt"
	"github.com/adelolmo/delauncher/magnet"
	"github.com/adelolmo/delauncher/notifications"
	"github.com/adelolmo/delugeclient"
	"github.com/leaanthony/mewn"
	"github.com/webview/webview"
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
	assets := mewn.Group("./assets/")
	configProperties, err := conf.Get()
	if err != nil {
		notifications.Message(fmt.Sprintf("Unable to read configuration. Error %s", err.Error()))
		os.Exit(1)
	}
	w := webview.New(true)
	defer w.Destroy()
	w.SetTitle("Delauncher")
	w.SetSize(500, 195, webview.HintNone)
	w.Bind("testConnection", func(serverUrl, password string) bool {
		deluge := delugeclient.NewDeluge(serverUrl, password)
		if err := deluge.Connect(); err == nil {
			return true
		}
		return false
	})
	w.Bind("save", func(serverUrl, password string) {
		conf.Save(serverUrl, password)
	})
	w.Bind("exit", func() {
		w.Terminate()
	})
	type Config struct {
		ServerUrl, Password string
	}
	w.Bind("showConfig", func() Config {
		return Config{
			ServerUrl: configProperties.ServerUrl,
			Password:  configProperties.Password,
		}
	})
	w.Navigate(fmt.Sprintf(`data:text/html,%s`, assets.MustString("html/index.html")))
	w.Run()
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
