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
	"math"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestOptions(t *testing.T) {

	tme, _ := time.Parse(time.RFC3339, "1987-06-22T08:00:00.123456789Z")

	Convey("Can create []bool array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []bool{true, false}
		var opt = &Handle{ArrType: make([]bool, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []int array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []int{1, 2, 3}
		var opt = &Handle{ArrType: make([]int, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []int8 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []int8{1, 2, 3}
		var opt = &Handle{ArrType: make([]int8, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []int16 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []int16{1, 2, 3}
		var opt = &Handle{ArrType: make([]int16, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []int32 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []int32{1, 2, 3}
		var opt = &Handle{ArrType: make([]int32, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []int64 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []int64{1, 2, 3}
		var opt = &Handle{ArrType: make([]int64, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []uint array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []uint{1, 2, 3}
		var opt = &Handle{ArrType: make([]uint, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []uint8 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []uint8{1, 2, 3}
		var opt = &Handle{ArrType: make([]uint8, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []uint16 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []uint16{1, 2, 3}
		var opt = &Handle{ArrType: make([]uint16, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []uint32 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []uint32{1, 2, 3}
		var opt = &Handle{ArrType: make([]uint32, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []uint64 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []uint64{1, 2, 3}
		var opt = &Handle{ArrType: make([]uint64, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []string array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []string{"hello", "world"}
		var opt = &Handle{ArrType: make([]string, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []float32 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []float32{math.Pi, math.Pi}
		var opt = &Handle{ArrType: make([]float32, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []float64 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []float64{math.Pi, math.Pi}
		var opt = &Handle{ArrType: make([]float64, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []complex64 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []complex64{math.Pi, math.Pi}
		var opt = &Handle{ArrType: make([]complex64, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []complex128 array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []complex128{math.Pi, math.Pi}
		var opt = &Handle{ArrType: make([]complex128, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []time.Time array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []time.Time{tme, tme}
		var opt = &Handle{ArrType: make([]time.Time, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create []interface{} array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []interface{}{true, 1, 2, 3, "test", math.Pi, tme}
		var opt = &Handle{ArrType: make([]interface{}, 0)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create reflect array type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = []Simple{}
		var opt = &Handle{ArrType: reflect.TypeOf(val)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create map[string]int map type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = map[string]int{"test": 1}
		var opt = &Handle{MapType: make(map[string]int)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create map[string]uint map type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = map[string]uint{"test": 1}
		var opt = &Handle{MapType: make(map[string]uint)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create map[string]bool map type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = map[string]bool{"test": true}
		var opt = &Handle{MapType: make(map[string]bool)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create map[string]string map type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = map[string]string{"test": "hello world"}
		var opt = &Handle{MapType: make(map[string]string)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create map[int]interface{} map type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = map[int]interface{}{1: "test"}
		var opt = &Handle{MapType: make(map[int]interface{})}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create map[uint]interface{} map type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = map[uint]interface{}{1: "test"}
		var opt = &Handle{MapType: make(map[uint]interface{})}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create map[string]interface{} map type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = map[string]interface{}{"test": "test"}
		var opt = &Handle{MapType: make(map[string]interface{})}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create map[time.Time]interface{} map type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = map[time.Time]interface{}{tme: "test"}
		var opt = &Handle{MapType: make(map[time.Time]interface{})}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create map[interface{}]interface{} map type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = map[interface{}]interface{}{tme: "test", 1: true, math.Pi: 3}
		var opt = &Handle{MapType: make(map[interface{}]interface{})}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

	Convey("Can create reflect map type", t, func() {
		var tmp interface{}
		var buf []byte
		var val = map[string]Tested{"test": {Data: []byte("test")}}
		var opt = &Handle{MapType: reflect.TypeOf(val)}
		NewEncoderBytes(&buf).Options(opt).Encode(val)
		NewDecoderBytes(buf).Options(opt).Decode(&tmp)
		So(tmp, ShouldResemble, val)
	})

}
