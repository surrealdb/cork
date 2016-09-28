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
	"encoding/binary"
	"io"
	"math"
	"reflect"
	"sort"
	"time"
)

type Encoder struct {
	w *writer
}

// Encode encodes a data object into a CORK.
func Encode(src interface{}) (dst []byte) {
	buf := bytes.NewBuffer(dst)
	NewEncoder(buf).Encode(src)
	return buf.Bytes()
}

// NewEncoder returns an Encoder for encoding into an io.Writer.
func NewEncoder(dst io.Writer) *Encoder {
	return &Encoder{w: newWriter(dst)}
}

// Encode writes an object into a stream.
//
// Encoding can be configured via the struct tag for the fields.
// The "codec" key in struct field's tag value is the key name,
// followed by an optional comma and options.
// Note that the "json" key is used in the absence of the "codec" key.
//
// Struct values "usually" encode as maps. Each exported struct field is encoded unless:
//    - the field's tag is "-", OR
//    - the field is empty (empty or the zero value) and its tag specifies the "omitempty" option.
//
// When encoding as a map, the first string in the tag (before the comma)
// is the map key string to use when encoding.
//
// However, struct values may encode as arrays. This happens when:
//    - StructToArray Encode option is set, OR
//    - the tag on the _struct field sets the "toarray" option
//
// Values with types that implement MapBySlice are encoded as stream maps.
//
// The empty values (for omitempty option) are false, 0, any nil pointer
// or interface value, and any array, slice, map, or string of length zero.
//
// Anonymous fields are encoded inline except:
//    - the struct tag specifies a replacement name (first value)
//    - the field is of an interface type
//
// Examples:
//
//      // NOTE: 'json:' can be used as struct tag key, in place 'codec:' below.
//      type MyStruct struct {
//          _struct bool    `codec:",omitempty"`   //set omitempty for every field
//          Field1 string   `codec:"-"`            //skip this field
//          Field2 int      `codec:"myName"`       //Use key "myName" in encode stream
//          Field3 int32    `codec:",omitempty"`   //use key "Field3". Omit if empty.
//          Field4 bool     `codec:"f4,omitempty"` //use key "f4". Omit if empty.
//          io.Reader                              //use key "Reader".
//          MyStruct        `codec:"my1"           //use key "my1".
//          MyStruct                               //inline it
//          ...
//      }
//
//      type MyStruct struct {
//          _struct bool    `codec:",omitempty,toarray"`   //set omitempty for every field
//                                                         //and encode struct as an array
//      }
//
// The mode of encoding is based on the type of the value. When a value is seen:
//   - If a Selfer, call its CodecEncodeSelf method
//   - If an extension is registered for it, call that extension function
//   - If it implements encoding.(Binary|Text|JSON)Marshaler, call its Marshal(Binary|Text|JSON) method
//   - Else encode it based on its reflect.Kind
//
// Note that struct field names and keys in map[string]XXX will be treated as symbols.
// Some formats support symbols (e.g. binc) and will properly encode the string
// only once in the stream, and use a tag to refer to it thereafter.
//
func (e *Encoder) Encode(src interface{}) (err error) {
	// TODO need to catch panics and enable errors in encoder
	e.encode(src)
	return
}

func (e *Encoder) encode(src interface{}) {

	switch val := src.(type) {

	case Corker:
		e.encodeExt(val)

	case nil:
		e.encodeBit(cNil)

	case bool:
		e.encodeVal(val)

	case []byte:
		e.encodeBin(val)

	case string:
		e.encodeStr(val)

	case int:
		e.encodeInt(val)

	case int8:
		e.encodeInt(int(val))

	case int16:
		e.encodeInt(int(val))

	case int32:
		e.encodeInt(int(val))

	case int64:
		e.encodeInt(int(val))

	case uint:
		e.encodeUint(val)

	case uint8:
		e.encodeUint(uint(val))

	case uint16:
		e.encodeUint(uint(val))

	case uint32:
		e.encodeUint(uint(val))

	case uint64:
		e.encodeUint(uint(val))

	case float32:
		e.encodeBit(cFloat32)
		e.encodeFloat32(val)

	case float64:
		e.encodeBit(cFloat64)
		e.encodeFloat64(val)

	case time.Time:
		e.encodeBit(cTime)
		e.encodeTime(val)

	case []bool:
		e.encodeBit(cArrBool)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeVal(v)
		}

	case []string:
		e.encodeBit(cArrStr)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeStr(v)
		}

	case []int:
		e.encodeBit(cArrInt)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeInt(int(v))
		}

	case []int8:
		e.encodeBit(cArrInt8)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeInt(int(v))
		}

	case []int16:
		e.encodeBit(cArrInt16)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeInt(int(v))
		}

	case []int32:
		e.encodeBit(cArrInt32)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeInt(int(v))
		}

	case []int64:
		e.encodeBit(cArrInt64)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeInt(int(v))
		}

	case []uint:
		e.encodeBit(cArrUint)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeUint(uint(v))
		}

	case []uint16:
		e.encodeBit(cArrUint16)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeUint(uint(v))
		}

	case []uint32:
		e.encodeBit(cArrUint32)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeUint(uint(v))
		}

	case []uint64:
		e.encodeBit(cArrUint64)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeUint(uint(v))
		}

	case []float32:
		e.encodeBit(cArrFloat32)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeFloat32(v)
		}

	case []float64:
		e.encodeBit(cArrFloat64)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeFloat64(v)
		}

	case []time.Time:
		e.encodeBit(cArrTime)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encodeTime(v)
		}

	case []interface{}:
		e.encodeBit(cArrNil)
		e.encodeLen(len(val))
		for _, v := range val {
			e.encode(v)
		}

	case map[string]bool:
		set := make([]string, 0, len(val))
		e.encodeBit(cMapStrBool)
		e.encodeLen(len(val))
		for k, _ := range val {
			set = append(set, k)
		}
		sort.Strings(set)
		for _, v := range set {
			e.encodeStr(v)
			e.encodeVal(val[v])
		}

	case map[string]int:
		set := make([]string, 0, len(val))
		e.encodeBit(cMapStrInt)
		e.encodeLen(len(val))
		for k, _ := range val {
			set = append(set, k)
		}
		sort.Strings(set)
		for _, v := range set {
			e.encodeStr(v)
			e.encodeInt(val[v])
		}

	case map[string]string:
		set := make([]string, 0, len(val))
		e.encodeBit(cMapStrStr)
		e.encodeLen(len(val))
		for k, _ := range val {
			set = append(set, k)
		}
		sort.Strings(set)
		for _, v := range set {
			e.encodeStr(v)
			e.encodeStr(val[v])
		}

	case map[string]interface{}:
		set := make([]string, 0, len(val))
		e.encodeBit(cMapStrNil)
		e.encodeLen(len(val))
		for k, _ := range val {
			set = append(set, k)
		}
		sort.Strings(set)
		for _, v := range set {
			e.encodeStr(v)
			e.encode(val[v])
		}

	case map[interface{}]interface{}:
		e.encodeBit(cMapNilNil)
		e.encodeLen(len(val))
		for k, v := range val {
			e.encode(k)
			e.encode(v)
		}

	default:

		item := reflect.ValueOf(src)
		kind := reflect.TypeOf(src)

		switch kind.Kind() {

		case reflect.Ptr:
			item := item.Elem()
			if !item.IsValid() {
				e.encodeBit(cNil)
				return
			}
			e.encode(item.Interface())

		case reflect.Struct:
			flds := make([]*field, 0)
			e.encodeBit(cMapStrNil)
			for i := 0; i < item.NumField(); i++ {
				if fld := newField(kind.Field(i), item.Field(i)); fld != nil {
					flds = append(flds, fld)
				}
			}
			e.encodeLen(len(flds))
			for _, fld := range flds {
				e.encode(fld.Show())
				e.encode(item.FieldByName(fld.Name()).Interface())
			}

		case reflect.Slice:
			e.encodeBit(cArrNil)
			e.encodeLen(item.Len())
			for i := 0; i < reflect.ValueOf(src).Len(); i++ {
				e.encode(item.Index(i).Interface())
			}

		case reflect.Map:
			e.encodeBit(cMapStrNil)
			e.encodeLen(item.Len())
			for _, k := range item.MapKeys() {
				e.encode(k)
				e.encode(item.MapIndex(k))
			}

		}

	}

}

func (e *Encoder) encodeBit(val byte) {
	e.w.WriteOne(val)
	return
}

func (e *Encoder) encodeLen(val int) {
	e.encodeInt(val)
	return
}

func (e *Encoder) encodeVal(val bool) {
	if val {
		e.encodeBit(cTrue)
	} else {
		e.encodeBit(cFalse)
	}
	return
}

func (e *Encoder) encodeBin(val []byte) {
	sze := len(val)
	switch {
	case sze <= fixedBin:
		e.encodeBit(cFixBin + byte(sze))
	case sze <= math.MaxInt8:
		e.encodeBit(cBin8)
		binary.Write(e.w, binary.BigEndian, int8(sze))
	case sze <= math.MaxInt16:
		e.encodeBit(cBin16)
		binary.Write(e.w, binary.BigEndian, int16(sze))
	case sze <= math.MaxInt32:
		e.encodeBit(cBin32)
		binary.Write(e.w, binary.BigEndian, int32(sze))
	case sze <= math.MaxInt64:
		e.encodeBit(cBin64)
		binary.Write(e.w, binary.BigEndian, int64(sze))
	}
	e.w.WriteMany(val)
	return
}

func (e *Encoder) encodeStr(val string) {
	sze := len(val)
	switch {
	case sze <= fixedStr:
		e.encodeBit(cFixStr + byte(sze))
	case sze <= math.MaxInt8:
		e.encodeBit(cStr8)
		binary.Write(e.w, binary.BigEndian, int8(sze))
	case sze <= math.MaxInt16:
		e.encodeBit(cStr16)
		binary.Write(e.w, binary.BigEndian, int16(sze))
	case sze <= math.MaxInt32:
		e.encodeBit(cStr32)
		binary.Write(e.w, binary.BigEndian, int32(sze))
	case sze <= math.MaxInt64:
		e.encodeBit(cStr64)
		binary.Write(e.w, binary.BigEndian, int64(sze))
	}
	e.w.WriteManys(val)
	return
}

func (e *Encoder) encodeExt(val Corker) {
	enc, err := val.MarshalCORK()
	if err != nil {
		panic(err)
	}
	sze := len(enc)
	switch {
	case sze <= math.MaxInt8:
		e.encodeBit(cExt8)
		binary.Write(e.w, binary.BigEndian, int8(sze))
	case sze <= math.MaxInt16:
		e.encodeBit(cExt16)
		binary.Write(e.w, binary.BigEndian, int16(sze))
	case sze <= math.MaxInt32:
		e.encodeBit(cExt32)
		binary.Write(e.w, binary.BigEndian, int32(sze))
	case sze <= math.MaxInt64:
		e.encodeBit(cExt64)
		binary.Write(e.w, binary.BigEndian, int64(sze))
	}
	e.w.WriteOne(val.ExtendCORK())
	e.w.WriteMany(enc)
	return
}

func (e *Encoder) encodeInt(val int) {
	switch {
	case val >= 0 && val <= fixedInt:
		e.encodeBit(byte(val))
	case val <= math.MaxInt8:
		e.encodeBit(cInt8)
		binary.Write(e.w, binary.BigEndian, int8(val))
	case val <= math.MaxInt16:
		e.encodeBit(cInt16)
		binary.Write(e.w, binary.BigEndian, int16(val))
	case val <= math.MaxInt32:
		e.encodeBit(cInt32)
		binary.Write(e.w, binary.BigEndian, int32(val))
	case val <= math.MaxInt64:
		e.encodeBit(cInt64)
		binary.Write(e.w, binary.BigEndian, int64(val))
	}
	return
}

func (e *Encoder) encodeUint(val uint) {
	switch {
	case val >= 0 && val <= fixedInt:
		e.encodeBit(byte(val))
	case val <= math.MaxUint8:
		e.encodeBit(cUint8)
		binary.Write(e.w, binary.BigEndian, uint8(val))
	case val <= math.MaxUint16:
		e.encodeBit(cUint16)
		binary.Write(e.w, binary.BigEndian, uint16(val))
	case val <= math.MaxUint32:
		e.encodeBit(cUint32)
		binary.Write(e.w, binary.BigEndian, uint32(val))
	case val <= math.MaxUint64:
		e.encodeBit(cUint64)
		binary.Write(e.w, binary.BigEndian, uint64(val))
	}
	return
}

func (e *Encoder) encodeTime(val time.Time) {
	binary.Write(e.w, binary.BigEndian, val.UTC().UnixNano())
	return
}

func (e *Encoder) encodeFloat32(val float32) {
	binary.Write(e.w, binary.BigEndian, val)
	return
}

func (e *Encoder) encodeFloat64(val float64) {
	binary.Write(e.w, binary.BigEndian, val)
	return
}
