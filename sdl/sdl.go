package sdl

// #cgo CFLAGS: -Iinclude
// #cgo windows LDFLAGS: -Llib/windows/amd64 -lSDL
// #include <SDL/SDL.h>
import "C"

import (
	"unsafe"
	//"os"
	//"runtime"
//	"sync"
	//"unsafe"
)

type (
	InitFlags uint32
)

const (
	InitTimer InitFlags = C.SDL_INIT_TIMER
	InitAudio InitFlags = C.SDL_INIT_AUDIO
	InitVideo InitFlags = C.SDL_INIT_VIDEO
	InitCDRom InitFlags = C.SDL_INIT_CDROM
	InitJoystick InitFlags = C.SDL_INIT_JOYSTICK
	InitEverything InitFlags = C.SDL_INIT_EVERYTHING
	InitNoParachute InitFlags = C.SDL_INIT_NOPARACHUTE
	InitEventThread InitFlags = C.SDL_INIT_EVENTTHREAD
)

// Walk all the values in a c array (ends with null)
func walkArray(ptr unsafe.Pointer, handler func(unsafe.Pointer)) {
	i := uintptr(ptr)
	for {
		v := unsafe.Pointer(i)
		if v == nil {
			break
		}
		handler(v)
		i++
	}
}

// The SDL_Init function initializes the Simple Directmedia Library and the subsystems specified by flags. It should be called before all other SDL functions.
func Init(flags InitFlags) error {
	c_flags := C.Uint32(flags)
	c_err := C.SDL_Init(c_flags)
	if c_err == 0 {
		return nil
	}
	return GetError()
}

func Quit() {
	C.SDL_Quit()
}
