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

func isVal(b byte) bool {
	return b == cTrue || b == cFalse
}

func isArr(b byte) bool {
	return b >= cArr && b <= cArr+0x0F
}

func isMap(b byte) bool {
	return b >= cMap && b <= cMap+0x0F
}

func isNum(b byte) bool {
	return b >= cFixInt && b <= cFixInt+0x7F
}

func isExt(b byte) bool {
	return b == cExt8 || b == cExt16 || b == cExt32 || b == cExt64
}

func isBin(b byte) bool {
	return b == cBin8 || b == cBin16 || b == cBin32 || b == cBin64 || (b >= cFixBin && b <= cFixBin+0x1F)
}

func isStr(b byte) bool {
	return b == cStr8 || b == cStr16 || b == cStr32 || b == cStr64 || (b >= cFixStr && b <= cFixStr+0x1F)
}

func isInt(b byte) bool {
	return b == cInt8 || b == cInt16 || b == cInt32 || b == cInt64 || isNum(b)
}

func isUint(b byte) bool {
	return b == cUint8 || b == cUint16 || b == cUint32 || b == cUint64 || isNum(b)
}

func isFloat(b byte) bool {
	return b == cFloat32 || b == cFloat64
}

func isTime(b byte) bool {
	return b == cTime
}
