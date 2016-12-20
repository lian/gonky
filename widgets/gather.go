package widgets

import (
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	psutil_cpu "github.com/shirou/gopsutil/cpu"
	psutil_mem "github.com/shirou/gopsutil/mem"
)

func NewStats() *Stats {
	s := &Stats{
		Updated:              make(chan bool),
		FanGraphMaxCount:     60,
		ThermalGraphMaxCount: 60,
		MemoryGraphMaxCount:  60,
		CpuGraphMaxCount:     60,

		ThermalValueMin: 0xffff,
		//FanValueMin:     0xffff,
		MemoryValueMin: 0xffff,
		CpuValueMin:    0xffff,

		FanValueMin: 0,
		FanValueMax: 10000,
	}
	return s
}

type Stats struct {
	Updated chan bool

	ThermalValue         int
	ThermalValueMax      int
	ThermalValueMin      int
	ThermalGraph         []int
	ThermalGraphMaxCount int

	FanLevel         int
	FanValue         int
	FanValueMax      int
	FanValueMin      int
	FanGraph         []int
	FanGraphMaxCount int

	MemoryValue         float64
	MemoryValueMax      float64
	MemoryValueMin      float64
	MemoryGraph         []float64
	MemoryGraphMaxCount int

	CpuValue         float64
	CpuValueMax      float64
	CpuValueMin      float64
	CpuGraph         []float64
	CpuGraphMaxCount int
}

func (s *Stats) Run() {
	s.UpdateMemory()
	s.UpdateCPU()
	s.UpdateThermal()
	s.UpdateFan()
	s.Updated <- true

	five := time.NewTicker(time.Second * 5)
	ten := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-five.C:
			s.UpdateCPU()
			s.UpdateThermal()
			s.UpdateFan()
			break
		case <-ten.C:
			s.UpdateMemory()
			break
		}
		s.Updated <- true
	}
}

var fanRegexp *regexp.Regexp = regexp.MustCompile("speed:\t\t(\\d+)\nlevel:\t\t(.+)")

func (s *Stats) UpdateFan() {
	var rpm int
	var level int
	var file string = "/proc/acpi/ibm/fan"

	if buf, err := ioutil.ReadFile(file); err == nil {
		m := fanRegexp.FindStringSubmatch(string(buf))
		if len(m) == 3 {
			rpm, _ = strconv.Atoi(m[1])
			if m[2] == "disengaged" {
				level = 8
			} else {
				level, _ = strconv.Atoi(m[2])
			}
		}
	}

	s.FanLevel = level
	s.FanValue = rpm

	if s.FanValue > s.FanValueMax {
		s.FanValue = s.FanValueMax
	}

	if s.FanValue < s.FanValueMin {
		s.FanValue = s.FanValueMin
	}

	if len(s.FanGraph) >= s.FanGraphMaxCount {
		s.FanGraph = append(s.FanGraph[1:], s.FanValue)
	} else {
		s.FanGraph = append(s.FanGraph, s.FanValue)
	}
}

func (s *Stats) UpdateThermal() {
	var max uint64 = 0

	sensors := []string{
		"/sys/devices/platform/coretemp.0/hwmon/hwmon0/temp1_input",
		"/sys/devices/platform/coretemp.0/hwmon/hwmon0/temp2_input",
		"/sys/devices/platform/coretemp.0/hwmon/hwmon0/temp3_input",
	}

	for _, path := range sensors {
		if buf, err := ioutil.ReadFile(path); err == nil {
			str := strings.Replace(string(buf), "\n", "", -1)
			value, err := strconv.ParseUint(str, 10, 64)
			if err == nil && value > max {
				max = value
			}
		}
	}

	s.ThermalValue = int(max / 1000)

	if s.ThermalValue > s.ThermalValueMax {
		s.ThermalValueMax = s.ThermalValue
	}

	if s.ThermalValue < s.ThermalValueMin {
		s.ThermalValueMin = s.ThermalValue
	}

	if len(s.ThermalGraph) >= s.ThermalGraphMaxCount {
		s.ThermalGraph = append(s.ThermalGraph[1:], s.ThermalValue)
	} else {
		s.ThermalGraph = append(s.ThermalGraph, s.ThermalValue)
	}
}

func (s *Stats) UpdateMemory() {
	v, _ := psutil_mem.VirtualMemory()
	s.MemoryValue = v.UsedPercent

	if s.MemoryValue > s.MemoryValueMax {
		s.MemoryValueMax = s.MemoryValue
	}

	if s.MemoryValue < s.MemoryValueMin {
		s.MemoryValueMin = s.MemoryValue
	}

	if len(s.MemoryGraph) >= s.MemoryGraphMaxCount {
		s.MemoryGraph = append(s.MemoryGraph[1:], s.MemoryValue)
	} else {
		s.MemoryGraph = append(s.MemoryGraph, s.MemoryValue)
	}
}

func (s *Stats) UpdateCPU() {
	//info, _ := psutil_cpu.Info(); spew.Dump(info)
	percent, err := psutil_cpu.Percent(0, false)
	if err != nil || len(percent) != 1 {
		return
	}

	s.CpuValue = percent[0]

	if s.CpuValue > s.CpuValueMax {
		s.CpuValueMax = s.CpuValue
	}

	if s.CpuValue < s.CpuValueMin {
		s.CpuValueMin = s.CpuValue
	}

	if len(s.CpuGraph) >= s.CpuGraphMaxCount {
		s.CpuGraph = append(s.CpuGraph[1:], s.CpuValue)
	} else {
		s.CpuGraph = append(s.CpuGraph, s.CpuValue)
	}
}
