package main

import (
	"fmt"
	"github.com/dplummer/wotstat/wn8"
	"github.com/dplummer/wotstat/wotapi"
	"log"
	"math"
)

type Wn8CombinedStat struct {
	TankId     int
	TankName   string
	ExpDamage  float64
	ExpDef     float64
	ExpFrag    float64
	ExpSpot    float64
	ExpWinRate float64
	Battles    int
	AvgDamage  float64
	AvgSpot    float64
	AvgFrag    float64
	AvgDef     float64
	AvgWinRate float64
}

type OverallWn8 struct {
	CombinedStats []Wn8CombinedStat
}

func (overall *OverallWn8) Add(tank wotapi.TankInfo, actual wotapi.VehicleBattleStat, expected wn8.ExpectedTank) {
	overall.CombinedStats = append(overall.CombinedStats, Wn8CombinedStat{
		tank.TankId,
		tank.NameI18n,
		expected.ExpDamage,
		expected.ExpDef,
		expected.ExpFrag,
		expected.ExpSpot,
		expected.ExpWinRate,
		actual.Battles,
		actual.AvgDamage(),
		actual.AvgSpot(),
		actual.AvgFrag(),
		actual.AvgDef(),
		actual.AvgWinRate(),
	})
}

func (overall *OverallWn8) Battles() int {
	var sum int

	for _, stat := range overall.CombinedStats {
		sum += stat.Battles
	}

	return sum
}

func (overall *OverallWn8) Wn8() float64 {
	var sumOfDamage, sumOfExpectedDamage,
		sumOfDef, sumOfExpectedDef,
		sumOfFrag, sumOfExpectedFrag,
		sumOfSpot, sumOfExpectedSpot,
		sumOfWinRate, sumOfExpectedWinRate,
		battles float64

	for _, stat := range overall.CombinedStats {
		battles = float64(stat.Battles)

		sumOfDamage += battles * stat.AvgDamage
		sumOfExpectedDamage += battles * stat.ExpDamage
		sumOfDef += battles * stat.AvgDef
		sumOfExpectedDef += battles * stat.ExpDef
		sumOfFrag += battles * stat.AvgFrag
		sumOfExpectedFrag += battles * stat.ExpFrag
		sumOfSpot += battles * stat.AvgSpot
		sumOfExpectedSpot += battles * stat.ExpSpot
		sumOfWinRate += battles * stat.AvgWinRate
		sumOfExpectedWinRate += battles * stat.ExpWinRate
	}

	totalBattles := float64(overall.Battles())

	rDamage := (sumOfDamage / totalBattles) / (sumOfExpectedDamage / totalBattles)
	rDef := (sumOfDef / totalBattles) / (sumOfExpectedDef / totalBattles)
	rFrag := (sumOfFrag / totalBattles) / (sumOfExpectedFrag / totalBattles)
	rSpot := (sumOfSpot / totalBattles) / (sumOfExpectedSpot / totalBattles)
	rWin := (sumOfWinRate / totalBattles) / (sumOfExpectedWinRate / totalBattles)

	log.Printf("rDamage: %f\n", rDamage)
	log.Printf("rDef: %f\n", rDef)
	log.Printf("rFrag: %f\n", rFrag)
	log.Printf("rSpot: %f\n", rSpot)
	log.Printf("rWin = %f / %f: %f\n", (sumOfWinRate / totalBattles), (sumOfExpectedWinRate / totalBattles), rWin)

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
		log.Fatalf("Infinite wn8, you win the game\n")
	}
	return wn8
}

func main() {
	tankInfo := wotapi.LoadTankInfo()
	playerVehicleStats := wotapi.LoadPlayerVehicleStats()
	expectedTankWn8s := wn8.LoadExpectedTankWn8()

	var overall OverallWn8

	for _, tank := range tankInfo.Data {
		vehicleStat := playerVehicleStats.Lookup(tank.TankId)

		if vehicleStat.Battles > 0 {
			expectedWn8 := expectedTankWn8s.Lookup(tank.TankId)

			if expectedWn8.TankId == 0 {
				log.Printf("Couldn't find expected WN8 for %s (%d)", tank.NameI18n, tank.TankId)
			} else {
				wn8 := expectedWn8.CalculateWn8(vehicleStat)
				fmt.Printf("%s (%d): Battles: %d WN8: %f Winrate: %f\n",
					tank.NameI18n,
					tank.TankId,
					vehicleStat.Battles,
					wn8,
					vehicleStat.AvgWinRate())

				overall.Add(tank, vehicleStat, expectedWn8)
			}
		}
	}

	fmt.Printf("battles: %d\nwn8: %v\n", overall.Battles(), overall.Wn8())
}
