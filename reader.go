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
	"bufio"
	"io"
	"reflect"
	"unsafe"
)

type reader struct {
	*bufio.Reader
}

func newReader(r io.Reader) *reader {
	return &reader{Reader: bufio.NewReader(r)}
}

func (r *reader) ReadOne() (val byte) {
	data, err := r.ReadByte()
	if err != nil {
		panic(err)
	}
	return data
}

func (r *reader) ReadMany(l int) (val []byte) {
	data := make([]byte, l)
	_, err := r.Read(data)
	if err != nil {
		panic(err)
	}
	return data
}

func (r *reader) ReadText(l int) (val string) {
	b := r.ReadMany(l)
	return *(*string)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&b))))
}
