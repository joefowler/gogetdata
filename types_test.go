package getdata

import (
	"testing"
)

func TestTypeInference(t *testing.T) {
	type TypeTest struct {
		arVal interface{}
		arPtr interface{}
		elVal interface{}
		elPtr interface{}
		tName string
		rType RetType
	}

	st := make([]string, 1)
	u8 := make([]uint8, 1)
	i8 := make([]int8, 1)
	u16 := make([]uint16, 1)
	i16 := make([]int16, 1)
	u32 := make([]uint32, 1)
	i32 := make([]int32, 1)
	u64 := make([]uint64, 1)
	i64 := make([]int64, 1)
	f32 := make([]float32, 1)
	f64 := make([]float64, 1)
	c64 := make([]complex64, 1)
	c128 := make([]complex128, 1)
	var tests = []TypeTest{
		{st, &st, st[0], &st[0], "STRING", STRING},
		{i8, &i8, i8[0], &i8[0], "INT8", INT8},
		{u8, &u8, u8[0], &u8[0], "UINT8", UINT8},
		{i16, &i16, i16[0], &i16[0], "INT16", INT16},
		{u16, &u16, u16[0], &u16[0], "UINT16", UINT16},
		{i32, &i32, i32[0], &i32[0], "INT32", INT32},
		{u32, &u32, u32[0], &u32[0], "UINT32", UINT32},
		{i64, &i64, i64[0], &i64[0], "INT64", INT64},
		{u64, &u64, u64[0], &u64[0], "UINT64", UINT64},
		{f32, &f32, f32[0], &f32[0], "FLOAT32", FLOAT32},
		{f64, &f64, f64[0], &f64[0], "FLOAT64", FLOAT64},
		{c64, &c64, c64[0], &c64[0], "COMPLEX64", COMPLEX64},
		{c128, &c128, c128[0], &c128[0], "COMPLEX128", COMPLEX128},
	}
	for _, test := range tests {
		if tval, ptr := pointer2type(test.elPtr); tval != test.rType {
			t.Errorf("pointer2type(&val) returns 0x%x, %p; want %s=0x%x, nil", tval, ptr, test.tName, test.rType)
		}
		if tval, ptr := value2type(test.elVal); tval != test.rType {
			t.Errorf("value2type(val) returns 0x%x, %p; want %s=0x%x, nil", tval, ptr, test.tName, test.rType)
		}
		if tval, _ := parray2type(test.arPtr); tval != test.rType {
			t.Errorf("parray2type(&[]float64) returns 0x%x, want %s=0x%x", tval, test.tName, test.rType)
		}
		if tval, _, _ := array2type(test.arVal); tval != test.rType {
			t.Errorf("array2type([]float64) returns 0x%x, want %s=0x%x", tval, test.tName, test.rType)
		}
	}

	type Incorrect int64
	var wrong Incorrect
	if tval, ptr := pointer2type(wrong); tval != UNKNOWN || ptr != nil {
		t.Errorf("pointer2type(val) returns 0x%x, %p; want UNKNOWN=0x%x, nil", tval, ptr, UNKNOWN)
	}
	if tval, ptr := value2type(wrong); tval != UNKNOWN || ptr != nil {
		t.Errorf("value2type(val) returns 0x%x, %p; want UNKNOWN=0x%x, nil", tval, ptr, UNKNOWN)
	}
	if tval, ptr := parray2type(wrong); tval != UNKNOWN || ptr != nil {
		t.Errorf("parray2type(val) returns 0x%x, %p; want UNKNOWN=0x%x, nil", tval, ptr, UNKNOWN)
	}
	if tval, ptr, _ := array2type(wrong); tval != UNKNOWN || ptr != nil {
		t.Errorf("array2type(val) returns 0x%x, %p; want UNKNOWN=0x%x, nil", tval, ptr, UNKNOWN)
	}
}
