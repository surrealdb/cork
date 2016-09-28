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
	name string
	show string
	opts []string
}

func parse(name string, tags string) *field {
	res := strings.Split(tags, ",")
	if res[0] == "" {
		res[0] = name
	}
	return &field{name, res[0], res[1:]}
}

func (t field) Name() string {
	return t.name
}

func (t field) Show() string {
	return t.show
}

func (t field) Opt(opt string) bool {
	for _, o := range t.opts {
		if o == opt {
			return true
		}
	}
	return false
}

func newField(kind reflect.StructField, item reflect.Value) *field {

	field := parse(kind.Name, kind.Tag.Get("cork"))

	if kind.PkgPath != "" {
		return nil
	}

	if field.Show() == "-" {
		return nil
	}

	if field.Opt("omitempty") && reflect.DeepEqual(item.Interface(), reflect.Zero(item.Type()).Interface()) {
		return nil
	}

	return field

}
