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
	"strings"
	"sync"
)

type field struct {
	name string
	show string
	item reflect.Value
}

var fields = sync.Pool{
	New: func() interface{} {
		return &field{}
	},
}

func (f *field) done() {
	fields.Put(f)
}

func (f *field) reset(name, show string, item reflect.Value) *field {
	f.name = name
	f.show = show
	f.item = item
	return f
}

func newField(kind reflect.StructField, item reflect.Value) *field {

	// Item is invalid
	if !item.IsValid() {
		return nil
	}

	// Item is private
	if kind.PkgPath != "" {
		return nil
	}

	opt := kind.Tag.Get(tag)

	if len(opt) == 0 {
		return fields.Get().(*field).reset(kind.Name, kind.Name, item)
	}

	// Item is ignored
	if opt == "-" {
		return nil
	}

	idx := strings.Index(opt, ",")

	// Item is renamed
	if idx < 0 {
		return fields.Get().(*field).reset(kind.Name, opt, item)
	}

	// Item is renamed
	if idx > 0 {
		if opt[idx+1:] == "omitempty" && isEmpty(item) {
			return nil
		}
		return fields.Get().(*field).reset(kind.Name, opt[:idx], item)
	}

	// Immediate comma
	if idx == 0 {
		if opt[idx+1:] == "omitempty" && isEmpty(item) {
			return nil
		}
		return fields.Get().(*field).reset(kind.Name, kind.Name, item)
	}

	return nil

}

func isEmpty(item reflect.Value) bool {
	switch item.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice, reflect.String:
		return item.Len() == 0
	case reflect.Bool:
		return !item.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return item.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return item.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return item.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return item.IsNil()
	}
	return false
}
