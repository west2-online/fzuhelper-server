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

package pack

import (
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	api "github.com/west2-online/fzuhelper-server/api/model/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func BuildUpYunFileDir(res *model.UpYunFileDir) *api.UpYunFileDir {
	return &api.UpYunFileDir{
		BasePath: res.BasePath,
		Folders:  res.Folders,
		Files:    res.Files,
	}
}

type RespWithDataInPaper struct {
	Code int    `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

func RespDataInPaper(c *app.RequestContext, data any) {
	c.JSON(consts.StatusOK, RespWithDataInPaper{
		Code: errno.SuccessCodePaper,
		Msg:  "Success",
		Data: data,
	})
}

func RespErrorInPaper(c *app.RequestContext, err error) {
	Errno := errno.ConvertErr(err)
	c.JSON(consts.StatusOK, RespWithDataInPaper{
		Code: int(Errno.ErrorCode),
		Msg:  Errno.ErrorMsg,
		Data: nil,
	})
}
