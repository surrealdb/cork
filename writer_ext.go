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
	"math"
)

// EncodeSelfer encodes a cork.Selfer value to the Writer.
func (w *Writer) EncodeSelfer(v Selfer) {
	w.writeOne(cSlf)
	w.writeOne(v.ExtendCORK())
	v.MarshalCORK(w)
}

// EncodeCorker encodes a cork.Corker value to the Writer.
func (w *Writer) EncodeCorker(v Corker) {
	enc, err := v.MarshalCORK()
	if err != nil {
		panic(err)
	}
	sze := len(enc)
	switch {
	case sze <= fixedExt:
		w.writeOne(cFixExt + byte(sze))
	case sze <= math.MaxUint8:
		w.writeOne(cExt8)
		w.writeLen8(uint8(sze))
	case sze <= math.MaxUint16:
		w.writeOne(cExt16)
		w.writeLen16(uint16(sze))
	case sze <= math.MaxUint32:
		w.writeOne(cExt32)
		w.writeLen32(uint32(sze))
	case sze <= math.MaxInt64:
		w.writeOne(cExt64)
		w.writeLen64(uint64(sze))
	}
	w.writeOne(v.ExtendCORK())
	w.writeMany(enc)
}
