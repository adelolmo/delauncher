package deluge

import (
	"fmt"
	"github.com/adelolmo/delugeclient"
)

type Deluge struct {
	Url, Password string
}

func NewDeluge(url, password string) Deluge {
	return Deluge{
		Url:      url,
		Password: password,
	}
}

func (d Deluge) Add(link string) error {
	client := delugeclient.NewDeluge(d.Url, d.Password)
	if err := client.Connect(); err != nil {
		return fmt.Errorf("unable to stablish connection to server %s", d.Url)
	}
	if err := client.AddMagnet(link); err != nil {
		return fmt.Errorf("unable to add link %s", link)
	}
	return nil
}
