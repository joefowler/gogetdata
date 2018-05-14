package getdata

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lgetdata
#include <getdata.h>
#include <stdlib.h>
*/
import "C"
import (
	"fmt"
	"unsafe"
)

// Dirfile wraps the GetData.DIRFILE opaque object.
type Dirfile struct {
	name string
	d    *C.DIRFILE
}

// Flags are dirfile opening flags, including encoding methods
type Flags uint

// RDONLY open read-only
const RDONLY Flags = C.GD_RDONLY

// RDWR open read/write
const RDWR Flags = C.GD_RDWR
const FORCEENDIAN Flags = C.GD_FORCE_ENDIAN
const BIGENDIAN Flags = C.GD_BIG_ENDIAN
const CREAT Flags = C.GD_CREAT
const EXCL Flags = C.GD_EXCL
const TRUNC Flags = C.GD_TRUNC
const PEDANTIC Flags = C.GD_PEDANTIC
const FORCEENCODING Flags = C.GD_FORCE_ENCODING
const VERBOSE Flags = C.GD_VERBOSE
const IGNOREDUPS Flags = C.GD_IGNORE_DUPS
const IGNOREREFS Flags = C.GD_IGNORE_REFS
const PRETTYPRINT Flags = C.GD_PRETTY_PRINT
const PERMISSIVE Flags = C.GD_PERMISSIVE
const TRUNCSUB Flags = C.GD_TRUNCSUB

// OpenDirfile returns an open Dirfile object, with read/write, encoding, and other flags
// given by the flags argument.
func OpenDirfile(name string, flags Flags) (Dirfile, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	result := C.gd_open(cname, C.ulong(flags))
	errcode := C.gd_error(result)
	dirfile := Dirfile{name: name}
	if errcode != 0 {
		return dirfile, fmt.Errorf("Error opening %s: error code %d", name, errcode)
	}
	dirfile.d = result
	return dirfile, nil
}