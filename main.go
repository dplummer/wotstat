package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
)

type Wn8Root struct {
	Data []ExpectedTank
}

type ExpectedTank struct {
	TankId     int     `json:"IDNum,string"`
	ExpDamage  float64 `json:",string"`
	ExpDef     float64 `json:",string"`
	ExpFrag    float64 `json:",string"`
	ExpSpot    float64 `json:",string"`
	ExpWinRate float64 `json:",string"`
}

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

func loadExpectedTankWn8() Wn8Root {
	wn8file, err := ioutil.ReadFile("./wn8.json")
	if err != nil {
		log.Fatalf("File error: %v\n", err)
	}

	var wn8Root Wn8Root
	err = json.Unmarshal(wn8file, &wn8Root)
	if err != nil {
		log.Fatalf("Json unmarshal error: %v\n", err)
	}

	return wn8Root
}

func loadTankInfo() TankInfoRoot {
	tankInfoFile, err := ioutil.ReadFile("./tankinfo.json")
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

func loadPlayerVehicleStats() PlayerVehicleStatsRoot {
	playerVehicleStatsFile, err := ioutil.ReadFile("./playerstats.json")
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

func (p *PlayerVehicleStatsRoot) lookup(tankId int) (notFound VehicleBattleStat) {
	for _, vehicleStat := range p.Data["1004751607"] {
		if vehicleStat.TankId == tankId {
			return vehicleStat.All
		}
	}
	return notFound
}

func (w *Wn8Root) lookup(tankId int) (notFound ExpectedTank) {
	for _, tank := range w.Data {
		if tank.TankId == tankId {
			return tank
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
	return float64(stat.Wins) / float64(stat.Battles)
}

func (exp *ExpectedTank) calculateWn8(tank VehicleBattleStat) float64 {
	rDamage := tank.AvgDamage() / exp.ExpDamage
	rSpot := tank.AvgSpot() / exp.ExpSpot
	rFrag := tank.AvgFrag() / exp.ExpFrag
	rDef := tank.AvgDef() / exp.ExpDef
	rWin := tank.AvgWinRate() / exp.ExpWinRate

	rWinC := math.Max(0, (rWin-0.71)/(1-0.71))
	rDamageC := math.Max(0, (rDamage-0.22)/(1-0.22))
	rFragC := math.Max(0, math.Min(rDamageC+0.2, (rFrag-0.12)/(1-0.12)))
	rSpotC := math.Max(0, math.Min(rDamageC+0.1, (rSpot-0.38)/(1-0.38)))
	rDefC := math.Max(0, math.Min(rDamageC+0.1, (rDef-0.10)/(1-0.10)))

	wn8 := 980*rDamageC +
		210*rDamageC*rFragC +
		155*rFragC*rSpotC +
		75*rDefC*rFragC +
		145*math.Min(1.8, rWinC)

	if math.IsInf(wn8, 1) {
		log.Fatalf("Infinite wn8, you win the game\n%+v\n%+v\n", exp, tank)
	}
	return wn8
}

type VehicleWn8 struct {
	Battles int
	Wn8     float64
}

type Wn8s []VehicleWn8

func (wn8s Wn8s) Battles() int {
	var sumOfBattles int

	for _, vehicle := range wn8s {
		sumOfBattles += vehicle.Battles
	}

	return sumOfBattles
}

func (wn8s Wn8s) AverageWn8() float64 {
	var sumOfWn8s float64
	var sumOfBattles int

	for _, vehicle := range wn8s {
		sumOfWn8s += (vehicle.Wn8 * float64(vehicle.Battles))
		sumOfBattles += vehicle.Battles
	}

	return sumOfWn8s / float64(sumOfBattles)
}

func main() {
	tankInfo := loadTankInfo()
	playerVehicleStats := loadPlayerVehicleStats()
	expectedTankWn8s := loadExpectedTankWn8()

	var wn8s Wn8s

	for _, tank := range tankInfo.Data {
		vehicleStat := playerVehicleStats.lookup(tank.TankId)
		if vehicleStat.Battles > 0 {
			expectedWn8 := expectedTankWn8s.lookup(tank.TankId)

			if expectedWn8.TankId == 0 {
				log.Printf("Couldn't find expected WN8 for %s (%d)", tank.NameI18n, tank.TankId)
			} else {
				wn8 := expectedWn8.calculateWn8(vehicleStat)
				fmt.Printf("%s (%d): Battles: %d WN8: %f Winrate: %f\n",
					tank.NameI18n,
					tank.TankId,
					vehicleStat.Battles,
					wn8,
					vehicleStat.AvgWinRate())
				wn8s = append(wn8s, VehicleWn8{vehicleStat.Battles, wn8})
			}
		}
	}

	fmt.Printf("battles: %d\nwn8: %v\n", wn8s.Battles(), wn8s.AverageWn8())
}
