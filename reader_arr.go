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
	"reflect"
	"time"
)

func (r *Reader) decodeArr(a reflect.Value) {
	t := a.Type()
	s := r.decodeArrLen()
	if s > 0 && (a.IsNil() || a.Len() < s) {
		a.Set(reflect.MakeSlice(t, s, s))
	}
	for i := 0; i < s; i++ {
		r.DecodeReflect(a.Index(i))
	}
}

func (r *Reader) decodeArrBool(a *[]bool) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]bool, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeBool(&(*a)[i])
	}
}

func (r *Reader) decodeArrInt(a *[]int) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]int, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeInt(&(*a)[i])
	}
}

func (r *Reader) decodeArrInt8(a *[]int8) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]int8, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeInt8(&(*a)[i])
	}
}

func (r *Reader) decodeArrInt16(a *[]int16) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]int16, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeInt16(&(*a)[i])
	}
}

func (r *Reader) decodeArrInt32(a *[]int32) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]int32, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeInt32(&(*a)[i])
	}
}

func (r *Reader) decodeArrInt64(a *[]int64) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]int64, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeInt64(&(*a)[i])
	}
}

func (r *Reader) decodeArrUint(a *[]uint) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]uint, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeUint(&(*a)[i])
	}
}

func (r *Reader) decodeArrUint8(a *[]uint8) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]uint8, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeUint8(&(*a)[i])
	}
}

func (r *Reader) decodeArrUint16(a *[]uint16) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]uint16, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeUint16(&(*a)[i])
	}
}

func (r *Reader) decodeArrUint32(a *[]uint32) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]uint32, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeUint32(&(*a)[i])
	}
}

func (r *Reader) decodeArrUint64(a *[]uint64) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]uint64, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeUint64(&(*a)[i])
	}
}

func (r *Reader) decodeArrString(a *[]string) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]string, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeString(&(*a)[i])
	}
}

func (r *Reader) decodeArrFloat32(a *[]float32) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]float32, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeFloat32(&(*a)[i])
	}
}

func (r *Reader) decodeArrFloat64(a *[]float64) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]float64, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeFloat64(&(*a)[i])
	}
}

func (r *Reader) decodeArrComplex64(a *[]complex64) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]complex64, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeComplex64(&(*a)[i])
	}
}

func (r *Reader) decodeArrComplex128(a *[]complex128) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]complex128, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeComplex128(&(*a)[i])
	}
}

func (r *Reader) decodeArrTime(a *[]time.Time) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]time.Time, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeTime(&(*a)[i])
	}
}

func (r *Reader) decodeArrAny(a *[]interface{}) {
	s := r.decodeArrLen()
	if s > 0 && (*a == nil || len(*a) < s) {
		*a = make([]interface{}, s)
	}
	for i := 0; i < s; i++ {
		r.DecodeAny(&(*a)[i])
	}
}
