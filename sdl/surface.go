package sdl

// #cgo CFLAGS: -Iinclude
// #cgo windows LDFLAGS: -Llib/windows/amd64 -lSDL
// #include <SDL/SDL.h>
import "C"

import (
	"unsafe"
)

type (
	SurfaceFlags uint32
	Surface struct {
		Ptr unsafe.Pointer
	}
)

// Surface Flags
const (
	Software SurfaceFlags = C.SDL_SWSURFACE
	Hardware SurfaceFlags = C.SDL_HWSURFACE
	AsyncBlit SurfaceFlags = C.SDL_ASYNCBLIT
	AnyFormat SurfaceFlags = C.SDL_ANYFORMAT
	HardwarePalette SurfaceFlags = C.SDL_HWPALETTE
	DoubleBuffer SurfaceFlags = C.SDL_DOUBLEBUF
	Fullscreen SurfaceFlags = C.SDL_FULLSCREEN
	OpenGL SurfaceFlags = C.SDL_OPENGL
	OpenGLBlit SurfaceFlags = C.SDL_OPENGLBLIT
	Resizable SurfaceFlags = C.SDL_RESIZABLE
	NoFrame SurfaceFlags = C.SDL_NOFRAME
)

func BlitSurface(src Surface, srcrect *Rect, dst Surface, dstrect *Rect) error {
	c_src := (*C.SDL_Surface)(src.Ptr)
	c_srcrect := (*C.SDL_Rect)(unsafe.Pointer(srcrect))
	c_dst := (*C.SDL_Surface)(dst.Ptr)
	c_dstrect := (*C.SDL_Rect)(unsafe.Pointer(dstrect))
	ret := C.SDL_BlitSurface(c_src, c_srcrect, c_dst, c_dstrect)
	if ret == 0 {
		return nil
	}
	return GetError()
}

func (this Surface) Flip() error {
	c_surface := (*C.SDL_Surface)(this.Ptr)
	ret := C.SDL_Flip(c_surface)
	if ret == 0 {
		return nil
	}
	return GetError()
}

func (this Surface) Free() {
	c_surface := (*C.SDL_Surface)(this.Ptr)
	C.SDL_FreeSurface(c_surface)
}

func (this Surface) SetColors(colors []Color) error {
	c_surface := (*C.SDL_Surface)(this.Ptr)
	var c_colors *C.SDL_Color
	if len(colors) > 0 {
		c_colors = (*C.SDL_Color)(unsafe.Pointer(&colors[0]))
	}
	c_firstcolor := C.int(0)
	c_ncolors := C.int(len(colors))
	C.SDL_SetColors(c_surface, c_colors, c_firstcolor, c_ncolors)
	return GetError()
}

func (this Surface) SetPalette(flags SetPaletteFlags, colors []Color) error {
	c_surface := (*C.SDL_Surface)(this.Ptr)
	c_flags := C.int(flags)
	var c_colors *C.SDL_Color
	if len(colors) > 0 {
		c_colors = (*C.SDL_Color)(unsafe.Pointer(&colors[0]))
	}
	c_firstcolor := C.int(0)
	c_ncolors := C.int(len(colors))
	C.SDL_SetPalette(c_surface, c_flags, c_colors, c_firstcolor, c_ncolors)
	return GetError()
}

func (this Surface) UpdateRect(rect Rect) {
	c_surface := (*C.SDL_Surface)(this.Ptr)
	c_x := C.Sint32(rect.X)
	c_y := C.Sint32(rect.Y)
	c_w := C.Uint32(rect.W)
	c_h := C.Uint32(rect.H)
	C.SDL_UpdateRect(c_surface, c_x, c_y, c_w, c_h)
}

func (this Surface) UpdateRects(rects []Rect) {
	c_surface := (*C.SDL_Surface)(this.Ptr)
	c_numrects := C.int(len(rects))
	var c_rects *C.SDL_Rect
	if len(rects) > 0 {
		c_rects = (*C.SDL_Rect)(unsafe.Pointer(&rects[0]))
	}
	C.SDL_UpdateRects(c_surface, c_numrects, c_rects)
}
