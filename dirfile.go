package getdata

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lgetdata
#include <getdata.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

type Dirfile struct {
	name string
	d    *C.DIRFILE
}

type Flags uint

const RDONLY Flags = 0
const (
	RDWR Flags = 1 << iota
	FORCE_ENDIAN
	BIG_ENDIAN
	CREAT
	EXCL
	TRUNC
	PEDANTIC
	FORCE_ENCODING
	VERBOSE
	IGNORE_DUPS
	INGORE_REFS
	PRETTY_PRINT
	NOR_ARM_ENDIAN
	PREMISSINVE
	TRUNCSUB
)

func OpenDirfile(name string, flags Flags) (Dirfile, error) {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cflags := C.ulong(flags)
	internal := C.gd_open(cname, cflags)
	// if internal
	d := Dirfile{name: name, d: internal}
	return d, nil
}
