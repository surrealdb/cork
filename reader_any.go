// Copyright © 2016 Abcum Ltd
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
	"encoding"
	"reflect"
	"time"
)

// DecodeAny decodes a value from the Reader.
func (r *Reader) DecodeAny(v interface{}) {

	switch v := v.(type) {

	case Selfer:
		r.DecodeSelfer(v)
	case Corker:
		r.DecodeCorker(v)

	// -------------------------

	case *bool:
		r.DecodeBool(v)
	case *byte:
		r.DecodeByte(v)
	case *[]byte:
		r.DecodeBytes(v)
	case *string:
		r.DecodeString(v)
	case *int:
		r.DecodeInt(v)
	case *int8:
		r.DecodeInt8(v)
	case *int16:
		r.DecodeInt16(v)
	case *int32:
		r.DecodeInt32(v)
	case *int64:
		r.DecodeInt64(v)
	case *uint:
		r.DecodeUint(v)
	case *uint16:
		r.DecodeUint16(v)
	case *uint32:
		r.DecodeUint32(v)
	case *uint64:
		r.DecodeUint64(v)
	case *float32:
		r.DecodeFloat32(v)
	case *float64:
		r.DecodeFloat64(v)
	case *complex64:
		r.DecodeComplex64(v)
	case *complex128:
		r.DecodeComplex128(v)
	case *time.Time:
		r.DecodeTime(v)

	// -------------------------

	case *[]bool:
		r.decodeArrBool(v)
	case *[]int:
		r.decodeArrInt(v)
	case *[]int8:
		r.decodeArrInt8(v)
	case *[]int16:
		r.decodeArrInt16(v)
	case *[]int32:
		r.decodeArrInt32(v)
	case *[]int64:
		r.decodeArrInt64(v)
	case *[]uint:
		r.decodeArrUint(v)
	case *[]uint16:
		r.decodeArrUint16(v)
	case *[]uint32:
		r.decodeArrUint32(v)
	case *[]uint64:
		r.decodeArrUint64(v)
	case *[]string:
		r.decodeArrString(v)
	case *[]float32:
		r.decodeArrFloat32(v)
	case *[]float64:
		r.decodeArrFloat64(v)
	case *[]complex64:
		r.decodeArrComplex64(v)
	case *[]complex128:
		r.decodeArrComplex128(v)
	case *[]time.Time:
		r.decodeArrTime(v)
	case *[]interface{}:
		r.decodeArrAny(v)

	// -------------------------

	case *map[string]int:
		r.decodeMapStringInt(v)
	case *map[string]uint:
		r.decodeMapStringUint(v)
	case *map[string]bool:
		r.decodeMapStringBool(v)
	case *map[string]string:
		r.decodeMapStringString(v)
	case *map[int]interface{}:
		r.decodeMapIntAny(v)
	case *map[uint]interface{}:
		r.decodeMapUintAny(v)
	case *map[string]interface{}:
		r.decodeMapStringAny(v)
	case *map[time.Time]interface{}:
		r.decodeMapTimeAny(v)
	case *map[interface{}]interface{}:
		r.decodeMapAnyAny(v)

	// -------------------------

	case *interface{}:
		r.DecodeInterface(v)

	// -------------------------

	case encoding.BinaryUnmarshaler:
		var enc []byte
		r.DecodeBytes(&enc)
		err := v.UnmarshalBinary(enc)
		if err != nil {
			panic(err)
		}

	case encoding.TextUnmarshaler:
		var enc []byte
		r.DecodeBytes(&enc)
		err := v.UnmarshalText(enc)
		if err != nil {
			panic(err)
		}

	default:
		r.DecodeReflect(reflect.ValueOf(v))

	}

}
