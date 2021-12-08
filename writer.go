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
	"math"
	"time"

	"github.com/surrealdb/bump"
)

// Writer is used when self-encoding a cork.Selfer item into binary form.
type Writer struct {
	h *Handle
	w *bump.Writer
}

func newWriter() *Writer {
	return &Writer{
		w: bump.NewWriter(nil),
	}
}

func (w *Writer) writeOne(v byte) {
	if err := w.w.WriteByte(v); err != nil {
		panic(err)
	}
}

func (w *Writer) writeMany(v []byte) {
	if err := w.w.WriteBytes(v); err != nil {
		panic(err)
	}
}

func (w *Writer) writeText(v string) {
	if err := w.w.WriteString(v); err != nil {
		panic(err)
	}
}

// ---------------------------------------------------------------------------

func (w *Writer) writeLen(v uint) {
	switch {
	case v >= 0 && v <= fixedInt:
		w.writeOne(byte(v))
	case v <= math.MaxUint8:
		w.writeOne(cUint8)
		w.writeOne(byte(v))
	case v <= math.MaxUint16:
		w.writeOne(cUint16)
		w.writeOne(byte(v >> 8))
		w.writeOne(byte(v))
	case v <= math.MaxUint32:
		w.writeOne(cUint32)
		w.writeOne(byte(v >> 24))
		w.writeOne(byte(v >> 16))
		w.writeOne(byte(v >> 8))
		w.writeOne(byte(v))
	case v <= math.MaxUint64:
		w.writeOne(cUint64)
		w.writeOne(byte(v >> 56))
		w.writeOne(byte(v >> 48))
		w.writeOne(byte(v >> 40))
		w.writeOne(byte(v >> 32))
		w.writeOne(byte(v >> 24))
		w.writeOne(byte(v >> 16))
		w.writeOne(byte(v >> 8))
		w.writeOne(byte(v))
	}
}

func (w *Writer) writeLen8(val uint8) {
	w.writeOne(byte(val))
}

func (w *Writer) writeLen16(val uint16) {
	w.writeOne(byte(val >> 8))
	w.writeOne(byte(val))
}

func (w *Writer) writeLen32(val uint32) {
	w.writeOne(byte(val >> 24))
	w.writeOne(byte(val >> 16))
	w.writeOne(byte(val >> 8))
	w.writeOne(byte(val))
}

func (w *Writer) writeLen64(val uint64) {
	w.writeOne(byte(val >> 56))
	w.writeOne(byte(val >> 48))
	w.writeOne(byte(val >> 40))
	w.writeOne(byte(val >> 32))
	w.writeOne(byte(val >> 24))
	w.writeOne(byte(val >> 16))
	w.writeOne(byte(val >> 8))
	w.writeOne(byte(val))
}

// ---------------------------------------------------------------------------

// EncodeBool writes a nil value to the Writer.
func (w *Writer) EncodeNil() {
	w.writeOne(cNil)
}

// EncodeBool encodes a boolean value to the Writer.
func (w *Writer) EncodeBool(v bool) {
	if v {
		w.writeOne(cTrue)
	} else {
		w.writeOne(cFalse)
	}
}

// EncodeByte encodes a byte value to the Writer.
func (w *Writer) EncodeByte(v byte) {
	w.writeOne(v)
}

// EncodeBytes encodes a byte slice value to the Writer.
func (w *Writer) EncodeBytes(v []byte) {
	sze := len(v)
	switch {
	case sze <= fixedBin:
		w.writeOne(cFixBin + byte(sze))
	case sze <= math.MaxUint8:
		w.writeOne(cBin8)
		w.writeLen8(uint8(sze))
	case sze <= math.MaxUint16:
		w.writeOne(cBin16)
		w.writeLen16(uint16(sze))
	case sze <= math.MaxUint32:
		w.writeOne(cBin32)
		w.writeLen32(uint32(sze))
	case sze <= math.MaxInt64:
		w.writeOne(cBin64)
		w.writeLen64(uint64(sze))
	}
	w.writeMany(v)
}

// EncodeString encodes a string value to the Writer.
func (w *Writer) EncodeString(v string) {
	sze := len(v)
	switch {
	case sze <= fixedStr:
		w.writeOne(cFixStr + byte(sze))
	case sze <= math.MaxUint8:
		w.writeOne(cStr8)
		w.writeLen8(uint8(sze))
	case sze <= math.MaxUint16:
		w.writeOne(cStr16)
		w.writeLen16(uint16(sze))
	case sze <= math.MaxUint32:
		w.writeOne(cStr32)
		w.writeLen32(uint32(sze))
	case sze <= math.MaxInt64:
		w.writeOne(cStr64)
		w.writeLen64(uint64(sze))
	}
	w.writeText(v)
}

// ---------------------------------------------------------------------------

// EncodeInt encodes an int value to the Writer.
func (w *Writer) EncodeInt(v int) {
	switch {
	case v >= 0 && v <= fixedInt:
		w.writeOne(byte(v))
	case v >= math.MinInt8 && v <= math.MaxInt8:
		w.writeOne(cInt8)
		w.writeOne(byte(v))
	case v >= math.MinInt16 && v <= math.MaxInt16:
		w.writeOne(cInt16)
		w.writeOne(byte(v >> 8))
		w.writeOne(byte(v))
	case v >= math.MinInt32 && v <= math.MaxInt32:
		w.writeOne(cInt32)
		w.writeOne(byte(v >> 24))
		w.writeOne(byte(v >> 16))
		w.writeOne(byte(v >> 8))
		w.writeOne(byte(v))
	case v >= math.MinInt64 && v <= math.MaxInt64:
		w.writeOne(cInt64)
		w.writeOne(byte(v >> 56))
		w.writeOne(byte(v >> 48))
		w.writeOne(byte(v >> 40))
		w.writeOne(byte(v >> 32))
		w.writeOne(byte(v >> 24))
		w.writeOne(byte(v >> 16))
		w.writeOne(byte(v >> 8))
		w.writeOne(byte(v))
	}
}

// EncodeInt8 encodes an int8 value to the Writer.
func (w *Writer) EncodeInt8(v int8) {
	w.EncodeInt(int(v))
}

// EncodeInt16 encodes an int16 value to the Writer.
func (w *Writer) EncodeInt16(v int16) {
	w.EncodeInt(int(v))
}

// EncodeInt32 encodes an int32 value to the Writer.
func (w *Writer) EncodeInt32(v int32) {
	w.EncodeInt(int(v))
}

// EncodeInt64 encodes an int64 value to the Writer.
func (w *Writer) EncodeInt64(v int64) {
	w.EncodeInt(int(v))
}

// ---------------------------------------------------------------------------

// EncodeUint encodes a uint value to the Writer.
func (w *Writer) EncodeUint(v uint) {
	switch {
	case v >= 0 && v <= fixedInt:
		w.writeOne(byte(v))
	case v <= math.MaxUint8:
		w.writeOne(cUint8)
		w.writeOne(byte(v))
	case v <= math.MaxUint16:
		w.writeOne(cUint16)
		w.writeOne(byte(v >> 8))
		w.writeOne(byte(v))
	case v <= math.MaxUint32:
		w.writeOne(cUint32)
		w.writeOne(byte(v >> 24))
		w.writeOne(byte(v >> 16))
		w.writeOne(byte(v >> 8))
		w.writeOne(byte(v))
	case v <= math.MaxUint64:
		w.writeOne(cUint64)
		w.writeOne(byte(v >> 56))
		w.writeOne(byte(v >> 48))
		w.writeOne(byte(v >> 40))
		w.writeOne(byte(v >> 32))
		w.writeOne(byte(v >> 24))
		w.writeOne(byte(v >> 16))
		w.writeOne(byte(v >> 8))
		w.writeOne(byte(v))
	}
}

// EncodeUint8 encodes a uint8 value to the Writer.
func (w *Writer) EncodeUint8(v uint8) {
	w.EncodeUint(uint(v))
}

// EncodeUint16 encodes a uint16 value to the Writer.
func (w *Writer) EncodeUint16(v uint16) {
	w.EncodeUint(uint(v))
}

// EncodeUint32 encodes a uint32 value to the Writer.
func (w *Writer) EncodeUint32(v uint32) {
	w.EncodeUint(uint(v))
}

// EncodeUint64 encodes a uint64 value to the Writer.
func (w *Writer) EncodeUint64(v uint64) {
	w.EncodeUint(uint(v))
}

// ---------------------------------------------------------------------------

// EncodeFloat32 encodes a float32 value to the Writer.
func (w *Writer) EncodeFloat32(v float32) {
	tmp := math.Float32bits(v)
	w.writeOne(cFloat32)
	w.writeOne(byte(tmp >> 24))
	w.writeOne(byte(tmp >> 16))
	w.writeOne(byte(tmp >> 8))
	w.writeOne(byte(tmp))
}

// EncodeFloat64 encodes a float64 value to the Writer.
func (w *Writer) EncodeFloat64(v float64) {
	tmp := math.Float64bits(v)
	w.writeOne(cFloat64)
	w.writeOne(byte(tmp >> 56))
	w.writeOne(byte(tmp >> 48))
	w.writeOne(byte(tmp >> 40))
	w.writeOne(byte(tmp >> 32))
	w.writeOne(byte(tmp >> 24))
	w.writeOne(byte(tmp >> 16))
	w.writeOne(byte(tmp >> 8))
	w.writeOne(byte(tmp))
}

// ---------------------------------------------------------------------------

// EncodeComplex64 encodes a complex64 value to the Writer.
func (w *Writer) EncodeComplex64(v complex64) {
	one := math.Float32bits(real(v))
	two := math.Float32bits(imag(v))
	w.writeOne(cComplex64)
	w.writeOne(byte(one >> 24))
	w.writeOne(byte(one >> 16))
	w.writeOne(byte(one >> 8))
	w.writeOne(byte(one))
	w.writeOne(byte(two >> 24))
	w.writeOne(byte(two >> 16))
	w.writeOne(byte(two >> 8))
	w.writeOne(byte(two))
}

// EncodeComplex128 encodes a complex128 value to the Writer.
func (w *Writer) EncodeComplex128(v complex128) {
	one := math.Float64bits(real(v))
	two := math.Float64bits(imag(v))
	w.writeOne(cComplex128)
	w.writeOne(byte(one >> 56))
	w.writeOne(byte(one >> 48))
	w.writeOne(byte(one >> 40))
	w.writeOne(byte(one >> 32))
	w.writeOne(byte(one >> 24))
	w.writeOne(byte(one >> 16))
	w.writeOne(byte(one >> 8))
	w.writeOne(byte(one))
	w.writeOne(byte(two >> 56))
	w.writeOne(byte(two >> 48))
	w.writeOne(byte(two >> 40))
	w.writeOne(byte(two >> 32))
	w.writeOne(byte(two >> 24))
	w.writeOne(byte(two >> 16))
	w.writeOne(byte(two >> 8))
	w.writeOne(byte(two))
}

// ---------------------------------------------------------------------------

// EncodeInt encodes a time.Time value to the Writer.
func (w *Writer) EncodeTime(v time.Time) {
	tmp := uint64(v.UTC().UnixNano())
	w.writeOne(cTime)
	w.writeOne(byte(tmp >> 56))
	w.writeOne(byte(tmp >> 48))
	w.writeOne(byte(tmp >> 40))
	w.writeOne(byte(tmp >> 32))
	w.writeOne(byte(tmp >> 24))
	w.writeOne(byte(tmp >> 16))
	w.writeOne(byte(tmp >> 8))
	w.writeOne(byte(tmp))
}

// ---------------------------------------------------------------------------

// EncodeArr encodes an array to the Writer.
func (w *Writer) EncodeArr(v interface{}) {
	switch a := v.(type) {
	case []bool:
		w.encodeArrBool(a)
	case []int:
		w.encodeArrInt(a)
	case []int8:
		w.encodeArrInt8(a)
	case []int16:
		w.encodeArrInt16(a)
	case []int32:
		w.encodeArrInt32(a)
	case []int64:
		w.encodeArrInt64(a)
	case []uint:
		w.encodeArrUint(a)
	case []uint8:
		w.encodeArrUint8(a)
	case []uint16:
		w.encodeArrUint16(a)
	case []uint32:
		w.encodeArrUint32(a)
	case []uint64:
		w.encodeArrUint64(a)
	case []string:
		w.encodeArrString(a)
	case []float32:
		w.encodeArrFloat32(a)
	case []float64:
		w.encodeArrFloat64(a)
	case []complex64:
		w.encodeArrComplex64(a)
	case []complex128:
		w.encodeArrComplex128(a)
	case []time.Time:
		w.encodeArrTime(a)
	case []interface{}:
		w.encodeArrAny(a)
	default:
		w.EncodeAny(v)
	}
}

func (w *Writer) encodeArrLen(v int) {
	switch {
	case v >= 0 && v <= fixedArr:
		w.writeOne(cFixArr + byte(v))
	default:
		w.writeOne(cArr)
		w.writeLen(uint(v))
	}
}

// EncodeMap encodes a map to the Writer.
func (w *Writer) EncodeMap(v interface{}) {
	switch m := v.(type) {
	case map[string]int:
		w.encodeMapStringInt(m)
	case map[string]uint:
		w.encodeMapStringUint(m)
	case map[string]bool:
		w.encodeMapStringBool(m)
	case map[string]string:
		w.encodeMapStringString(m)
	case map[int]interface{}:
		w.encodeMapIntAny(m)
	case map[uint]interface{}:
		w.encodeMapUintAny(m)
	case map[string]interface{}:
		w.encodeMapStringAny(m)
	case map[time.Time]interface{}:
		w.encodeMapTimeAny(m)
	case map[interface{}]interface{}:
		w.encodeMapAnyAny(m)
	default:
		w.EncodeAny(v)
	}
}

func (w *Writer) encodeMapLen(v int) {
	switch {
	case v >= 0 && v <= fixedMap:
		w.writeOne(cFixMap + byte(v))
	default:
		w.writeOne(cMap)
		w.writeLen(uint(v))
	}
}
