package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adelolmo/delauncher/crypt"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

const (
	configDir  string = ".config/delauncher"
	configFile string = "config.json"
)

type Config struct {
	Filename string
	Key      crypt.Key
}

type delugeConfig struct {
	ServerUrl string `json:"serverUrl"`
	Password  []byte `json:"password"`
}

type Properties struct {
	ServerUrl, Password string
}

func NewConfig(key crypt.Key) Config {
	return Config{
		Filename: filepath.Join(getHome(), configDir, configFile),
		Key:      key,
	}
}

func (c Config) Get() (Properties, error) {
	if err := os.MkdirAll(filepath.Join(getHome(), configDir), 0755); err != nil {
		fmt.Printf("Cannot create directory %s in home %s. Error: %s", c.Filename, getHome(), err)
	}
	if _, err := os.Stat(c.Filename); os.IsNotExist(err) {
		return Properties{}, err
	}

	configFile, err := os.OpenFile(c.Filename, os.O_RDONLY, 0700)
	if err != nil {
		return Properties{}, fmt.Errorf("cannot open file %s. Error: %s", c.Filename, err)
	}

	r := bufio.NewReader(configFile)
	var config delugeConfig
	if err := json.NewDecoder(r).Decode(&config); err != nil {
		return Properties{}, errors.New("cannot deserialize configuration")
	}

	result, err := c.Key.Decrypt(config.Password)
	if err != nil {
		return Properties{}, errors.New("cannot decrypt password")
	}
	password := string(result[:])

	err = configFile.Close()
	if err != nil {
		return Properties{}, fmt.Errorf("cannot close file %s. Error: %s", c.Filename, err)
	}

	return Properties{ServerUrl: config.ServerUrl, Password: password}, nil
}

func (c Config) Save(serverUrl, password string) {
	f, err := os.Create(c.Filename)
	if err != nil {
		panic(err)
	}
	w := bufio.NewWriter(f)

	encryptedPassword, err := c.Key.Encrypt([]byte(password))
	if err != nil {
		log.Fatal(err)
	}
	config := &delugeConfig{ServerUrl: serverUrl, Password: encryptedPassword}
	if err = json.NewEncoder(w).Encode(&config); err != nil {
		log.Fatalf("Cannot encode json configuration. Error: %s", err)
	}
	if err = w.Flush(); err != nil {
		log.Fatalf("Cannot flush into file %s. Error: %s", f.Name(), err)
	}

	if err = f.Close(); err != nil {
		log.Fatalf("Cannot close file %s. Error: %s", c.Filename, err)
	}
}

func getHome() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	return usr.HomeDir
}
