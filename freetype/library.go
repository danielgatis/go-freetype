package freetype

// #include <stdlib.h>
// #include <freetype2/ft2build.h>
// #include FT_FREETYPE_H
import "C"

import (
	"sync"
)

// Library Represents a library
type Library struct {
	sync.Mutex
	native C.FT_Library
}

// NewLibrary Creates a new library
func NewLibrary() (*Library, error) {
	l := &Library{}

	if errno := C.FT_Init_FreeType(&l.native); errno != 0 {
		return nil, newErr(errno)
	}

	return l, nil
}

// Done Done
func (l *Library) Done() error {
	l.Lock()
	defer l.Unlock()

	if errno := C.FT_Done_FreeType(l.native); errno != 0 {
		return newErr(errno)
	}

	l.native = nil
	return nil
}
