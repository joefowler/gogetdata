package getdata

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lgetdata
#include <getdata.h>
#include <stdlib.h>
*/
import "C"
import "unsafe"

// RetType enumerates the data vector return types available (see constants.go
// for specific instances).
type RetType uint

// UINT8 means unsigned 8-bit integers
const UINT8 RetType = C.GD_UINT8

// INT8 means signed 8-bit integers
const INT8 RetType = C.GD_INT8

// UINT16 means unsigned 16-bit integers
const UINT16 RetType = C.GD_UINT16

// INT16 means signed 16-bit integers
const INT16 RetType = C.GD_INT16

// UINT32 means unsigned 32-bit integers
const UINT32 RetType = C.GD_UINT32

// INT32 means signed 32-bit integers
const INT32 RetType = C.GD_INT32

// UINT64 means unsigned 64-bit integers
const UINT64 RetType = C.GD_UINT64

// INT64 means signed 64-bit integers
const INT64 RetType = C.GD_INT64

// FLOAT32 means 32-bit IEEE 754 floats
const FLOAT32 RetType = C.GD_FLOAT32

// FLOAT64 means 64-bit IEEE 754 floats
const FLOAT64 RetType = C.GD_FLOAT64

// COMPLEX64 means complex numbers made of 32-bit IEEE 754 floats
const COMPLEX64 RetType = C.GD_COMPLEX64

// COMPLEX128 means complex numbers made of 64-bit IEEE 754 floats
const COMPLEX128 RetType = C.GD_COMPLEX128

// STRING means character string data
const STRING RetType = C.GD_STRING

// NULLTYPE means the null type
const NULLTYPE RetType = C.GD_NULL

// UNKNOWN means the type could not be detected (passing to GetData always results
// in an error)
const UNKNOWN RetType = C.GD_UNKNOWN

// parray2type accepts a pointer to a slice of numeric values and returns the
// matching RetType from the GetData library and an unsafe.Pointer to the
// first value in the underlying array.
// Differs from array2type in that you use this if you want to be able to resize.
func parray2type(a interface{}) (RetType, unsafe.Pointer) {
	switch v := a.(type) {
	case *[]uint8:
		return UINT8, unsafe.Pointer(&(*v)[0])
	case *[]int8:
		return INT8, unsafe.Pointer(&(*v)[0])
	case *[]uint16:
		return UINT16, unsafe.Pointer(&(*v)[0])
	case *[]int16:
		return INT16, unsafe.Pointer(&(*v)[0])
	case *[]uint32:
		return UINT32, unsafe.Pointer(&(*v)[0])
	case *[]int32:
		return INT32, unsafe.Pointer(&(*v)[0])
	case *[]uint64:
		return UINT64, unsafe.Pointer(&(*v)[0])
	case *[]int64:
		return INT64, unsafe.Pointer(&(*v)[0])
	case *[]float32:
		return FLOAT32, unsafe.Pointer(&(*v)[0])
	case *[]float64:
		return FLOAT64, unsafe.Pointer(&(*v)[0])
	case *[]complex64:
		return COMPLEX64, unsafe.Pointer(&(*v)[0])
	case *[]complex128:
		return COMPLEX128, unsafe.Pointer(&(*v)[0])
	case *[]string:
		return STRING, unsafe.Pointer(&(*v)[0])
	default:
		return UNKNOWN, nil
	}
}

// array2type accepts a slice of numeric values and returns the
// matching RetType from the GetData library and an unsafe.Pointer to the
// first value in the underlying array.
func array2type(a interface{}) (RetType, unsafe.Pointer, int) {
	switch v := a.(type) {
	case []uint8:
		return UINT8, unsafe.Pointer(&v[0]), len(v)
	case []int8:
		return INT8, unsafe.Pointer(&v[0]), len(v)
	case []uint16:
		return UINT16, unsafe.Pointer(&v[0]), len(v)
	case []int16:
		return INT16, unsafe.Pointer(&v[0]), len(v)
	case []uint32:
		return UINT32, unsafe.Pointer(&v[0]), len(v)
	case []int32:
		return INT32, unsafe.Pointer(&v[0]), len(v)
	case []uint64:
		return UINT64, unsafe.Pointer(&v[0]), len(v)
	case []int64:
		return INT64, unsafe.Pointer(&v[0]), len(v)
	case []float32:
		return FLOAT32, unsafe.Pointer(&v[0]), len(v)
	case []float64:
		return FLOAT64, unsafe.Pointer(&v[0]), len(v)
	case []complex64:
		return COMPLEX64, unsafe.Pointer(&v[0]), len(v)
	case []complex128:
		return COMPLEX128, unsafe.Pointer(&v[0]), len(v)
	case []string:
		return STRING, unsafe.Pointer(&v[0]), len(v)
	default:
		return UNKNOWN, nil, 0
	}
}

// pointer2type accepts a pointer to a numeric value or string and returns the
// matching RetType from the GetData library and an unsafe.Pointer to the
// underlying value.
func pointer2type(object interface{}) (RetType, unsafe.Pointer) {
	switch p := object.(type) {
	case *uint8:
		return UINT8, unsafe.Pointer(p)
	case *int8:
		return INT8, unsafe.Pointer(p)
	case *uint16:
		return UINT16, unsafe.Pointer(p)
	case *int16:
		return INT16, unsafe.Pointer(p)
	case *uint32:
		return UINT32, unsafe.Pointer(p)
	case *int32:
		return INT32, unsafe.Pointer(p)
	case *uint64:
		return UINT64, unsafe.Pointer(p)
	case *int64:
		return INT64, unsafe.Pointer(p)
	case *float32:
		return FLOAT32, unsafe.Pointer(p)
	case *float64:
		return FLOAT64, unsafe.Pointer(p)
	case *complex64:
		return COMPLEX64, unsafe.Pointer(p)
	case *complex128:
		return COMPLEX128, unsafe.Pointer(p)
	case *string:
		return STRING, unsafe.Pointer(p)
	}
	return UNKNOWN, nil
}

// value2type accepts a numeric value or string and returns the
// matching RetType from the GetData library and an unsafe.Pointer to the
// underlying value.
func value2type(object interface{}) (RetType, unsafe.Pointer) {
	switch p := object.(type) {
	case uint8:
		return UINT8, unsafe.Pointer(&p)
	case int8:
		return INT8, unsafe.Pointer(&p)
	case uint16:
		return UINT16, unsafe.Pointer(&p)
	case int16:
		return INT16, unsafe.Pointer(&p)
	case uint32:
		return UINT32, unsafe.Pointer(&p)
	case int32:
		return INT32, unsafe.Pointer(&p)
	case uint64:
		return UINT64, unsafe.Pointer(&p)
	case int64:
		return INT64, unsafe.Pointer(&p)
	case float32:
		return FLOAT32, unsafe.Pointer(&p)
	case float64:
		return FLOAT64, unsafe.Pointer(&p)
	case complex64:
		return COMPLEX64, unsafe.Pointer(&p)
	case complex128:
		return COMPLEX128, unsafe.Pointer(&p)
	case string:
		return STRING, unsafe.Pointer(&p)
	}
	return UNKNOWN, nil
}
