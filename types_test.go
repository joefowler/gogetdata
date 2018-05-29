package getdata

import (
	"testing"
)

func TestArrayTypes(t *testing.T) {
	u := make([]uint8, 1)
	if tval, ptr := parray2type(u); tval != UNKNOWN || ptr != nil {
		t.Errorf("array2type(array) returns 0x%x, %p; want UNKNOWN=0x%x, nil", tval, ptr, UNKNOWN)
	}

	s := make([]string, 1)
	if tval, ptr := parray2type(&s); tval != UNKNOWN || ptr != nil {
		t.Errorf("parray2type([]string) returns 0x%x, %p; want UNKNOWN=0x%x, nil", tval, ptr, UNKNOWN)
	}

	u8 := make([]uint8, 1)
	if tval, _ := parray2type(&u8); tval != UINT8 {
		t.Errorf("parray2type(&[]uint8) returns 0x%x, want UINT8=0x%x", tval, UINT8)
	}
	if tval, _, _ := array2type(u8); tval != UINT8 {
		t.Errorf("array2type([]uint8) returns 0x%x, want UINT8=0x%x", tval, UINT8)
	}

	i8 := make([]int8, 1)
	if tval, _ := parray2type(&i8); tval != INT8 {
		t.Errorf("parray2type(&[]int8) returns 0x%x, want INT8=0x%x", tval, INT8)
	}
	if tval, _, _ := array2type(i8); tval != INT8 {
		t.Errorf("array2type([]int8) returns 0x%x, want INT8=0x%x", tval, INT8)
	}

	u16 := make([]uint16, 1)
	if tval, _ := parray2type(&u16); tval != UINT16 {
		t.Errorf("parray2type(&[]uint16) returns 0x%x, want UINT16=0x%x", tval, UINT16)
	}
	if tval, _, _ := array2type(u16); tval != UINT16 {
		t.Errorf("array2type([]uint16) returns 0x%x, want UINT16=0x%x", tval, UINT16)
	}

	i16 := make([]int16, 1)
	if tval, _ := parray2type(&i16); tval != INT16 {
		t.Errorf("parray2type(&[]int16) returns 0x%x, want INT16=0x%x", tval, INT16)
	}
	if tval, _, _ := array2type(i16); tval != INT16 {
		t.Errorf("array2type([]int16) returns 0x%x, want INT16=0x%x", tval, INT16)
	}

	u32 := make([]uint32, 1)
	if tval, _ := parray2type(&u32); tval != UINT32 {
		t.Errorf("parray2type(&[]uint32) returns 0x%x, want UINT32=0x%x", tval, UINT32)
	}
	if tval, _, _ := array2type(u32); tval != UINT32 {
		t.Errorf("array2type([]uint32) returns 0x%x, want UINT32=0x%x", tval, UINT32)
	}

	i32 := make([]int32, 1)
	if tval, _ := parray2type(&i32); tval != INT32 {
		t.Errorf("parray2type(&[]int32) returns 0x%x, want INT32=0x%x", tval, INT32)
	}
	if tval, _, _ := array2type(i32); tval != INT32 {
		t.Errorf("array2type([]int32) returns 0x%x, want INT32=0x%x", tval, INT32)
	}

	u64 := make([]uint64, 1)
	if tval, _ := parray2type(&u64); tval != UINT64 {
		t.Errorf("parray2type(&[]uint64) returns 0x%x, want UINT64=0x%x", tval, UINT64)
	}
	if tval, _, _ := array2type(u64); tval != UINT64 {
		t.Errorf("array2type([]uint64) returns 0x%x, want UINT64=0x%x", tval, UINT64)
	}

	i64 := make([]int64, 1)
	if tval, _ := parray2type(&i64); tval != INT64 {
		t.Errorf("parray2type(&[]int64) returns 0x%x, want INT64=0x%x", tval, INT64)
	}
	if tval, _, _ := array2type(i64); tval != INT64 {
		t.Errorf("array2type([]int64) returns 0x%x, want INT64=0x%x", tval, INT64)
	}

	f32 := make([]float32, 1)
	if tval, _ := parray2type(&f32); tval != FLOAT32 {
		t.Errorf("parray2type(&[]float32) returns 0x%x, want FLOAT32=0x%x", tval, FLOAT32)
	}
	if tval, _, _ := array2type(f32); tval != FLOAT32 {
		t.Errorf("array2type([]float32) returns 0x%x, want FLOAT32=0x%x", tval, FLOAT32)
	}

	f64 := make([]float64, 1)
	if tval, _ := parray2type(&f64); tval != FLOAT64 {
		t.Errorf("parray2type(&[]float64) returns 0x%x, want FLOAT64=0x%x", tval, FLOAT64)
	}
	if tval, _, _ := array2type(f64); tval != FLOAT64 {
		t.Errorf("array2type([]float64) returns 0x%x, want FLOAT64=0x%x", tval, FLOAT64)
	}

	c64 := make([]complex64, 1)
	if tval, _ := parray2type(&c64); tval != COMPLEX64 {
		t.Errorf("parray2type(&[]complex64) returns 0x%x, want COMPLEX64=0x%x", tval, COMPLEX64)
	}
	if tval, _, _ := array2type(c64); tval != COMPLEX64 {
		t.Errorf("array2type([]complex64) returns 0x%x, want COMPLEX64=0x%x", tval, COMPLEX64)
	}

	c128 := make([]complex128, 1)
	if tval, _ := parray2type(&c128); tval != COMPLEX128 {
		t.Errorf("parray2type(&[]complex128) returns 0x%x, want COMPLEX128=0x%x", tval, COMPLEX128)
	}
	if tval, _, _ := array2type(c128); tval != COMPLEX128 {
		t.Errorf("array2type([]complex128) returns 0x%x, want COMPLEX128=0x%x", tval, COMPLEX128)
	}
}

func TestPointerTypes(t *testing.T) {
	var u8 uint8
	if tval, ptr := pointer2type(&u8); tval != UINT8 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want UINT8=0x%x, nil", tval, ptr, UINT8)
	}
	var i8 int8
	if tval, ptr := pointer2type(&i8); tval != INT8 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want INT8=0x%x, nil", tval, ptr, INT8)
	}
	var u16 uint16
	if tval, ptr := pointer2type(&u16); tval != UINT16 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want UINT16=0x%x, nil", tval, ptr, UINT16)
	}
	var i16 int16
	if tval, ptr := pointer2type(&i16); tval != INT16 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want INT16=0x%x, nil", tval, ptr, INT16)
	}
	var u32 uint32
	if tval, ptr := pointer2type(&u32); tval != UINT32 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want UINT32=0x%x, nil", tval, ptr, UINT32)
	}
	var i32 int32
	if tval, ptr := pointer2type(&i32); tval != INT32 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want INT32=0x%x, nil", tval, ptr, INT32)
	}
	var u64 uint64
	if tval, ptr := pointer2type(&u64); tval != UINT64 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want UINT64=0x%x, nil", tval, ptr, UINT64)
	}
	var i64 int64
	if tval, ptr := pointer2type(&i64); tval != INT64 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want INT64=0x%x, nil", tval, ptr, INT64)
	}
	var f32 float32
	if tval, ptr := pointer2type(&f32); tval != FLOAT32 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want FLOAT32=0x%x, nil", tval, ptr, FLOAT32)
	}
	var f64 float64
	if tval, ptr := pointer2type(&f64); tval != FLOAT64 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want FLOAT64=0x%x, nil", tval, ptr, FLOAT64)
	}
	var c64 complex64
	if tval, ptr := pointer2type(&c64); tval != COMPLEX64 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want COMPLEX64=0x%x, nil", tval, ptr, COMPLEX64)
	}
	var c128 complex128
	if tval, ptr := pointer2type(&c128); tval != COMPLEX128 {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want COMPLEX128=0x%x, nil", tval, ptr, COMPLEX128)
	}
	var s string
	if tval, ptr := pointer2type(&s); tval != STRING {
		t.Errorf("pointer2type(&val) returns 0x%x, %p; want STRING=0x%x, nil", tval, ptr, STRING)
	}

	if tval, ptr := pointer2type(u8); tval != UNKNOWN || ptr != nil {
		t.Errorf("pointer2type(val) returns 0x%x, %p; want UNKNOWN=0x%x, nil", tval, ptr, UNKNOWN)
	}
}
