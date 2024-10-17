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

package rpc

import (
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
	"github.com/west2-online/fzuhelper-server/pkg/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitPaperRPC() {
	client, err := client.InitPaperRPC()
	if err != nil {
		logger.Fatalf("api.rpc.Paper InitPaperRPC failed, err  %v", err)
	}
	paperClient = *client
}

func GetDownloadUrlRPC(ctx context.Context, req *paper.GetDownloadUrlRequest) (url string, err error) {
	resp, err := paperClient.GetDownloadUrl(ctx, req)
	if err != nil {
		logger.Errorf("GetDownloadUrlRPC: RPC called failed: %v", err.Error())
		return "", errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return "", errno.BizError.WithMessage(resp.Base.Msg)
	}
	return resp.Url, nil
}

func GetDirFilesRPC(ctx context.Context, req *paper.ListDirFilesRequest) (files *model.UpYunFileDir, err error) {
	resp, err := paperClient.ListDirFiles(ctx, req)
	if err != nil {
		logger.Errorf("GetDirFilesRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return nil, errno.BizError.WithMessage(resp.Base.Msg)
	}
	return resp.Dir, nil
}

func UploadFileRPC(ctx context.Context, req *paper.UploadFileRequest) (err error) {
	resp, err := paperClient.UploadFile(ctx, req)
	if err != nil {
		logger.Errorf("UploadFileRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return errno.BizError.WithMessage(resp.Base.Msg)
	}

	return nil
}
