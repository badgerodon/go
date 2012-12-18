package sdl

// #cgo CFLAGS: -Iinclude
// #cgo windows LDFLAGS: -Llib/windows/amd64 -lSDL
// #include <SDL/SDL.h>
import "C"

import (
	"unsafe"
)

type (
	Transparency int32
	OverlayFormat int32
	SetPaletteFlags int32

	Rect struct {
		X, Y int16
		W, H uint16
	}
	Color struct {
		R, G, B, _ uint8
	}
	PixelFormat struct {
		Ptr unsafe.Pointer
	}
	VideoInfo struct {
		Ptr unsafe.Pointer
	}
	Overlay struct {
		Ptr unsafe.Pointer
	}
)

// Transparency
const (
	Opaque Transparency = C.SDL_ALPHA_OPAQUE
	Transparent Transparency = C.SDL_ALPHA_TRANSPARENT
)

// Overlay Formats
const (
	YV12 OverlayFormat = C.SDL_YV12_OVERLAY
	IYUV OverlayFormat = C.SDL_IYUV_OVERLAY
	YUY2 OverlayFormat = C.SDL_YUY2_OVERLAY
	UYVY OverlayFormat = C.SDL_UYVY_OVERLAY
	YVYU OverlayFormat = C.SDL_YVYU_OVERLAY
)

// Palette Flags
const (
	LogicalPalette SetPaletteFlags = C.SDL_LOGPAL
	PhysicalPalette SetPaletteFlags = C.SDL_PHYSPAL
)

func RGB(r, g, b uint8) Color {
	return Color{r,g,b,0}
}

func VideoInit(driver string, flags InitFlags) error {
	c_str := C.CString(driver)
	defer C.free(unsafe.Pointer(c_str))
	c_flags := C.Uint32(flags)
	ret := C.SDL_VideoInit(c_str, c_flags)
	if ret == 0 {
		return nil
	}
	return GetError()
}

func VideoQuit() {
	C.SDL_VideoQuit()
}

func VideoDriverName() string {
	ptr := C.malloc(512)
	defer C.free(ptr)
	nptr := C.SDL_VideoDriverName((*C.char)(ptr), 512)
	if nptr == nil {
		return ""
	}
	return C.GoString(nptr)
}

func GetVideoSurface() Surface {
	ptr := C.SDL_GetVideoSurface()
	return Surface{unsafe.Pointer(ptr)}
}

func GetVideoInfo() VideoInfo {
	ptr := C.SDL_GetVideoInfo()
	return VideoInfo{unsafe.Pointer(ptr)}
}

func VideoModeOK(width, height, bitsPerPixel int, flags SurfaceFlags) int {
	c_width := C.int(width)
	c_height := C.int(height)
	c_bpp := C.int(bitsPerPixel)
	c_flags := C.Uint32(flags)
	ret := C.SDL_VideoModeOK(c_width, c_height, c_bpp, c_flags)
	return int(ret)
}

func ListModes(pixelFormat *PixelFormat, flags SurfaceFlags) []Rect {
	var c_format *C.SDL_PixelFormat
	if pixelFormat != nil {
		c_format = (*C.SDL_PixelFormat)(unsafe.Pointer(pixelFormat))
	}
	c_flags := C.Uint32(flags)
	ret := C.SDL_ListModes(c_format, c_flags)
	modes := []Rect{}
	walkArray(unsafe.Pointer(ret), func(v unsafe.Pointer) {
		modes = append(modes, *(*Rect)(v))
	})
	return modes
}

func SetVideoMode(width, height, bitsPerPixel int, flags SurfaceFlags) (Surface, error) {
	c_width := C.int(width)
	c_height := C.int(height)
	c_bpp := C.int(bitsPerPixel)
	c_flags := C.Uint32(flags)
	ret := C.SDL_SetVideoMode(c_width, c_height, c_bpp, c_flags)
	if ret == nil {
		return Surface{}, GetError()
	}
	return Surface{unsafe.Pointer(ret)}, nil
}

func SetGamma(red, green, blue float32) error {
	c_red := C.float(red)
	c_green := C.float(green)
	c_blue := C.float(blue)
	ret := C.SDL_SetGamma(c_red, c_green, c_blue)
	if ret == 0 {
		return nil
	}
	return GetError()
}

func SetGammaRamp(red, green, blue *[256]uint16) error {
	var c_red *C.Uint16
	if red != nil {
		c_red = (*C.Uint16)(unsafe.Pointer(&red[0]))
	}
	var c_green *C.Uint16
	if green != nil {
		c_green = (*C.Uint16)(unsafe.Pointer(&green[0]))
	}
	var c_blue *C.Uint16
	if blue != nil {
		c_blue = (*C.Uint16)(unsafe.Pointer(&blue[0]))
	}
	ret := C.SDL_SetGammaRamp(c_red, c_green, c_blue)
	if ret == 0 {
		return nil
	}
	return GetError()
}

func GetGammaRamp(red, green, blue *[256]uint16) error {
	var c_red *C.Uint16
	if red != nil {
		c_red = (*C.Uint16)(unsafe.Pointer(&red[0]))
	}
	var c_green *C.Uint16
	if green != nil {
		c_green = (*C.Uint16)(unsafe.Pointer(&green[0]))
	}
	var c_blue *C.Uint16
	if blue != nil {
		c_blue = (*C.Uint16)(unsafe.Pointer(&blue[0]))
	}
	ret := C.SDL_GetGammaRamp(c_red, c_green, c_blue)
	if ret == 0 {
		return nil
	}
	return GetError()
}

func MapRGB(pixelFormat PixelFormat, r, g, b uint8) uint32 {
	c_format := (*C.SDL_PixelFormat)(pixelFormat.Ptr)
	c_r := C.Uint8(r)
	c_g := C.Uint8(g)
	c_b := C.Uint8(b)
	ret := C.SDL_MapRGB(c_format, c_r, c_g, c_b)
	return uint32(ret)
}

func MapRGBA(pixelFormat PixelFormat, r, g, b, a uint8) uint32 {
	c_format := (*C.SDL_PixelFormat)(pixelFormat.Ptr)
	c_r := C.Uint8(r)
	c_g := C.Uint8(g)
	c_b := C.Uint8(b)
	c_a := C.Uint8(a)
	ret := C.SDL_MapRGBA(c_format, c_r, c_g, c_b, c_a)
	return uint32(ret)
}
