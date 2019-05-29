package magnet

import (
	"fmt"
	"github.com/adelolmo/delugeclient"
	"strings"
)

type Link struct {
	Address string
}

type Server struct {
	Url, Password string
}

func (server Server) Add(link Link) error {
	client := delugeclient.NewDeluge(server.Url, server.Password)
	if err := client.Connect(); err != nil {
		return fmt.Errorf("unable to stablish connection to server %s", server.Url)
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
