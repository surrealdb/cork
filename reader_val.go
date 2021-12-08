// Copyright © SurrealDB Ltd
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://wwr.deache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cork

import (
	"encoding/binary"
	"reflect"
	"time"
)

// DecodeInterface decodes any value from the Reader, with
// whatever type is next in the stream.
func (r *Reader) DecodeInterface(v *interface{}) {

	b := r.peekOne()

	switch {
	case b == cNil:
		r.readOne()
		*v = nil
	case isBool(b):
		var x bool
		r.DecodeBool(&x)
		*v = x
	case isTime(b):
		var x time.Time
		r.DecodeTime(&x)
		*v = x
	case isBin(b):
		var x []byte
		r.DecodeBytes(&x)
		*v = x
	case isStr(b):
		var x string
		r.DecodeString(&x)
		*v = x
	case isNum(b):
		var x int
		r.DecodeInt(&x)
		*v = x
	case isInt(b):
		var x int
		r.DecodeInt(&x)
		*v = x
	case isUint(b):
		var x uint
		r.DecodeUint(&x)
		*v = x
	case b == cFloat32:
		var x float32
		r.DecodeFloat32(&x)
		*v = x
	case b == cFloat64:
		var x float64
		r.DecodeFloat64(&x)
		*v = x
	case b == cComplex64:
		var x complex64
		r.DecodeComplex64(&x)
		*v = x
	case b == cComplex128:
		var x complex128
		r.DecodeComplex128(&x)
		*v = x

	// -------------------------

	case isExt(b):
		*v = r.createExt()
	case isSlf(b):
		*v = r.createSlf()
	case isArr(b):
		*v = r.createArr()
	case isMap(b):
		*v = r.createMap()

	// -------------------------

	default:
		panic(fail)

	}

}

func (r *Reader) createExt() (v Corker) {
	s := 0
	b := r.readOne()
	switch {
	case b >= cFixExt && b <= cFixExt+fixedExt:
		s = int(b - cFixExt)
	case b == cExt8:
		s = int(r.readOne())
	case b == cExt16:
		s = int(binary.BigEndian.Uint16(r.readMany(2)))
	case b == cExt32:
		s = int(binary.BigEndian.Uint32(r.readMany(4)))
	case b == cExt64:
		s = int(binary.BigEndian.Uint64(r.readMany(8)))
	default:
		panic(fail)
	}
	e := r.readOne()
	d := r.readMany(s)
	v = reflect.New(registry[e]).Interface().(Corker)
	if err := v.UnmarshalCORK(d); err != nil {
		panic(err)
	}
	return
}

func (r *Reader) createSlf() (v Selfer) {
	if r.readOne() == cSlf {
		e := r.readOne()
		v = reflect.New(registry[e]).Interface().(Selfer)
		v.UnmarshalCORK(r)
		return
	}
	panic(fail)
}

func (r *Reader) createArr() (v interface{}) {
	if r.h != nil && r.h.ArrType != nil {
		switch a := r.h.ArrType.(type) {
		case []bool:
			var x []bool
			r.decodeArrBool(&x)
			return x
		case []int:
			var x []int
			r.decodeArrInt(&x)
			return x
		case []int8:
			var x []int8
			r.decodeArrInt8(&x)
			return x
		case []int16:
			var x []int16
			r.decodeArrInt16(&x)
			return x
		case []int32:
			var x []int32
			r.decodeArrInt32(&x)
			return x
		case []int64:
			var x []int64
			r.decodeArrInt64(&x)
			return x
		case []uint:
			var x []uint
			r.decodeArrUint(&x)
			return x
		case []uint16:
			var x []uint16
			r.decodeArrUint16(&x)
			return x
		case []uint32:
			var x []uint32
			r.decodeArrUint32(&x)
			return x
		case []uint64:
			var x []uint64
			r.decodeArrUint64(&x)
			return x
		case []string:
			var x []string
			r.decodeArrString(&x)
			return x
		case []float32:
			var x []float32
			r.decodeArrFloat32(&x)
			return x
		case []float64:
			var x []float64
			r.decodeArrFloat64(&x)
			return x
		case []complex64:
			var x []complex64
			r.decodeArrComplex64(&x)
			return x
		case []complex128:
			var x []complex128
			r.decodeArrComplex128(&x)
			return x
		case []time.Time:
			var x []time.Time
			r.decodeArrTime(&x)
			return x
		case reflect.Type:
			var x = reflect.MakeSlice(a, 0, 0)
			r.decodeArr(x)
			return x.Interface()
		}
	}
	var x []interface{}
	r.decodeArrAny(&x)
	return x
}

func (r *Reader) createMap() (v interface{}) {
	if r.h != nil && r.h.MapType != nil {
		switch m := r.h.MapType.(type) {
		default:
		case map[string]int:
			var x map[string]int
			r.decodeMapStringInt(&x)
			return x
		case map[string]uint:
			var x map[string]uint
			r.decodeMapStringUint(&x)
			return x
		case map[string]bool:
			var x map[string]bool
			r.decodeMapStringBool(&x)
			return x
		case map[string]string:
			var x map[string]string
			r.decodeMapStringString(&x)
			return x
		case map[int]interface{}:
			var x map[int]interface{}
			r.decodeMapIntAny(&x)
			return x
		case map[uint]interface{}:
			var x map[uint]interface{}
			r.decodeMapUintAny(&x)
			return x
		case map[string]interface{}:
			var x map[string]interface{}
			r.decodeMapStringAny(&x)
			return x
		case map[time.Time]interface{}:
			var x map[time.Time]interface{}
			r.decodeMapTimeAny(&x)
			return x
		case reflect.Type:
			var x = reflect.MakeMap(m)
			r.decodeMap(x)
			return x.Interface()
		}
	}
	var x map[interface{}]interface{}
	r.decodeMapAnyAny(&x)
	return x
}
