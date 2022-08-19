package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/adelolmo/delauncher/crypt"
	"io/fs"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
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

/*var secretKey = []byte{11, 22, 33, 44, 55, 66, 77, 88, 99, 00, 11, 22, 33, 44, 55, 66}
var key = crypt.Key{
	Value: secretKey,
}*/

func GenRandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return blk, err
}

func NewConfig() Config {
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatal(err)
	}
	configDir := filepath.Join(userConfigDir, "delauncher")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		log.Fatalf("cannot create directory %s. Error: %s", configDir, err)
	}
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		log.Fatal(err)
	}

	key, err := readSecret(configDir)
	return Config{
		Filename: filepath.Join(configDir, "config.json"),
		Key:      key,
	}
}

func readSecret(configDir string) (crypt.Key, error) {
	keyPath := filepath.Join(configDir, "delauncher.key")
	content, err := ioutil.ReadFile(keyPath)
	if err != nil {
		blk, err := GenRandomBytes(16)
		if err != nil {
			return crypt.Key{Value: blk}, err
		}
		err = createSecretFile(keyPath, blk)
		if err != nil {
			return crypt.Key{}, err
		}
	}
	return crypt.Key{Value: content}, nil
}

func createSecretFile(path string, data []byte) error {
	err := ioutil.WriteFile(path, data, fs.FileMode(0600))
	if err != nil {
		return err
	}
	return nil
}

func (c Config) Get() (Properties, error) {
	configFile, err := os.OpenFile(c.Filename, os.O_RDONLY, 0700)
	if err != nil {
		return Properties{}, nil
	}

	r := bufio.NewReader(configFile)
	var config delugeConfig
	if err := json.NewDecoder(r).Decode(&config); err != nil {
		return Properties{}, errors.New("cannot deserialize configuration")
	}

	decryptedPassword, err := c.decrypt(config.Password)
	if err != nil {
		return Properties{}, err
	}

	err = configFile.Close()
	if err != nil {
		return Properties{}, fmt.Errorf("cannot close file %s. Error: %s", c.Filename, err)
	}

	return Properties{ServerUrl: config.ServerUrl, Password: decryptedPassword}, nil
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

func (c Config) decrypt(encryptedPassword []byte) (string, error) {
	if len(encryptedPassword) == 0 {
		return "", nil
	}
	result, err := c.Key.Decrypt(encryptedPassword)
	if err != nil {
		return "", errors.New("cannot decrypt password")
	}
	password := string(result[:])
	return password, nil
}
