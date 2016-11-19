package widgets

import "github.com/llgcode/draw2d/draw2dimg"

type Widget interface {
	State() *State
	Render(gc *draw2dimg.GraphicContext) error
}
type Creator func() Widget

var Widgets = map[string]Creator{}

func Add(name string, creator Creator) {
	Widgets[name] = creator
}
