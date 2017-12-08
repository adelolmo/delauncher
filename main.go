package main

import (
	"os/exec"
	"bytes"
	"fmt"
	"os"
	"encoding/json"
	"bufio"
	"os/user"
	"path/filepath"
	"log"
	"strings"
	"github.com/ProtonMail/ui"
	"github.com/adelolmo/delugeclient"
)

const (
	CONFIG_DIR string = ".config/delauncher"
	CONFIG_FILE string = "config.json"
)

type DelugeConfig struct {
	ServerUrl string
	Password  string
}

func main() {
	switch len(os.Args) {
	case 1:
		config()
	case 2:
		addMagnet(os.Args[1])
	default:
		fmt.Fprint(os.Stderr, "Usage: delauncher (MAGNET_LINK)")
		os.Exit(1)
	}
}

func addMagnet(magnet string) {
	serverUrl, password := getConfig(filepath.Join(getHome(), CONFIG_DIR, CONFIG_FILE))
	client := delugeclient.NewDeluge(serverUrl, password)
	if err := client.Connect(); err != nil {
		fmt.Errorf("unable to stablish connection to server %s", serverUrl)
		notify(fmt.Sprintf("Unable to stablish connection to server %s", serverUrl))
		os.Exit(2)
	}
	if err := client.AddMagnet(magnet); err != nil {
		fmt.Errorf("unable to add magnet link %s", magnet)
		notify("Error! Can't add magnet link")
		os.Exit(2)
	}

	notify(fmt.Sprintf("Magnet added:\n%s", getLinkName(magnet)))
}

func getHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func getConfig(configFilename string) (string, string) {

	os.MkdirAll(filepath.Join(getHome(), CONFIG_DIR), 0755)
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		return "", ""
	}

	configFile, err := os.OpenFile(configFilename, os.O_RDONLY, 0700)
	if err != nil {
		log.Fatalf("Cannot open cache file %s. Error: %s", configFilename, err)
	}
	defer configFile.Close()

	r := bufio.NewReader(configFile)
	var config DelugeConfig
	if err := json.NewDecoder(r).Decode(&config); err != nil {
		panic(err)
	}
	return config.ServerUrl, config.Password
}

func config() {
	err := ui.Main(func() {
		configFilename := filepath.Join(getHome(), CONFIG_DIR, CONFIG_FILE)
		serverUrl, password := getConfig(configFilename)

		serverUrlField := ui.NewEntry()
		passwordField := ui.NewPasswordEntry()

		serverUrlField.SetText(serverUrl)
		passwordField.SetText(password)

		button := ui.NewButton("Save & Quit")
		box := ui.NewVerticalBox()
		box.Append(ui.NewLabel("Deluge server URL:"), false)
		box.Append(serverUrlField, false)
		box.Append(ui.NewLabel("Password:"), false)
		box.Append(passwordField, false)
		box.Append(button, false)
		window := ui.NewWindow("Delauncher", 400, 150, false)
		window.SetChild(box)
		button.OnClicked(func(*ui.Button) {
			f, err := os.Create(configFilename)
			if err != nil {
				panic(err)
			}
			defer f.Close()
			w := bufio.NewWriter(f)

			config := &DelugeConfig{ServerUrl:serverUrlField.Text(), Password:passwordField.Text()}
			json.NewEncoder(w).Encode(&config)
			w.Flush()

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

func notify(message string) {
	cmd := exec.Command("notify-send", "Deluge", message, "--icon=delauncher")
	cmdOutput := &bytes.Buffer{}
	cmd.Stdout = cmdOutput
	if err := cmd.Run(); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("==> Error: %s\n", err.Error()))
		os.Exit(3)
	}
}

func getLinkName(magnet string) string {
	params := magnet[61:]
	p := strings.Split(params, "&")
	return p[0][3:]
}
