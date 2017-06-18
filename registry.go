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
)

var registry = make(map[byte]reflect.Type)

// Register adds a Corker type to the registry, enabling the
// object type to be encoded and decoded using the Corker methods.
func Register(value interface{}) {

	switch val := value.(type) {
	case Corker:
		registry[val.ExtendCORK()] = reflect.TypeOf(val).Elem()
	case Selfer:
		registry[val.ExtendCORK()] = reflect.TypeOf(val).Elem()
	}

}
