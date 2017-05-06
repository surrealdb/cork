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

type writer struct {
	*bufio.Writer
}

func newWriter(w io.Writer) *writer {
	return &writer{Writer: bufio.NewWriter(w)}
}

func (w *writer) WriteOne(val byte) {
	err := w.WriteByte(val)
	if err != nil {
		panic(err)
	}
	return
}

func (w *writer) WriteMany(val []byte) {
	_, err := w.Write(val)
	if err != nil {
		panic(err)
	}
	return
}

func (w *writer) WriteText(val string) {
	b := *(*[]byte)(unsafe.Pointer((*reflect.SliceHeader)(unsafe.Pointer(&val))))
	w.WriteMany(b)
}
