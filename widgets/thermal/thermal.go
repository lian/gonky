package thermal

import (
	"fmt"
	"image"
	"image/color"

	"github.com/lian/gonky/shader"
	"github.com/lian/gonky/texture"
	"github.com/lian/gonky/widgets"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"

	font "github.com/lian/gonky/font/terminus"
)

type Graphs struct {
	Texture *texture.Texture
	Redraw  chan bool

	GraphPadding int
	Stats        *widgets.Stats
}

func New(program *shader.Program, stats *widgets.Stats) *Graphs {
	s := &Graphs{
		Texture:      &texture.Texture{X: 20, Y: 768 - (18 * 2), Width: 300, Height: 200},
		Redraw:       make(chan bool),
		GraphPadding: 8,
		Stats:        stats,
	}
	s.Texture.Setup(program)
	return s
}

func (s *Graphs) Render() {
	data := image.NewRGBA(image.Rect(0, 0, int(s.Texture.Width), int(s.Texture.Height)))
	gc := draw2dimg.NewGraphicContext(data)

	gc.SetFillColor(color.RGBA{0x33, 0x33, 0x33, 0xff})
	draw2dkit.Rectangle(gc, 0, 0, s.Texture.Width, s.Texture.Height)
	gc.Fill()

	//gc.SetFillColor(color.RGBA{0x66, 0x66, 0x66, 0xff})
	gc.SetStrokeColor(color.RGBA{0x66, 0x66, 0x66, 0xff})
	gc.SetLineWidth(1.0)

	s.DrawThermal(gc, data)
	s.DrawFan(gc, data)

	s.Texture.Write(&data.Pix)
}

func (s *Graphs) DrawThermal(gc *draw2dimg.GraphicContext, data *image.RGBA) {
	padding := s.GraphPadding
	graphHeight := 40.0
	yOffset := 0.0

	maxItems := (int(s.Texture.Width) - (font.Width * 5)) / padding
	start := len(s.Stats.ThermalGraph) - maxItems
	if start < 0 {
		start = 0
	}

	//gc.MoveTo(0, graphHeight+yOffset)
	var i, value int
	for i, value = range s.Stats.ThermalGraph[start:] {
		scaled := graphHeight - float64(int((float64(value-s.Stats.ThermalValueMin)/float64(s.Stats.ThermalValueMax-s.Stats.ThermalValueMin))*graphHeight))
		height := scaled + float64(yOffset)
		if i == 0 {
			gc.MoveTo(float64(i*padding), height)
		} else {
			gc.LineTo(float64(i*padding), height)
		}
		gc.LineTo(float64(i*padding)+float64(padding), height)
	}
	//gc.LineTo(float64(i*padding)+float64(padding), graphHeight+yOffset)
	//gc.Close()
	//gc.Fill()
	gc.Stroke()

	x := (int(s.Texture.Width) - (font.Width * 4))
	y := int(yOffset + ((graphHeight - font.Height) / 2))
	font.DrawString(data, x, y, fmt.Sprintf("%dC", s.Stats.ThermalValue), color.RGBA{0x66, 0x66, 0x66, 0xff})
}

func (s *Graphs) DrawFan(gc *draw2dimg.GraphicContext, data *image.RGBA) {
	padding := s.GraphPadding
	graphHeight := 40.0
	yOffset := 60.0

	maxItems := (int(s.Texture.Width) - (font.Width * 13)) / padding
	start := len(s.Stats.FanGraph) - maxItems
	if start < 0 {
		start = 0
	}

	//gc.MoveTo(0, graphHeight+yOffset)
	var i, value int
	for i, value = range s.Stats.FanGraph[start:] {
		scaled := graphHeight - float64(int((float64(value-s.Stats.FanValueMin)/float64(s.Stats.FanValueMax-s.Stats.FanValueMin))*graphHeight))
		height := scaled + float64(yOffset)
		if i == 0 {
			gc.MoveTo(float64(i*padding), height)
		} else {
			gc.LineTo(float64(i*padding), height)
		}
		gc.LineTo(float64(i*padding)+float64(padding), height)
	}
	//gc.LineTo(float64(i*padding)+float64(padding), graphHeight+yOffset)
	//gc.Close()
	//gc.Fill()
	gc.Stroke()

	x := (int(s.Texture.Width) - (font.Width * 12))
	y := int(yOffset + ((graphHeight - font.Height) / 2))
	font.DrawString(data, x, y, fmt.Sprintf("%d RPM L%d", s.Stats.FanValue, s.Stats.FanLevel), color.RGBA{0x66, 0x66, 0x66, 0xff})
}
