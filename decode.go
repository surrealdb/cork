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
	"bytes"
	"encoding"
	"encoding/binary"
	"fmt"
	"io"
	"reflect"
	"time"
)

// Decoder represents a CORK decoder.
type Decoder struct {
	h *Handle
	r *reader
}

// Decode decodes a CORK into a data object.
func Decode(src []byte) (dst interface{}) {
	buf := bytes.NewBuffer(src)
	NewDecoder(buf).Decode(&dst)
	return dst
}

// NewDecoder returns a Decoder for decoding from an io.Reader.
func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		r: newReader(r),
		h: &Handle{
			Precision: false,
		},
	}
}

// Options sets the configuration options that the Decoder should use.
func (d *Decoder) Options(h *Handle) *Decoder {
	d.h = h
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
			if catch, ok := r.(string); ok {
				err = fmt.Errorf(catch)
			}
			if catch, ok := r.(error); ok {
				err = catch
			}
		}
	}()
	d.decode(dst)
	return
}

func (d *Decoder) decode(dst interface{}) {
	d.proceed(dst, d.decodeBit())
}

func (d *Decoder) proceed(dst interface{}, b byte) {

	switch val := dst.(type) {

	case *bool:
		if isVal(b) {
			*val = d.decodeVal(b)
			return
		}

	case *[]byte:
		if isBin(b) {
			*val = d.decodeBin(b)
			return
		}

	case *string:
		if isStr(b) {
			*val = d.decodeStr(b)
			return
		}

	case *time.Time:
		if isTime(b) {
			*val = d.decodeTime(b)
			return
		}

	case *int:
		if isInt(b) {
			*val = int(d.decodeInt(b))
			return
		}

	case *int8:
		switch {
		case isNum(b):
			*val = int8(d.decodeInt1(b))
			return
		case b == cInt8:
			*val = int8(d.decodeInt8(b))
			return
		}

	case *int16:
		switch {
		case isNum(b):
			*val = int16(d.decodeInt1(b))
			return
		case b == cInt8:
			*val = int16(d.decodeInt8(b))
			return
		case b == cInt16:
			*val = int16(d.decodeInt16(b))
			return
		}

	case *int32:
		switch {
		case isNum(b):
			*val = int32(d.decodeInt1(b))
			return
		case b == cInt8:
			*val = int32(d.decodeInt8(b))
			return
		case b == cInt16:
			*val = int32(d.decodeInt16(b))
			return
		case b == cInt32:
			*val = int32(d.decodeInt32(b))
			return
		}

	case *int64:
		switch {
		case isNum(b):
			*val = int64(d.decodeInt1(b))
			return
		case b == cInt8:
			*val = int64(d.decodeInt8(b))
			return
		case b == cInt16:
			*val = int64(d.decodeInt16(b))
			return
		case b == cInt32:
			*val = int64(d.decodeInt32(b))
			return
		case b == cInt64:
			*val = int64(d.decodeInt64(b))
			return
		}

	case *uint:
		if isUint(b) {
			*val = uint(d.decodeUint(b))
			return
		}

	case *uint8:
		switch {
		case isNum(b):
			*val = uint8(d.decodeUint1(b))
			return
		case b == cUint8:
			*val = uint8(d.decodeUint8(b))
			return
		}

	case *uint16:
		switch {
		case isNum(b):
			*val = uint16(d.decodeUint1(b))
			return
		case b == cUint8:
			*val = uint16(d.decodeUint8(b))
			return
		case b == cUint16:
			*val = uint16(d.decodeUint16(b))
			return
		}

	case *uint32:
		switch {
		case isNum(b):
			*val = uint32(d.decodeUint1(b))
			return
		case b == cUint8:
			*val = uint32(d.decodeUint8(b))
			return
		case b == cUint16:
			*val = uint32(d.decodeUint16(b))
			return
		case b == cUint32:
			*val = uint32(d.decodeUint32(b))
			return
		}

	case *uint64:
		switch {
		case isNum(b):
			*val = uint64(d.decodeUint1(b))
			return
		case b == cUint8:
			*val = uint64(d.decodeUint8(b))
			return
		case b == cUint16:
			*val = uint64(d.decodeUint16(b))
			return
		case b == cUint32:
			*val = uint64(d.decodeUint32(b))
			return
		case b == cUint64:
			*val = uint64(d.decodeUint64(b))
			return
		}

	case *float32:
		switch {
		case b == cFloat32:
			*val = float32(d.decodeFloat32(b))
			return
		}

	case *float64:
		switch {
		case b == cFloat32:
			*val = float64(d.decodeFloat32(b))
			return
		case b == cFloat64:
			*val = float64(d.decodeFloat64(b))
			return
		}

	case *complex64:
		switch {
		case b == cComplex64:
			*val = complex64(d.decodeComplex64(b))
			return
		}

	case *complex128:
		switch {
		case b == cComplex64:
			*val = complex128(d.decodeComplex64(b))
			return
		case b == cComplex128:
			*val = complex128(d.decodeComplex128(b))
			return
		}

	// ---------------------------------------------
	// Include common slice types
	// ---------------------------------------------

	case *[]interface{}:
		if b == cArr {
			*val = d.decodeArrNil(b)
			return
		}
	case *[]bool:
		if b == cArr {
			*val = d.decodeArrVal(b)
			return
		}
	case *[]string:
		if b == cArr {
			*val = d.decodeArrStr(b)
			return
		}
	case *[]int:
		if b == cArr {
			*val = d.decodeArrInt(b)
			return
		}
	case *[]int8:
		if b == cArr {
			*val = d.decodeArrInt8(b)
			return
		}
	case *[]int16:
		if b == cArr {
			*val = d.decodeArrInt16(b)
			return
		}
	case *[]int32:
		if b == cArr {
			*val = d.decodeArrInt32(b)
			return
		}
	case *[]int64:
		if b == cArr {
			*val = d.decodeArrInt64(b)
			return
		}
	case *[]uint:
		if b == cArr {
			*val = d.decodeArrUint(b)
			return
		}
	case *[]uint16:
		if b == cArr {
			*val = d.decodeArrUint16(b)
			return
		}
	case *[]uint32:
		if b == cArr {
			*val = d.decodeArrUint32(b)
			return
		}
	case *[]uint64:
		if b == cArr {
			*val = d.decodeArrUint64(b)
			return
		}
	case *[]float32:
		if b == cArr {
			*val = d.decodeArrFloat32(b)
			return
		}
	case *[]float64:
		if b == cArr {
			*val = d.decodeArrFloat64(b)
			return
		}
	case *[]complex64:
		if b == cArr {
			*val = d.decodeArrComplex64(b)
			return
		}
	case *[]complex128:
		if b == cArr {
			*val = d.decodeArrComplex128(b)
			return
		}
	case *[]time.Time:
		if b == cArr {
			*val = d.decodeArrTime(b)
			return
		}

	// ---------------------------------------------
	// Include common map[string]<T> types
	// ---------------------------------------------

	case *map[string]int:
		if b == cMap {
			*val = d.decodeMapStrInt(b)
			return
		}
	case *map[string]uint:
		if b == cMap {
			*val = d.decodeMapStrUint(b)
			return
		}
	case *map[string]bool:
		if b == cMap {
			*val = d.decodeMapStrVal(b)
			return
		}
	case *map[string]string:
		if b == cMap {
			*val = d.decodeMapStrStr(b)
			return
		}
	case *map[string]interface{}:
		if b == cMap {
			*val = d.decodeMapStrNil(b)
			return
		}

	// ---------------------------------------------
	// Include common map[<T>]interface{} types
	// ---------------------------------------------

	case *map[interface{}]interface{}:
		if b == cMap {
			*val = d.decodeMapNilNil(b)
			return
		}

	// ---------------------------------------------
	// Include nil interface{} type decoding
	// ---------------------------------------------

	case *interface{}:

		switch {
		case b == cNil:
			*val = nil
			return
		case b == cTrue:
			*val = d.decodeVal(b)
			return
		case b == cFalse:
			*val = d.decodeVal(b)
			return
		case b == cTime:
			*val = d.decodeTime(b)
			return

		case isBin(b):
			*val = d.decodeBin(b)
			return
		case isStr(b):
			*val = d.decodeStr(b)
			return
		case isExt(b):
			*val = d.decodeExt(b)
			return

		case isNum(b):
			if d.h.Precision {
				*val = d.decodeInt1(b)
			} else {
				*val = d.decodeIntAny(b)
			}
			return
		case b == cInt8:
			if d.h.Precision {
				*val = d.decodeInt8(b)
			} else {
				*val = d.decodeIntAny(b)
			}
			return
		case b == cInt16:
			if d.h.Precision {
				*val = d.decodeInt16(b)
			} else {
				*val = d.decodeIntAny(b)
			}
			return
		case b == cInt32:
			if d.h.Precision {
				*val = d.decodeInt32(b)
			} else {
				*val = d.decodeIntAny(b)
			}
			return
		case b == cInt64:
			if d.h.Precision {
				*val = d.decodeInt64(b)
			} else {
				*val = d.decodeIntAny(b)
			}
			return
		case b == cUint8:
			if d.h.Precision {
				*val = d.decodeUint8(b)
			} else {
				*val = d.decodeUintAny(b)
			}
			return
		case b == cUint16:
			if d.h.Precision {
				*val = d.decodeUint16(b)
			} else {
				*val = d.decodeUintAny(b)
			}
			return
		case b == cUint32:
			if d.h.Precision {
				*val = d.decodeUint32(b)
			} else {
				*val = d.decodeUintAny(b)
			}
			return
		case b == cUint64:
			if d.h.Precision {
				*val = d.decodeUint64(b)
			} else {
				*val = d.decodeUintAny(b)
			}
			return
		case b == cFloat32:
			*val = d.decodeFloat32(b)
			return
		case b == cFloat64:
			*val = d.decodeFloat64(b)
			return
		case b == cComplex64:
			*val = d.decodeComplex64(b)
			return
		case b == cComplex128:
			*val = d.decodeComplex128(b)
			return

		case isArr(b):
			*val = d.decodeArr(b)
			return
		case isMap(b):
			*val = d.decodeMap(b)
			return

		}

	// ---------------------------------------------
	// Include self decoders
	// ---------------------------------------------

	case Corker:
		d.decodeCrk(b, val)
		return

	case encoding.BinaryUnmarshaler:
		if isBin(b) {
			bin := d.decodeBin(b)
			err := val.UnmarshalBinary(bin)
			if err != nil {
				panic(err)
			}
			return
		}

	case encoding.TextUnmarshaler:
		if isStr(b) {
			bin := d.decodeTxt(b)
			err := val.UnmarshalText(bin)
			if err != nil {
				panic(err)
			}
			return
		}

	// ---------------------------------------------
	// Use reflect for any remaining types
	// ---------------------------------------------

	default:
		d.decodeRef(b, reflect.ValueOf(dst))
		return

	}

	panic(fmt.Errorf("Can't decode into %T", dst))

}

func (d *Decoder) decodeBit() (val byte) {
	return d.r.ReadOne()
}

func (d *Decoder) decodeVal(b byte) (val bool) {
	return b == cTrue
}

func (d *Decoder) decodeBin(b byte) (val []byte) {
	var sze int
	switch {
	case b >= cFixBin && b <= cFixBin+fixedBin:
		sze = int(b - cFixBin)
	case b == cBin8:
		var tmp uint8
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cBin16:
		var tmp uint16
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cBin32:
		var tmp uint32
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cBin64:
		var tmp uint64
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	}
	return d.r.ReadMany(sze)
}

func (d *Decoder) decodeStr(b byte) (val string) {
	var sze int
	switch {
	case b >= cFixStr && b <= cFixStr+fixedStr:
		sze = int(b - cFixStr)
	case b == cStr8:
		var tmp uint8
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr16:
		var tmp uint16
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr32:
		var tmp uint32
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr64:
		var tmp uint64
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	}
	return string(d.r.ReadMany(sze))
}

func (d *Decoder) decodeExt(b byte) (val Corker) {
	var sze int
	switch {
	case b >= cFixExt && b <= cFixExt+fixedExt:
		sze = int(b - cFixExt)
	case b == cExt8:
		var tmp uint8
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cExt16:
		var tmp uint16
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cExt32:
		var tmp uint32
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cExt64:
		var tmp uint64
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	}
	bit := d.r.ReadOne()
	bin := d.r.ReadMany(sze)
	obj := reflect.New(registry[bit]).Interface().(Corker)
	err := obj.UnmarshalCORK(bin)
	if err != nil {
		panic(err)
	}
	return obj
}

func (d *Decoder) decodeCrk(b byte, val Corker) {
	var sze int
	switch {
	case b >= cFixExt && b <= cFixExt+fixedExt:
		sze = int(b - cFixExt)
	case b == cExt8:
		var tmp uint8
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cExt16:
		var tmp uint16
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cExt32:
		var tmp uint32
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cExt64:
		var tmp uint64
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	}
	d.r.ReadOne()
	bin := d.r.ReadMany(sze)
	err := val.UnmarshalCORK(bin)
	if err != nil {
		panic(err)
	}
	return
}

func (d *Decoder) decodeTxt(b byte) (val []byte) {
	var sze int
	switch {
	case b >= cFixStr && b <= cFixStr+fixedStr:
		sze = int(b - cFixStr)
	case b == cStr8:
		var tmp uint8
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr16:
		var tmp uint16
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr32:
		var tmp uint32
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr64:
		var tmp uint64
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	}
	return d.r.ReadMany(sze)
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func (d *Decoder) decodeInt(b byte) (val int) {
	switch {
	case isNum(b):
		val = int(d.decodeInt1(b))
	case b == cInt8:
		val = int(d.decodeInt8(b))
	case b == cInt16:
		val = int(d.decodeInt16(b))
	case b == cInt32:
		val = int(d.decodeInt32(b))
	case b == cInt64:
		val = int(d.decodeInt64(b))
	}
	return
}

func (d *Decoder) decodeIntAny(b byte) (val int64) {
	return int64(d.decodeInt(b))
}

func (d *Decoder) decodeInt1(b byte) (val int8) {
	return int8(b)
}

func (d *Decoder) decodeInt8(b byte) (val int8) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeInt16(b byte) (val int16) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeInt32(b byte) (val int32) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeInt64(b byte) (val int64) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeUint(b byte) (val uint) {
	switch {
	case isNum(b):
		val = uint(d.decodeUint1(b))
	case b == cUint8:
		val = uint(d.decodeUint8(b))
	case b == cUint16:
		val = uint(d.decodeUint16(b))
	case b == cUint32:
		val = uint(d.decodeUint32(b))
	case b == cUint64:
		val = uint(d.decodeUint64(b))
	}
	return
}

func (d *Decoder) decodeUintAny(b byte) (val uint64) {
	return uint64(d.decodeUint(b))
}

func (d *Decoder) decodeUint1(b byte) (val uint8) {
	return uint8(b)
}

func (d *Decoder) decodeUint8(b byte) (val uint8) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeUint16(b byte) (val uint16) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeUint32(b byte) (val uint32) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeUint64(b byte) (val uint64) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeTime(b byte) (val time.Time) {
	var tmp int64
	binary.Read(d.r, binary.BigEndian, &tmp)
	return time.Unix(0, tmp).UTC()
}

func (d *Decoder) decodeFloat32(b byte) (val float32) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeFloat64(b byte) (val float64) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeComplex64(b byte) (val complex64) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

func (d *Decoder) decodeComplex128(b byte) (val complex128) {
	binary.Read(d.r, binary.BigEndian, &val)
	return
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func (d *Decoder) decodeArr(b byte) (val interface{}) {
	if d.h.ArrType != nil {
		switch d.h.ArrType.(type) {
		default:
			return d.decodeArrNil(b)
		case []bool:
			return d.decodeArrVal(b)
		case []string:
			return d.decodeArrStr(b)
		case []int:
			return d.decodeArrInt(b)
		case []int8:
			return d.decodeArrInt8(b)
		case []int16:
			return d.decodeArrInt16(b)
		case []int32:
			return d.decodeArrInt32(b)
		case []int64:
			return d.decodeArrInt64(b)
		case []uint:
			return d.decodeArrUint(b)
		case []uint8:
			return d.decodeArrUint8(b)
		case []uint16:
			return d.decodeArrUint16(b)
		case []uint32:
			return d.decodeArrUint32(b)
		case []uint64:
			return d.decodeArrUint64(b)
		case []float32:
			return d.decodeArrFloat32(b)
		case []float64:
			return d.decodeArrFloat64(b)
		case []complex64:
			return d.decodeArrComplex64(b)
		case []complex128:
			return d.decodeArrComplex128(b)
		case []time.Time:
			return d.decodeArrTime(b)
		case reflect.Type:
			return d.decodeArrUnk(b)
		}
	}
	return d.decodeArrNil(b)
}

func (d *Decoder) decodeArrNil(b byte) (val []interface{}) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]interface{}, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrVal(b byte) (val []bool) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]bool, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrStr(b byte) (val []string) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]string, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrInt(b byte) (val []int) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]int, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrInt8(b byte) (val []int8) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]int8, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrInt16(b byte) (val []int16) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]int16, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrInt32(b byte) (val []int32) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]int32, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrInt64(b byte) (val []int64) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]int64, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrUint(b byte) (val []uint) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]uint, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrUint8(b byte) (val []uint8) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]uint8, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrUint16(b byte) (val []uint16) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]uint16, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrUint32(b byte) (val []uint32) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]uint32, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrUint64(b byte) (val []uint64) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]uint64, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrFloat32(b byte) (val []float32) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]float32, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrFloat64(b byte) (val []float64) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]float64, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrComplex64(b byte) (val []complex64) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]complex64, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrComplex128(b byte) (val []complex128) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]complex128, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrTime(b byte) (val []time.Time) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]time.Time, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func (d *Decoder) decodeMap(b byte) (val interface{}) {
	if d.h.MapType != nil {
		switch d.h.MapType.(type) {
		default:
			return d.decodeMapNilNil(b)
		case map[string]bool:
			return d.decodeMapStrVal(b)
		case map[string]string:
			return d.decodeMapStrStr(b)
		case map[string]int:
			return d.decodeMapStrInt(b)
		case map[string]uint:
			return d.decodeMapStrUint(b)
		case map[string]interface{}:
			return d.decodeMapStrNil(b)
		case reflect.Type:
			return d.decodeMapUnk(b)
		}
	}
	return d.decodeMapNilNil(b)
}

func (d *Decoder) decodeMapNilNil(b byte) (val map[interface{}]interface{}) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[interface{}]interface{}, tot)
	for i := 0; i < tot; i++ {
		var k interface{}
		var v interface{}
		d.decode(&k)
		d.decode(&v)
		obj[k] = v
	}
	return obj
}

func (d *Decoder) decodeMapStrNil(b byte) (val map[string]interface{}) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]interface{}, tot)
	for i := 0; i < tot; i++ {
		var k string
		var v interface{}
		d.decode(&k)
		d.decode(&v)
		obj[k] = v
	}
	return obj
}

func (d *Decoder) decodeMapStrVal(b byte) (val map[string]bool) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]bool, tot)
	for i := 0; i < tot; i++ {
		var k string
		var v bool
		d.decode(&k)
		d.decode(&v)
		obj[k] = v
	}
	return obj
}

func (d *Decoder) decodeMapStrStr(b byte) (val map[string]string) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]string, tot)
	for i := 0; i < tot; i++ {
		var k string
		var v string
		d.decode(&k)
		d.decode(&v)
		obj[k] = v
	}
	return obj
}

func (d *Decoder) decodeMapStrInt(b byte) (val map[string]int) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]int, tot)
	for i := 0; i < tot; i++ {
		var k string
		var v int
		d.decode(&k)
		d.decode(&v)
		obj[k] = v
	}
	return obj
}

func (d *Decoder) decodeMapStrUint(b byte) (val map[string]uint) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]uint, tot)
	for i := 0; i < tot; i++ {
		var k string
		var v uint
		d.decode(&k)
		d.decode(&v)
		obj[k] = v
	}
	return obj
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func (d *Decoder) decodeRef(b byte, val reflect.Value) {

	switch val.Type().Kind() {

	case reflect.Ptr:
		if val.Elem().IsValid() {
			d.decodeRef(b, val.Elem())
		} else {
			val.Set(reflect.ValueOf(d.decodeExt(b)))
		}
		return

	case reflect.Map:
		d.decodeMapAny(b, val)
		return

	case reflect.Slice:
		d.decodeArrAny(b, val)
		return

	case reflect.Struct:
		if _, ok := val.Addr().Interface().(Corker); ok {
			val.Set(reflect.ValueOf(d.decodeExt(b)).Elem())
		} else {
			d.decodeStructMap(b, val)
		}
		return

	case reflect.Interface:
		var i interface{}
		d.proceed(&i, b)
		val.Set(reflect.ValueOf(i))
		return

	case reflect.Bool:
		var i bool
		d.proceed(&i, b)
		val.SetBool(i)

	case reflect.String:
		var i string
		d.proceed(&i, b)
		val.SetString(i)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var i int
		d.proceed(&i, b)
		val.SetInt(int64(i))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		var i uint
		d.proceed(&i, b)
		val.SetUint(uint64(i))

	case reflect.Float32:
		var i float32
		d.proceed(&i, b)
		val.SetFloat(float64(i))

	case reflect.Float64:
		var i float64
		d.proceed(&i, b)
		val.SetFloat(float64(i))

	case reflect.Complex64:
		var i complex64
		d.proceed(&i, b)
		val.SetComplex(complex128(i))

	case reflect.Complex128:
		var i complex128
		d.proceed(&i, b)
		val.SetComplex(complex128(i))

	}

}

func (d *Decoder) decodeArrUnk(b byte) (val interface{}) {
	typ := d.h.ArrType.(reflect.Type)
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := reflect.MakeSlice(typ, tot, tot)
	for i := 0; i < tot; i++ {
		if typ.Elem().Kind() == reflect.Ptr {
			v := reflect.New(typ.Elem().Elem())
			d.decode(v.Interface())
			arr.Index(i).Set(v)
		} else {
			v := reflect.New(typ.Elem())
			d.decode(v.Interface())
			arr.Index(i).Set(v.Elem())
		}
	}
	return arr.Interface()
}

func (d *Decoder) decodeArrAny(b byte, arr reflect.Value) {
	typ := arr.Type()
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	if arr.IsNil() {
		arr.Set(reflect.MakeSlice(typ, tot, tot))
	}
	for i := 0; i < tot; i++ {
		if typ.Elem().Kind() == reflect.Ptr {
			v := reflect.New(typ.Elem().Elem())
			d.decode(v.Interface())
			arr.Index(i).Set(v)
		} else {
			v := reflect.New(typ.Elem())
			d.decode(v.Interface())
			arr.Index(i).Set(v.Elem())
		}
	}
}

func (d *Decoder) decodeMapUnk(b byte) (val interface{}) {
	typ := d.h.MapType.(reflect.Type)
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := reflect.MakeMap(typ)
	for i := 0; i < tot; i++ {
		k := reflect.New(typ.Key())
		d.decode(k.Interface())
		v := reflect.New(typ.Elem())
		d.decode(v.Interface())
		obj.SetMapIndex(k.Elem(), v.Elem())
	}
	return obj.Interface()
}

func (d *Decoder) decodeMapAny(b byte, obj reflect.Value) {
	typ := obj.Type()
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	if obj.IsNil() {
		obj.Set(reflect.MakeMap(typ))
	}
	for i := 0; i < tot; i++ {
		k := reflect.New(typ.Key())
		d.decode(k.Interface())
		v := reflect.New(typ.Elem())
		d.decode(v.Interface())
		obj.SetMapIndex(k.Elem(), v.Elem())
	}
}

func (d *Decoder) decodeStructMap(b byte, item reflect.Value) {

	typ := item.Type()
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]*field, tot)

	for i := 0; i < item.NumField(); i++ {
		if fld := newField(typ.Field(i), item.Field(i)); fld != nil {
			obj[fld.Show()] = fld
		}
	}

	for i := 0; i < tot; i++ {

		nxt = d.decodeBit()
		k := d.decodeStr(nxt)

		if itm, ok := obj[k]; ok {

			fld := item.FieldByName(itm.Name())

			// Println(i, itm, fld, fld.Interface())
			// Printf("%T %v \n", itm, item)
			// Printf("%T %v \n", fld, fld)
			// Printf("%T %v \n", fld.Interface(), fld.Interface())
			// Printf("%T %v \n", fld.Addr().Interface(), fld.Addr().Interface())
			// Println("---")

			if fld.CanSet() {
				if fld.CanAddr() {
					d.decode(fld.Addr().Interface())
				} else {
					d.decode(fld.Interface())
				}
			}

		}

	}

}
