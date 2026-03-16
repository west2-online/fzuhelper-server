// Copyright 2025 CloudWeGo Authors
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

package bufiox

import (
	"errors"
	"math/bits"

	"github.com/bytedance/gopkg/lang/dirtmake"
)

var (
	errNoRemainingData = errors.New("bufiox: no remaining data left")
)

var _ Reader = &BytesReader{}

type BytesReader struct {
	buf []byte // buf[ri:] is the buffer for reading.
	ri  int    // buf read positions
}

// NewBytesReader returns a new DefaultReader that reads from buf[:len(buf)].
// Its operation on buf is read-only.
func NewBytesReader(buf []byte) *BytesReader {
	r := &BytesReader{buf: buf}
	return r
}

func (r *BytesReader) Next(n int) (buf []byte, err error) {
	if n < 0 {
		err = errNegativeCount
		return
	}
	if n > len(r.buf)-r.ri {
		err = errNoRemainingData
		return
	}
	// nocopy read
	buf = r.buf[r.ri : r.ri+n]
	r.ri += n
	return
}

func (r *BytesReader) Peek(n int) (buf []byte, err error) {
	if n < 0 {
		err = errNegativeCount
		return
	}
	if n > len(r.buf)-r.ri {
		err = errNoRemainingData
		return
	}
	// nocopy read
	buf = r.buf[r.ri : r.ri+n]
	return
}

func (r *BytesReader) Skip(n int) (err error) {
	if n < 0 {
		err = errNegativeCount
		return
	}
	if n > len(r.buf)-r.ri {
		err = errNoRemainingData
		return
	}
	r.ri += n
	return
}

func (r *BytesReader) ReadLen() (n int) {
	return r.ri
}

func (r *BytesReader) ReadBinary(bs []byte) (n int, err error) {
	if len(bs) > len(r.buf)-r.ri {
		err = errNoRemainingData
		return
	}
	n = copy(bs, r.buf[r.ri:])
	r.ri += n
	return
}

func (r *BytesReader) Release(e error) error {
	r.buf = r.buf[r.ri:]
	r.ri = 0
	return nil
}

// BytesWriter implements Writer and builds a []byte result.
//
// It uses a deferred-copy scheme to avoid copying on buffer growth:
// when the buffer needs to grow, the old buffer is saved to oldBuf and
// a new buffer is allocated WITHOUT copying the old data. Slices returned
// by Malloc still point into the old buffer's backing array, so writes
// to them remain valid. At Flush time, data is reconstructed by copying
// each oldBuf entry's delta into the final buffer.
//
// BytesWriter can be flushed multiple times; each Flush outputs the
// accumulated data (including pre-existing data) and resets WrittenLen to 0.
type BytesWriter struct {
	wn     int      // bytes written since last Flush
	buf    []byte   // current write buffer; buf[:len(buf)] is the logical data
	oldBuf [][]byte // snapshots of buf before each grow, for deferred copy
	toBuf  *[]byte  // output destination, set by Flush
}

// NewBytesWriter returns a new BytesWriter that appends to buf[len(buf):cap(buf)].
// Existing data in buf[:len(buf)] is preserved.
func NewBytesWriter(buf *[]byte) *BytesWriter {
	w := &BytesWriter{toBuf: buf, buf: *buf}
	return w
}

func (w *BytesWriter) acquire(n int) {
	// fast path, for inline
	if len(w.buf)+n <= cap(w.buf) {
		return
	}
	w.acquireSlow(n)
}

func (w *BytesWriter) acquireSlow(n int) {
	need := len(w.buf) + n
	ncap := 1 << bits.Len(uint(need-1)) // smallest power of 2 >= need
	if ncap < defaultBufSize {
		ncap = defaultBufSize
	}
	// deltaLen is the number of new bytes in w.buf since the last snapshot.
	// If positive, w.buf has data that must be preserved via deferred copy.
	deltaLen := len(w.buf)
	if len(w.oldBuf) > 0 {
		deltaLen -= len(w.oldBuf[len(w.oldBuf)-1])
	}
	if deltaLen > 0 {
		w.oldBuf = append(w.oldBuf, w.buf)
	}
	// Allocate new buffer, set len to preserve the logical data length.
	// The region [0:len] is dirty and will be reconstructed from oldBuf at Flush.
	nbuf := dirtmake.Bytes(ncap, ncap)
	w.buf = nbuf[:len(w.buf)]
}

func (w *BytesWriter) Malloc(n int) (buf []byte, err error) {
	if n < 0 {
		err = errNegativeCount
		return
	}
	w.acquire(n)
	buf = w.buf[len(w.buf) : len(w.buf)+n]
	w.buf = w.buf[:len(w.buf)+n]
	w.wn += n
	return
}

func (w *BytesWriter) WriteBinary(bs []byte) (n int, err error) {
	w.acquire(len(bs))
	n = copy(w.buf[len(w.buf):cap(w.buf)], bs)
	w.buf = w.buf[:len(w.buf)+n]
	w.wn += n
	return
}

func (w *BytesWriter) WrittenLen() int {
	return w.wn
}

func (w *BytesWriter) Flush() (err error) {
	var offset int
	for _, old := range w.oldBuf {
		offset += copy(w.buf[offset:], old[offset:])
	}
	*w.toBuf = w.buf[:len(w.buf):len(w.buf)]
	w.oldBuf = nil
	w.wn = 0
	return nil
}
