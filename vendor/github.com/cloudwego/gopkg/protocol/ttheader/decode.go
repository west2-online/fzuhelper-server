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
	"context"
	"encoding/binary"
	"fmt"
	"io"

	"github.com/cloudwego/gopkg/bufiox"
)

const (
	// MagicMask is bit mask for checking version.
	MagicMask = 0xffff0000
)

// DecodeParam is used to return the ttheader info after decoding.
type DecodeParam struct {
	// Flags is used to set up header flags, default is 0.
	Flags HeaderFlags

	// SeqID is used to set up sequence id of a request/response.
	// MUST be unique for each request/response.
	SeqID int32

	// ProtocolID is used to set up protocol id of a request/response.
	// Default is ProtocolIDThriftBinary.
	ProtocolID ProtocolID

	// IntInfo is used to set up int key-value info into InfoIDIntKeyValue.
	// You can refer to metakey.go for more details.
	IntInfo map[uint16]string

	// StrInfo is used to set up string key-value info into InfoIDKeyValue.
	// You can refer to metakey.go for more details.
	StrInfo map[string]string

	// HeaderLen is used to set up header length of a request/response.
	HeaderLen int

	// PayloadLen is used to set up payload length of a request/response.
	PayloadLen int
}

// DecodeFromBytes decodes ttheader param from bytes.
func DecodeFromBytes(ctx context.Context, bs []byte) (param DecodeParam, err error) {
	in := bufiox.NewBytesReader(bs)
	param, err = Decode(ctx, in)
	_ = in.Release(nil)
	return
}

// Decode decodes ttheader param from bufiox.Reader.
func Decode(ctx context.Context, in bufiox.Reader) (param DecodeParam, err error) {
	var headerMeta []byte
	headerMeta, err = in.Next(TTHeaderMetaSize)
	if err != nil {
		return
	}
	if !IsTTHeader(headerMeta) {
		err = fmt.Errorf("not TTHeader protocol (first4Bytes=%#x, second4Bytes=%#x)", headerMeta[:4], headerMeta[4:8])
		return
	}
	totalLen := Bytes2Uint32NoCheck(headerMeta[:Size32])

	flags := Bytes2Uint16NoCheck(headerMeta[Size16*3:])
	param.Flags = HeaderFlags(flags)

	seqID := Bytes2Uint32NoCheck(headerMeta[Size32*2 : Size32*3])
	param.SeqID = int32(seqID)

	// avoid uint16 * 4 overflow, using uint32 to convert the result of Bytes2Uint16NoCheck
	headerInfoSize := uint32(Bytes2Uint16NoCheck(headerMeta[Size32*3:TTHeaderMetaSize])) * 4
	if headerInfoSize > MaxHeaderSize || headerInfoSize < 2 {
		err = fmt.Errorf("invalid header length[%d]", headerInfoSize)
		return
	}

	var headerInfo []byte
	if headerInfo, err = in.Next(int(headerInfoSize)); err != nil {
		return
	}
	param.ProtocolID = ProtocolID(headerInfo[0])
	if err = checkProtocolID(headerInfo[0]); err != nil {
		return
	}
	hdIdx := 2
	transformIDNum := int(headerInfo[1])
	if int(headerInfoSize)-hdIdx < transformIDNum {
		err = fmt.Errorf("need read %d transformIDs, but not enough", transformIDNum)
		return
	}
	transformIDs := make([]uint8, transformIDNum)
	for i := 0; i < transformIDNum; i++ {
		transformIDs[i] = headerInfo[hdIdx]
		hdIdx++
	}

	param.IntInfo, param.StrInfo, err = readKVInfo(hdIdx, headerInfo)
	if err != nil {
		err = fmt.Errorf("ttHeader read kv info failed, %s, headerInfo=%#x", err.Error(), headerInfo)
		return
	}

	param.HeaderLen = int(uint32(headerInfoSize) + TTHeaderMetaSize)
	param.PayloadLen = int(totalLen) + Size32 - param.HeaderLen
	return
}

/**
 * +------------------------------------------------------------+
 * |                  4Byte                 |       2Byte       |
 * +------------------------------------------------------------+
 * |   			     Length			    	|   HEADER MAGIC    |
 * +------------------------------------------------------------+
 */

func IsTTHeader(flagBuf []byte) bool {
	return binary.BigEndian.Uint32(flagBuf[Size32:])&MagicMask == TTHeaderMagic
}

func readKVInfo(idx int, buf []byte) (intKVMap map[uint16]string, strKVMap map[string]string, err error) {
	for {
		var infoID uint8
		infoID, err = Bytes2Uint8(buf, idx)
		idx++
		if err != nil {
			// this is the last field, read until there is no more padding
			if err == io.EOF {
				err = nil
			}
			return
		}
		switch InfoIDType(infoID) {
		case InfoIDPadding:
			continue
		case InfoIDKeyValue:
			if strKVMap == nil {
				strKVMap = make(map[string]string)
			}
			_, err = readStrKVInfo(&idx, buf, strKVMap)
			if err != nil {
				return
			}
		case InfoIDIntKeyValue:
			if intKVMap == nil {
				intKVMap = make(map[uint16]string)
			}
			_, err = readIntKVInfo(&idx, buf, intKVMap)
			if err != nil {
				return
			}
		case InfoIDACLToken:
			if strKVMap == nil {
				strKVMap = make(map[string]string)
			}
			if err = readACLToken(&idx, buf, strKVMap); err != nil {
				return
			}
		default:
			err = fmt.Errorf("invalid infoIDType[%#x]", infoID)
			return
		}
	}
}

func readIntKVInfo(idx *int, buf []byte, info map[uint16]string) (has bool, err error) {
	kvSize, err := Bytes2Uint16(buf, *idx)
	*idx += 2
	if err != nil {
		return false, fmt.Errorf("error reading int kv info size: %s", err.Error())
	}
	if kvSize <= 0 {
		return false, nil
	}
	for i := uint16(0); i < kvSize; i++ {
		key, err := Bytes2Uint16(buf, *idx)
		*idx += 2
		if err != nil {
			return false, fmt.Errorf("error reading int kv info: %s", err.Error())
		}
		val, n, err := ReadString2BLen(buf, *idx)
		*idx += n
		if err != nil {
			return false, fmt.Errorf("error reading int kv info: %s", err.Error())
		}
		info[key] = val
	}
	return true, nil
}

func readStrKVInfo(idx *int, buf []byte, info map[string]string) (has bool, err error) {
	kvSize, err := Bytes2Uint16(buf, *idx)
	*idx += 2
	if err != nil {
		return false, fmt.Errorf("error reading str kv info size: %s", err.Error())
	}
	if kvSize <= 0 {
		return false, nil
	}
	for i := uint16(0); i < kvSize; i++ {
		key, n, err := ReadString2BLen(buf, *idx)
		*idx += n
		if err != nil {
			return false, fmt.Errorf("error reading str kv info: %s", err.Error())
		}
		val, n, err := ReadString2BLen(buf, *idx)
		*idx += n
		if err != nil {
			return false, fmt.Errorf("error reading str kv info: %s", err.Error())
		}
		info[key] = val
	}
	return true, nil
}

// readACLToken reads acl token
func readACLToken(idx *int, buf []byte, info map[string]string) error {
	val, n, err := ReadString2BLen(buf, *idx)
	*idx += n
	if err != nil {
		return fmt.Errorf("error reading acl token: %s", err.Error())
	}
	info[GDPRToken] = val
	return nil
}

// protoID just for ttheader
func checkProtocolID(protoID uint8) error {
	switch protoID {
	case uint8(ProtocolIDThriftBinary):
	case uint8(ProtocolIDKitexProtobuf):
	case uint8(ProtocolIDThriftCompactV2):
		// just for compatibility
	case uint8(ProtocolIDThriftStruct):
	case uint8(ProtocolIDProtobufStruct):
	default:
		return fmt.Errorf("unsupported ProtocolID[%d]", protoID)
	}
	return nil
}
