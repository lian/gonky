package status

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lian/gonky/font/terminus"
	"github.com/lian/gonky/shader"
	"github.com/lian/gonky/texture"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	psutil_cpu "github.com/shirou/gopsutil/cpu"
	psutil_mem "github.com/shirou/gopsutil/mem"
	psutil_net "github.com/shirou/gopsutil/net"
)

type Status struct {
	Texture *texture.Texture
	Redraw  chan bool
	Time    string
	Memory  string
	CPU     string
	Network string
	Battery string
	Thermal string
	Fan     string
}

func New(x, y, width, height float64, program *shader.Program) *Status {
	status := &Status{
		Texture: &texture.Texture{X: x, Y: y, Width: width, Height: height},
		Redraw:  make(chan bool),
	}
	status.Texture.Setup(program)
	return status
}

func (s *Status) Render() {
	data := image.NewRGBA(image.Rect(0, 0, int(s.Texture.Width), int(s.Texture.Height)))
	gc := draw2dimg.NewGraphicContext(data)

	gc.SetFillColor(color.RGBA{0xcc, 0xcc, 0xcc, 0xff})
	draw2dkit.Rectangle(gc, 0, 0, s.Texture.Width, s.Texture.Height)
	gc.Fill()

	text_height := 3
	terminus.DrawString(data, terminus.Width, text_height, s.Time, color.Black)

	buf := strings.Join([]string{s.Memory, s.Fan, s.Thermal, s.CPU, s.Network, s.Battery}, "  |  ")
	right := int(s.Texture.Width) - ((len(buf) * terminus.Width) + terminus.Width)
	terminus.DrawString(data, right, text_height, buf, color.Black)

	s.Texture.Write(&data.Pix)
}

func (s *Status) Run() {
	s.UpdateTime()
	s.UpdateMemory()
	s.UpdateCPU()
	s.UpdateNetwork()
	s.UpdateBattery()
	s.UpdateThermal()
	s.UpdateFan()
	s.Redraw <- true

	one := time.NewTicker(time.Second * 1)
	five := time.NewTicker(time.Second * 5)
	ten := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-one.C:
			s.UpdateTime()
			break
		case <-five.C:
			s.UpdateCPU()
			s.UpdateNetwork()
			break
		case <-ten.C:
			s.UpdateMemory()
			s.UpdateBattery()
			s.UpdateThermal()
			s.UpdateFan()
			break
		}
		s.Redraw <- true
	}
}

func (s *Status) UpdateTime() {
	s.Time = time.Now().Format("15:04:05 02.01.2006")
}

func (s *Status) UpdateMemory() {
	v, _ := psutil_mem.VirtualMemory()
	s.Memory = fmt.Sprintf("%.2f%% RAM", v.UsedPercent)
}

func (s *Status) UpdateCPU() {
	//info, _ := psutil_cpu.Info(); spew.Dump(info)
	percent, err := psutil_cpu.Percent(0, false)
	if err == nil && len(percent) == 1 {
		s.CPU = fmt.Sprintf("%.2f%% CPU", percent[0])
	}
}

func (s *Status) UpdateNetwork() {
	stats, _ := psutil_net.IOCounters(true)
	networks := []string{}

	for _, v := range stats {
		switch v.Name {
		case "lo":
			continue
		case "enp0s25":
			v.Name = "lan"
		case "wlp3s0":
			v.Name = "wifi"
		}

		if v.BytesRecv == 0 {
			continue
		}

		buf := fmt.Sprintf("%.1f-%s-%.1f", float64(v.BytesRecv), v.Name, float64(v.BytesSent))
		networks = append(networks, buf)
	}

	s.Network = strings.Join(networks, " | ")
}

func (s *Status) UpdateBattery() {
	b, err := ReadBattery("BAT0")
	if err == nil {
		if b.Status == "Idle" {
			s.Battery = fmt.Sprintf("idle %.0f%%", b.Percent)
		} else {
			s.Battery = fmt.Sprintf("%s %sh %.0fmA %.0f%%", strings.ToLower(b.Status), b.Remaining, b.Amps, b.Percent)
		}
	}
}

func (s *Status) UpdateThermal() {
	sensors := []string{
		"/sys/devices/platform/coretemp.0/hwmon/hwmon0/temp1_input",
		"/sys/devices/platform/coretemp.0/hwmon/hwmon0/temp2_input",
		"/sys/devices/platform/coretemp.0/hwmon/hwmon0/temp3_input",
	}
	var result uint64 = 0

	for _, path := range sensors {
		if buf, err := ioutil.ReadFile(path); err == nil {
			str := strings.Replace(string(buf), "\n", "", -1)
			value, err := strconv.ParseUint(str, 10, 64)
			if err == nil && value > result {
				result = value
			}
		}
	}

	result = result / 1000
	s.Thermal = fmt.Sprintf("%dC", result)
}

var fanRegexp *regexp.Regexp = regexp.MustCompile("speed:\t\t(\\d+)\nlevel:\t\t(.+)")

func (s *Status) UpdateFan() {
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
			s.Fan = fmt.Sprintf("%d RPM L%d", rpm, level)
		}
	}
}
