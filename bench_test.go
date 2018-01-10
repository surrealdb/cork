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
	"reflect"
	"testing"
	"time"

	"github.com/ugorji/go/codec"
)

var hl = Handle{
	ArrType: make([]interface{}, 0),
	MapType: make(map[string]interface{}),
}

var (
	jl codec.JsonHandle
	cl codec.CborHandle
	ml codec.MsgpackHandle
)

var bit []byte

var dec interface{}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

type Str struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

type Crk struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

func (this *Crk) ExtendCORK() byte {
	return 0xB0
}

func (this *Crk) MarshalCORK() (val []byte, err error) {
	return
}

func (this *Crk) UnmarshalCORK(val []byte) (err error) {
	return
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

type Slf struct {
	Name     string
	BirthDay time.Time
	Phone    string
	Siblings int
	Spouse   bool
	Money    float64
}

func (this *Slf) ExtendCORK() byte {
	return 0xB1
}

func (this *Slf) MarshalCORK(w *Writer) error {
	w.EncodeString(this.Name)
	w.EncodeTime(this.BirthDay)
	w.EncodeString(this.Phone)
	w.EncodeInt(this.Siblings)
	w.EncodeBool(this.Spouse)
	w.EncodeFloat64(this.Money)
	return nil
}

func (this *Slf) UnmarshalCORK(r *Reader) error {
	r.DecodeString(&this.Name)
	r.DecodeTime(&this.BirthDay)
	r.DecodeString(&this.Phone)
	r.DecodeInt(&this.Siblings)
	r.DecodeBool(&this.Spouse)
	r.DecodeFloat64(&this.Money)
	return nil
}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func init() {

	Register(&Str{})
	Register(&Crk{})
	Register(&Slf{})

	jl.Canonical = true
	jl.InternString = true
	jl.HTMLCharsAsIs = true
	jl.CheckCircularRef = false
	jl.SliceType = reflect.TypeOf([]interface{}(nil))
	jl.MapType = reflect.TypeOf(map[string]interface{}(nil))

	cl.Canonical = true
	cl.InternString = true
	cl.CheckCircularRef = false
	cl.SliceType = reflect.TypeOf([]interface{}(nil))
	cl.MapType = reflect.TypeOf(map[string]interface{}(nil))

	ml.WriteExt = true
	ml.Canonical = true
	ml.RawToString = true
	ml.InternString = true
	ml.CheckCircularRef = false
	ml.SliceType = reflect.TypeOf([]interface{}(nil))
	ml.MapType = reflect.TypeOf(map[string]interface{}(nil))

}

// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------
// --------------------------------------------------

func BenchmarkCorkEncodeObject(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	var obj = map[string]interface{}{"Name": "test", "data": []byte("test"), "Temp": []interface{}{"1", "2"}, "Count": int64(25)}
	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(nil)
		crk := NewEncoderFromPool(buf).Options(&hl)
		crk.Encode(obj)
		crk.Reset()
		bit = buf.Bytes()
	}
}

func BenchmarkCorkDecodeObject(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		buf := bytes.NewReader(bit)
		crk := NewDecoderFromPool(buf).Options(&hl)
		crk.Decode(&dec)
		crk.Reset()
	}
}

func BenchmarkCborEncodeObject(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	var obj = map[string]interface{}{"Name": "test", "data": []byte("test"), "Temp": []interface{}{"1", "2"}, "Count": int64(25)}
	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(nil)
		crk := codec.NewEncoder(buf, &cl)
		crk.Encode(obj)
		bit = buf.Bytes()
	}
}

func BenchmarkCborDecodeObject(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		buf := bytes.NewReader(bit)
		crk := codec.NewDecoder(buf, &cl)
		crk.Decode(&dec)
	}
}

func BenchmarkPackEncodeObject(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	var obj = map[string]interface{}{"Name": "test", "data": []byte("test"), "Temp": []interface{}{"1", "2"}, "Count": int64(25)}
	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(nil)
		crk := codec.NewEncoder(buf, &ml)
		crk.Encode(obj)
		bit = buf.Bytes()
	}
}

func BenchmarkPackDecodeObject(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		buf := bytes.NewReader(bit)
		crk := codec.NewDecoder(buf, &ml)
		crk.Decode(&dec)
	}
}

// ----------------------------------------------------------------------------------------------------

func BenchmarkCorkEncodeStruct(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	var obj = &Str{Name: "Tobie Morgan Hitchcock", BirthDay: time.Now(), Phone: "+44 7931 739579", Siblings: 3, Spouse: true, Money: 13336183.419}
	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(nil)
		crk := NewEncoderFromPool(buf).Options(&hl)
		crk.Encode(obj)
		crk.Reset()
		bit = buf.Bytes()
	}
}

func BenchmarkCorkDecodeStruct(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		buf := bytes.NewReader(bit)
		crk := NewDecoderFromPool(buf).Options(&hl)
		crk.Decode(&dec)
		crk.Reset()
	}
}

func BenchmarkCorkEncodeCorker(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	var obj = &Crk{Name: "Tobie Morgan Hitchcock", BirthDay: time.Now(), Phone: "+44 7931 739579", Siblings: 3, Spouse: true, Money: 13336183.419}
	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(nil)
		crk := NewEncoderFromPool(buf).Options(&hl)
		crk.Encode(obj)
		crk.Reset()
		bit = buf.Bytes()
	}
}

func BenchmarkCorkDecodeCorker(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		buf := bytes.NewReader(bit)
		crk := NewDecoderFromPool(buf).Options(&hl)
		crk.Decode(&dec)
		crk.Reset()
	}
}

func BenchmarkCorkEncodeSelfer(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	var obj = &Slf{Name: "Tobie Morgan Hitchcock", BirthDay: time.Now(), Phone: "+44 7931 739579", Siblings: 3, Spouse: true, Money: 13336183.419}
	for n := 0; n < b.N; n++ {
		buf := bytes.NewBuffer(nil)
		crk := NewEncoderFromPool(buf).Options(&hl)
		crk.Encode(obj)
		crk.Reset()
		bit = buf.Bytes()
	}
}

func BenchmarkCorkDecodeSelfer(b *testing.B) {
	b.SetBytes(2)
	b.ReportAllocs()
	for n := 0; n < b.N; n++ {
		buf := bytes.NewReader(bit)
		crk := NewDecoderFromPool(buf).Options(&hl)
		crk.Decode(&dec)
		crk.Reset()
	}
}
