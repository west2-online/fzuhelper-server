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

package service

import (
	"fmt"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
)

func (s *PaperService) GetDir(req *paper.ListDirFilesRequest) (*model.UpYunFileDir, error) {
	key := s.cache.Paper.GetFileDirKey(req.Path)

	if s.cache.IsKeyExist(s.ctx, key) {
		fileDir, err := s.cache.Paper.GetFileDirCache(s.ctx, key)
		if err != nil {
			return nil, errno.Errorf(errno.InternalRedisErrorCode, "Paper.GetDir: get dir info failed: %v", err)
		}
		return fileDir, nil

	}

	fileDir, err := upyun.GetDir(req.Path)
	if err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "Paper.GetDir: get dir info failed: %v", err)
	}

	for i := len(fileDir.Folders) - 1; i >= 0; i-- {
		if constants.IgnoreUpyunDir[fileDir.Folders[i]] {
			fileDir.Folders = append(fileDir.Folders[:i], fileDir.Folders[i+1:]...)
		}
	}

	s.taskQueue.Add(fmt.Sprintf("SetFileDirCache:%s", key), taskqueue.QueueTask{Execute: func() error {
		return s.cache.Paper.SetFileDirCache(s.ctx, key, *fileDir)
	}})

	return fileDir, nil
}
