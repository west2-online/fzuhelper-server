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

// Package apache contains code for working with apache thrift indirectly
//
// It acts as a bridge between generated code which relies on apache codec like:
//
//	Write(p thrift.TProtocol) error
//	Read(p thrift.TProtocol) error
//
// and kitex ecosystem.
//
// Because we're deprecating apache thrift, all kitex ecosystem code will not rely on apache thrift
// except one pkg: `github.com/cloudwego/kitex/pkg/protocol/bthrift`. Why is the package chosen?
// All legacy generated code relies on it, and we may not be able to update the code in a brief timeframe.
// So the package is chosen to register helper methods to this package in order to use it
// without importing `github.com/apache/thrift`
//
// Users must call `RegisterCheckTStruct`, `RegisterThriftRead` and `RegisterThriftWrite` for
// using `CheckTStruct`, `ThriftRead`, and `ThriftWrite`
//
// see README.md of `bthrift` above for more details
package apache

import (
	"errors"

	"github.com/cloudwego/gopkg/bufiox"
)

var (
	fnCheckTStruct func(v interface{}) error

	fnThriftRead  func(r bufiox.Reader, v interface{}) error
	fnThriftWrite func(w bufiox.Writer, v interface{}) error
)

// RegisterCheckTStruct accepts `thrift.TStruct check` func and save it for later use.
func RegisterCheckTStruct(fn func(v interface{}) error) {
	fnCheckTStruct = fn
}

// RegisterThriftRead ...
func RegisterThriftRead(fn func(r bufiox.Reader, v interface{}) error) {
	fnThriftRead = fn
}

// RegisterThriftWrite ...
func RegisterThriftWrite(fn func(w bufiox.Writer, v interface{}) error) {
	fnThriftWrite = fn
}

var (
	errCheckTStructNotRegistered = errors.New("func `RegisterCheckTStruct` not called")
	errThriftReadNotRegistered   = errors.New("func `RegisterThriftRead` not called")
	errThriftWriteNotRegistered  = errors.New("func `RegisterThriftWrite` not called")
)

// CheckTStruct ...
func CheckTStruct(v interface{}) error {
	if fnCheckTStruct == nil {
		return errCheckTStructNotRegistered
	}
	return fnCheckTStruct(v)
}

// ThriftRead ...
func ThriftRead(r bufiox.Reader, v interface{}) error {
	if fnThriftRead == nil {
		return errThriftReadNotRegistered
	}
	return fnThriftRead(r, v)
}

// ThriftWrite ...
func ThriftWrite(w bufiox.Writer, v interface{}) error {
	if fnThriftWrite == nil {
		return errThriftWriteNotRegistered
	}
	return fnThriftWrite(w, v)
}
