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
	"io"
	"sync"

	"github.com/bytedance/gopkg/lang/mcache"
	"github.com/cloudwego/gopkg/bufiox"
)

var poolSkipDecoder = sync.Pool{
	New: func() interface{} {
		return &SkipDecoder{}
	},
}

// SkipDecoder scans the underlying io.Reader and returns the bytes of a type
type SkipDecoder struct {
	r bufiox.Reader

	rn int
}

// NewSkipDecoder ...
//
// call Release if no longer use
func NewSkipDecoder(r bufiox.Reader) *SkipDecoder {
	p := poolSkipDecoder.Get().(*SkipDecoder)
	p.r = r
	return p
}

// Release puts SkipDecoder back to pool and reuse it next time.
//
// DO NOT USE SkipDecoder after calling Release.
func (p *SkipDecoder) Release() {
	*p = SkipDecoder{}
	poolSkipDecoder.Put(p)
}

// Next skips a specific type and returns its bytes.
//
// The returned buf is directly from bufiox.Reader with the same lifecycle.
func (p *SkipDecoder) Next(t TType) (buf []byte, err error) {
	p.rn = 0
	if err = skipDecoderImpl(p, t, defaultRecursionDepth); err != nil {
		return
	}
	buf, err = p.r.Next(p.rn)
	return
}

// SkipN implements SkipDecoderIface
func (p *SkipDecoder) SkipN(n int) (buf []byte, err error) {
	// old version netpoll may have performance issue when using Peek
	// see: https://github.com/cloudwego/netpoll/pull/335
	if buf, err = p.r.Peek(p.rn + n); err == nil {
		buf = buf[p.rn:]
		p.rn += n
	}
	return
}

// BytesSkipDecoder ...
type BytesSkipDecoder struct {
	n int
	b []byte
}

var poolBytesSkipDecoder = sync.Pool{
	New: func() interface{} {
		return &BytesSkipDecoder{}
	},
}

// NewBytesSkipDecoder ...
//
// call Release if no longer use
func NewBytesSkipDecoder(b []byte) *BytesSkipDecoder {
	p := poolBytesSkipDecoder.Get().(*BytesSkipDecoder)
	p.Reset(b)
	return p
}

// Release puts BytesSkipDecoder back to pool and reuse it next time.
//
// DO NOT USE BytesSkipDecoder after calling Release.
func (p *BytesSkipDecoder) Release() {
	p.Reset(nil)
	poolBytesSkipDecoder.Put(p)
}

// Reset ...
func (p *BytesSkipDecoder) Reset(b []byte) {
	p.n = 0
	p.b = b
}

// Next skips a specific type and returns its bytes.
//
// The returned buf refers to the input []byte without copy
func (p *BytesSkipDecoder) Next(t TType) (b []byte, err error) {
	if err = skipDecoderImpl(p, t, defaultRecursionDepth); err != nil {
		return
	}
	b = p.b[:p.n]
	p.b = p.b[p.n:]
	p.n = 0
	return
}

// SkipN implements skipDecoderIface
func (p *BytesSkipDecoder) SkipN(n int) ([]byte, error) {
	if len(p.b) >= p.n+n {
		p.n += n
		return p.b[p.n-n : p.n], nil
	}
	return nil, io.EOF
}

// ReaderSkipDecoder ...
type ReaderSkipDecoder struct {
	r io.Reader

	n int // bytes read, n <= len(b)
	b []byte
}

var poolReaderSkipDecoder = sync.Pool{
	New: func() interface{} {
		return &ReaderSkipDecoder{}
	},
}

// NewReaderSkipDecoder creates a ReaderSkipDecoder from pool
//
// call Release if no longer use
func NewReaderSkipDecoder(r io.Reader) *ReaderSkipDecoder {
	p := poolReaderSkipDecoder.Get().(*ReaderSkipDecoder)
	p.Reset(r)
	return p
}

// Release puts ReaderSkipDecoder back to pool and reuse it next time.
//
// DO NOT USE ReaderSkipDecoder after calling Release.
func (p *ReaderSkipDecoder) Release() {
	// no need to free p.b
	// will make use of p.b without reallcation
	p.Reset(nil)
	poolReaderSkipDecoder.Put(p)
}

// Reset ...
func (p *ReaderSkipDecoder) Reset(r io.Reader) {
	p.r = r
	p.n = 0
}

// Grow grows the underlying buffer to fit n bytes
func (p *ReaderSkipDecoder) Grow(n int) {
	if len(p.b)-p.n >= n {
		return
	}
	p.growSlow(n)
}

func (p *ReaderSkipDecoder) growSlow(n int) {
	// mcache will take care of the new cap of newb to be power of 2
	newb := mcache.Malloc((len(p.b) + n) | 255) // at least 255 bytes
	newb = newb[:cap(newb)]
	if p.n > 0 {
		copy(newb, p.b[:p.n])
	}
	if p.b != nil {
		mcache.Free(p.b)
	}
	p.b = newb
}

// Next skips a specific type and returns its bytes.
//
// The returned []byte is valid before the next `Next` call or `Release`
func (p *ReaderSkipDecoder) Next(t TType) (b []byte, err error) {
	p.n = 0
	if err = skipDecoderImpl(p, t, defaultRecursionDepth); err != nil {
		return
	}
	return p.b[:p.n], nil
}

// SkipN implements SkipDecoderIface
func (p *ReaderSkipDecoder) SkipN(n int) (buf []byte, err error) {
	p.Grow(n)
	buf = p.b[p.n : p.n+n]
	for i := 0; i < n && err == nil; { // io.ReadFull(buf)
		var nn int
		nn, err = p.r.Read(buf[i:])
		i += nn
	}
	if err != nil {
		return
	}
	p.n += n
	return
}
