package getdata

/*
#cgo CFLAGS: -I/usr/local/include -DGD_C89_API
#cgo LDFLAGS: -L/usr/local/lib -lgetdata
#include <getdata.h>
#include <stdlib.h>
*/
import "C"
import (
	"unsafe"
)

type raw struct {
	spf      uint
	dataType RetType
}

// MAXLINCOM is the maximum number of linear combination inputs
const MAXLINCOM int = int(C.GD_MAX_LINCOM)

type lincom struct {
	nFields int
	cm      [MAXLINCOM]complex128
	m       [MAXLINCOM]float64
	cb      [MAXLINCOM]complex128
	b       [MAXLINCOM]float64
}

// MAXPOLYORD is the maximum polynomial order
const MAXPOLYORD int = int(C.GD_MAX_POLYORD)

type polynomial struct {
	polyOrder int
	a         [MAXPOLYORD + 1]float64
	ca        [MAXPOLYORD + 1]complex128
}

type bits struct {
	bitnum  int
	numbits int
}

type reciprocal struct {
	dividend  float64
	cdividend complex128
}

type mplex struct {
	countVal int
	period   int
}

type window struct {
	windOp     WindowOps
	iThreshold int64
	uThreshold uint64
	fThreshold float64
}

// Entry wraps the gd_entry_t object, and is used to access field metadata
type Entry struct {
	df        *Dirfile
	name      string
	fieldType EntryType
	flags     uint
	fragment  int
	raw
	inFields [MAXLINCOM]string
	lincom
	polynomial
	table string
	bits
	phaseShift int64
	reciprocal
	mplex
	window
	constType RetType
	arrayLen  int
}

// func entryToC(e *Entry) *C.gd_entry_t {
// 	cp := (*C.gd_entry_t)(C.malloc(C.sizeof_gd_entry_t))
// 	ce := *cp
// 	ce.field = C.CString(e.name)
// 	ce.field_type = C.gd_entype_t(e.fieldType)
// 	ce.flags = C.uint(e.flags)
// 	for i := 0; i < MAXLINCOM; i++ {
// 		ce.in_fields[i] = C.CString(e.inFields[i])
// 	}
// 	return cp
// }

func entryFromC(df *Dirfile, ce *C.gd_entry_t) Entry {
	e := Entry{
		df:        df,
		name:      C.GoString(ce.field),
		fieldType: EntryType(ce.field_type),
		flags:     uint(ce.flags),
	}
	for i := 0; i < MAXLINCOM; i++ {
		e.inFields[i] = C.GoString(ce.in_fields[i])
	}

	base := uintptr(unsafe.Pointer(&ce.flags)) + unsafe.Sizeof(ce.flags)
	switch e.fieldType {
	case RAWENTRY:
		e.spf = uint(*(*C.uint)(unsafe.Pointer(base)))
		base += unsafe.Sizeof(C.uint(0))
		e.dataType = RetType(*(*C.gd_type_t)(unsafe.Pointer(base)))

	case LINCOMENTRY:
		e.nFields = int(*(*C.long)(unsafe.Pointer(base)))
		base += unsafe.Sizeof(C.long(0))
		for i := 0; i < MAXLINCOM; i++ {
			e.m[i] = float64(*(*C.double)(unsafe.Pointer(base)))
			base += unsafe.Sizeof(C.double(0))
		}
		for i := 0; i < MAXLINCOM; i++ {
			e.cm[i] = complex128(*(*C.complexdouble)(unsafe.Pointer(base)))
			base += unsafe.Sizeof(C.complexdouble(0))
		}
		for i := 0; i < MAXLINCOM; i++ {
			e.b[i] = float64(*(*C.double)(unsafe.Pointer(base)))
			base += unsafe.Sizeof(C.double(0))
		}
		for i := 0; i < MAXLINCOM; i++ {
			e.cb[i] = complex128(*(*C.complexdouble)(unsafe.Pointer(base)))
			base += unsafe.Sizeof(C.complexdouble(0))
		}

	case POLYNOMENTRY:
		e.polyOrder = int(*(*C.int)(unsafe.Pointer(base)))
		base += unsafe.Sizeof(C.int(0))
		base += 4 // because of alignment
		for i := 0; i < MAXPOLYORD+1; i++ {
			e.a[i] = float64(*(*C.double)(unsafe.Pointer(base)))
			base += unsafe.Sizeof(C.double(0))
		}
		for i := 0; i < MAXPOLYORD+1; i++ {
			e.ca[i] = complex128(*(*C.complexdouble)(unsafe.Pointer(base)))
			base += unsafe.Sizeof(C.complexdouble(0))
		}

	case LINTERPENTRY:
		e.table = C.GoString(*(**C.char)(unsafe.Pointer(base)))

	case BITENTRY, SBITENTRY:
		e.bitnum = int(*(*C.int)(unsafe.Pointer(base)))
		base += unsafe.Sizeof(C.int(0))
		e.numbits = int(*(*C.int)(unsafe.Pointer(base)))

	case RECIPENTRY:
		e.dividend = float64(*(*C.double)(unsafe.Pointer(base)))
		base += unsafe.Sizeof(C.double(0))
		e.cdividend = complex128(*(*C.complexdouble)(unsafe.Pointer(base)))

	case PHASEENTRY:
		e.phaseShift = int64(*(*C.gd_int64_t)(unsafe.Pointer(base)))

	case MPLEXENTRY:
		e.countVal = int(*(*C.int)(unsafe.Pointer(base)))
		base += unsafe.Sizeof(C.int(0))
		e.period = int(*(*C.int)(unsafe.Pointer(base)))

	case CONSTENTRY, SARRAYENTRY:
		e.constType = RetType(*(*C.gd_type_t)(unsafe.Pointer(base)))
		base += unsafe.Sizeof(C.long(0))
		e.arrayLen = int(*(*C.size_t)(unsafe.Pointer(base)))

	case WINDOWENTRY:
		e.windOp = WindowOps(*(*C.gd_windop_t)(unsafe.Pointer(base)))
		base += unsafe.Sizeof(C.gd_windop_t(0))
		// Handle the 3-way union called gd_triplet_t
		e.iThreshold = int64(*(*int64)(unsafe.Pointer(base)))
		e.uThreshold = uint64(*(*uint64)(unsafe.Pointer(base)))
		e.fThreshold = float64(*(*float64)(unsafe.Pointer(base)))

		// for i := 0; i < 19; i++ {
		// 	fmt.Printf("%16.16x\n", *(*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(base)) + uintptr(i*8))))
		// }
	}
	return e
}

// RawEntry creates an Entry of RAW type without adding to any Dirfile
func RawEntry(name string, fragmentIndex int, samplesPerFrame uint, dType RetType) Entry {
	var e = Entry{
		name:      name,
		fieldType: RAWENTRY,
		fragment:  fragmentIndex,
	}
	e.spf = samplesPerFrame
	e.dataType = dType
	return e
}

// BitEntry creates an Entry of BIT type without adding to any Dirfile
func BitEntry(name, inField string, bitnum, numbits, fragmentIndex int) Entry {
	var e = Entry{
		name:      name,
		fieldType: BITENTRY,
		fragment:  fragmentIndex,
	}
	e.inFields[0] = inField
	e.numbits = numbits
	e.bitnum = bitnum
	return e
}

// Filename returns the raw dirfile's filename
func (e Entry) Filename() (string, error) {
	return e.df.Filename(e.name)
}

// func (e *Entry) SetType(t RetType)

// Move moves this entry to a new fragment number
func (e *Entry) Move(newfrag int, flags RenameFlags) error {
	fcode := C.CString(e.name)
	defer C.free(unsafe.Pointer(fcode))
	cidx := C.int(newfrag)
	result := C.gd_move(e.df.d, fcode, cidx, C.uint(flags))
	if result < 0 {
		return e.df.Error()
	}
	e.fragment = newfrag
	return nil
}

// Rename renames this entry to the indicated name
func (e *Entry) Rename(newname string, flags RenameFlags) error {
	fcode := C.CString(e.name)
	defer C.free(unsafe.Pointer(fcode))
	ncode := C.CString(newname)
	defer C.free(unsafe.Pointer(ncode))
	result := C.gd_rename(e.df.d, fcode, ncode, C.uint(flags))
	if result < 0 {
		return e.df.Error()
	}
	e.name = newname
	return nil
}
