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

package cork

import (
	"io"
	"sync"
)

// Encoder represents a CORK encoder.
type Encoder struct {
	h *Handle
	w *Writer
	p bool
}

var encoders = sync.Pool{
	New: func() interface{} {
		return &Encoder{
			w: newWriter(),
			h: new(Handle),
			p: true,
		}
	},
}

// Encode encodes a Go object.
func Encode(src interface{}) (dst []byte) {
	enc := NewEncoderBytesFromPool(&dst)
	enc.Encode(src)
	enc.Reset()
	return
}

// EncodeInto encodes a Go object into a byte slice.
func EncodeInto(src interface{}, dst *[]byte) {
	enc := NewEncoderBytesFromPool(dst)
	enc.Encode(src)
	enc.Reset()
	return
}

// NewEncoder returns an Encoder for encoding into an io.Writer.
func NewEncoder(w io.Writer) *Encoder {
	e := &Encoder{w: newWriter(), h: new(Handle)}
	e.w.w.Reset(w)
	return e
}

// NewEncoderBytes returns an Encoder for encoding directly into a byte slice.
func NewEncoderBytes(b *[]byte) *Encoder {
	e := &Encoder{w: newWriter(), h: new(Handle)}
	e.w.w.ResetBytes(b)
	return e
}

// NewEncoderFromPool returns an Encoder for encoding into an
// io.Writer. The Encoder is taken from a pool of encoders, and
// must be put back when finished, using e.Reset().
func NewEncoderFromPool(w io.Writer) *Encoder {
	e := encoders.Get().(*Encoder)
	e.w.w.Reset(w)
	return e
}

// NewEncoderBytesFromPool returns an Encoder for encoding directly
// into a byte slice. The Encoder is taken from a pool of encoders,
// and must be put back when finished, using e.Reset().
func NewEncoderBytesFromPool(b *[]byte) *Encoder {
	e := encoders.Get().(*Encoder)
	e.w.w.ResetBytes(b)
	return e
}

// Reset flushes any remaing data and adds the Encoder back into
// the sync pool. If the Encoder was not originally from the
// sync pool, then the Encoder is discarded.
func (e *Encoder) Reset() {
	if e.p {
		encoders.Put(e)
	}
}

// Options sets the configuration options that the Encoder should use.
func (e *Encoder) Options(h *Handle) *Encoder {
	e.w.h, e.h = h, h
	return e
}

/*
Encode encodes the 'src' object into the stream.

The decoder can be configured using struct tags. The 'cork' key if found will be
analysed for any configuration options when encoding.

Each exported struct field is encoded unless:
	- the field's tag is "-"
	- the field is empty and its tag specifies the "omitempty" option.

When encoding a struct as a map, the first string in the tag (before the comma)
will be used for the map key, and if not specified will default to the struct key
name.

The empty values (for omitempty option) are false, 0, any nil pointer or
interface value, and any array, slice, map, or string of length zero.

	type Tester struct {
		Test bool   `cork:"-"`              // Skip this field
		Name string `cork:"name"`           // Use key "name" in encode stream
		Size int32  `cork:"size"`           // Use key "size" in encode stream
		Data []byte `cork:"data,omitempty"` // Use key data in encode stream, and omit if empty
	}

Example:

	// Encoding a typed value
	var s string = "Hello"
	buf := bytes.NewBuffer(nil)
	err = cork.NewEncoder(buf).Encode(s)

	// Encoding a struct
	var t &Tester{Name: "Temp", Size: 0}
	buf := bytes.NewBuffer(nil)
	err = cork.NewEncoder(buf).Encode(t)

*/
func (e *Encoder) Encode(src interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if catch, ok := r.(error); ok {
				err = catch
			}
		}
	}()
	e.w.EncodeAny(src)
	e.w.w.Flush()
	return
}
