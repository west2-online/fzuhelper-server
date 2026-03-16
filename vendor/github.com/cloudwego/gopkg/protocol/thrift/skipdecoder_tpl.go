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
)

// skipDecoderIface represent the generics constraint of a SkipDecoder.
//
// It's used by skipDecoderImpl
type skipDecoderIface interface {
	// SkipN reads and skips n bytes
	//
	// SkipDecoderIface will not hold or modify the bytes between two `SkipN` calls.
	// It's safe to reuse buffer for next `SkipN` call.
	SkipN(n int) ([]byte, error)
}

// NOTE: At the time of writing Go generics doesn't fully support monomorphization, and
// it doesn't generate code copies for specific types which means
// inline calling of SkipN is not working ...
// This would be fixed in the future hopefully, so we use generics here.
func skipDecoderImpl[T skipDecoderIface](r T, t TType, maxdepth int) error {
	if maxdepth == 0 {
		return errDepthLimitExceeded
	}
	if sz := typeToSize[t]; sz > 0 {
		_, err := r.SkipN(int(sz))
		return err
	}
	switch t {
	case STRING:
		b, err := r.SkipN(4)
		if err != nil {
			return err
		}
		sz := int(binary.BigEndian.Uint32(b))
		if sz < 0 {
			return errDataLength
		}
		if _, err := r.SkipN(sz); err != nil {
			return err
		}
	case STRUCT:
		for {
			b, err := r.SkipN(1) // TType
			if err != nil {
				return err
			}
			tp := TType(b[0])
			if tp == STOP {
				break
			}
			if sz := typeToSize[tp]; sz > 0 {
				// fastpath
				// Field ID + Value
				if _, err := r.SkipN(2 + int(sz)); err != nil {
					return err
				}
				continue
			}

			// Field ID
			if _, err := r.SkipN(2); err != nil {
				return err
			}
			// Field Value
			if err := skipDecoderImpl(r, tp, maxdepth-1); err != nil {
				return err
			}
		}
	case MAP:
		b, err := r.SkipN(6) // 1 byte key TType, 1 byte value TType, 4 bytes Len
		if err != nil {
			return err
		}
		kt, vt, sz := TType(b[0]), TType(b[1]), int32(binary.BigEndian.Uint32(b[2:]))
		if sz < 0 {
			return errDataLength
		}
		ksz, vsz := int(typeToSize[kt]), int(typeToSize[vt])
		if ksz > 0 && vsz > 0 {
			_, err := r.SkipN(int(sz) * (ksz + vsz))
			return err
		}
		for i := int32(0); i < sz; i++ {
			if ksz > 0 {
				if _, err := r.SkipN(ksz); err != nil {
					return err
				}
			} else {
				if err := skipDecoderImpl(r, kt, maxdepth-1); err != nil {
					return err
				}
			}
			if vsz > 0 {
				if _, err := r.SkipN(vsz); err != nil {
					return err
				}
			} else {
				if err := skipDecoderImpl(r, vt, maxdepth-1); err != nil {
					return err
				}
			}
		}
	case SET, LIST:
		b, err := r.SkipN(5) // 1 byte value type, 4 bytes Len
		if err != nil {
			return err
		}
		vt, sz := TType(b[0]), int32(binary.BigEndian.Uint32(b[1:]))
		if sz < 0 {
			return errDataLength
		}
		if vsz := typeToSize[vt]; vsz > 0 {
			_, err := r.SkipN(int(sz) * int(vsz))
			return err
		}
		for i := int32(0); i < sz; i++ {
			if err := skipDecoderImpl(r, vt, maxdepth-1); err != nil {
				return err
			}
		}
	default:
		return NewProtocolException(INVALID_DATA, fmt.Sprintf("unknown data type %d", t))
	}
	return nil
}
