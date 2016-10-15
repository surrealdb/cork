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
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEncode(t *testing.T) {

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
		enc []byte
	}{
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
			dec: str,
			enc: append([]byte{cStr8, 58}, bin...),
		},
		{
			dec: mstr,
			enc: append([]byte{cStr16, 1, 34}, mbin...),
		},
		{
			dec: []byte{72, 101, 108, 108, 111},
			enc: []byte{cFixBin + 0x05, 72, 101, 108, 108, 111},
		},
		{
			dec: bin,
			enc: append([]byte{cBin8, 58}, bin...),
		},
		{
			dec: mbin,
			enc: append([]byte{cBin16, 1, 34}, mbin...),
		},
		{
			dec: clock,
			enc: []byte{cTime, 7, 166, 199, 91, 123, 67, 205, 21},
		},
		{
			dec: Tested{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""},
			enc: []byte{cMap, cFixInt + 0x04, /**/
				cFixStr + 0x04, 78, 97, 109, 101, cFixStr + 0x04, 116, 101, 115, 116, /**/
				cFixStr + 0x04, 100, 97, 116, 97, cFixBin + 0x04, 116, 101, 115, 116, /**/
				cFixStr + 0x04, 84, 101, 109, 112, cArr, cFixInt + 0x02, 129, 49, 129, 50, /**/
				cFixStr + 0x05, 67, 111, 117, 110, 116 /**/, 25,
			},
		},
		{
			dec: &Tested{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""},
			enc: []byte{cMap, cFixInt + 0x04, /**/
				cFixStr + 0x04, 78, 97, 109, 101, cFixStr + 0x04, 116, 101, 115, 116, /**/
				cFixStr + 0x04, 100, 97, 116, 97, cFixBin + 0x04, 116, 101, 115, 116, /**/
				cFixStr + 0x04, 84, 101, 109, 112, cArr, cFixInt + 0x02, 129, 49, 129, 50, /**/
				cFixStr + 0x05, 67, 111, 117, 110, 116 /**/, 25,
			},
		},
		// {
		// 	dec: Corked{"test", []byte("test"), false, 25, "", ""},
		// 	enc: []byte{cFixExt + 0x0B, 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixBin + 0x04, 116, 101, 115, 116 /**/, 25},
		// },
		// {
		// 	dec: &Corked{"test", []byte("test"), false, 25, "", ""},
		// 	enc: []byte{cFixExt + 0x0B, 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixBin + 0x04, 116, 101, 115, 116 /**/, 25},
		// },
		// --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		{
			dec: []bool{true, false},
			enc: []byte{cArr, cFixInt + 0x02, cTrue, cFalse},
		},
		{
			dec: []string{"Hello", "World"},
			enc: []byte{cArr, cFixInt + 0x02 /**/, cFixStr + 0x05, 72, 101, 108, 108, 111 /**/, cFixStr + 0x05, 87, 111, 114, 108, 100},
		},
		{
			dec: []time.Time{clock, clock},
			enc: []byte{cArr, cFixInt + 0x02 /**/, cTime, 7, 166, 199, 91, 123, 67, 205, 21 /**/, cTime, 7, 166, 199, 91, 123, 67, 205, 21},
		},
		{
			dec: []int8{1, math.MaxInt8},
			enc: []byte{cArr, cFixInt + 0x02 /**/, 1 /**/, 127},
		},
		{
			dec: []int16{1, math.MaxInt8, math.MaxInt16},
			enc: []byte{cArr, cFixInt + 0x03 /**/, 1 /**/, 127 /**/, cInt16, 127, 255},
		},
		{
			dec: []int32{1, math.MaxInt8, math.MaxInt16, math.MaxInt32},
			enc: []byte{cArr, cFixInt + 0x04 /**/, 1 /**/, 127 /**/, cInt16, 127, 255 /**/, cInt32, 127, 255, 255, 255},
		},
		{
			dec: []int64{1, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64},
			enc: []byte{cArr, cFixInt + 0x05 /**/, 1 /**/, 127 /**/, cInt16, 127, 255 /**/, cInt32, 127, 255, 255, 255 /**/, cInt64, 127, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: []int{0, 1, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64},
			enc: []byte{cArr, cFixInt + 0x06 /**/, 0 /**/, 1 /**/, 127 /**/, cInt16, 127, 255 /**/, cInt32, 127, 255, 255, 255 /**/, cInt64, 127, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: []uint16{1, math.MaxUint8, math.MaxUint16},
			enc: []byte{cArr, cFixInt + 0x03 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255},
		},
		{
			dec: []uint32{1, math.MaxUint8, math.MaxUint16, math.MaxUint32},
			enc: []byte{cArr, cFixInt + 0x04 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255 /**/, cUint32, 255, 255, 255, 255},
		},
		{
			dec: []uint64{1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64},
			enc: []byte{cArr, cFixInt + 0x05 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255 /**/, cUint32, 255, 255, 255, 255 /**/, cUint64, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: []uint{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64},
			enc: []byte{cArr, cFixInt + 0x06 /**/, 0 /**/, 1 /**/, cUint8, 255 /**/, cUint16, 255, 255 /**/, cUint32, 255, 255, 255, 255 /**/, cUint64, 255, 255, 255, 255, 255, 255, 255, 255},
		},
		{
			dec: []float32{math.Pi, math.Pi},
			enc: []byte{cArr, cFixInt + 0x02 /**/, cFloat32, 64, 73, 15, 219 /**/, cFloat32, 64, 73, 15, 219},
		},
		{
			dec: []float64{math.Pi, math.Pi},
			enc: []byte{cArr, cFixInt + 0x02 /**/, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24 /**/, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			dec: []complex64{math.Pi, math.Pi},
			enc: []byte{cArr, cFixInt + 0x02 /**/, cComplex64, 64, 73, 15, 219, 0, 0, 0, 0 /**/, cComplex64, 64, 73, 15, 219, 0, 0, 0, 0},
		},
		{
			dec: []complex128{math.Pi, math.Pi},
			enc: []byte{cArr, cFixInt + 0x02 /**/, cComplex128, 64, 9, 33, 251, 84, 68, 45, 24, 0, 0, 0, 0, 0, 0, 0, 0 /**/, cComplex128, 64, 9, 33, 251, 84, 68, 45, 24, 0, 0, 0, 0, 0, 0, 0, 0},
		},
		{
			dec: []interface{}{nil, true, false, "test", []byte("test"), int8(77), uint8(177), float64(math.Pi)},
			enc: []byte{cArr, cFixInt + 0x08 /**/, cNil, cTrue, cFalse /**/, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixBin + 0x04, 116, 101, 115, 116 /**/, 77 /**/, cUint8, 177 /**/, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			dec: []Tested{{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}},
			enc: []byte{cArr, cFixInt + 0x01, cMap, cFixInt + 0x04, /**/
				cFixStr + 0x04, 78, 97, 109, 101, cFixStr + 0x04, 116, 101, 115, 116, /**/
				cFixStr + 0x04, 100, 97, 116, 97, cFixBin + 0x04, 116, 101, 115, 116, /**/
				cFixStr + 0x04, 84, 101, 109, 112, cArr, cFixInt + 0x02, 129, 49, 129, 50, /**/
				cFixStr + 0x05, 67, 111, 117, 110, 116 /**/, 25,
			},
		},
		{
			dec: []*Tested{{"test", []byte("test"), []string{"1", "2"}, false, 25, "", ""}},
			enc: []byte{cArr, cFixInt + 0x01, cMap, cFixInt + 0x04, /**/ /**/
				cFixStr + 0x04, 78, 97, 109, 101, cFixStr + 0x04, 116, 101, 115, 116, /**/
				cFixStr + 0x04, 100, 97, 116, 97, cFixBin + 0x04, 116, 101, 115, 116, /**/
				cFixStr + 0x04, 84, 101, 109, 112, cArr, cFixInt + 0x02, 129, 49, 129, 50, /**/
				cFixStr + 0x05, 67, 111, 117, 110, 116 /**/, 25,
			},
		},
		// {
		// 	dec: []Corked{{"test", []byte("test"), false, 25, "", ""}},
		// 	enc: []byte{cArr, cFixInt + 0x01, cFixExt + 0x0B, 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixBin + 0x04, 116, 101, 115, 116 /**/, 25},
		// },
		// {
		// 	dec: []Corked{{"test", []byte("test"), false, 25, "", ""}},
		// 	enc: []byte{cArr, cFixInt + 0x01, cFixExt + 0x0B, 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116 /**/, cFixBin + 0x04, 116, 101, 115, 116 /**/, 25},
		// },
		// --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		{
			dec: [][]float64{{math.Pi}, {math.Pi}},
			enc: []byte{cArr, cFixInt + 0x02 /**/, cArr, cFixInt + 0x01, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24 /**/, cArr, cFixInt + 0x01, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			dec: [][][]float64{{{math.Pi}}, {{math.Pi}}},
			enc: []byte{cArr, cFixInt + 0x02 /**/, cArr, cFixInt + 0x01, cArr, cFixInt + 0x01, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24 /**/, cArr, cFixInt + 0x01, cArr, cFixInt + 0x01, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		// --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
		{
			dec: map[string]bool{"test": true},
			enc: []byte{cMap, cFixInt + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cTrue},
		},
		{
			dec: map[string]int{"test": math.MaxInt32},
			enc: []byte{cMap, cFixInt + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cInt32, 127, 255, 255, 255},
		},
		{
			dec: map[string]string{"test": "Hello"},
			enc: []byte{cMap, cFixInt + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cFixStr + 0x05, 72, 101, 108, 108, 111},
		},
		{
			dec: map[string]interface{}{"test": math.Pi},
			enc: []byte{cMap, cFixInt + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			dec: map[interface{}]interface{}{"test": math.Pi},
			enc: []byte{cMap, cFixInt + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116, cFloat64, 64, 9, 33, 251, 84, 68, 45, 24},
		},
		{
			dec: map[string]interface{}{"test": map[string]interface{}{"test": "Embedded"}},
			enc: []byte{cMap, cFixInt + 0x01 /**/, cFixStr + 0x04, 116, 101, 115, 116 /**/, cMap, cFixInt + 0x01, cFixStr + 0x04, 116, 101, 115, 116, cFixStr + 0x08, 69, 109, 98, 101, 100, 100, 101, 100},
		},
	}

	// ----------------------------------------------------------------------------------------------------

	for _, test := range tests {

		Convey(fmt.Sprintf("%T will encode --- %v", test.dec, test.dec), t, func() {
			buf := bytes.NewBuffer(nil)
			err := NewEncoder(buf).Encode(test.dec)
			So(err, ShouldBeNil)
			So(buf.Bytes(), ShouldResemble, test.enc)
		})

	}

}
