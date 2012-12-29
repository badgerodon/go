package cairo

// #cgo CFLAGS: -Iinclude
// #cgo darwin LDFLAGS: -Llib/darwin/amd64 -lcairo
// #cgo windows LDFLAGS: -Llib/windows/amd64 -lcairo
// #include <cairo/cairo.h>
// #include <cairo/cairo-win32.h>
import "C"

import (
  "unsafe"
)