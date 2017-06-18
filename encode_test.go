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

func TestEncode(t *testing.T) {

	var src = "Hello"

	var dst = []byte{cFixStr + 0x05, 72, 101, 108, 108, 111}

	Convey("Can use Encode", t, func() {
		out := Encode(src)
		So(out, ShouldResemble, dst)
	})

	Convey("Can use EncodeInto", t, func() {
		var buf []byte
		EncodeInto(src, &buf)
		So(buf, ShouldResemble, dst)
	})

	Convey("Can use NewEncoder", t, func() {
		var buf = bytes.NewBuffer(nil)
		NewEncoder(buf).Encode(src)
		So(buf.Bytes(), ShouldResemble, dst)
	})

	Convey("Can use NewEncoderBytes", t, func() {
		var buf []byte
		NewEncoderBytes(&buf).Encode(src)
		So(buf, ShouldResemble, dst)
	})

	Convey("Can use NewEncoderFromPool", t, func() {
		var buf = bytes.NewBuffer(nil)
		NewEncoderFromPool(buf).Encode(src)
		So(buf.Bytes(), ShouldResemble, dst)
	})

	Convey("Can use NewEncoderBytesFromPool", t, func() {
		var buf []byte
		NewEncoderBytesFromPool(&buf).Encode(src)
		So(buf, ShouldResemble, dst)
	})

}
