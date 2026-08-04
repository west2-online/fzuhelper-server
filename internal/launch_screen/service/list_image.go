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
	"math"

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	db "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

const (
	defaultLaunchScreenPageNum  = 1
	defaultLaunchScreenPageSize = 20
	maxLaunchScreenPageSize     = 100
)

func normalizeLaunchScreenListPage(pageNum, pageSize int64) (int, int, error) {
	if pageNum <= 0 {
		pageNum = defaultLaunchScreenPageNum
	}
	if pageSize <= 0 || pageSize > maxLaunchScreenPageSize {
		pageSize = defaultLaunchScreenPageSize
	}

	if pageNum-1 > math.MaxInt/pageSize {
		return 0, 0, fmt.Errorf("LaunchScreenService.ListImage error: page offset is too large")
	}

	return int(pageNum), int(pageSize), nil
}

// ListImage returns one page of admin-visible launch screen pictures.
func (s *LaunchScreenService) ListImage(req *launch_screen.ListImageRequest) (*[]db.Picture, int64, error) {
	if !utils.CheckPwd(req.Secret) {
		return nil, 0, fmt.Errorf("LaunchScreenService.ListImage error: AuthFailedError")
	}

	pageNum, pageSize, err := normalizeLaunchScreenListPage(req.GetPageNum(), req.GetPageSize())
	if err != nil {
		return nil, 0, err
	}

	pictures, total, err := s.db.LaunchScreen.ListImage(s.ctx, pageNum, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("LaunchScreenService.ListImage error:%w", err)
	}
	return pictures, total, nil
}
