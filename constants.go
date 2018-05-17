package getdata

/*
#cgo CFLAGS: -I/usr/local/include
#cgo LDFLAGS: -L/usr/local/lib -lgetdata
#include <getdata.h>
#include <stdlib.h>
*/
import "C"

// Flags are dirfile opening flags, including encoding methods
type Flags uint64

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

// UNENCODED means data are raw binary, not compressed
const UNENCODED Flags = C.GD_UNENCODED

// **

// EntryType signifies the field type given for entries in the FORMAT files
type EntryType uint64

// NOENTRY denotes an invalid entry type
const NOENTRY = C.GD_NO_ENTRY
