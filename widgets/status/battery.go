package status

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type BatteryStatus struct {
	BatteryID    string
	Status       string
	Capacity     float64
	CapacityFull float64
	Percent      float64
	Amps         float64
	Remaining    string
}

const batteryPath = "/sys/class/power_supply"

func ReadBatteries() ([]BatteryStatus, error) {
	dirs, err := ioutil.ReadDir(batteryPath)
	if err != nil {
		return nil, err
	}

	var batteries []BatteryStatus
	for _, dir := range dirs {
		if strings.ToLower(dir.Name()) == "ac" {
			continue
		}

		battery, err := ReadBattery(dir.Name())
		if err != nil {
			return nil, err
		}
		batteries = append(batteries, *battery)
	}

	return batteries, nil
}

func ReadBattery(name string) (*BatteryStatus, error) {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/uevent", batteryPath, name))
	if err != nil {
		return nil, err
	}

	vars := map[string]string{}

	for _, line := range strings.Split(string(file), "\n") {
		buf := strings.Split(line, "=")
		if len(buf) == 2 {
			vars[buf[0]] = buf[1]
		}
	}

	battery := &BatteryStatus{
		BatteryID: name,
		Status:    vars["POWER_SUPPLY_STATUS"],
	}

	energyFull, _ := strconv.ParseUint(vars["POWER_SUPPLY_ENERGY_FULL"], 10, 64)
	energyNow, _ := strconv.ParseUint(vars["POWER_SUPPLY_ENERGY_NOW"], 10, 64)
	powerNow, _ := strconv.ParseUint(vars["POWER_SUPPLY_POWER_NOW"], 10, 64)

	battery.Amps = float64(powerNow / 10000)
	battery.Capacity = float64(energyNow / 10000)
	battery.CapacityFull = float64(energyFull / 10000)
	battery.Percent = (battery.Capacity * 100.0) / battery.CapacityFull

	if battery.Amps > 0 {
		remaining := 0.0
		if battery.Status == "Charging" {
			remaining = (battery.CapacityFull - battery.Capacity) / battery.Amps
		} else {
			remaining = battery.Capacity / battery.Amps
		}

		seconds := int(remaining * 3600)
		hours := seconds / 3600
		minutes := (seconds - (hours * 3600)) / 60

		battery.Remaining = fmt.Sprintf("%.2d:%.2d", hours, minutes)
	} else {
		battery.Remaining = "00:00"
	}

	if battery.Status == "Unknown" {
		battery.Status = "Idle"
	}

	return battery, nil
}
