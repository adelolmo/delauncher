package magnet

import (
	"fmt"
	"github.com/adelolmo/delugeclient"
	"github.com/anacrolix/torrent/metainfo"
	"strings"
)

type Link struct {
	Address string
	Name    string
}

type Server struct {
	Url, Password string
}

func NewLink(string string) (Link, error) {
	if strings.HasSuffix(string, ".torrent") {
		mi, err := metainfo.LoadFromFile(string)
		if err != nil {
			return Link{}, fmt.Errorf("error reading metainfo from stdin: %s", err)
		}
		info, err := mi.UnmarshalInfo()
		if err != nil {
			return Link{}, fmt.Errorf("error unmarshalling info: %s", err)
		}
		linkAddress := fmt.Sprintf("%s", mi.Magnet(info.Name, mi.HashInfoBytes()).String())
		return Link{
			Address: linkAddress,
			Name:    info.Name,
		}, nil
	}

	params := string[61:]
	p := strings.Split(params, "&")
	name := p[0][3:]

	return Link{
		Address: string,
		Name:    name,
	}, nil
}

func (server Server) Add(link Link) error {
	client := delugeclient.NewDeluge(server.Url, server.Password)
	if err := client.Connect(); err != nil {
		return fmt.Errorf("unable to stablish connection to server %s", server.Url)
	}
	if err := client.AddMagnet(link.Address); err != nil {
		return fmt.Errorf("unable to add link %s", link.Address)
	}
	return nil
}
