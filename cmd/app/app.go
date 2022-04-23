package main

import (
	"log"
	"time"

	"github.com/inkyblackness/imgui-go/v4"
	"github.com/scrpi/go-sdl-imgui/internal/app"
	"github.com/scrpi/go-sdl-imgui/internal/ui"
	"github.com/veandco/go-sdl2/sdl"
)

func handleInput(running *bool) {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		uiCapturedMouse, uiCapturedKB := ui.ProcessEvent(event)
		_ = uiCapturedKB
		_ = uiCapturedMouse

		switch event.GetType() {
		case sdl.QUIT:
			*running = false
		}
	}
}

func main() {
	err := app.Init("Golang SDL ImGUI Template", 1200, 800)
	if err != nil {
		log.Fatalf("failed to initialize OpenGL: %v", err)
	}
	defer app.Dispose()

	ui.Init()
	defer ui.Dispose()

	running := true
	showDemoWindow := true

	for running {
		dt := app.GetDeltaTime()

		handleInput(&running)
		if !running {
			break
		}

		ui.NewFrame(dt, app.DisplaySize())

		imgui.ShowDemoWindow(&showDemoWindow)

		// Creates the draw data list
		imgui.Render()

		app.PreRender()
		// Application performs it's own rendering here...
		// xyz.RenderScene()
		ui.Render(app.DisplaySize(), app.FramebufferSize(), imgui.RenderedDrawData())
		app.PostRender()

		<-time.After(time.Millisecond * 5)
	}
}
