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
	"io"
)

type writer struct {
	io.Writer
}

func newWriter(dst io.Writer) *writer {
	return &writer{dst}
}

func (w *writer) WriteOne(val byte) {
	_, err := w.Write([]byte{val})
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

func (w *writer) WriteManys(val string) {
	_, err := w.Write([]byte(val))
	if err != nil {
		panic(err)
	}
	return
}
