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
	"math"

	"github.com/cloudwego/gopkg/bufiox"
)

/**
 *	TTHeader Protocol
 *	+-------------2Byte--------------|-------------2Byte-------------+
 *	+----------------------------------------------------------------+
 *	| 0|                          LENGTH                             |
 *	+----------------------------------------------------------------+
 *	| 0|       HEADER MAGIC          |            FLAGS              |
 *	+----------------------------------------------------------------+
 *	|                         SEQUENCE NUMBER                        |
 *	+----------------------------------------------------------------+
 *	| 0|     Header Size(/32)        | ...
 *	+---------------------------------
 *
 *	Header is of variable size:
 *	(and starts at offset 14)
 *
 *	+----------------------------------------------------------------+
 *	| PROTOCOL ID  |NUM TRANSFORMS . |TRANSFORM 0 ID (uint8)|
 *	+----------------------------------------------------------------+
 *	|  TRANSFORM 0 DATA ...
 *	+----------------------------------------------------------------+
 *	|         ...                              ...                   |
 *	+----------------------------------------------------------------+
 *	|        INFO 0 ID (uint8)      |       INFO 0  DATA ...
 *	+----------------------------------------------------------------+
 *	|         ...                              ...                   |
 *	+----------------------------------------------------------------+
 *	|                                                                |
 *	|                              PAYLOAD                           |
 *	|                                                                |
 *	+----------------------------------------------------------------+
 */

// Header keys
const (
	// Meta Size
	TTHeaderMetaSize = 14

	// Header Magics
	// 0 and 16th bits must be 0 to differentiate from framed & unframed
	TTHeaderMagic     uint32 = 0x10000000
	MeshHeaderMagic   uint32 = 0xFFAF0000
	MeshHeaderLenMask uint32 = 0x0000FFFF

	// HeaderMask        uint32 = 0xFFFF0000
	FlagsMask           uint32 = 0x0000FFFF
	MethodMask          uint32 = 0x41000000 // method first byte [A-Za-z_]
	MaxFrameSize        uint32 = 0x3FFFFFFF
	MaxHeaderSize       uint32 = 4 * math.MaxUint16
	maxHeaderStringSize int    = math.MaxUint16
)

type HeaderFlags uint16

const (
	HeaderFlagsStreaming        HeaderFlags = 0b0000_0000_0000_0010
	HeaderFlagSupportOutOfOrder HeaderFlags = 0x01
	HeaderFlagDuplexReverse     HeaderFlags = 0x08
	HeaderFlagSASL              HeaderFlags = 0x10
)

// ProtocolID is the wrapped protocol id used in THeader.
type ProtocolID uint8

// Supported ProtocolID values.
const (
	ProtocolIDThriftBinary    ProtocolID = 0x00
	ProtocolIDThriftCompact   ProtocolID = 0x02 // Kitex not support
	ProtocolIDThriftCompactV2 ProtocolID = 0x03 // Kitex not support
	ProtocolIDKitexProtobuf   ProtocolID = 0x04
	ProtocolIDThriftStruct    ProtocolID = 0x10 // TTHeader Streaming: only thrift struct encoded, no magic
	ProtocolIDProtobufStruct  ProtocolID = 0x11 // TTHeader Streaming: only protobuf struct encoded, no magic
	ProtocolIDDefault                    = ProtocolIDThriftBinary
)

type InfoIDType uint8 // uint8

const (
	InfoIDPadding     InfoIDType = 0
	InfoIDKeyValue    InfoIDType = 0x01
	InfoIDIntKeyValue InfoIDType = 0x10
	InfoIDACLToken    InfoIDType = 0x11
)

const (
	FrameTypeMeta    = "1"
	FrameTypeHeader  = "2"
	FrameTypeData    = "3"
	FrameTypeTrailer = "4"
	FrameTypeRst     = "5"
	FrameTypeInvalid = ""
)

// EncodeParam is used to set up params to encode ttheader.
type EncodeParam struct {
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
}

// EncodeToBytes encode ttheader to bytes.
// NOTICE: Must call
//
//	`binary.BigEndian.PutUint32(buf, uint32(totalLen))`
//
// after encoding both header and payload data to set total length of a request/response.
// And `totalLen` should be the length of header + payload - 4.
// You may refer to unit tests for examples.
func EncodeToBytes(ctx context.Context, param EncodeParam) (buf []byte, err error) {
	out := bufiox.NewBytesWriter(&buf)
	if _, err = Encode(ctx, param, out); err != nil {
		return
	}
	if err = out.Flush(); err != nil {
		return
	}
	return
}

// Encode encode ttheader to bufiox.Writer.
// NOTICE: Must call
//
//	`binary.BigEndian.PutUint32(totalLenField, uint32(totalLen))`
//
// after encoding both header and payload data to set total length of a request/response.
// And `totalLen` should be the length of header + payload - 4.
// You may refer to unit tests for examples.
func Encode(ctx context.Context, param EncodeParam, out bufiox.Writer) (totalLenField []byte, err error) {
	// 1. header meta
	var headerMeta []byte
	headerMeta, err = out.Malloc(TTHeaderMetaSize)
	if err != nil {
		return nil, fmt.Errorf("ttHeader malloc header meta failed, %s", err.Error())
	}

	totalLenField = headerMeta[0:4]
	headerInfoSizeField := headerMeta[12:14]
	binary.BigEndian.PutUint32(headerMeta[4:8], TTHeaderMagic+uint32(param.Flags))
	binary.BigEndian.PutUint32(headerMeta[8:12], uint32(param.SeqID))

	var transformIDs []uint8 // transformIDs not support TODO compress
	// 2.  header info, malloc and write
	if err = WriteByte(byte(param.ProtocolID), out); err != nil {
		return nil, fmt.Errorf("ttHeader write protocol id failed, %s", err.Error())
	}
	if err = WriteByte(byte(len(transformIDs)), out); err != nil {
		return nil, fmt.Errorf("ttHeader write transformIDs length failed, %s", err.Error())
	}
	for tid := range transformIDs {
		if err = WriteByte(byte(tid), out); err != nil {
			return nil, fmt.Errorf("ttHeader write transformIDs failed, %s", err.Error())
		}
	}
	// PROTOCOL ID(u8) + NUM TRANSFORMS(always 0)(u8) + TRANSFORM IDs([]u8)
	headerInfoSize := 1 + 1 + len(transformIDs)
	headerInfoSize, err = writeKVInfo(headerInfoSize, param.IntInfo, param.StrInfo, out)
	if err != nil {
		return nil, fmt.Errorf("ttHeader write kv info failed, %s", err.Error())
	}

	if uint32(headerInfoSize) > MaxHeaderSize {
		return nil, fmt.Errorf("invalid header length[%d]", headerInfoSize)
	}
	binary.BigEndian.PutUint16(headerInfoSizeField, uint16(headerInfoSize/4))
	return totalLenField, nil
}

func writeKVInfo(writtenSize int, intKVMap map[uint16]string, strKVMap map[string]string, out bufiox.Writer) (writeSize int, err error) {
	writeSize = writtenSize
	// str kv info
	strKVSize := len(strKVMap)
	// write gdpr token into InfoIDACLToken
	// supplementary doc: https://www.cloudwego.io/docs/kitex/reference/transport_protocol_ttheader/
	if gdprToken, ok := strKVMap[GDPRToken]; ok {
		strKVSize--
		// INFO ID TYPE(u8)
		if err = WriteByte(byte(InfoIDACLToken), out); err != nil {
			return writeSize, err
		}
		writeSize += 1

		wLen, err := WriteString2BLen(gdprToken, out)
		if err != nil {
			return writeSize, err
		}
		writeSize += wLen
	}

	if strKVSize > 0 {
		// INFO ID TYPE(u8) + NUM HEADERS(u16)
		if err = WriteByte(byte(InfoIDKeyValue), out); err != nil {
			return writeSize, err
		}
		if err = WriteUint16(uint16(strKVSize), out); err != nil {
			return writeSize, err
		}
		writeSize += 3
		for key, val := range strKVMap {
			if key == GDPRToken {
				continue
			}
			keyWLen, err := WriteString2BLen(key, out)
			if err != nil {
				return writeSize, err
			}
			valWLen, err := WriteString2BLen(val, out)
			if err != nil {
				return writeSize, err
			}
			writeSize = writeSize + keyWLen + valWLen
		}
	}

	// int kv info
	intKVSize := len(intKVMap)
	if intKVSize > 0 {
		// INFO ID TYPE(u8) + NUM HEADERS(u16)
		if err = WriteByte(byte(InfoIDIntKeyValue), out); err != nil {
			return writeSize, err
		}
		if err = WriteUint16(uint16(intKVSize), out); err != nil {
			return writeSize, err
		}
		writeSize += 3
		for key, val := range intKVMap {
			if err = WriteUint16(key, out); err != nil {
				return writeSize, err
			}
			valWLen, err := WriteString2BLen(val, out)
			if err != nil {
				return writeSize, err
			}
			writeSize = writeSize + 2 + valWLen
		}
	}

	// padding = (4 - headerInfoSize%4) % 4
	padding := (4 - writeSize%4) % 4
	paddingBuf, err := out.Malloc(padding)
	if err != nil {
		return writeSize, err
	}
	for i := 0; i < len(paddingBuf); i++ {
		paddingBuf[i] = byte(0)
	}
	writeSize += padding
	return
}
