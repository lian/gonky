package widgets

import (
	"fmt"
	"image"
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/llgcode/draw2d/draw2dimg"
)

type GL_State struct {
	texture      uint32
	vao          uint32
	vbo          uint32
	model        mgl32.Mat4
	modelUniform int32
}

type Base struct {
	Position  image.Point
	Geometry  image.Point
	_gl_state GL_State
}

func (widget *Base) gl_Draw() {
	gl.UniformMatrix4fv(widget._gl_state.modelUniform, 1, false, &widget._gl_state.model[0])
	gl.BindVertexArray(widget._gl_state.vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, widget._gl_state.texture)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}

func (widget *Base) gl_Setup(program, vertexAttr textureAttr, modelUni uint32) {
	gl.GenVertexArrays(1, &widget._gl_state.vao)
	gl.BindVertexArray(widget._gl_state.vao)

	planeVertices := []float32{
		//  X, Y, Z, U, V
		0.0, float32(widget.Geometry.Y), 0.0, 0.0, 0.0,
		float32(widget.Geometry.X), float32(widget.Geometry.Y), 0.0, 1.0, 0.0,
		float32(widget.Geometry.X), 0.0, 0.0, 1.0, 1.0,
		0.0, 0.0, 0.0, 0.0, 1.0,
	}

	gl.GenBuffers(1, &widget._gl_state.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, widget._gl_state.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(planeVertices)*4, gl.Ptr(planeVertices), gl.STATIC_DRAW)

	//vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertexAttr)
	gl.VertexAttribPointer(vertexAttr, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	//texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(textureAttr)
	gl.VertexAttribPointer(textureAttr, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	//w.model = mgl32.Ident4()
	w.model = mgl32.Translate3D(widget.Position.X, widget.Position.Y, 0.0)
	//w.modelUniform = gl.GetUniformLocation(program, gl.Str("model\x00"))
	w.modelUniform = modelUni
}

func (widget *Base) gl_ClearTexture() {
	if widget._gl_state.texture != 0 {
		gl.DeleteTextures(1, &widget._gl_state.texture)
		widget._gl_state.texture = 0
	}
}

func RenderWidget(widget Widget) {
	//surface := cairo.NewSurface(cairo.FORMAT_ARGB32, int(state.Width), int(state.Height))
	dest := image.NewRGBA(image.Rect(0, 0, int(widget._gl_state.Geometry.X), int(widget._gl_state.Geometry.Y)))
	gc := draw2dimg.NewGraphicContext(dest)

	widget.Render(gc)

	buf := gl.Ptr(dest.Pix)

	if widget._gl_state.texture == 0 {
		fmt.Printf("generate texture\n")
		widget._gl_state.texture = glTextureFromCairoData(int(widget._gl_state.Geometry.X), int(widget._gl_state.Geometry.Y), buf)
	} else {
		fmt.Printf("reuse texture\n")
		glSetTextureFromCairoData(int(widget._gl_state.Geometry.X), int(widget._gl_state.Geometry.Y), buf)
	}
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

type State struct {
	X            float32
	Y            float32
	Width        float64
	Height       float64
	texture      uint32
	vao          uint32
	vbo          uint32
	model        mgl32.Mat4
	modelUniform int32
}

func (s *State) ClearTexture() {
	if s.texture != 0 {
		gl.DeleteTextures(1, &s.texture)
		s.texture = 0
	}
}

func (w *State) Setup(program uint32) {
	gl.GenVertexArrays(1, &w.vao)
	gl.BindVertexArray(w.vao)

	planeVertices := []float32{
		//  X, Y, Z, U, V
		0.0, float32(w.Height), 0.0, 0.0, 0.0,
		float32(w.Width), float32(w.Height), 0.0, 1.0, 0.0,
		float32(w.Width), 0.0, 0.0, 1.0, 1.0,
		0.0, 0.0, 0.0, 0.0, 1.0,
	}

	gl.GenBuffers(1, &w.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, w.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(planeVertices)*4, gl.Ptr(planeVertices), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	//w.model = mgl32.Ident4()
	w.model = mgl32.Translate3D(w.X, w.Y, 0.0)
	w.modelUniform = gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(w.modelUniform, 1, false, &w.model[0])
}

func (s *State) Draw() {
	gl.UniformMatrix4fv(s.modelUniform, 1, false, &s.model[0])
	gl.BindVertexArray(s.vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, s.texture)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}
