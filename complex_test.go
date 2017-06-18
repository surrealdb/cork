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
	"math"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestComplex(t *testing.T) {

	tme, _ := time.Parse(time.RFC3339, "1987-06-22T08:00:00.123456789Z")

	Convey("*Custom will encode and decode", t, func() {

		var bit []byte

		var tmp Custom

		var val = &Custom{
			Null:            nil,
			Bool:            true,
			String:          "test",
			Bytes:           []byte("test"),
			Time:            tme,
			Int:             int(1),
			Int8:            int8(math.MaxInt8),
			Int16:           int16(math.MaxInt16),
			Int32:           int32(math.MaxInt32),
			Int64:           int64(math.MaxInt64),
			Uint:            uint(1),
			Uint8:           uint8(math.MaxUint8),
			Uint16:          uint16(math.MaxUint16),
			Uint32:          uint32(math.MaxUint32),
			Uint64:          uint64(math.MaxUint64),
			Float32:         math.MaxFloat32,
			Float64:         math.MaxFloat64,
			Complex64:       complex64(math.Pi),
			Complex128:      complex128(math.Pi),
			Any:             "test",
			ArrBool:         []bool{true, false},
			ArrString:       []string{"test", "test"},
			ArrInt:          []int{math.MinInt64, math.MinInt32, math.MinInt16, math.MinInt8, 0, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64},
			ArrInt8:         []int8{math.MinInt8, 0, math.MaxInt8},
			ArrInt16:        []int16{math.MinInt16, math.MinInt8, 0, math.MaxInt8, math.MaxInt16},
			ArrInt32:        []int32{math.MinInt32, math.MinInt16, math.MinInt8, 0, math.MaxInt8, math.MaxInt16, math.MaxInt32},
			ArrInt64:        []int64{math.MinInt64, math.MinInt32, math.MinInt16, math.MinInt8, 0, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64},
			ArrUint:         []uint{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64},
			ArrUint8:        []uint8{0, 1, math.MaxUint8},
			ArrUint16:       []uint16{0, 1, math.MaxUint8, math.MaxUint16},
			ArrUint32:       []uint32{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32},
			ArrUint64:       []uint64{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64},
			ArrFloat32:      []float32{math.Pi, math.Pi},
			ArrFloat64:      []float64{math.Pi, math.Pi},
			ArrComplex64:    []complex64{math.Pi, math.Pi},
			ArrComplex128:   []complex128{math.Pi, math.Pi},
			ArrTime:         []time.Time{tme, tme, tme},
			ArrAny:          []interface{}{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18"},
			ArrNo:           "test",
			MapStringInt:    map[string]int{"test": math.MaxInt64},
			MapStringUint:   map[string]uint{"test": math.MaxUint64},
			MapStringBool:   map[string]bool{"test": true},
			MapStringString: map[string]string{"test": "test"},
			MapIntAny:       map[int]interface{}{math.MaxInt64: "test"},
			MapUintAny:      map[uint]interface{}{math.MaxUint64: "test"},
			MapStringAny:    map[string]interface{}{"test": "test"},
			MapTimeAny:      map[time.Time]interface{}{tme: "test"},
			MapAnyAny: map[interface{}]interface{}{
				"1": "1", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6", "7": "7", "8": "8", "9": "9",
				"10": "10", "11": "11", "12": "12", "13": "13", "14": "14", "15": "15", "16": "16", "17": "17", "18": "18",
			},
			MapNo:         "test",
			CustomBool:    CustomBool(true),
			CustomString:  CustomString("test"),
			CustomInt:     CustomInt(1),
			CustomInt8:    CustomInt8(math.MaxInt8),
			CustomInt16:   CustomInt16(math.MaxInt16),
			CustomInt32:   CustomInt32(math.MaxInt32),
			CustomInt64:   CustomInt64(math.MaxInt64),
			CustomUint:    CustomUint(1),
			CustomUint8:   CustomUint8(math.MaxUint8),
			CustomUint16:  CustomUint16(math.MaxUint16),
			CustomUint32:  CustomUint32(math.MaxUint32),
			CustomUint64:  CustomUint64(math.MaxUint64),
			CustomFloat32: CustomFloat32(math.MaxFloat32),
			CustomFloat64: CustomFloat64(math.MaxFloat64),
			CustomAny:     CustomAny("test"),
			Tignored:      nil,
			Signored:      nil,
			Cignored:      nil,
			Testable:      &Tested{Data: []byte("test")},
			Selfable:      &Selfed{Data: []byte("test")},
			Corkable:      &Corked{Data: []byte("test")},
			Drrayble:      []Corked{{Data: []byte("test")}},
			Arrayble:      []*Corked{{Data: []byte("test")}},
		}

		val.Embedded.One = "test"
		val.Embedded.Two = 666
		val.Embedded.Ced = &Corked{Data: []byte("test")}
		val.Embedded.Sed = &Selfed{Data: []byte("test")}

		enc := NewEncoderBytes(&bit)
		eer := enc.Encode(val)
		So(eer, ShouldBeNil)

		dec := NewDecoderBytes(bit)
		der := dec.Decode(&tmp)
		So(der, ShouldBeNil)

		So(tmp, ShouldResemble, *val)

	})

	Convey("*Complex will encode and decode", t, func() {

		var val = &Complex{
			Bool:            true,
			String:          "test",
			Bytes:           []byte("test"),
			Time:            tme,
			Int:             -1,
			Int8:            math.MaxInt8,
			Int16:           math.MaxInt16,
			Int32:           math.MaxInt32,
			Int64:           math.MaxInt64,
			Uint:            1,
			Uint8:           math.MaxUint8,
			Uint16:          math.MaxUint16,
			Uint32:          math.MaxUint32,
			Uint64:          math.MaxUint64,
			Float32:         math.MaxFloat32,
			Float64:         math.MaxFloat64,
			Complex64:       complex64(math.Pi),
			Complex128:      complex128(math.Pi),
			Any:             "test",
			ArrBool:         []bool{true, false},
			ArrString:       []string{"test", "test"},
			ArrInt:          []int{math.MinInt64, math.MinInt32, math.MinInt16, math.MinInt8, 0, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64},
			ArrInt8:         []int8{math.MinInt8, 0, math.MaxInt8},
			ArrInt16:        []int16{math.MinInt16, math.MinInt8, 0, math.MaxInt8, math.MaxInt16},
			ArrInt32:        []int32{math.MinInt32, math.MinInt16, math.MinInt8, 0, math.MaxInt8, math.MaxInt16, math.MaxInt32},
			ArrInt64:        []int64{math.MinInt64, math.MinInt32, math.MinInt16, math.MinInt8, 0, math.MaxInt8, math.MaxInt16, math.MaxInt32, math.MaxInt64},
			ArrUint:         []uint{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64},
			ArrUint8:        []uint8{0, 1, math.MaxUint8},
			ArrUint16:       []uint16{0, 1, math.MaxUint8, math.MaxUint16},
			ArrUint32:       []uint32{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32},
			ArrUint64:       []uint64{0, 1, math.MaxUint8, math.MaxUint16, math.MaxUint32, math.MaxUint64},
			ArrFloat32:      []float32{math.Pi, math.Pi},
			ArrFloat64:      []float64{math.Pi, math.Pi},
			ArrComplex64:    []complex64{math.Pi, math.Pi},
			ArrComplex128:   []complex128{math.Pi, math.Pi},
			ArrTime:         []time.Time{tme, tme, tme},
			ArrAny:          []interface{}{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14", "15", "16", "17", "18"},
			ArrNo:           "test",
			MapStringInt:    map[string]int{"test": math.MaxInt64},
			MapStringUint:   map[string]uint{"test": math.MaxUint64},
			MapStringBool:   map[string]bool{"test": true},
			MapStringString: map[string]string{"test": "test"},
			MapIntAny:       map[int]interface{}{math.MaxInt64: "test"},
			MapUintAny:      map[uint]interface{}{math.MaxUint64: "test"},
			MapStringAny:    map[string]interface{}{"test": "test"},
			MapTimeAny:      map[time.Time]interface{}{tme: "test"},
			MapAnyAny: map[interface{}]interface{}{
				"1": "1", "2": "2", "3": "3", "4": "4", "5": "5", "6": "6", "7": "7", "8": "8", "9": "9",
				"10": "10", "11": "11", "12": "12", "13": "13", "14": "14", "15": "15", "16": "16", "17": "17", "18": "18",
			},
			MapNo: "test",
		}

		var dst Complex
		var tmp interface{}

		enb := bytes.NewBuffer(nil)
		enc := NewEncoder(enb)
		eer := enc.Encode(val)
		src := enb.Bytes()

		deb := bytes.NewReader(src)
		dec := NewDecoder(deb)
		der := dec.Decode(&dst)

		teb := bytes.NewReader(src)
		tec := NewDecoder(teb)
		ter := tec.Decode(&tmp)

		So(eer, ShouldBeNil)
		So(der, ShouldBeNil)
		So(ter, ShouldBeNil)

		So(src, ShouldNotBeNil)
		So(dst, ShouldNotBeNil)
		So(tmp, ShouldNotBeNil)

		So(dst, ShouldResemble, *val)
		So(tmp, ShouldResemble, val)

	})

}
