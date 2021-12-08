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
	"encoding"
	"reflect"
	"time"
)

func (w *Writer) EncodeAny(v interface{}) {

	switch v := v.(type) {

	case Selfer:
		w.EncodeSelfer(v)
	case Corker:
		w.EncodeCorker(v)

	// -------------------------

	case nil:
		w.EncodeNil()
	case bool:
		w.EncodeBool(v)
	case byte:
		w.EncodeByte(v)
	case []byte:
		w.EncodeBytes(v)
	case string:
		w.EncodeString(v)
	case int:
		w.EncodeInt(v)
	case int8:
		w.EncodeInt(int(v))
	case int16:
		w.EncodeInt(int(v))
	case int32:
		w.EncodeInt(int(v))
	case int64:
		w.EncodeInt(int(v))
	case uint:
		w.EncodeUint(v)
	case uint16:
		w.EncodeUint(uint(v))
	case uint32:
		w.EncodeUint(uint(v))
	case uint64:
		w.EncodeUint(uint(v))
	case float32:
		w.EncodeFloat32(v)
	case float64:
		w.EncodeFloat64(v)
	case complex64:
		w.EncodeComplex64(v)
	case complex128:
		w.EncodeComplex128(v)
	case time.Time:
		w.EncodeTime(v)

	// -------------------------

	case []bool:
		w.encodeArrBool(v)
	case []int:
		w.encodeArrInt(v)
	case []int8:
		w.encodeArrInt8(v)
	case []int16:
		w.encodeArrInt16(v)
	case []int32:
		w.encodeArrInt32(v)
	case []int64:
		w.encodeArrInt64(v)
	case []uint:
		w.encodeArrUint(v)
	case []uint16:
		w.encodeArrUint16(v)
	case []uint32:
		w.encodeArrUint32(v)
	case []uint64:
		w.encodeArrUint64(v)
	case []string:
		w.encodeArrString(v)
	case []float32:
		w.encodeArrFloat32(v)
	case []float64:
		w.encodeArrFloat64(v)
	case []complex64:
		w.encodeArrComplex64(v)
	case []complex128:
		w.encodeArrComplex128(v)
	case []time.Time:
		w.encodeArrTime(v)
	case []interface{}:
		w.encodeArrAny(v)

	// -------------------------

	case map[string]int:
		w.encodeMapStringInt(v)
	case map[string]uint:
		w.encodeMapStringUint(v)
	case map[string]bool:
		w.encodeMapStringBool(v)
	case map[string]string:
		w.encodeMapStringString(v)
	case map[int]interface{}:
		w.encodeMapIntAny(v)
	case map[uint]interface{}:
		w.encodeMapUintAny(v)
	case map[string]interface{}:
		w.encodeMapStringAny(v)
	case map[time.Time]interface{}:
		w.encodeMapTimeAny(v)
	case map[interface{}]interface{}:
		w.encodeMapAnyAny(v)

	// -------------------------

	case encoding.BinaryMarshaler:
		enc, err := v.MarshalBinary()
		if err != nil {
			panic(err)
		}
		w.EncodeBytes(enc)

	case encoding.TextMarshaler:
		enc, err := v.MarshalText()
		if err != nil {
			panic(err)
		}
		w.EncodeBytes(enc)

	default:
		w.EncodeReflect(reflect.ValueOf(v))

	}

}
