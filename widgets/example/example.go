package example

import (
  "fmt"
  "github.com/lian/gonky/widgets"
  "github.com/ungerik/go-cairo"
)

type Example struct {
  state *widgets.State
}

func (s *Example) State() *widgets.State {
  return s.state
}

func (s *Example) Render(surface *cairo.Surface) error {
  fmt.Printf("render example widget\n")

  surface.SetSourceRGB(0.3, 0.3, 0.3)
  surface.Rectangle(0, 0, s.state.Width, s.state.Height)
  surface.Fill()

  surface.SetSourceRGB(0.5, 0.5, 0.5)
  surface.Rectangle(10, 10, s.state.Width-20, s.state.Height-20)
  surface.Fill()

  return nil
}

func init() {
    widgets.Add("example", func() widgets.Widget {
        return &Example{
          state: &widgets.State{
            Width: 200,
            Height: 200,
            X: 200,
            Y: 200,
          },
        }
    })
}
