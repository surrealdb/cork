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
	"bytes"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDecode(t *testing.T) {

	var src = "Hello"

	var val = interface{}("Hello")

	var dst = []byte{cFixStr + 0x05, 72, 101, 108, 108, 111}

	Convey("Can use Decode", t, func() {
		out := Decode(dst)
		So(out, ShouldResemble, val)
	})

	Convey("Can use DecodeInto", t, func() {
		var out string
		DecodeInto(dst, &out)
		So(out, ShouldResemble, src)
	})

	Convey("Can use NewDecoder", t, func() {
		var out string
		var buf = bytes.NewReader(dst)
		NewDecoder(buf).Decode(&out)
		So(out, ShouldResemble, src)
	})

	Convey("Can use NewDecoderBytes", t, func() {
		var out string
		NewDecoderBytes(dst).Decode(&out)
		So(out, ShouldResemble, src)
	})

	Convey("Can use NewDecoderFromPool", t, func() {
		var out string
		var buf = bytes.NewReader(dst)
		NewDecoderFromPool(buf).Decode(&out)
		So(out, ShouldResemble, src)
	})

	Convey("Can use NewDecoderBytesFromPool", t, func() {
		var out string
		NewDecoderBytesFromPool(dst).Decode(&out)
		So(out, ShouldResemble, src)
	})

}
