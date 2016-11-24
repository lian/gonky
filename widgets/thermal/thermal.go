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
	Texture      *texture.Texture
	Redraw       chan bool
	Sensors      []string
	SensorValue  uint64
	SensorsGraph []uint64
	GraphPadding int
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
		SensorsGraph: []uint64{},
		GraphPadding: 10,
	}
	s.Texture.Setup(program)
	return s
}

func (s *ThermalGraph) Render() {
	data := image.NewRGBA(image.Rect(0, 0, int(s.Texture.Width), int(s.Texture.Height)))
	gc := draw2dimg.NewGraphicContext(data)

	gc.SetFillColor(color.RGBA{0x33, 0x33, 0x33, 0xff})
	//gc.SetFillColor(color.RGBA{0x66, 0x66, 0x66, 0xff})
	draw2dkit.Rectangle(gc, 0, 0, s.Texture.Width, s.Texture.Height)
	gc.Fill()

	padding := s.GraphPadding
	//gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.SetFillColor(color.RGBA{0x66, 0x66, 0x66, 0xff})
	//gc.SetFillColor(color.RGBA{0x33, 0x33, 0x33, 0xff})

	/*
		for i, v := range s.SensorsGraph {
			//fmt.Println(i, v)
			height := (float64(v[0]) / 100.0) * 40
			draw2dkit.Rectangle(gc, float64(i*padding), float64(int(height)), float64(i*padding)+float64(padding), 40)
			gc.Fill()
		}
	*/

	/*
		//gc.SetLineWidth(1.0)
		gc.Current.LineWidth = 10
		for i, v := range s.SensorsGraph {
			height := float64(int((float64(v[0]) / 100.0) * 40))
			if i == 0 {
				gc.MoveTo(float64(i*padding), height)
			} else {
				gc.LineTo(float64(i*padding), height)
			}
			gc.LineTo(float64(i*padding)+float64(padding), height)
		}
		gc.Stroke()
	*/

	gc.MoveTo(0, 40)
	var i int
	var value uint64
	for i, value = range s.SensorsGraph {
		height := float64(int((float64(value) / 100.0) * 40))
		gc.LineTo(float64(i*padding), height)
		gc.LineTo(float64(i*padding)+float64(padding), height)
	}
	gc.LineTo(float64(i*padding)+float64(padding), 40)
	gc.Close()
	gc.Fill()

	x := (int(s.Texture.Width) - (font.Width * 4))
	y := (40 - font.Height) / 2
	//font.DrawString(data, x, y, fmt.Sprintf("%dC", s.SensorValue), color.Black)
	//font.DrawString(data, x, y, fmt.Sprintf("%dC", s.SensorValue), color.White)
	font.DrawString(data, x, y, fmt.Sprintf("%dC", s.SensorValue), color.RGBA{0x66, 0x66, 0x66, 0xff})

	s.Texture.Write(&data.Pix)
}

func (s *ThermalGraph) Run() {
	s.UpdateThermal()
	s.Redraw <- true

	ten := time.NewTicker(time.Second * 2)
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

	s.SensorValue = max / 1000

	maxItems := (int(s.Texture.Width) - (font.Width * 5)) / s.GraphPadding
	if len(s.SensorsGraph) >= maxItems {
		s.SensorsGraph = append(s.SensorsGraph[1:], s.SensorValue)
	} else {
		s.SensorsGraph = append(s.SensorsGraph, s.SensorValue)
	}
}
