package thermal

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/lian/gonky/shader"
	"github.com/lian/gonky/texture"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	font "github.com/lian/gonky/font/terminus"
)

type ThermalGraph struct {
	Texture        *texture.Texture
	Redraw         chan bool
	Sensors        []string
	SensorValue    int
	SensorValueMax int
	SensorValueMin int
	SensorsGraph   []int
	GraphPadding   int

	FanLevel    int
	FanValue    int
	FanValueMax int
	FanValueMin int
	FanGraph    []int
}

func New(program *shader.Program) *ThermalGraph {
	s := &ThermalGraph{
		Texture:      &texture.Texture{X: 20, Y: 768 - (18 * 2), Width: 300, Height: 200},
		Redraw:       make(chan bool),
		GraphPadding: 8,

		SensorsGraph:   []int{},
		SensorValueMax: 0,
		SensorValueMin: 100,

		FanGraph:    []int{},
		FanValueMax: 0,
		FanValueMin: 10000,
	}
	s.Texture.Setup(program)
	return s
}

func (s *ThermalGraph) Render() {
	data := image.NewRGBA(image.Rect(0, 0, int(s.Texture.Width), int(s.Texture.Height)))
	gc := draw2dimg.NewGraphicContext(data)

	gc.SetFillColor(color.RGBA{0x33, 0x33, 0x33, 0xff})
	draw2dkit.Rectangle(gc, 0, 0, s.Texture.Width, s.Texture.Height)
	gc.Fill()

	padding := s.GraphPadding
	gc.SetFillColor(color.RGBA{0x66, 0x66, 0x66, 0xff})
	gc.SetStrokeColor(color.RGBA{0x66, 0x66, 0x66, 0xff})

	graphHeight := 40.0
	yOffset := 0.0

	gc.MoveTo(0, graphHeight+yOffset)
	var i, value int
	for i, value = range s.SensorsGraph {
		scaled := graphHeight - float64(int((float64(value-s.SensorValueMin)/float64(s.SensorValueMax-s.SensorValueMin))*graphHeight))
		height := scaled + float64(yOffset)
		gc.LineTo(float64(i*padding), height)
		gc.LineTo(float64(i*padding)+float64(padding), height)
	}
	gc.LineTo(float64(i*padding)+float64(padding), graphHeight+yOffset)
	gc.Close()
	gc.Fill()
	//gc.Stroke()

	x := (int(s.Texture.Width) - (font.Width * 4))
	y := int(yOffset + (graphHeight-font.Height)/2)
	font.DrawString(data, x, y, fmt.Sprintf("%dC", s.SensorValue), color.RGBA{0x66, 0x66, 0x66, 0xff})

	yOffset = 60.0

	gc.MoveTo(0, graphHeight+yOffset)
	for i, value = range s.FanGraph {
		scaled := graphHeight - float64(int((float64(value-s.FanValueMin)/float64(s.FanValueMax-s.FanValueMin))*graphHeight))
		height := scaled + float64(yOffset)
		gc.LineTo(float64(i*padding), height)
		gc.LineTo(float64(i*padding)+float64(padding), height)
	}
	gc.LineTo(float64(i*padding)+float64(padding), graphHeight+yOffset)
	gc.Close()
	gc.Fill()
	//gc.Stroke()

	x = (int(s.Texture.Width) - (font.Width * 12))
	y = int(yOffset + (graphHeight-font.Height)/2)
	font.DrawString(data, x, y, fmt.Sprintf("%d RPM L%d", s.FanValue, s.FanLevel), color.RGBA{0x66, 0x66, 0x66, 0xff})

	s.Texture.Write(&data.Pix)
}

func (s *ThermalGraph) Run() {
	s.UpdateThermal()
	s.UpdateFan()
	s.Redraw <- true

	ten := time.NewTicker(time.Second * 2)
	for {
		select {
		case <-ten.C:
			s.UpdateThermal()
			s.UpdateFan()
			break
		}
		s.Redraw <- true
	}
}

func (s *ThermalGraph) UpdateThermal() {
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

	s.SensorValue = int(max / 1000)

	if s.SensorValue > s.SensorValueMax {
		s.SensorValueMax = s.SensorValue + 6
	}

	if s.SensorValue < s.SensorValueMin {
		s.SensorValueMin = s.SensorValue - 6
	}

	maxItems := (int(s.Texture.Width) - (font.Width * 5)) / s.GraphPadding
	if len(s.SensorsGraph) >= maxItems {
		s.SensorsGraph = append(s.SensorsGraph[1:], s.SensorValue)
	} else {
		s.SensorsGraph = append(s.SensorsGraph, s.SensorValue)
	}
}

var fanRegexp *regexp.Regexp = regexp.MustCompile("speed:\t\t(\\d+)\nlevel:\t\t(.+)")

func (s *ThermalGraph) UpdateFan() {
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
		s.FanValueMax = s.FanValue + 100
	}

	if s.FanValue < s.FanValueMin {
		s.FanValueMin = s.FanValue - 100
	}

	maxItems := (int(s.Texture.Width) - (font.Width * 13)) / s.GraphPadding
	if len(s.FanGraph) >= maxItems {
		s.FanGraph = append(s.FanGraph[1:], s.FanValue)
	} else {
		s.FanGraph = append(s.FanGraph, s.FanValue)
	}
}
