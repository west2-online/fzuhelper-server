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

	"github.com/bytedance/gopkg/lang/dirtmake"
)

// nocopyWriteThreshold represents the threshold of using `NocopyWriter` for binary or string
//
// It's used by `WriteBinaryNocopy` and `WriteStringNocopy` of `BinaryProtocol`
// which are relied by kitex tool or thriftgo
const nocopyWriteThreshold = 4096

// BinaryWriter represents the method used in thrift encoding for nocopy writes
// It supports netpoll nocopy feature, see: https://github.com/cloudwego/netpoll/blob/develop/nocopy.go
type NocopyWriter interface {
	WriteDirect(b []byte, remainCap int) error
}

// FastCodec represents the interface of thrift fastcodec generated structs
type FastCodec interface {
	BLength() int
	FastWriteNocopy(buf []byte, bw NocopyWriter) int
	FastRead(buf []byte) (int, error)
}

// FastMarshal marshals the msg to buf. The msg should implement FastCodec.
func FastMarshal(msg FastCodec) []byte {
	sz := msg.BLength()
	buf := dirtmake.Bytes(sz, sz)
	msg.FastWriteNocopy(buf, nil)
	return buf
}

// FastUnmarshal unmarshal the buf into msg. The msg should implement FastCodec.
func FastUnmarshal(buf []byte, msg FastCodec) error {
	_, err := msg.FastRead(buf)
	return err
}

// MarshalFastMsg encodes the given msg to buf for generic thrift RPC.
func MarshalFastMsg(method string, msgType TMessageType, seq int32, msg FastCodec) ([]byte, error) {
	if method == "" {
		return nil, errors.New("method not set")
	}
	sz := Binary.MessageBeginLength(method) + msg.BLength()
	b := dirtmake.Bytes(sz, sz)
	i := Binary.WriteMessageBegin(b, method, msgType, seq)
	_ = msg.FastWriteNocopy(b[i:], nil)
	return b, nil
}

// UnmarshalFastMsg parses the given buf and stores the result to msg for generic thrift RPC.
// for EXCEPTION msgType, it will returns `err` with *ApplicationException type without storing the result to msg.
func UnmarshalFastMsg(b []byte, msg FastCodec) (method string, seq int32, err error) {
	method, msgType, seq, i, err := Binary.ReadMessageBegin(b)
	if err != nil {
		return "", 0, err
	}
	b = b[i:]

	if msgType == EXCEPTION {
		ex := NewApplicationException(UNKNOWN_APPLICATION_EXCEPTION, "")
		_, err = ex.FastRead(b)
		if err != nil {
			return method, seq, err
		}
		return method, seq, ex
	}
	_, err = msg.FastRead(b)
	return method, seq, err
}
