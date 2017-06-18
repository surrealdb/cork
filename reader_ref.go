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

// DecodeReflect decodes a reflect.Value value from the Reader.
func (r *Reader) DecodeReflect(v reflect.Value) {

	b := r.peekOne()

	if b == cNil {
		r.readOne()
		return
	}

	t := v.Type()
	k := v.Kind()

	// First let's check to see if this is
	// a nil pointer, and if it is then we
	// will create a new value for the
	// underlying type, and set the pointer
	// to this value.

	if v.Kind() == reflect.Ptr && v.IsNil() {
		n := reflect.New(t.Elem())
		v.Set(n)
	}

	// Next let's check to see if the type
	// implements either the Selfer or Corker
	// interfaces, and if it does then decode
	// it directly. Caching the interface
	// detection speeds up the decoding.

	if c.Selfable(t) {
		n := reflect.New(t.Elem())
		r.DecodeSelfer(n.Interface().(Selfer))
		v.Set(n)
		return
	}

	if c.Corkable(t) {
		n := reflect.New(t.Elem())
		r.DecodeCorker(n.Interface().(Corker))
		v.Set(n)
		return
	}

	// It wasn't a self describing interface
	// so let's now see if it is a string or
	// a byte slice, and if it is, then
	// decode it immediately.

	switch t {

	case typeStr:
		var x string
		r.DecodeString(&x)
		v.SetString(x)
		return

	case typeBit:
		var x []byte
		r.DecodeBytes(&x)
		v.SetBytes(x)
		return

	case typeTime:
		var x time.Time
		r.DecodeTime(&x)
		v.Set(reflect.ValueOf(x))
		return

	}

	// Otherwise let's switch over all of the
	// possible types that this item can be
	// and decode it into the correct type.
	// For structs, we will cache the struct
	// fields, so that we do not have to parse
	// these for every item that we process.

	switch k {

	case reflect.Ptr:
		r.DecodeReflect(v.Elem())

	case reflect.Map:
		r.decodeMap(v)

	case reflect.Slice:
		r.decodeArr(v)

	case reflect.Bool:
		var x bool
		r.DecodeBool(&x)
		v.SetBool(x)

	case reflect.String:
		var x string
		r.DecodeString(&x)
		v.SetString(x)

	case reflect.Int:
		var x int
		r.DecodeInt(&x)
		v.SetInt(int64(x))

	case reflect.Int8:
		var x int8
		r.DecodeInt8(&x)
		v.SetInt(int64(x))

	case reflect.Int16:
		var x int16
		r.DecodeInt16(&x)
		v.SetInt(int64(x))

	case reflect.Int32:
		var x int32
		r.DecodeInt32(&x)
		v.SetInt(int64(x))

	case reflect.Int64:
		var x int64
		r.DecodeInt64(&x)
		v.SetInt(x)

	case reflect.Uint:
		var x uint
		r.DecodeUint(&x)
		v.SetUint(uint64(x))

	case reflect.Uint8:
		var x uint8
		r.DecodeUint8(&x)
		v.SetUint(uint64(x))

	case reflect.Uint16:
		var x uint16
		r.DecodeUint16(&x)
		v.SetUint(uint64(x))

	case reflect.Uint32:
		var x uint32
		r.DecodeUint32(&x)
		v.SetUint(uint64(x))

	case reflect.Uint64:
		var x uint64
		r.DecodeUint64(&x)
		v.SetUint(x)

	case reflect.Float32:
		var x float32
		r.DecodeFloat32(&x)
		v.SetFloat(float64(x))

	case reflect.Float64:
		var x float64
		r.DecodeFloat64(&x)
		v.SetFloat(x)

	case reflect.Complex64:
		var x complex64
		r.DecodeComplex64(&x)
		v.SetComplex(complex128(x))

	case reflect.Complex128:
		var x complex128
		r.DecodeComplex128(&x)
		v.SetComplex(x)

	case reflect.Interface:
		var x interface{}
		r.DecodeInterface(&x)
		if reflect.ValueOf(x).IsValid() {
			v.Set(reflect.ValueOf(x))
		}

	case reflect.Struct:

		if !c.Has(t) {
			tot := 0
			fls := make([]*field, t.NumField())
			for i := 0; i < t.NumField(); i++ {
				if f := newField(t.Field(i)); f != nil {
					fls[tot] = f
					tot++
				}
			}
			c.Set(t, fls[:tot])
		}

		x := c.Get(t)
		s := r.decodeMapLen()

		for i := 0; i < s; i++ {

			var k string
			r.DecodeString(&k)

			for _, f := range x {
				if k == f.Name() {
					if f := v.FieldByIndex(f.indx); f.CanSet() {
						if v.CanAddr() {
							r.DecodeReflect(f.Addr())
						} else {
							r.DecodeReflect(f)
						}
					}
					continue
				}
			}

		}

	}

}
