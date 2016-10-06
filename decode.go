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
	"encoding/binary"
	"io"
	"reflect"
	"time"
)

// Decoder represents a CORK decoder.
type Decoder struct {
	r *reader
}

// Decode decodes a CORK into a data object.
func Decode(src []byte) (dst interface{}) {
	buf := bytes.NewBuffer(src)
	NewDecoder(buf).Decode(&dst)
	return dst
}

// NewDecoder returns a Decoder for decoding from an io.Reader.
func NewDecoder(src io.Reader) *Decoder {
	return &Decoder{r: newReader(src)}
}

// Decode decodes the stream from reader and stores the result in the
// value pointed to by v. v cannot be a nil pointer. v can also be
// a reflect.Value of a pointer.
//
// Note that a pointer to a nil interface is not a nil pointer.
// If you do not know what type of stream it is, pass in a pointer to a nil interface.
// We will decode and store a value in that nil interface.
//
// Sample usages:
//   // Decoding into a non-nil typed value
//   var f float32
//   err = codec.NewDecoder(r, handle).Decode(&f)
//
//   // Decoding into nil interface
//   var v interface{}
//   dec := codec.NewDecoder(r, handle)
//   err = dec.Decode(&v)
//
// When decoding into a nil interface{}, we will decode into an appropriate value based
// on the contents of the stream:
//   - Numbers are decoded as float64, int64 or uint64.
//   - Other values are decoded appropriately depending on the type:
//     bool, string, []byte, time.Time, etc
//   - Extensions are decoded as RawExt (if no ext function registered for the tag)
// Configurations exist on the Handle to override defaults
// (e.g. for MapType, SliceType and how to decode raw bytes).
//
// When decoding into a non-nil interface{} value, the mode of encoding is based on the
// type of the value. When a value is seen:
//   - If an extension is registered for it, call that extension function
//   - If it implements BinaryUnmarshaler, call its UnmarshalBinary(data []byte) error
//   - Else decode it based on its reflect.Kind
//
// There are some special rules when decoding into containers (slice/array/map/struct).
// Decode will typically use the stream contents to UPDATE the container.
//   - A map can be decoded from a stream map, by updating matching keys.
//   - A slice can be decoded from a stream array,
//     by updating the first n elements, where n is length of the stream.
//   - A slice can be decoded from a stream map, by decoding as if
//     it contains a sequence of key-value pairs.
//   - A struct can be decoded from a stream map, by updating matching fields.
//   - A struct can be decoded from a stream array,
//     by updating fields as they occur in the struct (by index).
//
// When decoding a stream map or array with length of 0 into a nil map or slice,
// we reset the destination map or slice to a zero-length value.
//
// However, when decoding a stream nil, we reset the destination container
// to its "zero" value (e.g. nil for slice/map, etc).
//
func (d *Decoder) Decode(dst interface{}) (err error) {
	// TODO need to catch panics and enable errors in decoder
	d.decode(dst)
	return
}

func (d *Decoder) decode(dst interface{}, bits ...byte) {

	var b byte

	if len(bits) == 0 {
		b = d.decodeBit()
	} else {
		b = bits[0]
	}

	switch val := dst.(type) {

	case *bool:
		if isVal(b) {
			*val = d.decodeVal(b)
		}

	case *[]byte:
		if isBin(b) {
			*val = d.decodeBin(b)
		}

	case *string:
		if isStr(b) {
			*val = d.decodeStr(b)
		}

	case *time.Time:
		if isTime(b) {
			*val = d.decodeTime(b)
		}

	case *int:
		if isInt(b) {
			*val = int(d.decodeInt(b))
		}

	case *int8:
		switch {
		case isNum(b):
			*val = int8(d.decodeInt1(b))
		case b == cInt8:
			*val = int8(d.decodeInt8(b))
		}

	case *int16:
		switch {
		case isNum(b):
			*val = int16(d.decodeInt1(b))
		case b == cInt8:
			*val = int16(d.decodeInt8(b))
		case b == cInt16:
			*val = int16(d.decodeInt16(b))
		}

	case *int32:
		switch {
		case isNum(b):
			*val = int32(d.decodeInt1(b))
		case b == cInt8:
			*val = int32(d.decodeInt8(b))
		case b == cInt16:
			*val = int32(d.decodeInt16(b))
		case b == cInt32:
			*val = int32(d.decodeInt32(b))
		}

	case *int64:
		switch {
		case isNum(b):
			*val = int64(d.decodeInt1(b))
		case b == cInt8:
			*val = int64(d.decodeInt8(b))
		case b == cInt16:
			*val = int64(d.decodeInt16(b))
		case b == cInt32:
			*val = int64(d.decodeInt32(b))
		case b == cInt64:
			*val = int64(d.decodeInt64(b))
		}

	case *uint:
		if isUint(b) {
			*val = uint(d.decodeUint(b))
		}

	case *uint8:
		switch {
		case isNum(b):
			*val = uint8(d.decodeUint1(b))
		case b == cUint8:
			*val = uint8(d.decodeUint8(b))
		}

	case *uint16:
		switch {
		case isNum(b):
			*val = uint16(d.decodeUint1(b))
		case b == cUint8:
			*val = uint16(d.decodeUint8(b))
		case b == cUint16:
			*val = uint16(d.decodeUint16(b))
		}

	case *uint32:
		switch {
		case isNum(b):
			*val = uint32(d.decodeUint1(b))
		case b == cUint8:
			*val = uint32(d.decodeUint8(b))
		case b == cUint16:
			*val = uint32(d.decodeUint16(b))
		case b == cUint32:
			*val = uint32(d.decodeUint32(b))
		}

	case *uint64:
		switch {
		case isNum(b):
			*val = uint64(d.decodeUint1(b))
		case b == cUint8:
			*val = uint64(d.decodeUint8(b))
		case b == cUint16:
			*val = uint64(d.decodeUint16(b))
		case b == cUint32:
			*val = uint64(d.decodeUint32(b))
		case b == cUint64:
			*val = uint64(d.decodeUint64(b))
		}

	case *float32:
		switch {
		case b == cFloat32:
			*val = float32(d.decodeFloat32(b))
		}

	case *float64:
		switch {
		case b == cFloat32:
			*val = float64(d.decodeFloat32(b))
		case b == cFloat64:
			*val = float64(d.decodeFloat64(b))
		}

	case *[]interface{}:
		if b == cArrNil {
			*val = d.decodeArrNil(b)
		}
	case *[]bool:
		if b == cArrBool {
			*val = d.decodeArrBool(b)
		}
	case *[]string:
		if b == cArrStr {
			*val = d.decodeArrStr(b)
		}
	case *[]int:
		if b == cArrInt {
			*val = d.decodeArrInt(b)
		}
	case *[]int8:
		if b == cArrInt8 {
			*val = d.decodeArrInt8(b)
		}
	case *[]int16:
		if b == cArrInt16 {
			*val = d.decodeArrInt16(b)
		}
	case *[]int32:
		if b == cArrInt32 {
			*val = d.decodeArrInt32(b)
		}
	case *[]int64:
		if b == cArrInt64 {
			*val = d.decodeArrInt64(b)
		}
	case *[]uint:
		if b == cArrUint {
			*val = d.decodeArrUint(b)
		}
	case *[]uint16:
		if b == cArrUint16 {
			*val = d.decodeArrUint16(b)
		}
	case *[]uint32:
		if b == cArrUint32 {
			*val = d.decodeArrUint32(b)
		}
	case *[]uint64:
		if b == cArrUint64 {
			*val = d.decodeArrUint64(b)
		}
	case *[]float32:
		if b == cArrFloat32 {
			*val = d.decodeArrFloat32(b)
		}
	case *[]float64:
		if b == cArrFloat64 {
			*val = d.decodeArrFloat64(b)
		}
	case *[]time.Time:
		if b == cArrTime {
			*val = d.decodeArrTime(b)
		}

	case *map[string]int:
		if b == cMapStrInt {
			*val = d.decodeMapStrInt(b)
		}
	case *map[string]bool:
		if b == cMapStrBool {
			*val = d.decodeMapStrBool(b)
		}
	case *map[string]string:
		if b == cMapStrStr {
			*val = d.decodeMapStrStr(b)
		}
	case *map[string]interface{}:
		if b == cMapStrNil {
			*val = d.decodeMapStrNil(b)
		}
	case *map[interface{}]interface{}:
		if b == cMapNilNil {
			*val = d.decodeMapNilNil(b)
		}

	case *interface{}:

		switch {
		case b == cNil:
			*val = nil
		case b == cTrue:
			*val = d.decodeVal(b)
		case b == cFalse:
			*val = d.decodeVal(b)
		case b == cTime:
			*val = d.decodeTime(b)

		case isBin(b):
			*val = d.decodeBin(b)
		case isStr(b):
			*val = d.decodeStr(b)
		case isExt(b):
			*val = d.decodeExt(b)

		case isNum(b):
			*val = d.decodeInt1(b)
		case b == cInt8:
			*val = d.decodeInt8(b)
		case b == cInt16:
			*val = d.decodeInt16(b)
		case b == cInt32:
			*val = d.decodeInt32(b)
		case b == cInt64:
			*val = d.decodeInt64(b)
		case b == cUint8:
			*val = d.decodeUint8(b)
		case b == cUint16:
			*val = d.decodeUint16(b)
		case b == cUint32:
			*val = d.decodeUint32(b)
		case b == cUint64:
			*val = d.decodeUint64(b)
		case b == cFloat32:
			*val = d.decodeFloat32(b)
		case b == cFloat64:
			*val = d.decodeFloat64(b)

		case b == cArr:
		case b == cArrNil:
			*val = d.decodeArrNil(b)
		case b == cArrBool:
			*val = d.decodeArrBool(b)
		case b == cArrStr:
			*val = d.decodeArrStr(b)
		case b == cArrInt:
			*val = d.decodeArrInt(b)
		case b == cArrInt8:
			*val = d.decodeArrInt8(b)
		case b == cArrInt16:
			*val = d.decodeArrInt16(b)
		case b == cArrInt32:
			*val = d.decodeArrInt32(b)
		case b == cArrInt64:
			*val = d.decodeArrInt64(b)
		case b == cArrUint:
			*val = d.decodeArrUint(b)
		case b == cArrUint16:
			*val = d.decodeArrUint16(b)
		case b == cArrUint32:
			*val = d.decodeArrUint32(b)
		case b == cArrUint64:
			*val = d.decodeArrUint64(b)
		case b == cArrFloat32:
			*val = d.decodeArrFloat32(b)
		case b == cArrFloat64:
			*val = d.decodeArrFloat64(b)
		case b == cArrTime:
			*val = d.decodeArrTime(b)

		case b == cMap:
		case b == cMapStrInt:
			*val = d.decodeMapStrInt(b)
		case b == cMapStrBool:
			*val = d.decodeMapStrBool(b)
		case b == cMapStrStr:
			*val = d.decodeMapStrStr(b)
		case b == cMapStrNil:
			*val = d.decodeMapStrNil(b)
		case b == cMapNilNil:
			*val = d.decodeMapNilNil(b)
		}

	default:

		if reflect.TypeOf(dst).Kind() == reflect.Ptr {
			item := reflect.ValueOf(dst).Elem()
			kind := item.Type()
			d.decodeStructMap(b, kind, item)
		}

	}

	return

}

func (d *Decoder) decodeBit() (val byte) {
	return d.r.ReadOne()
}

func (d *Decoder) decodeVal(b byte) (val bool) {
	switch {
	case b == cTrue:
		return true
	case b == cFalse:
		return false
	}
	return
}

func (d *Decoder) decodeBin(b byte) (val []byte) {
	var sze int
	switch {
	case b >= cFixBin && b <= cFixBin+0x1F:
		sze = int(b - cFixBin)
	case b == cBin8:
		var tmp int8
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cBin16:
		var tmp int16
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cBin32:
		var tmp int32
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cBin64:
		var tmp int64
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	}
	bin := d.r.ReadMany(sze)
	return bin
}

func (d *Decoder) decodeStr(b byte) (val string) {
	var sze int
	switch {
	case b >= cFixStr && b <= cFixStr+0x1F:
		sze = int(b - cFixStr)
	case b == cStr8:
		var tmp int8
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr16:
		var tmp int16
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr32:
		var tmp int32
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr64:
		var tmp int64
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	}
	bin := d.r.ReadMany(sze)
	return string(bin)
}

func (d *Decoder) decodeExt(b byte) (val interface{}) {

	var sze int

	switch {
	case b >= cFixStr && b <= cFixStr+0x1F:
		sze = int(b - cFixStr)
	case b == cStr8:
		var tmp int8
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr16:
		var tmp int16
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr32:
		var tmp int32
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	case b == cStr64:
		var tmp int64
		binary.Read(d.r, binary.BigEndian, &tmp)
		sze = int(tmp)
	}

	bin := d.r.ReadMany(sze)

	return string(bin)

}

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

func (d *Decoder) decodeArrNil(b byte) (val []interface{}) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]interface{}, tot)
	for i := 0; i < tot; i++ {
		d.decode(&arr[i])
	}
	return arr
}

func (d *Decoder) decodeArrBool(b byte) (val []bool) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]bool, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = d.decodeVal(nxt)
	}
	return arr
}

func (d *Decoder) decodeArrStr(b byte) (val []string) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]string, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = d.decodeStr(nxt)
	}
	return arr
}

func (d *Decoder) decodeArrInt(b byte) (val []int) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]int, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = int(d.decodeInt(nxt))
	}
	return arr
}

func (d *Decoder) decodeArrInt8(b byte) (val []int8) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]int8, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = int8(d.decodeInt(nxt))
	}
	return arr
}

func (d *Decoder) decodeArrInt16(b byte) (val []int16) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]int16, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = int16(d.decodeInt(nxt))
	}
	return arr
}

func (d *Decoder) decodeArrInt32(b byte) (val []int32) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]int32, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = int32(d.decodeInt(nxt))
	}
	return arr
}

func (d *Decoder) decodeArrInt64(b byte) (val []int64) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]int64, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = int64(d.decodeInt(nxt))
	}
	return arr
}

func (d *Decoder) decodeArrUint(b byte) (val []uint) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]uint, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = uint(d.decodeUint(nxt))
	}
	return arr
}

func (d *Decoder) decodeArrUint16(b byte) (val []uint16) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]uint16, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = uint16(d.decodeUint(nxt))
	}
	return arr
}

func (d *Decoder) decodeArrUint32(b byte) (val []uint32) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]uint32, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = uint32(d.decodeUint(nxt))
	}
	return arr
}

func (d *Decoder) decodeArrUint64(b byte) (val []uint64) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]uint64, tot)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		arr[i] = uint64(d.decodeUint(nxt))
	}
	return arr
}

func (d *Decoder) decodeArrFloat32(b byte) (val []float32) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]float32, tot)
	for i := 0; i < tot; i++ {
		arr[i] = d.decodeFloat32(nxt)
	}
	return arr
}

func (d *Decoder) decodeArrFloat64(b byte) (val []float64) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]float64, tot)
	for i := 0; i < tot; i++ {
		arr[i] = d.decodeFloat64(nxt)
	}
	return arr
}

func (d *Decoder) decodeArrTime(b byte) (val []time.Time) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	arr := make([]time.Time, tot)
	for i := 0; i < tot; i++ {
		arr[i] = d.decodeTime(nxt)
	}
	return arr
}

func (d *Decoder) decodeMapStrInt(b byte) (val map[string]int) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]int)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		k := d.decodeStr(nxt)
		nxt = d.decodeBit()
		v := d.decodeInt(nxt)
		obj[k] = v
	}
	return obj
}

func (d *Decoder) decodeMapStrBool(b byte) (val map[string]bool) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]bool)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		k := d.decodeStr(nxt)
		nxt = d.decodeBit()
		v := d.decodeVal(nxt)
		obj[k] = v
	}
	return obj
}

func (d *Decoder) decodeMapStrStr(b byte) (val map[string]string) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]string)
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		k := d.decodeStr(nxt)
		nxt = d.decodeBit()
		v := d.decodeStr(nxt)
		obj[k] = v
	}
	return obj
}

func (d *Decoder) decodeMapStrNil(b byte) (val map[string]interface{}) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]interface{})
	for i := 0; i < tot; i++ {
		nxt = d.decodeBit()
		k := d.decodeStr(nxt)
		var v interface{}
		d.decode(&v)
		obj[k] = v
	}
	return obj
}

func (d *Decoder) decodeMapNilNil(b byte) (val map[interface{}]interface{}) {
	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[interface{}]interface{})
	for i := 0; i < tot; i++ {
		var k interface{}
		d.decode(&k)
		var v interface{}
		d.decode(&v)
		obj[k] = v
	}
	return obj
}

func (d *Decoder) decodeStructMap(b byte, kind reflect.Type, item reflect.Value) {

	nxt := d.decodeBit()
	tot := d.decodeInt(nxt)
	obj := make(map[string]*field)

	for i := 0; i < item.NumField(); i++ {
		if fld := newField(kind.Field(i), item.Field(i)); fld != nil {
			obj[fld.Show()] = fld
		}
	}

	for i := 0; i < tot; i++ {

		nxt = d.decodeBit()
		k := d.decodeStr(nxt)
		var v interface{}
		d.decode(&v)

		if itm, ok := obj[k]; ok {
			fld := item.FieldByName(itm.Name())
			if fld.CanSet() {

				l := fld.Type().Kind()
				r := reflect.TypeOf(v).Kind()

				switch l {
				case r:
					fld.Set(reflect.ValueOf(v))
				case reflect.Int:
					switch o := v.(type) {
					case int:
						fld.Set(reflect.ValueOf(int(o)))
					case int8:
						fld.Set(reflect.ValueOf(int(o)))
					case int16:
						fld.Set(reflect.ValueOf(int(o)))
					case int32:
						fld.Set(reflect.ValueOf(int(o)))
					}
				case reflect.Int16:
					switch o := v.(type) {
					case int8:
						fld.Set(reflect.ValueOf(int16(o)))
					}
				case reflect.Int32:
					switch o := v.(type) {
					case int8:
						fld.Set(reflect.ValueOf(int32(o)))
					case int16:
						fld.Set(reflect.ValueOf(int32(o)))
					}
				case reflect.Int64:
					switch o := v.(type) {
					case int:
						fld.Set(reflect.ValueOf(int64(o)))
					case int8:
						fld.Set(reflect.ValueOf(int64(o)))
					case int16:
						fld.Set(reflect.ValueOf(int64(o)))
					case int32:
						fld.Set(reflect.ValueOf(int64(o)))
					}
				case reflect.Uint:
					switch o := v.(type) {
					case uint:
						fld.Set(reflect.ValueOf(uint(o)))
					case uint8:
						fld.Set(reflect.ValueOf(uint(o)))
					case uint16:
						fld.Set(reflect.ValueOf(uint(o)))
					case uint32:
						fld.Set(reflect.ValueOf(uint(o)))
					}
				case reflect.Uint16:
					switch o := v.(type) {
					case uint8:
						fld.Set(reflect.ValueOf(uint16(o)))
					}
				case reflect.Uint32:
					switch o := v.(type) {
					case uint8:
						fld.Set(reflect.ValueOf(uint32(o)))
					case uint16:
						fld.Set(reflect.ValueOf(uint32(o)))
					}
				case reflect.Uint64:
					switch o := v.(type) {
					case uint:
						fld.Set(reflect.ValueOf(uint64(o)))
					case uint8:
						fld.Set(reflect.ValueOf(uint64(o)))
					case uint16:
						fld.Set(reflect.ValueOf(uint64(o)))
					case uint32:
						fld.Set(reflect.ValueOf(uint64(o)))
					}
				}

			}
		}

	}

}
