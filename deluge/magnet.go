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
		return torrentFileToMagnetLink(string)
	}

	magnetUri, err := metainfo.ParseMagnetUri(string)
	if err != nil {
		return MagnetLink{}, fmt.Errorf("invalid magnet link: %w", err)
	}

	return MagnetLink{
		Address: magnetUri.String(),
		Name:    magnetUri.DisplayName,
	}, nil
}

func torrentFileToMagnetLink(string string) (MagnetLink, error) {
	mi, err := metainfo.LoadFromFile(string)
	if err != nil {
		return MagnetLink{}, fmt.Errorf("error reading metainfo from stdin: %w", err)
	}

	info, err := mi.UnmarshalInfo()
	if err != nil {
		return MagnetLink{}, fmt.Errorf("error unmarshalling info: %w", err)
	}

	hash := mi.HashInfoBytes()
	linkAddress := fmt.Sprintf("%s", mi.Magnet(&hash, &info).String())
	return MagnetLink{
		Address: linkAddress,
		Name:    info.Name,
	}, nil
}
