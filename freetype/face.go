package freetype

// #include <stdlib.h>
// #include <freetype2/ft2build.h>
// #include FT_FREETYPE_H
import "C"

import (
	"fmt"
	"image"
	"reflect"
	"sync"
	"unsafe"
)

// Metrics Represents the glyph metrics
type Metrics struct {
	Width              int
	Height             int
	HorizontalBearingX int
	HorizontalBearingY int
	AdvanceWidth       int
	VerticalBearingX   int
	VerticalBearingY   int
	AdvanceHeight      int
}

// Face Represents a font face
type Face struct {
	sync.Mutex
	library *Library
	native  C.FT_Face
}

// NewFace Creates a new font face
func NewFace(l *Library, fontBytes []byte, idx int) (*Face, error) {
	f := &Face{
		library: l,
	}

	n := len(fontBytes)
	b := C.CBytes(fontBytes)

	if errno := C.FT_New_Memory_Face(l.native, (*C.uchar)(b), C.long(n), C.long(idx), &f.native); errno != 0 {
		return nil, newErr(errno)
	}

	return f, nil
}

// Pt Changes the size (in points)
func (f *Face) Pt(pt, dpi int) error {
	f.Lock()
	defer f.Unlock()

	if errno := C.FT_Set_Char_Size(
		f.native,
		0,
		C.FT_F26Dot6(pt<<6),
		0,
		C.FT_UInt(dpi),
	); errno != 0 {
		return newErr(errno)
	}

	return nil
}

// Glyph Gets the glyph as an image
func (f *Face) Glyph(ch rune) (*image.RGBA, *Metrics, error) {
	f.Lock()
	defer f.Unlock()

	loadFlags := C.FT_LOAD_RENDER

	if f.hasColor() {
		loadFlags = loadFlags | C.FT_LOAD_COLOR
	}

	errno := C.FT_Load_Char(
		f.native,
		C.ulong(ch),
		C.int(loadFlags),
	)

	if errno != 0 {
		return nil, nil, newErr(errno)
	}

	pix, err := rasterize(&f.native.glyph.bitmap)

	if err != nil {
		return nil, nil, err
	}

	m := f.native.glyph.metrics
	metrics := &Metrics{}

	metrics.Width = int(m.width >> 6)
	metrics.Height = int(m.height >> 6)
	metrics.HorizontalBearingX = int(m.horiBearingX >> 6)
	metrics.HorizontalBearingY = int(m.horiBearingY >> 6)
	metrics.AdvanceWidth = int(m.horiAdvance >> 6)
	metrics.VerticalBearingX = int(m.vertBearingX >> 6)
	metrics.VerticalBearingY = int(m.vertBearingY >> 6)
	metrics.AdvanceHeight = int(m.vertAdvance >> 6)

	img := &image.RGBA{
		Pix:    pix,
		Stride: int(f.native.glyph.bitmap.width) * 4,
		Rect: image.Rect(
			0,
			0,
			int(f.native.glyph.bitmap.width),
			int(f.native.glyph.bitmap.rows),
		),
	}

	return img, metrics, nil
}

// Done Done
func (f *Face) Done() error {
	f.Lock()
	defer f.Unlock()

	if errno := C.FT_Done_Face(f.native); errno != 0 {
		return newErr(errno)
	}

	f.native = nil
	return nil
}

func rasterize(bitmap *C.FT_Bitmap) ([]byte, error) {
	size := int(bitmap.rows) * int(bitmap.width) * int(bitmap.pitch)

	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(bitmap.buffer)),
		Len:  size,
		Cap:  size,
	}

	data := *(*[]byte)(unsafe.Pointer(&hdr))
	pix := make([]byte, 0)

	switch bitmap.pixel_mode {
	case C.FT_PIXEL_MODE_MONO:
		unpackByte := func(b byte, count int) {
			bit := 7

			for count != 0 {
				v := ((b >> bit) & 1) * 255

				pix = append(pix, v)   //R
				pix = append(pix, v)   //G
				pix = append(pix, v)   //B
				pix = append(pix, 255) //A

				count -= 1
				bit -= 1
			}
		}

		for i := 0; i < int(bitmap.rows); i++ {
			columns := int(bitmap.width)
			b := 0
			offset := i * int(abs(int32(bitmap.pitch)))

			for columns != 0 {
				bits := min(8, columns)
				unpackByte(data[offset+b], bits)
				columns -= bits
				b += 1
			}
		}

		return pix, nil

	case C.FT_PIXEL_MODE_GRAY:
		for i := 0; i < int(bitmap.rows); i++ {
			start := i * int(bitmap.pitch)
			stop := start + int(bitmap.width)

			for _, b := range data[start:stop] {
				pix = append(pix, b)   //R
				pix = append(pix, b)   //G
				pix = append(pix, b)   //B
				pix = append(pix, 255) //A
			}
		}

		return pix, nil

	case C.FT_PIXEL_MODE_BGRA:
		p := int(bitmap.pitch) / int(bitmap.width)
		s := int(bitmap.rows) * int(bitmap.width) * p

		for i := 0; i < s; i = i + p {
			pix = append(pix, data[i+2]) //R
			pix = append(pix, data[i+1]) //G
			pix = append(pix, data[i+0]) //B
			pix = append(pix, data[i+3]) //A
		}

		return pix, nil

	default:
		return pix, fmt.Errorf("pixel_model 0x%02x not implemented", bitmap.pixel_mode)
	}
}

func (f *Face) hasColor() bool {
	return f.native.face_flags&C.FT_FACE_FLAG_COLOR != 0
}

func abs(n int32) int32 {
	y := n >> 31
	return (n ^ y) - y
}

func min(a, b int) int {
	if a < b {
		return a
	}

	return b
}
