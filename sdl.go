package sdl

// #cgo windows CFLAGS: -Iinclude/windows/x64
// #cgo windows LDFLAGS: -Llib/windows/x64 -lSDL -lSDL_image
// #include <SDL/SDL.h>
// #include <SDL/SDL_image.h>
import "C"

import (
	"errors"
	"unsafe"
	//"os"
	//"runtime"
//	"sync"
	//"unsafe"
)

type (
	InitFlags uint32
	Rectangle struct {
		ptr *C.SDL_Rect
	}
	SurfaceFlags uint32
	Surface struct {
		ptr *C.SDL_Surface
	}
)

const (
	INIT_TIMER InitFlags = C.SDL_INIT_TIMER
	INIT_AUDIO InitFlags = C.SDL_INIT_AUDIO
	INIT_VIDEO InitFlags = C.SDL_INIT_VIDEO
	INIT_CDROM InitFlags = C.SDL_INIT_CDROM
	INIT_JOYSTICK InitFlags = C.SDL_INIT_JOYSTICK
	INIT_EVERYTHING InitFlags = C.SDL_INIT_EVERYTHING
	INIT_NOPARACHUTE InitFlags = C.SDL_INIT_NOPARACHUTE
	INIT_EVENTTHREAD InitFlags = C.SDL_INIT_EVENTTHREAD

	// Create the video surface in system memory
	SWSURFACE SurfaceFlags = C.SDL_SWSURFACE
	// Create the video surface in video memory
	HWSURFACE SurfaceFlags = C.SDL_HWSURFACE
	// Enables the use of asynchronous updates of the display surface. This will usually slow down blitting on single CPU machines, but may provide a speed increase on SMP systems.
	ASYNCBLIT SurfaceFlags = C.SDL_ASYNCBLIT
	// Normally, if a video surface of the requested bits-per-pixel (bpp) is not available, SDL will emulate one with a shadow surface. Passing SDL_ANYFORMAT prevents this and causes SDL to use the video surface, regardless of its pixel depth.
	ANYFORMAT SurfaceFlags = C.SDL_ANYFORMAT
	// Give SDL exclusive palette access. Without this flag you may not always get the colors you request with SDL_SetColors or SDL_SetPalette.
	HWPALETTE SurfaceFlags = C.SDL_HWPALETTE
	// Enable hardware double buffering; only valid with SDL_HWSURFACE. Calling SDL_Flip will flip the buffers and update the screen. All drawing will take place on the surface that is not displayed at the moment. If double buffering could not be enabled then SDL_Flip will just perform a SDL_UpdateRect on the entire screen.SDL_HWPALETTE
	DOUBLEBUF SurfaceFlags = C.SDL_DOUBLEBUF
	// SDL will attempt to use a fullscreen mode. If a hardware resolution change is not possible (for whatever reason), the next higher resolution will be used and the display window centered on a black background.
	FULLSCREEN SurfaceFlags = C.SDL_FULLSCREEN
	// Create an OpenGL rendering context. You should have previously set OpenGL video attributes with SDL_GL_SetAttribute.
	OPENGL SurfaceFlags = C.SDL_OPENGL
	// Create an OpenGL rendering context, like above, but allow normal blitting operations. The screen (2D) surface may have an alpha channel, and SDL_UpdateRects must be used for updating changes to the screen surface. NOTE: This option is kept for compatibility only, and will be removed in next versions. Is not recommended for new code.
	OPENGLBLIT SurfaceFlags = C.SDL_OPENGLBLIT
	// Create a resizable window. When the window is resized by the user a SDL_VIDEORESIZE event is generated and SDL_SetVideoMode can be called again with the new size.
	RESIZABLE SurfaceFlags = C.SDL_RESIZABLE
	// If possible, SDL_NOFRAME causes SDL to create a window with no title bar or frame decoration. Fullscreen modes automatically have this flag set.
	NOFRAME SurfaceFlags = C.SDL_NOFRAME
)

func (this *Surface) toC() *C.SDL_Surface {
	if this == nil {
		return nil
	}
	return this.ptr
}
func (this *Rectangle) toC() *C.SDL_Rect {
	if this == nil {
		return nil
	}
	return this.ptr
}

func lastError() error {
	str := C.SDL_GetError()
	return errors.New(C.GoString(str))
}

// The SDL_Init function initializes the Simple Directmedia Library and the subsystems specified by flags. It should be called before all other SDL functions.
func Init(flags InitFlags) error {
	c_flags := C.Uint32(flags)
	c_err := C.SDL_Init(c_flags)
	if c_err == 0 {
		return nil
	}
	return lastError()
}

func Load(filename string) (*Surface, error) {
	c_string := C.CString(filename)
	defer C.free(unsafe.Pointer(c_string))
	c_surface := C.IMG_Load(c_string)
	if c_surface == nil {
		return nil, lastError()
	}
	return &Surface{c_surface}, nil
}

// This performs a fast blit from the source surface to the destination surface.
func BlitSurface(src *Surface, srcrect *Rectangle, dst *Surface, dstrect *Rectangle) error {
	c_err := C.SDL_BlitSurface(src.toC(), srcrect.toC(), dst.toC(), dstrect.toC())
	if c_err == 0 {
		return nil
	}
	return lastError()
}

func Delay(ms uint32) {
	c_ms := C.Uint32(ms)
	C.SDL_Delay(c_ms)
}

func Flip(surface *Surface) error {
	c_err := C.SDL_Flip(surface.ptr)
	if c_err == 0 {
		return nil
	}
	return lastError()
}


func SetVideoMode(width, height, bitsPerPixel int, flags SurfaceFlags) (*Surface, error) {
	c_width := C.int(width)
	c_height := C.int(height)
	c_bitsPerPixel := C.int(bitsPerPixel)
	c_flags := C.Uint32(flags)
	c_surface := C.SDL_SetVideoMode(c_width, c_height, c_bitsPerPixel, c_flags)
	if c_surface == nil {
		return nil, lastError()
	}
	return &Surface{c_surface}, nil
}

func Quit() {
	C.SDL_Quit()
}
