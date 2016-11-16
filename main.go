package main

import (
	"fmt"
	"log"
	"runtime"

	"math"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.1/glfw"
	"github.com/go-gl/mathgl/mgl32"

	"github.com/lian/gonky/shader"

	"github.com/lian/gonky/widgets"
	_ "github.com/lian/gonky/widgets/example"
)

func init() {
	runtime.LockOSThread()
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	fmt.Printf("%v %d, %v %v\n", key, scancode, action, mods)
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
}

const WindowWidth = 800
const WindowHeight = 600

func setupPerspective(width, height int, program uint32) {
	fov := float32(60.0)
	eyeX := float32(WindowWidth) / 2.0
	eyeY := float32(WindowHeight) / 2.0
	ratio := float32(WindowWidth) / float32(WindowHeight)
	halfFov := (math.Pi * fov) / 360.0
	theTan := math.Tan(float64(halfFov))
	dist := eyeY / float32(theTan)
	nearDist := dist / 10.0
	farDist := dist * 10.0

	projection := mgl32.Perspective(mgl32.DegToRad(fov), ratio, nearDist, farDist)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{eyeX, eyeY, dist}, mgl32.Vec3{eyeX, eyeY, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	//model := mgl32.Ident4()
	//modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	//gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))
}

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	//glfw.WindowHint(glfw.Samples, 4)

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Derp", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	window.SetKeyCallback(keyCallback)
	//window.SetCursorPosCallback(cursorPosCallback)
	//window.SetScrollCallback(scrollCallback)
	//window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	program, err := shader.DefaultShader()
	if err != nil {
		panic(err)
	}
	fmt.Printf("program: %d\n", program)
	gl.UseProgram(program)

	setupPerspective(WindowWidth, WindowHeight, program)

	manager := widgets.NewManger(program)

	for name, creator := range widgets.Widgets {
		manager.Add(name, creator)
	}

	manager.Render()

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Render
		gl.UseProgram(program)
		manager.Draw()

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}
