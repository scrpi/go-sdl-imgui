package ui

import (
	_ "embed"

	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/inkyblackness/imgui-go/v4"
)

//go:embed ui-shader/main.vert
var vertexShader string

//go:embed ui-shader/main.frag
var fragmentShader string

type rendererContext struct {
	fontTexture            uint32
	shaderHandle           uint32
	vertHandle             uint32
	fragHandle             uint32
	attribLocationTex      int32
	attribLocationProjMtx  int32
	attribLocationPosition int32
	attribLocationUV       int32
	attribLocationColor    int32
	vboHandle              uint32
	elementsHandle         uint32
}

var ctx rendererContext

func rendererInit() {
	// Backup GL state
	var lastTexture int32
	var lastArrayBuffer int32
	var lastVertexArray int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)

	ctx.shaderHandle = gl.CreateProgram()
	ctx.vertHandle = gl.CreateShader(gl.VERTEX_SHADER)
	ctx.fragHandle = gl.CreateShader(gl.FRAGMENT_SHADER)

	glShaderSource := func(handle uint32, source string) {
		csource, free := gl.Strs(source + "\x00")
		defer free()

		gl.ShaderSource(handle, 1, csource, nil)
	}

	glShaderSource(ctx.vertHandle, vertexShader)
	glShaderSource(ctx.fragHandle, fragmentShader)
	gl.CompileShader(ctx.vertHandle)
	gl.CompileShader(ctx.fragHandle)
	gl.AttachShader(ctx.shaderHandle, ctx.vertHandle)
	gl.AttachShader(ctx.shaderHandle, ctx.fragHandle)
	gl.LinkProgram(ctx.shaderHandle)

	ctx.attribLocationTex = gl.GetUniformLocation(ctx.shaderHandle, gl.Str("Texture"+"\x00"))
	ctx.attribLocationProjMtx = gl.GetUniformLocation(ctx.shaderHandle, gl.Str("ProjMtx"+"\x00"))
	ctx.attribLocationPosition = gl.GetAttribLocation(ctx.shaderHandle, gl.Str("Position"+"\x00"))
	ctx.attribLocationUV = gl.GetAttribLocation(ctx.shaderHandle, gl.Str("UV"+"\x00"))
	ctx.attribLocationColor = gl.GetAttribLocation(ctx.shaderHandle, gl.Str("Color"+"\x00"))

	gl.GenBuffers(1, &ctx.vboHandle)
	gl.GenBuffers(1, &ctx.elementsHandle)

	createFontsTexture()

	// Restore modified GL state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindVertexArray(uint32(lastVertexArray))
}

func rendererDispose() {
	if ctx.vboHandle != 0 {
		gl.DeleteBuffers(1, &ctx.vboHandle)
	}
	ctx.vboHandle = 0
	if ctx.elementsHandle != 0 {
		gl.DeleteBuffers(1, &ctx.elementsHandle)
	}
	ctx.elementsHandle = 0

	if (ctx.shaderHandle != 0) && (ctx.vertHandle != 0) {
		gl.DetachShader(ctx.shaderHandle, ctx.vertHandle)
	}
	if ctx.vertHandle != 0 {
		gl.DeleteShader(ctx.vertHandle)
	}
	ctx.vertHandle = 0

	if (ctx.shaderHandle != 0) && (ctx.fragHandle != 0) {
		gl.DetachShader(ctx.shaderHandle, ctx.fragHandle)
	}
	if ctx.fragHandle != 0 {
		gl.DeleteShader(ctx.fragHandle)
	}
	ctx.fragHandle = 0

	if ctx.shaderHandle != 0 {
		gl.DeleteProgram(ctx.shaderHandle)
	}
	ctx.shaderHandle = 0

	if ctx.fontTexture != 0 {
		gl.DeleteTextures(1, &ctx.fontTexture)
		imgui.CurrentIO().Fonts().SetTextureID(0)
		ctx.fontTexture = 0
	}
}

func createFontsTexture() {
	// Build texture atlas
	io := imgui.CurrentIO()
	image := io.Fonts().TextureDataAlpha8()

	// Upload texture to graphics system
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	gl.GenTextures(1, &ctx.fontTexture)
	gl.BindTexture(gl.TEXTURE_2D, ctx.fontTexture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.PixelStorei(gl.UNPACK_ROW_LENGTH, 0)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RED, int32(image.Width), int32(image.Height),
		0, gl.RED, gl.UNSIGNED_BYTE, image.Pixels)

	// Store our identifier
	io.Fonts().SetTextureID(imgui.TextureID(ctx.fontTexture))

	// Restore state
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
}

// Render translates the ImGui draw data to ctx3 commands.
func Render(displaySize [2]float32, framebufferSize [2]float32, drawData imgui.DrawData) {
	// Avoid rendering when minimized, scale coordinates for retina displays (screen coordinates != framebuffer coordinates)
	displayWidth, displayHeight := displaySize[0], displaySize[1]
	fbWidth, fbHeight := framebufferSize[0], framebufferSize[1]
	if (fbWidth <= 0) || (fbHeight <= 0) {
		return
	}
	drawData.ScaleClipRects(imgui.Vec2{
		X: fbWidth / displayWidth,
		Y: fbHeight / displayHeight,
	})

	// Backup GL state
	var lastActiveTexture int32
	gl.GetIntegerv(gl.ACTIVE_TEXTURE, &lastActiveTexture)
	gl.ActiveTexture(gl.TEXTURE0)
	var lastProgram int32
	gl.GetIntegerv(gl.CURRENT_PROGRAM, &lastProgram)
	var lastTexture int32
	gl.GetIntegerv(gl.TEXTURE_BINDING_2D, &lastTexture)
	var lastSampler int32
	gl.GetIntegerv(gl.SAMPLER_BINDING, &lastSampler)
	var lastArrayBuffer int32
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &lastArrayBuffer)
	var lastElementArrayBuffer int32
	gl.GetIntegerv(gl.ELEMENT_ARRAY_BUFFER_BINDING, &lastElementArrayBuffer)
	var lastVertexArray int32
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &lastVertexArray)
	var lastPolygonMode [2]int32
	gl.GetIntegerv(gl.POLYGON_MODE, &lastPolygonMode[0])
	var lastViewport [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &lastViewport[0])
	var lastScissorBox [4]int32
	gl.GetIntegerv(gl.SCISSOR_BOX, &lastScissorBox[0])
	var lastBlendSrcRgb int32
	gl.GetIntegerv(gl.BLEND_SRC_RGB, &lastBlendSrcRgb)
	var lastBlendDstRgb int32
	gl.GetIntegerv(gl.BLEND_DST_RGB, &lastBlendDstRgb)
	var lastBlendSrcAlpha int32
	gl.GetIntegerv(gl.BLEND_SRC_ALPHA, &lastBlendSrcAlpha)
	var lastBlendDstAlpha int32
	gl.GetIntegerv(gl.BLEND_DST_ALPHA, &lastBlendDstAlpha)
	var lastBlendEquationRgb int32
	gl.GetIntegerv(gl.BLEND_EQUATION_RGB, &lastBlendEquationRgb)
	var lastBlendEquationAlpha int32
	gl.GetIntegerv(gl.BLEND_EQUATION_ALPHA, &lastBlendEquationAlpha)
	lastEnableBlend := gl.IsEnabled(gl.BLEND)
	lastEnableCullFace := gl.IsEnabled(gl.CULL_FACE)
	lastEnableDepthTest := gl.IsEnabled(gl.DEPTH_TEST)
	lastEnableScissorTest := gl.IsEnabled(gl.SCISSOR_TEST)

	// Setup render state: alpha-blending enabled, no face culling, no depth testing, scissor enabled, polygon fill
	gl.Enable(gl.BLEND)
	gl.BlendEquation(gl.FUNC_ADD)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	gl.Disable(gl.CULL_FACE)
	gl.Disable(gl.DEPTH_TEST)
	gl.Enable(gl.SCISSOR_TEST)
	gl.PolygonMode(gl.FRONT_AND_BACK, gl.FILL)

	// Setup viewport, orthographic projection matrix
	// Our visible imgui space lies from draw_data->DisplayPos (top left) to draw_data->DisplayPos+data_data->DisplaySize (bottom right).
	// DisplayMin is typically (0,0) for single viewport apps.
	gl.Viewport(0, 0, int32(fbWidth), int32(fbHeight))
	orthoProjection := [4][4]float32{
		{2.0 / displayWidth, 0.0, 0.0, 0.0},
		{0.0, 2.0 / -displayHeight, 0.0, 0.0},
		{0.0, 0.0, -1.0, 0.0},
		{-1.0, 1.0, 0.0, 1.0},
	}
	gl.UseProgram(ctx.shaderHandle)
	gl.Uniform1i(ctx.attribLocationTex, 0)
	gl.UniformMatrix4fv(ctx.attribLocationProjMtx, 1, false, &orthoProjection[0][0])
	gl.BindSampler(0, 0) // Rely on combined texture/sampler state.

	// Recreate the VAO every time
	// (This is to easily allow multiple GL contexts. VAO are not shared among GL contexts, and
	// we don't track creation/deletion of windows so we don't have an obvious key to use to cache them.)
	var vaoHandle uint32
	gl.GenVertexArrays(1, &vaoHandle)
	gl.BindVertexArray(vaoHandle)
	gl.BindBuffer(gl.ARRAY_BUFFER, ctx.vboHandle)
	gl.EnableVertexAttribArray(uint32(ctx.attribLocationPosition))
	gl.EnableVertexAttribArray(uint32(ctx.attribLocationUV))
	gl.EnableVertexAttribArray(uint32(ctx.attribLocationColor))
	vertexSize, vertexOffsetPos, vertexOffsetUv, vertexOffsetCol := imgui.VertexBufferLayout()
	gl.VertexAttribPointerWithOffset(uint32(ctx.attribLocationPosition), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetPos))
	gl.VertexAttribPointerWithOffset(uint32(ctx.attribLocationUV), 2, gl.FLOAT, false, int32(vertexSize), uintptr(vertexOffsetUv))
	gl.VertexAttribPointerWithOffset(uint32(ctx.attribLocationColor), 4, gl.UNSIGNED_BYTE, true, int32(vertexSize), uintptr(vertexOffsetCol))
	indexSize := imgui.IndexBufferLayout()
	drawType := gl.UNSIGNED_SHORT
	const bytesPerUint32 = 4
	if indexSize == bytesPerUint32 {
		drawType = gl.UNSIGNED_INT
	}

	// Draw
	for _, list := range drawData.CommandLists() {
		vertexBuffer, vertexBufferSize := list.VertexBuffer()
		gl.BindBuffer(gl.ARRAY_BUFFER, ctx.vboHandle)
		gl.BufferData(gl.ARRAY_BUFFER, vertexBufferSize, vertexBuffer, gl.STREAM_DRAW)

		indexBuffer, indexBufferSize := list.IndexBuffer()
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ctx.elementsHandle)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, indexBufferSize, indexBuffer, gl.STREAM_DRAW)

		for _, cmd := range list.Commands() {
			if cmd.HasUserCallback() {
				cmd.CallUserCallback(list)
			} else {
				gl.BindTexture(gl.TEXTURE_2D, uint32(cmd.TextureID()))
				clipRect := cmd.ClipRect()
				gl.Scissor(int32(clipRect.X), int32(fbHeight)-int32(clipRect.W), int32(clipRect.Z-clipRect.X), int32(clipRect.W-clipRect.Y))
				gl.DrawElementsBaseVertexWithOffset(gl.TRIANGLES, int32(cmd.ElementCount()), uint32(drawType),
					uintptr(cmd.IndexOffset()*indexSize), int32(cmd.VertexOffset()))
			}
		}
	}
	gl.DeleteVertexArrays(1, &vaoHandle)

	// Restore modified GL state
	gl.UseProgram(uint32(lastProgram))
	gl.BindTexture(gl.TEXTURE_2D, uint32(lastTexture))
	gl.BindSampler(0, uint32(lastSampler))
	gl.ActiveTexture(uint32(lastActiveTexture))
	gl.BindVertexArray(uint32(lastVertexArray))
	gl.BindBuffer(gl.ARRAY_BUFFER, uint32(lastArrayBuffer))
	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, uint32(lastElementArrayBuffer))
	gl.BlendEquationSeparate(uint32(lastBlendEquationRgb), uint32(lastBlendEquationAlpha))
	gl.BlendFuncSeparate(uint32(lastBlendSrcRgb), uint32(lastBlendDstRgb), uint32(lastBlendSrcAlpha), uint32(lastBlendDstAlpha))
	if lastEnableBlend {
		gl.Enable(gl.BLEND)
	} else {
		gl.Disable(gl.BLEND)
	}
	if lastEnableCullFace {
		gl.Enable(gl.CULL_FACE)
	} else {
		gl.Disable(gl.CULL_FACE)
	}
	if lastEnableDepthTest {
		gl.Enable(gl.DEPTH_TEST)
	} else {
		gl.Disable(gl.DEPTH_TEST)
	}
	if lastEnableScissorTest {
		gl.Enable(gl.SCISSOR_TEST)
	} else {
		gl.Disable(gl.SCISSOR_TEST)
	}
	gl.PolygonMode(gl.FRONT_AND_BACK, uint32(lastPolygonMode[0]))
	gl.Viewport(lastViewport[0], lastViewport[1], lastViewport[2], lastViewport[3])
	gl.Scissor(lastScissorBox[0], lastScissorBox[1], lastScissorBox[2], lastScissorBox[3])
}
