/*
 * Copyright 2024 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package thrift

import (
	"unsafe"
)

// p2i32, used by skipType which implements a fast skip with unsafe.Pointer without bounds check
func p2i32(p unsafe.Pointer) int32 {
	return int32(uint32(*(*byte)(unsafe.Add(p, 3))) |
		uint32(*(*byte)(unsafe.Add(p, 2)))<<8 |
		uint32(*(*byte)(unsafe.Add(p, 1)))<<16 |
		uint32(*(*byte)(p))<<24)
}
