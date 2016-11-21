package shader

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

func DefaultShader() (*Program, error) {
	program, err := NewProgram(DefaultVertexShader, DefaultFragmentShader)
	return program, err
}

type Program struct {
	ID uint32
}

func (p *Program) Use() {
	gl.UseProgram(p.ID)
}

func (p *Program) AttributeLocation(name string) uint32 {
	return uint32(gl.GetAttribLocation(p.ID, gl.Str(name+"\x00")))
}

func (p *Program) UniformLocation(name string) int32 {
	return int32(gl.GetUniformLocation(p.ID, gl.Str(name+"\x00")))
}

func NewProgram(vertexShaderSource, fragmentShaderSource string) (*Program, error) {
	p := &Program{}
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return p, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return p, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return p, errors.New(fmt.Sprintf("failed to link program: %v", log))
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	p.ID = program

	return p, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)

	csource := gl.Str(source)
	gl.ShaderSource(shader, 1, &csource, nil)
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func SetupPerspective(width, height int, program *Program) {
	program.Use()

	fov := float32(60.0)
	eyeX := float32(width) / 2.0
	eyeY := float32(height) / 2.0
	ratio := float32(width) / float32(height)
	halfFov := (math.Pi * fov) / 360.0
	theTan := math.Tan(float64(halfFov))
	dist := eyeY / float32(theTan)
	nearDist := dist / 10.0
	farDist := dist * 10.0

	projection := mgl32.Perspective(mgl32.DegToRad(fov), ratio, nearDist, farDist)
	projectionUniform := program.UniformLocation("projection")
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{eyeX, eyeY, dist}, mgl32.Vec3{eyeX, eyeY, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := program.UniformLocation("camera")
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	//model := mgl32.Ident4()
	//modelUniform := program.UniformLocation("model")
	//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := program.UniformLocation("tex")
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program.ID, 0, gl.Str("outputColor\x00"))

	gl.Viewport(0, 0, int32(width), int32(height))
}

var DefaultVertexShader string = `
#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var DefaultFragmentShader = `
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"
