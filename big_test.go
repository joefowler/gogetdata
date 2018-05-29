package getdata

import (
	"fmt"
	"log"
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

	// #6: getdata (int64) check
	u2 := make([]uint64, 8)
	n, err = d.GetData("data", 5, 0, 1, 0, &u2)
	if err != nil {
		t.Error("Could not GetData: ", err)
	} else if len(u2) < 8 {
		t.Errorf("GetData out has len=%d (cap=%d), want 8", len(u2), cap(u2))
	} else if n != 8 {
		t.Errorf("GetData returned %d, want 8", n)
	} else {
		for i := 0; i < 8; i++ {
			if u2[i] != uint64(41+i) {
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
	}
	if i32 != 5 {
		t.Errorf("GetConstantInt32 returns %d, want 5", i32)
	}
	i32, err = d.GetConstantInt32("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantInt32 returned %d on non-existent field, want error", i32)
	}

	// #15: constant (int64) check
	i64, err := d.GetConstantInt64("const")
	if err != nil {
		t.Error("Could not GetConstantInt64: ", err)
	}
	if i64 != 5 {
		t.Errorf("GetConstantInt64 returns %d, want 5", i64)
	}
	i64, err = d.GetConstantInt64("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantInt64 returned %d on non-existent field, want error", i64)
	}

	// #17: constant (float) check
	f32, err := d.GetConstantFloat32("const")
	if err != nil {
		t.Error("Could not GetConstantFloat32: ", err)
	}
	if f32 != 5.5 {
		t.Errorf("GetConstantFloat32 returns %f, want 5.5", f32)
	}
	f32, err = d.GetConstantFloat32("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantFloat32 returned %f on non-existent field, want error", f32)
	}
	f64, err := d.GetConstantFloat64("const")
	if err != nil {
		t.Errorf("Could not GetConstantFloat64")
	}
	if f64 != 5.5 {
		t.Errorf("GetConstantFloat64 returns %f, want 5.5", f64)
	}
	f64, err = d.GetConstantFloat64("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantFloat64 returned %f on non-existent field, want error", f64)
	}

	// #19: constant (complex) check
	c64, err := d.GetConstantComplex64("const")
	if err != nil {
		t.Errorf("Could not GetConstantComplex64")
	}
	if c64 != 5.5 {
		t.Errorf("GetConstantComplex64 returns %f, want 5.5", c64)
	}
	c64, err = d.GetConstantComplex64("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantComplex64 returned %f on non-existent field, want error", c64)
	}
	c128, err := d.GetConstantComplex128("const")
	if err != nil {
		t.Errorf("Could not GetConstantComplex128")
	}
	if c128 != 5.5 {
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
	if nlincom != 1 { // TODO: update to 3 when we've added vectors in earlier tests
		t.Errorf("NFieldsByType(LINCOMENTRY) returned %d, want 1", nlincom)
	}

	// #68: FieldListByType check
	ftnames := d.FieldListByType(LINCOMENTRY)
	// trueftnames := []string{"new2","new3","lincom"} TODO: use this line
	trueftnames := []string{"lincom"}
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
	if nvec != 15 { // TODO: update to 25 when we've added vectors in earlier tests
		t.Errorf("NVectors = %d, want 20", nvec)
	}

	// #70: VectorList check
	vectors := d.VectorList()
	if len(vectors) != int(nvec) {
		t.Errorf("VectorList length is %d, want %d", len(vectors), nvec)
	} else {
		// truevnames := []string{"bit", "div", "data", "mult", "new1", "new2", "new3",
		// 	"new4", "new5", "new6", "new7", "new8", "new9", "sbit", "INDEX",
		// 	"alias", "indir", "mplex", "new10", "phase", "recip", "lincom",
		// 	"window", "linterp", "polynom"}
		// TODO: fix when all tests are in place. Above is the true answer.
		truevnames := []string{"bit", "div", "data", "mult", "sbit", "INDEX",
			"alias", "indir", "mplex", "phase", "recip", "lincom",
			"window", "linterp", "polynom"}
		for i := 0; i < int(nvec); i++ {
			if vectors[i] != truevnames[i] {
				t.Errorf("FieldList[%d]=\"%s\", want \"%s\"", i, vectors[i], truevnames[i])
			}
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

	// #138: PutData (auto) check duplicates #37, so skip it.

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

	// #208 sync check
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

	// #201 metaflush check
	err = d.MetaFlush()
	if err != nil {
		t.Errorf("Could not call MetaFlush")
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
	if ne != 14 { // TODO: eventually 27
		t.Errorf("d.NEntries counts %d VECTOR entries, want 14", ne)
	}

	// #239: EntryList check
	entryList := d.EntryList("", VECTORENTRIES, HIDDENENTRIES|NOALIASENTRIES)
	if len(entryList) != int(ne) {
		t.Errorf("d.EntryList return %d entries, want %d", len(entryList), ne)
	} else {
		// trueEntries := []string{"bit", "div", "data", "mult", "new1", "new2", "new3",
		// 	"new4", "new5", "new6", "new7", "new8", "sbit", "INDEX",
		// 	"indir", "mplex", "new14", "new15", "new16", "new18", "new21", "phase", "recip", "lincom",
		// 	"window", "linterp", "polynom"}
		// TODO: use the above
		trueEntries := []string{"bit", "div", "data", "mult", "sbit", "INDEX",
			"indir", "mplex", "phase", "recip", "lincom",
			"window", "linterp", "polynom"}
		for i := 0; i < int(ne); i++ {
			if entryList[i] != trueEntries[i] {
				t.Errorf("EntryList[%d]=\"%s\", want \"%s\"", i, entryList[i], trueEntries[i])
			}
		}
	}

	// #240: MplexLookback test (returns nothing, so simply call it)
	d.MplexLookback(LOOKBACKALL)

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
