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

type Handle struct {
	// SortMaps specifies whether maps should be sorted before
	// being encoded into CORK. This guarantees that the same
	// input data is always encoded into the same binary data.
	SortMaps bool
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
}
