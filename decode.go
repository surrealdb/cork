// Copyright © 2016 Abcum Ltd
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

// Decoder represents a CORK decoder.
type Decoder struct {
	h *Handle
	r *Reader
	p bool
}

var decoders = sync.Pool{
	New: func() interface{} {
		return &Decoder{
			r: newReader(),
			h: new(Handle),
			p: true,
		}
	},
}

// Decode decodes binary data.
func Decode(src []byte) (dst interface{}) {
	dec := NewDecoderBytesFromPool(src)
	dec.Decode(&dst)
	dec.Reset()
	return
}

// DecodeInto decodes a byte slice into a Go object.
func DecodeInto(src []byte, dst interface{}) {
	dec := NewDecoderBytesFromPool(src)
	dec.Decode(dst)
	dec.Reset()
	return
}

// NewDecoder returns a Decoder for decoding from an io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	d := &Decoder{r: newReader(), h: new(Handle)}
	d.r.r.Reset(r)
	return d
}

// NewDecoderBytes returns a Decoder for decoding directly from a byte slice.
func NewDecoderBytes(b []byte) *Decoder {
	d := &Decoder{r: newReader(), h: new(Handle)}
	d.r.r.ResetBytes(b)
	return d
}

// NewDecoderFromPool returns a Decoder for decoding into an
// io.Reader. The Decoder is taken from a pool of decoders, and
// must be put back when finished, using d.Reset().
func NewDecoderFromPool(r io.Reader) *Decoder {
	d := decoders.Get().(*Decoder)
	d.r.r.Reset(r)
	return d
}

// NewDecoderBytesFromPool returns a Decoder for decoding directly
// from a byte slice. The Decoder is taken from a pool of decoders,
// and must be put back when finished, using d.Reset().
func NewDecoderBytesFromPool(b []byte) *Decoder {
	d := decoders.Get().(*Decoder)
	d.r.r.ResetBytes(b)
	return d
}

// Reset flushes adds the Decoder back into the sync pool. If the
// Decoder was not originally from the sync pool, then the
// Decoder is discarded.
func (d *Decoder) Reset() {
	if d.p {
		decoders.Put(d)
	}
}

// Options sets the configuration options that the Decoder should use.
func (d *Decoder) Options(h *Handle) *Decoder {
	d.r.h, d.h = h, h
	return d
}

/*
Decode decodes the stream into the 'dst' object.

The decoder can not decode into a nil pointer, but can decode into a nil interface.
If you do not know what type of stream it is, pass in a pointer to a nil interface.
We will decode and store a value in that nil interface.

When decoding into a nil interface{}, we will decode into an appropriate value based
on the contents of the stream. When decoding into a non-nil interface{} value, the
mode of encoding is based on the type of the value.

Example:

	// Decoding into a non-nil typed value
	var s string
	buf := bytes.NewBuffer(src)
	err = cork.NewDecoder(buf).Decode(&s)

	// Decoding into nil interface
	var v interface{}
	buf := bytes.NewBuffer(src)
	err := cork.NewDecoder(buf).Decode(&v)

*/
func (d *Decoder) Decode(dst interface{}) (err error) {
	defer func() {
		if r := recover(); r != nil {
			if catch, ok := r.(error); ok {
				err = catch
			}
		}
	}()
	d.r.DecodeAny(dst)
	return
}
