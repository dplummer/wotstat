package wotapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type TankInfoRoot struct {
	Count int
	Data  map[string]TankInfo
}

type TankInfo struct {
	ContourImage  string `json:"contour_image"`
	Image         string
	ImageSmall    string `json:"image_small"`
	IsPremium     bool   `json:"is_premium"`
	Level         int
	Name          string
	NameI18n      string `json:"name_i18n"`
	Nation        string
	NationI18n    string `json:"nation_i18n"`
	ShortNameI18n string `json:"short_name_i18n"`
	TankId        int    `json:"tank_id"`
	Type          string
	TypeI18n      string `json:"type_i18n"`
}

func LoadTankInfo() TankInfoRoot {
	tankInfoFile, err := ioutil.ReadFile("./data/tankinfo.json")
	if err != nil {
		log.Fatalf("File error: %v\n", err)
	}

	var tankInfo TankInfoRoot
	err = json.Unmarshal(tankInfoFile, &tankInfo)
	if err != nil {
		log.Fatalf("Json unmarshal error: %v\n", err)
	}

	return tankInfo
}
