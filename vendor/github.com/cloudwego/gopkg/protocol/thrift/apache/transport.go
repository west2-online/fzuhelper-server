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

package apache

import (
	"bytes"
	"context"
	"io"
	"unsafe"
)

// TTransport is identical with thrift.TTransport.
type TTransport interface {
	io.ReadWriteCloser
	RemainingBytes() (num_bytes uint64)
	Flush(ctx context.Context) (err error)
	Open() error
	IsOpen() bool
}

type defaultTransport struct {
	io.ReadWriter
}

// NewDefaultTransport converts io.ReadWriter to TTransport.
// Use NewBufferTransport if using *bytes.Buffer for better performance.
func NewDefaultTransport(rw io.ReadWriter) TTransport {
	if buf, ok := rw.(*bytes.Buffer); ok {
		return NewBufferTransport(buf)
	}
	return defaultTransport{rw}
}

// remoteByteBuffer represents remote.ByteBuffer in kitex
type remoteByteBuffer interface {
	ReadableLen() (n int)
}

func (p defaultTransport) IsOpen() bool                  { return true }
func (p defaultTransport) Open() error                   { return nil }
func (p defaultTransport) Close() error                  { return nil }
func (p defaultTransport) Flush(_ context.Context) error { return nil }

func (p defaultTransport) RemainingBytes() uint64 {
	if v, ok := p.ReadWriter.(remoteByteBuffer); ok {
		n := v.ReadableLen()
		if n > 0 {
			return uint64(n)
		}
	}
	return ^uint64(0)
}

type bufferTransport struct {
	bytes.Buffer
}

// NewBufferTransport extends bytes.Buffer to support TTransport
func NewBufferTransport(buf *bytes.Buffer) TTransport {
	// reuse buf's pointer with more methods
	return (*bufferTransport)(unsafe.Pointer(buf))
}

func (p *bufferTransport) IsOpen() bool                  { return true }
func (p *bufferTransport) Open() error                   { return nil }
func (p *bufferTransport) Close() error                  { p.Reset(); return nil }
func (p *bufferTransport) Flush(_ context.Context) error { return nil }
func (p *bufferTransport) RemainingBytes() uint64        { return uint64(p.Len()) }
