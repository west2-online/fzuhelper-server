/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package base

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func TestBuildBaseResp(t *testing.T) {
	nilError := BuildBaseResp(nil)
	assert.Equal(t, int64(errno.SuccessCode), nilError.Code)
	assert.Equal(t, errno.Success.ErrorMsg, nilError.Msg)

	normalError := BuildBaseResp(fmt.Errorf("ok"))
	assert.Equal(t, int64(errno.InternalServiceErrorCode), normalError.Code)
	assert.Equal(t, "ok", normalError.Msg)

	errnoError := BuildBaseResp(errno.NewErrNo(200, "ok"))
	assert.Equal(t, int64(200), errnoError.Code)
	assert.Equal(t, "ok", errnoError.Msg)
}

func TestBuildSuccessResp(t *testing.T) {
	r := BuildSuccessResp()
	assert.Equal(t, int64(errno.SuccessCode), r.Code)
	assert.Equal(t, errno.Success.ErrorMsg, r.Msg)
}

func TestBuildRespAndLog(t *testing.T) {
	nilError := BuildBaseResp(nil)
	assert.Equal(t, int64(errno.SuccessCode), nilError.Code)
	assert.Equal(t, errno.Success.ErrorMsg, nilError.Msg)

	normalError := BuildBaseResp(fmt.Errorf("ok"))
	assert.Equal(t, int64(errno.InternalServiceErrorCode), normalError.Code)
	assert.Equal(t, "ok", normalError.Msg)

	errnoError := BuildBaseResp(errno.NewErrNo(200, "ok"))
	assert.Equal(t, int64(200), errnoError.Code)
	assert.Equal(t, "ok", errnoError.Msg)
}
