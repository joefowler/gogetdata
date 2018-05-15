package getdata

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lgetdata
#include <getdata.h>
#include <stdlib.h>
*/
import "C"
import (
	"errors"
	"unsafe"
)

// Dirfile wraps the GetData.DIRFILE opaque object.
type Dirfile struct {
	name string
	nerr int
	d    *C.DIRFILE
}

// Flags are dirfile opening flags, including encoding methods
type Flags uint

// RDONLY open read-only
const RDONLY Flags = C.GD_RDONLY

// RDWR open read/write
const RDWR Flags = C.GD_RDWR

// FORCEENDIAN override endianness
const FORCEENDIAN Flags = C.GD_FORCE_ENDIAN

// BIGENDIAN assume big-endian raw data
const BIGENDIAN Flags = C.GD_BIG_ENDIAN

// LITTLEENDIAN assume big-endian raw data
const LITTLEENDIAN Flags = C.GD_LITTLE_ENDIAN

// CREAT create dirfile if it doesn't exist
const CREAT Flags = C.GD_CREAT

// EXCL forces creation of dirfile (and fail if it exists)
const EXCL Flags = C.GD_EXCL

// TRUNC truncates the dirfile contents to be empty
const TRUNC Flags = C.GD_TRUNC

// PEDANTIC makes the dirfile instist on strict adherence to standards
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
	dirfile := Dirfile{name: name, d: result}

	errcode := C.gd_error(result)
	if errcode != C.GD_E_OK {
		return dirfile, dirfile.Error()
	}
	return dirfile, nil
}

type RetType uint

const UNKNOWN RetType = C.GD_UNKNOWN
const UINT8 RetType = C.GD_UINT8
const INT8 RetType = C.GD_INT8
const UINT16 RetType = C.GD_UINT16
const INT16 RetType = C.GD_INT16
const UINT32 RetType = C.GD_UINT32
const INT32 RetType = C.GD_INT32
const UINT64 RetType = C.GD_UINT64
const INT64 RetType = C.GD_INT64
const FLOAT32 RetType = C.GD_FLOAT32
const FLOAT64 RetType = C.GD_FLOAT64
const COMPLEX64 RetType = C.GD_COMPLEX64
const COMPLEX128 RetType = C.GD_COMPLEX128
const STRING RetType = C.GD_STRING

func object2type(object interface{}) RetType {
	switch object.(type) {
	case uint8:
		return UINT8
	case int8:
		return INT8
	case uint16:
		return UINT16
	case int16:
		return INT16
	case uint32:
		return UINT32
	case int32:
		return INT32
	case uint64:
		return UINT64
	case int64:
		return INT64
	case float32:
		return FLOAT32
	case float64:
		return FLOAT64
	case complex64:
		return COMPLEX64
	case complex128:
		return COMPLEX128
	}
	return UNKNOWN
}

// Error returns the latest error as a golang error type.
// It uses C API gd_error_string to generate the underlying string.
func (df *Dirfile) Error() error {
	df.nerr += int(C.gd_error_count(df.d))
	cmsg := C.gd_error_string(df.d, (*C.char)(C.NULL), 0)
	defer C.free(unsafe.Pointer(cmsg))
	return errors.New(C.GoString(cmsg))
}

func (df *Dirfile) ErrorCount() int {
	c := df.nerr
	df.nerr = 0
	return c
}

// func (df Dirfile) GetData(fieldcode string, firstFrame, firstSample, numFrames, numSamples int, out []interface{}) int {
// 	fcode := C.CString(fieldcode)
// 	defer C.free(unsafe.Pointer(fcode))
// 	retType := object2type(out[0])
// 	n := C.gd_getdata(df.d, fcode, C.off_t(firstFrame), C.off_t(firstSample),
// 		C.size_t(numFrames), C.size_t(numSamples), C.gd_type_t(retType), unsafe.Pointer(&out[0]))
// 	return int(n)
// }

func (df Dirfile) GetConstantInt32(fieldcode string) (int32, error) {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	result := C.int(32)
	errcode := C.gd_get_constant(df.d, fcode, C.GD_INT32, unsafe.Pointer(&result))
	if errcode != C.GD_E_OK {
		return 0.0, df.Error()
	}
	return int32(result), nil
}

func (df Dirfile) GetConstantFloat32(fieldcode string) (float32, error) {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	result := C.float(3.2)
	errcode := C.gd_get_constant(df.d, fcode, C.GD_FLOAT32, unsafe.Pointer(&result))
	if errcode != C.GD_E_OK {
		return 0.0, df.Error()
	}
	return float32(result), nil
}

func (df Dirfile) GetConstantFloat64(fieldcode string) (float64, error) {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	result := C.double(6.4)
	errcode := C.gd_get_constant(df.d, fcode, C.GD_FLOAT64, unsafe.Pointer(&result))
	if errcode != C.GD_E_OK {
		return 0.0, df.Error()
	}
	return float64(result), nil
}
