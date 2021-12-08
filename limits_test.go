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
	"fmt"
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLimits(t *testing.T) {

	var sizes = []uint{1, 127, 128, math.MaxUint8, math.MaxUint8 + 1, math.MaxUint16, math.MaxUint16 + 1}

	Convey("Can create larger types", t, func() {

		for _, l := range sizes {

			Convey(fmt.Sprintf("Can encode map of length %d", l), func() {

				var bit []byte
				var val = make(map[uint]interface{})
				var tmp = make(map[uint]interface{})

				for i := uint(0); i < l; i++ {
					val[i] = nil
				}

				enc := NewEncoderBytes(&bit)
				eer := enc.Encode(val)
				So(eer, ShouldBeNil)

				dec := NewDecoderBytes(bit)
				der := dec.Decode(&tmp)
				So(der, ShouldBeNil)

				So(tmp, ShouldResemble, val)

			})

			Convey(fmt.Sprintf("Can encode array of length %d", l), func() {

				var bit []byte
				var val = make([]interface{}, 0)
				var tmp = make([]interface{}, 0)

				for i := uint(0); i < l; i++ {
					val = append(val, nil)
				}

				enc := NewEncoderBytes(&bit)
				eer := enc.Encode(val)
				So(eer, ShouldBeNil)

				dec := NewDecoderBytes(bit)
				der := dec.Decode(&tmp)
				So(der, ShouldBeNil)

				So(tmp, ShouldResemble, val)

			})

		}

	})

}
