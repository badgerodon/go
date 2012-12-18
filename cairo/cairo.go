// Go binding for the cairo graphics library
package cairo

// #cgo CFLAGS: -Iinclude
// #cgo darwin LDFLAGS: -Llib/darwin/amd64 -lcairo
// #cgo windows LDFLAGS: -Llib/windows/amd64 -lcairo
// #include <cairo/cairo.h>
import "C"


type FontType int
// cairo_font_type_t
const (
	FONT_TYPE_TOY FontType = iota
	FONT_TYPE_FT
	FONT_TYPE_WIN32
	FONT_TYPE_QUARTZ
	FONT_TYPE_USER
)

type RegionOverlap int

const (
	REGION_OVERLAP_IN RegionOverlap = iota
	REGION_OVERLAP_OUT
	REGION_OVERLAP_PART
)

type DeviceType int

const (
	DEVICE_TYPE_DRM DeviceType = iota
	DEVICE_TYPE_GL
	DEVICE_TYPE_SCRIPT
	DEVICE_TYPE_XCB
	DEVICE_TYPE_XLIB
	DEVICE_TYPE_XML
)

const (
	MIME_TYPE_JPEG = "image/jpeg"
	MIME_TYPE_PNG  = "image/png"
	MIME_TYPE_JP2  = "image/jp2"
	MIME_TYPE_URI  = "text/x-uri"
)







type FontFace struct {
	// todo
}

type FontOptions struct {
	// todo
}

type ScaledFont struct {
	// todo
}

type Device struct {
	// todo
}

func Version() int {
	return int(C.cairo_version())
}

func VersionString() string {
	return C.GoString(C.cairo_version_string())
}

func cairobool2bool(flag C.cairo_bool_t) bool {
	if int(flag) > 0 {
		return true
	}
	return false;
}
