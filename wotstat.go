package main

import (
	"fmt"
	"github.com/dplummer/wotstat/wn8"
	"github.com/dplummer/wotstat/wotapi"
	"log"
)

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
	tankInfo := wotapi.LoadTankInfo()
	playerVehicleStats := wotapi.LoadPlayerVehicleStats()
	expectedTankWn8s := wn8.LoadExpectedTankWn8()

	var wn8s Wn8s

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
				wn8s = append(wn8s, VehicleWn8{vehicleStat.Battles, wn8})
			}
		}
	}

	fmt.Printf("battles: %d\nwn8: %v\n", wn8s.Battles(), wn8s.AverageWn8())
}
