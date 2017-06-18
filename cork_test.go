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
	"errors"
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

// ----------------------------------------------------------------------

var handle = &Handle{
	ArrType: make([]interface{}, 0),
	MapType: make(map[string]interface{}),
}

// ----------------------------------------------------------------------

type CustomAny interface{}
type CustomBool bool
type CustomInt int
type CustomInt8 int8
type CustomInt16 int16
type CustomInt32 int32
type CustomInt64 int64
type CustomUint uint
type CustomUint8 uint8
type CustomUint16 uint16
type CustomUint32 uint32
type CustomUint64 uint64
type CustomFloat32 float32
type CustomFloat64 float64
type CustomString string

type Custom struct {
	Nil             interface{}
	Null            interface{}
	Bool            bool
	String          string
	Bytes           []byte
	Time            time.Time
	Int             int
	Int8            int8
	Int16           int16
	Int32           int32
	Int64           int64
	Uint            uint
	Uint8           uint8
	Uint16          uint16
	Uint32          uint32
	Uint64          uint64
	Float32         float32
	Float64         float64
	Complex64       complex64
	Complex128      complex128
	Any             interface{}
	ArrBool         []bool
	ArrString       []string
	ArrInt          []int
	ArrInt8         []int8
	ArrInt16        []int16
	ArrInt32        []int32
	ArrInt64        []int64
	ArrUint         []uint
	ArrUint8        []uint8
	ArrUint16       []uint16
	ArrUint32       []uint32
	ArrUint64       []uint64
	ArrFloat32      []float32
	ArrFloat64      []float64
	ArrComplex64    []complex64
	ArrComplex128   []complex128
	ArrTime         []time.Time
	ArrAny          []interface{}
	ArrNo           string
	MapStringInt    map[string]int
	MapStringUint   map[string]uint
	MapStringBool   map[string]bool
	MapStringString map[string]string
	MapIntAny       map[int]interface{}
	MapUintAny      map[uint]interface{}
	MapStringAny    map[string]interface{}
	MapTimeAny      map[time.Time]interface{}
	MapAnyAny       map[interface{}]interface{}
	MapNo           string
	CustomBool      CustomBool
	CustomString    CustomString
	CustomInt       CustomInt
	CustomInt8      CustomInt8
	CustomInt16     CustomInt16
	CustomInt32     CustomInt32
	CustomInt64     CustomInt64
	CustomUint      CustomUint
	CustomUint8     CustomUint8
	CustomUint16    CustomUint16
	CustomUint32    CustomUint32
	CustomUint64    CustomUint64
	CustomFloat32   CustomFloat32
	CustomFloat64   CustomFloat64
	CustomAny       CustomAny
	Tignored        *Tested
	Signored        *Selfed
	Cignored        *Corked
	Testable        *Tested
	Selfable        *Selfed
	Corkable        *Corked
	Funcable        func()
	Chanable        chan int
	Drrayble        []Corked
	Arrayble        []*Corked
	Embedded        struct {
		One string
		Two int
		Ced *Corked
		Sed *Selfed
	}
}

// ----------------------------------------------------------------------

type CustomTextFailer struct {
	Field string
}

func (this *CustomTextFailer) MarshalText() ([]byte, error) {
	return nil, errors.New("Marshal error")
}

func (this *CustomTextFailer) UnmarshalText(v []byte) error {
	return errors.New("Unmarshal error")
}

type CustomBinaryFailer struct {
	Field string
}

func (this *CustomBinaryFailer) MarshalBinary() ([]byte, error) {
	return nil, errors.New("Marshal error")
}

func (this *CustomBinaryFailer) UnmarshalBinary(v []byte) error {
	return errors.New("Unmarshal error")
}

// ----------------------------------------------------------------------

type CustomTextMarshaler struct {
	Field string
}

func (this *CustomTextMarshaler) MarshalText() ([]byte, error) {
	return []byte(this.Field), nil
}

func (this *CustomTextMarshaler) UnmarshalText(v []byte) error {
	this.Field = string(v)
	return nil
}

type CustomBinaryMarshaler struct {
	Field string
}

func (this *CustomBinaryMarshaler) MarshalBinary() ([]byte, error) {
	return []byte(this.Field), nil
}

func (this *CustomBinaryMarshaler) UnmarshalBinary(v []byte) error {
	this.Field = string(v)
	return nil
}

// ----------------------------------------------------------------------

type Tested struct {
	Name  string
	Data  []byte `cork:"data"`
	Temp  []string
	Test  map[string]string
	priv  bool
	Count int
	Omit  string `cork:"-"`
	Empty string `cork:",omitempty"`
}

// ----------------------------------------------------------------------

type Errord struct{}

func (this *Errord) ExtendCORK() byte {
	return 0x00
}

func (this *Errord) MarshalCORK() (dst []byte, err error) {
	return nil, errors.New("Marshal error")
}

func (this *Errord) UnmarshalCORK(src []byte) (err error) {
	return errors.New("Unmarshal error")
}

// ----------------------------------------------------------------------

type Simple struct{}

func (this *Simple) ExtendCORK() byte {
	return 0x01
}

func (this *Simple) MarshalCORK() (dst []byte, err error) {
	return
}

func (this *Simple) UnmarshalCORK(src []byte) (err error) {
	return
}

// ----------------------------------------------------------------------

type Corked struct {
	Name  string
	Data  []byte   `cork:"data"`
	Temp  []string `cork:"-"`
	Test  map[string]string
	priv  bool
	Count int
	Omit  string `cork:"-"`
	Empty string `cork:",omitempty"`
}

func (this *Corked) ExtendCORK() byte {
	return 0x02
}

func (this *Corked) MarshalCORK() (dst []byte, err error) {
	b := bytes.NewBuffer(dst)
	e := NewEncoder(b)
	e.Encode(this.Name)
	e.Encode(this.Data)
	e.Encode(this.Temp)
	e.Encode(this.Test)
	e.Encode(this.Count)
	return b.Bytes(), nil
}

func (this *Corked) UnmarshalCORK(src []byte) (err error) {
	b := bytes.NewBuffer(src)
	d := NewDecoder(b)
	d.Decode(&this.Name)
	d.Decode(&this.Data)
	d.Decode(&this.Temp)
	d.Decode(&this.Test)
	d.Decode(&this.Count)
	return
}

// ----------------------------------------------------------------------

type Selfed struct {
	Name  string
	Data  []byte   `cork:"data"`
	Temp  []string `cork:"-"`
	Test  map[string]string
	priv  bool
	Count int
	Omit  string `cork:"-"`
	Empty string `cork:",omitempty"`
}

func (this *Selfed) ExtendCORK() byte {
	return 0x03
}

func (this *Selfed) MarshalCORK(w *Writer) (err error) {
	w.EncodeString(this.Name)
	w.EncodeBytes(this.Data)
	w.EncodeArr(this.Temp)
	w.EncodeMap(this.Test)
	w.EncodeInt(this.Count)
	return
}

func (this *Selfed) UnmarshalCORK(r *Reader) (err error) {
	r.DecodeString(&this.Name)
	r.DecodeBytes(&this.Data)
	r.DecodeArr(&this.Temp)
	r.DecodeMap(&this.Test)
	r.DecodeInt(&this.Count)
	return
}

// ----------------------------------------------------------------------

type Complex struct {
	Bool            bool
	String          string
	Bytes           []byte
	Time            time.Time
	Int             int
	Int8            int8
	Int16           int16
	Int32           int32
	Int64           int64
	Uint            uint
	Uint8           uint8
	Uint16          uint16
	Uint32          uint32
	Uint64          uint64
	Float32         float32
	Float64         float64
	Complex64       complex64
	Complex128      complex128
	Any             interface{}
	ArrBool         []bool
	ArrString       []string
	ArrInt          []int
	ArrInt8         []int8
	ArrInt16        []int16
	ArrInt32        []int32
	ArrInt64        []int64
	ArrUint         []uint
	ArrUint8        []uint8
	ArrUint16       []uint16
	ArrUint32       []uint32
	ArrUint64       []uint64
	ArrFloat32      []float32
	ArrFloat64      []float64
	ArrComplex64    []complex64
	ArrComplex128   []complex128
	ArrTime         []time.Time
	ArrAny          []interface{}
	ArrNo           string
	MapStringInt    map[string]int
	MapStringUint   map[string]uint
	MapStringBool   map[string]bool
	MapStringString map[string]string
	MapIntAny       map[int]interface{}
	MapUintAny      map[uint]interface{}
	MapStringAny    map[string]interface{}
	MapTimeAny      map[time.Time]interface{}
	MapAnyAny       map[interface{}]interface{}
	MapNo           string
}

func (this *Complex) ExtendCORK() byte {
	return 0x04
}

func (this *Complex) MarshalCORK(w *Writer) (err error) {
	w.EncodeBool(this.Bool)
	w.EncodeString(this.String)
	w.EncodeBytes(this.Bytes)
	w.EncodeTime(this.Time)
	w.EncodeInt(this.Int)
	w.EncodeInt8(this.Int8)
	w.EncodeInt16(this.Int16)
	w.EncodeInt32(this.Int32)
	w.EncodeInt64(this.Int64)
	w.EncodeUint(this.Uint)
	w.EncodeUint8(this.Uint8)
	w.EncodeUint16(this.Uint16)
	w.EncodeUint32(this.Uint32)
	w.EncodeUint64(this.Uint64)
	w.EncodeFloat32(this.Float32)
	w.EncodeFloat64(this.Float64)
	w.EncodeComplex64(this.Complex64)
	w.EncodeComplex128(this.Complex128)
	w.EncodeAny(this.Any)
	w.EncodeArr(this.ArrBool)
	w.EncodeArr(this.ArrString)
	w.EncodeArr(this.ArrInt)
	w.EncodeArr(this.ArrInt8)
	w.EncodeArr(this.ArrInt16)
	w.EncodeArr(this.ArrInt32)
	w.EncodeArr(this.ArrInt64)
	w.EncodeArr(this.ArrUint)
	w.EncodeArr(this.ArrUint8)
	w.EncodeArr(this.ArrUint16)
	w.EncodeArr(this.ArrUint32)
	w.EncodeArr(this.ArrUint64)
	w.EncodeArr(this.ArrFloat32)
	w.EncodeArr(this.ArrFloat64)
	w.EncodeArr(this.ArrComplex64)
	w.EncodeArr(this.ArrComplex128)
	w.EncodeArr(this.ArrTime)
	w.EncodeArr(this.ArrAny)
	w.EncodeArr(this.ArrNo)
	w.EncodeMap(this.MapStringInt)
	w.EncodeMap(this.MapStringUint)
	w.EncodeMap(this.MapStringBool)
	w.EncodeMap(this.MapStringString)
	w.EncodeMap(this.MapIntAny)
	w.EncodeMap(this.MapUintAny)
	w.EncodeMap(this.MapStringAny)
	w.EncodeMap(this.MapTimeAny)
	w.EncodeMap(this.MapAnyAny)
	w.EncodeMap(this.MapNo)
	return
}

func (this *Complex) UnmarshalCORK(r *Reader) (err error) {
	r.DecodeBool(&this.Bool)
	r.DecodeString(&this.String)
	r.DecodeBytes(&this.Bytes)
	r.DecodeTime(&this.Time)
	r.DecodeInt(&this.Int)
	r.DecodeInt8(&this.Int8)
	r.DecodeInt16(&this.Int16)
	r.DecodeInt32(&this.Int32)
	r.DecodeInt64(&this.Int64)
	r.DecodeUint(&this.Uint)
	r.DecodeUint8(&this.Uint8)
	r.DecodeUint16(&this.Uint16)
	r.DecodeUint32(&this.Uint32)
	r.DecodeUint64(&this.Uint64)
	r.DecodeFloat32(&this.Float32)
	r.DecodeFloat64(&this.Float64)
	r.DecodeComplex64(&this.Complex64)
	r.DecodeComplex128(&this.Complex128)
	r.DecodeAny(&this.Any)
	r.DecodeArr(&this.ArrBool)
	r.DecodeArr(&this.ArrString)
	r.DecodeArr(&this.ArrInt)
	r.DecodeArr(&this.ArrInt8)
	r.DecodeArr(&this.ArrInt16)
	r.DecodeArr(&this.ArrInt32)
	r.DecodeArr(&this.ArrInt64)
	r.DecodeArr(&this.ArrUint)
	r.DecodeArr(&this.ArrUint8)
	r.DecodeArr(&this.ArrUint16)
	r.DecodeArr(&this.ArrUint32)
	r.DecodeArr(&this.ArrUint64)
	r.DecodeArr(&this.ArrFloat32)
	r.DecodeArr(&this.ArrFloat64)
	r.DecodeArr(&this.ArrComplex64)
	r.DecodeArr(&this.ArrComplex128)
	r.DecodeArr(&this.ArrTime)
	r.DecodeArr(&this.ArrAny)
	r.DecodeArr(&this.ArrNo)
	r.DecodeMap(&this.MapStringInt)
	r.DecodeMap(&this.MapStringUint)
	r.DecodeMap(&this.MapStringBool)
	r.DecodeMap(&this.MapStringString)
	r.DecodeMap(&this.MapIntAny)
	r.DecodeMap(&this.MapUintAny)
	r.DecodeMap(&this.MapStringAny)
	r.DecodeMap(&this.MapTimeAny)
	r.DecodeMap(&this.MapAnyAny)
	r.DecodeMap(&this.MapNo)
	return
}

// ----------------------------------------------------------------------

var str = "This is the very last time that I should have to test this"

var bin = []byte{
	84, 104, 105, 115, 32, 105, 115, 32, 116, 104, 101, 32, 118, 101,
	114, 121, 32, 108, 97, 115, 116, 32, 116, 105, 109, 101, 32, 116,
	104, 97, 116, 32, 73, 32, 115, 104, 111, 117, 108, 100, 32, 104, 97,
	118, 101, 32, 116, 111, 32, 116, 101, 115, 116, 32, 116, 104, 105, 115,
}

var lng = append(bin, append(bin, append(bin, append(bin, bin...)...)...)...)

// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------
// ----------------------------------------------------------------------

func init() {
	Register(&Simple{})
	Register(&Corked{})
	Register(&Selfed{})
	Register(&Complex{})
}

func TestGeneral(t *testing.T) {

	tme, _ := time.Parse(time.RFC3339, "1987-06-22T08:00:00.123456789Z")

	Convey("nil will encode and decode", t, func() {
		var tmp interface{}
		var val = interface{}(nil)
		var bit = []byte{cNil}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("true will encode and decode", t, func() {
		var tmp bool
		var val = true
		var bit = []byte{cTrue}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("false will encode and decode", t, func() {
		var tmp bool
		var val = false
		var bit = []byte{cFalse}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("byte will encode and decode", t, func() {
		var tmp byte
		var val = byte('a')
		var bit = []byte{97}
		var oth = int(97)
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, oth)
	})

	Convey("string will encode and decode", t, func() {
		var tmp string
		var val = "Hello"
		var bit = []byte{cFixStr + 0x05, 72, 101, 108, 108, 111}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("string8 will encode and decode", t, func() {
		var tmp string
		var val = str
		var bit = append([]byte{cStr8, 58}, bin...)
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("string16 will encode and decode", t, func() {
		var tmp string
		var val = str + str + str + str + str
		var bit = append([]byte{cStr16, 1, 34}, lng...)
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("[]byte will encode and decode", t, func() {
		var tmp []byte
		var val = []byte{72, 101, 108, 108, 111}
		var bit = []byte{cFixBin + 0x05, 72, 101, 108, 108, 111}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("[]byte8 will encode and decode", t, func() {
		var tmp []byte
		var val = bin
		var bit = append([]byte{cBin8, 58}, bin...)
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("[]byte16 will encode and decode", t, func() {
		var tmp []byte
		var val = lng
		var bit = append([]byte{cBin16, 1, 34}, lng...)
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("int will encode and decode", t, func() {
		var tmp int
		var val = int(1)
		var bit = []byte{1}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("int8 will encode and decode", t, func() {
		var tmp int8
		var val = int8(math.MaxInt8)
		var bit = []byte{127}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, int(val))
	})

	Convey("int16 will encode and decode", t, func() {
		var tmp int16
		var val = int16(math.MaxInt16)
		var bit = []byte{cInt16, 127, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, int(val))
	})

	Convey("int32 will encode and decode", t, func() {
		var tmp int32
		var val = int32(math.MaxInt32)
		var bit = []byte{cInt32, 127, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, int(val))
	})

	Convey("int64 will encode and decode", t, func() {
		var tmp int64
		var val = int64(math.MaxInt64)
		var bit = []byte{cInt64, 127, 255, 255, 255, 255, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, int(val))
	})

	Convey("uint will encode and decode", t, func() {
		var tmp uint
		var val = uint(1)
		var bit = []byte{1}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, int(val))
	})

	Convey("uint16 will encode and decode", t, func() {
		var tmp uint16
		var val = uint16(math.MaxUint16)
		var bit = []byte{cUint16, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, uint(val))
	})

	Convey("uint32 will encode and decode", t, func() {
		var tmp uint32
		var val = uint32(math.MaxUint32)
		var bit = []byte{cUint32, 255, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, uint(val))
	})

	Convey("uint64 will encode and decode", t, func() {
		var tmp uint64
		var val = uint64(math.MaxUint64)
		var bit = []byte{cUint64, 255, 255, 255, 255, 255, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, uint(val))
	})

	Convey("float32 will encode and decode", t, func() {
		var tmp float32
		var val = float32(math.Pi)
		var bit = []byte{cFloat32, 64, 73, 15, 219}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("float32 will encode and decode into float64", t, func() {
		var tmp float64
		var val = float32(math.Pi)
		var gen = float64(float32(math.Pi))
		var bit = []byte{cFloat32, 64, 73, 15, 219}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, gen)
		So(dec, ShouldResemble, val)
	})

	Convey("float64 will encode and decode", t, func() {
		var tmp float64
		var val = float64(math.Pi)
		var bit = []byte{cFloat64, 64, 9, 33, 251, 84, 68, 45, 24}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("complex64 will encode and decode", t, func() {
		var tmp complex64
		var val = complex64(math.Pi)
		var bit = []byte{cComplex64, 64, 73, 15, 219, 0, 0, 0, 0}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("complex128 will encode and decode", t, func() {
		var tmp complex128
		var val = complex128(math.Pi)
		var bit = []byte{cComplex128, 64, 9, 33, 251, 84, 68, 45, 24, 0, 0, 0, 0, 0, 0, 0, 0}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("time.Time will encode and decode", t, func() {
		var tmp time.Time
		var val = tme
		var bit = []byte{cTime, 7, 166, 199, 91, 123, 67, 205, 21}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, val)
	})

	Convey("CustomInt will encode and decode", t, func() {
		var tmp CustomInt
		var val = CustomInt(666)
		var gen = int(666)
		var bit = []byte{cInt16, 2, 154}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("CustomString will encode and decode", t, func() {
		var tmp CustomString
		var val = CustomString("Hello")
		var gen = "Hello"
		var bit = []byte{cFixStr + 0x05, 72, 101, 108, 108, 111}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("Tested will encode and decode", t, func() {
		var tmp Tested
		var val = Tested{"test", []byte("test"), []string{"1", "2"}, map[string]string{"1": "2"}, false, 25, "Omitted", ""}
		var oth = Tested{"test", []byte("test"), []string{"1", "2"}, map[string]string{"1": "2"}, false, 25, "", ""}
		var gen = map[interface{}]interface{}{"Name": "test", "data": []byte("test"), "Temp": []interface{}{"1", "2"}, "Test": map[interface{}]interface{}{"1": "2"}, "Count": int(25)}
		var bit = []byte{cFixMap + 0x05, /**/
			cFixStr + 0x04, 78, 97, 109, 101, cFixStr + 0x04, 116, 101, 115, 116, /**/
			cFixStr + 0x04, 100, 97, 116, 97, cFixBin + 0x04, 116, 101, 115, 116, /**/
			cFixStr + 0x04, 84, 101, 109, 112, cFixArr + 0x02, 129, 49, 129, 50, /**/
			cFixStr + 0x04, 84, 101, 115, 116, cFixMap + 0x01, 129, 49, 129, 50, /**/
			cFixStr + 0x05, 67, 111, 117, 110, 116 /**/, 25,
		}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, oth)
		So(dec, ShouldResemble, gen)
	})

	Convey("*Tested will encode and decode", t, func() {
		var tmp Tested
		var val = Tested{"test", []byte("test"), []string{"1", "2"}, map[string]string{"1": "2"}, false, 25, "Omitted", ""}
		var oth = Tested{"test", []byte("test"), []string{"1", "2"}, map[string]string{"1": "2"}, false, 25, "", ""}
		var gen = map[interface{}]interface{}{"Name": "test", "data": []byte("test"), "Temp": []interface{}{"1", "2"}, "Test": map[interface{}]interface{}{"1": "2"}, "Count": int(25)}
		var bit = []byte{cFixMap + 0x05, /**/
			cFixStr + 0x04, 78, 97, 109, 101, cFixStr + 0x04, 116, 101, 115, 116, /**/
			cFixStr + 0x04, 100, 97, 116, 97, cFixBin + 0x04, 116, 101, 115, 116, /**/
			cFixStr + 0x04, 84, 101, 109, 112, cFixArr + 0x02, 129, 49, 129, 50, /**/
			cFixStr + 0x04, 84, 101, 115, 116, cFixMap + 0x01, 129, 49, 129, 50, /**/
			cFixStr + 0x05, 67, 111, 117, 110, 116 /**/, 25,
		}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, oth)
		So(dec, ShouldResemble, gen)
	})

	Convey("*Errord will not encode and decode", t, func() {
		var tmp Errord
		var val = &Errord{}
		var bit = []byte{cFixExt + 0x00}
		eer := NewEncoder(bytes.NewBuffer(nil)).Encode(val)
		der := NewDecoder(bytes.NewReader(bit)).Decode(&tmp)
		So(eer, ShouldNotBeNil)
		So(der, ShouldNotBeNil)
	})

	Convey("*Simple will encode and decode", t, func() {
		var tmp Simple
		var val = &Simple{}
		var gen = *val
		var bit = []byte{cFixExt, 0x01}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, gen)
		So(dec, ShouldResemble, val)
	})

	Convey("*Corked will encode and decode", t, func() {
		var tmp Corked
		var val = &Corked{"test", []byte("test"), []string{"1", "2"}, map[string]string{"1": "2"}, false, 25, "", ""}
		var gen = *val
		var bit = []byte{cExt8, 0x15, 0x02, /**/
			cFixStr + 0x04, 116, 101, 115, 116, /**/
			cFixBin + 0x04, 116, 101, 115, 116, /**/
			cFixArr + 0x02, 129, 49, 129, 50, /**/
			cFixMap + 0x01, 129, 49, 129, 50, /**/
			25,
		}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, gen)
		So(dec, ShouldResemble, val)
	})

	Convey("*Selfed will encode and decode", t, func() {
		var tmp Selfed
		var val = &Selfed{"test", []byte("test"), []string{"1", "2"}, map[string]string{"1": "2"}, false, 25, "", ""}
		var gen = *val
		var bit = []byte{cSlf, 0x03, /**/
			cFixStr + 0x04, 116, 101, 115, 116, /**/
			cFixBin + 0x04, 116, 101, 115, 116, /**/
			cFixArr + 0x02, 129, 49, 129, 50, /**/
			cFixMap + 0x01, 129, 49, 129, 50, /**/
			25,
		}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, gen)
		So(dec, ShouldResemble, val)
	})

	Convey("[]bool will encode and decode", t, func() {
		var tmp []bool
		var val = []bool{true, false}
		var gen = []interface{}{true, false}
		var bit = []byte{cFixArr + 0x02, cTrue, cFalse}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]string will encode and decode", t, func() {
		var tmp []string
		var val = []string{"Hello", "World"}
		var gen = []interface{}{"Hello", "World"}
		var bit = []byte{cFixArr + 0x02 /**/, cFixStr + 0x05, 72, 101, 108, 108, 111 /**/, cFixStr + 0x05, 87, 111, 114, 108, 100}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]int will encode and decode", t, func() {
		var tmp []int
		var val = []int{0, 1, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64}
		var gen = []interface{}{0, 1, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64}
		var bit = []byte{cFixArr + 0x06 /**/, 0 /**/, 1 /**/, 127 /**/, cInt16, 127, 255 /**/, cInt32, 127, 255, 255, 255 /**/, cInt64, 127, 255, 255, 255, 255, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]int8 will encode and decode", t, func() {
		var tmp []int8
		var val = []int8{1, math.MaxInt8}
		var gen = []interface{}{1, math.MaxInt8}
		var bit = []byte{cFixArr + 0x02 /**/, 1 /**/, 127}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]int16 will encode and decode", t, func() {
		var tmp []int16
		var val = []int16{1, math.MaxInt8, math.MaxInt16}
		var gen = []interface{}{1, math.MaxInt8, math.MaxInt16}
		var bit = []byte{cFixArr + 0x03 /**/, 1 /**/, 127 /**/, cInt16, 127, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]int32 will encode and decode", t, func() {
		var tmp []int32
		var val = []int32{1, math.MaxInt8, math.MaxInt16, math.MaxInt32}
		var gen = []interface{}{1, math.MaxInt8, math.MaxInt16, math.MaxInt32}
		var bit = []byte{cFixArr + 0x04 /**/, 1 /**/, 127 /**/, cInt16, 127, 255 /**/, cInt32, 127, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]int64 will encode and decode", t, func() {
		var tmp []int64
		var val = []int64{1, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64}
		var gen = []interface{}{1, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64}
		var bit = []byte{cFixArr + 0x05 /**/, 1 /**/, 127 /**/, cInt16, 127, 255 /**/, cInt32, 127, 255, 255, 255 /**/, cInt64, 127, 255, 255, 255, 255, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]uint will encode and decode", t, func() {
		var tmp []uint
		var val = []uint{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64}
		var gen = []interface{}{int(0), int(1), uint(math.MaxUint8), uint(math.MaxUint16), uint(math.MaxUint32), uint(math.MaxUint64)}
		var bit = []byte{cFixArr + 0x06 /**/, 0 /**/, 1, cUint8, 255 /**/, cUint16, 255, 255 /**/, cUint32, 255, 255, 255, 255 /**/, cUint64, 255, 255, 255, 255, 255, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]uint16 will encode and decode", t, func() {
		var tmp []uint16
		var val = []uint16{1, math.MaxUint8, math.MaxUint16}
		var gen = []interface{}{int(1), uint(math.MaxUint8), uint(math.MaxUint16)}
		var bit = []byte{cFixArr + 0x03 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]uint32 will encode and decode", t, func() {
		var tmp []uint32
		var val = []uint32{1, math.MaxUint8, math.MaxUint16, math.MaxUint32}
		var gen = []interface{}{int(1), uint(math.MaxUint8), uint(math.MaxUint16), uint(math.MaxUint32)}
		var bit = []byte{cFixArr + 0x04 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255 /**/, cUint32, 255, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]uint64 will encode and decode", t, func() {
		var tmp []uint64
		var val = []uint64{1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64}
		var gen = []interface{}{int(1), uint(math.MaxUint8), uint(math.MaxUint16), uint(math.MaxUint32), uint(math.MaxUint64)}
		var bit = []byte{cFixArr + 0x05 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255 /**/, cUint32, 255, 255, 255, 255 /**/, cUint64, 255, 255, 255, 255, 255, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]float32 will encode and decode", t, func() {
		var tmp []float32
		var val = []float32{math.Pi, math.Pi}
		var gen = []interface{}{float32(math.Pi), float32(math.Pi)}
		var bit = []byte{cFixArr + 0x02 /**/, cFloat32, 64, 73, 15, 219 /**/, cFloat32, 64, 73, 15, 219}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]float64 will encode and decode", t, func() {
		var tmp []float64
		var val = []float64{math.Pi, math.Pi}
		var gen = []interface{}{float64(math.Pi), float64(math.Pi)}
		var bit = []byte{cFixArr + 0x02 /**/, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24 /**/, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]complex64 will encode and decode", t, func() {
		var tmp []complex64
		var val = []complex64{math.Pi, math.Pi}
		var gen = []interface{}{complex64(math.Pi), complex64(math.Pi)}
		var bit = []byte{cFixArr + 0x02 /**/, cComplex64, 64, 73, 15, 219, 0, 0, 0, 0 /**/, cComplex64, 64, 73, 15, 219, 0, 0, 0, 0}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]complex128 will encode and decode", t, func() {
		var tmp []complex128
		var val = []complex128{math.Pi, math.Pi}
		var gen = []interface{}{complex128(math.Pi), complex128(math.Pi)}
		var bit = []byte{cFixArr + 0x02 /**/, cComplex128, 64, 9, 33, 251, 84, 68, 45, 24, 0, 0, 0, 0, 0, 0, 0, 0 /**/, cComplex128, 64, 9, 33, 251, 84, 68, 45, 24, 0, 0, 0, 0, 0, 0, 0, 0}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]time.Time will encode and decode", t, func() {
		var tmp []time.Time
		var val = []time.Time{tme, tme}
		var gen = []interface{}{tme, tme}
		var bit = []byte{cFixArr + 0x02 /**/, cTime, 7, 166, 199, 91, 123, 67, 205, 21 /**/, cTime, 7, 166, 199, 91, 123, 67, 205, 21}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]interface{} will encode and decode", t, func() {
		var tmp []interface{}
		var val = []interface{}{nil, true, false, "test", []byte("test"), int8(77), uint16(177), float64(math.Pi)}
		var gen = []interface{}{nil, true, false, "test", []byte("test"), int(77), uint(177), float64(math.Pi)}
		var bit = []byte{cFixArr + 0x08 /**/, cNil, cTrue, cFalse /**/, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixBin + 0x04, 116, 101, 115, 116 /**/, 77 /**/, cUint8, 177 /**/, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, gen)
		So(dec, ShouldResemble, gen)
	})

	Convey("[]interface{} will encode and decode", t, func() {
		var tmp []interface{}
		var val = [][][]float64{{{math.Pi}}, {{math.Pi}}}
		var gen = []interface{}{[]interface{}{[]interface{}{math.Pi}}, []interface{}{[]interface{}{math.Pi}}}
		var bit = []byte{cFixArr + 0x02 /**/, cFixArr + 0x01, cFixArr + 0x01, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24 /**/, cFixArr + 0x01, cFixArr + 0x01, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, gen)
		So(dec, ShouldResemble, gen)
	})

	Convey("map[string]bool will encode and decode", t, func() {
		var tmp map[string]bool
		var val = map[string]bool{"test": true}
		var gen = map[interface{}]interface{}{"test": true}
		var bit = []byte{cFixMap + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cTrue}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("map[string]int will encode and decode", t, func() {
		var tmp map[string]int
		var val = map[string]int{"test": math.MaxInt32}
		var gen = map[interface{}]interface{}{"test": math.MaxInt32}
		var bit = []byte{cFixMap + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cInt32, 127, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("map[string]uint will encode and decode", t, func() {
		var tmp map[string]uint
		var val = map[string]uint{"test": math.MaxUint32}
		var gen = map[interface{}]interface{}{"test": uint(math.MaxUint32)}
		var bit = []byte{cFixMap + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cUint32, 255, 255, 255, 255}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("map[string]string will encode and decode", t, func() {
		var tmp map[string]string
		var val = map[string]string{"test": "Hello"}
		var gen = map[interface{}]interface{}{"test": "Hello"}
		var bit = []byte{cFixMap + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cFixStr + 0x05, 72, 101, 108, 108, 111}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("map[int]interface{} will encode and decode", t, func() {
		var tmp map[int]interface{}
		var val = map[int]interface{}{math.MaxInt8: math.Pi}
		var gen = map[interface{}]interface{}{int(math.MaxInt8): float64(math.Pi)}
		var bit = []byte{cFixMap + 0x01 /**/, 127, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("map[uint]interface{} will encode and decode", t, func() {
		var tmp map[uint]interface{}
		var val = map[uint]interface{}{math.MaxUint8: math.Pi}
		var gen = map[interface{}]interface{}{uint(math.MaxUint8): float64(math.Pi)}
		var bit = []byte{cFixMap + 0x01 /**/, cUint8, 255, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("map[string]interface{} will encode and decode", t, func() {
		var tmp map[string]interface{}
		var val = map[string]interface{}{"test": math.Pi}
		var gen = map[interface{}]interface{}{"test": math.Pi}
		var bit = []byte{cFixMap + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("map[time.Time]interface{} will encode and decode", t, func() {
		var tmp map[time.Time]interface{}
		var val = map[time.Time]interface{}{tme: math.Pi}
		var gen = map[interface{}]interface{}{tme: math.Pi}
		var bit = []byte{cFixMap + 0x01 /**/, cTime, 7, 166, 199, 91, 123, 67, 205, 21, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("map[interface{}]interface{} will encode and decode", t, func() {
		var tmp map[interface{}]interface{}
		var val = map[interface{}]interface{}{"test": math.Pi}
		var gen = map[interface{}]interface{}{"test": math.Pi}
		var bit = []byte{cFixMap + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, val)
		So(dec, ShouldResemble, gen)
	})

	Convey("embedded map[string]interface{} will encode and decode", t, func() {
		var tmp map[string]interface{}
		var val = map[string]interface{}{"test": map[string]interface{}{"test": "Embedded"}}
		var oth = map[string]interface{}{"test": map[interface{}]interface{}{"test": "Embedded"}}
		var gen = map[interface{}]interface{}{"test": map[interface{}]interface{}{"test": "Embedded"}}
		var bit = []byte{cFixMap + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixMap + 0x01, cFixStr + 0x04, 116, 101, 115, 116, cFixStr + 0x08, 69, 109, 98, 101, 100, 100, 101, 100}
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, oth)
		So(dec, ShouldResemble, gen)
	})

	Convey("CustomTextMarshaler will encode and decode", t, func() {
		var tmp CustomTextMarshaler
		var val = &CustomTextMarshaler{Field: "TEXT"}
		var bit = []byte{cFixBin + 0x04, 84, 69, 88, 84}
		var oth = []byte("TEXT")
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(&tmp, ShouldResemble, val)
		So(dec, ShouldResemble, oth)
	})

	Convey("CustomBinaryMarshaler will encode and decode", t, func() {
		var tmp CustomBinaryMarshaler
		var val = &CustomBinaryMarshaler{Field: "DATA"}
		var bit = []byte{cFixBin + 0x04, 68, 65, 84, 65}
		var oth = []byte("DATA")
		var enc = Encode(val)
		var dec = Decode(bit)
		DecodeInto(bit, &tmp)
		So(enc, ShouldResemble, bit)
		So(&tmp, ShouldResemble, val)
		So(dec, ShouldResemble, oth)
	})

	Convey("CustomTextFailer will not encode and decode", t, func() {
		var tmp CustomTextFailer
		var val = &CustomTextFailer{Field: "TEXT"}
		var bit = []byte{cFixBin + 0x04, 84, 69, 88, 84}
		eer := NewEncoder(bytes.NewBuffer(nil)).Encode(val)
		der := NewDecoder(bytes.NewReader(bit)).Decode(&tmp)
		So(eer, ShouldNotBeNil)
		So(der, ShouldNotBeNil)
	})

	Convey("CustomBinaryFailer will not encode and decode", t, func() {
		var tmp CustomBinaryFailer
		var val = &CustomBinaryFailer{Field: "DATA"}
		var bit = []byte{cFixBin + 0x04, 68, 65, 84, 65}
		eer := NewEncoder(bytes.NewBuffer(nil)).Encode(val)
		der := NewDecoder(bytes.NewReader(bit)).Decode(&tmp)
		So(eer, ShouldNotBeNil)
		So(der, ShouldNotBeNil)
	})

	Convey("func will not encode and decode", t, func() {
		var tmp interface{}
		var val = func() {}
		var bit = []byte{cNil}
		var oth = interface{}(nil)
		var enc = Encode(val)
		var dec = Decode(bit)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, oth)
		So(dec, ShouldResemble, oth)
	})

	Convey("channel will not encode and decode", t, func() {
		var tmp interface{}
		var val = make(chan int)
		var bit = []byte{cNil}
		var oth = interface{}(nil)
		var enc = Encode(val)
		var dec = Decode(bit)
		So(enc, ShouldResemble, bit)
		So(tmp, ShouldResemble, oth)
		So(dec, ShouldResemble, oth)
	})

}

func isNotNum(v reflect.Kind) bool {
	switch v {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return false
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return false
	}
	return true
}

func isNotMap(v reflect.Kind) bool {
	return v != reflect.Map && v != reflect.Struct
}

func isNotFloat(v reflect.Kind) bool {
	return v != reflect.Float32 && v != reflect.Float64
}

func isNotComplex(v reflect.Kind) bool {
	return v != reflect.Complex64 && v != reflect.Complex128
}

func tester(out, val interface{}) {
	buf := bytes.NewBuffer(nil)
	enc := NewEncoder(buf)
	eer := enc.Encode(val)
	dec := NewDecoder(buf)
	der := dec.Decode(out)
	So(eer, ShouldBeNil)
	So(der, ShouldNotBeNil)
}

func TestInvariants(t *testing.T) {

	tme, _ := time.Parse(time.RFC3339, "1987-06-22T08:00:00.123456789Z")

	obj := []interface{}{
		true,
		false,
		str,
		bin,
		lng,
		tme,
		CustomInt(1),
		CustomString("Hi"),
		float32(math.Pi),
		float64(math.Pi),
		complex64(math.Pi),
		complex128(math.Pi),
		[]bool{true, false},
		[]string{"one", "two"},
		[]int{1, 2, 3, math.MaxInt8},
		[]int8{1, 2, 3, math.MaxInt8},
		[]int16{1, 2, 3, math.MaxInt16},
		[]int32{1, 2, 3, math.MaxInt32},
		[]int64{1, 2, 3, math.MaxInt64},
		[]uint{1, 2, 3, math.MaxUint16},
		[]uint16{1, 2, 3, math.MaxUint16},
		[]uint32{1, 2, 3, math.MaxUint32},
		[]uint64{1, 2, 3, math.MaxUint64},
		[]float32{1, 2, 3, math.Pi},
		[]float64{1, 2, 3, math.Pi},
		[]complex64{1, 2, 3, math.MaxUint64},
		[]complex128{1, 2, 3, math.MaxUint64},
		[]time.Time{tme, tme, tme, tme, tme, tme},
		[]interface{}{int8(1), "2", int16(3), int32(4), uint64(5)},
		Tested{"test", []byte("test"), []string{"1", "2"}, map[string]string{"1": "2"}, false, 25, "", ""},
		Corked{"test", []byte("test"), []string{"1", "2"}, map[string]string{"1": "2"}, false, 25, "", ""},
		Selfed{"test", []byte("test"), []string{"1", "2"}, map[string]string{"1": "2"}, false, 25, "", ""},
		map[string]int{},
		map[string]uint{},
		map[string]bool{},
		map[string]string{},
		map[int]interface{}{},
		map[uint]interface{}{},
		map[bool]interface{}{},
		map[string]interface{}{},
		map[time.Time]interface{}{},
		map[string]interface{}{
			"1": "test",
			"2": map[interface{}]interface{}{
				true: []byte("Check"),
			},
			"3": map[interface{}]interface{}{
				"str": str,
				"bin": bin,
				"lng": lng,
			},
		},
	}

	for _, v := range obj {

		Convey(fmt.Sprintf("Attempt to incorrectly decode %T into incorrect types...", v), t, func() {

			if reflect.TypeOf(v).Kind() != reflect.Bool {
				var out bool
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.String {
				var out string
				tester(&out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out int
				tester(&out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out int8
				tester(&out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out int16
				tester(&out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out int32
				tester(&out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out int64
				tester(&out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out uint
				tester(&out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out uint16
				tester(&out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out uint32
				tester(&out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out uint64
				tester(&out, v)
			}

			if isNotFloat(reflect.TypeOf(v).Kind()) {
				var out float32
				tester(&out, v)
			}

			if isNotFloat(reflect.TypeOf(v).Kind()) {
				var out float64
				tester(&out, v)
			}

			if isNotComplex(reflect.TypeOf(v).Kind()) {
				var out complex64
				tester(&out, v)
			}

			if isNotComplex(reflect.TypeOf(v).Kind()) {
				var out complex128
				tester(&out, v)
			}

			if reflect.TypeOf(v) != reflect.TypeOf(time.Time{}) {
				var out time.Time
				tester(&out, v)
			}

			// --------------------------------------------------

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []bool
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []string
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []int
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []int8
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []int16
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []int32
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []int64
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []uint
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []uint8
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []uint16
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []uint32
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []uint64
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []float32
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []float64
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []complex64
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []complex128
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []time.Time
				tester(&out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []interface{}
				tester(&out, v)
			}

			// --------------------------------------------------

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[string]int
				tester(&out, v)
			}

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[string]uint
				tester(&out, v)
			}

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[string]bool
				tester(&out, v)
			}

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[string]string
				tester(&out, v)
			}

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[string]interface{}
				tester(&out, v)
			}

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[interface{}]interface{}
				tester(&out, v)
			}

		})

	}

}

func sxpand(arg string, num int) (out string) {
	for i := 0; i < num; i++ {
		out = out + arg
	}
	return
}
func bxpand(arg []byte, num int) (out []byte) {
	for i := 0; i < num; i++ {
		out = append(out, arg...)
	}
	return
}

func extend(args ...interface{}) (out []byte) {
	for _, a := range args {
		switch v := a.(type) {
		case int:
			out = append(out, byte(v))
		case byte:
			out = append(out, v)
		case []byte:
			out = append(out, v...)
		}
	}
	return
}

func TestLargeEncoderDecoder(t *testing.T) {

	var obj = []interface{}{
		sxpand(str, 5),
		sxpand(str, 10),
		sxpand(str, 50),
		sxpand(str, 100),
		sxpand(str, 1000),
		sxpand(str, 5000),
		bxpand(bin, 5),
		bxpand(bin, 10),
		bxpand(bin, 50),
		bxpand(bin, 100),
		bxpand(bin, 1000),
		bxpand(bin, 5000),
	}

	var src = extend(
		cFixArr+0x0C,
		cStr16, 1, 34, bxpand(bin, 5),
		cStr16, 2, 68, bxpand(bin, 10),
		cStr16, 11, 84, bxpand(bin, 50),
		cStr16, 22, 168, bxpand(bin, 100),
		cStr16, 226, 144, bxpand(bin, 1000),
		cStr32, 0, 4, 108, 208, bxpand(bin, 5000),
		cBin16, 1, 34, bxpand(bin, 5),
		cBin16, 2, 68, bxpand(bin, 10),
		cBin16, 11, 84, bxpand(bin, 50),
		cBin16, 22, 168, bxpand(bin, 100),
		cBin16, 226, 144, bxpand(bin, 1000),
		cBin32, 0, 4, 108, 208, bxpand(bin, 5000),
	)

	Convey("Large encoder data encodes correctly", t, func() {
		var dst []byte
		enc := NewEncoderBytesFromPool(&dst)
		enc.Encode(obj)
		enc.Reset()
		So(dst, ShouldResemble, src)
	})

	Convey("Large decoder data decodes correctly", t, func() {
		var dst []interface{}
		dec := NewDecoderBytesFromPool(src)
		dec.Decode(&dst)
		dec.Reset()
		So(dst, ShouldResemble, obj)
	})

}
