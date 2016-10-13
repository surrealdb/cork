# CORK

CORK is a dynamic and static binary object serialization specification
similar to GOB and MessagePack. CORK manages encoding and decoding binary
data streams, for use in transmission of, or storage of data.

### Basics

A CORK data stream is self-describing. Each data item in the stream is preceded
by a description of its type, expressed in terms of a small set of predefined
types. Pointers to values are not transmitted, but the contents are flattened
and transmitted.

To use gobs, create an Encoder and present it with a series of data items as
values or addresses that can be dereferenced to values. The Encoder makes sure
all type information is sent before it is needed. At the receive side, a Decoder 
retrieves values from the encoded stream and unpacks them into local variables.

### Types

CORK has built in support for the following types:

`nil`, `bool`, `string`, `[]byte`, `int8`, `int16`, `int32`, `int64`, `uint8`, 
`uint16`, `uint32`, `uint64`, `float32`, `float64`, `complex64`, `complex128`, 
`time.Time`, `[]<T>`, `map[<T>]<T>`

### Structs

When a struct is encountered whilst encoding (and that struct does not satisfy
the Corker interface) then the struct will be encoded into the stream as a map
with the keys encoded as strings, and the values as the relevant type. Any
struct tags describing how the struct should be encoded will be used.

### Corkers

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

### Encoding types

| Encoding name | Byte range        | Usage 
| :------------ | :---------------- | :---------------------------------------------- |
| fixint        | 0x00 ... 0x7F     | A positive number from `0` to `127`
| fixstr        | 0x80 ... 0x9F     | A string which is `0` to `(1<<5)-1` in length
| fixbin        | 0xA0 ... 0xAF     | Binary data which is `0` to `(1<<4)-1` in length
| fixext        | 0xB0 ... 0xBF     | A custom type which is `0` to `(1<<4)-1` in length
| fixarr        | 0xC0 ... 0xCF     | A set whose length is between `0` to `(1<<4)-1`
| fixmap        | 0xD0 ... 0xDF     | A map whose length is between `0` to `(1<<4)-1`
| -             | -                 | -
| nil           | 0xE0              | A `nil` value
| true          | 0xE1              | Boolean value which is `true`
| false         | 0xE2              | Boolean value which is `false`
| time          | 0xE3              | Timestamp with nanosecond precision
| -             | -                 | -
| str8          | 0xE4              | A string up to `(1<<8)-1` in length
| str16         | 0xE5              | A string up to `(1<<16)-1` in length
| str32         | 0xE6              | A string up to `(1<<32)-1` in length
| str64         | 0xE7              | A string up to `(1<<63)-1` in length
| -             | -                 | -
| bin8          | 0xE8              | Binary data up to `(1<<8)-1` in length
| bin16         | 0xE9              | Binary data up to `(1<<16)-1` in length
| bin32         | 0xEA              | Binary data up to `(1<<32)-1` in length
| bin64         | 0xEB              | Binary data up to `(1<<63)-1` in length
| -             | -                 | -
| ext8          | 0xEC              | A custom type up to `(1<<8)-1` in length
| ext16         | 0xED              | A custom type up to `(1<<16)-1` in length
| ext32         | 0xEE              | A custom type up to `(1<<32)-1` in length
| ext64         | 0xEF              | A custom type up to `(1<<63)-1` in length
| -             | -                 | 
| int8          | 0xF0              | A signed integer less than `(1<<7)-1`
| int16         | 0xF1              | A signed integer less than `(1<<15)-1`
| int32         | 0xF2              | A signed integer less than `(1<<31)-1`
| int64         | 0xF3              | A signed integer less than `(1<<63)-1`
| -             | -                 | -
| uint8         | 0xF4              | An unsigned integer less than `(1<<8)-1`
| uint16        | 0xF5              | An unsigned integer less than `(1<<16)-1`
| uint32        | 0xF6              | An unsigned integer less than `(1<<32)-1`
| uint64        | 0xF7              | An unsigned integer less than `(1<<64)-1`
| -             | -                 | -
| float32       | 0xF8              | A IEEE-754 32bit floating point number
| float64       | 0xF9              | A IEEE-754 64bit floating point number
| -             | -                 | -
| complex64     | 0xFA              | A 64bit complex number
| complex128    | 0xFB              | A 128bit complex number
| -             | -                 | -
| arr           | 0xFC              | A set whose length is greater than `(1<<4)-1`
| map           | 0xFD              | A map whose length is greater than `(1<<4)-1`
| sym           | 0xFE              | *Reserved for internal use*
| alt           | 0xFF              | *Reserved for internal use*

### Encoding methods

##### nil

A `nil` value is stored in `1` byte:

	+--------+
	|  0xE0  |
	+--------+

##### bool

A `bool` value is stored in `1` byte:

	true:
	+--------+
	|  0xE1  |
	+--------+

	false:
	+--------+
	|  0xE2  |
	+--------+

##### time

A `time` value is stored in `9` bytes:

	time:
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+
	|  0xE3  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+

##### numbers

A `number` value is stored in `1`, `2`, `3`, `5`, or `9` bytes:

	fixint stores a positive integer from 0 to 127:
	+--------+
	|  0x??  |
	+--------+

	int8 stores a signed integer upto (1<<7)-1:
	+--------+--------+
	|  0xF1  |  ----  |
	+--------+--------+

	int16 stores a signed integer upto (1<<15)-1:
	+--------+--------+--------+
	|  0xF2  |  ----  |  ----  |
	+--------+--------+--------+

	int32 stores a signed integer upto (1<<31)-1:
	+--------+--------+--------+--------+--------+
	|  0xF3  |  ----  |  ----  |  ----  |  ----  |
	+--------+--------+--------+--------+--------+

	int64 stores a signed integer upto (1<<63)-1:
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+
	|  0xF4  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+

	uint8 stores a signed integer upto (1<<8)-1:
	+--------+--------+
	|  0xF6  |  ----  |
	+--------+--------+

	uint16 stores a signed integer upto (1<<16)-1:
	+--------+--------+--------+
	|  0xF7  |  ----  |  ----  |
	+--------+--------+--------+

	uint32 stores a signed integer upto (1<<32)-1:
	+--------+--------+--------+--------+--------+
	|  0xF8  |  ----  |  ----  |  ----  |  ----  |
	+--------+--------+--------+--------+--------+

	uint64 stores a signed integer upto (1<<64)-1:
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+
	|  0xF9  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+

##### floats

A `float` value is stored in `5`, or `9` bytes:

	float32 stores a floating point number in IEEE 754 single precision format:
	+--------+--------+--------+--------+--------+
	|  0xFA  |  ----  |  ----  |  ----  |  ----  |
	+--------+--------+--------+--------+--------+

	float64 stores a floating point number in IEEE 754 double precision format:
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+
	|  0xFB  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+

##### strings

A `string` value is stored in `1`, `2`, `3`, `5`, or `9` descriptive bytes in addition to the size of the string data:

	fixstr stores a string whose length is upto 31 bytes:
	+--------+========+
	|  0x??  |  data  |
	+--------+========+

	str8 stores a string whose length is upto (1<<8)-1 bytes:
	+--------+--------+========+
	|  0xE4  |  ----  |  data  |
	+--------+--------+========+

	str16 stores a string whose length is upto (1<<16)-1 bytes:
	+--------+--------+--------+========+
	|  0xE5  |  ----  |  ----  |  data  |
	+--------+--------+--------+========+

	str32 stores a string whose length is upto (1<<32)-1 bytes:
	+--------+--------+--------+--------+--------+========+
	|  0xE6  |  ----  |  ----  |  ----  |  ----  |  data  |
	+--------+--------+--------+--------+--------+========+

	str64 stores a string whose length is upto (1<<64)-1 bytes:
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+========+
	|  0xE7  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  data  |
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+========+

##### binary

A `binary` value is stored in `1`, `2`, `3`, `5`, or `9` descriptive bytes in addition to the size of the binary data:

	fixbin stores a byte array whose length is upto 15 bytes:
	+--------+========+
	|  0xA?  |  data  |
	+--------+========+

	bin8 stores a byte array whose length is upto (1<<8)-1 bytes:
	+--------+--------+========+
	|  0xE8  |  ----  |  data  |
	+--------+--------+========+

	bin16 stores a byte array whose length is upto (1<<16)-1 bytes:
	+--------+--------+--------+========+
	|  0xE9  |  ----  |  ----  |  data  |
	+--------+--------+--------+========+

	bin32 stores a byte array whose length is upto (1<<32)-1 bytes:
	+--------+--------+--------+--------+--------+========+
	|  0xEA  |  ----  |  ----  |  ----  |  ----  |  data  |
	+--------+--------+--------+--------+--------+========+

	bin64 stores a byte array whose length is upto (1<<64)-1 bytes:
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+========+
	|  0xEB  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  data  |
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+========+

##### custom

A `custom` value is stored in `1`, `2`, `3`, `5`, or `9` descriptive bytes in addition to the size of the custom data:

	fixext stores a custom type whose length is upto 15 bytes:
	+--------+--------+========+
	|  0xB?  |  type  |  data  |
	+--------+--------+========+

	ext8 stores a custom type whose length is upto (1<<8)-1 bytes:
	+--------+--------+--------+========+
	|  0xEC  |  type  |  ----  |  data  |
	+--------+--------+--------+========+

	ext16 stores a custom type whose length is upto (1<<16)-1 bytes:
	+--------+--------+--------+--------+========+
	|  0xED  |  type  |  ----  |  ----  |  data  |
	+--------+--------+--------+--------+========+

	ext32 stores a custom type whose length is upto (1<<32)-1 bytes:
	+--------+--------+--------+--------+--------+--------+========+
	|  0xEE  |  type  |  ----  |  ----  |  ----  |  ----  |  data  |
	+--------+--------+--------+--------+--------+--------+========+

	ext64 stores a custom type whose length is upto (1<<64)-1 bytes:
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+--------+========+
	|  0xEF  |  type  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  ----  |  data  |
	+--------+--------+--------+--------+--------+--------+--------+--------+--------+--------+========+

##### set

A `set` value is stored in `1`, or `2` descriptive bytes in addition to the set elements:

	fixarr stores a set whose length is upto 15 elements:
	+--------+ - - - - - - - -+
	|  0xC?  |    Elements    |
	+--------+ - - - - - - - -+

	arr stores a set whose length is upto (1<<64)-1 elements:
	+--------+ - - - - - - - -+ - - - - - - - -+
	|  0xFC  |     Length     |    Elements    |
	+--------+ - - - - - - - -+ - - - - - - - -+

##### map

A `map` value is stored in `1`, or `2` descriptive bytes in addition to the map key-value pairs:

	fixmap stores a set whose length is upto 15 elements:
	+--------+ - - - - - - - -+
	|  0xD?  |    Elements    |
	+--------+ - - - - - - - -+

	map stores a set whose length is upto (1<<64)-1 key-value pairs:
	+--------+ - - - - - - - -+ - - - - - - - -+
	|  0xFD  |     Length     |    Elements    |
	+--------+ - - - - - - - -+ - - - - - - - -+
