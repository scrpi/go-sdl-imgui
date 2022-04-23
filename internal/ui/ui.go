package ui

import (
	"github.com/inkyblackness/imgui-go/v4"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	buttonLeft   = 0
	buttonRight  = 1
	buttonMiddle = 2
	buttonMax_   = 3
)

var mousePressed [buttonMax_]bool

func Init() {
	imgui.CreateContext(nil)
	setKeyMapping()

	io := imgui.CurrentIO()
	io.SetBackendFlags(io.GetBackendFlags() | imgui.BackendFlagsRendererHasVtxOffset)

	rendererInit()
}

func Dispose() {
	rendererDispose()

	ctx, err := imgui.CurrentContext()
	if err != nil {
		ctx.Destroy()
	}
}

func setKeyMapping() {
	keys := map[int]int{
		imgui.KeyTab:        sdl.SCANCODE_TAB,
		imgui.KeyLeftArrow:  sdl.SCANCODE_LEFT,
		imgui.KeyRightArrow: sdl.SCANCODE_RIGHT,
		imgui.KeyUpArrow:    sdl.SCANCODE_UP,
		imgui.KeyDownArrow:  sdl.SCANCODE_DOWN,
		imgui.KeyPageUp:     sdl.SCANCODE_PAGEUP,
		imgui.KeyPageDown:   sdl.SCANCODE_PAGEDOWN,
		imgui.KeyHome:       sdl.SCANCODE_HOME,
		imgui.KeyEnd:        sdl.SCANCODE_END,
		imgui.KeyInsert:     sdl.SCANCODE_INSERT,
		imgui.KeyDelete:     sdl.SCANCODE_DELETE,
		imgui.KeyBackspace:  sdl.SCANCODE_BACKSPACE,
		imgui.KeySpace:      sdl.SCANCODE_BACKSPACE,
		imgui.KeyEnter:      sdl.SCANCODE_RETURN,
		imgui.KeyEscape:     sdl.SCANCODE_ESCAPE,
		imgui.KeyA:          sdl.SCANCODE_A,
		imgui.KeyC:          sdl.SCANCODE_C,
		imgui.KeyV:          sdl.SCANCODE_V,
		imgui.KeyX:          sdl.SCANCODE_X,
		imgui.KeyY:          sdl.SCANCODE_Y,
		imgui.KeyZ:          sdl.SCANCODE_Z,
	}

	io := imgui.CurrentIO()

	// Keyboard mapping. ImGui will use those indices to peek into the io.KeysDown[] array.
	for imguiKey, nativeKey := range keys {
		io.KeyMap(imguiKey, nativeKey)
	}
}

func ProcessEvent(event sdl.Event) (bool, bool) {
	io := imgui.CurrentIO()
	switch event.GetType() {
	case sdl.MOUSEBUTTONDOWN:
		buttonEvent := event.(*sdl.MouseButtonEvent)
		switch buttonEvent.Button {
		case sdl.BUTTON_LEFT:
			mousePressed[buttonLeft] = true
		case sdl.BUTTON_RIGHT:
			mousePressed[buttonRight] = true
		case sdl.BUTTON_MIDDLE:
			mousePressed[buttonMiddle] = true
		}
	case sdl.MOUSEWHEEL:
		wheelEvent := event.(*sdl.MouseWheelEvent)
		var deltaX, deltaY float32
		if wheelEvent.X > 0 {
			deltaX++
		}
		if wheelEvent.X < 0 {
			deltaX--
		}
		if wheelEvent.Y > 0 {
			deltaY++
		}
		if wheelEvent.Y < 0 {
			deltaY--
		}
		io.AddMouseWheelDelta(deltaX, deltaY)
	case sdl.TEXTINPUT:
		inputEvent := event.(*sdl.TextInputEvent)
		io.AddInputCharacters(string(inputEvent.Text[:]))
	case sdl.KEYDOWN:
		keyEvent := event.(*sdl.KeyboardEvent)
		io.KeyPress(int(keyEvent.Keysym.Scancode))
		updateKeyModifier()
	case sdl.KEYUP:
		keyEvent := event.(*sdl.KeyboardEvent)
		io.KeyRelease(int(keyEvent.Keysym.Scancode))
		updateKeyModifier()
	}

	return io.WantCaptureMouse(), io.WantCaptureKeyboard()
}

func updateKeyModifier() {
	modState := sdl.GetModState()
	mapModifier := func(lMask sdl.Keymod, lKey int, rMask sdl.Keymod, rKey int) (lResult int, rResult int) {
		if (modState & lMask) != 0 {
			lResult = lKey
		}
		if (modState & rMask) != 0 {
			rResult = rKey
		}
		return
	}
	io := imgui.CurrentIO()
	io.KeyShift(mapModifier(sdl.KMOD_LSHIFT, sdl.SCANCODE_LSHIFT, sdl.KMOD_RSHIFT, sdl.SCANCODE_RSHIFT))
	io.KeyCtrl(mapModifier(sdl.KMOD_LCTRL, sdl.SCANCODE_LCTRL, sdl.KMOD_RCTRL, sdl.SCANCODE_RCTRL))
	io.KeyAlt(mapModifier(sdl.KMOD_LALT, sdl.SCANCODE_LALT, sdl.KMOD_RALT, sdl.SCANCODE_RALT))
}

func NewFrame(dt float32, displaySize [2]float32) {
	io := imgui.CurrentIO()

	io.SetDisplaySize(imgui.Vec2{displaySize[0], displaySize[1]})
	io.SetDeltaTime(dt)

	// If a mouse press event came, always pass it as "mouse held this frame", so we don't miss click-release events that are shorter than 1 frame.
	x, y, state := sdl.GetMouseState()
	io.SetMousePosition(imgui.Vec2{X: float32(x), Y: float32(y)})

	io.SetMouseButtonDown(buttonLeft, mousePressed[buttonLeft] || (state&sdl.Button(sdl.BUTTON_LEFT)) != 0)
	io.SetMouseButtonDown(buttonRight, mousePressed[buttonRight] || (state&sdl.Button(sdl.BUTTON_RIGHT)) != 0)
	io.SetMouseButtonDown(buttonMiddle, mousePressed[buttonMiddle] || (state&sdl.Button(sdl.BUTTON_MIDDLE)) != 0)
	for i := 0; i < buttonMax_; i++ {
		mousePressed[i] = false
	}

	imgui.NewFrame()
}
