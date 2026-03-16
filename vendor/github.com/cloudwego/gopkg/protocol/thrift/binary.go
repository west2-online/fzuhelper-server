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
	"encoding/binary"
	"fmt"
	"math"
	"unsafe"

	"github.com/bytedance/gopkg/lang/span"

	"github.com/cloudwego/gopkg/unsafex"
)

type spanCacheImpl interface {
	Copy(buf []byte) []byte
}

var (
	spanCache       spanCacheImpl // span.NewSpanCache does not return exported Type
	spanCacheEnable = false
)

// SetSpanCache enable/disable binary protocol bytes/string allocator
func SetSpanCache(enable bool) {
	spanCacheEnable = enable
	if enable && spanCache == nil {
		spanCache = span.NewSpanCache(1024 * 1024)
	}
}

var Binary BinaryProtocol

type BinaryProtocol struct{}

func (BinaryProtocol) WriteMessageBegin(buf []byte, name string, typeID TMessageType, seq int32) int {
	binary.BigEndian.PutUint32(buf, uint32(msgVersion1)|uint32(typeID&msgTypeMask))
	binary.BigEndian.PutUint32(buf[4:], uint32(len(name)))
	off := 8 + copy(buf[8:], name)
	binary.BigEndian.PutUint32(buf[off:], uint32(seq))
	return off + 4
}

func (BinaryProtocol) WriteFieldBegin(buf []byte, typeID TType, id int16) int {
	buf[0] = byte(typeID)
	binary.BigEndian.PutUint16(buf[1:], uint16(id))
	return 3
}

func (BinaryProtocol) WriteFieldStop(buf []byte) int {
	buf[0] = byte(STOP)
	return 1
}

func (BinaryProtocol) WriteMapBegin(buf []byte, kt, vt TType, size int) int {
	buf[0] = byte(kt)
	buf[1] = byte(vt)
	binary.BigEndian.PutUint32(buf[2:], uint32(size))
	return 6
}

func (BinaryProtocol) WriteListBegin(buf []byte, et TType, size int) int {
	buf[0] = byte(et)
	binary.BigEndian.PutUint32(buf[1:], uint32(size))
	return 5
}

func (BinaryProtocol) WriteSetBegin(buf []byte, et TType, size int) int {
	buf[0] = byte(et)
	binary.BigEndian.PutUint32(buf[1:], uint32(size))
	return 5
}

func (BinaryProtocol) WriteBool(buf []byte, v bool) int {
	if v {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return 1
}

func (BinaryProtocol) WriteByte(buf []byte, v int8) int {
	buf[0] = byte(v)
	return 1
}

func (BinaryProtocol) WriteI16(buf []byte, v int16) int {
	binary.BigEndian.PutUint16(buf, uint16(v))
	return 2
}

func (BinaryProtocol) WriteI32(buf []byte, v int32) int {
	binary.BigEndian.PutUint32(buf, uint32(v))
	return 4
}

func (BinaryProtocol) WriteI64(buf []byte, v int64) int {
	binary.BigEndian.PutUint64(buf, uint64(v))
	return 8
}

func (BinaryProtocol) WriteDouble(buf []byte, v float64) int {
	binary.BigEndian.PutUint64(buf, math.Float64bits(v))
	return 8
}

func (BinaryProtocol) WriteBinary(buf, v []byte) int {
	n := copy(buf[4:], v)
	binary.BigEndian.PutUint32(buf, uint32(n))
	return 4 + n
}

func (p BinaryProtocol) WriteBinaryNocopy(buf []byte, w NocopyWriter, v []byte) int {
	if w == nil || len(v) < nocopyWriteThreshold {
		return p.WriteBinary(buf, v)
	}
	binary.BigEndian.PutUint32(buf, uint32(len(v)))
	_ = w.WriteDirect(v, len(buf[4:])) // always err == nil ?
	return 4
}

func (BinaryProtocol) WriteString(buf []byte, v string) int {
	n := copy(buf[4:], v)
	binary.BigEndian.PutUint32(buf, uint32(n))
	return 4 + n
}

func (p BinaryProtocol) WriteStringNocopy(buf []byte, w NocopyWriter, v string) int {
	if w == nil || len(v) < nocopyWriteThreshold {
		return p.WriteString(buf, v)
	}
	binary.BigEndian.PutUint32(buf, uint32(len(v)))
	_ = w.WriteDirect(unsafex.StringToBinary(v), len(buf[4:])) // always err == nil ?
	return 4
}

// Append methods

func (p BinaryProtocol) AppendMessageBegin(buf []byte, name string, typeID TMessageType, seq int32) []byte {
	buf = appendUint32(buf, uint32(msgVersion1)|uint32(typeID&msgTypeMask))
	buf = p.AppendString(buf, name)
	return p.AppendI32(buf, seq)
}

func (BinaryProtocol) AppendFieldBegin(buf []byte, typeID TType, id int16) []byte {
	return append(buf, byte(typeID), byte(uint16(id>>8)), byte(id))
}

func (BinaryProtocol) AppendFieldStop(buf []byte) []byte {
	return append(buf, byte(STOP))
}

func (p BinaryProtocol) AppendMapBegin(buf []byte, kt, vt TType, size int) []byte {
	return p.AppendI32(append(buf, byte(kt), byte(vt)), int32(size))
}

func (p BinaryProtocol) AppendListBegin(buf []byte, et TType, size int) []byte {
	return p.AppendI32(append(buf, byte(et)), int32(size))
}

func (p BinaryProtocol) AppendSetBegin(buf []byte, et TType, size int) []byte {
	return p.AppendI32(append(buf, byte(et)), int32(size))
}

func (p BinaryProtocol) AppendBinary(buf, v []byte) []byte {
	return append(p.AppendI32(buf, int32(len(v))), v...)
}

func (p BinaryProtocol) AppendString(buf []byte, v string) []byte {
	return append(p.AppendI32(buf, int32(len(v))), v...)
}

func (BinaryProtocol) AppendBool(buf []byte, v bool) []byte {
	if v {
		return append(buf, 1)
	} else {
		return append(buf, 0)
	}
}

func (BinaryProtocol) AppendByte(buf []byte, v int8) []byte {
	return append(buf, byte(v))
}

func (BinaryProtocol) AppendI16(buf []byte, v int16) []byte {
	return append(buf, byte(uint16(v)>>8), byte(v))
}

func (BinaryProtocol) AppendI32(buf []byte, v int32) []byte {
	return appendUint32(buf, uint32(v))
}

func (BinaryProtocol) AppendI64(buf []byte, v int64) []byte {
	return appendUint64(buf, uint64(v))
}

func (BinaryProtocol) AppendDouble(buf []byte, v float64) []byte {
	return appendUint64(buf, math.Float64bits(v))
}

func appendUint32(buf []byte, v uint32) []byte {
	return append(buf, byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

func appendUint64(buf []byte, v uint64) []byte {
	return append(buf, byte(v>>56), byte(v>>48), byte(v>>40), byte(v>>32),
		byte(v>>24), byte(v>>16), byte(v>>8), byte(v))
}

// Length methods

func (BinaryProtocol) MessageBeginLength(method string) int {
	return 4 + (4 + len(method)) + 4
}

func (BinaryProtocol) FieldBeginLength() int           { return 3 }
func (BinaryProtocol) FieldStopLength() int            { return 1 }
func (BinaryProtocol) MapBeginLength() int             { return 6 }
func (BinaryProtocol) ListBeginLength() int            { return 5 }
func (BinaryProtocol) SetBeginLength() int             { return 5 }
func (BinaryProtocol) BoolLength() int                 { return 1 }
func (BinaryProtocol) ByteLength() int                 { return 1 }
func (BinaryProtocol) I16Length() int                  { return 2 }
func (BinaryProtocol) I32Length() int                  { return 4 }
func (BinaryProtocol) I64Length() int                  { return 8 }
func (BinaryProtocol) DoubleLength() int               { return 8 }
func (BinaryProtocol) StringLength(v string) int       { return 4 + len(v) }
func (BinaryProtocol) BinaryLength(v []byte) int       { return 4 + len(v) }
func (BinaryProtocol) StringLengthNocopy(v string) int { return 4 + len(v) }
func (BinaryProtocol) BinaryLengthNocopy(v []byte) int { return 4 + len(v) }

// Read methods

var (
	errReadMessage = NewProtocolException(INVALID_DATA, "ReadMessageBegin: buf too small")
	errBadVersion  = NewProtocolException(BAD_VERSION, "ReadMessageBegin: bad version")
)

func (p BinaryProtocol) ReadMessageBegin(buf []byte) (name string, typeID TMessageType, seq int32, l int, err error) {
	if len(buf) < 4 { // version+type header + name header
		return "", 0, 0, 0, errReadMessage
	}

	// read header for version and type
	header := binary.BigEndian.Uint32(buf)
	if header&msgVersionMask != msgVersion1 {
		return "", 0, 0, 0, errBadVersion
	}
	typeID = TMessageType(header & msgTypeMask)

	off := 4

	// read method name
	name, l, err1 := p.ReadString(buf[off:])
	if err1 != nil {
		return "", 0, 0, 0, errReadMessage
	}
	off += l

	// read seq
	seq, l, err2 := p.ReadI32(buf[off:])
	if err2 != nil {
		return "", 0, 0, 0, errReadMessage
	}
	off += l
	return name, typeID, seq, off, nil
}

var (
	errReadField = NewProtocolException(INVALID_DATA, "ReadFieldBegin: buf too small")
	errReadMap   = NewProtocolException(INVALID_DATA, "ReadMapBegin: buf too small")
	errReadList  = NewProtocolException(INVALID_DATA, "ReadListBegin: buf too small")
	errReadSet   = NewProtocolException(INVALID_DATA, "ReadSetBegin: buf too small")
	errReadStr   = NewProtocolException(INVALID_DATA, "ReadString: buf too small")
	errReadBin   = NewProtocolException(INVALID_DATA, "ReadBinary: buf too small")

	errReadBool   = NewProtocolException(INVALID_DATA, "ReadBool: len(buf) < 1")
	errReadByte   = NewProtocolException(INVALID_DATA, "ReadByte: len(buf) < 1")
	errReadI16    = NewProtocolException(INVALID_DATA, "ReadI16: len(buf) < 2")
	errReadI32    = NewProtocolException(INVALID_DATA, "ReadI32: len(buf) < 4")
	errReadI64    = NewProtocolException(INVALID_DATA, "ReadI64: len(buf) < 8")
	errReadDouble = NewProtocolException(INVALID_DATA, "ReadDouble: len(buf) < 8")
)

func (BinaryProtocol) ReadFieldBegin(buf []byte) (typeID TType, id int16, l int, err error) {
	if len(buf) < 1 {
		return 0, 0, 0, errReadField
	}
	typeID = TType(buf[0])
	if typeID == STOP {
		return STOP, 0, 1, nil
	}
	if len(buf) < 3 {
		return 0, 0, 0, errReadField
	}
	return typeID, int16(binary.BigEndian.Uint16(buf[1:])), 3, nil
}

func (BinaryProtocol) ReadMapBegin(buf []byte) (kt, vt TType, size, l int, err error) {
	if len(buf) < 6 {
		return 0, 0, 0, 0, errReadMap
	}
	l = 6
	kt, vt = TType(buf[0]), TType(buf[1])
	size = int(int32(binary.BigEndian.Uint32(buf[2:])))
	if size < 0 {
		err = errDataLength
	}
	return
}

func (BinaryProtocol) ReadListBegin(buf []byte) (et TType, size, l int, err error) {
	if len(buf) < 5 {
		return 0, 0, 0, errReadList
	}
	l = 5
	et = TType(buf[0])
	size = int(int32(binary.BigEndian.Uint32(buf[1:])))
	if size < 0 {
		err = errDataLength
	}
	return
}

func (BinaryProtocol) ReadSetBegin(buf []byte) (et TType, size, l int, err error) {
	if len(buf) < 5 {
		return 0, 0, 0, errReadSet
	}
	l = 5
	et = TType(buf[0])
	size = int(int32(binary.BigEndian.Uint32(buf[1:])))
	if size < 0 {
		err = errDataLength
	}
	return
}

func (p BinaryProtocol) ReadBinary(buf []byte) (b []byte, l int, err error) {
	sz, _, err := p.ReadI32(buf)
	if err != nil {
		return nil, 0, errReadBin
	}
	if sz < 0 {
		return nil, 0, errDataLength
	}
	l = 4 + int(sz)
	if len(buf) < l {
		return nil, 4, errReadBin
	}
	if spanCacheEnable {
		b = spanCache.Copy(buf[4:l])
	} else {
		b = []byte(string(buf[4:l]))
	}
	return b, l, nil
}

func (p BinaryProtocol) ReadString(buf []byte) (s string, l int, err error) {
	sz, _, err := p.ReadI32(buf)
	if err != nil {
		return "", 0, errReadStr
	}
	if sz < 0 {
		return "", 0, errDataLength
	}
	l = 4 + int(sz)
	if len(buf) < l {
		return "", 4, errReadStr
	}
	if spanCacheEnable {
		data := spanCache.Copy(buf[4:l])
		s = unsafex.BinaryToString(data)
	} else {
		s = string(buf[4:l])
	}
	return s, l, nil
}

func (BinaryProtocol) ReadBool(buf []byte) (v bool, l int, err error) {
	if len(buf) < 1 {
		return false, 0, errReadBool
	}
	if buf[0] == 1 {
		return true, 1, nil
	}
	return false, 1, nil
}

func (BinaryProtocol) ReadByte(buf []byte) (v int8, l int, err error) {
	if len(buf) < 1 {
		return 0, 0, errReadByte
	}
	return int8(buf[0]), 1, nil
}

func (BinaryProtocol) ReadI16(buf []byte) (v int16, l int, err error) {
	if len(buf) < 2 {
		return 0, 0, errReadI16
	}
	return int16(binary.BigEndian.Uint16(buf)), 2, nil
}

func (BinaryProtocol) ReadI32(buf []byte) (v int32, l int, err error) {
	if len(buf) < 4 {
		return 0, 0, errReadI32
	}
	return int32(binary.BigEndian.Uint32(buf)), 4, nil
}

func (BinaryProtocol) ReadI64(buf []byte) (v int64, l int, err error) {
	if len(buf) < 8 {
		return 0, 0, errReadI64
	}
	return int64(binary.BigEndian.Uint64(buf)), 8, nil
}

func (BinaryProtocol) ReadDouble(buf []byte) (v float64, l int, err error) {
	if len(buf) < 8 {
		return 0, 0, errReadDouble
	}
	return math.Float64frombits(binary.BigEndian.Uint64(buf)), 8, nil
}

var errDepthLimitExceeded = NewProtocolException(DEPTH_LIMIT, "depth limit exceeded")

var typeToSize = [256]int8{
	BOOL:   1,
	BYTE:   1,
	DOUBLE: 8,
	I16:    2,
	I32:    4,
	I64:    8,
}

func skipstr(p unsafe.Pointer, e uintptr) (int, error) {
	if uintptr(p)+uintptr(4) <= e {
		n := int(p2i32(p))
		if n < 0 {
			return 0, errDataLength
		}
		if uintptr(p)+uintptr(4+n) <= e {
			return 4 + n, nil
		}
	}
	return 0, errBufferTooShort
}

// Skip skips over the value for the given type using Go implementation.
func (BinaryProtocol) Skip(b []byte, t TType) (int, error) {
	if len(b) == 0 {
		return 0, errBufferTooShort
	}
	p := unsafe.Pointer(&b[0])
	e := uintptr(p) + uintptr(len(b))
	return skipType(p, e, t, defaultRecursionDepth)
}

func skipType(p unsafe.Pointer, e uintptr, t TType, maxdepth int) (int, error) {
	if maxdepth == 0 {
		return 0, errDepthLimitExceeded
	}
	if n := typeToSize[t]; n > 0 {
		if uintptr(p)+uintptr(n) > e {
			return 0, errBufferTooShort
		}
		return int(n), nil
	}
	var err error
	switch t {
	case STRING:
		return skipstr(p, e)
	case MAP:
		if uintptr(p)+uintptr(6) > e {
			return 0, errBufferTooShort
		}
		kt, vt, sz := TType(*(*byte)(p)), TType(*(*byte)(unsafe.Add(p, 1))), p2i32(unsafe.Add(p, 2))
		if sz < 0 {
			return 0, errDataLength
		}
		ksz, vsz := int(typeToSize[kt]), int(typeToSize[vt])
		if ksz > 0 && vsz > 0 { // fast path, fast skip
			mapkvsize := (int(sz) * (ksz + vsz))
			if uintptr(p)+uintptr(6+mapkvsize) > e {
				return 0, errBufferTooShort
			}
			return 6 + mapkvsize, nil
		}
		i := 6
		for j := int32(0); j < sz; j++ {
			if uintptr(p)+uintptr(i) >= e {
				return 0, errBufferTooShort
			}
			ki := 0
			if ksz > 0 {
				ki = ksz
			} else if kt == STRING {
				ki, err = skipstr(unsafe.Add(p, i), e)
			} else {
				ki, err = skipType(unsafe.Add(p, i), e, kt, maxdepth-1)
			}
			if err != nil {
				return i, err
			}
			i += ki
			if uintptr(p)+uintptr(i) >= e {
				return 0, errBufferTooShort
			}
			vi := 0
			if vsz > 0 {
				vi = vsz
			} else if vt == STRING {
				vi, err = skipstr(unsafe.Add(p, i), e)
			} else {
				vi, err = skipType(unsafe.Add(p, i), e, vt, maxdepth-1)
			}
			if err != nil {
				return i, err
			}
			i += vi
		}
		return i, nil
	case LIST, SET:
		if uintptr(p)+uintptr(5) > e {
			return 0, errBufferTooShort
		}
		vt, sz := TType(*(*byte)(p)), p2i32(unsafe.Add(p, 1))
		if sz < 0 {
			return 0, errDataLength
		}
		vsz := int(typeToSize[vt])
		if vsz > 0 { // fast path, fast skip
			listvsize := int(sz) * vsz
			if uintptr(p)+uintptr(5+listvsize) > e {
				return 0, errBufferTooShort
			}
			return 5 + listvsize, nil
		}
		i := 5
		for j := int32(0); j < sz; j++ {
			if uintptr(p)+uintptr(i) >= e {
				return 0, errBufferTooShort
			}
			vi := 0
			if vsz > 0 {
				vi = vsz
			} else if vt == STRING {
				vi, err = skipstr(unsafe.Add(p, i), e)
			} else {
				vi, err = skipType(unsafe.Add(p, i), e, vt, maxdepth-1)
			}
			if err != nil {
				return i, err
			}
			i += vi
		}
		return i, nil
	case STRUCT:
		i := 0
		for {
			if uintptr(p)+uintptr(i) >= e {
				return i, errBufferTooShort
			}
			ft := TType(*(*byte)(unsafe.Add(p, i)))
			i += 1 // TType
			if ft == STOP {
				return i, nil
			}
			i += 2 // Field ID
			if uintptr(p)+uintptr(i) >= e {
				return i, errBufferTooShort
			}
			fi := 0
			if typeToSize[ft] > 0 {
				fi = int(typeToSize[ft])
			} else if ft == STRING {
				fi, err = skipstr(unsafe.Add(p, i), e)
			} else {
				fi, err = skipType(unsafe.Add(p, i), e, ft, maxdepth-1)
			}
			if err != nil {
				return i, err
			}
			i += fi
		}
	default:
		return 0, NewProtocolException(INVALID_DATA, fmt.Sprintf("unknown data type %d", t))
	}
}
