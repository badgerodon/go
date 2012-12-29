package ttf

// #cgo CFLAGS: -Iinclude/
// #cgo darwin LDFLAGS: -lSDL_ttf
// #cgo windows LDFLAGS: -Llib/windows/amd64 -lSDL_ttf
// #include <SDL/SDL_ttf.h>
import "C"

import (
	"unsafe"

	"github.com/badgerodon/go/sdl"
)

type (
	Font struct {
		Ptr unsafe.Pointer
	}
	FontStyle int32
	Hinting int32
)

const (
	Normal FontStyle = C.TTF_STYLE_NORMAL
	Bold FontStyle = C.TTF_STYLE_BOLD
	Italic FontStyle = C.TTF_STYLE_ITALIC
	Underline FontStyle = C.TTF_STYLE_UNDERLINE
	StrikeThrough FontStyle = C.TTF_STYLE_STRIKETHROUGH

	NormalHinting Hinting = C.TTF_HINTING_NORMAL
	LightHinting Hinting = C.TTF_HINTING_LIGHT
	MonoHinting Hinting = C.TTF_HINTING_MONO
	NoHinting Hinting = C.TTF_HINTING_NONE
)

func ByteSwappedUnicode(isSwapped bool) {
	c_swapped := C.int(0)
	if isSwapped {
		c_swapped = C.int(1)
	}
	C.TTF_ByteSwappedUNICODE(c_swapped)
}

func Init() error {
	ret := C.TTF_Init()
	if ret == 0 {
		return nil
	}
	return sdl.GetError()
}
func WasInit() bool {
	ret := C.TTF_WasInit()
	if ret == 1 {
		return true
	}
	return false
}
func Quit() {
	C.TTF_Quit()
}

func OpenFont(file string, ptsize int) (Font, error) {
	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))
	c_ptsize := C.int(ptsize)
	ret := C.TTF_OpenFont(c_file, c_ptsize)
	if ret == nil {
		return Font{}, sdl.GetError()
	}
	return Font{unsafe.Pointer(ret)}, nil
}
func OpenFontIndex(file string, ptsize int, index int64) (Font, error) {
	c_file := C.CString(file)
	defer C.free(unsafe.Pointer(c_file))
	c_ptsize := C.int(ptsize)
	c_index := C.long(index)
	ret := C.TTF_OpenFontIndex(c_file, c_ptsize, c_index)
	if ret == nil {
		return Font{}, sdl.GetError()
	}
	return Font{unsafe.Pointer(ret)}, nil
}/*
func OpenFontRW(src sdl.RWops, freesrc int, ptsize int) Font {

}
func OpenFontIndexRW(src sdl.RWops, freesrc int, ptsize int, index int64) Font {

}
*/
func (this Font) GetStyle() FontStyle {
	c_font := (*C.TTF_Font)(this.Ptr)
	ret := C.TTF_GetFontStyle(c_font)
	return FontStyle(ret)
}
func (this Font) SetStyle(style FontStyle) {
	c_font := (*C.TTF_Font)(this.Ptr)
	c_style := C.int(style)
	C.TTF_SetFontStyle(c_font, c_style)
}
/*
func (this Font) GetOutline() int {

}
func (this Font) SetOutline(outline int) {

}
func (this Font) GetHinting() Hinting {

}
func (this Font) SetHinting(hinting Hinting) {

}
func (this Font) Height() int {

}
func (this Font) Ascent() int {

}
func (this Font) Descent() int {

}
func (this Font) LineSkip() int {

}
func (this Font) GetKerning() int {

}
func (this Font) SetKerning(allowed int) {

}
func (this Font) Faces() int64 {

}
func (this Font) FaceIsFixedWidth() bool {

}
func (this Font) FaceFamilyName() string {

}
func (this Font) FaceStyleName() string {

}
func (this Font) SizeText(text string) (width int, height int, err error) {

}
*/
func (this Font) RenderTextSolid(text string, color sdl.Color) (sdl.Surface, error) {
	c_font := (*C.TTF_Font)(this.Ptr)
	c_text := C.CString(text)
	defer C.free(unsafe.Pointer(c_text))
	c_color := *(*C.SDL_Color)(unsafe.Pointer(&color))
	ret := C.TTF_RenderUTF8_Solid(c_font, c_text, c_color)
	if ret == nil {
		return sdl.Surface{}, sdl.GetError()
	}
	return sdl.Surface{unsafe.Pointer(ret)}, nil
}
func (this Font) RenderTextShaded(text string, fg, bg sdl.Color) (sdl.Surface, error) {
	c_font := (*C.TTF_Font)(this.Ptr)
	c_text := C.CString(text)
	defer C.free(unsafe.Pointer(c_text))
	c_fg := *(*C.SDL_Color)(unsafe.Pointer(&fg))
	c_bg := *(*C.SDL_Color)(unsafe.Pointer(&bg))
	ret := C.TTF_RenderUTF8_Shaded(c_font, c_text, c_fg, c_bg)
	if ret == nil {
		return sdl.Surface{}, sdl.GetError()
	}
	return sdl.Surface{unsafe.Pointer(ret)}, nil
}
func (this Font) RenderTextBlended(text string, fg sdl.Color) (sdl.Surface, error) {
	c_font := (*C.TTF_Font)(this.Ptr)
	c_text := C.CString(text)
	defer C.free(unsafe.Pointer(c_text))
	c_fg := *(*C.SDL_Color)(unsafe.Pointer(&fg))
	ret := C.TTF_RenderUTF8_Blended(c_font, c_text, c_fg)
	if ret == nil {
		return sdl.Surface{}, sdl.GetError()
	}
	return sdl.Surface{unsafe.Pointer(ret)}, nil
}
func (this Font) Close() {
	c_font := (*C.TTF_Font)(this.Ptr)
	C.TTF_CloseFont(c_font)
}
