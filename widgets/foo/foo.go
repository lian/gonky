package foo

import (
	"image"
	"image/color"
	"io/ioutil"

	"github.com/golang/freetype/truetype"
	"github.com/lian/gonky/terminus"
	"github.com/lian/gonky/texture"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"github.com/llgcode/draw2d/draw2dkit"
)

func loadFont() *truetype.Font {
	var data []byte
	var font *truetype.Font
	var err error

	if data, err = ioutil.ReadFile("./TerminusTTF-4.40.1.ttf"); err != nil {
		return nil
	}

	if font, err = truetype.Parse(data); err != nil {
		return nil
	}

	return font
}

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

	fontData := draw2d.FontData{
		Name:   "Terminus",
		Family: draw2d.FontFamilyMono,
		Style:  draw2d.FontStyleNormal,
	}
	draw2d.RegisterFont(fontData, loadFont())

	gc.SetFontData(fontData)
	//gc.SetFontSize(10)

	gc.SetFillColor(color.RGBA{0x00, 0x00, 0x00, 0xff})
	gc.FillStringAt("| 0.0-coinbase-0.0 |", 25, 25)

	//pixfont.DrawString(data, 20, 150, "Hello, World!", color.Black)
	//terminus.Font.DrawString(data, 20, 150, "Hello, World!", color.Black)
	//terminus.Font.DrawString(data, 20, 150, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789", color.Black)
	terminus.DrawString(data, 20, 150, "!#$%&'()*+,-./0123456789:;<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZ[\\]^_`abcdefghijklmnopqrstuvwxyz{|}~\"", color.Black)

	//terminus.DrawString(data, 20, 100, "hello world! { foo: 0 }", color.Black)

	//terminus.DrawString(data, 20, 80, "lianju\nBar", color.Black)
	terminus.DrawString(data, 20, 80, "Go's standard library provides strong support for \ninterpreting UTF-8 text. If a for range loop isn't sufficient for your purposes,\nchances are the facility you need is provided by a package in the library.", color.Black)

	s.Texture.Write(&data.Pix)
}
