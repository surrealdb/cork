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

func TestDecode(t *testing.T) {

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

	tests := []struct {
		dec interface{}
		obj interface{}
	}{
		{
			dec: true,
			obj: true,
		},
		{
			dec: false,
			obj: false,
		},
		{
			dec: "Hello",
			obj: "Hello",
		},
		{
			dec: str,
			obj: str,
		},
		{
			dec: mstr,
			obj: mstr,
		},
		{
			dec: []byte{72, 101, 108, 108, 111},
			obj: []byte{72, 101, 108, 108, 111},
		},
		{
			dec: bin,
			obj: bin,
		},
		{
			dec: mbin,
			obj: mbin,
		},
		{
			dec: clock,
			obj: clock,
		},
		{
			dec: Tested{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""},
			obj: map[interface{}]interface{}{"Name": "test", "data": []byte("test"), "Temp": []interface{}{"1", "2"}, "Count": int64(25)},
		},
		// ------------------------------------------------------------------------------------------------------------------------
		// ------------------------------------------------------------------------------------------------------------------------
		// ------------------------------------------------------------------------------------------------------------------------
		{
			dec: []bool{true, false},
			obj: []interface{}{true, false},
		},
		{
			dec: []string{"Hello", "World"},
			obj: []interface{}{"Hello", "World"},
		},
		{
			dec: []time.Time{clock, clock},
			obj: []interface{}{clock, clock},
		},
		{
			dec: []int8{1, math.MaxInt8},
			obj: []interface{}{int64(1), int64(math.MaxInt8)},
		},
		{
			dec: []int16{1, math.MaxInt8, math.MaxInt16},
			obj: []interface{}{int64(1), int64(math.MaxInt8), int64(math.MaxInt16)},
		},
		{
			dec: []int32{1, math.MaxInt8, math.MaxInt16, math.MaxInt32},
			obj: []interface{}{int64(1), int64(math.MaxInt8), int64(math.MaxInt16), int64(math.MaxInt32)},
		},
		{
			dec: []int64{1, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64},
			obj: []interface{}{int64(1), int64(math.MaxInt8), int64(math.MaxInt16), int64(math.MaxInt32), int64(math.MaxInt64)},
		},
		{
			dec: []int{0, 1, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64},
			obj: []interface{}{int64(0), int64(1), int64(math.MaxInt8), int64(math.MaxInt16), int64(math.MaxInt32), int64(math.MaxInt64)},
		},
		{
			dec: []uint16{1, math.MaxUint8, math.MaxUint16},
			obj: []interface{}{int64(1), uint64(math.MaxUint8), uint64(math.MaxUint16)},
		},
		{
			dec: []uint32{1, math.MaxUint8, math.MaxUint16, math.MaxUint32},
			obj: []interface{}{int64(1), uint64(math.MaxUint8), uint64(math.MaxUint16), uint64(math.MaxUint32)},
		},
		{
			dec: []uint64{1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64},
			obj: []interface{}{int64(1), uint64(math.MaxUint8), uint64(math.MaxUint16), uint64(math.MaxUint32), uint64(math.MaxUint64)},
		},
		{
			dec: []uint{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64},
			obj: []interface{}{int64(0), int64(1), uint64(math.MaxUint8), uint64(math.MaxUint16), uint64(math.MaxUint32), uint64(math.MaxUint64)},
		},
		{
			dec: []float32{math.Pi, math.Pi},
			obj: []interface{}{float32(math.Pi), float32(math.Pi)},
		},
		{
			dec: []float64{math.Pi, math.Pi},
			obj: []interface{}{float64(math.Pi), float64(math.Pi)},
		},
		{
			dec: []complex64{math.Pi, math.Pi},
			obj: []interface{}{complex64(math.Pi), complex64(math.Pi)},
		},
		{
			dec: []complex128{math.Pi, math.Pi},
			obj: []interface{}{complex128(math.Pi), complex128(math.Pi)},
		},
		{
			dec: []interface{}{nil, true, false, "test", []byte("test"), int64(77), uint64(177), float64(math.Pi)},
			obj: []interface{}{nil, true, false, "test", []byte("test"), int64(77), uint64(177), float64(math.Pi)},
		},
		{
			dec: []Tested{{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}},
			obj: []interface{}{map[interface{}]interface{}{"Name": "test", "data": []byte("test"), "Temp": []interface{}{"1", "2"}, "Count": int64(25)}},
		},
		// ------------------------------------------------------------------------------------------------------------------------
		// ------------------------------------------------------------------------------------------------------------------------
		// ------------------------------------------------------------------------------------------------------------------------
		{
			dec: [][]float64{{math.Pi}, {math.Pi}},
			obj: []interface{}{[]interface{}{float64(math.Pi)}, []interface{}{float64(math.Pi)}},
		},
		{
			dec: [][][]float64{{{math.Pi}}, {{math.Pi}}},
			obj: []interface{}{[]interface{}{[]interface{}{float64(math.Pi)}}, []interface{}{[]interface{}{float64(math.Pi)}}},
		},
		// ------------------------------------------------------------------------------------------------------------------------
		// ------------------------------------------------------------------------------------------------------------------------
		// ------------------------------------------------------------------------------------------------------------------------
		{
			dec: map[float32]bool{math.Pi: true},
			obj: map[interface{}]interface{}{float32(math.Pi): true},
		},
		{
			dec: map[float64]bool{math.Pi: true},
			obj: map[interface{}]interface{}{float64(math.Pi): true},
		},
		// ------------------------------------------------------------------------------------------------------------------------
		{
			dec: map[int]int{1: 1, 2: math.MaxInt64},
			obj: map[interface{}]interface{}{int64(1): int64(1), int64(2): int64(math.MaxInt64)},
		},
		{
			dec: map[int]uint{1: 1, 2: math.MaxUint64},
			obj: map[interface{}]interface{}{int64(1): int64(1), int64(2): uint64(math.MaxUint64)},
		},
		{
			dec: map[int]bool{1: true, 2: false},
			obj: map[interface{}]interface{}{int64(1): true, int64(2): false},
		},
		{
			dec: map[int]string{1: "Hello", 2: "World"},
			obj: map[interface{}]interface{}{int64(1): "Hello", int64(2): "World"},
		},
		{
			dec: map[int]interface{}{1: "Hello", 2: math.Pi},
			obj: map[interface{}]interface{}{int64(1): "Hello", int64(2): float64(math.Pi)},
		},
		// ------------------------------------------------------------------------------------------------------------------------
		{
			dec: map[uint]int{1: 1, 2: math.MaxInt64},
			obj: map[interface{}]interface{}{int64(1): int64(1), int64(2): int64(math.MaxInt64)},
		},
		{
			dec: map[uint]uint{1: 1, 2: math.MaxUint64},
			obj: map[interface{}]interface{}{int64(1): int64(1), int64(2): uint64(math.MaxUint64)},
		},
		{
			dec: map[uint]bool{1: true, 2: false},
			obj: map[interface{}]interface{}{int64(1): true, int64(2): false},
		},
		{
			dec: map[uint]string{1: "Hello", 2: "World"},
			obj: map[interface{}]interface{}{int64(1): "Hello", int64(2): "World"},
		},
		{
			dec: map[uint]interface{}{1: "Hello", 2: math.Pi},
			obj: map[interface{}]interface{}{int64(1): "Hello", int64(2): float64(math.Pi)},
		},
		// ------------------------------------------------------------------------------------------------------------------------
		{
			dec: map[string]int{"one": 1, "two": math.MaxInt64},
			obj: map[interface{}]interface{}{"one": int64(1), "two": int64(math.MaxInt64)},
		},
		{
			dec: map[string]uint{"one": 1, "two": math.MaxUint64},
			obj: map[interface{}]interface{}{"one": int64(1), "two": uint64(math.MaxUint64)},
		},
		{
			dec: map[string]bool{"one": true, "two": false},
			obj: map[interface{}]interface{}{"one": true, "two": false},
		},
		{
			dec: map[string]string{"one": "Hello", "two": "World"},
			obj: map[interface{}]interface{}{"one": "Hello", "two": "World"},
		},
		{
			dec: map[string]interface{}{"one": "Hello", "two": math.Pi},
			obj: map[interface{}]interface{}{"one": "Hello", "two": float64(math.Pi)},
		},
		// ------------------------------------------------------------------------------------------------------------------------
		{
			dec: map[interface{}]interface{}{"one": "Hello", int64(2): float64(math.Pi)},
			obj: map[interface{}]interface{}{"one": "Hello", int64(2): float64(math.Pi)},
		},
		{
			dec: map[interface{}]interface{}{"one": "Hello", "two": map[interface{}]interface{}{"three": "Test", "four": "Embedded"}},
			obj: map[interface{}]interface{}{"one": "Hello", "two": map[interface{}]interface{}{"three": "Test", "four": "Embedded"}},
		},
	}

	// ----------------------------------------------------------------------------------------------------

	for _, test := range tests {

		Convey(fmt.Sprintf("%T will decode into interface --- %v", test.dec, test.dec), t, func() {
			var out interface{}
			buf := bytes.NewBuffer(nil)
			eer := NewEncoder(buf).Encode(test.dec)
			der := NewDecoder(buf).Decode(&out)
			So(eer, ShouldBeNil)
			So(der, ShouldBeNil)
			So(out, ShouldResemble, test.obj)
		})

		Convey(fmt.Sprintf("%T will decode into new object --- %v", test.dec, test.dec), t, func() {
			buf := bytes.NewBuffer(nil)
			out := reflect.New(reflect.TypeOf(test.dec))
			eer := NewEncoder(buf).Encode(test.dec)
			der := NewDecoder(buf).Decode(out.Interface())
			So(eer, ShouldBeNil)
			So(der, ShouldBeNil)
			So(out.Elem().Interface(), ShouldResemble, test.dec)
		})

		Convey(fmt.Sprintf("%T will decode into same object --- %v", test.dec, test.dec), t, func() {
			buf := bytes.NewBuffer(nil)
			eer := NewEncoder(buf).Encode(test.dec)
			der := NewDecoder(buf).Decode(&test.dec)
			So(eer, ShouldBeNil)
			So(der, ShouldBeNil)
			So(test.dec, ShouldResemble, test.dec)
		})

	}

}
