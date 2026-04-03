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

	"github.com/west2-online/fzuhelper-server/internal/paper/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
)

// PaperServiceImpl implements the last service interface defined in the IDL.
type PaperServiceImpl struct {
	ClientSet *base.ClientSet
	taskQueue taskqueue.TaskQueue
}

func NewPaperService(clientSet *base.ClientSet, taskQueue taskqueue.TaskQueue) *PaperServiceImpl {
	return &PaperServiceImpl{
		ClientSet: clientSet,
		taskQueue: taskQueue,
	}
}

// ListDirFiles implements the PaperServiceImpl interface.
func (s *PaperServiceImpl) ListDirFiles(ctx context.Context, req *paper.ListDirFilesRequest) (resp *paper.ListDirFilesResponse, err error) {
	resp = new(paper.ListDirFilesResponse)

	fileDir, err := service.NewPaperService(ctx, s.ClientSet, s.taskQueue).GetDir(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Dir = fileDir
	return resp, nil
}

// GetDownloadUrl implements the PaperServiceImpl interface.
func (s *PaperServiceImpl) GetDownloadUrl(ctx context.Context, req *paper.GetDownloadUrlRequest) (resp *paper.GetDownloadUrlResponse, err error) {
	resp = new(paper.GetDownloadUrlResponse)
	url, err := service.NewPaperService(ctx, s.ClientSet, nil).GetDownloadUrl(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Url = url
	return resp, nil
}
