package thermal

import (
	"fmt"
	"image"
	"image/color"
	"io/ioutil"
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
}

func New(program *shader.Program) *ThermalGraph {
	s := &ThermalGraph{
		Texture: &texture.Texture{X: 20, Y: 768 - (18 * 2), Width: 300, Height: 40},
		Redraw:  make(chan bool),
		Sensors: []string{
			"/sys/devices/platform/coretemp.0/hwmon/hwmon0/temp1_input",
			"/sys/devices/platform/coretemp.0/hwmon/hwmon0/temp2_input",
			"/sys/devices/platform/coretemp.0/hwmon/hwmon0/temp3_input",
		},
		SensorsGraph:   []int{},
		GraphPadding:   8,
		SensorValueMax: 0,
		SensorValueMin: 100,
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

	gc.MoveTo(0, 40)
	var i, value int
	for i, value = range s.SensorsGraph {
		height := 40 - float64(int((float64(value-s.SensorValueMin)/float64(s.SensorValueMax-s.SensorValueMin))*40))
		gc.LineTo(float64(i*padding), height)
		gc.LineTo(float64(i*padding)+float64(padding), height)
	}
	gc.LineTo(float64(i*padding)+float64(padding), 40)
	gc.Close()
	gc.Fill()

	x := (int(s.Texture.Width) - (font.Width * 4))
	y := (40 - font.Height) / 2
	font.DrawString(data, x, y, fmt.Sprintf("%dC", s.SensorValue), color.RGBA{0x66, 0x66, 0x66, 0xff})

	s.Texture.Write(&data.Pix)
}

func (s *ThermalGraph) Run() {
	s.UpdateThermal()
	s.Redraw <- true

	ten := time.NewTicker(time.Second * 10)
	for {
		select {
		case <-ten.C:
			s.UpdateThermal()
			break
		}
		s.Redraw <- true
	}
}

func (s *ThermalGraph) UpdateThermal() {
	var max uint64 = 0

	for _, path := range s.Sensors {
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
