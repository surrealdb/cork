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

func (w *Writer) encodeMap(m reflect.Value) {
	w.encodeMapLen(m.Len())
	if w.h != nil && w.h.SortMaps {
		for _, v := range sortMap(m) {
			w.writeMany(v.key)
			w.EncodeReflect(v.ref)
		}
	} else {
		for _, k := range m.MapKeys() {
			w.EncodeReflect(k)
			w.EncodeReflect(m.MapIndex(k))
		}
	}
}

func (w *Writer) encodeMapStringInt(m map[string]int) {
	w.encodeMapLen(len(m))
	if w.h != nil && w.h.SortMaps {
		for _, v := range sortMapStringInt(m) {
			w.EncodeString(v)
			w.EncodeInt(m[v])
		}
	} else {
		for k, v := range m {
			w.EncodeString(k)
			w.EncodeInt(v)
		}
	}
}

func (w *Writer) encodeMapStringUint(m map[string]uint) {
	w.encodeMapLen(len(m))
	if w.h != nil && w.h.SortMaps {
		for _, v := range sortMapStringUint(m) {
			w.EncodeString(v)
			w.EncodeUint(m[v])
		}
	} else {
		for k, v := range m {
			w.EncodeString(k)
			w.EncodeUint(v)
		}
	}
}

func (w *Writer) encodeMapStringBool(m map[string]bool) {
	w.encodeMapLen(len(m))
	if w.h != nil && w.h.SortMaps {
		for _, v := range sortMapStringBool(m) {
			w.EncodeString(v)
			w.EncodeBool(m[v])
		}
	} else {
		for k, v := range m {
			w.EncodeString(k)
			w.EncodeBool(v)
		}
	}
}

func (w *Writer) encodeMapStringString(m map[string]string) {
	w.encodeMapLen(len(m))
	if w.h != nil && w.h.SortMaps {
		for _, v := range sortMapStringString(m) {
			w.EncodeString(v)
			w.EncodeString(m[v])
		}
	} else {
		for k, v := range m {
			w.EncodeString(k)
			w.EncodeString(v)
		}
	}
}

func (w *Writer) encodeMapIntAny(m map[int]interface{}) {
	w.encodeMapLen(len(m))
	if w.h != nil && w.h.SortMaps {
		for _, v := range sortMapIntAny(m) {
			w.EncodeInt(v)
			w.EncodeAny(m[v])
		}
	} else {
		for k, v := range m {
			w.EncodeInt(k)
			w.EncodeAny(v)
		}
	}
}

func (w *Writer) encodeMapUintAny(m map[uint]interface{}) {
	w.encodeMapLen(len(m))
	if w.h != nil && w.h.SortMaps {
		for _, v := range sortMapUintAny(m) {
			w.EncodeUint(v)
			w.EncodeAny(m[v])
		}
	} else {
		for k, v := range m {
			w.EncodeUint(k)
			w.EncodeAny(v)
		}
	}
}

func (w *Writer) encodeMapStringAny(m map[string]interface{}) {
	w.encodeMapLen(len(m))
	if w.h != nil && w.h.SortMaps {
		for _, v := range sortMapStringAny(m) {
			w.EncodeString(v)
			w.EncodeAny(m[v])
		}
	} else {
		for k, v := range m {
			w.EncodeString(k)
			w.EncodeAny(v)
		}
	}
}

func (w *Writer) encodeMapTimeAny(m map[time.Time]interface{}) {
	w.encodeMapLen(len(m))
	if w.h != nil && w.h.SortMaps {
		for _, v := range sortMapTimeAny(m) {
			w.EncodeTime(v)
			w.EncodeAny(m[v])
		}
	} else {
		for k, v := range m {
			w.EncodeTime(k)
			w.EncodeAny(v)
		}
	}
}

func (w *Writer) encodeMapAnyAny(m map[interface{}]interface{}) {
	w.encodeMapLen(len(m))
	if w.h != nil && w.h.SortMaps {
		for _, v := range sortMapAnyAny(m) {
			w.writeMany(v.key)
			w.EncodeAny(v.val)
		}
	} else {
		for k, v := range m {
			w.EncodeAny(k)
			w.EncodeAny(v)
		}
	}
}
