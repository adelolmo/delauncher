package deluge

import (
	"fmt"
	"github.com/anacrolix/torrent/metainfo"
	"strings"
)

type MagnetLink struct {
	Address string
	Name    string
}

func NewLink(string string) (MagnetLink, error) {
	if strings.HasSuffix(string, ".torrent") {
		mi, err := metainfo.LoadFromFile(string)
		if err != nil {
			return MagnetLink{}, fmt.Errorf("error reading metainfo from stdin: %s", err)
		}
		info, err := mi.UnmarshalInfo()
		if err != nil {
			return MagnetLink{}, fmt.Errorf("error unmarshalling info: %s", err)
		}
		hash := mi.HashInfoBytes()
		linkAddress := fmt.Sprintf("%s", mi.Magnet(&hash, &info).String())
		return MagnetLink{
			Address: linkAddress,
			Name:    info.Name,
		}, nil
	}

	params := string[61:]
	p := strings.Split(params, "&")
	name := p[0][3:]

	return MagnetLink{
		Address: string,
		Name:    name,
	}, nil
}
