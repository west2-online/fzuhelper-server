// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ttheader

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/cloudwego/gopkg/bufiox"
	"github.com/cloudwego/gopkg/unsafex"
)

// The byte count of 32 and 16 integer values.
const (
	Size32 = 4
	Size16 = 2
)

// Bytes2Uint32NoCheck ...
func Bytes2Uint32NoCheck(bytes []byte) uint32 {
	return binary.BigEndian.Uint32(bytes)
}

// Bytes2Uint16NoCheck ...
func Bytes2Uint16NoCheck(bytes []byte) uint16 {
	return binary.BigEndian.Uint16(bytes)
}

// Bytes2Uint8 ...
func Bytes2Uint8(bytes []byte, off int) (uint8, error) {
	if len(bytes)-off < 1 {
		return 0, io.EOF
	}
	return bytes[off], nil
}

// Bytes2Uint16 ...
func Bytes2Uint16(bytes []byte, off int) (uint16, error) {
	if len(bytes)-off < 2 {
		return 0, io.EOF
	}
	return binary.BigEndian.Uint16(bytes[off:]), nil
}

// ReadString2BLen ...
func ReadString2BLen(bytes []byte, off int) (string, int, error) {
	length, err := Bytes2Uint16(bytes, off)
	strLen := int(length)
	if err != nil {
		return "", 0, err
	}
	off += 2
	if len(bytes)-off < strLen {
		return "", 0, io.EOF
	}

	buf := bytes[off : off+strLen]
	return string(buf), int(length) + 2, nil
}

// WriteByte ...
func WriteByte(val byte, out bufiox.Writer) error {
	var buf []byte
	var err error
	if buf, err = out.Malloc(1); err != nil {
		return err
	}
	buf[0] = val
	return nil
}

// WriteUint32 ...
func WriteUint32(val uint32, out bufiox.Writer) error {
	var buf []byte
	var err error
	if buf, err = out.Malloc(Size32); err != nil {
		return err
	}
	binary.BigEndian.PutUint32(buf, val)
	return nil
}

// WriteUint16 ...
func WriteUint16(val uint16, out bufiox.Writer) error {
	var buf []byte
	var err error
	if buf, err = out.Malloc(Size16); err != nil {
		return err
	}
	binary.BigEndian.PutUint16(buf, val)
	return nil
}

// WriteString ...
func WriteString(val string, out bufiox.Writer) (int, error) {
	strLen := len(val)
	if err := WriteUint32(uint32(strLen), out); err != nil {
		return 0, err
	}
	n, err := out.WriteBinary(unsafex.StringToBinary(val))
	if err != nil {
		return 0, err
	}
	return n + 4, nil
}

// WriteString2BLen ...
func WriteString2BLen(val string, out bufiox.Writer) (int, error) {
	strLen := len(val)
	if strLen > maxHeaderStringSize {
		// printing first 100 bytes is enough to troubleshooting
		return 0, fmt.Errorf("string exceeded %dB max size (actual: %dB, preview: %q)", maxHeaderStringSize, strLen, val[:100]+"...")
	}
	if err := WriteUint16(uint16(strLen), out); err != nil {
		return 0, err
	}
	n, err := out.WriteBinary(unsafex.StringToBinary(val))
	if err != nil {
		return 0, err
	}
	return n + 2, nil
}

func IsStreaming(bytes []byte) bool {
	if len(bytes) < 8 {
		return false
	}
	return binary.BigEndian.Uint16(bytes[Size32:]) == uint16(TTHeaderMagic>>16) &&
		binary.BigEndian.Uint16(bytes[Size32+Size16:])&uint16(HeaderFlagsStreaming) != 0
}
