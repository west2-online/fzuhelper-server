//go:build arm64 || amd64

// Copyright 2024 CloudWeGo Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtimex

import (
	"fmt"
	"runtime"
	"unsafe"

	"github.com/modern-go/reflect2"
)

var (
	// Runtime internals variables
	gType    reflect2.StructType
	gidField reflect2.StructField

	ErrNotSupported = fmt.Errorf("%s is not supported", runtime.Version())
)

func init() {
	gType, _ = reflect2.TypeByName("runtime.g").(reflect2.StructType)
	if gType != nil {
		gidField = gType.FieldByName("goid")
	}
}

func getg() unsafe.Pointer

//go:noescape
//go:linkname runtime_procPin runtime.procPin
func runtime_procPin() int

//go:noescape
//go:linkname runtime_procUnpin runtime.procUnpin
func runtime_procUnpin()

// GID return current goroutine's ID
func GID() (int, error) {
	if gidField == nil {
		return 0, ErrNotSupported
	}
	gp := getg()
	gid := *(*int64)(unsafe.Add(gp, gidField.Offset()))
	return int(gid), nil
}

// PID return current processor's ID
// It may not 100% accurate because the scheduler may preempt current goroutine and re-schedule it to another P after the function call return,
// so the caller should tolerate the preemption.
func PID() (int, error) {
	id := runtime_procPin()
	runtime_procUnpin()
	return id, nil
}
