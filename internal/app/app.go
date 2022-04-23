package app

import (
	"fmt"
	"runtime"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/veandco/go-sdl2/sdl"
)

var ctx struct {
	window    *sdl.Window
	time      uint64
	glContext sdl.GLContext
}

func Init(title string, w, h int32) error {
	runtime.LockOSThread()

	//sdl.SetHint(sdl.HINT_VIDEO_HIGHDPI_DISABLED, "0")
	err := sdl.Init(sdl.INIT_VIDEO)
	if err != nil {
		return fmt.Errorf("failed to initialize SDL2: %w", err)
	}

	window, err := sdl.CreateWindow(
		title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		w, h,
		sdl.WINDOW_OPENGL|sdl.WINDOW_SHOWN|sdl.WINDOW_ALLOW_HIGHDPI|sdl.WINDOW_RESIZABLE,
	)
	if err != nil {
		return fmt.Errorf("failed to create window: %w", err)
	}
	ctx.window = window

	sdl.GLSetAttribute(sdl.GL_CONTEXT_MAJOR_VERSION, 3)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_MINOR_VERSION, 2)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_FLAGS, sdl.GL_CONTEXT_FORWARD_COMPATIBLE_FLAG)
	sdl.GLSetAttribute(sdl.GL_CONTEXT_PROFILE_MASK, sdl.GL_CONTEXT_PROFILE_CORE)
	sdl.GLSetAttribute(sdl.GL_DOUBLEBUFFER, 1)
	sdl.GLSetAttribute(sdl.GL_DEPTH_SIZE, 24)
	sdl.GLSetAttribute(sdl.GL_STENCIL_SIZE, 8)

	glContext, err := window.GLCreateContext()
	if err != nil {
		return fmt.Errorf("failed to create OpenGL context: %w", err)
	}
	err = window.GLMakeCurrent(glContext)
	if err != nil {
		return fmt.Errorf("failed to set current OpenGL context: %w", err)
	}
	ctx.glContext = glContext

	sdl.GLSetSwapInterval(1)

	err = gl.Init()
	if err != nil {
		return fmt.Errorf("failed to initialize OpenGL: %w", err)
	}

	return nil
}

func Dispose() {
	sdl.GLDeleteContext(ctx.glContext)
	ctx.window.Destroy()
	sdl.Quit()
}

func DisplaySize() [2]float32 {
	w, h := ctx.window.GetSize()
	return [2]float32{float32(w), float32(h)}
}

func FramebufferSize() [2]float32 {
	w, h := ctx.window.GLGetDrawableSize()
	return [2]float32{float32(w), float32(h)}
}

func GetDeltaTime() float32 {
	// Don't use SDL_GetTicks() because it is using millisecond resolution
	frequency := sdl.GetPerformanceFrequency()
	currentTime := sdl.GetPerformanceCounter()

	var dt float32 = 1.0 / 60.0

	if ctx.time > 0 {
		dt = float32(currentTime-ctx.time) / float32(frequency)
	}
	ctx.time = currentTime

	return dt
}

func PreRender() {
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
}

func PostRender() {
	ctx.window.GLSwap()
}
