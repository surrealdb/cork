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
	"bytes"
	"reflect"
	"sort"
	"time"
)

type sortable struct {
	key []byte
	val interface{}
	ref reflect.Value
}

func sortMap(m reflect.Value) (a []*sortable) {
	for _, k := range m.MapKeys() {
		s := &sortable{ref: m.MapIndex(k)}
		NewEncoderBytes(&s.key).w.EncodeReflect(k)
		a = append(a, s)
	}
	sort.Slice(a, func(x, y int) bool {
		return bytes.Compare(a[x].key, a[y].key) < 0
	})
	return
}

func sortMapStringInt(m map[string]int) (a []string) {
	for k := range m {
		a = append(a, k)
	}
	sort.Strings(a)
	return
}

func sortMapStringUint(m map[string]uint) (a []string) {
	for k := range m {
		a = append(a, k)
	}
	sort.Strings(a)
	return
}

func sortMapStringBool(m map[string]bool) (a []string) {
	for k := range m {
		a = append(a, k)
	}
	sort.Strings(a)
	return
}

func sortMapStringString(m map[string]string) (a []string) {
	for k := range m {
		a = append(a, k)
	}
	sort.Strings(a)
	return
}

func sortMapIntAny(m map[int]interface{}) (a []int) {
	for k := range m {
		a = append(a, k)
	}
	sort.Ints(a)
	return
}

func sortMapUintAny(m map[uint]interface{}) (a []uint) {
	for k := range m {
		a = append(a, k)
	}
	sort.Slice(a, func(x, y int) bool {
		return a[x] < a[y]
	})
	return
}

func sortMapStringAny(m map[string]interface{}) (a []string) {
	for k := range m {
		a = append(a, k)
	}
	sort.Strings(a)
	return
}

func sortMapTimeAny(m map[time.Time]interface{}) (a []time.Time) {
	for k := range m {
		a = append(a, k)
	}
	sort.Slice(a, func(x, y int) bool {
		return a[x].UnixNano() < a[y].UnixNano()
	})
	return
}

func sortMapAnyAny(m map[interface{}]interface{}) (a []*sortable) {
	for k, v := range m {
		s := &sortable{val: v}
		NewEncoderBytes(&s.key).w.EncodeAny(k)
		a = append(a, s)
	}
	sort.Slice(a, func(x, y int) bool {
		return bytes.Compare(a[x].key, a[y].key) < 0
	})
	return
}
