package wn8

import (
	"encoding/json"
	"github.com/dplummer/wotstat/wotapi"
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

func LoadExpectedTankWn8() Wn8Root {
	wn8file, err := ioutil.ReadFile("./data/wn8.json")
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

func (w *Wn8Root) Lookup(tankId int) (notFound ExpectedTank) {
	for _, tank := range w.Data {
		if tank.TankId == tankId {
			return tank
		}
	}
	return notFound
}

func (exp *ExpectedTank) CalculateWn8(tank wotapi.VehicleBattleStat) float64 {
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
