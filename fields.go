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
)

type field struct {
	omit bool
	indx []int
	name string
	show string
}

func (f *field) Name() string {
	if len(f.show) > 0 {
		return f.show
	}
	return f.name
}

func newField(kind reflect.StructField) *field {

	// Field is private
	if len(kind.PkgPath) > 0 {
		return nil
	}

	// Field not supported
	switch kind.Type.Kind() {
	case reflect.Chan:
		return nil
	case reflect.Func:
		return nil
	}

	// Retrieve the tag
	tag := kind.Tag.Get(tag)

	// No tag specified
	if len(tag) == 0 {
		return &field{
			omit: false,
			name: kind.Name,
			indx: kind.Index,
		}
	}

	// Field is ignored
	if tag == "-" {
		return nil
	}

	idx := strings.Index(tag, ",")

	// Field is renamed
	if idx < 0 {
		return &field{
			omit: false,
			name: kind.Name,
			show: tag,
			indx: kind.Index,
		}
	}

	// Field is renamed
	if idx > 0 {
		return &field{
			omit: tag[idx+1:] == "omitempty",
			name: kind.Name,
			show: tag[:idx],
			indx: kind.Index,
		}
	}

	// Immediate comma
	if idx == 0 {
		return &field{
			omit: tag[idx+1:] == "omitempty",
			name: kind.Name,
			indx: kind.Index,
		}
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
