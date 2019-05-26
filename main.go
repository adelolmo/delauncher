package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/adelolmo/delauncher/crypt"
	"github.com/adelolmo/delauncher/magnet"
	"github.com/andlabs/ui"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

const (
	ConfigDir  string = ".config/delauncher"
	ConfigFile string = "config.json"
)

var SecretKey = []byte{11, 22, 33, 44, 55, 66, 77, 88, 99, 00, 11, 22, 33, 44, 55, 66}

type DelugeConfig struct {
	ServerUrl string
	Password  []byte
}

func main() {
	switch len(os.Args) {
	case 1:
		config()
	case 2:
		addMagnet(magnet.Link{
			Address: os.Args[1],
		})
	default:
		fmt.Print("usage: delauncher (MAGNET_LINK)")
		os.Exit(1)
	}
}

func addMagnet(magnetLink magnet.Link) {
	serverUrl, password := getConfig(filepath.Join(getHome(), ConfigDir, ConfigFile))
	if err := magnetLink.Add(serverUrl, password); err != nil {
		fmt.Printf(err.Error())
		notify(err.Error())
		os.Exit(2)
	}
	notify(fmt.Sprintf("Magnet added:\n%s", magnetLink.Name()))
}

func getHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}

func getConfig(configFilename string) (string, string) {

	if err := os.MkdirAll(filepath.Join(getHome(), ConfigDir), 0755); err != nil {
		fmt.Printf("Cannot create directory %s in home %s. Error: %s", configFilename, getHome(), err)
	}
	if _, err := os.Stat(configFilename); os.IsNotExist(err) {
		return "", ""
	}

	configFile, err := os.OpenFile(configFilename, os.O_RDONLY, 0700)
	if err != nil {
		fmt.Printf("Cannot open cache file %s. Error: %s", configFilename, err)
	}

	r := bufio.NewReader(configFile)
	var config DelugeConfig
	if err := json.NewDecoder(r).Decode(&config); err != nil {
		panic(err)
	}

	result, err := crypt.Decrypt(SecretKey, config.Password)
	if err != nil {
		log.Fatal("Cannot decrypt secret")
	}
	password := string(result[:])

	err = configFile.Close()
	if err != nil {
		log.Fatalf("Cannot close file %s. Error: %s", configFilename, err)
	}

	return config.ServerUrl, password
}

func config() {
	err := ui.Main(func() {
		configFilename := filepath.Join(getHome(), ConfigDir, ConfigFile)
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
			w := bufio.NewWriter(f)

			encryptedPassword, err := crypt.Encrypt(SecretKey, []byte(passwordField.Text()))
			if err != nil {
				log.Fatal(err)
			}
			config := &DelugeConfig{ServerUrl: serverUrlField.Text(), Password: encryptedPassword}
			if err = json.NewEncoder(w).Encode(&config); err != nil {
				log.Fatalf("Cannot encode json configuration. Error: %s", err)
			}
			if err = w.Flush(); err != nil {
				log.Fatalf("Cannot flush into file %s. Error: %s", f.Name(), err)
			}

			if err = f.Close(); err != nil {
				log.Fatalf("Cannot close file %s. Error: %s", configFilename, err)
			}

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
		fmt.Printf("==> Error: %s\n", err.Error())
		os.Exit(3)
	}
}
