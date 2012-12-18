package sdl

// #cgo CFLAGS: -Iinclude
// #cgo windows LDFLAGS: -Llib/windows/amd64 -lSDL
// #include <SDL/SDL.h>
import "C"

type (
	ActiveEvent struct {
		Type, Gain, State uint8
	}
	KeySym struct {
		ScanCode uint8
		Sym      Key
		Mod      Mod
		Unicde   uint16
	}
	KeyboardEvent struct {
		Type, State uint8
		KeySym      KeySym
	}
)
