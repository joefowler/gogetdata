package getdata

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lgetdata
#include <getdata.h>
#include <stdlib.h>
*/
import "C"

///

// Flags are dirfile-opening flags, including encoding methods
type Flags int64

// RDONLY open read-only
const RDONLY Flags = C.GD_RDONLY

// RDWR open read/write
const RDWR Flags = C.GD_RDWR

// FORCEENDIAN override endianness
const FORCEENDIAN Flags = C.GD_FORCE_ENDIAN

// BIGENDIAN specifies big-endian raw data
const BIGENDIAN Flags = C.GD_BIG_ENDIAN

// LITTLEENDIAN specifies little-endian raw data
const LITTLEENDIAN Flags = C.GD_LITTLE_ENDIAN

// NATIVEENDIAN specifies native-endian raw data
const NATIVEENDIAN Flags = 0

// NONNATIVEENDIAN specifies the opposite of native-endian raw data
const NONNATIVEENDIAN Flags = BIGENDIAN | LITTLEENDIAN

// CREAT create dirfile if it doesn't exist
const CREAT Flags = C.GD_CREAT

// EXCL forces creation of dirfile (and fail if it exists)
const EXCL Flags = C.GD_EXCL

// TRUNC if the dirfile already exists, truncate it before opening it. Truncating a
// dirfile deletes all files in the specified dirfile directory, so use this flag with caution.
const TRUNC Flags = C.GD_TRUNC

// PEDANTIC makes the dirfile instist on strict adherence to standards
const PEDANTIC Flags = C.GD_PEDANTIC

// FORCEENCODING makes dirfile ignore any encoding specified in the dirfile itself: just use the encoding specified by these flags.
const FORCEENCODING Flags = C.GD_FORCE_ENCODING

// VERBOSE writes error messages to standard error automatically when errors are triggered
const VERBOSE Flags = C.GD_VERBOSE

// IGNOREDUPS ignore duplicate field names while parsing the dirfile metadata
const IGNOREDUPS Flags = C.GD_IGNORE_DUPS

// IGNOREREFS ignore /REFERENCE directives while parsing the dirfile metadata
const IGNOREREFS Flags = C.GD_IGNORE_REFS

// PRETTYPRINT attempt to make a nicer looking format specification (in the
// human-readable sense) when writing metadata to disk.
const PRETTYPRINT Flags = C.GD_PRETTY_PRINT

// PERMISSIVE accepts non-compliant syntax, even if the dirfile contains a /VERSION directive.
const PERMISSIVE Flags = C.GD_PERMISSIVE

// TRUNCSUB if truncating a dirfile, also delete subdirectories. Ignored if TRUNC is not also specified.
const TRUNCSUB Flags = C.GD_TRUNCSUB

///

// EntryType signifies the field type given for entries in the FORMAT files
type EntryType int64

// NOENTRY denotes an invalid entry type
const NOENTRY EntryType = C.GD_NO_ENTRY

// BITENTRY denotes one or more bits out of an input vector field, treating the result as unsigned
const BITENTRY EntryType = C.GD_BIT_ENTRY

// CARRAYENTRY denotes array of constants fully specified in the format file metadata
const CARRAYENTRY EntryType = C.GD_CARRAY_ENTRY

// CONSTENTRY denotes a scalar constant fully specified in the format file metadata
const CONSTENTRY EntryType = C.GD_CONST_ENTRY

// DIVIDEENTRY denotes the quotient of two vector fields
const DIVIDEENTRY EntryType = C.GD_DIVIDE_ENTRY

// INDIRENTRY denotes indirection of a CARRAY scalar indexed by a vector field
const INDIRENTRY EntryType = C.GD_INDIR_ENTRY

// LINCOMENTRY denotes a linear combination of 1, 2, or 3 vector fields
const LINCOMENTRY EntryType = C.GD_LINCOM_ENTRY

// LINTERPENTRY denotes a linearly interpolated table lookup
const LINTERPENTRY EntryType = C.GD_LINTERP_ENTRY

// MPLEXENTRY denotes multiplexing of several low-rate fields into a single one
const MPLEXENTRY EntryType = C.GD_MPLEX_ENTRY

// MULTIPLYENTRY denotes the product of two vector fields
const MULTIPLYENTRY EntryType = C.GD_MULTIPLY_ENTRY

// PHASEENTRY denotes a vector field shifted in time by a specified number of samples
const PHASEENTRY EntryType = C.GD_PHASE_ENTRY

// POLYNOMENTRY denotes a polynomial function of a single input field
const POLYNOMENTRY EntryType = C.GD_POLYNOM_ENTRY

// RAWENTRY denotes time streams on disk
const RAWENTRY EntryType = C.GD_RAW_ENTRY

// RECIPENTRY denotes the reciprocal of an input field
const RECIPENTRY EntryType = C.GD_RECIP_ENTRY

// SARRAYENTRY denotes array of strings fully specified in the format file metadata
const SARRAYENTRY EntryType = C.GD_SARRAY_ENTRY

// SBITENTRY denotes one or more bits out of an input vector field, treating the result as signed
const SBITENTRY EntryType = C.GD_SBIT_ENTRY

// SINDIRENTRY denotes indirection of a SARRAY scalar indexed by a vector field
const SINDIRENTRY EntryType = C.GD_SINDIR_ENTRY

// STRINGENTRY denotes a single character string fully specified in the format file metadata
const STRINGENTRY EntryType = C.GD_STRING_ENTRY

// WINDOWENTRY denotes a portion of an input vector based on a comparison
const WINDOWENTRY EntryType = C.GD_WINDOW_ENTRY

// INDEXENTRY denotes the field type of the implicit INDEX field
const INDEXENTRY EntryType = C.GD_INDEX_ENTRY

// ALLENTRIES denotes that all entry types should be counted/listed
const ALLENTRIES EntryType = 0

// ALIASENTRIES denotes that only aliases should be counted/listed
const ALIASENTRIES EntryType = C.GD_ALIAS_ENTRIES

// SCALARENTRIES denotes that only scalar fields should be counted/listed
// (That is, CONST, CARRAY, and STRING)
const SCALARENTRIES EntryType = C.GD_SCALAR_ENTRIES

// VECTORENTRIES denotes that only vector fields should be counted/listed
const VECTORENTRIES EntryType = C.GD_VECTOR_ENTRIES

// HIDDENENTRIES denotes that hidden entries should be counted/listed
const HIDDENENTRIES EntryType = C.GD_ENTRIES_HIDDEN

// NOALIASENTRIES denotes that alias fields should NOT be counted/listed
const NOALIASENTRIES EntryType = C.GD_ENTRIES_NOALIAS

// REGEXPCRE use the Perl-Compatible Regular Expression library instead of the POSIX
const REGEXPCRE EntryType = C.GD_REGEX_PCRE

// REGEXCASELESS do case-insensitive matching
const REGEXCASELESS EntryType = C.GD_REGEX_CASELESS

// REGEXICASE do case-insensitive matching (synonym of above)
const REGEXICASE EntryType = C.GD_REGEX_ICASE

// REGEXJAVASCRIPT (PCRE only): use Javascript-compatible reg exp grammar
const REGEXJAVASCRIPT EntryType = C.GD_REGEX_JAVASCRIPT

// REGEXUNICODE (PCRE only): use UTF-8
const REGEXUNICODE EntryType = C.GD_REGEX_UNICODE

// LOOKBACKALL searches backwards all the way to the start of a data source.
const LOOKBACKALL int = -1

// FRAMEHERE indicates the current location of the I/O pointer
const FRAMEHERE int = C.GD_HERE

///

// DeleteFlags are flags for the Dirfile.Delete function
type DeleteFlags uint64

// DELETEDATA indicates delete the binary file associated with the RAW field
const DELETEDATA DeleteFlags = C.GD_DEL_DATA

// DELETEDEREF indicates dereference a CONST field used as a field parameter
const DELETEDEREF DeleteFlags = C.GD_DEL_DEREF

// DELETEFORCE indicates delete the field even if it's an input to other fields
const DELETEFORCE DeleteFlags = C.GD_DEL_FORCE

// DELETEMETA indicates delete metafields attached to the field
const DELETEMETA DeleteFlags = C.GD_DEL_META

///

// WindowOps are operations used in a WINDOW field
type WindowOps uint64

// WINDOPEQ means check field equals threshold
const WINDOPEQ WindowOps = C.GD_WINDOP_EQ

// WINDOPNE means check field does not equal threshold
const WINDOPNE WindowOps = C.GD_WINDOP_NE

// WINDOPSET means at least one bit set in threshold is also set in the check field
const WINDOPSET WindowOps = C.GD_WINDOP_SET

// WINDOPCLR means at least one bit set in threshold is not set in the check field
const WINDOPCLR WindowOps = C.GD_WINDOP_CLR

// WINDOPGE means the check field is greater than or equal to threshold
const WINDOPGE WindowOps = C.GD_WINDOP_GE

// WINDOPGT means the check field is strictly greater than threshold
const WINDOPGT WindowOps = C.GD_WINDOP_GT

// WINDOPLE means the check field is less than or equal to threshold
const WINDOPLE WindowOps = C.GD_WINDOP_LE

// WINDOPLT means the check field is strictly less than threshold
const WINDOPLT WindowOps = C.GD_WINDOP_LT

// WINDOPUNK means an invalid value
const WINDOPUNK WindowOps = C.GD_WINDOP_UNK

///

// SeekFlags are flags for Dirfile.Seek
type SeekFlags uint

// SEEKSET means file position is relative to the beginning of field
const SEEKSET SeekFlags = C.GD_SEEK_SET

// SEEKCUR means file position is relative to the current position
const SEEKCUR SeekFlags = C.GD_SEEK_CUR

// SEEKEND means file position is relative to the end of field
const SEEKEND SeekFlags = C.GD_SEEK_END

// SEEKWRITE means the next operation on the field will be a write via PutData
const SEEKWRITE SeekFlags = C.GD_SEEK_WRITE

/// Encoding methods. See http://getdata.sourceforge.net/dirfile.html for
/// further information on the methods

// AUTOENCODED means encoding should be detected by getdata library.
const AUTOENCODED Flags = C.GD_AUTO_ENCODED

// UNENCODED means raw data are not encoded
const UNENCODED Flags = C.GD_UNENCODED

// TEXTENCODED means raw data are text encoded
const TEXTENCODED Flags = C.GD_TEXT_ENCODED

// SLIMENCODED means raw data are encoded by slimlib
const SLIMENCODED Flags = C.GD_SLIM_ENCODED

// GZIPENCODED means raw data are gzip encoded by zlib
const GZIPENCODED Flags = C.GD_GZIP_ENCODED

// BZIP2ENCODED means raw data are bzip2 encoded
const BZIP2ENCODED Flags = C.GD_BZIP2_ENCODED

// LZMAENCODED means raw data are lzma encoded
const LZMAENCODED Flags = C.GD_LZMA_ENCODED

// SIEENCODED means raw data are sample-index encoded
const SIEENCODED Flags = C.GD_SIE_ENCODED

// ZZIPENCODED means raw data are zzip encoded
const ZZIPENCODED Flags = C.GD_ZZIP_ENCODED

// ZZSLIMENCODED means raw data are encoded by a combination of zzip and slimlib
const ZZSLIMENCODED Flags = C.GD_ZZSLIM_ENCODED

// FLACENCODED means raw data are FLAC encoded
const FLACENCODED Flags = C.GD_FLAC_ENCODED

// RenameFlags are used in Entry.Move and Entry.Rename
type RenameFlags uint

// RENAMEDANGLE means don't update ALIAS entries, but turn them into dangling aliases
const RENAMEDANGLE RenameFlags = C.GD_REN_DANGLE

// RENAMEDATA if renaming a RAW field, also rename the data file on disk
const RENAMEDATA RenameFlags = C.GD_REN_DATA

// RENAMEFORCE means instead of having the call fail, just skip updating field codes which would become invalid
const RENAMEFORCE RenameFlags = C.GD_REN_FORCE

// RENAMEUPDATEDB means update references to the renamed field to use its new name
const RENAMEUPDATEDB RenameFlags = C.GD_REN_UPDB
