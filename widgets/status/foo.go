package status

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"strings"
	"time"

	"github.com/lian/gonky/shader"
	"github.com/lian/gonky/texture"
	"github.com/lian/gonky/widgets"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	font "github.com/lian/gonky/font/terminus"

	psutil_net "github.com/shirou/gopsutil/net"
)

type Net struct {
	Name          string
	Time          time.Time
	LastBytesRecv uint64
	LastBytesSent uint64
	RateRecv      float64
	RateSent      float64
}

type Status struct {
	Texture    *texture.Texture
	Redraw     chan bool
	Time       string
	Network    string
	NetworkMap map[string]*Net
	Battery    string
	Stats      *widgets.Stats
}

var FontPadding int = 3

func New(windowWidth, windowHeight int, program *shader.Program, stats *widgets.Stats) *Status {
	height := float64(font.Height + (2 * FontPadding))
	status := &Status{
		Texture:    &texture.Texture{X: 0, Y: float64(windowHeight), Width: float64(windowWidth), Height: height},
		Redraw:     make(chan bool),
		NetworkMap: map[string]*Net{},
		Stats:      stats,
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

	text_height := FontPadding
	font.DrawString(data, font.Width, text_height, s.Time, color.Black)

	thermalText := fmt.Sprintf("%dC", s.Stats.ThermalValue)
	fanText := fmt.Sprintf("%d RPM L%d", s.Stats.FanValue, s.Stats.FanLevel)
	memoryText := fmt.Sprintf("%.2f%% RAM", s.Stats.MemoryValue)
	cpuText := fmt.Sprintf("%.2f%% CPU", s.Stats.CpuValue)

	buf := strings.Join([]string{memoryText, fanText, thermalText, cpuText, s.Network, s.Battery}, "  |  ")
	right := int(s.Texture.Width) - ((len(buf) * font.Width) + font.Width)
	font.DrawString(data, right, text_height, buf, color.Black)

	s.Texture.Write(&data.Pix)
}

func (s *Status) Run() {
	s.UpdateTime()
	s.UpdateNetwork()
	s.UpdateBattery()
	s.Redraw <- true

	five := time.NewTicker(time.Second * 5)
	ten := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-five.C:
			s.UpdateTime()
			s.UpdateNetwork()
			break
		case <-ten.C:
			s.UpdateBattery()
			break
		}
		s.Redraw <- true
	}
}

func (s *Status) UpdateTime() {
	//s.Time = time.Now().Format("15:04:05 02.01.2006")
	s.Time = time.Now().Format("15:04 02.01.2006")
}

var NetworkNamesMap map[string]string = map[string]string{
	"enp0s25": "lan",
	"wlp3s0":  "wifi",
}

func (s *Status) UpdateNetwork() {
	stats, _ := psutil_net.IOCounters(true)
	networks := []string{}

	isAvailable := map[string]bool{}
	for id, _ := range s.NetworkMap {
		isAvailable[id] = false
	}

	for _, v := range stats {
		if v.Name == "lo" || v.BytesRecv == 0 {
			continue
		}
		if alias, ok := NetworkNamesMap[v.Name]; ok {
			v.Name = alias
		}

		isAvailable[v.Name] = true

		var net *Net
		var ok bool

		if net, ok = s.NetworkMap[v.Name]; ok {
			now := time.Now()
			timeDiff := now.Sub(net.Time)
			net.Time = now

			recvDiff := math.Abs(float64(net.LastBytesRecv) - float64(v.BytesRecv))
			sentDiff := math.Abs(float64(net.LastBytesSent) - float64(v.BytesSent))
			net.LastBytesRecv = v.BytesRecv
			net.LastBytesSent = v.BytesSent

			net.RateRecv = recvDiff / timeDiff.Seconds()
			net.RateSent = sentDiff / timeDiff.Seconds()
		} else {
			net = &Net{
				Time:          time.Now(),
				Name:          v.Name,
				LastBytesRecv: v.BytesRecv,
				LastBytesSent: v.BytesSent,
				RateRecv:      0,
				RateSent:      0,
			}
			s.NetworkMap[v.Name] = net
		}

		buf := fmt.Sprintf("%.1f-%s-%.1f", net.RateRecv/1024, v.Name, net.RateSent/1024)
		networks = append(networks, buf)
	}

	for id, state := range isAvailable {
		if state == false {
			delete(s.NetworkMap, id)
		}
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
