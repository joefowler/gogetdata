package getdata

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lgetdata
#include <getdata.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

// PROTECTNONE protects neither data nor format for a fragment
const PROTECTNONE Flags = C.GD_PROTECT_NONE

// PROTECTFORMAT protects format but not data for a fragment
const PROTECTFORMAT Flags = C.GD_PROTECT_FORMAT

// PROTECTDATA protects data but not format for a fragment
const PROTECTDATA Flags = C.GD_PROTECT_DATA

// PROTECTALL protects both data and format for a fragment
const PROTECTALL Flags = C.GD_PROTECT_ALL

// Fragment is used to access and modify dirfile metadata with fragment scope
// (ie., byte sex, encoding scheme, frame offset and protection levels).
type Fragment struct {
	df         *Dirfile
	index      int
	encoding   Flags
	endianness Flags
	frameoff   uint
	protection Flags
	name       string
	namespace  string
	parent     int
	prefix     string
	suffix     string
}

// NewFragment creates a pointer to fragment number index in the given Dirfile.
func NewFragment(df *Dirfile, index int) (*Fragment, error) {
	frag := &Fragment{df: df, index: index}
	cidx := C.int(index)
	frag.encoding = Flags(C.gd_encoding(df.d, cidx))
	frag.endianness = Flags(C.gd_endianness(df.d, cidx))
	frag.frameoff = uint(C.gd_frameoffset(df.d, cidx))
	frag.protection = Flags(C.gd_protection(df.d, cidx))
	frag.name = C.GoString(C.gd_fragmentname(df.d, cidx))
	frag.namespace = C.GoString(C.gd_fragment_namespace(df.d, cidx, (*C.char)(C.NULL)))
	frag.parent = -1
	if index > 0 {
		frag.parent = int(C.gd_parent_fragment(df.d, cidx))
		if frag.parent < 0 {
			return nil, df.Error()
		}
	}

	// Find the gd_verbose_prefix and suffix
	var pfx, sfx *C.char
	result := C.gd_fragment_affixes(df.d, cidx, &pfx, &sfx)
	if result != 0 {
		return nil, df.Error()
	}
	frag.prefix = C.GoString(pfx)
	frag.suffix = C.GoString(sfx)
	C.free(unsafe.Pointer(pfx))
	C.free(unsafe.Pointer(sfx))

	return frag, nil
}

// Rewrite forces GetData to rewrite a format specification fragment, even if unchanged.
func (frag *Fragment) Rewrite() error {
	result := C.gd_rewrite_fragment(frag.df.d, C.int(frag.index))
	if result != C.GD_E_OK {
		return frag.df.Error()
	}
	return nil
}

// SetNamespace sets the namespace for this fragment.
func (frag *Fragment) SetNamespace(ns string) error {
	namespace := C.CString(ns)
	defer C.free(unsafe.Pointer(namespace))
	cidx := C.int(frag.index)
	result := C.gd_fragment_namespace(frag.df.d, cidx, namespace)
	if result == nullCString {
		return frag.df.Error()
	}
	frag.namespace = ns
	return nil
}

// Prefix returns the fragment field name prefix
func (frag Fragment) Prefix() string {
	return frag.prefix
}

// SetPrefix changes the fragment's prefix to the given value
func (frag *Fragment) SetPrefix(prefix string) error {
	cprefix := C.CString(prefix)
	defer C.free(unsafe.Pointer(cprefix))
	cidx := C.int(frag.index)
	result := C.gd_alter_affixes(frag.df.d, cidx, cprefix, nullCString)
	if result < 0 {
		return frag.df.Error()
	}
	frag.prefix = prefix
	return nil
}

// Suffix returns the fragment field name suffix
func (frag Fragment) Suffix() string {
	return frag.suffix
}

// SetSuffix changes the fragment's suffix to the given value
func (frag *Fragment) SetSuffix(suffix string) error {
	csuffix := C.CString(suffix)
	defer C.free(unsafe.Pointer(csuffix))
	cidx := C.int(frag.index)
	result := C.gd_alter_affixes(frag.df.d, cidx, nullCString, csuffix)
	if result < 0 {
		return frag.df.Error()
	}
	frag.suffix = suffix
	return nil
}

// Encoding returns the encoding system of all RAW entries in the fragment.
func (frag Fragment) Encoding() Flags {
	return frag.encoding
}

// SetEncoding changes the encoding system of all RAW entries in the fragment to the
// given encoding scheme. If recode is true, then associated binary files will be re-encoded.
func (frag *Fragment) SetEncoding(encoding Flags, recode bool) error {
	cidx := C.int(frag.index)
	var rc C.int
	if recode {
		rc = 1
	}
	result := C.gd_alter_encoding(frag.df.d, C.ulong(encoding), cidx, rc)
	if result < 0 {
		return frag.df.Error()
	}
	frag.encoding = encoding
	return nil
}

// Endianness returns the endianness of all RAW entries in the fragment.
func (frag Fragment) Endianness() Flags {
	return frag.endianness
}

// SetEndianness changes the byte sex of all RAW entries in the fragment to the
// given scheme. The bytesex should be one of BIGENDIAN, LITTLEENDIAN, NATIVEENDIAN,
// or NONNATIVEENDIAN. If recode is true, then associated binary files will be re-encoded.
func (frag *Fragment) SetEndianness(bytesex Flags, recode bool) error {
	cidx := C.int(frag.index)
	var rc C.int
	if recode {
		rc = 1
	}
	result := C.gd_alter_endianness(frag.df.d, C.ulong(bytesex), cidx, rc)
	if result < 0 {
		return frag.df.Error()
	}
	frag.endianness = bytesex
	return nil
}

// FrameOffset returns the fragment frame offset
func (frag Fragment) FrameOffset() uint {
	return frag.frameoff
}

// SetFrameOffset changes the frame offset of RAW fields in a given fragment to the
// given offset. If recode is true, then associated binary files will be re-encoded.
func (frag *Fragment) SetFrameOffset(offset uint, recode bool) error {
	cidx := C.int(frag.index)
	var rc C.int
	if recode {
		rc = 1
	}
	result := C.gd_alter_frameoffset(frag.df.d, C.off_t(offset), cidx, rc)
	if result < 0 {
		return frag.df.Error()
	}
	frag.frameoff = offset
	return nil
}

// SetProtection sets the protection level for this fragment.
func (frag *Fragment) SetProtection(level Flags) error {
	result := C.gd_alter_protection(frag.df.d, C.int(level), C.int(frag.index))
	if result != C.GD_E_OK {
		return frag.df.Error()
	}
	frag.protection = level
	return nil
}
