package foo

import (
	"image"
	"image/color"

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

	gc.SetFillColor(color.RGBA{0xff, 0x00, 0x00, 0xff})
	draw2dkit.Rectangle(gc, 0, 0, s.Texture.Width, s.Texture.Height)
	gc.Fill()

	gc.SetFillColor(color.RGBA{0x00, 0x00, 0xff, 0xff})
	draw2dkit.Rectangle(gc, 10, 10, s.Texture.Width-10, s.Texture.Height-10)
	gc.Fill()

	s.Texture.Write(&data.Pix)
}
