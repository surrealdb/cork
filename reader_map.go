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

func (r *Reader) decodeMap(m reflect.Value) {
	t := m.Type()
	s := r.decodeMapLen()
	if s > 0 && m.IsNil() {
		m.Set(reflect.MakeMap(t))
	}
	for i := 0; i < s; i++ {
		k := reflect.New(t.Key())
		r.DecodeReflect(k)
		v := reflect.New(t.Elem())
		r.DecodeReflect(v)
		m.SetMapIndex(k.Elem(), v.Elem())
	}
}

func (r *Reader) decodeMapStringInt(m *map[string]int) {
	s := r.decodeMapLen()
	if s > 0 && *m == nil {
		*m = make(map[string]int, s)
	}
	for i := 0; i < s; i++ {
		var k string
		var v int
		r.DecodeString(&k)
		r.DecodeInt(&v)
		(*m)[k] = v
	}
}

func (r *Reader) decodeMapStringUint(m *map[string]uint) {
	s := r.decodeMapLen()
	if s > 0 && *m == nil {
		*m = make(map[string]uint, s)
	}
	for i := 0; i < s; i++ {
		var k string
		var v uint
		r.DecodeString(&k)
		r.DecodeUint(&v)
		(*m)[k] = v
	}
}

func (r *Reader) decodeMapStringBool(m *map[string]bool) {
	s := r.decodeMapLen()
	if s > 0 && *m == nil {
		*m = make(map[string]bool, s)
	}
	for i := 0; i < s; i++ {
		var k string
		var v bool
		r.DecodeString(&k)
		r.DecodeBool(&v)
		(*m)[k] = v
	}
}

func (r *Reader) decodeMapStringString(m *map[string]string) {
	s := r.decodeMapLen()
	if s > 0 && *m == nil {
		*m = make(map[string]string, s)
	}
	for i := 0; i < s; i++ {
		var k string
		var v string
		r.DecodeString(&k)
		r.DecodeString(&v)
		(*m)[k] = v
	}
}

func (r *Reader) decodeMapIntAny(m *map[int]interface{}) {
	s := r.decodeMapLen()
	if s > 0 && *m == nil {
		*m = make(map[int]interface{}, s)
	}
	for i := 0; i < s; i++ {
		var k int
		var v interface{}
		r.DecodeInt(&k)
		r.DecodeAny(&v)
		(*m)[k] = v
	}
}

func (r *Reader) decodeMapUintAny(m *map[uint]interface{}) {
	s := r.decodeMapLen()
	if s > 0 && *m == nil {
		*m = make(map[uint]interface{}, s)
	}
	for i := 0; i < s; i++ {
		var k uint
		var v interface{}
		r.DecodeUint(&k)
		r.DecodeAny(&v)
		(*m)[k] = v
	}
}

func (r *Reader) decodeMapStringAny(m *map[string]interface{}) {
	s := r.decodeMapLen()
	if s > 0 && *m == nil {
		*m = make(map[string]interface{}, s)
	}
	for i := 0; i < s; i++ {
		var k string
		var v interface{}
		r.DecodeString(&k)
		r.DecodeAny(&v)
		(*m)[k] = v
	}
}

func (r *Reader) decodeMapTimeAny(m *map[time.Time]interface{}) {
	s := r.decodeMapLen()
	if s > 0 && *m == nil {
		*m = make(map[time.Time]interface{}, s)
	}
	for i := 0; i < s; i++ {
		var k time.Time
		var v interface{}
		r.DecodeTime(&k)
		r.DecodeAny(&v)
		(*m)[k] = v
	}
}

func (r *Reader) decodeMapAnyAny(m *map[interface{}]interface{}) {
	s := r.decodeMapLen()
	if s > 0 && *m == nil {
		*m = make(map[interface{}]interface{}, s)
	}
	for i := 0; i < s; i++ {
		var k interface{}
		var v interface{}
		r.DecodeAny(&k)
		r.DecodeAny(&v)
		(*m)[k] = v
	}
}
