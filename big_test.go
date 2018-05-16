package getdata

import (
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func createTestDirfile(dir string) {
	err := os.Mkdir(dir, 0775)
	if err != nil && !os.IsExist(err) {
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

	// #3-5: getdata (int) check
	// var n int
	// out = d.GetData("data", 5, 0, 1, 0, out)

	// #12: constant (int) check
	i32, err := d.GetConstantInt32("const")
	if err != nil {
		t.Errorf("Could not GetConstantInt32")
	}
	if i32 != 5 {
		t.Errorf("GetConstantInt32 returns %d, want 5", i32)
	}
	i32, err = d.GetConstantInt32("doesnt exist")
	if err == nil {
		t.Errorf("GetConstantInt32 returned %d on non-existent field, want error", i32)
	}

	// #17: constant (float) check
	f32, err := d.GetConstantFloat32("const")
	if err != nil {
		t.Errorf("Could not GetConstantFloat32")
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

	// #28: nframes check
	nf := d.NFrames()
	if nf != 10 {
		t.Errorf("NFrames returned %d, want 10", nf)
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

	// #235: flags check
	d.Flags(VERBOSE, 0)
	flags := d.Flags(PRETTYPRINT, 0)
	if flags&PRETTYPRINT == 0 {
		t.Errorf("Flags(0x%x, 0x0) returned 0x%x, want that flag set", PRETTYPRINT, flags)
	}
	flags = d.Flags(0, PRETTYPRINT)
	if flags&PRETTYPRINT != 0 {
		t.Errorf("Flags(0x0, 0x%x) returned 0x%x, want that flag clear", PRETTYPRINT, flags)
	}

	// #236: verbose_prefix check
	err = d.VerbosePrefix("big_test: ")
	if err != nil {
		t.Errorf("Could not set VerbosePrefix()")
	}

	// No #: test discard
	err = d.Discard()
	if err != nil {
		t.Errorf("Could not discard dirfile read-only")
	}
}
