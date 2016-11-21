package foo

import (
	"image"
	"image/color"

	"github.com/lian/gonky/font/mono6x13"
	"github.com/lian/gonky/font/terminus"
	"github.com/lian/gonky/texture"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

type Foo struct {
	Texture *texture.Texture
}

func (s *Foo) Render() {
	data := image.NewRGBA(image.Rect(0, 0, int(s.Texture.Width), int(s.Texture.Height)))
	gc := draw2dimg.NewGraphicContext(data)

	gc.SetFillColor(color.RGBA{0xdd, 0xdd, 0xdd, 0xff})
	draw2dkit.Rectangle(gc, 0, 0, s.Texture.Width, s.Texture.Height)
	gc.Fill()

	/*
		gc.SetFillColor(color.RGBA{0x00, 0x00, 0xff, 0xff})
		draw2dkit.Rectangle(gc, 10, 10, s.Texture.Width-10, s.Texture.Height-10)
		gc.Fill()
	*/

	terminus.DrawString(data, 20, 10, "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"", color.Black)
	mono6x13.DrawString(data, 20, 30, "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"", color.Black)

	terminus.DrawString(data, 20, 50, "Go's standard library provides strong support for \ninterpreting UTF-8 text. If a for range loop isn't sufficient for your purposes,\nchances are the facility you need is provided by a package in the library.", color.Black)

	mono6x13.DrawString(data, 20, 100, "Go's standard library provides strong support for \ninterpreting UTF-8 text. If a for range loop isn't sufficient for your purposes,\nchances are the facility you need is provided by a package in the library.", color.Black)

	s.Texture.Write(&data.Pix)
}
