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

package handler

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"

	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/api/rpc"
	commonrpc "github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// TracePing .
// @router /api/v1/trace/ping [GET]
func TracePing(ctx context.Context, c *app.RequestContext) {
	// log with trace context
	logger.WithCtx(ctx).Info("HTTP trace ping request received")

	message, err := rpc.TracePingRPC(ctx, &commonrpc.TracePingRequest{})
	if err != nil {
		pack.RespError(c, err)
		return
	}

	pack.RespData(c, utils.H{
		"message": message,
	})
}
