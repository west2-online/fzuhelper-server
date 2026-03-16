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
	"sync"

	"github.com/bytedance/gopkg/lang/dirtmake"
	"github.com/cloudwego/gopkg/bufiox"
	"github.com/cloudwego/gopkg/unsafex"
)

// BufferReader represents a reader for binary protocol
type BufferReader struct {
	r bufiox.Reader
}

var poolBufferReader = sync.Pool{
	New: func() interface{} {
		return &BufferReader{}
	},
}

// NewBufferReader ... call Release if no longer use for reusing
func NewBufferReader(r bufiox.Reader) *BufferReader {
	ret := poolBufferReader.Get().(*BufferReader)
	ret.r = r
	return ret
}

// Recycle ...
func (r *BufferReader) Recycle() {
	r.r = nil
	poolBufferReader.Put(r)
}

func (r *BufferReader) next(n int) (b []byte, err error) {
	b, err = r.r.Next(n)
	if err != nil {
		err = NewProtocolExceptionWithErr(err)
	}
	return
}

func (r *BufferReader) readBinary(bs []byte) (n int, err error) {
	n, err = r.r.ReadBinary(bs)
	if err != nil {
		err = NewProtocolExceptionWithErr(err)
	}
	return
}

func (r *BufferReader) skipn(n int) (err error) {
	if n < 0 {
		return errDataLength
	}
	if err = r.r.Skip(n); err != nil {
		return NewProtocolExceptionWithErr(err)
	}
	return nil
}

// Readn returns total bytes read from underlying reader
func (r *BufferReader) Readn() int64 {
	return int64(r.r.ReadLen())
}

// ReadBool ...
func (r *BufferReader) ReadBool() (v bool, err error) {
	b, err := r.next(1)
	if err != nil {
		return false, err
	}
	v = b[0] == 1
	return
}

// ReadByte ...
func (r *BufferReader) ReadByte() (v int8, err error) {
	b, err := r.next(1)
	if err != nil {
		return 0, err
	}
	v = int8(b[0])
	return
}

// ReadI16 ...
func (r *BufferReader) ReadI16() (v int16, err error) {
	b, err := r.next(2)
	if err != nil {
		return 0, err
	}
	v = int16(binary.BigEndian.Uint16(b))
	return
}

// ReadI32 ...
func (r *BufferReader) ReadI32() (v int32, err error) {
	b, err := r.next(4)
	if err != nil {
		return 0, err
	}
	v = int32(binary.BigEndian.Uint32(b))
	return
}

// ReadI64 ...
func (r *BufferReader) ReadI64() (v int64, err error) {
	b, err := r.next(8)
	if err != nil {
		return 0, err
	}
	v = int64(binary.BigEndian.Uint64(b))
	return
}

// ReadDouble ...
func (r *BufferReader) ReadDouble() (v float64, err error) {
	b, err := r.next(8)
	if err != nil {
		return 0, err
	}
	v = math.Float64frombits(binary.BigEndian.Uint64(b))
	return
}

// ReadBinary ...
func (r *BufferReader) ReadBinary() (b []byte, err error) {
	sz, err := r.ReadI32()
	if err != nil {
		return nil, err
	}
	if sz < 0 {
		return nil, errDataLength
	}
	b = dirtmake.Bytes(int(sz), int(sz))
	_, err = r.readBinary(b)
	return
}

// ReadString ...
func (r *BufferReader) ReadString() (s string, err error) {
	b, err := r.ReadBinary()
	if err != nil {
		return "", err
	}
	return unsafex.BinaryToString(b), nil
}

// ReadMessageBegin ...
func (r *BufferReader) ReadMessageBegin() (name string, typeID TMessageType, seq int32, err error) {
	var header int32
	header, err = r.ReadI32()
	if err != nil {
		return
	}
	// read header for version and type
	if uint32(header)&msgVersionMask != msgVersion1 {
		err = errBadVersion
		return
	}
	typeID = TMessageType(uint32(header) & msgTypeMask)

	// read method name
	name, err = r.ReadString()
	if err != nil {
		return
	}

	// read seq
	seq, err = r.ReadI32()
	if err != nil {
		return
	}
	return
}

// ReadFieldBegin ...
func (r *BufferReader) ReadFieldBegin() (typeID TType, id int16, err error) {
	b, err := r.next(1)
	if err != nil {
		return 0, 0, err
	}
	typeID = TType(b[0])
	if typeID == STOP {
		return STOP, 0, nil
	}
	b, err = r.next(2)
	if err != nil {
		return 0, 0, err
	}
	id = int16(binary.BigEndian.Uint16(b))
	return
}

// ReadMapBegin ...
func (r *BufferReader) ReadMapBegin() (kt, vt TType, size int, err error) {
	b, err := r.next(6)
	if err != nil {
		return 0, 0, 0, err
	}
	kt, vt, size = TType(b[0]), TType(b[1]), int(binary.BigEndian.Uint32(b[2:]))
	return
}

// ReadListBegin ...
func (r *BufferReader) ReadListBegin() (et TType, size int, err error) {
	b, err := r.next(5)
	if err != nil {
		return 0, 0, err
	}
	et, size = TType(b[0]), int(binary.BigEndian.Uint32(b[1:]))
	return
}

// ReadSetBegin ...
func (r *BufferReader) ReadSetBegin() (et TType, size int, err error) {
	b, err := r.next(5)
	if err != nil {
		return 0, 0, err
	}
	et, size = TType(b[0]), int(binary.BigEndian.Uint32(b[1:]))
	return
}

// Skip ...
func (r *BufferReader) Skip(t TType) error {
	return r.skipType(t, defaultRecursionDepth)
}

func (r *BufferReader) skipstr() error {
	n, err := r.ReadI32()
	if err != nil {
		return err
	}
	return r.skipn(int(n))
}

func (r *BufferReader) skipType(t TType, maxdepth int) error {
	if maxdepth == 0 {
		return errDepthLimitExceeded
	}
	if n := typeToSize[t]; n > 0 {
		return r.skipn(int(n))
	}
	switch t {
	case STRING:
		return r.skipstr()
	case MAP:
		kt, vt, sz, err := r.ReadMapBegin()
		if err != nil {
			return err
		}
		if sz < 0 {
			return errDataLength
		}
		ksz, vsz := int(typeToSize[kt]), int(typeToSize[vt])
		if ksz > 0 && vsz > 0 {
			return r.skipn(sz * (ksz + vsz))
		}
		for j := 0; j < sz; j++ {
			if ksz > 0 {
				err = r.skipn(ksz)
			} else if kt == STRING {
				err = r.skipstr()
			} else {
				err = r.skipType(kt, maxdepth-1)
			}
			if err != nil {
				return err
			}
			if vsz > 0 {
				err = r.skipn(vsz)
			} else if vt == STRING {
				err = r.skipstr()
			} else {
				err = r.skipType(vt, maxdepth-1)
			}
			if err != nil {
				return err
			}
		}
		return nil
	case LIST, SET:
		vt, sz, err := r.ReadListBegin()
		if err != nil {
			return err
		}
		if sz < 0 {
			return errDataLength
		}
		if vsz := typeToSize[vt]; vsz > 0 {
			return r.skipn(sz * int(vsz))
		}
		for j := 0; j < sz; j++ {
			if vt == STRING {
				err = r.skipstr()
			} else {
				err = r.skipType(vt, maxdepth-1)
			}
			if err != nil {
				return err
			}
		}
		return nil
	case STRUCT:
		for {
			ft, _, err := r.ReadFieldBegin()
			if err != nil {
				return err
			}
			if ft == STOP {
				return nil
			}
			if fsz := typeToSize[ft]; fsz > 0 {
				err = r.skipn(int(fsz))
			} else {
				err = r.skipType(ft, maxdepth-1)
			}
			if err != nil {
				return err
			}
		}
	default:
		return NewProtocolException(INVALID_DATA, fmt.Sprintf("unknown data type %d", t))
	}
}
