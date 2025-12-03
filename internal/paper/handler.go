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

package paper

import (
	"context"
	"errors"
	"fmt"

	"github.com/west2-online/fzuhelper-server/internal/paper/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
	"github.com/west2-online/fzuhelper-server/pkg/base"
)

// PaperServiceImpl implements the last service interface defined in the IDL.
type PaperServiceImpl struct {
	ClientSet *base.ClientSet
}

func NewPaperService(clientSet *base.ClientSet) *PaperServiceImpl {
	return &PaperServiceImpl{
		ClientSet: clientSet,
	}
}

// ListDirFiles implements the PaperServiceImpl interface.
func (s *PaperServiceImpl) ListDirFiles(ctx context.Context, req *paper.ListDirFilesRequest) (resp *paper.ListDirFilesResponse, err error) {
	resp = new(paper.ListDirFilesResponse)
	fileDir := new(model.UpYunFileDir) //nolint:ineffassign

	var success bool

	success, fileDir, err = service.NewPaperService(ctx, s.ClientSet).GetDir(req)
	if err != nil {
		base.LogError(fmt.Errorf("Paper.ListDirFiles: get dir info failed: %w", err))
	}
	if !success {
		resp.Base = base.BuildBaseResp(errors.New("Paper.ListDirFiles: failed to get files info"))
		return resp, nil
	}

	resp.Base = base.BuildSuccessResp()
	resp.Dir = fileDir
	return resp, err
}

// GetDownloadUrl implements the PaperServiceImpl interface.
func (s *PaperServiceImpl) GetDownloadUrl(ctx context.Context, req *paper.GetDownloadUrlRequest) (resp *paper.GetDownloadUrlResponse, err error) {
	resp = new(paper.GetDownloadUrlResponse)

	url, err := service.NewPaperService(ctx, s.ClientSet).GetDownloadUrl(req)
	if err != nil {
		resp.Base = base.BuildRespAndLog(fmt.Errorf("Paper.GetDownloadUrl: get download url failed: %w", err))
		return resp, nil
	}

	resp.Base = base.BuildSuccessResp()
	resp.Url = url
	return resp, err
}
