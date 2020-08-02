package services

import (
	"github.com/stac47/myroomies/pkg/models"
)

const (
	name    = "MyRoomies"
	licence = "GPL Version 3"
	creator = "stac47 - https://github.com/stac47"
)

func GetGlobalInfo() models.GlobalInfo {
	return models.GlobalInfo{
		Name:    name,
		Version: version,
		Licence: licence,
		Creator: creator,
	}
}
