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

// EncodeReflect encodes a reflect.Value value to the Writer.
func (w *Writer) EncodeReflect(v reflect.Value) {

	// If the element is a function or a
	// channel, then we can ignore these
	// types as these are not able to be
	// encoded. In addition, if the value
	// is a nil pointer or interface then
	// we can encode nil immediately.

	k := v.Kind()

	switch k {
	case reflect.Func, reflect.Chan:
		w.EncodeNil()
		return
	case reflect.Ptr, reflect.Interface:
		if v.IsNil() {
			w.EncodeNil()
			return
		}
	}

	// Next let's check to see if the type
	// implements either the Selfer or Corker
	// interfaces, and if it does then encode
	// it directly. Caching the interface
	// detection speeds up the encoding.

	t := v.Type()

	if c.Selfable(t) {
		w.EncodeSelfer(v.Interface().(Selfer))
		return
	}

	if c.Corkable(t) {
		w.EncodeCorker(v.Interface().(Corker))
		return
	}

	// It wasn't a self describing interface
	// so let's now see if it is a string or
	// a byte slice, and if it is, then
	// encode it immediately.

	switch t {

	case typeStr:
		w.EncodeString(v.String())
		return

	case typeBit:
		w.EncodeBytes(v.Bytes())
		return

	case typeTime:
		w.EncodeTime(v.Interface().(time.Time))
		return

	}

	// Otherwise let's switch over all of the
	// possible types that this item can be
	// and encode it into the correct type.
	// For structs, we will cache the struct
	// fields, so that we do not have to parse
	// these for every item that we process.

	switch k {

	case reflect.Ptr:
		w.EncodeReflect(v.Elem())

	case reflect.Map:
		w.encodeMap(v)

	case reflect.Slice:
		w.encodeArr(v)

	case reflect.Bool:
		w.EncodeBool(v.Bool())

	case reflect.String:
		w.EncodeString(v.String())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		w.EncodeInt(int(v.Int()))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		w.EncodeUint(uint(v.Uint()))

	case reflect.Float32:
		w.EncodeFloat32(float32(v.Float()))

	case reflect.Float64:
		w.EncodeFloat64(v.Float())

	case reflect.Complex64:
		w.EncodeComplex64(complex64(v.Complex()))

	case reflect.Complex128:
		w.EncodeComplex128(v.Complex())

	case reflect.Interface:
		w.EncodeAny(v.Interface())

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

		sze, fls := 0, c.Get(t)

		for _, f := range fls {
			if v := v.FieldByIndex(f.indx); v.IsValid() {
				if !f.omit || (f.omit && !isEmpty(v)) {
					sze++
				}
			}
		}

		w.encodeMapLen(sze)

		for _, f := range fls {
			if v := v.FieldByIndex(f.indx); v.IsValid() {
				if !f.omit || (f.omit && !isEmpty(v)) {
					w.EncodeString(f.Name())
					w.EncodeReflect(v)
				}
			}
		}

	}

}
