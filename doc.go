// Copyright © SurrealDB Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

/*

CORK is a dynamic and static binary object serialization specification
similar to GOB and MessagePack. CORK manages encoding and decoding binary
data streams, for use in transmission of, or storage of data.

Basics

A CORK data stream is self-describing. Each data item in the stream is preceded
by a description of its type, expressed in terms of a small set of predefined
types. Pointers to values are not transmitted, but the contents are flattened
and transmitted.

To use cork, create an Encoder and present it with a series of data items as
values or addresses that can be dereferenced to values. The Encoder makes sure
all type information is sent before it is needed. At the receive side, a
Decoder retrieves values from the encoded stream and unpacks them into local
variables.

Types

CORK has built in support for the built-in Golang types

	nil
	bool
	string
	[]byte
	int
	int8
	int16
	int32
	int64
	uint
	uint8
	uint16
	uint32
	uint64
	float32
	float64
	complex64
	complex128
	time.Time
	interface{}
	[]<T>
	map[<T>]<T>

Structs

When a struct is encountered whilst encoding (and that struct does not satisfy
the Corker interface) then the struct will be encoded into the stream as a map
with the keys encoded as strings, and the values as the relevant type. Any
struct tags describing how the struct should be encoded will be used.

Corkers

CORK allows applications to define application-specific types to be added
to the encoding format. Each extended type must be assigned a unique byte
from 0x00 upto 0xFF. Application-specific types (otherwise known as Corkers)
are able to encode themselves into a binary data value, and are able to
decode themselves from that same binary data value.

To define a custom type, an application must ensure that the type satisfies
the Corker interface, and must then register the type using the Register
method.

If a custom type is found in the stream when decoding, but no type with the
specified unique byte is registered, then the binary data value will be
decoded as a raw binary data value.

Types and Values

The source and destination values/types need not correspond exactly.  For structs,
fields (identified by name) that are in the source but absent from the receiving
variable will be ignored.  Fields that are in the receiving variable but missing
from the transmitted type or value will be ignored in the destination.  If a field
with the same name is present in both, their types must be compatible. Both the
receiver and transmitter will do all necessary indirection and dereferencing to
convert between cork encoded data and actual Go values. For instance, a cork type
that is schematically,

	struct { A, B int }

can be sent from or received into any of these Go types:

	struct { A, B int }	// the same
	*struct { A, B int }	// extra indirection of the struct
	struct { *A, **B int }	// extra indirection of the fields
	struct { A, B int64 }	// different concrete value type; see below

It may also be received into any of these:

	struct { A, B int }	// the same
	struct { B, A int }	// ordering doesn't matter; matching is by name
	struct { A, B, C int }	// extra field (C) ignored
	struct { B int }	// missing field (A) ignored; data will be dropped
	struct { B, C int }	// missing field (A) ignored; extra field (C) ignored.

Attempting to receive into these types will draw a decode error:

	struct { A int; B uint }	// change of signedness for B
	struct { A int; B float }	// change of type for B
	struct { }			// no field names in common
	struct { C, D int }		// no field names in common

Integers can be serialized using two different methods: arbritrary precision or
full precision. When using arbritrary precision, integers (int, int8, int16, int32,
int64) and unsinged integers (uint, uint8, uint16, uint32, uint64) are encoded
with a variable-length encoding using as few bytes as possible, and are decoded
into the destination variable, or an integer with the necessary capacity (when
decoding into a nil interface). When using full precision, all integers (int8,
int16, int32, int64) and unsinged integers (uint8, uint16, uint32, uint64) are
encoded with a fixed-length encoding format, and are able to be decoded into the
corresponding variable type when decoding into a nil interface.

Signed integers may be received into any signed integer variable: int, int16, etc.;
unsigned integers may be received into any unsigned integer variable; and floating
point values may be received into any floating point variable.  However,
the destination variable must be able to represent the value or the decode
operation will fail.

Structs, arrays and slices are also supported. Structs encode and decode only
exported fields. Struct tags (using the 'cork' descriptor) can specify custom key
names to be used when serailizing, and can be omitted entirely, or when empty using
the 'omitempty' tag keyword. Strings and arrays of bytes are supported with a special,
efficient representation. When a slice is decoded, if the existing slice has
capacity the slice will be extended in place; if not, a new array is allocated.
Regardless, the length of the resulting slice reports the number of elements decoded.

In general, if allocation is required, the decoder will allocate memory. If not,
it will update the destination variables with values read from the stream. It does
not initialize them first, so if the destination is a compound value such as a
map, struct, or slice, the decoded values will be merged elementwise into the
existing variables.

Functions and channels will not be encoded into a CORK. Attempting to encode
such a value at the top level will fail. A struct field of chan or func type
is treated exactly like an unexported field and is ignored.

Specification

The full specification can be found at http://github.com/surrealdb/cork/SPEC.md

*/
package cork
