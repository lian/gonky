package status

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

type BatteryStatus struct {
	BatteryID    string
	EnergyFull   uint64
	EnergyNow    uint64
	PowerNow     uint64
	Status       string
	Capacity     float64
	CapacityFull float64
	Percent      float64
	Remaining    string
	Amps         float64
	//EnergyFullDesign uint64
	//Timestamp time.Time
	//CycleCount       uint64
	//Manufacturer     string
	//ModelName        string
	//SerialNumber     string
	//Technology       string
}

const batteryPath = "/sys/class/power_supply"

func Batteries() ([]BatteryStatus, error) {
	dirs, err := ioutil.ReadDir(batteryPath)
	if err != nil {
		return nil, err
	}

	var batteries []BatteryStatus
	for _, dir := range dirs {
		if strings.ToLower(dir.Name()) == "ac" {
			continue
		}

		battery, err := readBattery(dir.Name())
		if err != nil {
			return nil, err
		}
		batteries = append(batteries, *battery)
	}

	return batteries, nil
}

func readBattery(name string) (*BatteryStatus, error) {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/%s/uevent", batteryPath, name))
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(file), "\n")
	vars := make(map[string]string)

	for _, line := range lines {
		data := strings.Split(line, "=")
		if len(data) == 2 {
			vars[data[0]] = data[1]
		}
	}

	battery := &BatteryStatus{}
	battery.BatteryID = name
	battery.Status = vars["POWER_SUPPLY_STATUS"]
	battery.EnergyFull, _ = strconv.ParseUint(vars["POWER_SUPPLY_ENERGY_FULL"], 10, 64)
	//battery.EnergyFullDesign, _ = strconv.ParseUint(vars["POWER_SUPPLY_ENERGY_FULL_DESIGN"], 10, 64)
	battery.EnergyNow, _ = strconv.ParseUint(vars["POWER_SUPPLY_ENERGY_NOW"], 10, 64)
	battery.PowerNow, _ = strconv.ParseUint(vars["POWER_SUPPLY_POWER_NOW"], 10, 64)
	//battery.Timestamp = time.Now()
	//battery.CycleCount, _ = strconv.ParseUint(vars["POWER_SUPPLY_CYCLE_COUNT"], 10, 64)
	//battery.Manufacturer = vars["POWER_SUPPLY_MANUFACTURER"]
	//battery.ModelName = vars["POWER_SUPPLY_MODEL_NAME"]
	//battery.SerialNumber = vars["POWER_SUPPLY_SERIAL_NUMBER"]
	//battery.Technology = vars["POWER_SUPPLY_TECHNOLOGY"]

	battery.Amps = float64(battery.PowerNow / 10000.0)
	battery.Capacity = float64(battery.EnergyNow) / 10000.0
	battery.CapacityFull = float64(battery.EnergyFull) / 10000.0
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
