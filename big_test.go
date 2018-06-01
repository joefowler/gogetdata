package getdata

import (
	"fmt"
	"log"
	"math"
	"os"
	"strings"
	"testing"
)

func createTestDirfile(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Fatal("Could not remove dirfile: ", err)
	}
	err = os.Mkdir(dir, 0775)
	if err != nil {
		log.Fatal("Could not create dirfile: ", err)
	}
	data := make([]byte, 80)
	for i := 0; i < 80; i++ {
		data[i] = byte(i + 1)
	}
	f, err := os.Create(fmt.Sprintf("%s/data", dir))
	if err != nil {
		log.Fatal("Could not create dirfile/data")
	}
	f.Write(data)
	f.Close()

	f, err = os.Create(fmt.Sprintf("%s/format", dir))
	if err != nil {
		log.Fatal("Could not create dirfile/format")
	}
	fmt.Fprintf(f, `/ENDIAN little
data RAW INT8 8
lincom LINCOM data 1.1 2.2 INDEX 2.2 3.3;4.4 linterp const const
/META data mstr STRING "This is a string constant."
/META data mconst CONST COMPLEX128 3.3;4.4
/META data mcarray CARRAY FLOAT64 1.9 2.8 3.7 4.6 5.5
/META data mlut LINTERP DATA ./lut
const CONST FLOAT64 5.5
carray CARRAY FLOAT64 1.1 2.2 3.3 4.4 5.5 6.6
linterp LINTERP data ./lut
polynom POLYNOM data 1.1 2.2 2.2 3.3;4.4 const const
bit BIT data 3 4
sbit SBIT data 5 6
mplex MPLEX data sbit 1 10
mult MULTIPLY data sbit
div DIVIDE mult bit
recip RECIP div 6.5;4.3
phase PHASE data 11
window WINDOW linterp mult LT 4.1
/ALIAS alias data
string STRING "Zaphod Beeblebrox"
sarray SARRAY one two three four five six seven
data/msarray SARRAY eight nine ten eleven twelve
indir INDIR data carray
sindir SINDIR data sarray
`)
	f.Close()

	f, err = os.Create(fmt.Sprintf("%s/form2", dir))
	if err != nil {
		log.Fatal("Could not create dirfile/form2")
	}
	fmt.Fprintf(f, "const2 CONST INT8 -19")
	f.Close()
}

func removeTestDirfile(dir string) {
	// TODO: eventually want to remove this dirfile when the test is done.
}

func TestRead(t *testing.T) {
	dir := "dirfile"
	createTestDirfile(dir)
	defer removeTestDirfile(dir)

	// #1: read-only open check
	d, err := OpenDirfile(dir, RDONLY)
	if err != nil {
		t.Errorf("Could not open dirfile read-only")
	}
	if c := d.ErrorCount(); c > 0 {
		t.Errorf("Error count %d when open dirfile read-only, want 0", c)
	}
	err = d.Close()
	if err != nil {
		t.Errorf("Could not close dirfile read-only")
	}

	// #1b: read-only open check on non-existing file
	d, err = OpenDirfile("randomfile", RDONLY)
	if err == nil {
		t.Errorf("Could open a non-existent dirfile, want error")
	}
	if c := d.ErrorCount(); c != 1 {
		t.Errorf("Error count %d when open dirfile read-only, want 1", c)
	}

	// #2: read-write open check
	// d, err := OpenDirfile(dir, RDWR)
	d, err = OpenDirfile(dir, RDWR)
	if err != nil {
		t.Errorf("Could not open dirfile read-write")
	}
	if c := d.ErrorCount(); c > 0 {
		t.Errorf("Error count %d when open dirfile read-write, want 0", c)
	}

	// #3: getdata (int) check
	u1 := make([]int8, 8)
	n, err := d.GetData("data", 5, 0, 1, 0, &u1)
	if err != nil {
		t.Error("Could not GetData: ", err)
	} else if len(u1) < 8 {
		t.Errorf("GetData out has len=%d (cap=%d), want 8", len(u1), cap(u1))
	} else if n != 8 {
		t.Errorf("GetData returned %d, want 8", n)
	} else {
		for i := 0; i < 8; i++ {
			if u1[i] != int8(41+i) {
				t.Errorf("GetData out[%d]=%d, want %d", i, u1[i], 41+i)
			}
		}
	}

	// #6: getdata (int64) check and check FRAMEHERE
	u2 := make([]uint64, 8)
	n, err = d.GetData("data", FRAMEHERE, 0, 1, 0, &u2)
	if err != nil {
		t.Error("Could not GetData: ", err)
	} else if len(u2) < 8 {
		t.Errorf("GetData out has len=%d (cap=%d), want 8", len(u2), cap(u2))
	} else if n != 8 {
		t.Errorf("GetData returned %d, want 8", n)
	} else {
		for i := 0; i < 8; i++ {
			if u2[i] != uint64(49+i) {
				t.Errorf("GetData out[%d]=%d, want %d", i, u2[i], 41+i)
			}
		}
	}

	// #8: getdata (float64) check
	u3 := make([]float64, 8)
	n, err = d.GetData("data", 5, 0, 1, 0, &u3)
	if err != nil {
		t.Error("Could not GetData: ", err)
	} else if len(u3) < 8 {
		t.Errorf("GetData out has len=%d (cap=%d), want 8", len(u3), cap(u3))

	} else if n != 8 {
		t.Errorf("GetData returned %d, want 8", n)

	} else {
		for i := 0; i < 8; i++ {
			if u3[i] != float64(41+i) {
				t.Errorf("GetData out[%d]=%.5f, want %d.00000", i, u3[i], 41+i)
			}
		}
	}

	// #10: getdata (complex128) check
	u4 := make([]complex128, 8)
	n, err = d.GetData("data", 5, 0, 1, 0, &u4)
	if err != nil {
		t.Error("Could not GetData: ", err)
	} else if len(u4) < 8 {
		t.Errorf("GetData out has len=%d (cap=%d), want 8", len(u4), cap(u4))

	} else if n != 8 {
		t.Errorf("GetData returned %d, want 8", n)

	} else {
		for i := 0; i < 8; i++ {
			if u4[i] != complex(41.0+float64(i), 0.0) {
				t.Errorf("GetData out[%d]=%.5f+%5fi, want %d.00000", i, real(u4[i]), imag(u4[i]), 41+i)
			}
		}
	}

	// #11: Check for appropriate errors in GetData
	n, err = d.GetData("data", 5, 0, 0, 0, &u3)
	if err == nil || n > 0 {
		t.Errorf("GetData with 0 frames/samples requested returns (%d, %s), want (0, error)", n, err)
	}
	n, err = d.GetData("data", 5, 0, 1, 0, u3)
	if err == nil || n > 0 {
		t.Errorf("GetData with out argument not a slice pointer returns (%d, %s), want (0, error)", n, err)
	}

	// #12: constant (int) check
	i32, err := d.GetConstantInt32("const")
	if err != nil {
		t.Error("Could not GetConstantInt32: ", err)
	} else if i32 != 5 {
		t.Errorf("GetConstantInt32 returns %d, want 5", i32)
	}
	i32 = -9
	err = d.GetConstant("const", &i32)
	if err != nil {
		t.Error("Could not GetConstant for &int32: ", err)
	} else if i32 != 5 {
		t.Errorf("GetConstant for int32 returns %d, want 5", i32)
	}
	i32, err = d.GetConstantInt32("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantInt32 returned %d on non-existent field, want error", i32)
	}

	// #15: constant (int64) check
	i64, err := d.GetConstantInt64("const")
	if err != nil {
		t.Error("Could not GetConstantInt64: ", err)
	} else if i64 != 5 {
		t.Errorf("GetConstantInt64 returns %d, want 5", i64)
	}
	i64 = -9
	err = d.GetConstant("const", &i64)
	if err != nil {
		t.Error("Could not GetConstant for int64: ", err)
	} else if i64 != 5 {
		t.Errorf("GetConstant for int64 returns %d, want 5", i64)
	}
	i64, err = d.GetConstantInt64("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantInt64 returned %d on non-existent field, want error", i64)
	}

	// #17: constant (float) check
	f32, err := d.GetConstantFloat32("const")
	if err != nil {
		t.Error("Could not GetConstantFloat32: ", err)
	} else if f32 != 5.5 {
		t.Errorf("GetConstantFloat32 returns %f, want 5.5", f32)
	}
	f32 = -9
	err = d.GetConstant("const", &f32)
	if err != nil {
		t.Error("Could not GetConstant for float32: ", err)
	} else if f32 != 5.5 {
		t.Errorf("GetConstant for float32 returns %f, want 5.5", f32)
	}
	f32, err = d.GetConstantFloat32("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantFloat32 returned %f on non-existent field, want error", f32)
	}

	f64, err := d.GetConstantFloat64("const")
	if err != nil {
		t.Errorf("Could not GetConstantFloat64")
	} else if f64 != 5.5 {
		t.Errorf("GetConstantFloat64 returns %f, want 5.5", f64)
	}
	f64 = -9.0
	err = d.GetConstant("const", &f64)
	if err != nil {
		t.Errorf("Could not GetConstant for float64")
	} else if f64 != 5.5 {
		t.Errorf("GetConstant for float64 returns %f, want 5.5", f64)
	}
	f64, err = d.GetConstantFloat64("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantFloat64 returned %f on non-existent field, want error", f64)
	}

	// #19: constant (complex) check
	c64, err := d.GetConstantComplex64("const")
	if err != nil {
		t.Errorf("Could not GetConstantComplex64")
	} else if c64 != 5.5 {
		t.Errorf("GetConstantComplex64 returns %f, want 5.5", c64)
	}
	c64 = complex(-9, 4)
	err = d.GetConstant("const", &c64)
	if err != nil {
		t.Errorf("Could not GetConstant for complex64")
	} else if c64 != 5.5 {
		t.Errorf("GetConstant for complex64 returns %f, want 5.5", c64)
	}
	c64, err = d.GetConstantComplex64("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantComplex64 returned %f on non-existent field, want error", c64)
	}

	c128, err := d.GetConstantComplex128("const")
	if err != nil {
		t.Errorf("Could not GetConstantComplex128")
	} else if c128 != 5.5 {
		t.Errorf("GetConstant for complex128 returns %f, want 5.5", c128)
	}
	c128 = complex(-7, 5)
	err = d.GetConstant("const", &c128)
	if err != nil {
		t.Errorf("Could not GetConstant for complex128")
	} else if c128 != 5.5 {
		t.Errorf("GetConstantComplex128 returns %f, want 5.5", c128)
	}
	c128, err = d.GetConstantComplex128("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantComplex128 returned %f on non-existent field, want error", c128)
	}

	// #23: NFields check
	nfields := d.NFields()
	if nfields != 20 {
		t.Errorf("Nfields = %d, want 20", nfields)
	}

	// #25: FieldList check
	fields := d.FieldList()
	if len(fields) != int(nfields) {
		t.Errorf("FieldList length is %d, want %d", len(fields), nfields)
	} else {
		truenames := []string{"bit", "div", "data", "mult", "sbit", "INDEX",
			"alias", "const", "indir", "mplex", "phase", "recip", "carray", "lincom",
			"sarray", "sindir", "string", "window", "linterp", "polynom"}
		for i := 0; i < int(nfields); i++ {
			if fields[i] != truenames[i] {
				t.Errorf("FieldList[%d]=\"%s\", want \"%s\"", i, fields[i], truenames[i])
			}
		}
	}

	// #26: NMFields check
	nmfields := d.NMFields("data")
	if nmfields != 5 {
		t.Errorf("Nfields(\"data\") returned %d, want 5", nmfields)
	}

	// #27: MFieldList check
	mfields := d.MFieldList("data")
	if len(mfields) != int(nmfields) {
		t.Errorf("MFieldList length is %d, want %d", len(mfields), nmfields)
	} else {
		truemnames := []string{"mstr", "mconst", "mcarray", "mlut", "msarray"}
		for i := 0; i < int(nmfields); i++ {
			if mfields[i] != truemnames[i] {
				t.Errorf("FieldList[%d]=\"%s\", want \"%s\"", i, mfields[i], truemnames[i])
			}
		}
	}

	// #28: NFrames check
	nf := d.NFrames()
	if nf != 10 {
		t.Errorf("NFrames returned %d, want 10", nf)
	}

	// #30: PutData int32 check
	dataToPut := []int32{13, 14, 15, 16}
	n, err = d.PutData("data", 5, 1, dataToPut)
	if err != nil {
		t.Errorf("Could not PutData in test 30")
	} else if n != len(dataToPut) {
		t.Errorf("PutData returned %d in test 30, want %d", n, len(dataToPut))
	}
	testData := make([]int32, 8)
	n, err = d.GetData("data", 5, 0, 1, 0, &testData)
	if err != nil {
		t.Errorf("Could not GetData in test 30")
	} else if n != len(testData) {
		t.Errorf("GetData returned %d in test 30, want %d", n, len(testData))
	} else {
		expectedData := []int32{41, 13, 14, 15, 16, 46, 47, 48}
		for i := 0; i < 8; i++ {
			if testData[i] != expectedData[i] {
				t.Errorf("GetData returned d[%d]=%d in test 30, want %d", i, testData[i], expectedData[i])
			}
		}
	}

	// #33: PutData int64 check
	dataToPut64 := []int64{13, 14, 15, 16}
	n, err = d.PutData("data", 5, 1, dataToPut64)
	if err != nil {
		t.Errorf("Could not PutData in test 33")
	} else if n != len(dataToPut64) {
		t.Errorf("PutData returned %d in test 33, want %d", n, len(dataToPut64))
	}
	n, err = d.GetData("data", 5, 0, 1, 0, &testData)
	if err != nil {
		t.Errorf("Could not GetData in test 33")
	} else if n != len(testData) {
		t.Errorf("GetData returned %d in test 33, want %d", n, len(testData))
	} else {
		expectedData := []int32{41, 13, 14, 15, 16, 46, 47, 48}
		for i := 0; i < 8; i++ {
			if testData[i] != expectedData[i] {
				t.Errorf("GetData returned d[%d]=%d in test 33, want %d", i, testData[i], expectedData[i])
			}
		}
	}

	// #35: PutData float64 check
	dataToPutF := []float64{13, 14, 15, 16}
	n, err = d.PutData("data", 5, 1, dataToPutF)
	if err != nil {
		t.Errorf("Could not PutData in test 35")
	} else if n != len(dataToPutF) {
		t.Errorf("PutData returned %d in test 35, want %d", n, len(dataToPutF))
	}
	n, err = d.GetData("data", 5, 0, 1, 0, &testData)
	if err != nil {
		t.Errorf("Could not GetData in test 35")
	} else if n != len(testData) {
		t.Errorf("GetData returned %d in test 35, want %d", n, len(testData))
	} else {
		expectedData := []int32{41, 13, 14, 15, 16, 46, 47, 48}
		for i := 0; i < 8; i++ {
			if testData[i] != expectedData[i] {
				t.Errorf("GetData returned d[%d]=%d in test 35, want %d", i, testData[i], expectedData[i])
			}
		}
	}

	// #37: PutData complex128 check
	dataToPutC := []complex128{13, 14, 15, 16}
	n, err = d.PutData("data", 5, 1, dataToPutC)
	if err != nil {
		t.Errorf("Could not PutData in test 37")
	} else if n != len(dataToPutC) {
		t.Errorf("PutData returned %d in test 37, want %d", n, len(dataToPutC))
	}
	n, err = d.GetData("data", 5, 0, 1, 0, &testData)
	if err != nil {
		t.Errorf("Could not GetData in test 37")
	} else if n != len(testData) {
		t.Errorf("GetData returned %d in test 37, want %d", n, len(testData))
	} else {
		expectedData := []int32{41, 13, 14, 15, 16, 46, 47, 48}
		for i := 0; i < 8; i++ {
			if testData[i] != expectedData[i] {
				t.Errorf("GetData returned d[%d]=%d in test 37, want %d", i, testData[i], expectedData[i])
			}
		}
	}

	// #40: Entry (raw) check
	ent, err := d.Entry("data")
	if err != nil {
		t.Error("Could not get Entry for raw type: ", err)
	} else {
		if ent.fieldType != RAWENTRY {
			t.Errorf("Entry gets field type 0x%x, want RAW=0x%x", ent.fieldType, RAWENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		if ent.dataType != INT8 {
			t.Errorf("Entry gets data type 0x%x, want INT8=0x%x", ent.dataType, INT8)
		}
		if ent.spf != 8 {
			t.Errorf("Entry gets spf=%d, want 8", ent.spf)
		}
	}

	// #42: Entry (lincom) check
	ent, err = d.Entry("lincom")
	if err != nil {
		t.Error("Could not get Entry for lincom type: ", err)
	} else {
		if ent.fieldType != LINCOMENTRY {
			t.Errorf("Entry gets field type 0x%x, want LINCOM=0x%x", ent.fieldType, LINCOMENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		if ent.nFields != 3 {
			t.Errorf("Entry gets nFields=%d, want 3", ent.nFields)
		}
		expectedf := []string{"data", "INDEX", "linterp"}
		expectedm := []complex128{1.1, 2.2, 5.5}
		expectedb := []complex128{2.2, complex(3.3, 4.4), 5.5}
		for i := 0; i < ent.nFields; i++ {
			if ent.inFields[i] != expectedf[i] {
				t.Errorf("Entry gets inFields[%d]=\"%s\", want \"%s\"", i, ent.inFields[i], expectedf[i])
			}
			if ent.m[i] != real(expectedm[i]) {
				t.Errorf("Entry gets m[%d]=%f, want %f", i, ent.m[i], expectedm[i])
			}
			if ent.cm[i] != expectedm[i] {
				t.Errorf("Entry gets cm[%d]=%f, want %f", i, ent.cm[i], expectedm[i])
			}
			if ent.b[i] != real(expectedb[i]) {
				t.Errorf("Entry gets b[%d]=%f, want %f", i, ent.b[i], expectedb[i])
			}
			if ent.cb[i] != expectedb[i] {
				t.Errorf("Entry gets cb[%d]=%f, want %f", i, ent.cb[i], expectedb[i])
			}
		}
	}

	// #44: Entry (polynom) check
	ent, err = d.Entry("polynom")
	if err != nil {
		t.Error("Could not get Entry for polynom type: ", err)
	} else {
		if ent.fieldType != POLYNOMENTRY {
			t.Errorf("Entry gets field type 0x%x, want POLYNOM=0x%x", ent.fieldType, POLYNOMENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		if ent.polyOrder != 5 {
			t.Errorf("Entry gets polyOrder=%d, want 5", ent.polyOrder)
		}
		expectedf := "data"
		expecteda := []complex128{1.1, 2.2, 2.2, complex(3.3, 4.4), 5.5, 5.5}
		if ent.inFields[0] != expectedf {
			t.Errorf("Entry gets inFields[%d]=\"%s\", want \"%s\"", 0, ent.inFields[0], expectedf)
		}
		for i := 0; i < ent.polyOrder; i++ {
			if ent.a[i] != real(expecteda[i]) {
				t.Errorf("Entry gets a[%d]=%f, want %f", i, ent.a[i], expecteda[i])
			}
			if ent.ca[i] != expecteda[i] {
				t.Errorf("Entry gets ca[%d]=%f, want %f", i, ent.ca[i], expecteda[i])
			}
		}
	}

	// #45: Entry (linterp) check
	ent, err = d.Entry("linterp")
	if err != nil {
		t.Error("Could not get Entry for linterp type: ", err)
	} else {
		if ent.fieldType != LINTERPENTRY {
			t.Errorf("Entry gets field type 0x%x, want LINTERP=0x%x", ent.fieldType, LINTERPENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		expectedf := "data"
		if ent.inFields[0] != expectedf {
			t.Errorf("Entry gets inFields[%d]=\"%s\", want \"%s\"", 0, ent.inFields[0], expectedf)
		}
		if ent.table != "./lut" {
			t.Errorf("Entry gets table=\"%s\", want \"%s\"", ent.table, "./lut")
		}
	}

	// #46: Entry (bits) check
	ent, err = d.Entry("bit")
	if err != nil {
		t.Error("Could not get Entry for bit type: ", err)
	} else {
		if ent.fieldType != BITENTRY {
			t.Errorf("Entry gets field type 0x%x, want BIT=0x%x", ent.fieldType, BITENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		expectedf := "data"
		if ent.inFields[0] != expectedf {
			t.Errorf("Entry gets inFields[%d]=\"%s\", want \"%s\"", 0, ent.inFields[0], expectedf)
		}
		if ent.numbits != 4 {
			t.Errorf("Entry gets numbits=%d, want 4", ent.numbits)
		}
		if ent.bitnum != 3 {
			t.Errorf("Entry gets bitnum=%d, want 3", ent.bitnum)
		}
	}

	// #47: Entry (sbits) check
	ent, err = d.Entry("sbit")
	if err != nil {
		t.Error("Could not get Entry for sbit type: ", err)
	} else {
		if ent.fieldType != SBITENTRY {
			t.Errorf("Entry gets field type 0x%x, want SBIT=0x%x", ent.fieldType, SBITENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		expectedf := "data"
		if ent.inFields[0] != expectedf {
			t.Errorf("Entry gets inFields[%d]=\"%s\", want \"%s\"", 0, ent.inFields[0], expectedf)
		}
		if ent.numbits != 6 {
			t.Errorf("Entry gets numbits=%d, want 6", ent.numbits)
		}
		if ent.bitnum != 5 {
			t.Errorf("Entry gets bitnum=%d, want 5", ent.bitnum)
		}
	}

	// #48: Entry (mult) check
	ent, err = d.Entry("mult")
	if err != nil {
		t.Error("Could not get Entry for mult type: ", err)
	} else {
		if ent.fieldType != MULTIPLYENTRY {
			t.Errorf("Entry gets field type 0x%x, want MULTIPLY=0x%x", ent.fieldType, MULTIPLYENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		expectedf := []string{"data", "sbit"}
		for i := 0; i < 2; i++ {
			if ent.inFields[i] != expectedf[i] {
				t.Errorf("Entry gets inFields[%d]=\"%s\", want \"%s\"", i, ent.inFields[i], expectedf[i])
			}
		}
	}

	// #48a: Entry (div) check
	ent, err = d.Entry("div")
	if err != nil {
		t.Error("Could not get Entry for div type: ", err)
	} else {
		if ent.fieldType != DIVIDEENTRY {
			t.Errorf("Entry gets field type 0x%x, want DIVIDE=0x%x", ent.fieldType, DIVIDEENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		expectedf := []string{"mult", "bit"}
		for i := 0; i < 2; i++ {
			if ent.inFields[i] != expectedf[i] {
				t.Errorf("Entry gets inFields[%d]=\"%s\", want \"%s\"", i, ent.inFields[i], expectedf[i])
			}
		}
	}

	// #48b: Entry (recip) check
	ent, err = d.Entry("recip")
	if err != nil {
		t.Error("Could not get Entry for recip type: ", err)
	} else {
		if ent.fieldType != RECIPENTRY {
			t.Errorf("Entry gets field type 0x%x, want RECIP=0x%x", ent.fieldType, RECIPENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		expectedf := []string{"div"}
		for i := 0; i < 1; i++ {
			if ent.inFields[i] != expectedf[i] {
				t.Errorf("Entry gets inFields[%d]=\"%s\", want \"%s\"", i, ent.inFields[i], expectedf[i])
			}
		}
		expected := complex(6.5, 4.3)
		if ent.dividend != real(expected) {
			t.Errorf("Entry recip gets dividend=%f, want %f", ent.dividend, real(expected))
		}
		if ent.cdividend != expected {
			t.Errorf("Entry recip gets cdividend=%f, want %f", ent.cdividend, expected)
		}
	}

	// #49: Entry (phase) check
	ent, err = d.Entry("phase")
	if err != nil {
		t.Error("Could not get Entry for phase type: ", err)
	} else {
		if ent.fieldType != PHASEENTRY {
			t.Errorf("Entry gets field type 0x%x, want PHASE=0x%x", ent.fieldType, PHASEENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		expectedf := []string{"data"}
		for i := 0; i < 1; i++ {
			if ent.inFields[i] != expectedf[i] {
				t.Errorf("Entry gets inFields[%d]=\"%s\", want \"%s\"", i, ent.inFields[i], expectedf[i])
			}
		}
		if ent.phaseShift != 11 {
			t.Errorf("Entry recip gets phaseShift=%d, want 11", ent.phaseShift)
		}
	}

	// #50: Entry (const) check
	ent, err = d.Entry("const")
	if err != nil {
		t.Error("Could not get Entry for const type: ", err)
	} else {
		if ent.fieldType != CONSTENTRY {
			t.Errorf("Entry gets field type 0x%x, want CONST=0x%x", ent.fieldType, CONSTENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
		if ent.constType != FLOAT64 {
			t.Errorf("Entry recip gets constType=0x%x, want FLOAT64=0x%x", ent.constType, FLOAT64)
		}
	}

	// #51: Entry (string) check
	ent, err = d.Entry("string")
	if err != nil {
		t.Error("Could not get Entry for string type: ", err)
	} else {
		if ent.fieldType != STRINGENTRY {
			t.Errorf("Entry gets field type 0x%x, want STRING=0x%x", ent.fieldType, STRINGENTRY)
		}
		if ent.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", ent.fragment)
		}
	}

	// #52: FragmentIndex check
	n52, err := d.FragmentIndex("data")
	if err != nil {
		t.Error("Could not get FragmentIndex: ", err.Error())
	} else if n52 != 0 {
		t.Errorf("FragmentIndex returns %d, want 0", n52)
	}

	// #53: RawEntry check
	e53 := RawEntry("new1", 0, 3, FLOAT64)
	err = d.AddEntry(&e53)
	if err != nil {
		t.Error("Could not AddEntry in test 53:", err)
	}
	e53b, err := d.Entry("new1")
	if err != nil {
		t.Error("Could not read Entry in test 53:", err)
	} else {
		if e53b.fieldType != RAWENTRY {
			t.Errorf("Entry gets field type 0x%x, want RAW=0x%x", e53b.fieldType, RAWENTRY)
		}
		if e53b.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", e53b.fragment)
		}
		if e53b.dataType != FLOAT64 {
			t.Errorf("Entry gets data type 0x%x, want FLOAT64=0x%x", e53b.dataType, FLOAT64)
		}
		if e53b.spf != 3 {
			t.Errorf("Entry gets spf=%d, want 3", e53b.spf)
		}
	}

	// #54: Lincom check
	in54 := []string{"in1", "in2"}
	m54 := []float64{9.9, 7.7}
	b54 := []float64{8.8, 6.6}
	err = d.AddLincom("new2", in54, m54, b54, 0)
	if err != nil {
		t.Error("Could not AddLincom in test 54:", err)
	}
	e54, err := d.Entry("new2")
	if err != nil {
		t.Error("Could not read Entry new2 in test 54:", err)
	} else {
		if e54.fieldType != LINCOMENTRY {
			t.Errorf("Entry new2 gets field type 0x%x, want LINCOM=0x%x", e54.fieldType, LINCOMENTRY)
		}
		if e54.fragment != 0 {
			t.Errorf("Entry new2 gets fragment index=%d, want 0", e54.fragment)
		}
		if e54.nFields != 2 {
			t.Errorf("Entry new2 gets %d fields, want 2", e54.nFields)
		}
		for i := 0; i < len(in54); i++ {
			if e54.inFields[i] != in54[i] {
				t.Errorf("Entry new2 inFields[%d]=%s, want %s", i, e54.inFields[i], in54[i])
			}
			if e54.m[i] != m54[i] {
				t.Errorf("Entry new2 m[%d]=%f, want %f", i, e54.m[i], m54[i])
			}
			if e54.b[i] != b54[i] {
				t.Errorf("Entry new2 b[%d]=%f, want %f", i, e54.b[i], b54[i])
			}
		}
	}

	// #55: CLincom (complex) check
	in55 := []string{"in1", "in2"}
	m55 := []complex128{complex(1.1, 1.2), complex(1.4, 1.5)}
	b55 := []complex128{complex(1.3, 1.4), complex(1.6, 1.7)}

	err = d.AddCLincom("new3", in55, m55, b55, 0)
	if err != nil {
		t.Error("Could not AddCLincom in test 55:", err)
	}
	e55, err := d.Entry("new3")
	if err != nil {
		t.Error("Could not read Entry new3 in test 55:", err)
	} else {
		if e55.fieldType != LINCOMENTRY {
			t.Errorf("Entry new3 gets field type 0x%x, want LINCOM=0x%x", e55.fieldType, LINCOMENTRY)
		}
		if e55.fragment != 0 {
			t.Errorf("Entry new3 gets fragment index=%d, want 0", e55.fragment)
		}
		if e55.nFields != 2 {
			t.Errorf("Entry new3 gets %d fields, want 2", e55.nFields)
		}
		for i := 0; i < len(in55); i++ {
			if e55.inFields[i] != in55[i] {
				t.Errorf("Entry new3 inFields[%d]=%s, want %s", i, e55.inFields[i], in55[i])
			}
			if e55.cm[i] != m55[i] {
				t.Errorf("Entry new3 m[%d]=%f, want %f", i, e55.cm[i], m55[i])
			}
			if e55.cb[i] != b55[i] {
				t.Errorf("Entry new3 b[%d]=%f, want %f", i, e55.cb[i], b55[i])
			}
		}
	}

	// #56: Polynom check
	in56 := "in1"
	a56 := []float64{3.9, 4.8, 5.7, 6.6}
	err = d.AddPolynom("new4", in56, a56, 0)
	if err != nil {
		t.Error("Could not AddPolynom in test 56:", err)
	}
	e56, err := d.Entry("new4")
	if err != nil {
		t.Error("Could not read Entry new4 in test 56:", err)
	} else {
		if e56.fieldType != POLYNOMENTRY {
			t.Errorf("Entry new4 gets field type 0x%x, want POLYNOM=0x%x", e56.fieldType, POLYNOMENTRY)
		}
		if e56.fragment != 0 {
			t.Errorf("Entry new4 gets fragment index=%d, want 0", e56.fragment)
		}
		if e56.polyOrder != len(a56)-1 {
			t.Errorf("Entry new4 has polyOrder=%d, want %d", e56.polyOrder, len(a56)-1)
		}
		if e56.inFields[0] != in56 {
			t.Errorf("Entry new4 inFields[0]=%s, want %s", e56.inFields[0], in56)
		}
		for i := 0; i < len(a56); i++ {
			if e56.a[i] != a56[i] {
				t.Errorf("Entry new4 a[%d]=%f, want %f", i, e56.a[i], a56[i])
			}
		}
	}

	// #57: CPolynom check
	in57 := "in3"
	a57 := []complex128{complex(3.1, 9), complex(4.2, 8), complex(5.2, 9), complex(6.3, 4.4)}
	err = d.AddCPolynom("new5", in57, a57, 0)
	if err != nil {
		t.Error("Could not AddCPolynom in test 57:", err)
	}
	e57, err := d.Entry("new5")
	if err != nil {
		t.Error("Could not read Entry new5 in test 57:", err)
	} else {
		if e57.fieldType != POLYNOMENTRY {
			t.Errorf("Entry new5 gets field type 0x%x, want POLYNOM=0x%x", e57.fieldType, POLYNOMENTRY)
		}
		if e57.fragment != 0 {
			t.Errorf("Entry new5 gets fragment index=%d, want 0", e57.fragment)
		}
		if e57.polyOrder != len(a57)-1 {
			t.Errorf("Entry new5 has polyOrder=%d, want %d", e57.polyOrder, len(a57)-1)
		}
		if e57.inFields[0] != in57 {
			t.Errorf("Entry new5 inFields[0]=%s, want %s", e57.inFields[0], in57)
		}
		for i := 0; i < len(a57); i++ {
			if e57.ca[i] != a57[i] {
				t.Errorf("Entry new5 ca[%d]=%f, want %f", i, e57.ca[i], a57[i])
			}
		}
	}

	// #58: Linterp check
	in58 := "in1"
	t58 := "./some/table"
	err = d.AddLinterp("new6", in58, t58, 0)
	if err != nil {
		t.Error("Could not AddLinterp in test 58:", err)
	}
	e58, err := d.Entry("new6")
	if err != nil {
		t.Error("Could not read Entry new6 in test 58:", err)
	} else {
		if e58.fieldType != LINTERPENTRY {
			t.Errorf("Entry new6 gets field type 0x%x, want LINTERP=0x%x", e58.fieldType, LINTERPENTRY)
		}
		if e58.fragment != 0 {
			t.Errorf("Entry new6 gets fragment index=%d, want 0", e58.fragment)
		}
		if e58.inFields[0] != in58 {
			t.Errorf("Entry new6 inFields[0]=%s, want %s", e58.inFields[0], in58)
		}
		if e58.table != t58 {
			t.Errorf("Entry new6 table=%s, want %s", e58.table, t58)
		}
	}

	// #59: BitEntry check
	e59 := BitEntry("new7", "in", 13, 12, 0)
	err = d.AddEntry(&e59)
	if err != nil {
		t.Error("Could not AddEntry in test 59:", err)
	}
	e59b, err := d.Entry("new7")
	if err != nil {
		t.Error("Could not read Entry in test 59:", err)
	} else {
		if e59b.fieldType != BITENTRY {
			t.Errorf("Entry gets field type 0x%x, want BIT=0x%x", e59b.fieldType, BITENTRY)
		}
		if e59b.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", e59b.fragment)
		}
		if e59b.inFields[0] != "in" {
			t.Errorf("Entry in_fields[0]=%s, want %s", e59b.inFields[0], "in")
		}
		if e59b.bitnum != 13 {
			t.Errorf("Entry gets bitnum=%d, want 13", e59b.bitnum)
		}
		if e59b.numbits != 12 {
			t.Errorf("Entry gets numbits=%d, want 12", e59b.numbits)
		}
	}

	// #60: SBit Entry check
	err = d.AddSbit("new8", "in2", 14, 15, 0)
	if err != nil {
		t.Error("Could not AddEntry in test 60:", err)
	}
	e60, err := d.Entry("new8")
	if err != nil {
		t.Error("Could not read Entry in test 60:", err)
	} else {
		if e60.fieldType != SBITENTRY {
			t.Errorf("Entry gets field type 0x%x, want SBIT=0x%x", e60.fieldType, SBITENTRY)
		}
		if e60.fragment != 0 {
			t.Errorf("Entry gets fragment index=%d, want 0", e60.fragment)
		}
		if e60.inFields[0] != "in2" {
			t.Errorf("Entry in_fields[0]=%s, want %s", e60.inFields[0], "in")
		}
		if e60.bitnum != 14 {
			t.Errorf("Entry gets bitnum=%d, want 14", e60.bitnum)
		}
		if e60.numbits != 15 {
			t.Errorf("Entry gets numbits=%d, want 15", e60.numbits)
		}
	}

	// #61: Mult check
	in61 := []string{"in1", "in2"}
	err = d.AddMultiply("new9", in61[0], in61[1], 0)
	if err != nil {
		t.Error("Could not AddMultiply in test 61:", err)
	}
	e61, err := d.Entry("new9")
	if err != nil {
		t.Error("Could not read Entry new9 in test 61:", err)
	} else {
		if e61.fieldType != MULTIPLYENTRY {
			t.Errorf("Entry new9 gets field type 0x%x, want MULTIPLY=0x%x", e61.fieldType, MULTIPLYENTRY)
		}
		if e61.fragment != 0 {
			t.Errorf("Entry new9 gets fragment index=%d, want 0", e61.fragment)
		}
		for i := 0; i < len(in61); i++ {
			if e61.inFields[i] != in61[i] {
				t.Errorf("Entry new9 inFields[%d]=%s, want %s", i, e61.inFields[i], in61[i])
			}
		}
	}

	// #62: Phase check
	in62 := []string{"in1"}
	shift62 := int64(22)
	err = d.AddPhase("new10", in62[0], shift62, 0)
	if err != nil {
		t.Error("Could not AddPhase in test 62:", err)
	}
	e62, err := d.Entry("new10")
	if err != nil {
		t.Error("Could not read Entry new10 in test 62:", err)
	} else {
		if e62.fieldType != PHASEENTRY {
			t.Errorf("Entry new10 gets field type 0x%x, want PHASE=0x%x", e62.fieldType, PHASEENTRY)
		}
		if e62.fragment != 0 {
			t.Errorf("Entry new10 gets fragment index=%d, want 0", e62.fragment)
		}
		for i := 0; i < len(in62); i++ {
			if e62.inFields[i] != in62[i] {
				t.Errorf("Entry new10 inFields[%d]=%s, want %s", i, e62.inFields[i], in62[i])
			}
		}
		if e62.phaseShift != shift62 {
			t.Errorf("Entry new10 gets phaseShift=%d, want %d", e62.fieldType, shift62)
		}
	}

	// #63: Const check
	v63 := float32(5.6)
	err = d.AddConst("new11", FLOAT64, v63, 0)
	if err != nil {
		t.Error("Could not AddConst in test 63:", err)
	}
	e63, err := d.Entry("new11")
	if err != nil {
		t.Error("Could not read Entry new11 in test 63:", err)
	} else {
		if e63.fieldType != CONSTENTRY {
			t.Errorf("Entry new11 gets field type 0x%x, want CONST=0x%x", e63.fieldType, CONSTENTRY)
		}
		if e63.fragment != 0 {
			t.Errorf("Entry new11 gets fragment index=%d, want 0", e63.fragment)
		}
		if e63.constType != FLOAT64 {
			t.Errorf("Entry new11 gets const type 0x%x, want FLOAT64=0x%x", e63.constType, FLOAT64)
		}
	}

	// #65: nfragments check
	nfrag := d.NFragments()
	if nfrag != 1 {
		t.Errorf("NFragments returned %d, want 1", nfrag)
	}

	// #66: include check
	idx, err := d.Include("form2", 0)
	if err != nil {
		t.Errorf("Could not Include(\"form2\")")
	} else if idx != 1 {
		t.Errorf("Include(\"form2\") returned %d, want 1", idx)
	}
	c2, _ := d.GetConstantInt32("const2")
	if c2 != -19 {
		t.Errorf("Failed to read form2 fragment const2: get %d, want -19", c2)
	}

	// #67: NFieldsByType check
	nlincom := d.NFieldsByType(LINCOMENTRY)
	if nlincom != 3 {
		t.Errorf("NFieldsByType(LINCOMENTRY) returned %d, want 1", nlincom)
	}

	// #68: FieldListByType check
	ftnames := d.FieldListByType(LINCOMENTRY)
	trueftnames := []string{"new2", "new3", "lincom"}
	if len(ftnames) != len(trueftnames) {
		t.Errorf("FieldListByType length is %d, want %d", len(ftnames), len(trueftnames))
	} else {
		for i := 0; i < len(trueftnames); i++ {
			if ftnames[i] != trueftnames[i] {
				t.Errorf("FieldListByType[%d]=\"%s\", want \"%s\"", i, ftnames[i], trueftnames[i])
			}
		}
	}

	// #69: NVectors check
	nvec := d.NVectors()
	if nvec != 25 {
		t.Errorf("NVectors = %d, want 25", nvec)
	}

	// #70: VectorList check
	vectors := d.VectorList()
	if len(vectors) != int(nvec) {
		t.Errorf("VectorList length is %d, want %d", len(vectors), nvec)
	} else {
		truevnames := []string{"bit", "div", "data", "mult", "new1", "new2", "new3",
			"new4", "new5", "new6", "new7", "new8", "new9", "sbit", "INDEX",
			"alias", "indir", "mplex", "new10", "phase", "recip", "lincom",
			"window", "linterp", "polynom"}
		for i := 0; i < int(nvec); i++ {
			if vectors[i] != truevnames[i] {
				t.Errorf("FieldList[%d]=\"%s\", want \"%s\"", i, vectors[i], truevnames[i])
			}
		}
	}

	// #81: GetString check
	st, err := d.GetString("string")
	if err != nil {
		t.Error("GetString error in test 81: ", err)
	} else {
		if st != "Zaphod Beeblebrox" {
			t.Errorf("GetString returned \"%s\", want \"Zaphod Beeblebrox\"", st)
		}
	}

	// #82: AddString check
	err = d.AddString("new12", "glob", 0)
	if err != nil {
		t.Error("AddString error in test 82: ", err)
	}
	e82, err := d.Entry("new12")
	if err != nil {
		t.Error("Entry error in test 82: ", err)
	} else {
		if e82.fieldType != STRINGENTRY {
			t.Errorf("Entry new12 gets field type 0x%x, want STRING=0x%x", e82.fieldType, STRINGENTRY)
		}
		if e82.fragment != 0 {
			t.Errorf("Entry new21 gets fragment index=%d, want 0", e82.fragment)
		}
	}

	// #86: PutConstant int32 check
	err = d.PutConstant("const", int32(86))
	if err != nil {
		t.Error("PutConstant error in test 86: ", err)
	} else {
		var v int32
		err = d.GetConstant("const", &v)
		if err != nil {
			t.Error("GetConstant error in test 86: ", err)
		} else if v != 86 {
			t.Errorf("PutConstant(int32(86)) then GetConstant reads %d, want 86", v)
		}
	}

	// #88: PutConstant uint64 check
	err = d.PutConstant("const", uint64(88))
	if err != nil {
		t.Error("PutConstant error in test 88: ", err)
	} else {
		var v uint64
		err = d.GetConstant("const", &v)
		if err != nil {
			t.Error("GetConstant error in test 88: ", err)
		} else if v != 88 {
			t.Errorf("PutConstant(uint64(88)) then GetConstant reads %d, want 88", v)
		}
	}

	// #89: PutConstant int64 check
	err = d.PutConstant("const", int64(89))
	if err != nil {
		t.Error("PutConstant error in test 89: ", err)
	} else {
		var v int64
		err = d.GetConstant("const", &v)
		if err != nil {
			t.Error("GetConstant error in test 89: ", err)
		} else if v != 89 {
			t.Errorf("PutConstant(int64(89)) then GetConstant reads %d, want 89", v)
		}
	}

	// #91: PutConstant float64 check
	err = d.PutConstant("const", float64(91))
	if err != nil {
		t.Error("PutConstant error in test 91: ", err)
	} else {
		var v float64
		err = d.GetConstant("const", &v)
		if err != nil {
			t.Error("GetConstant error in test 91: ", err)
		} else if v != 91 {
			t.Errorf("PutConstant(float64(91)) then GetConstant reads %.2f, want 91.00", v)
		}
	}

	// #93: PutConstant complex64 check
	err = d.PutConstant("const", complex64(93))
	if err != nil {
		t.Error("PutConstant error in test 93: ", err)
	} else {
		var v complex64
		err = d.GetConstant("const", &v)
		if err != nil {
			t.Error("GetConstant error in test 93: ", err)
		} else if v != 93 {
			t.Errorf("PutConstant(complex64(93)) then GetConstant reads %.2f, want 93.00", v)
		}
	}

	// #94: PutString check
	err = d.PutString("string", "Arthur Dent")
	if err != nil {
		t.Error("PutString error in test 94: ", err)
	} else {
		s, err2 := d.GetString("string")
		if err2 != nil {
			t.Error("GetString error in test 94: ", err2)
		} else if s != "Arthur Dent" {
			t.Errorf("PutString() then GetString reads %s, want \"Arthur Dent\"", s)
		}
	}

	// #95: NMFieldsByType check
	nlinterp := d.NMFieldsByType("data", LINTERPENTRY)
	if nlinterp != 1 {
		t.Errorf("NMFieldsByType(\"data\", LINTERPENTRY) returned %d, want 0", nlinterp)
	}

	// #96: MFieldListByType check
	mtfields := d.MFieldListByType("data", LINTERPENTRY)
	if len(mtfields) != int(nlinterp) {
		t.Errorf("MVectorList(\"data\", LINTERPENTRY) length is %d, want %d", len(mtfields), nlinterp)
	} else {
		truemtnames := []string{"mlut"}
		for i := 0; i < int(nlinterp); i++ {
			if mtfields[i] != truemtnames[i] {
				t.Errorf("MFieldListByType[%d]=\"%s\", want \"%s\"", i, mtfields[i], truemtnames[i])
			}
		}
	}

	// #97: NMVectors check
	mnvec := d.NMVectors("data")
	if mnvec != 1 {
		t.Errorf("NMVectors(\"data\") returned %d, want 1", mnvec)
	}

	// #98: MVectorList check
	mvectors := d.MVectorList("data")
	if len(mvectors) != int(mnvec) {
		t.Errorf("MVectorList length is %d, want %d", len(mvectors), mnvec)
	} else {
		truemvnames := []string{"mlut"}
		for i := 0; i < int(mnvec); i++ {
			if mvectors[i] != truemvnames[i] {
				t.Errorf("MVectorList[%d]=\"%s\", want \"%s\"", i, vectors[i], truemvnames[i])
			}
		}
	}

	// #110: Fragment encoding check
	frag, err := d.Fragment(0)
	if err != nil {
		t.Errorf("Could not create a Dirfile.Fragment(0)")
	}
	if frag.encoding != UNENCODED {
		t.Errorf("frag.encoding is %d, want %d", frag.encoding, UNENCODED)
	}

	// #111: Fragment endianness check
	if frag.endianness != LITTLEENDIAN {
		t.Errorf("frag.endianness is 0x%x, want 0x%x", frag.endianness, LITTLEENDIAN)
	}

	// #112: dirfilename check
	name := d.Dirfilename()
	if !strings.HasSuffix(name, d.name) {
		t.Errorf("d.Dirfilename() returns '%s', want suffix to be '%s'", name, d.name)
	}

	// #113: Fragment parent check
	if frag.parent != -1 {
		t.Errorf("frag.parent is %d, want -1", frag.parent)
	}
	frag1, err := d.Fragment(1)
	if err != nil {
		t.Errorf("Could not run a Dirfile.Fragment(1)")
	} else {
		if frag1.parent != 0 {
			t.Errorf("frag1.parent is %d, want 0", frag1.parent)
		}
		if !strings.HasSuffix(frag1.name, "form2") {
			t.Errorf("frag1.name is %s, want suffix to be \"form2\"", frag1.name)
		}

		// #114: Fragment.SetProtection check
		err = frag1.SetProtection(PROTECTDATA)
		if err != nil {
			t.Errorf("Could not SetProtection(PROTECTDATA): %s", err)
		}
	}

	// #115: Fragment.protection check
	frag1, err = d.Fragment(1)
	if err != nil {
		t.Errorf("Could not run a Dirfile.Fragment(1)")
	} else if frag1.protection != PROTECTDATA {
		t.Errorf("frag1.protection is 0x%x, want 0x%x", frag1.protection, PROTECTDATA)
	}

	// #122: Dirfile.Delete check
	err = d.Delete("new10", 0)
	if err != nil {
		t.Error("Could not Delete in test 122:", err)
	}
	_, err = d.Entry("new10")
	if err == nil {
		t.Error("Loaded entry new10 after it was deleted in test 122:", err)
	}

	// #138: PutData (auto) check duplicates #37, so skip it.

	// #146: Divide check
	in146 := []string{"in1", "in2"}
	err = d.AddDivide("new14", in146[0], in146[1], 0)
	if err != nil {
		t.Error("Could not AddDivide in test 146:", err)
	}
	e146, err := d.Entry("new14")
	if err != nil {
		t.Error("Could not read Entry new14 in test 146:", err)
	} else {
		if e146.fieldType != DIVIDEENTRY {
			t.Errorf("Entry new14 gets field type 0x%x, want DIVIDE=0x%x", e146.fieldType, DIVIDEENTRY)
		}
		if e146.fragment != 0 {
			t.Errorf("Entry new14 gets fragment index=%d, want 0", e146.fragment)
		}
		for i := 0; i < len(in146); i++ {
			if e146.inFields[i] != in146[i] {
				t.Errorf("Entry new14 inFields[%d]=%s, want %s", i, e146.inFields[i], in146[i])
			}
		}
	}

	// #147: Recip check
	in147 := []string{"in3"}
	err = d.AddRecip("new15", in147[0], 31.9, 0)
	if err != nil {
		t.Error("Could not AddRecip in test 147:", err)
	}
	e147, err := d.Entry("new15")
	if err != nil {
		t.Error("Could not read Entry new15 in test 147:", err)
	} else {
		if e147.fieldType != RECIPENTRY {
			t.Errorf("Entry new15 gets field type 0x%x, want RECIP=0x%x", e147.fieldType, RECIPENTRY)
		}
		if e147.fragment != 0 {
			t.Errorf("Entry new15 gets fragment index=%d, want 0", e147.fragment)
		}
		for i := 0; i < len(in147); i++ {
			if e147.inFields[i] != in147[i] {
				t.Errorf("Entry new15 inFields[%d]=%s, want %s", i, e147.inFields[i], in147[i])
			}
		}
	}

	// #148: CRecip check
	in148 := []string{"in2"}
	d148 := complex(33.3, 44.4)
	err = d.AddCRecip("new16", in148[0], d148, 0)
	if err != nil {
		t.Error("Could not AddCRecip in test 148:", err)
	}
	e148, err := d.Entry("new16")
	if err != nil {
		t.Error("Could not read Entry new16 in test 148:", err)
	} else {
		if e148.fieldType != RECIPENTRY {
			t.Errorf("Entry new16 gets field type 0x%x, want RECIP=0x%x", e148.fieldType, RECIPENTRY)
		}
		if e148.fragment != 0 {
			t.Errorf("Entry new16 gets fragment index=%d, want 0", e148.fragment)
		}
		for i := 0; i < len(in148); i++ {
			if e148.inFields[i] != in148[i] {
				t.Errorf("Entry new16 inFields[%d]=%s, want %s", i, e148.inFields[i], in148[i])
			}
		}
	}

	// #155: Fragment.Rewrite check
	err = frag.Rewrite()
	if err != nil {
		t.Errorf("Could not run Fragment.Rewrite()")
	}

	// #156: invalid dirfile check
	invalid := InvalidDirfile()
	if invalid.d == nil {
		t.Errorf("InvalidDirfile returned a nil dirfile")
	}
	err = invalid.Flush("data")
	if err == nil {
		t.Errorf("InvalidDirfile().Flush() did not return error")
	}
	err = invalid.FlushAll()
	if err == nil {
		t.Errorf("InvalidDirfile().FlushAll() did not return error")
	}
	err = invalid.Sync("data")
	if err == nil {
		t.Errorf("InvalidDirfile().Sync() did not return error")
	}
	err = invalid.SyncAll()
	if err == nil {
		t.Errorf("InvalidDirfile().SyncAll() did not return error")
	}
	err = invalid.RawClose("data")
	if err == nil {
		t.Errorf("InvalidDirfile().RawClose() did not return error")
	}
	err = invalid.RawCloseAll()
	if err == nil {
		t.Errorf("InvalidDirfile().RawCloseAll() did not return error")
	}
	err = invalid.MetaFlush()
	if err == nil {
		t.Errorf("InvalidDirfile().MetaFlush() did not return error")
	}
	err = invalid.Close()
	if err != nil {
		t.Errorf("Could not Close an invalid dirfile")
	}

	// #158: GetCarray test
	a158 := make([]float64, 8)
	err = d.GetCarray("carray", &a158)
	if err != nil {
		t.Error("Could not GetCarray: ", err)
	} else {
		L := d.ArrayLen("carray")
		for i := 0; i < L; i++ {
			if math.Abs(a158[i]-1.1*float64(i+1)) > 1e-6 {
				t.Errorf("GetCarray returns v[%d]=%.2f, want %.2f", i, a158[i], 1.1*float64(i+1))
			}
		}
	}
	err = d.GetCarray("sdfsdfasf", &a158)
	if err == nil {
		t.Error("GetCarray did not error when passed an invalid field name")
	}
	err = d.GetCarray("carray", &invalid)
	if err == nil {
		t.Error("GetCarray did not error when passed an invalid pointer")
	}

	// #159: GetCarray test
	err = d.GetCarraySlice("carray", 2, 2, &a158)
	if err != nil {
		t.Error("Could not GetCarraySlice: ", err)
	} else {
		for i := 0; i < 2; i++ {
			if math.Abs(a158[i]-1.1*float64(i+3)) > 1e-6 {
				t.Errorf("GetCarraySlice returns v[%d]=%.2f, want %.2f", i, a158[i], 1.1*float64(i+1))
			}
		}
	}
	err = d.GetCarraySlice("sdfsdfasf", 2, 2, &a158)
	if err == nil {
		t.Error("GetCarraySlice did not error when passed an invalid field name")
	}
	err = d.GetCarraySlice("carray", 2, 2, &invalid)
	if err == nil {
		t.Error("GetCarraySlice did not error when passed an invalid pointer")
	}

	// #168: PutCarray test
	p168 := []float64{9.6, 8.5, 7.4, 6.3, 5.2, 4.1}
	err = d.PutCarray("carray", p168)
	if err != nil {
		t.Error("PutCarray failed: ", err)
	}
	err = d.GetCarray("carray", &a158)
	if err != nil {
		t.Error("GetCarray failed on test 168: ", err)
	} else {
		for i := 0; i < 6; i++ {
			if math.Abs(a158[i]-9.6+1.1*float64(i)) > 1e-6 {
				t.Errorf("GetCarray returns v[%d]=%.2f, want %.2f", i, a158[i], 9.6+1.1*float64(i))
			}
		}
	}
	err = d.PutCarray("sdfsdfasf", p168)
	if err == nil {
		t.Error("PutCarray did not error when passed an invalid field name")
	}

	// #177: ArrayLen test
	l177 := d.ArrayLen("carray")
	if l177 != 6 {
		t.Errorf("ArrayLen(\"carray\") returned %d, want 6", l177)
	}

	// #199: Strings test
	// s199, err := d.Strings()
	// if err != nil {
	// 	t.Error("Strings() failed: ", err)
	// } else {
	// 	// expected := []string{"Lorem ipsum", "", "Arthur Dent"}
	// 	// TODO: use the above
	// 	expected := []string{"Arthur Dent"}
	// 	for i := 0; i < len(expected); i++ {
	// 		if s199[i] != expected[i] {
	// 			t.Errorf("Strings returned s[%d]=%s, want %s", i, s199[i], expected[i])
	// 		}
	// 	}
	// }

	// #208: Sync check
	err = d.Sync("data")
	if err != nil {
		t.Errorf("Could not call Sync on a field")
	}
	err = d.Sync("")
	if err != nil {
		t.Errorf("Could not call Sync on all fields")
	}
	err = d.SyncAll()
	if err != nil {
		t.Errorf("Could not call SyncAll")
	}

	// #209 flush check
	err = d.Flush("data")
	if err != nil {
		t.Errorf("Could not call Flush on a field")
	}
	err = d.Flush("")
	if err != nil {
		t.Errorf("Could not call Flush on all fields")
	}
	err = d.FlushAll()
	if err != nil {
		t.Errorf("Could not call FlushAll")
	}

	// #210 metaflush check
	err = d.MetaFlush()
	if err != nil {
		t.Errorf("Could not call MetaFlush")
	}

	// #212: AddWindow check
	err = d.AddWindow("new18", "in1", "in2", WINDOPNE, 32, 0)
	if err != nil {
		t.Error("AddWindow error in test 212: ", err)
	}
	e212, err := d.Entry("new18")
	if err != nil {
		t.Error("Entry error in test 212: ", err)
	} else {
		if e212.fieldType != WINDOWENTRY {
			t.Errorf("Entry new20 (alias) gets field type 0x%x, want WINDOW=0x%x", e212.fieldType, WINDOWENTRY)
		}
		if e212.fragment != 0 {
			t.Errorf("Entry new20 gets fragment index=%d, want 0", e212.fragment)
		}
	}

	// #219: AddAlias check
	err = d.AddAlias("new20", "data", 0)
	if err != nil {
		t.Error("AddAlias error in test 219: ", err)
	}
	e219, err := d.Entry("new20")
	if err != nil {
		t.Error("Entry error in test 219: ", err)
	} else {
		if e219.fieldType != RAWENTRY {
			t.Errorf("Entry new20 (alias) gets field type 0x%x, want RAW=0x%x", e219.fieldType, RAWENTRY)
		}
		if e219.fragment != 0 {
			t.Errorf("Entry new20 gets fragment index=%d, want 0", e219.fragment)
		}
	}

	// #229: AddMplex check
	in229 := []string{"in1", "in2"}
	err = d.AddMplex("new21", in229[0], in229[1], 5, 6, 0)
	if err != nil {
		t.Error("Could not AddMplex in test 229:", err)
	}
	e229, err := d.Entry("new21")
	if err != nil {
		t.Error("Could not read Entry new21 in test 229:", err)
	} else {
		if e229.fieldType != MPLEXENTRY {
			t.Errorf("Entry new21 gets field type 0x%x, want MPLEX=0x%x", e229.fieldType, MPLEXENTRY)
		}
		if e229.fragment != 0 {
			t.Errorf("Entry new21 gets fragment index=%d, want 0", e229.fragment)
		}
		for i := 0; i < len(in229); i++ {
			if e229.inFields[i] != in229[i] {
				t.Errorf("Entry new21 inFields[%d]=%s, want %s", i, e229.inFields[i], in229[i])
			}
		}
		if e229.countVal != 5 {
			t.Errorf("Entry new21 gets countVal=%d, want 5", e229.countVal)
		}
		if e229.period != 6 {
			t.Errorf("Entry new21 gets period=%d, want 6", e229.period)
		}
	}

	// #233: raw_close check
	err = d.RawClose("data")
	if err != nil {
		t.Errorf("Could not call RawClose on a field")
	}
	err = d.RawClose("")
	if err != nil {
		t.Errorf("Could not call RawClose on all fields")
	}
	err = d.RawCloseAll()
	if err != nil {
		t.Errorf("Could not call RawCloseAll")
	}

	// #234: desync check
	_, err = d.Desync(true, true)
	if err != nil {
		t.Errorf("Could not call Desync(true, true)")
	}

	// #235: Flags check
	d.Flags(VERBOSE, 0)
	flags := d.Flags(PRETTYPRINT, 0)
	if flags&PRETTYPRINT == 0 {
		t.Errorf("Flags(0x%x, 0x0) returned 0x%x, want that flag set", PRETTYPRINT, flags)
	}
	flags = d.Flags(0, PRETTYPRINT)
	if flags&PRETTYPRINT != 0 {
		t.Errorf("Flags(0x0, 0x%x) returned 0x%x, want that flag clear", PRETTYPRINT, flags)
	}

	// #236: VerbosePrefix check
	err = d.VerbosePrefix("big_test: ")
	if err != nil {
		t.Errorf("Could not set VerbosePrefix()")
	}

	// #237: NEntries check
	ne := d.NEntries("data", SCALARENTRIES, HIDDENENTRIES|NOALIASENTRIES)
	if ne != 4 { // TODO: eventually 5
		t.Errorf("d.NEntries counts %d SCALAR entries, want 4", ne)
	}
	ne = d.NEntries("", VECTORENTRIES, HIDDENENTRIES|NOALIASENTRIES)
	if ne != 28 {
		t.Errorf("d.NEntries counts %d VECTOR entries, want 28", ne)
	}

	// #239: EntryList check
	entryList := d.EntryList("", VECTORENTRIES, HIDDENENTRIES|NOALIASENTRIES)
	if len(entryList) != int(ne) {
		t.Errorf("d.EntryList return %d entries, want %d", len(entryList), ne)
	} else {
		trueEntries := []string{"bit", "div", "data", "mult", "new1", "new2", "new3",
			"new4", "new5", "new6", "new7", "new8", "new9", "sbit", "INDEX",
			"indir", "mplex", "new14", "new15", "new16", "new18", "new21", "phase", "recip", "lincom",
			"window", "linterp", "polynom"}
		for i := 0; i < int(ne); i++ {
			if entryList[i] != trueEntries[i] {
				t.Errorf("EntryList[%d]=\"%s\", want \"%s\"", i, entryList[i], trueEntries[i])
			}
		}
	}

	// #240: MplexLookback test (returns nothing, so simply call it)
	d.MplexLookback(LOOKBACKALL)

	// #283: Sarray check
	in283 := []string{"str1", "str2", "str4", "str8"}
	err = d.AddSarray("new283", in283, 0)
	if err != nil {
		t.Error("Could not AddSarray in test 283:", err)
	}
	e283, err := d.Entry("new283")
	if err != nil {
		t.Error("Could not read Entry new283 in test 283:", err)
	} else {
		if e283.fieldType != SARRAYENTRY {
			t.Errorf("Entry new283 gets field type 0x%x, want SARRAY=0x%x", e283.fieldType, SARRAYENTRY)
		}
		if e283.fragment != 0 {
			t.Errorf("Entry new283 gets fragment index=%d, want 0", e283.fragment)
		}
		if e283.arrayLen != len(in283) {
			t.Errorf("Entry new283 gets %d fields, want %d", e283.nFields, len(in283))
		}
	}

	// #289: Indir check
	in289 := []string{"in1", "in2"}
	err = d.AddIndir("new289", in289[0], in289[1], 0)
	if err != nil {
		t.Error("Could not AddIndir in test 289:", err)
	}
	e289, err := d.Entry("new289")
	if err != nil {
		t.Error("Could not read Entry new289 in test 289:", err)
	} else {
		if e289.fieldType != INDIRENTRY {
			t.Errorf("Entry new289 gets field type 0x%x, want INDIR=0x%x", e289.fieldType, INDIRENTRY)
		}
		if e289.fragment != 0 {
			t.Errorf("Entry new289 gets fragment index=%d, want 0", e289.fragment)
		}
		for i := 0; i < len(in289); i++ {
			if e289.inFields[i] != in289[i] {
				t.Errorf("Entry new289 inFields[%d]=%s, want %s", i, e289.inFields[i], in289[i])
			}
		}
	}

	// #293: Sindir check
	in293 := []string{"in1", "in2"}
	err = d.AddSindir("new293", in293[0], in293[1], 0)
	if err != nil {
		t.Error("Could not AddSindir in test 293:", err)
	}
	e293, err := d.Entry("new293")
	if err != nil {
		t.Error("Could not read Entry new293 in test 293:", err)
	} else {
		if e293.fieldType != SINDIRENTRY {
			t.Errorf("Entry new293 gets field type 0x%x, want SINDIR=0x%x", e293.fieldType, SINDIRENTRY)
		}
		if e293.fragment != 0 {
			t.Errorf("Entry new293 gets fragment index=%d, want 0", e293.fragment)
		}
		for i := 0; i < len(in293); i++ {
			if e293.inFields[i] != in293[i] {
				t.Errorf("Entry new293 inFields[%d]=%s, want %s", i, e293.inFields[i], in293[i])
			}
		}
	}

	// #302: IncludeNS
	idxfrag2, err := d.IncludeNS("format2", 0, "ns", CREAT|EXCL)
	if err != nil {
		t.Errorf("Could not Dirfile.IncludeNS")
	}
	if idxfrag2 != 2 {
		t.Errorf("IncludeNS returned fragment index %d, want 2", idxfrag2)
	}

	// #303: get namespace
	frag2, err := d.Fragment(2)
	if err != nil {
		t.Errorf("Could not open Fragment(2)")
	} else {
		if frag2.namespace != "ns" {
			t.Errorf("Fragment(2) namespace is %s, want \"ns\"", frag2.namespace)
		}

		// #304: SetNamespace
		err = frag2.SetNamespace("ns2")
		if err != nil {
			t.Errorf("Could not Fragment.SetNamespace()")
		}
		if frag2.namespace != "ns2" {
			t.Errorf("Fragment(2) namespace is %s, want \"ns2\"", frag2.namespace)
		}
	}

	// No #: test Uninclude, with and without del=true
	const DONTDELETE bool = false
	d.Uninclude(2, DONTDELETE)
	idxfragX, err := d.IncludeNS("formatXX", 0, "nsblah", CREAT|EXCL)
	if err != nil {
		t.Errorf("Problem unincluding Fragment 2")
	}
	fname := fmt.Sprintf("%s/formatXX", dir)
	_, err = os.Stat(fname)
	if err != nil {
		t.Errorf("Problem looking for %s", fname)
	}
	const DELETE bool = true
	d.Uninclude(idxfragX, DELETE)
	_, err = os.Stat(fname)
	if err == nil || !os.IsNotExist(err) {
		t.Errorf("Problem with %s, should not exist", fname)
	}

	// #305: MatchEntries
	matchedList, err := d.MatchEntries("^lin", 0, 0, 0)
	if err != nil {
		t.Errorf("Could not d.MatchEntries")
	} else if len(matchedList) != 2 {
		t.Errorf("d.MatchEntries returns list of length %d, want 2", len(matchedList))
	} else {
		trueList := []string{"lincom", "linterp"}
		for i := 0; i < len(trueList); i++ {
			if matchedList[i] != trueList[i] {
				t.Errorf("d.MatchEntries()[%d] = \"%s\", want \"%s\"",
					i, matchedList[i], trueList[i])
			}
		}
	}

	// No #: test discard
	err = d.Discard()
	if err != nil {
		t.Errorf("Could not discard dirfile read-only")
	}
}
