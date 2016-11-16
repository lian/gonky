package widgets

import (
	"fmt"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/ungerik/go-cairo"
)

func NewManger(defaultProgram uint32) *Manager {
	return &Manager{
		DefaultProgram: defaultProgram,
		Widgets:        make(map[string]Widget),
	}
}

type Manager struct {
	DefaultProgram uint32
	States         map[string]*State
	Widgets        map[string]Widget
}

func (m *Manager) Add(name string, creator Creator) {
	widget := creator()
	state := widget.State()
	state.Setup(m.DefaultProgram)
	m.Widgets[name] = widget
}

func (m *Manager) Draw() {
	for _, widget := range m.Widgets {
		widget.State().Draw()
	}
}

func (m *Manager) Render() {
	for _, widget := range m.Widgets {
		m.RenderWidget(widget)
	}
}

func (m *Manager) RenderWidget(widget Widget) {
	state := widget.State()

	surface := cairo.NewSurface(cairo.FORMAT_ARGB32, int(state.Width), int(state.Height))

	widget.Render(surface)

	buf := gl.Ptr(surface.GetData())

	if state.texture == 0 {
		fmt.Printf("generate texture\n")
		state.texture = glTextureFromCairoData(int(state.Width), int(state.Height), buf)
	} else {
		fmt.Printf("reuse texture\n")
		glSetTextureFromCairoData(int(state.Width), int(state.Height), buf)
	}

	surface.Destroy()
}

func glTextureFromCairoData(width int, height int, data unsafe.Pointer) uint32 {
	//glValidTextureSize(width, height)
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.BGRA, gl.UNSIGNED_BYTE, data)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

func glSetTextureFromCairoData(width int, height int, data unsafe.Pointer) {
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(width), int32(height), 0, gl.BGRA, gl.UNSIGNED_BYTE, data)
}

/*
func glValidTextureSize(width, height int) bool {
  var maxSize int32
  gl.GetIntegerv(gl.MAX_TEXTURE_SIZE, &maxSize)
  if (int32(width) < maxSize) && (int32(height) < maxSize) {
    return true
  }

  var texWidth int32
  gl.GetTexLevelParameteriv(gl.PROXY_TEXTURE_2D, 0, gl.TEXTURE_WIDTH, &texWidth)
  if texWidth > 0 {
    return true
  }

  fmt.Printf("TEXTURE TOO LARGE for OpenGL! width:%d height:%d maxSize:%d texWidth:%d\n", width, height, maxSize, texWidth)
  return false
}
*/
