package widgets

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type State struct {
  X float32
  Y float32
  Width float64
  Height float64
  texture uint32
  vao uint32
  vbo uint32
  model mgl32.Mat4
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
        0.0, float32(w.Height),    0.0,  0.0, 0.0,
        float32(w.Width),   float32(w.Height),    0.0,  1.0, 0.0,
        float32(w.Width),   0.0,  0.0,  1.0, 1.0,
        0.0, 0.0,  0.0,  0.0, 1.0,
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
