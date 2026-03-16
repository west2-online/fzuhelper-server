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

// TMessageType represents message type constants in the Thrift protocol.
// originally from github.com/apache/thrift
type TMessageType = int32 // use alias for better flexibility of interfaces

const (
	INVALID_TMESSAGE_TYPE TMessageType = 0
	CALL                  TMessageType = 1
	REPLY                 TMessageType = 2
	EXCEPTION             TMessageType = 3
	ONEWAY                TMessageType = 4
)

// TType represents field type constants in the Thrift protocol
// originally from github.com/apache/thrift
type TType = int8 // use alias for better flexibility of interfaces

const (
	STOP   TType = 0
	VOID   TType = 1
	BOOL   TType = 2
	BYTE   TType = 3
	I08    TType = 3
	DOUBLE TType = 4
	I16    TType = 6
	I32    TType = 8
	I64    TType = 10
	STRING TType = 11
	UTF7   TType = 11
	STRUCT TType = 12
	MAP    TType = 13
	SET    TType = 14
	LIST   TType = 15
	UTF8   TType = 16
	UTF16  TType = 17
)

const defaultRecursionDepth = 64 // for skip

const ( // for Write/ReadMessage
	msgVersion1    = 0x80010000
	msgVersionMask = 0xffff0000
	msgTypeMask    = 0x0000ffff // for TMessageType
)
