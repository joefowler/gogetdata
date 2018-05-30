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
	"fmt"
	"unsafe"
)

// Dirfile wraps the GetData.DIRFILE opaque object.
type Dirfile struct {
	name string
	nerr int
	d    *C.DIRFILE
}

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

// Error returns the latest error as a golang error type.
// It uses C API gd_error_string to generate the underlying string.
func (df *Dirfile) Error() error {
	df.nerr += int(C.gd_error_count(df.d))
	cmsg := C.gd_error_string(df.d, (*C.char)(C.NULL), 0)
	defer C.free(unsafe.Pointer(cmsg))
	return errors.New(C.GoString(cmsg))
}

// ErrorCount returns the number of errors raised by this Dirfile since the last
// call (i.e., this call resets the counter to zero).
func (df *Dirfile) ErrorCount() int {
	c := df.nerr
	df.nerr = 0
	return c
}

// Flags lets you modify flags which affect long-term operation only.
// These are VERBOSE and PRETTYPRINT.
func (df *Dirfile) Flags(set Flags, reset Flags) Flags {
	retval := C.gd_flags(df.d, C.ulong(set), C.ulong(reset))
	return Flags(retval)
}

// Close closes all open file handles and flushes all metadata.
func (df *Dirfile) Close() error {
	errcode := C.gd_close(df.d)
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	df.d = nil
	return nil
}

// Discard closes all open file handles but discards all metadata rather than
// flushing it to disk.
func (df *Dirfile) Discard() error {
	errcode := C.gd_discard(df.d)
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	df.d = nil
	return nil
}

// Flush will flush and close file descriptors associated with a field.
// Use Dirfile.Flush("") or (eqivalently) Dirfile.FlushAll() to cover all
// fields.
func (df *Dirfile) Flush(fieldcode string) error {
	if len(fieldcode) == 0 {
		return df.FlushAll()
	}
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	errcode := C.gd_flush(df.d, fcode)
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	return nil
}

// FlushAll will flush and close file descriptors associated with all fields.
// Use Dirfile.Flush("fieldname") to flush and close only one field.
func (df *Dirfile) FlushAll() error {
	errcode := C.gd_flush(df.d, (*C.char)(C.NULL))
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	return nil
}

// Sync flushes file descriptors associated with a field without closing them
func (df *Dirfile) Sync(fieldcode string) error {
	if len(fieldcode) == 0 {
		return df.SyncAll()
	}
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	errcode := C.gd_sync(df.d, fcode)
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	return nil
}

// SyncAll will flush file descriptors associated with all fields without closing them.
// Use Dirfile.Sync("fieldname") to flush only one field.
func (df *Dirfile) SyncAll() error {
	errcode := C.gd_sync(df.d, (*C.char)(C.NULL))
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	return nil
}

// RawClose closes file descriptors associated with a field without flushing them
func (df *Dirfile) RawClose(fieldcode string) error {
	if len(fieldcode) == 0 {
		return df.RawCloseAll()
	}
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	errcode := C.gd_raw_close(df.d, fcode)
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	return nil
}

// RawCloseAll will close file descriptors associated with all fields without flushing them.
// Use Dirfile.RawClose("fieldname") to close only one field.
func (df *Dirfile) RawCloseAll() error {
	errcode := C.gd_raw_close(df.d, (*C.char)(C.NULL))
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	return nil
}

// MetaFlush flushes the dirfile metadata to disk, without flushing the field data.
func (df *Dirfile) MetaFlush() error {
	errcode := C.gd_metaflush(df.d)
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	return nil
}

// VerbosePrefix sets the prefix used when dirfile is in VERBOSE mode.
func (df *Dirfile) VerbosePrefix(prefix string) error {
	vpre := C.CString(prefix)
	defer C.free(unsafe.Pointer(vpre))
	errcode := C.gd_verbose_prefix(df.d, vpre)
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	return nil
}

// InvalidDirfile creates a Dirfile instance whose methods will always produce
// a GD_E_BAD_DIRFILE error.
func InvalidDirfile() Dirfile {
	df := C.gd_invalid_dirfile()
	return Dirfile{name: "invalid", d: df}
}

// Desync detects desynchronization of a dirfile stored on disk.
// See C API for pathcheck and reopen meaning.
func (df *Dirfile) Desync(pathcheck, reopen bool) (bool, error) {
	flags := C.uint(0)
	if pathcheck {
		flags |= C.GD_DESYNC_PATHCHECK
	}
	if reopen {
		flags |= C.GD_DESYNC_REOPEN
	}
	result := int(C.gd_desync(df.d, flags))
	if result < 0 {
		return false, df.Error()
	}
	return result > 0, nil
}

// GetData fetches data from a vector field in the dirfile (incl. metafields)
// out should be a *pointer to* a slice of numeric data of adequate size, e.g.
// d := make([]int32, 20)
// df.GetData("field", 5, 0, 2, 0, &d)
// Returns (n, err) where n is the number of samples read.
func (df Dirfile) GetData(fieldcode string, firstFrame, firstSample, numFrames, numSamples int, out interface{}) (int, error) {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	// expectedSamples := numSamples + numFrames*df.
	retType, ptr := parray2type(out)
	if retType == UNKNOWN || retType == STRING || ptr == C.NULL {
		return 0, fmt.Errorf("GetData out variable was not a pointer to numeric slice")
	}
	n := C.gd_getdata(df.d, fcode, C.off_t(firstFrame), C.off_t(firstSample),
		C.size_t(numFrames), C.size_t(numSamples), C.gd_type_t(retType), ptr)
	if n == 0 {
		return 0, df.Error()
	}
	return int(n), nil
}

// MplexLookback changes how far GetData searches backwards for the initial
// value of a field when reading a MPLEX field
func (df *Dirfile) MplexLookback(lookback int) {
	C.gd_mplex_lookback(df.d, C.int(lookback))
}

// GetConstant fills the numeric type pointed to by inptry with the constant or metadata field named fieldcode
func (df Dirfile) GetConstant(fieldcode string, inptr interface{}) error {
	typecode, uptr := pointer2type(inptr)
	if typecode == UNKNOWN {
		return fmt.Errorf("GetConstant called with ptr not a pointer to string or numeric type")
	}

	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	errcode := C.gd_get_constant(df.d, fcode, C.gd_type_t(typecode), uptr)
	if errcode != C.GD_E_OK {
		return df.Error()
	}
	return nil
}

// GetConstantInt32 returns an int32 for the constant or metadata field named fieldcode
func (df Dirfile) GetConstantInt32(fieldcode string) (int32, error) {
	var c int32
	return c, df.GetConstant(fieldcode, &c)
}

// GetConstantInt64 returns an int64 for the constant or metadata field named fieldcode
func (df Dirfile) GetConstantInt64(fieldcode string) (int64, error) {
	var c int64
	return c, df.GetConstant(fieldcode, &c)
}

// GetConstantFloat32 returns a float32 for the constant or metadata field named fieldcode
func (df Dirfile) GetConstantFloat32(fieldcode string) (float32, error) {
	var c float32
	return c, df.GetConstant(fieldcode, &c)
}

// GetConstantFloat64 returns a float64 for the constant or metadata field named fieldcode
func (df Dirfile) GetConstantFloat64(fieldcode string) (float64, error) {
	var c float64
	return c, df.GetConstant(fieldcode, &c)
}

// GetConstantComplex64 returns a complex64 for the constant or metadata field named fieldcode
func (df Dirfile) GetConstantComplex64(fieldcode string) (complex64, error) {
	var c complex64
	return c, df.GetConstant(fieldcode, &c)
}

// GetConstantComplex128 returns a complex128 for the constant or metadata field named fieldcode
func (df Dirfile) GetConstantComplex128(fieldcode string) (complex128, error) {
	var c complex128
	return c, df.GetConstant(fieldcode, &c)
}

// GetCarray fills the numeric array pointed to by out with a list of the values
// of all elements in a CARRAY field (including metafields).
func (df Dirfile) GetCarray(fieldcode string, out interface{}) error {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	retType, ptr := parray2type(out)
	if retType == UNKNOWN || retType == STRING || ptr == C.NULL {
		return fmt.Errorf("GetCarray out variable was not a pointer to numeric slice")
	}
	result := int(C.gd_get_carray(df.d, fcode, C.gd_type_t(retType), ptr))
	if result < 0 {
		return df.Error()
	}
	return nil
}

// GetCarraySlice fills the numeric array pointed to by out with a list a portion of
// the elements in a CARRAY field (including metafields).
func (df Dirfile) GetCarraySlice(fieldcode string, start, n uint, out interface{}) error {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	retType, ptr := parray2type(out)
	if retType == UNKNOWN || retType == STRING || ptr == C.NULL {
		return fmt.Errorf("GetCarray out variable was not a pointer to numeric slice")
	}
	result := int(C.gd_get_carray_slice(df.d, fcode, C.ulong(start), C.size_t(n), C.gd_type_t(retType), ptr))
	if result < 0 {
		return df.Error()
	}
	return nil
}

// func (df Dirfile) GetCarrays(fieldcode string, retType RetType) []
// Hmm. Not sure how to do this one!

// GetString returns the value of a STRING field (including metafields)
func (df Dirfile) GetString(fieldcode string) (string, error) {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	bsize := C.size_t(256)
	cresult := (*C.char)(C.malloc(bsize))
	n := int(C.gd_get_string(df.d, fcode, bsize, cresult))
	if n == 0 {
		return "", df.Error()
	}
	return C.GoString(cresult), nil
}

// Strings returns the value of all STRING fields (including metafields)
func (df Dirfile) Strings() ([]string, error) {
	cptr := C.gd_strings(df.d)
	if cptr == (**C.char)(C.NULL) {
		return nil, fmt.Errorf("Dirfile.Strings returned NULL")
	}
	var result []string
	listend := (*C.char)(C.NULL)
	cstr0 := *(**C.char)(cptr)
	for cstr0 != listend {
		result = append(result, C.GoString(cstr0))
		// fmt.Printf("cptr=%p, cstr0=%p, result=%s ", cptr, cstr0, result)
		cptr = (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(&cstr0)) + unsafe.Sizeof(cstr0)))
		cstr0 = *(**C.char)(cptr)
		// fmt.Printf("new cptr=%p, cstr0=%p\n", cptr, cstr0)
	}
	return result, nil
}

// PutData stores data to a vector field in the dirfile (incl. metafields)
// data should be a slice of numeric data, e.g.
// var d []int32{4,5,6,7,8}
// n, err := df.PutData("field", 5, 0, d)
// if n != len(d) || err != nil {...<problem>...}
func (df *Dirfile) PutData(fieldcode string, firstFrame, firstSample int, data interface{}) (int, error) {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	dType, ptr, lenData := array2type(data)
	if dType == UNKNOWN || dType == NULLTYPE || ptr == C.NULL {
		return 0, fmt.Errorf("PutData data variable was not a numeric slice")
	}
	n := C.gd_putdata(df.d, fcode, C.off_t(firstFrame), C.off_t(firstSample),
		C.size_t(0), C.size_t(lenData), C.gd_type_t(dType), ptr)
	if n == 0 {
		return 0, df.Error()
	}
	return int(n), nil
}

// PutConstant stores the value of a CONST field (including metafields)
// data should be a value of some numeric type.
func (df *Dirfile) PutConstant(fieldcode string, data interface{}) error {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	dType, ptr := value2type(data)
	if dType == UNKNOWN || dType == NULLTYPE || ptr == C.NULL {
		return fmt.Errorf("PutConstant data variable was not a numeric value")
	}
	n := C.gd_put_constant(df.d, fcode, C.gd_type_t(dType), ptr)
	if n != 0 {
		return df.Error()
	}
	return nil
}

// PutCarray stores an entire CARRAY field (including metafields)
func (df *Dirfile) PutCarray(fieldcode string, array interface{}) error {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	dType, ptr, _ := array2type(array)
	n := C.gd_put_carray(df.d, fcode, C.gd_type_t(dType), ptr)
	if n != 0 {
		return df.Error()
	}
	return nil
}

// PutString stores the value of a STRING field (including metafields)
func (df *Dirfile) PutString(fieldcode, value string) error {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	newval := C.CString(value)
	defer C.free(unsafe.Pointer(newval))
	n := C.gd_put_string(df.d, fcode, newval)
	if n != 0 {
		return df.Error()
	}
	return nil
}

// Dirfilename returns the full path to the dirfile
func (df Dirfile) Dirfilename() string {
	result := C.gd_dirfilename(df.d)
	return C.GoString(result)
}

// NFrames returns the number of frames in the dirfile (or on error, 0)
func (df Dirfile) NFrames() int {
	return int(C.gd_nframes(df.d))
}

// NFragments returns the number of fragments in the dirfile (or on error, 0)
func (df Dirfile) NFragments() int {
	return int(C.gd_nfragments(df.d))
}

// Fragment returns a Fragment pointer to the nth dirfile fragment
func (df Dirfile) Fragment(n int) (*Fragment, error) {
	return NewFragment(&df, n)
}

// ArrayLen returns the number of elements in a scalar field (CARRAY, CONST,
// or STRING)
func (df Dirfile) ArrayLen(fieldcode string) int {
	fcode := C.CString(fieldcode)
	defer C.free(unsafe.Pointer(fcode))
	return int(C.gd_array_len(df.d, fcode))
}

// NEntries returns the number of fields in the dirfile satisfying various criteria.
func (df Dirfile) NEntries(parent string, etype EntryType, flags EntryType) uint {
	cparent := C.CString(parent)
	defer C.free(unsafe.Pointer(cparent))
	if len(parent) == 0 {
		cparent = (*C.char)(C.NULL)
	}
	return uint(C.gd_nentries(df.d, cparent, C.int(etype), C.uint(flags)))
}

// NFields returns the number of fields in the dirfile
func (df Dirfile) NFields() uint {
	return uint(C.gd_nfields(df.d))
}

// NVectors returns the number of vector fields (that is all field types except
// CONST, CARRAY, and STRING) in the dirfile
func (df Dirfile) NVectors() uint {
	return uint(C.gd_nvectors(df.d))
}

// NFieldsByType returns the number of fields in the dirfile.
func (df Dirfile) NFieldsByType(etype EntryType) uint {
	return uint(C.gd_nfields_by_type(df.d, C.gd_entype_t(etype)))
}

// NMFields returns the number of metafields in the dirfile for a particular parent.
func (df Dirfile) NMFields(parent string) uint {
	cparent := C.CString(parent)
	defer C.free(unsafe.Pointer(cparent))
	if len(parent) == 0 {
		cparent = (*C.char)(C.NULL)
	}
	return uint(C.gd_nmfields(df.d, cparent))
}

// NMVectors returns the number of vector metafields in the dirfile for a particular parent.
// (That is, all field types except CONST, CARRAY, and STRING.)
func (df Dirfile) NMVectors(parent string) uint {
	cparent := C.CString(parent)
	defer C.free(unsafe.Pointer(cparent))
	if len(parent) == 0 {
		cparent = (*C.char)(C.NULL)
	}
	return uint(C.gd_nmvectors(df.d, cparent))
}

// NMFieldsByType returns the number of metafields in the dirfile for a particular
// parent and a specified field type.
func (df Dirfile) NMFieldsByType(parent string, etype EntryType) uint {
	cparent := C.CString(parent)
	defer C.free(unsafe.Pointer(cparent))
	if len(parent) == 0 {
		cparent = (*C.char)(C.NULL)
	}
	return uint(C.gd_nmfields_by_type(df.d, cparent, C.gd_entype_t(etype)))
}

func ppchar2stringSlice(listptr unsafe.Pointer) []string {
	var fields []string
	for i := 0; ; i++ {
		fields = append(fields, C.GoString(*(**C.char)(listptr)))
		listptr = unsafe.Pointer(uintptr(listptr) + unsafe.Sizeof(uintptr(0)))
		if unsafe.Pointer(*(**C.char)(listptr)) == (C.NULL) {
			break
		}
	}
	return fields
}

// MatchEntries returns a list of entries in the dirfile satisfying various criteria
func (df Dirfile) MatchEntries(regex string, fragment int, et EntryType, flags Flags) ([]string, error) {
	cregex := C.CString(regex)
	defer C.free(unsafe.Pointer(cregex))
	var ptr (**C.char)
	n := C.gd_match_entries(df.d, cregex, C.int(fragment), C.int(et), C.uint(flags), &ptr)
	if n < 0 {
		return nil, df.Error()
	}
	return ppchar2stringSlice(unsafe.Pointer(ptr)), nil
}

// EntryList returns a slice of strings listing all fields meeting various criteria.
func (df Dirfile) EntryList(parent string, et EntryType, flags EntryType) []string {
	cparent := C.CString(parent)
	defer C.free(unsafe.Pointer(cparent))
	if len(parent) == 0 {
		cparent = (*C.char)(C.NULL)
	}
	return ppchar2stringSlice(unsafe.Pointer(C.gd_entry_list(df.d, cparent, C.int(et), C.uint(flags))))
}

// FieldList returns a slice of strings listing all fields (no metafields).
func (df Dirfile) FieldList() []string {
	return ppchar2stringSlice(unsafe.Pointer(C.gd_field_list(df.d)))
}

// VectorList returns a slice of strings listing all vector fields (no metafields).
func (df Dirfile) VectorList() []string {
	return ppchar2stringSlice(unsafe.Pointer(C.gd_vector_list(df.d)))
}

// FieldListByType returns a slice of strings listing all fields (no metafields).
func (df Dirfile) FieldListByType(et EntryType) []string {
	return ppchar2stringSlice(unsafe.Pointer(C.gd_field_list_by_type(df.d, C.gd_entype_t(et))))
}

// MFieldList returns a slice of strings listing all metafields in the dirfile for a particular parent.
func (df Dirfile) MFieldList(parent string) []string {
	cparent := C.CString(parent)
	defer C.free(unsafe.Pointer(cparent))
	if len(parent) == 0 {
		cparent = (*C.char)(C.NULL)
	}
	return ppchar2stringSlice(unsafe.Pointer(C.gd_mfield_list(df.d, cparent)))
}

// MVectorList returns a slice of strings listing all vector metafields in the dirfile for a particular parent.
func (df Dirfile) MVectorList(parent string) []string {
	cparent := C.CString(parent)
	defer C.free(unsafe.Pointer(cparent))
	if len(parent) == 0 {
		cparent = (*C.char)(C.NULL)
	}
	return ppchar2stringSlice(unsafe.Pointer(C.gd_mvector_list(df.d, cparent)))
}

// MFieldListByType returns a slice of strings listing all metafields of a specified type for a particular parent.
func (df Dirfile) MFieldListByType(parent string, et EntryType) []string {
	cparent := C.CString(parent)
	defer C.free(unsafe.Pointer(cparent))
	if len(parent) == 0 {
		cparent = (*C.char)(C.NULL)
	}
	return ppchar2stringSlice(unsafe.Pointer(
		C.gd_mfield_list_by_type(df.d, cparent, C.gd_entype_t(et))))
}

// Include adds the named fragment to the dirfile.
func (df *Dirfile) Include(file string, flags Flags) (int, error) {
	return df.IncludeAtIndex(file, 0, flags)
}

// IncludeAtIndex adds the named fragment to the dirfile at the given index.
func (df *Dirfile) IncludeAtIndex(file string, index int, flags Flags) (int, error) {
	fragmentname := C.CString(file)
	defer C.free(unsafe.Pointer(fragmentname))
	result := int(C.gd_include(df.d, fragmentname, C.int(index), C.ulong(flags)))
	if result < 0 {
		return result, df.Error()
	}
	return result, nil
}

// IncludeNS adds the named fragment to the dirfile at the given index, adding a namespace.
func (df *Dirfile) IncludeNS(file string, index int, namespace string, flags Flags) (int, error) {
	fragmentname := C.CString(file)
	defer C.free(unsafe.Pointer(fragmentname))
	cnamespace := C.CString(namespace)
	defer C.free(unsafe.Pointer(cnamespace))
	result := int(C.gd_include_ns(df.d, fragmentname, C.int(index), cnamespace, C.ulong(flags)))
	if result < 0 {
		return result, df.Error()
	}
	return result, nil
}

// Uninclude removes the fragment from the dirfile at the given index.
func (df *Dirfile) Uninclude(index int, del bool) error {
	cdel := C.int(0)
	if del {
		cdel = C.int(1)
	}
	result := int(C.gd_uninclude(df.d, C.int(index), cdel))
	if result != C.GD_E_OK {
		return df.Error()
	}
	return nil
}
