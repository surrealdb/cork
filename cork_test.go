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
	"fmt"
	"math"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type TestInt int

type TestStr string

type Tested struct {
	Name  string
	Data  []byte `cork:"data"`
	Temp  []string
	priv  bool
	Count int
	Omit  string `cork:"-"`
	Empty string `cork:",omitempty"`
}

type Corked struct {
	Name  string
	Data  []byte `cork:"data"`
	Temp  []string
	priv  bool
	Count int
	Omit  string `cork:"-"`
	Empty string `cork:",omitempty"`
}

func (this *Corked) ExtendCORK() byte {
	return 0x01
}

func (this *Corked) MarshalCORK() (dst []byte, err error) {
	b := bytes.NewBuffer(dst)
	e := NewEncoder(b)
	e.Encode(this.Name)
	e.Encode(this.Data)
	e.Encode(this.Temp)
	e.Encode(this.Count)
	return b.Bytes(), nil
}

func (this *Corked) UnmarshalCORK(src []byte) (err error) {
	b := bytes.NewBuffer(src)
	d := NewDecoder(b)
	d.Decode(&this.Name)
	d.Decode(&this.Data)
	d.Decode(&this.Temp)
	d.Decode(&this.Count)
	return
}

func TestBasic(t *testing.T) {

	Convey("Static methods wil encode <=> decode", t, func() {
		obj := "test"
		enc := Encode(obj)
		dec := Decode(enc)
		So(dec, ShouldResemble, obj)
	})

}

func TestCorkers(t *testing.T) {

	Register(&Corked{})

	enc := Encode(&Corked{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""})

	tst := []byte{
		cExt8, cFixInt + 0x11, 0x01,
		cFixStr + 0x04, 116, 101, 115, 116,
		cFixBin + 0x04, 116, 101, 115, 116,
		cArr, cFixInt + 0x02, 129, 49, 129, 50,
		25,
	}

	Convey("Just check", t, func() {
		So(enc, ShouldResemble, tst)
	})

	Convey("Test decode into: interface{}", t, func() {
		var out interface{}
		chk := &Corked{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}
		buf := bytes.NewBuffer(enc)
		err := NewDecoder(buf).Decode(&out)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, chk)
	})

	Convey("Test decode into: Corked", t, func() {
		var out Corked
		chk := Corked{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}
		buf := bytes.NewBuffer(enc)
		err := NewDecoder(buf).Decode(&out)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, chk)
	})

	Convey("Test decode into: Corked{}", t, func() {
		out := Corked{}
		chk := Corked{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}
		buf := bytes.NewBuffer(enc)
		err := NewDecoder(buf).Decode(&out)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, chk)
	})

	// TODO fix this failing test
	/*Convey("Test decode into: *Corked (direct)", t, func() {
		var out *Corked
		chk := &Corked{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}
		buf := bytes.NewBuffer(enc)
		err := NewDecoder(buf).Decode(out)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, chk)
	})*/

	Convey("Test decode into: *Corked (pointer)", t, func() {
		var out *Corked
		chk := &Corked{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}
		buf := bytes.NewBuffer(enc)
		err := NewDecoder(buf).Decode(&out)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, chk)
	})

	Convey("Test decode into: &Corked{} (direct)", t, func() {
		out := &Corked{}
		chk := &Corked{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}
		buf := bytes.NewBuffer(enc)
		err := NewDecoder(buf).Decode(out)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, chk)
	})

	Convey("Test decode into: &Corked{} (pointer)", t, func() {
		out := &Corked{}
		chk := &Corked{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}
		buf := bytes.NewBuffer(enc)
		err := NewDecoder(buf).Decode(&out)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, chk)
	})

	Convey("Test decode into: new(Corked) (direct)", t, func() {
		out := new(Corked)
		chk := &Corked{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}
		buf := bytes.NewBuffer(enc)
		err := NewDecoder(buf).Decode(out)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, chk)
	})

	Convey("Test decode into: new(Corked) (pointer)", t, func() {
		out := new(Corked)
		chk := &Corked{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}
		buf := bytes.NewBuffer(enc)
		err := NewDecoder(buf).Decode(&out)
		So(err, ShouldBeNil)
		So(out, ShouldResemble, chk)
	})

}

func TestComplex(t *testing.T) {

	bef := Tested{"test", []byte("test"), []string{"1", "2"}, true, 25, "test", ""}
	err := Tested{"", []byte(nil), nil, false, 0, "", ""}
	aft := Tested{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}
	oth := map[string]interface{}{"Name": "test", "data": []byte("test"), "Temp": []interface{}{"1", "2"}, "Count": int8(25)}
	gen := map[interface{}]interface{}{"Name": "test", "data": []byte("test"), "Temp": []interface{}{"1", "2"}, "Count": int8(25)}

	Convey("Nil will decode into interface{}", t, func() {
		var out interface{}
		enc := Encode(nil)
		buf := bytes.NewBuffer(enc)
		NewDecoder(buf).Decode(&out)
		So(out, ShouldResemble, nil)
	})

	// ----------------------------------------------------------------------------------------------------

	Convey("Object will not decode into interface{} (direct)", t, func() {
		var out interface{}
		enc := Encode(bef)
		buf := bytes.NewBuffer(enc)
		NewDecoder(buf).Decode(out)
		So(out, ShouldResemble, nil)
	})

	Convey("Object (direct) will not decode into Tested (direct)", t, func() {
		var out Tested
		enc := Encode(bef)
		buf := bytes.NewBuffer(enc)
		NewDecoder(buf).Decode(out)
		So(out, ShouldResemble, err)
	})

	Convey("Object (pointer) will not decode into Tested (direct)", t, func() {
		var out Tested
		enc := Encode(&bef)
		buf := bytes.NewBuffer(enc)
		NewDecoder(buf).Decode(out)
		So(out, ShouldResemble, err)
	})

	// ----------------------------------------------------------------------------------------------------

	Convey("Object will decode into interface{} (pointer)", t, func() {
		var out interface{}
		enc := Encode(bef)
		buf := bytes.NewBuffer(enc)
		NewDecoder(buf).Decode(&out)
		So(out, ShouldResemble, gen)
	})

	Convey("Object (direct) will decode into Tested (pointer)", t, func() {
		var out Tested
		enc := Encode(bef)
		buf := bytes.NewBuffer(enc)
		NewDecoder(buf).Decode(&out)
		So(out, ShouldResemble, aft)
	})

	Convey("Object (pointer) will decode into Tested (pointer)", t, func() {
		var out Tested
		enc := Encode(&bef)
		buf := bytes.NewBuffer(enc)
		NewDecoder(buf).Decode(&out)
		So(out, ShouldResemble, aft)
	})

	// ----------------------------------------------------------------------------------------------------

	Convey("Object will decode into MapType", t, func() {
		var opt Handle
		var out interface{}
		opt.MapType = reflect.TypeOf(map[string]interface{}{})
		enc := Encode(bef)
		buf := bytes.NewBuffer(enc)
		NewDecoder(buf).Options(&opt).Decode(&out)
		So(out, ShouldResemble, oth)
	})

}

func TestComplete(t *testing.T) {

	clock, _ := time.Parse(time.RFC3339, "1987-06-22T08:00:00.123456789Z")

	str := "This is the very last time that I should have to test this"

	bin := []byte{
		84, 104, 105, 115, 32, 105, 115, 32, 116, 104, 101, 32, 118, 101,
		114, 121, 32, 108, 97, 115, 116, 32, 116, 105, 109, 101, 32, 116,
		104, 97, 116, 32, 73, 32, 115, 104, 111, 117, 108, 100, 32, 104, 97,
		118, 101, 32, 116, 111, 32, 116, 101, 115, 116, 32, 116, 104, 105, 115,
	}

	mstr := str + str + str + str + str
	mbin := append(bin, append(bin, append(bin, append(bin, bin...)...)...)...)

	var opt Handle
	opt.Precision = true
	opt.ArrType = reflect.TypeOf([]interface{}{})
	opt.MapType = reflect.TypeOf(map[interface{}]interface{}{})

	dec := []interface{}{
		true,
		false,
		"one",
		clock,
		TestInt(1),
		TestStr("Hi"),
		float32(math.Pi),
		float64(math.Pi),
		complex64(math.Pi),
		complex128(math.Pi),
		[]bool{true, false},
		[]string{"one", "two"},
		[]int8{1, 2, 3, math.MaxInt8},
		[]int16{1, 2, 3, math.MaxInt16},
		[]int32{1, 2, 3, math.MaxInt32},
		[]int64{1, 2, 3, math.MaxInt64},
		[]uint16{1, 2, 3, math.MaxUint16},
		[]uint32{1, 2, 3, math.MaxUint32},
		[]uint64{1, 2, 3, math.MaxUint64},
		[]float32{1, 2, 3, math.Pi},
		[]float64{1, 2, 3, math.Pi},
		[]complex64{1, 2, 3, math.MaxUint64},
		[]complex128{1, 2, 3, math.MaxUint64},
		[]time.Time{clock, clock, clock, clock},
		[]interface{}{int8(1), "2", int16(3), int32(4), uint64(5)},
		Tested{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""},
		// Corked{"test", []byte("test"), false, 25, "", ""},
		map[string]interface{}{
			"1": "test",
			"2": map[interface{}]interface{}{
				true: []byte("Check"),
			},
			"3": map[interface{}]interface{}{
				"s-str": str,
				"b-str": bin,
				"s-bin": mstr,
				"b-bin": mbin,
			},
		},
	}

	rev := []interface{}{
		true,
		false,
		"one",
		clock,
		int64(1),
		string("Hi"),
		float32(math.Pi),
		float64(math.Pi),
		complex64(math.Pi),
		complex128(math.Pi),
		[]interface{}{true, false},
		[]interface{}{"one", "two"},
		[]interface{}{int8(1), int8(2), int8(3), int8(math.MaxInt8)},
		[]interface{}{int16(1), int16(2), int16(3), int16(math.MaxInt16)},
		[]interface{}{int32(1), int32(2), int32(3), int32(math.MaxInt32)},
		[]interface{}{int64(1), int64(2), int64(3), int64(math.MaxInt64)},
		[]interface{}{uint16(1), uint16(2), uint16(3), uint16(math.MaxUint16)},
		[]interface{}{uint32(1), uint32(2), uint32(3), uint32(math.MaxUint32)},
		[]interface{}{uint64(1), uint64(2), uint64(3), uint64(math.MaxUint64)},
		[]interface{}{float32(1), float32(2), float32(3), float32(math.Pi)},
		[]interface{}{float64(1), float64(2), float64(3), float64(math.Pi)},
		[]interface{}{complex64(1), complex64(2), complex64(3), complex64(math.MaxUint64)},
		[]interface{}{complex128(1), complex128(2), complex128(3), complex128(math.MaxUint64)},
		[]interface{}{clock, clock, clock, clock},
		[]interface{}{int8(1), "2", int16(3), int32(4), uint64(5)},
		map[interface{}]interface{}{"Name": "test", "data": []byte("test"), "Temp": []interface{}{"1", "2"}, "Count": int8(25)},
		// &Corked{"test", []byte("test"), false, 25, "", ""},
		map[interface{}]interface{}{
			"1": "test",
			"2": map[interface{}]interface{}{
				true: []byte("Check"),
			},
			"3": map[interface{}]interface{}{
				"s-str": str,
				"b-str": bin,
				"s-bin": mstr,
				"b-bin": mbin,
			},
		},
	}

	Convey("Object should encode <=> decode into interface{}", t, func() {
		var out interface{}
		buf := bytes.NewBuffer(nil)
		eer := NewEncoder(buf).Options(&opt).Encode(dec)
		der := NewDecoder(buf).Options(&opt).Decode(&out)
		So(eer, ShouldBeNil)
		So(der, ShouldBeNil)
		So(out, ShouldResemble, rev)
	})

	Convey("Object should encode <=> decode into []interface{}", t, func() {
		var out []interface{}
		buf := bytes.NewBuffer(nil)
		eer := NewEncoder(buf).Options(&opt).Encode(dec)
		der := NewDecoder(buf).Options(&opt).Decode(&out)
		So(eer, ShouldBeNil)
		So(der, ShouldBeNil)
		So(out, ShouldResemble, rev)
	})

	Convey("Object should encode <=> decode into []interface{}{}", t, func() {
		out := []interface{}{}
		buf := bytes.NewBuffer(nil)
		eer := NewEncoder(buf).Options(&opt).Encode(dec)
		der := NewDecoder(buf).Options(&opt).Decode(&out)
		So(eer, ShouldBeNil)
		So(der, ShouldBeNil)
		So(out, ShouldResemble, rev)
	})

	for k, v := range dec {

		Convey(fmt.Sprintf("Object element: %T - should encode <=> decode into %T", v, v), t, func() {
			buf := bytes.NewBuffer(nil)
			out := reflect.New(reflect.TypeOf(v))
			eer := NewEncoder(buf).Options(&opt).Encode(v)
			der := NewDecoder(buf).Options(&opt).Decode(out.Interface())
			So(eer, ShouldBeNil)
			So(der, ShouldBeNil)
			So(out.Elem().Interface(), ShouldResemble, dec[k])
		})

		Convey(fmt.Sprintf("Object element: %T - should encode <=> decode into interface{}", v), t, func() {
			var out interface{}
			buf := bytes.NewBuffer(nil)
			eer := NewEncoder(buf).Options(&opt).Encode(v)
			der := NewDecoder(buf).Options(&opt).Decode(&out)
			So(eer, ShouldBeNil)
			So(der, ShouldBeNil)
			So(out, ShouldResemble, rev[k])
		})

		Convey(fmt.Sprintf("Object element: %T - attempt decoding into inverse objects...", v), t, func() {

			if reflect.TypeOf(v).Kind() != reflect.Bool {
				var out bool
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.String {
				var out string
				tester(&opt, &out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out int
				tester(&opt, &out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out int8
				tester(&opt, &out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out int16
				tester(&opt, &out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out int32
				tester(&opt, &out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out int64
				tester(&opt, &out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out uint
				tester(&opt, &out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out uint8
				tester(&opt, &out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out uint16
				tester(&opt, &out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out uint32
				tester(&opt, &out, v)
			}

			if isNotNum(reflect.TypeOf(v).Kind()) {
				var out uint64
				tester(&opt, &out, v)
			}

			if isNotFloat(reflect.TypeOf(v).Kind()) {
				var out float32
				tester(&opt, &out, v)
			}

			if isNotFloat(reflect.TypeOf(v).Kind()) {
				var out float64
				tester(&opt, &out, v)
			}

			if isNotComplex(reflect.TypeOf(v).Kind()) {
				var out complex64
				tester(&opt, &out, v)
			}

			if isNotComplex(reflect.TypeOf(v).Kind()) {
				var out complex128
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v) != reflect.TypeOf(time.Time{}) {
				var out time.Time
				tester(&opt, &out, v)
			}

			// --------------------------------------------------

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []bool
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []string
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []int
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []int8
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []int16
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []int32
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []int64
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []uint
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []uint8
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []uint16
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []uint32
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []uint64
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []float32
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []float64
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []complex64
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []complex128
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []time.Time
				tester(&opt, &out, v)
			}

			if reflect.TypeOf(v).Kind() != reflect.Slice {
				var out []interface{}
				tester(&opt, &out, v)
			}

			// --------------------------------------------------

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[string]int
				tester(&opt, &out, v)
			}

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[string]uint
				tester(&opt, &out, v)
			}

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[string]bool
				tester(&opt, &out, v)
			}

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[string]string
				tester(&opt, &out, v)
			}

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[string]interface{}
				tester(&opt, &out, v)
			}

			if isNotMap(reflect.TypeOf(v).Kind()) {
				var out map[interface{}]interface{}
				tester(&opt, &out, v)
			}

		})

	}

}

func TestPrecision(t *testing.T) {

	tests := []struct {
		dec interface{}
		rlx []byte
		pre []byte
	}{
		{
			dec: int(1),
			rlx: []byte{1},
			pre: []byte{1},
		},
		{
			dec: int8(1),
			rlx: []byte{1},
			pre: []byte{cInt8, 1},
		},
		{
			dec: int8(math.MaxInt8),
			rlx: []byte{127},
			pre: []byte{cInt8, 127},
		},
		{
			dec: int16(1),
			rlx: []byte{1},
			pre: []byte{cInt16, 0, 1},
		},
		{
			dec: int16(math.MaxInt8),
			rlx: []byte{127},
			pre: []byte{cInt16, 0, 127},
		},
		{
			dec: int16(math.MaxInt16),
			rlx: []byte{cInt16, 127, 255},
			pre: []byte{cInt16, 127, 255},
		},
		{
			dec: int32(1),
			rlx: []byte{1},
			pre: []byte{cInt32, 0, 0, 0, 1},
		},
		{
			dec: int32(math.MaxInt8),
			rlx: []byte{127},
			pre: []byte{cInt32, 0, 0, 0, 127},
		},
		{
			dec: int32(math.MaxInt16),
			rlx: []byte{cInt16, 127, 255},
			pre: []byte{cInt32, 0, 0, 127, 255},
		},
		{
			dec: int32(math.MaxInt32),
			rlx: []byte{cInt32, 127, 255, 255, 255},
			pre: []byte{cInt32, 127, 255, 255, 255},
		},
		{
			dec: int64(1),
			rlx: []byte{1},
			pre: []byte{cInt64, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			dec: int64(math.MaxInt8),
			rlx: []byte{127},
			pre: []byte{cInt64, 0, 0, 0, 0, 0, 0, 0, 127},
		},
		{
			dec: int64(math.MaxInt16),
			rlx: []byte{cInt16, 127, 255},
			pre: []byte{cInt64, 0, 0, 0, 0, 0, 0, 127, 255},
		},
		{
			dec: int64(math.MaxInt32),
			rlx: []byte{cInt32, 127, 255, 255, 255},
			pre: []byte{cInt64, 0, 0, 0, 0, 127, 255, 255, 255},
		},
		{
			dec: int64(math.MaxInt64),
			rlx: []byte{cInt64, 127, 255, 255, 255, 255, 255, 255, 255},
			pre: []byte{cInt64, 127, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: uint(1),
			rlx: []byte{1},
			pre: []byte{1},
		},
		{
			dec: uint8(1),
			rlx: []byte{1},
			pre: []byte{cUint8, 1},
		},
		{
			dec: uint8(math.MaxUint8),
			rlx: []byte{cUint8, 255},
			pre: []byte{cUint8, 255},
		},
		{
			dec: uint16(1),
			rlx: []byte{1},
			pre: []byte{cUint16, 0, 1},
		},
		{
			dec: uint16(math.MaxUint8),
			rlx: []byte{cUint8, 255},
			pre: []byte{cUint16, 0, 255},
		},
		{
			dec: uint16(math.MaxUint16),
			rlx: []byte{cUint16, 255, 255},
			pre: []byte{cUint16, 255, 255},
		},
		{
			dec: uint32(1),
			rlx: []byte{1},
			pre: []byte{cUint32, 0, 0, 0, 1},
		},
		{
			dec: uint32(math.MaxUint8),
			rlx: []byte{cUint8, 255},
			pre: []byte{cUint32, 0, 0, 0, 255},
		},
		{
			dec: uint32(math.MaxUint16),
			rlx: []byte{cUint16, 255, 255},
			pre: []byte{cUint32, 0, 0, 255, 255},
		},
		{
			dec: uint32(math.MaxUint32),
			rlx: []byte{cUint32, 255, 255, 255, 255},
			pre: []byte{cUint32, 255, 255, 255, 255},
		},
		{
			dec: uint64(1),
			rlx: []byte{1},
			pre: []byte{cUint64, 0, 0, 0, 0, 0, 0, 0, 1},
		},
		{
			dec: uint64(math.MaxUint8),
			rlx: []byte{cUint8, 255},
			pre: []byte{cUint64, 0, 0, 0, 0, 0, 0, 0, 255},
		},
		{
			dec: uint64(math.MaxUint16),
			rlx: []byte{cUint16, 255, 255},
			pre: []byte{cUint64, 0, 0, 0, 0, 0, 0, 255, 255},
		},
		{
			dec: uint64(math.MaxUint32),
			rlx: []byte{cUint32, 255, 255, 255, 255},
			pre: []byte{cUint64, 0, 0, 0, 0, 255, 255, 255, 255},
		},
		{
			dec: uint64(math.MaxUint64),
			rlx: []byte{cUint64, 255, 255, 255, 255, 255, 255, 255, 255},
			pre: []byte{cUint64, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: float32(math.Pi),
			rlx: []byte{cFloat32, 64, 73, 15, 219},
			pre: []byte{cFloat32, 64, 73, 15, 219},
		},
		{
			dec: float64(math.Pi),
			rlx: []byte{cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
			pre: []byte{cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			dec: complex64(math.Pi),
			rlx: []byte{cComplex64, 64, 73, 15, 219, 0, 0, 0, 0},
			pre: []byte{cComplex64, 64, 73, 15, 219, 0, 0, 0, 0},
		},
		{
			dec: complex128(math.Pi),
			rlx: []byte{cComplex128, 64, 9, 33, 251, 84, 68, 45, 24, 0, 0, 0, 0, 0, 0, 0, 0},
			pre: []byte{cComplex128, 64, 9, 33, 251, 84, 68, 45, 24, 0, 0, 0, 0, 0, 0, 0, 0},
		},
	}

	// ----------------------------------------------------------------------------------------------------

	for _, test := range tests {

		var opt Handle

		Convey(fmt.Sprintf("Object should encode and decode (with Precision)  --- %T : %v", test.dec, test.dec), t, func() {
			opt.Precision = true
			buf := bytes.NewBuffer(nil)
			out := reflect.New(reflect.TypeOf(test.dec))
			NewEncoder(buf).Options(&opt).Encode(test.dec)
			So(buf.Bytes(), ShouldResemble, test.pre)
			NewDecoder(buf).Options(&opt).Decode(out.Interface())
			So(out.Elem().Interface(), ShouldResemble, test.dec)
		})

		Convey(fmt.Sprintf("Object should encode and decode (without Precision) --- %T : %v", test.dec, test.dec), t, func() {
			opt.Precision = false
			buf := bytes.NewBuffer(nil)
			out := reflect.New(reflect.TypeOf(test.dec))
			NewEncoder(buf).Options(&opt).Encode(test.dec)
			So(buf.Bytes(), ShouldResemble, test.rlx)
			NewDecoder(buf).Options(&opt).Decode(out.Interface())
			So(out.Elem().Interface(), ShouldResemble, test.dec)
		})

		Convey(fmt.Sprintf("Object should encode and decode (with Precision) into interface  --- %T : %v", test.dec, test.dec), t, func() {
			var out interface{}
			opt.Precision = true
			buf := bytes.NewBuffer(nil)
			NewEncoder(buf).Options(&opt).Encode(test.dec)
			So(buf.Bytes(), ShouldResemble, test.pre)
			NewDecoder(buf).Options(&opt).Decode(&out)
			So(out, ShouldEqual, test.dec)
		})

		Convey(fmt.Sprintf("Object should encode and decode (without Precision) into interface  --- %T : %v", test.dec, test.dec), t, func() {
			var out interface{}
			opt.Precision = false
			buf := bytes.NewBuffer(nil)
			NewEncoder(buf).Options(&opt).Encode(test.dec)
			So(buf.Bytes(), ShouldResemble, test.rlx)
			NewDecoder(buf).Options(&opt).Decode(&out)
			So(out, ShouldEqual, test.dec)
		})

	}

}

func tester(opt *Handle, out, val interface{}) {
	Convey(fmt.Sprintf("Object element: %T - will fail to encode <=> decode into %T", val, out), func() {
		buf := bytes.NewBuffer(nil)
		eer := NewEncoder(buf).Options(opt).Encode(val)
		der := NewDecoder(buf).Options(opt).Decode(out)
		So(eer, ShouldBeNil)
		So(der, ShouldNotBeNil)
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
