package example

import (
	"fmt"
	"image/color"

	"github.com/lian/gonky/widgets"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

type Example struct {
	state *widgets.State
}

func (s *Example) State() *widgets.State {
	return s.state
}

func (s *Example) Render(gc *draw2dimg.GraphicContext) error {
	fmt.Printf("render example widget\n")

	gc.SetFillColor(color.RGBA{0x44, 0xff, 0x44, 0xff})
	draw2dkit.Rectangle(gc, 0, 0, s.state.Width, s.state.Height)
	gc.Fill()

	gc.SetFillColor(color.RGBA{0x00, 0xff, 0x22, 0xff})
	draw2dkit.Rectangle(gc, 10, 10, s.state.Width-20, s.state.Height-20)
	gc.Fill()

	return nil
}

func init() {
	widgets.Add("example", func() widgets.Widget {
		return &Example{
			state: &widgets.State{
				Width:  200,
				Height: 200,
				X:      200,
				Y:      200,
			},
		}
	})
}
