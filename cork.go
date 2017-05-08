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

const tag = "cork"

const (
	fixedInt = 1<<7 - 1 // 127
	fixedStr = 1<<5 - 1 // 31
	fixedBin = 1<<4 - 1 // 15
	fixedExt = 1<<4 - 1 // 15
	fixedArr = 1<<4 - 1 // 16
	fixedMap = 1<<4 - 1 // 16
)

const (
	cFixInt     byte = 0x00 // -> 0x7F = 128
	cFixStr          = 0x80 // -> 0x9F = 32
	cFixBin          = 0xA0 // -> 0xAf = 16
	cFixExt          = 0xB0 // -> 0xBF = 16
	cFixArr          = 0xC0 // -> 0xCF = 16
	cFixMap          = 0xD0 // -> 0xDF = 16
	_                = 0
	cNil             = 0xE0
	cTrue            = 0xE1
	cFalse           = 0xE2
	cTime            = 0xE3
	_                = 0
	cStr8            = 0xE4
	cStr16           = 0xE5
	cStr32           = 0xE6
	cStr64           = 0xE7
	_                = 0
	cBin8            = 0xE8
	cBin16           = 0xE9
	cBin32           = 0xEA
	cBin64           = 0xEB
	_                = 0
	cExt8            = 0xEC
	cExt16           = 0xED
	cExt32           = 0xEE
	cExt64           = 0xEF
	_                = 0
	cInt8            = 0xF0
	cInt16           = 0xF1
	cInt32           = 0xF2
	cInt64           = 0xF3
	_                = 0
	cUint8           = 0xF4
	cUint16          = 0xF5
	cUint32          = 0xF6
	cUint64          = 0xF7
	_                = 0
	cFloat32         = 0xF8
	cFloat64         = 0xF9
	_                = 0
	cComplex64       = 0xFA
	cComplex128      = 0xFB
	_                = 0
	cArr             = 0xFC
	cMap             = 0xFD
	cSym             = 0xFE
	cAlt             = 0xFF
	_                = 0
)

// Corker represents an object which can encode and decode itself.
type Corker interface {
	ExtendCORK() byte
	MarshalCORK() ([]byte, error)
	UnmarshalCORK([]byte) error
}
