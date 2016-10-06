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

type Tester struct {
	Name  string
	Data  []byte `cork:"data"`
	priv  bool
	Count int
	Omit  string `cork:"-"`
	Empty string `cork:",omitempty"`
}

type Mapper struct {
	Name  string
	Data  []byte `cork:"data"`
	priv  bool
	Count int
	Omit  string `cork:"-"`
	Empty string `cork:",omitempty"`
}

var tests []struct {
	out interface{}
	com interface{}
	alt interface{}
	dec interface{}
	enc []byte
}

func TestMain(t *testing.T) {

	clock, _ := time.Parse(time.RFC3339, "1987-06-22T08:00:00.123456789Z")

	tests = []struct {
		out interface{}
		com interface{}
		alt interface{}
		dec interface{}
		enc []byte
	}{
		{
			dec: nil,
			enc: []byte{cNil},
		},
		{
			dec: true,
			enc: []byte{cTrue},
		},
		{
			dec: false,
			enc: []byte{cFalse},
		},
		{
			dec: "Hello",
			enc: []byte{cFixStr + 0x05, 72, 101, 108, 108, 111},
		},
		{
			dec: "This is the very last time that I should have to test this",
			enc: []byte{cStr8, 58,
				84, 104, 105, 115, 32, 105, 115, 32, 116, 104, 101, 32, 118, 101,
				114, 121, 32, 108, 97, 115, 116, 32, 116, 105, 109, 101, 32, 116,
				104, 97, 116, 32, 73, 32, 115, 104, 111, 117, 108, 100, 32, 104, 97,
				118, 101, 32, 116, 111, 32, 116, 101, 115, 116, 32, 116, 104, 105, 115,
			},
		},
		{
			dec: []byte{72, 101, 108, 108, 111},
			enc: []byte{cFixBin + 0x05, 72, 101, 108, 108, 111},
		},
		{
			dec: clock,
			enc: []byte{cTime, 7, 166, 199, 91, 123, 67, 205, 21},
		},
		{
			dec: int(1),
			enc: []byte{1},
		},
		{
			dec: int8(1),
			enc: []byte{1},
		},
		{
			dec: int8(math.MaxInt8),
			enc: []byte{127},
		},
		{
			dec: int16(1),
			enc: []byte{1},
		},
		{
			dec: int16(math.MaxInt8),
			enc: []byte{127},
		},
		{
			dec: int16(math.MaxInt16),
			enc: []byte{cInt16, 127, 255},
		},
		{
			dec: int32(1),
			enc: []byte{1},
		},
		{
			dec: int32(math.MaxInt8),
			enc: []byte{127},
		},
		{
			dec: int32(math.MaxInt16),
			enc: []byte{cInt16, 127, 255},
		},
		{
			dec: int32(math.MaxInt32),
			enc: []byte{cInt32, 127, 255, 255, 255},
		},
		{
			dec: int64(1),
			enc: []byte{1},
		},
		{
			dec: int64(math.MaxInt8),
			enc: []byte{127},
		},
		{
			dec: int64(math.MaxInt16),
			enc: []byte{cInt16, 127, 255},
		},
		{
			dec: int64(math.MaxInt32),
			enc: []byte{cInt32, 127, 255, 255, 255},
		},
		{
			dec: int64(math.MaxInt64),
			enc: []byte{cInt64, 127, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: uint(1),
			enc: []byte{1},
		},
		{
			dec: uint8(1),
			enc: []byte{1},
		},
		{
			dec: uint8(math.MaxUint8),
			enc: []byte{cUint8, 255},
		},
		{
			dec: uint16(1),
			enc: []byte{1},
		},
		{
			dec: uint16(math.MaxUint8),
			enc: []byte{cUint8, 255},
		},
		{
			dec: uint16(math.MaxUint16),
			enc: []byte{cUint16, 255, 255},
		},
		{
			dec: uint32(1),
			enc: []byte{1},
		},
		{
			dec: uint32(math.MaxUint8),
			enc: []byte{cUint8, 255},
		},
		{
			dec: uint32(math.MaxUint16),
			enc: []byte{cUint16, 255, 255},
		},
		{
			dec: uint32(math.MaxUint32),
			enc: []byte{cUint32, 255, 255, 255, 255},
		},
		{
			dec: uint64(1),
			enc: []byte{1},
		},
		{
			dec: uint64(math.MaxUint8),
			enc: []byte{cUint8, 255},
		},
		{
			dec: uint64(math.MaxUint16),
			enc: []byte{cUint16, 255, 255},
		},
		{
			dec: uint64(math.MaxUint32),
			enc: []byte{cUint32, 255, 255, 255, 255},
		},
		{
			dec: uint64(math.MaxUint64),
			enc: []byte{cUint64, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: float32(math.Pi),
			enc: []byte{cFloat32, 64, 73, 15, 219},
		},
		{
			dec: float64(math.Pi),
			enc: []byte{cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			dec: []bool{true, false},
			enc: []byte{cArrBool, cFixInt + 0x02, cTrue, cFalse},
		},
		{
			dec: []string{"Hello", "World"},
			enc: []byte{cArrStr, cFixInt + 0x02 /**/, cFixStr + 0x05, 72, 101, 108, 108, 111 /**/, cFixStr + 0x05, 87, 111, 114, 108, 100},
		},
		{
			dec: []time.Time{clock, clock},
			enc: []byte{cArrTime, cFixInt + 0x02 /**/, 7, 166, 199, 91, 123, 67, 205, 21 /**/, 7, 166, 199, 91, 123, 67, 205, 21},
		},
		{
			dec: []int8{1, math.MaxInt8},
			enc: []byte{cArrInt8, cFixInt + 0x02 /**/, 1 /**/, 127},
		},
		{
			dec: []int16{1, math.MaxInt8, math.MaxInt16},
			enc: []byte{cArrInt16, cFixInt + 0x03 /**/, 1 /**/, 127 /**/, cInt16, 127, 255},
		},
		{
			dec: []int32{1, math.MaxInt8, math.MaxInt16, math.MaxInt32},
			enc: []byte{cArrInt32, cFixInt + 0x04 /**/, 1 /**/, 127 /**/, cInt16, 127, 255 /**/, cInt32, 127, 255, 255, 255},
		},
		{
			dec: []int64{1, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64},
			enc: []byte{cArrInt64, cFixInt + 0x05 /**/, 1 /**/, 127 /**/, cInt16, 127, 255 /**/, cInt32, 127, 255, 255, 255 /**/, cInt64, 127, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: []int{0, 1, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64},
			enc: []byte{cArrInt, cFixInt + 0x06 /**/, 0 /**/, 1 /**/, 127 /**/, cInt16, 127, 255 /**/, cInt32, 127, 255, 255, 255 /**/, cInt64, 127, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: []uint16{1, math.MaxUint8, math.MaxUint16},
			enc: []byte{cArrUint16, cFixInt + 0x03 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255},
		},
		{
			dec: []uint32{1, math.MaxUint8, math.MaxUint16, math.MaxUint32},
			enc: []byte{cArrUint32, cFixInt + 0x04 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255 /**/, cUint32, 255, 255, 255, 255},
		},
		{
			dec: []uint64{1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64},
			enc: []byte{cArrUint64, cFixInt + 0x05 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255 /**/, cUint32, 255, 255, 255, 255 /**/, cUint64, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: []uint{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64},
			enc: []byte{cArrUint, cFixInt + 0x06 /**/, 0 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255 /**/, cUint32, 255, 255, 255, 255 /**/, cUint64, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: []float32{math.Pi, math.Pi},
			enc: []byte{cArrFloat32, cFixInt + 0x02 /**/, 64, 73, 15, 219 /**/, 64, 73, 15, 219},
		},
		{
			dec: []float64{math.Pi, math.Pi},
			enc: []byte{cArrFloat64, cFixInt + 0x02 /**/, 64, 9, 33, 251, 84, 68, 45, 24 /**/, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			dec: []interface{}{nil, true, false, "test", []byte("test"), int8(77), uint8(177), float64(math.Pi)},
			enc: []byte{cArrNil, cFixInt + 0x08 /**/, cNil, cTrue, cFalse /**/, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixBin + 0x04, 116, 101, 115, 116 /**/, 77 /**/, cUint8, 177 /**/, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			dec: map[string]bool{"one": true, "two": false},
			enc: []byte{cMapStrBool, cFixInt + 0x02 /**/, cFixStr + 0x03, 111, 110, 101, cTrue, cFixStr + 0x03, 116, 119, 111, cFalse},
		},
		{
			dec: map[string]int{"one": 1, "two": math.MaxInt64},
			enc: []byte{cMapStrInt, cFixInt + 0x02 /**/, cFixStr + 0x03, 111, 110, 101, cFixInt + 0x01, cFixStr + 0x03, 116, 119, 111, cInt64, 127, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: map[string]string{"one": "Hello", "two": "World"},
			enc: []byte{cMapStrStr, cFixInt + 0x02 /**/, cFixStr + 0x03, 111, 110, 101, cFixStr + 0x05, 72, 101, 108, 108, 111 /**/, cFixStr + 0x03, 116, 119, 111, cFixStr + 0x05, 87, 111, 114, 108, 100},
		},
		{
			dec: map[string]interface{}{"one": "Hello", "two": math.Pi},
			enc: []byte{cMapStrNil, cFixInt + 0x02 /**/, cFixStr + 0x03, 111, 110, 101, cFixStr + 0x05, 72, 101, 108, 108, 111 /**/, cFixStr + 0x03, 116, 119, 111, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			dec: map[interface{}]interface{}{"one": "Hello", int8(2): float64(math.Pi)},
			enc: []byte{cMapNilNil, cFixInt + 0x02 /**/, cFixStr + 0x03, 111, 110, 101, cFixStr + 0x05, 72, 101, 108, 108, 111 /**/, cFixInt + 0x02, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			out: &Mapper{},
			com: &Mapper{"test", []byte("test"), false, 25, "", ""},
			alt: map[string]interface{}{"Name": "test", "data": []byte("test"), "Count": int8(25)},
			dec: Mapper{"test", []byte("test"), true, 25, "test", ""},
			enc: []byte{cMapStrNil, cFixInt + 0x03 /**/, cFixStr + 0x04, 78, 97, 109, 101, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixStr + 0x04, 100, 97, 116, 97, cFixBin + 0x04, 116, 101, 115, 116 /**/, cFixStr + 0x05, 67, 111, 117, 110, 116 /**/, 25},
		},
		{
			out: &Mapper{},
			com: &Mapper{"test", []byte("test"), false, 25, "", ""},
			alt: map[string]interface{}{"Name": "test", "data": []byte("test"), "Count": int8(25)},
			dec: &Mapper{"test", []byte("test"), true, 25, "test", ""},
			enc: []byte{cMapStrNil, cFixInt + 0x03 /**/, cFixStr + 0x04, 78, 97, 109, 101, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixStr + 0x04, 100, 97, 116, 97, cFixBin + 0x04, 116, 101, 115, 116 /**/, cFixStr + 0x05, 67, 111, 117, 110, 116 /**/, 25},
		},
		// {
		// 	alt: Tester{"test", []byte("test"), false, 25, "test", ""},
		// 	dec: Tester{"test", []byte("test"), true, 25, "test", ""},
		// 	enc: []byte{cMapStrNil, cFixInt + 0x03 /**/, cFixStr + 0x04, 78, 97, 109, 101, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixStr + 0x04, 68, 97, 116, 97, cFixBin + 0x04, 116, 101, 115, 116 /**/, cFixStr + 0x05, 67, 111, 117, 110, 116 /**/, 25},
		// },
		// {
		// 	alt: &Tester{"test", []byte("test"), false, 25, "test", ""},
		// 	dec: &Tester{"test", []byte("test"), true, 25, "test", ""},
		// 	enc: []byte{cMapStrNil, cFixInt + 0x03 /**/, cFixStr + 0x04, 78, 97, 109, 101, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixStr + 0x04, 68, 97, 116, 97, cFixBin + 0x04, 116, 101, 115, 116 /**/, cFixStr + 0x05, 67, 111, 117, 110, 116 /**/, 25},
		// },
	}

}

// ----------------------------------------------------------------------------------------------------

func TestDefined(t *testing.T) {

	for _, test := range tests {

		if test.dec != nil && test.out != nil {
			Convey(fmt.Sprintf("Object should encode and decode into type --- %T : %v", test.dec, test.dec), t, func() {
				enc := Encode(test.dec)
				buf := bytes.NewBuffer(enc)
				NewDecoder(buf).Decode(test.out)
				So(test.out, ShouldResemble, test.com)
			})
		}

		if test.dec != nil && test.out == nil {
			Convey(fmt.Sprintf("Object should encode and decode into type --- %T : %v", test.dec, test.dec), t, func() {
				dup := reflect.New(reflect.ValueOf(test.dec).Type())
				enc := Encode(test.dec)
				buf := bytes.NewBuffer(enc)
				NewDecoder(buf).Decode(dup.Interface())
				So(dup.Elem().Interface(), ShouldResemble, test.dec)
			})
		}

	}

}

func TestUndefined(t *testing.T) {

	for _, test := range tests {

		Convey(fmt.Sprintf("Object should encode and decode into interface --- %T : %v", test.dec, test.dec), t, func() {

			enc := Encode(test.dec)
			dec := Decode(enc)

			switch test.dec.(type) {
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				So(enc, ShouldResemble, test.enc)
				So(dec, ShouldEqual, test.dec)
			case map[string]int, map[string]bool, map[string]string, map[string]interface{}, map[interface{}]interface{}:
				So(dec, ShouldResemble, test.dec)
			case Tester, *Tester, Mapper, *Mapper:
				So(enc, ShouldResemble, test.enc)
				So(dec, ShouldResemble, test.alt)
			default:
				So(enc, ShouldResemble, test.enc)
				So(dec, ShouldResemble, test.dec)
			}

		})

	}

}
