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

// DecodeSelfer decodes a cork.Selfer value from the Reader.
func (r *Reader) DecodeSelfer(v Selfer) {
	if r.readOne() == cSlf {
		if r.readOne() == v.ExtendCORK() {
			v.UnmarshalCORK(r)
			return
		}
	}
	panic(fail)
}

// DecodeCorker decodes a cork.Corker value from the Reader.
func (r *Reader) DecodeCorker(v Corker) {
	b := r.readOne()
	var s int
	switch {
	case b >= cFixExt && b <= cFixExt+fixedExt:
		s = int(b - cFixExt)
	case b == cExt8:
		s = int(r.readLen8())
	case b == cExt16:
		s = int(r.readLen16())
	case b == cExt32:
		s = int(r.readLen32())
	case b == cExt64:
		s = int(r.readLen64())
	default:
		panic(fail)
	}
	if r.readOne() != v.ExtendCORK() {
		panic(fail)
	}
	if err := v.UnmarshalCORK(r.readMany(s)); err != nil {
		panic(err)
	}
}
