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

// Code generated by hertz generator.

package api

import (
	"context"
	"errors"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/fzuhelper-server/api/model/api"
	"github.com/west2-online/fzuhelper-server/api/pack"
	"github.com/west2-online/fzuhelper-server/api/rpc"

	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

// ListDirFiles .
// @router /api/v1/paper/list [GET]
func ListDirFiles(ctx context.Context, c *app.RequestContext) {
	var err error

	path := c.DefaultQuery("path", "/") // 默认根目录
	if path == "" {
		logger.Errorf("api.ListDirFiles: path is empty")
		pack.RespError(c, errno.ParamError.WithError(errors.New("path is empty")))
		return
	}

	res, err := rpc.GetDirFilesRPC(ctx, &paper.ListDirFilesRequest{
		Path: path,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}

	resp := new(api.ListDirFilesResponse)
	resp.Dir = pack.BuildUpYunFileDir(res)
	pack.RespData(c, resp.Dir)
}

// GetDownloadUrl .
// @router /api/v1/paper/download [GET]
func GetDownloadUrl(ctx context.Context, c *app.RequestContext) {
	var err error

	filepath := c.DefaultQuery("filepath", "/C语言/10份练习.zip")

	url, err := rpc.GetDownloadUrlRPC(ctx, &paper.GetDownloadUrlRequest{
		Filepath: filepath,
	})
	if err != nil {
		pack.RespError(c, err)
		return
	}

	resp := new(api.GetDownloadUrlResponse)
	resp.URL = url

	pack.RespData(c, resp)
}
