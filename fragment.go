package getdata

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lgetdata
#include <getdata.h>
#include <stdlib.h>
*/
import "C"

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
	frag.parent = -1
	if index > 0 {
		frag.parent = int(C.gd_parent_fragment(df.d, cidx))
		if frag.parent < 0 {
			return nil, df.Error()
		}
	}

	return frag, nil
}

// SetProtection sets the protection level for this fragment.
func (frag *Fragment) SetProtection(level Flags) error {
	result := C.gd_alter_protection(frag.df.d, C.int(level), C.int(frag.index))
	if result == C.GD_E_OK {
		frag.protection = level
		return nil
	}
	return frag.df.Error()
}
