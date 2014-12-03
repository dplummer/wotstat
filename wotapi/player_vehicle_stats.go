package wotapi

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type PlayerVehicleStatsRoot struct {
	Count int
	Data  map[string][]PlayerVehicleStat
}

type PlayerVehicleStat struct {
	AccountId     int `json:"account_id"`
	All           VehicleBattleStat
	Clan          VehicleBattleStat
	Company       VehicleBattleStat
	Frags         map[string]int
	InGarage      bool `json:"in_garage"`
	MarkOfMastery int  `json:"mark_of_mastery"`
	MaxFrags      int  `json:"max_frags"`
	MaxXp         int  `json:"max_xp"`
	TankId        int  `json:"tank_id"`
	Team          VehicleBattleStat
}

type VehicleBattleStat struct {
	BattleAvgXp          int `json:"battle_avg_xp"`
	Battles              int
	CapturePoints        int `json:"capture_points"`
	DamageDealt          int `json:"damage_dealt"`
	DamageReceived       int `json:"damage_received"`
	Draws                int
	DroppedCapturePoints int `json:"dropped_capture_points"`
	Frags                int
	Hits                 int
	HitsPercents         int `json:"hits_percents"`
	Losses               int
	Shots                int
	Spotted              int
	SurvivedBattles      int `json:"survived_battles"`
	Wins                 int
	Xp                   int
}

func LoadPlayerVehicleStats() PlayerVehicleStatsRoot {
	playerVehicleStatsFile, err := ioutil.ReadFile("./data/playerstats.json")
	if err != nil {
		log.Fatalf("File error: %v\n", err)
	}

	var playerVehicleStats PlayerVehicleStatsRoot
	err = json.Unmarshal(playerVehicleStatsFile, &playerVehicleStats)
	if err != nil {
		log.Fatalf("Json unmarshal error: %v\n", err)
	}

	return playerVehicleStats
}

func (p *PlayerVehicleStatsRoot) Lookup(tankId int) (notFound VehicleBattleStat) {
	for _, vehicleStat := range p.Data["1004751607"] {
		if vehicleStat.TankId == tankId {
			return vehicleStat.All
		}
	}
	return notFound
}

func (stat *VehicleBattleStat) AvgDamage() float64 {
	return float64(stat.DamageDealt) / float64(stat.Battles)
}

func (stat *VehicleBattleStat) AvgSpot() float64 {
	return float64(stat.Spotted) / float64(stat.Battles)
}

func (stat *VehicleBattleStat) AvgFrag() float64 {
	return float64(stat.Frags) / float64(stat.Battles)
}

func (stat *VehicleBattleStat) AvgDef() float64 {
	return float64(stat.DroppedCapturePoints) / float64(stat.Battles)
}

func (stat *VehicleBattleStat) AvgWinRate() float64 {
	return float64(stat.Wins) / float64(stat.Battles) * 100.0
}
