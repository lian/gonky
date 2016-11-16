package widgets

import (
    "github.com/ungerik/go-cairo"
)

type Widget interface {
	State() *State
	Render(surface *cairo.Surface) error
}
type Creator func() Widget

var Widgets = map[string]Creator{}

func Add(name string, creator Creator) {
	Widgets[name] = creator
}
