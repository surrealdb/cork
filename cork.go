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

const (
	fixedInt = 1<<7 - 1
	fixedBin = 1<<5 - 1
	fixedStr = 1<<5 - 1
)

const (
	cFixInt     byte = 0x00
	cFixBin          = 0x80
	cFixStr          = 0xA0
	_                = 0
	cNil             = 0xC0
	cTrue            = 0xC1
	cFalse           = 0xC2
	cTime            = 0xC3
	_                = 0
	cBin8            = 0xC4
	cBin16           = 0xC5
	cBin32           = 0xC6
	cBin64           = 0xC7
	_                = 0
	cStr8            = 0xC8
	cStr16           = 0xC9
	cStr32           = 0xCA
	cStr64           = 0xCB
	_                = 0
	cExt8            = 0xCC
	cExt16           = 0xCD
	cExt32           = 0xCE
	cExt64           = 0xCF
	_                = 0
	cInt8            = 0xD0
	cInt16           = 0xD1
	cInt32           = 0xD2
	cInt64           = 0xD3
	_                = 0
	cUint8           = 0xD4
	cUint16          = 0xD5
	cUint32          = 0xD6
	cUint64          = 0xD7
	_                = 0
	cFloat32         = 0xD8
	cFloat64         = 0xD9
	_                = 0
	cArr             = 0xDA
	cArrNil          = 0xDB
	cArrBool         = 0xDC
	cArrTime         = 0xDD
	cArrStr          = 0xDE
	cArrInt          = 0xDF
	cArrInt8         = 0xE0
	cArrInt16        = 0xE1
	cArrInt32        = 0xE2
	cArrInt64        = 0xE3
	cArrUint         = 0xE4
	cArrUint16       = 0xE5
	cArrUint32       = 0xE6
	cArrUint64       = 0xE7
	cArrFloat32      = 0xE8
	cArrFloat64      = 0xE9
	_                = 0
	cMap             = 0xF0
	cStruct          = 0xF1
	cMapStrNil       = 0xF2
	cMapStrBool      = 0xF3
	cMapStrStr       = 0xF4
	cMapStrInt       = 0xF5
	cMapNilNil       = 0xF6
	_                = 0xF7
	_                = 0xF8
	_                = 0xF9
	_                = 0xFA
	_                = 0xFB
	_                = 0xFC
	_                = 0xFD
	_                = 0xFE
	_                = 0xFF
	_                = 0
)

// Register adds a Corker type to the registry, enabling the
// object type to be encoded and decoded using the Corker methods.
func Register(value interface{}) {

}

// Corker represents an object which can encode and decode itself.
type Corker interface {
	ExtendCORK() byte
	MarshalCORK() ([]byte, error)
	UnmarshalCORK([]byte) error
}
