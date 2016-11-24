package texture

import (
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/lian/gonky/shader"
)

type Texture struct {
	//ID     string
	X      float64
	Y      float64
	Width  float64
	Height float64

	texture      uint32
	vao          uint32
	vbo          uint32
	model        mgl32.Mat4
	modelUniform int32
	Program      *shader.Program
}

func (t *Texture) Setup(program *shader.Program) {
	t.Program = program
	vertexAttrLocation := t.Program.AttributeLocation("vert")
	textureAttrLocation := t.Program.AttributeLocation("vertTexCoord")
	modelUniformLocation := t.Program.UniformLocation("model")

	gl.GenVertexArrays(1, &t.vao)
	gl.BindVertexArray(t.vao)

	planeVertices := []float32{
		//  X, Y, Z, U, V
		0.0, float32(t.Height), 0.0, 0.0, 0.0,
		float32(t.Width), float32(t.Height), 0.0, 1.0, 0.0,
		float32(t.Width), 0.0, 0.0, 1.0, 1.0,
		0.0, 0.0, 0.0, 0.0, 1.0,
	}

	gl.GenBuffers(1, &t.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, t.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(planeVertices)*4, gl.Ptr(planeVertices), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(vertexAttrLocation)
	gl.VertexAttribPointer(vertexAttrLocation, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	gl.EnableVertexAttribArray(textureAttrLocation)
	gl.VertexAttribPointer(textureAttrLocation, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	t.model = mgl32.Translate3D(float32(t.X), float32(t.Y-t.Height), 0.0)
	t.modelUniform = modelUniformLocation
}

func (t *Texture) Draw() {
	gl.UniformMatrix4fv(t.modelUniform, 1, false, &t.model[0])
	gl.BindVertexArray(t.vao)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.texture)
	gl.DrawArrays(gl.TRIANGLE_FAN, 0, 4)
}

func (t *Texture) Clear() {
	if t.texture != 0 {
		gl.DeleteTextures(1, &t.texture)
		t.texture = 0
	}
}

func (t *Texture) Write(data *[]uint8) {
	buf := gl.Ptr(*data)

	if t.texture == 0 {
		t.texture = newTextureData(int32(t.Width), int32(t.Height), buf)
	} else {
		gl.BindTexture(gl.TEXTURE_2D, t.texture)
		updateTextureData(int32(t.Width), int32(t.Height), buf)
	}
}

func newTextureData(width, height int32, data unsafe.Pointer) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.REPEAT)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	return texture
}

func updateTextureData(width, height int32, data unsafe.Pointer) {
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, width, height, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)
}
