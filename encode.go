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
	"fmt"
	"io"
	"math"
	"reflect"
	"time"
)

// Encoder represents a CORK encoder.
type Encoder struct {
	h *Handle
	w *writer
}

// Encode encodes a data object into a CORK.
func Encode(src interface{}) (dst []byte) {
	buf := bytes.NewBuffer(dst)
	NewEncoder(buf).Encode(src)
	return buf.Bytes()
}

// NewEncoder returns an Encoder for encoding into an io.Writer.
func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		w: newWriter(w),
		h: &Handle{
			Precision: false,
		},
	}
}

// Options sets the configuration options that the Encoder should use.
func (e *Encoder) Options(h *Handle) *Encoder {
	e.h = h
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
	err = codec.NewEncoder(buf).Encode(s)

	// Encoding a struct
	var t &Tester{Name: "Temp", Size: 0}
	buf := bytes.NewBuffer(nil)
	err = codec.NewEncoder(buf).Encode(t)

*/
func (e *Encoder) Encode(src interface{}) (err error) {
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
	e.encode(src)
	return
}

func (e *Encoder) encode(src interface{}) {

	switch val := src.(type) {

	case nil:
		e.encodeBit(cNil)

	case bool:
		e.encodeVal(val)

	case []byte:
		e.encodeBin(val)

	case string:
		e.encodeStr(val)

	case int:
		e.encodeInt(val)

	case int8:
		e.encodeInt8(val)

	case int16:
		e.encodeInt16(val)

	case int32:
		e.encodeInt32(val)

	case int64:
		e.encodeInt64(val)

	case uint:
		e.encodeUint(val)

	case uint8:
		e.encodeUint8(val)

	case uint16:
		e.encodeUint16(val)

	case uint32:
		e.encodeUint32(val)

	case uint64:
		e.encodeUint64(val)

	case float32:
		e.encodeFloat32(val)

	case float64:
		e.encodeFloat64(val)

	case complex64:
		e.encodeComplex64(val)

	case complex128:
		e.encodeComplex128(val)

	case time.Time:
		e.encodeTime(val)

	// ---------------------------------------------
	// Include common slice types
	// ---------------------------------------------

	case []bool:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeVal(v)
		}

	case []string:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeStr(v)
		}

	case []int:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeInt(v)
		}

	case []int8:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeInt8(v)
		}

	case []int16:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeInt16(v)
		}

	case []int32:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeInt32(v)
		}

	case []int64:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeInt64(v)
		}

	case []uint:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeUint(v)
		}

	case []uint16:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeUint16(v)
		}

	case []uint32:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeUint32(v)
		}

	case []uint64:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeUint64(v)
		}

	case []float32:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeFloat32(v)
		}

	case []float64:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeFloat64(v)
		}

	case []complex64:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeComplex64(v)
		}

	case []complex128:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeComplex128(v)
		}

	case []time.Time:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeTime(v)
		}

	case []interface{}:
		e.encodeBit(cArr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encode(v)
		}

	// ---------------------------------------------
	// Include common map[string]<T> types
	// ---------------------------------------------

	case map[string]int:
		e.encodeBit(cMap)
		e.encodeLen(len(val))
		for k, v := range val {
			e.encodeStr(k)
			e.encodeInt(v)
		}

	case map[string]uint:
		e.encodeBit(cMap)
		e.encodeLen(len(val))
		for k, v := range val {
			e.encodeStr(k)
			e.encodeUint(v)
		}

	case map[string]bool:
		e.encodeBit(cMap)
		e.encodeLen(len(val))
		for k, v := range val {
			e.encodeStr(k)
			e.encodeVal(v)
		}

	case map[string]string:
		e.encodeBit(cMap)
		e.encodeLen(len(val))
		for k, v := range val {
			e.encodeStr(k)
			e.encodeStr(v)
		}

	// ---------------------------------------------
	// Include common map[<T>]interface{} types
	// ---------------------------------------------

	case map[int]interface{}:
		e.encodeBit(cMap)
		e.encodeLen(len(val))
		for k, v := range val {
			e.encodeInt(k)
			e.encode(v)
		}

	case map[uint]interface{}:
		e.encodeBit(cMap)
		e.encodeLen(len(val))
		for k, v := range val {
			e.encodeUint(k)
			e.encode(v)
		}

	case map[string]interface{}:
		e.encodeBit(cMap)
		e.encodeLen(len(val))
		for k, v := range val {
			e.encodeStr(k)
			e.encode(v)
		}

	// ---------------------------------------------
	// Include map[interface{}]interface{} type
	// ---------------------------------------------

	case map[interface{}]interface{}:
		e.encodeBit(cMap)
		e.encodeLen(len(val))
		for k, v := range val {
			e.encode(k)
			e.encode(v)
		}

	// ---------------------------------------------
	// Include map[interface{}]interface{} type
	// ---------------------------------------------

	case Corker:
		e.encodeExt(val)

	// ---------------------------------------------
	// Use reflect for any remaining types
	// ---------------------------------------------

	default:

		item := reflect.ValueOf(src)
		kind := reflect.TypeOf(src)

		switch kind.Kind() {

		case reflect.Ptr:
			item := item.Elem()
			if !item.IsValid() {
				e.encodeBit(cNil)
				return
			}
			e.encode(item.Interface())

		case reflect.Struct:
			flds := make([]*field, 0)
			for i := 0; i < item.NumField(); i++ {
				if fld := newField(kind.Field(i), item.Field(i)); fld != nil {
					flds = append(flds, fld)
				}
			}
			e.encodeBit(cMap)
			e.encodeLen(len(flds))
			for _, fld := range flds {
				e.encode(fld.Show())
				e.encode(item.FieldByName(fld.Name()).Interface())
			}

		case reflect.Slice:
			e.encodeArr(src, kind, item)

		case reflect.Map:
			e.encodeMap(src, kind, item)

		case reflect.Bool:
			e.encode(bool(item.Bool()))

		case reflect.String:
			e.encode(string(item.String()))

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			e.encode(item.Int())

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			e.encode(item.Uint())

		case reflect.Float32:
			e.encode(float32(item.Float()))

		case reflect.Float64:
			e.encode(float64(item.Float()))

		case reflect.Complex64:
			e.encode(complex64(item.Complex()))

		case reflect.Complex128:
			e.encode(complex128(item.Complex()))

		default:
			e.encodeBit(cNil)

		}

	}

}

func (e *Encoder) encodeBit(val byte) {
	e.w.WriteOne(val)
	return
}

func (e *Encoder) encodeVal(val bool) {
	if val {
		e.encodeBit(cTrue)
	} else {
		e.encodeBit(cFalse)
	}
	return
}

func (e *Encoder) encodeBin(val []byte) {
	sze := len(val)
	switch {
	case sze <= fixedBin:
		e.encodeBit(cFixBin + byte(sze))
	case sze <= math.MaxUint8:
		e.encodeBit(cBin8)
		e.encodeLen8(uint8(sze))
	case sze <= math.MaxUint16:
		e.encodeBit(cBin16)
		e.encodeLen16(uint16(sze))
	case sze <= math.MaxUint32:
		e.encodeBit(cBin32)
		e.encodeLen32(uint32(sze))
	case sze <= math.MaxInt64:
		e.encodeBit(cBin64)
		e.encodeLen64(uint64(sze))
	}
	e.w.WriteMany(val)
	return
}

func (e *Encoder) encodeStr(val string) {
	sze := len(val)
	switch {
	case sze <= fixedStr:
		e.encodeBit(cFixStr + byte(sze))
	case sze <= math.MaxUint8:
		e.encodeBit(cStr8)
		e.encodeLen8(uint8(sze))
	case sze <= math.MaxUint16:
		e.encodeBit(cStr16)
		e.encodeLen16(uint16(sze))
	case sze <= math.MaxUint32:
		e.encodeBit(cStr32)
		e.encodeLen32(uint32(sze))
	case sze <= math.MaxInt64:
		e.encodeBit(cStr64)
		e.encodeLen64(uint64(sze))
	}
	e.w.WriteText(val)
	return
}

func (e *Encoder) encodeExt(val Corker) {
	bit := val.ExtendCORK()
	enc, err := val.MarshalCORK()
	if err != nil {
		panic(err)
	}
	sze := len(enc)
	switch {
	case sze <= fixedExt:
		e.encodeBit(cFixExt + byte(sze))
	case sze <= math.MaxUint8:
		e.encodeBit(cExt8)
		e.encodeLen8(uint8(sze))
	case sze <= math.MaxUint16:
		e.encodeBit(cExt16)
		e.encodeLen16(uint16(sze))
	case sze <= math.MaxUint32:
		e.encodeBit(cExt32)
		e.encodeLen32(uint32(sze))
	case sze <= math.MaxInt64:
		e.encodeBit(cExt64)
		e.encodeLen64(uint64(sze))
	}
	e.w.WriteOne(bit)
	e.w.WriteMany(enc)
	return
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func (e *Encoder) encodeLen(val int) {
	e.encodeUint(uint(val))
	return
}

// --------------------------------------------------

func (e *Encoder) encodeLen8(val uint8) {
	e.encodeBit(byte(val))
}

func (e *Encoder) encodeLen16(val uint16) {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	e.w.Write(buf)
}

func (e *Encoder) encodeLen32(val uint32) {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	e.w.Write(buf)
}

func (e *Encoder) encodeLen64(val uint64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	e.w.Write(buf)
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func (e *Encoder) encodeInt(val int) {
	switch {
	case val >= 0 && val <= fixedInt:
		e.encodeInt1Fixed(int8(val))
	case val <= math.MaxInt8:
		e.encodeInt8Fixed(int8(val))
	case val <= math.MaxInt16:
		e.encodeInt16Fixed(int16(val))
	case val <= math.MaxInt32:
		e.encodeInt32Fixed(int32(val))
	case val <= math.MaxInt64:
		e.encodeInt64Fixed(int64(val))
	}
	return
}

// --------------------------------------------------

func (e *Encoder) encodeInt8(val int8) {
	if e.h.Precision {
		e.encodeInt8Fixed(val)
	} else {
		e.encodeInt(int(val))
	}
}

func (e *Encoder) encodeInt16(val int16) {
	if e.h.Precision {
		e.encodeInt16Fixed(val)
	} else {
		e.encodeInt(int(val))
	}
}

func (e *Encoder) encodeInt32(val int32) {
	if e.h.Precision {
		e.encodeInt32Fixed(val)
	} else {
		e.encodeInt(int(val))
	}
}

func (e *Encoder) encodeInt64(val int64) {
	if e.h.Precision {
		e.encodeInt64Fixed(val)
	} else {
		e.encodeInt(int(val))
	}
}

// --------------------------------------------------

func (e *Encoder) encodeInt1Fixed(val int8) {
	e.encodeBit(byte(val))
}

func (e *Encoder) encodeInt8Fixed(val int8) {
	e.encodeBit(cInt8)
	e.encodeBit(byte(val))
}

func (e *Encoder) encodeInt16Fixed(val int16) {
	e.encodeBit(cInt16)
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(val))
	e.w.Write(buf)
}

func (e *Encoder) encodeInt32Fixed(val int32) {
	e.encodeBit(cInt32)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(val))
	e.w.Write(buf)
}

func (e *Encoder) encodeInt64Fixed(val int64) {
	e.encodeBit(cInt64)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(val))
	e.w.Write(buf)
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func (e *Encoder) encodeUint(val uint) {
	switch {
	case val >= 0 && val <= fixedInt:
		e.encodeUint1Fixed(uint8(val))
	case val <= math.MaxUint8:
		e.encodeUint8Fixed(uint8(val))
	case val <= math.MaxUint16:
		e.encodeUint16Fixed(uint16(val))
	case val <= math.MaxUint32:
		e.encodeUint32Fixed(uint32(val))
	case val <= math.MaxUint64:
		e.encodeUint64Fixed(uint64(val))
	}
	return
}

// --------------------------------------------------

func (e *Encoder) encodeUint8(val uint8) {
	if e.h.Precision {
		e.encodeUint8Fixed(val)
	} else {
		e.encodeUint(uint(val))
	}
}

func (e *Encoder) encodeUint16(val uint16) {
	if e.h.Precision {
		e.encodeUint16Fixed(val)
	} else {
		e.encodeUint(uint(val))
	}
}

func (e *Encoder) encodeUint32(val uint32) {
	if e.h.Precision {
		e.encodeUint32Fixed(val)
	} else {
		e.encodeUint(uint(val))
	}
}

func (e *Encoder) encodeUint64(val uint64) {
	if e.h.Precision {
		e.encodeUint64Fixed(val)
	} else {
		e.encodeUint(uint(val))
	}
}

// --------------------------------------------------

func (e *Encoder) encodeUint1Fixed(val uint8) {
	e.encodeBit(byte(val))
}

func (e *Encoder) encodeUint8Fixed(val uint8) {
	e.encodeBit(cUint8)
	e.encodeBit(byte(val))
}

func (e *Encoder) encodeUint16Fixed(val uint16) {
	e.encodeBit(cUint16)
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	e.w.Write(buf)
}

func (e *Encoder) encodeUint32Fixed(val uint32) {
	e.encodeBit(cUint32)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	e.w.Write(buf)
}

func (e *Encoder) encodeUint64Fixed(val uint64) {
	e.encodeBit(cUint64)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	e.w.Write(buf)
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func (e *Encoder) encodeFloat32(val float32) {
	e.encodeBit(cFloat32)
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, math.Float32bits(val))
	e.w.Write(buf)
	return
}

func (e *Encoder) encodeFloat64(val float64) {
	e.encodeBit(cFloat64)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, math.Float64bits(val))
	e.w.Write(buf)
	return
}

func (e *Encoder) encodeComplex64(val complex64) {
	e.encodeBit(cComplex64)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint32(buf[:4], math.Float32bits(float32(real(val))))
	binary.BigEndian.PutUint32(buf[4:], math.Float32bits(float32(imag(val))))
	e.w.Write(buf)
	return
}

func (e *Encoder) encodeComplex128(val complex128) {
	e.encodeBit(cComplex128)
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[:8], math.Float64bits(float64(real(val))))
	binary.BigEndian.PutUint64(buf[8:], math.Float64bits(float64(imag(val))))
	e.w.Write(buf)
	return
}

func (e *Encoder) encodeTime(val time.Time) {
	e.encodeBit(cTime)
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(val.UTC().UnixNano()))
	e.w.Write(buf)
	return
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func (e *Encoder) encodeArr(src interface{}, kind reflect.Type, item reflect.Value) {
	e.encodeBit(cArr)
	e.encodeLen(item.Len())
	for i := 0; i < reflect.ValueOf(src).Len(); i++ {
		e.encode(item.Index(i).Interface())
	}
}

func (e *Encoder) encodeMap(src interface{}, kind reflect.Type, item reflect.Value) {
	e.encodeBit(cMap)
	e.encodeLen(item.Len())
	for _, k := range item.MapKeys() {
		e.encode(k.Interface())
		e.encode(item.MapIndex(k).Interface())
	}
}
