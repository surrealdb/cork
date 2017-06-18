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
	"encoding/binary"
	"math"
	"reflect"
	"time"
	"unsafe"

	"github.com/abcum/bump"
)

// Reader is used when self-decoding a cork.Selfer item from binary form.
type Reader struct {
	h *Handle
	r *bump.Reader
}

func newReader() *Reader {
	return &Reader{
		r: bump.NewReader(nil),
	}
}

func (r *Reader) peekOne() (val byte) {
	val, err := r.r.PeekByte()
	if err != nil {
		panic(err)
	}
	return val
}

func (r *Reader) readOne() (val byte) {
	val, err := r.r.ReadByte()
	if err != nil {
		panic(err)
	}
	return val
}

func (r *Reader) readMany(l int) (val []byte) {
	val, err := r.r.ReadBytes(l)
	if err != nil {
		panic(err)
	}
	return val
}

func (r *Reader) readText(l int) (val string) {
	b := r.readMany(l)
	return *(*string)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&b))))
}

// ---------------------------------------------------------------------------

func (r *Reader) readLen() int {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		return int(b)
	case b == cInt8:
		return int(r.readOne())
	case b == cInt16:
		return int(binary.BigEndian.Uint16(r.readMany(2)))
	case b == cInt32:
		return int(binary.BigEndian.Uint32(r.readMany(4)))
	case b == cInt64:
		return int(binary.BigEndian.Uint64(r.readMany(8)))
	default:
		panic(fail)
	}
}

func (r *Reader) readLen8() int {
	return int(r.readOne())
}

func (r *Reader) readLen16() int {
	return int(binary.BigEndian.Uint16(r.readMany(2)))
}

func (r *Reader) readLen32() int {
	return int(binary.BigEndian.Uint32(r.readMany(4)))
}

func (r *Reader) readLen64() int {
	return int(binary.BigEndian.Uint64(r.readMany(8)))
}

// ---------------------------------------------------------------------------

// DecodeBool decodes a boolean value from the Reader.
func (r *Reader) DecodeBool(v *bool) {
	switch r.readOne() {
	case cTrue:
		*v = true
	case cFalse:
		*v = false
	default:
		panic(fail)
	}
}

// DecodeByte decodes a byte value from the Reader.
func (r *Reader) DecodeByte(v *byte) {
	*v = r.readOne()
}

// DecodeBytes decodes a byte slice value from the Reader.
func (r *Reader) DecodeBytes(v *[]byte) {
	b := r.readOne()
	switch {
	case b >= cFixBin && b <= cFixBin+fixedBin:
		*v = r.readMany(int(b - cFixBin))
	case b == cBin8:
		*v = r.readMany(int(r.readLen8()))
	case b == cBin16:
		*v = r.readMany(int(r.readLen16()))
	case b == cBin32:
		*v = r.readMany(int(r.readLen32()))
	case b == cBin64:
		*v = r.readMany(int(r.readLen64()))
	default:
		panic(fail)
	}
}

// DecodeString decodes a string value from the Reader.
func (r *Reader) DecodeString(v *string) {
	b := r.readOne()
	switch {
	case b >= cFixStr && b <= cFixStr+fixedStr:
		*v = r.readText(int(b - cFixStr))
	case b == cStr8:
		*v = r.readText(int(r.readLen8()))
	case b == cStr16:
		*v = r.readText(int(r.readLen16()))
	case b == cStr32:
		*v = r.readText(int(r.readLen32()))
	case b == cStr64:
		*v = r.readText(int(r.readLen64()))
	default:
		panic(fail)
	}
}

// ---------------------------------------------------------------------------

// DecodeInt decodes an int value from the Reader.
func (r *Reader) DecodeInt(v *int) {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		*v = int(b)
	case b == cInt8:
		*v = int(int8(r.readOne()))
	case b == cInt16:
		*v = int(int16(binary.BigEndian.Uint16(r.readMany(2))))
	case b == cInt32:
		*v = int(int32(binary.BigEndian.Uint32(r.readMany(4))))
	case b == cInt64:
		*v = int(int64(binary.BigEndian.Uint64(r.readMany(8))))
	default:
		panic(fail)
	}
}

// DecodeInt8 decodes an int8 value from the Reader.
func (r *Reader) DecodeInt8(v *int8) {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		*v = int8(b)
	case b == cInt8:
		*v = int8(r.readOne())
	default:
		panic(fail)
	}
}

// DecodeInt16 decodes an int16 value from the Reader.
func (r *Reader) DecodeInt16(v *int16) {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		*v = int16(b)
	case b == cInt8:
		*v = int16(int8(r.readOne()))
	case b == cInt16:
		*v = int16(binary.BigEndian.Uint16(r.readMany(2)))
	default:
		panic(fail)
	}
}

// DecodeInt32 decodes an int32 value from the Reader.
func (r *Reader) DecodeInt32(v *int32) {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		*v = int32(b)
	case b == cInt8:
		*v = int32(int8(r.readOne()))
	case b == cInt16:
		*v = int32(int16(binary.BigEndian.Uint16(r.readMany(2))))
	case b == cInt32:
		*v = int32(binary.BigEndian.Uint32(r.readMany(4)))
	default:
		panic(fail)
	}
}

// DecodeInt64 decodes an int64 value from the Reader.
func (r *Reader) DecodeInt64(v *int64) {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		*v = int64(b)
	case b == cInt8:
		*v = int64(int8(r.readOne()))
	case b == cInt16:
		*v = int64(int16(binary.BigEndian.Uint16(r.readMany(2))))
	case b == cInt32:
		*v = int64(int32(binary.BigEndian.Uint32(r.readMany(4))))
	case b == cInt64:
		*v = int64(binary.BigEndian.Uint64(r.readMany(8)))
	default:
		panic(fail)
	}
}

// ---------------------------------------------------------------------------

// DecodeUint decodes a uint value from the Reader.
func (r *Reader) DecodeUint(v *uint) {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		*v = uint(b)
	case b == cUint8:
		*v = uint(r.readOne())
	case b == cUint16:
		*v = uint(binary.BigEndian.Uint16(r.readMany(2)))
	case b == cUint32:
		*v = uint(binary.BigEndian.Uint32(r.readMany(4)))
	case b == cUint64:
		*v = uint(binary.BigEndian.Uint64(r.readMany(8)))
	default:
		panic(fail)
	}
}

// DecodeUint8 decodes a uint8 value from the Reader.
func (r *Reader) DecodeUint8(v *uint8) {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		*v = uint8(b)
	case b == cUint8:
		*v = uint8(r.readOne())
	default:
		panic(fail)
	}
}

// DecodeUint16 decodes a uint16 value from the Reader.
func (r *Reader) DecodeUint16(v *uint16) {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		*v = uint16(b)
	case b == cUint8:
		*v = uint16(r.readOne())
	case b == cUint16:
		*v = uint16(binary.BigEndian.Uint16(r.readMany(2)))
	default:
		panic(fail)
	}
}

// DecodeUint32 decodes a uint32 value from the Reader.
func (r *Reader) DecodeUint32(v *uint32) {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		*v = uint32(b)
	case b == cUint8:
		*v = uint32(r.readOne())
	case b == cUint16:
		*v = uint32(binary.BigEndian.Uint16(r.readMany(2)))
	case b == cUint32:
		*v = uint32(binary.BigEndian.Uint32(r.readMany(4)))
	default:
		panic(fail)
	}
}

// DecodeUint64 decodes a uint64 value from the Reader.
func (r *Reader) DecodeUint64(v *uint64) {
	b := r.readOne()
	switch {
	case b >= cFixInt && b <= cFixInt+fixedInt:
		*v = uint64(b)
	case b == cUint8:
		*v = uint64(r.readOne())
	case b == cUint16:
		*v = uint64(binary.BigEndian.Uint16(r.readMany(2)))
	case b == cUint32:
		*v = uint64(binary.BigEndian.Uint32(r.readMany(4)))
	case b == cUint64:
		*v = uint64(binary.BigEndian.Uint64(r.readMany(8)))
	default:
		panic(fail)
	}
}

// ---------------------------------------------------------------------------

// DecodeFloat32 decodes a float32 value from the Reader.
func (r *Reader) DecodeFloat32(v *float32) {
	if r.readOne() == cFloat32 {
		b := binary.BigEndian.Uint32(r.readMany(4))
		*v = math.Float32frombits(b)
		return
	}
	panic(fail)
}

// DecodeFloat64 decodes a float64 value from the Reader.
func (r *Reader) DecodeFloat64(v *float64) {
	switch r.readOne() {
	case cFloat32:
		b := binary.BigEndian.Uint32(r.readMany(4))
		*v = float64(math.Float32frombits(b))
	case cFloat64:
		b := uint64(binary.BigEndian.Uint64(r.readMany(8)))
		*v = math.Float64frombits(b)
	default:
		panic(fail)
	}
}

// ---------------------------------------------------------------------------

// DecodeComplex64 decodes a complex64 value from the Reader.
func (r *Reader) DecodeComplex64(v *complex64) {
	if r.readOne() == cComplex64 {
		one := binary.BigEndian.Uint32(r.readMany(4))
		two := binary.BigEndian.Uint32(r.readMany(4))
		*v = complex(math.Float32frombits(one), math.Float32frombits(two))
		return
	}
	panic(fail)
}

// DecodeComplex128 decodes a complex128 value from the Reader.
func (r *Reader) DecodeComplex128(v *complex128) {
	if r.readOne() == cComplex128 {
		one := binary.BigEndian.Uint64(r.readMany(8))
		two := binary.BigEndian.Uint64(r.readMany(8))
		*v = complex(math.Float64frombits(one), math.Float64frombits(two))
		return
	}
	panic(fail)
}

// ---------------------------------------------------------------------------

// DecodeTime decodes a time.Time value from the Reader.
func (r *Reader) DecodeTime(v *time.Time) {
	if r.readOne() == cTime {
		b := int64(binary.BigEndian.Uint64(r.readMany(8)))
		*v = time.Unix(0, b).UTC()
		return
	}
	panic(fail)
}

// ---------------------------------------------------------------------------

// DecodeArr decodes an array from the Reader.
func (r *Reader) DecodeArr(v interface{}) {
	switch a := v.(type) {
	case *[]bool:
		r.decodeArrBool(a)
	case *[]int:
		r.decodeArrInt(a)
	case *[]int8:
		r.decodeArrInt8(a)
	case *[]int16:
		r.decodeArrInt16(a)
	case *[]int32:
		r.decodeArrInt32(a)
	case *[]int64:
		r.decodeArrInt64(a)
	case *[]uint:
		r.decodeArrUint(a)
	case *[]uint8:
		r.decodeArrUint8(a)
	case *[]uint16:
		r.decodeArrUint16(a)
	case *[]uint32:
		r.decodeArrUint32(a)
	case *[]uint64:
		r.decodeArrUint64(a)
	case *[]string:
		r.decodeArrString(a)
	case *[]float32:
		r.decodeArrFloat32(a)
	case *[]float64:
		r.decodeArrFloat64(a)
	case *[]complex64:
		r.decodeArrComplex64(a)
	case *[]complex128:
		r.decodeArrComplex128(a)
	case *[]time.Time:
		r.decodeArrTime(a)
	case *[]interface{}:
		r.decodeArrAny(a)
	case reflect.Value:
		r.decodeArr(a)
	default:
		r.DecodeAny(v)
	}
}

func (r *Reader) decodeArrLen() int {
	b := r.readOne()
	switch {
	case b >= cFixArr && b <= cFixArr+fixedArr:
		return int(b - cFixArr)
	case b == cArr:
		return r.readLen()
	default:
		panic(fail)
	}
}

// DecodeMap decodes a map from the Reader.
func (r *Reader) DecodeMap(v interface{}) {
	switch m := v.(type) {
	case *map[string]int:
		r.decodeMapStringInt(m)
	case *map[string]uint:
		r.decodeMapStringUint(m)
	case *map[string]bool:
		r.decodeMapStringBool(m)
	case *map[string]string:
		r.decodeMapStringString(m)
	case *map[int]interface{}:
		r.decodeMapIntAny(m)
	case *map[uint]interface{}:
		r.decodeMapUintAny(m)
	case *map[string]interface{}:
		r.decodeMapStringAny(m)
	case *map[time.Time]interface{}:
		r.decodeMapTimeAny(m)
	case *map[interface{}]interface{}:
		r.decodeMapAnyAny(m)
	case reflect.Value:
		r.decodeMap(m)
	default:
		r.DecodeAny(v)
	}
}

func (r *Reader) decodeMapLen() int {
	b := r.readOne()
	switch {
	case b >= cFixMap && b <= cFixMap+fixedMap:
		return int(b - cFixMap)
	case b == cMap:
		return r.readLen()
	default:
		panic(fail)
	}
}
