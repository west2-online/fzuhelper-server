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
	"errors"
	"fmt"
)

const ( // ApplicationException codes from apache thrift
	UNKNOWN_APPLICATION_EXCEPTION  = 0
	UNKNOWN_METHOD                 = 1
	INVALID_MESSAGE_TYPE_EXCEPTION = 2
	WRONG_METHOD_NAME              = 3
	BAD_SEQUENCE_ID                = 4
	MISSING_RESULT                 = 5
	INTERNAL_ERROR                 = 6
	PROTOCOL_ERROR                 = 7
	INVALID_TRANSFORM              = 8
	INVALID_PROTOCOL               = 9
	UNSUPPORTED_CLIENT_TYPE        = 10
)

// ApplicationException is for replacing apache.TApplicationException
// it implements ThriftFastCodec interface.
type ApplicationException struct {
	t int32
	m string
}

// NewApplicationException creates an ApplicationException instance
func NewApplicationException(t int32, msg string) *ApplicationException {
	return &ApplicationException{t: t, m: msg}
}

// Msg ...
func (e *ApplicationException) Msg() string { return e.m }

// TypeID ... for kitex
func (e *ApplicationException) TypeID() int32 { return e.t }

// TypeId ... for apache ApplicationException compatibility
func (e *ApplicationException) TypeId() int32 { return e.t }

// BLength returns the len of encoded buffer.
func (e *ApplicationException) BLength() int {
	return Binary.FieldBeginLength() + Binary.StringLength(e.m) + // e.m
		Binary.FieldBeginLength() + Binary.I32Length() + // e.t
		Binary.FieldStopLength() // STOP
}

// FastRead ...
func (e *ApplicationException) FastRead(b []byte) (off int, err error) {
	for {
		tp, id, l, err := Binary.ReadFieldBegin(b[off:])
		if err != nil {
			return off, err
		}
		off += l
		if tp == STOP {
			break
		}
		switch {
		case id == 1 && tp == STRING: // Msg
			e.m, l, err = Binary.ReadString(b[off:])
		case id == 2 && tp == I32: // TypeID
			e.t, l, err = Binary.ReadI32(b[off:])
		default:
			l, err = Binary.Skip(b, tp)
		}
		if err != nil {
			return off, err
		}
		off += l
	}
	return off, nil
}

// FastWrite ...
func (e *ApplicationException) FastWrite(b []byte) (off int) {
	off += Binary.WriteFieldBegin(b[off:], STRING, 1)
	off += Binary.WriteString(b[off:], e.m)
	off += Binary.WriteFieldBegin(b[off:], I32, 2)
	off += Binary.WriteI32(b[off:], e.t)
	off += Binary.WriteByte(b[off:], STOP)
	return off
}

// FastWriteNocopy ...
func (e *ApplicationException) FastWriteNocopy(b []byte, _ NocopyWriter) int {
	return e.FastWrite(b)
}

// originally from github.com/apache/thrift@v0.13.0/lib/go/thrift/exception.go
var defaultApplicationExceptionMessage = map[int32]string{
	UNKNOWN_APPLICATION_EXCEPTION:  "unknown application exception",
	UNKNOWN_METHOD:                 "unknown method",
	INVALID_MESSAGE_TYPE_EXCEPTION: "invalid message type",
	WRONG_METHOD_NAME:              "wrong method name",
	BAD_SEQUENCE_ID:                "bad sequence ID",
	MISSING_RESULT:                 "missing result",
	INTERNAL_ERROR:                 "unknown internal error",
	PROTOCOL_ERROR:                 "unknown protocol error",
	INVALID_TRANSFORM:              "Invalid transform",
	INVALID_PROTOCOL:               "Invalid protocol",
	UNSUPPORTED_CLIENT_TYPE:        "Unsupported client type",
}

// Error ...
func (e *ApplicationException) Error() string {
	if e.m != "" {
		return e.m
	}
	if m, ok := defaultApplicationExceptionMessage[e.t]; ok {
		return m
	}
	return fmt.Sprintf("unknown exception type [%d]", e.t)
}

// String ...
func (e *ApplicationException) String() string {
	return fmt.Sprintf("ApplicationException(%d): %q", e.t, e.m)
}

// TransportException is for replacing apache.TransportException
// it implements ThriftFastCodec interface.
type TransportException struct {
	ApplicationException // same implementation ...
}

// NewTransportException ...
func NewTransportException(t int32, m string) *TransportException {
	ret := TransportException{}
	ret.t = t
	ret.m = m
	return &ret
}

// ProtocolException is for replacing apache.ProtocolException
// it implements ThriftFastCodec interface.
type ProtocolException struct {
	ApplicationException // same implementation ...

	err error
}

const ( // ProtocolException codes from apache thrift
	UNKNOWN_PROTOCOL_EXCEPTION = 0
	INVALID_DATA               = 1
	NEGATIVE_SIZE              = 2
	SIZE_LIMIT                 = 3
	BAD_VERSION                = 4
	NOT_IMPLEMENTED            = 5
	DEPTH_LIMIT                = 6
)

var (
	errBufferTooShort = NewProtocolException(INVALID_DATA, "buffer too short")
	errDataLength     = NewProtocolException(INVALID_DATA, "invalid data length")
)

// NewTransportExceptionWithType
func NewProtocolException(t int32, m string) *ProtocolException {
	ret := ProtocolException{}
	ret.t = t
	ret.m = m
	return &ret
}

// NewProtocolException ...
func NewProtocolExceptionWithErr(err error) *ProtocolException {
	e, ok := err.(*ProtocolException)
	if ok {
		return e
	}
	ret := NewProtocolException(UNKNOWN_PROTOCOL_EXCEPTION, err.Error())
	ret.err = err
	return ret
}

// Unwrap ... for errors pkg
func (e *ProtocolException) Unwrap() error { return e.err }

// Is ... for errors pkg
func (e *ProtocolException) Is(err error) bool {
	t, ok := err.(tException)
	if ok && t.TypeId() == e.t && t.Error() == e.m {
		return true
	}
	return errors.Is(e.err, err)
}

// Generic Thrift exception with TypeId method
type tException interface {
	Error() string
	TypeId() int32
}

// Prepends additional information to an error without losing the Thrift exception interface
func PrependError(prepend string, err error) error {
	if t, ok := err.(*TransportException); ok {
		return NewTransportException(t.TypeID(), prepend+t.Error())
	}
	if t, ok := err.(*ProtocolException); ok {
		return NewProtocolException(t.TypeID(), prepend+err.Error())
	}
	if t, ok := err.(*ApplicationException); ok {
		return NewApplicationException(t.TypeID(), prepend+t.Error())
	}
	if t, ok := err.(tException); ok { // apache thrift exception?
		return NewApplicationException(t.TypeId(), prepend+t.Error())
	}
	return errors.New(prepend + err.Error())
}
