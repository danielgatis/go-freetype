package freetype

// #include <stdlib.h>
// #include <freetype2/ft2build.h>
// #include FT_FREETYPE_H
// #include FT_ERRORS_H
import "C"

import (
	"fmt"
	"unsafe"
)

func newErr(errno C.int) error {
	var err error

	s := C.FT_Error_String(errno)

	if s != nil {
		err = fmt.Errorf(C.GoString(s))
		C.free(unsafe.Pointer(s))
	} else {
		err = fmt.Errorf("Errno 0x%02x", int(errno))
	}

	return err
}
