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
	"math"
	"sync"

	"github.com/cloudwego/gopkg/bufiox"
	"github.com/cloudwego/gopkg/unsafex"
)

type BufferWriter struct {
	w bufiox.Writer
}

var poolBufferWriter = sync.Pool{
	New: func() interface{} {
		return &BufferWriter{}
	},
}

func NewBufferWriter(iw bufiox.Writer) *BufferWriter {
	w := poolBufferWriter.Get().(*BufferWriter)
	w.w = iw
	return w
}

func (w *BufferWriter) Recycle() {
	w.w = nil
	poolBufferWriter.Put(w)
}

func (w *BufferWriter) WriteMessageBegin(name string, typeID TMessageType, seq int32) error {
	buf, err := w.w.Malloc(Binary.MessageBeginLength(name))
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(buf, uint32(msgVersion1)|uint32(typeID&msgTypeMask))
	binary.BigEndian.PutUint32(buf[4:], uint32(len(name)))
	copy(buf[8:], name)
	binary.BigEndian.PutUint32(buf[8+len(name):], uint32(seq))
	return nil
}

func (w *BufferWriter) WriteFieldBegin(typeID TType, id int16) error {
	buf, err := w.w.Malloc(3)
	if err != nil {
		return err
	}
	buf[0], buf[1], buf[2] = byte(typeID), byte(uint16(id>>8)), byte(id)
	return nil
}

func (w *BufferWriter) WriteFieldStop() error {
	buf, err := w.w.Malloc(1)
	if err != nil {
		return err
	}
	buf[0] = byte(STOP)
	return nil
}

func (w *BufferWriter) WriteMapBegin(kt, vt TType, size int) error {
	buf, err := w.w.Malloc(6)
	if err != nil {
		return err
	}
	buf[0], buf[1] = byte(kt), byte(vt)
	binary.BigEndian.PutUint32(buf[2:], uint32(size))
	return nil
}

func (w *BufferWriter) WriteListBegin(et TType, size int) error {
	buf, err := w.w.Malloc(5)
	if err != nil {
		return err
	}
	buf[0] = byte(et)
	binary.BigEndian.PutUint32(buf[1:], uint32(size))
	return nil
}

func (w *BufferWriter) WriteSetBegin(et TType, size int) error {
	buf, err := w.w.Malloc(5)
	if err != nil {
		return err
	}
	buf[0] = byte(et)
	binary.BigEndian.PutUint32(buf[1:], uint32(size))
	return nil
}

func (w *BufferWriter) WriteBinary(v []byte) error {
	buf, err := w.w.Malloc(4)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(buf, uint32(len(v)))
	_, err = w.w.WriteBinary(v)
	return err
}

func (w *BufferWriter) WriteString(v string) error {
	return w.WriteBinary(unsafex.StringToBinary(v))
}

func (w *BufferWriter) WriteBool(v bool) error {
	buf, err := w.w.Malloc(1)
	if err != nil {
		return err
	}
	if v {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return nil
}

func (w *BufferWriter) WriteByte(v int8) error {
	buf, err := w.w.Malloc(1)
	if err != nil {
		return err
	}
	buf[0] = byte(v)
	return nil
}

func (w *BufferWriter) WriteI16(v int16) error {
	buf, err := w.w.Malloc(2)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint16(buf, uint16(v))
	return nil
}

func (w *BufferWriter) WriteI32(v int32) error {
	buf, err := w.w.Malloc(4)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint32(buf, uint32(v))
	return nil
}

func (w *BufferWriter) WriteI64(v int64) error {
	buf, err := w.w.Malloc(8)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint64(buf, uint64(v))
	return nil
}

func (w *BufferWriter) WriteDouble(v float64) error {
	buf, err := w.w.Malloc(8)
	if err != nil {
		return err
	}
	binary.BigEndian.PutUint64(buf, math.Float64bits(v))
	return nil
}
