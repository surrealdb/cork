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

type Handle struct {

	// ArrType specifies the type of slice to use when decoding
	// into a nil interface during schema-less decoding of a
	// slice in the stream.
	//
	// If not specified, we use []interface{}
	ArrType interface{}

	// MapType specifies the type of map to use when decoding
	// into a nil interface during schema-less decoding of a
	// map in the stream.
	//
	// If not specified, we use map[interface{}]interface{}
	MapType interface{}

	// Precision controls whether integers are encoded with
	// exact precision, or using only the required precision
	// necessary to encode the integer size.
	//
	// This is useful if you want to decode integers into a
	// nil interface, but need the same type as the encoded value.
	//
	// Enabling this is also useful if you want to use cork
	// for storage, and need to support in-place updating of
	// serialized data, without the need to move bytes around.
	Precision bool
}
