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
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
)

func (s *PaperService) GetDir(req *paper.ListDirFilesRequest) (bool, *model.UpYunFileDir, error) {
	var (
		success bool
		err     error
		fileDir *model.UpYunFileDir
	)

	key := s.cache.Paper.GetFileDirKey(req.Path)

	if ok := s.cache.IsKeyExist(s.ctx, key); ok {
		success, fileDir, err = s.cache.Paper.GetFileDirCache(s.ctx, key)
		if success {
			return true, fileDir, nil
		}

		if err != nil {
			return false, nil, fmt.Errorf("service.GetDir: get dir info failed: %w", err)
		}
	}

	fileDir, err = upyun.GetDir(req.Path)
	if err != nil {
		return false, nil, fmt.Errorf("service.GetDir: get dir info failed: %w", err)
	}

	if err = s.cache.Paper.SetFileDirCache(s.ctx, key, *fileDir); err != nil {
		return true, fileDir, fmt.Errorf("service.GetDir: set file dir cache failed: %w", err)
	}

	return true, fileDir, err
}
