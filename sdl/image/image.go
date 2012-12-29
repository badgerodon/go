package image

// #cgo CFLAGS: -Iinclude
// #cgo darwin LDFLAGS: -lSDL_image
// #cgo windows LDFLAGS: -Llib/windows/amd64 -lSDL_image
// #include <SDL/SDL_image.h>
import "C"

import (
	"unsafe"

	"github.com/badgerodon/go/sdl"
)

type (
	InitFlags int32
)

const (
	InitJPG InitFlags = C.IMG_INIT_JPG
	InitPNG InitFlags = C.IMG_INIT_PNG
	InitTIF InitFlags = C.IMG_INIT_TIF
	InitWEBP InitFlags = C.IMG_INIT_WEBP
)

func Init(flags InitFlags) error {
	c_flags := C.int(flags)
	ret := C.IMG_Init(c_flags)
	if ret == 0 {
		return sdl.GetError()
	}
	return nil
}

func Quit() {
	C.IMG_Quit()
}

//func LoadTypedRW(src sdl.RWops, freesrc int, type string) (sdl.Surface, error) {

//}

func Load(file string) (sdl.Surface, error) {
	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))
	ret := C.IMG_Load(c_file)
	if ret == nil {
		return sdl.Surface{}, sdl.GetError()
	}
	return sdl.Surface{unsafe.Pointer(ret)}, nil
}

//func LoadRW(src sdl.RWops, freesrc int) (sdl.Surface, error) {

//}
