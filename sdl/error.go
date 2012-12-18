package sdl

// #cgo CFLAGS: -Iinclude
// #cgo windows LDFLAGS: -Llib/windows/amd64 -lSDL
// #include <SDL/SDL.h>
import "C"

import (
	"errors"
)

func GetError() error {
	ret := C.SDL_GetError()
	if ret == nil {
		return nil
	}
	return errors.New(C.GoString(ret))
}
