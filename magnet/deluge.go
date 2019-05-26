package magnet

import (
	"fmt"
	"github.com/adelolmo/delugeclient"
	"strings"
)

type Link struct {
	Address string
}

func (link Link) Add(serverUrl, password string) error {
	client := delugeclient.NewDeluge(serverUrl, password)
	if err := client.Connect(); err != nil {
		return fmt.Errorf("unable to stablish connection to server %s", serverUrl)
	}
	if err := client.AddMagnet(link.Address); err != nil {
		return fmt.Errorf("unable to add magnet link %s", link.Address)
	}
	return nil
}

func (link Link) Name() string {
	params := link.Address[61:]
	p := strings.Split(params, "&")
	return p[0][3:]
}
