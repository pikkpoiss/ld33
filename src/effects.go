// Copyright 2015 Pikkpoiss
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"../lib/twodee"
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

const EFFECTS_FRAGMENT = `#version 150
precision mediump float;

uniform sampler2D u_TextureUnit;
uniform vec3 v_Color;
in vec2 v_TextureCoordinates;
out vec4 v_FragData;

void main() {
  v_FragData = clamp(
    vec4(clamp(v_Color, 0.0, 0.5), 1.0) +
      texture(u_TextureUnit, v_TextureCoordinates),
    0.0,
    1.0);
}` + "\x00"

const EFFECTS_VERTEX = `#version 150

in vec4 a_Position;
in vec2 a_TextureCoordinates;

out vec2 v_TextureCoordinates;

void main() {
    v_TextureCoordinates = a_TextureCoordinates;
    gl_Position = a_Position;
}` + "\x00"

type EffectsRenderer struct {
	framebuffer         uint32
	texture             uint32
	shader              uint32
	positionLoc         uint32
	textureLoc          uint32
	colorLoc            int32
	bufferDimensionsLoc int32
	textureUnitLoc      int32
	coords              uint32
	width               int
	height              int
	oldwidth            int
	oldheight           int
	Color               mgl32.Vec3
}

func NewEffectsRenderer(w, h int) (r *EffectsRenderer, err error) {
	r = &EffectsRenderer{
		width:  w,
		height: h,
		Color:  mgl32.Vec3{0.0, 0.0, 0.0},
	}
	_, _, r.oldwidth, r.oldheight = twodee.GetInteger4(gl.VIEWPORT)
	if r.shader, err = twodee.BuildProgram(EFFECTS_VERTEX, EFFECTS_FRAGMENT); err != nil {
		return
	}
	r.positionLoc = uint32(gl.GetAttribLocation(r.shader, gl.Str("a_Position\x00")))
	r.textureLoc = uint32(gl.GetAttribLocation(r.shader, gl.Str("a_TextureCoordinates\x00")))
	r.textureUnitLoc = gl.GetUniformLocation(r.shader, gl.Str("u_TextureUnit\x00"))
	r.colorLoc = gl.GetUniformLocation(r.shader, gl.Str("v_Color\x00"))
	gl.BindFragDataLocation(r.shader, 0, gl.Str("v_FragData\x00"))
	var size float32 = 1.0
	var rect = []float32{
		-size, -size, 0.0, 0, 0,
		-size, size, 0.0, 0, 1,
		size, -size, 0.0, 1, 0,
		size, size, 0.0, 1, 1,
	}
	if r.coords, err = twodee.CreateVBO(len(rect)*4, rect, gl.STATIC_DRAW); err != nil {
		return
	}

	if r.framebuffer, r.texture, err = r.initFramebuffer(w, h); err != nil {
		return
	}
	return
}

func (r *EffectsRenderer) initFramebuffer(w, h int) (fb uint32, tex uint32, err error) {
	gl.GenFramebuffers(1, &fb)
	gl.BindFramebuffer(gl.FRAMEBUFFER, fb)

	gl.GenTextures(1, &tex)
	gl.BindTexture(gl.TEXTURE_2D, tex)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, int32(w), int32(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)

	gl.FramebufferTexture2D(gl.DRAW_FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, tex, 0)
	if err = r.GetError(); err != nil {
		return
	}
	buffers := []uint32{gl.COLOR_ATTACHMENT0}
	gl.DrawBuffers(1, &buffers[0])

	var rb uint32
	gl.GenRenderbuffers(1, &rb)
	gl.BindRenderbuffer(gl.RENDERBUFFER, rb)
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.STENCIL_INDEX8, int32(w), int32(h))
	gl.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.STENCIL_ATTACHMENT, gl.RENDERBUFFER, rb)

	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	return
}

func (r *EffectsRenderer) GetError() error {
	if e := gl.GetError(); e != 0 {
		return fmt.Errorf("OpenGL error: %X", e)
	}
	var status = gl.CheckFramebufferStatus(gl.DRAW_FRAMEBUFFER)
	switch status {
	case gl.FRAMEBUFFER_COMPLETE:
		return nil
	case gl.FRAMEBUFFER_INCOMPLETE_ATTACHMENT:
		return fmt.Errorf("Attachment point unconnected")
	case gl.FRAMEBUFFER_INCOMPLETE_MISSING_ATTACHMENT:
		return fmt.Errorf("Missing attachment")
	case gl.FRAMEBUFFER_INCOMPLETE_DRAW_BUFFER:
		return fmt.Errorf("Draw buffer")
	case gl.FRAMEBUFFER_INCOMPLETE_READ_BUFFER:
		return fmt.Errorf("Read buffer")
	case gl.FRAMEBUFFER_UNSUPPORTED:
		return fmt.Errorf("Unsupported config")
	default:
		return fmt.Errorf("Unknown framebuffer error: %X", status)
	}
}

func (r *EffectsRenderer) Delete() error {
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindRenderbuffer(gl.RENDERBUFFER, 0)
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.DeleteFramebuffers(1, &r.framebuffer)
	gl.DeleteTextures(1, &r.texture)
	gl.DeleteBuffers(1, &r.coords)
	return r.GetError()
}

func (r *EffectsRenderer) Bind() error {
	gl.BindFramebuffer(gl.FRAMEBUFFER, r.framebuffer)
	gl.Enable(gl.STENCIL_TEST)
	gl.Viewport(0, 0, int32(r.width), int32(r.height))
	gl.ClearStencil(0)
	gl.ClearColor(0.0, 0.0, 0.0, 0.0)
	gl.StencilMask(0xFF) // Write to buffer
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.STENCIL_BUFFER_BIT)
	gl.StencilMask(0x00) // Don't write to buffer
	return nil
}

func (r *EffectsRenderer) Draw() (err error) {
	gl.UseProgram(r.shader)
	gl.Uniform1i(r.textureUnitLoc, 0)
	gl.Uniform3f(r.colorLoc, r.Color[0], r.Color[1], r.Color[2])
	gl.BindBuffer(gl.ARRAY_BUFFER, r.coords)
	gl.EnableVertexAttribArray(r.positionLoc)
	gl.VertexAttribPointer(r.positionLoc, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(r.textureLoc)
	gl.VertexAttribPointer(r.textureLoc, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, int32(r.oldwidth), int32(r.oldheight))
	gl.BlendFunc(gl.ONE, gl.ONE)
	gl.BindTexture(gl.TEXTURE_2D, r.texture)
	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
	return nil
}

func (r *EffectsRenderer) Unbind() error {
	gl.Viewport(0, 0, int32(r.oldwidth), int32(r.oldheight))
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Disable(gl.STENCIL_TEST)
	return nil
}

func (r *EffectsRenderer) DisableOutput() {
	gl.ColorMask(false, false, false, false)
	gl.StencilFunc(gl.NEVER, 1, 0xFF)                // Never pass
	gl.StencilOp(gl.REPLACE, gl.REPLACE, gl.REPLACE) // Replace to ref=1
	gl.StencilMask(0xFF)                             // Write to buffer
}

func (r *EffectsRenderer) EnableOutput() {
	gl.ColorMask(true, true, true, true)
	gl.StencilMask(0x00)              // No more writing
	gl.StencilFunc(gl.EQUAL, 0, 0xFF) // Only pass where stencil is 0
}
