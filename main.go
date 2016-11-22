package main

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"

	"github.com/lian/gonky/shader"
	"github.com/lian/gonky/texture"
	"github.com/lian/gonky/widgets/foo"
	"github.com/lian/gonky/widgets/status"
)

func init() {
	runtime.LockOSThread()
}

var redrawChan chan bool = make(chan bool, 10)

const redrawChanHalfLen = 5

func triggerRedraw() {
	if len(redrawChan) < redrawChanHalfLen {
		redrawChan <- true
	}
}

func keyCallback(window *glfw.Window, key glfw.Key, scancode int, action glfw.Action, mods glfw.ModifierKey) {
	//fmt.Printf("%v %d, %v %v\n", key, scancode, action, mods)
	if key == glfw.KeyEscape && action == glfw.Press {
		window.SetShouldClose(true)
	}
	triggerRedraw()
}

func focusCallback(window *glfw.Window, focused bool) {
	//fmt.Println("focus:", focused)
	triggerRedraw()
}

func refreshCallback(window *glfw.Window) {
	//fmt.Println("refreshCallback")
	triggerRedraw()
}

func resizeCallback(w *glfw.Window, width int, height int) {
	//fmt.Println("RESIZE", width, height)
	WindowWidth = width
	WindowHeight = height
	shader.SetupPerspective(width, height, program)
}

var WindowWidth int = 1366
var WindowHeight int = 768

var program *shader.Program

func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.Resizable, glfw.True)
	//glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	//glfw.WindowHint(glfw.Samples, 4)

	window, err := glfw.CreateWindow(WindowWidth, WindowHeight, "Derp", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()
	//window.SetInputMode(glfw.CursorMode, glfw.CursorDisabled)
	//window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)

	window.SetSizeCallback(resizeCallback)

	window.SetRefreshCallback(refreshCallback)
	window.SetFocusCallback(focusCallback)
	window.SetKeyCallback(keyCallback)
	//window.SetCursorPosCallback(cursorPosCallback)
	//window.SetScrollCallback(scrollCallback)

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	program, err = shader.DefaultShader()
	if err != nil {
		panic(err)
	}
	//fmt.Printf("program: %v\n", program)
	program.Use()

	shader.SetupPerspective(WindowWidth, WindowHeight, program)

	foo := &foo.Foo{
		Texture: &texture.Texture{X: 20, Y: 20, Width: 1024, Height: 256},
	}

	foo.Texture.Setup(program)
	foo.Render()

	status := status.New(0, 768-18, 1366, 18, program)
	go status.Run()

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(0.2, 0.2, 0.2, 1.0)

	pollEventsTimer := time.NewTicker(time.Millisecond * 100)
	maxRenderDelayTimer := time.NewTicker(time.Second * 20)

	for !window.ShouldClose() {
		select {
		case <-pollEventsTimer.C:
			glfw.PollEvents()
			continue
		case <-status.Redraw:
			status.Render()
		case <-maxRenderDelayTimer.C:
			//fmt.Println("max delay tick")
		case <-redrawChan:
			//fmt.Println("redraw tick")
		}

		//fmt.Println("DRAW")
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		program.Use()
		foo.Texture.Draw()
		status.Texture.Draw()

		window.SwapBuffers()
		glfw.PollEvents()
	}
}
