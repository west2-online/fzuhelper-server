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
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
)

func (s *LaunchScreenService) GetImageById(id int64) (*db.Picture, error) {
	img, err := db.GetImageById(s.ctx, id)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// GetImagesByStuId 不使用
func (s *LaunchScreenService) GetImagesByStuId(uid int64) (*[]db.Picture, int64, error) {
	imgList, cnt, err := db.ListImageByUid(s.ctx, 1, uid)
	if err != nil {
		return nil, 0, err
	}
	return imgList, cnt, nil
}