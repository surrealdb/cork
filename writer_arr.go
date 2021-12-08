// Copyright © SurrealDB Ltd
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
	"reflect"
	"time"
)

func (w *Writer) encodeArr(a reflect.Value) {
	w.encodeArrLen(a.Len())
	for i := 0; i < a.Len(); i++ {
		w.EncodeReflect(a.Index(i))
	}
}

func (w *Writer) encodeArrBool(a []bool) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeBool(v)
	}
}

func (w *Writer) encodeArrInt(a []int) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeInt(v)
	}
}

func (w *Writer) encodeArrInt8(a []int8) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeInt(int(v))
	}
}

func (w *Writer) encodeArrInt16(a []int16) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeInt(int(v))
	}
}

func (w *Writer) encodeArrInt32(a []int32) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeInt(int(v))
	}
}

func (w *Writer) encodeArrInt64(a []int64) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeInt(int(v))
	}
}

func (w *Writer) encodeArrUint(a []uint) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeUint(v)
	}
}

func (w *Writer) encodeArrUint8(a []uint8) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeUint(uint(v))
	}
}

func (w *Writer) encodeArrUint16(a []uint16) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeUint(uint(v))
	}
}

func (w *Writer) encodeArrUint32(a []uint32) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeUint(uint(v))
	}
}

func (w *Writer) encodeArrUint64(a []uint64) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeUint(uint(v))
	}
}

func (w *Writer) encodeArrString(a []string) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeString(v)
	}
}

func (w *Writer) encodeArrFloat32(a []float32) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeFloat32(v)
	}
}

func (w *Writer) encodeArrFloat64(a []float64) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeFloat64(v)
	}
}

func (w *Writer) encodeArrComplex64(a []complex64) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeComplex64(v)
	}
}

func (w *Writer) encodeArrComplex128(a []complex128) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeComplex128(v)
	}
}

func (w *Writer) encodeArrTime(a []time.Time) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeTime(v)
	}
}

func (w *Writer) encodeArrAny(a []interface{}) {
	w.encodeArrLen(len(a))
	for _, v := range a {
		w.EncodeAny(v)
	}
}
