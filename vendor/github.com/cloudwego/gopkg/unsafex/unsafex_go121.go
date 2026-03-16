//go:build go1.21

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

package unsafex

import "unsafe"

// XXX: this file is built >=go1.21 instead of go1.20 for fixing build issue in go1.20:
//
// unsafe.SliceData requires go1.20 or later (-lang was set to go1.18; check go.mod)
//
// see:
// 	https://github.com/golang/go/issues/59033
// 	https://github.com/golang/go/issues/58554

// BinaryToString converts []byte to string without copy
// Deprecated: use unsafe package instead if you're using >go1.20
func BinaryToString(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}

// StringToBinary converts string to []byte without copy
// Deprecated: use unsafe package instead if you're using >go1.20
func StringToBinary(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}
