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
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/cache"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *LaunchScreenService) MobileGetImage(req *launch_screen.MobileGetImageRequest) (*[]db.Picture, int64, error, bool) {
	getFromMysql := false
	if !cache.IsLaunchScreenCacheExist(s.ctx, utils.GenerateRedisKeyByStuId(req.StudentId, req.SType)) {
		// no cache
		getFromMysql = true
	}

	if !getFromMysql {
		if cache.IsLastLaunchScreenIdCacheExist(s.ctx) {
			id, err := db.GetLastImageId(s.ctx)
			if err != nil {
				return nil, -1, err, true
			}
			cacheId, err := cache.GetLastLaunchScreenIdCache(s.ctx)
			if err != nil {
				return nil, -1, err, true
			}
			if cacheId != id {
				// expire
				getFromMysql = true
			}
		}
	}

	if getFromMysql {
		imageModel := &db.Picture{
			SType: req.SType,
		}
		imgList, cnt, err := db.GetImageBySType(s.ctx, imageModel)
		if err != nil {
			return nil, -1, err, true
		}
		if cnt == 0 {
			return nil, 0, nil, true
		}

		return imgList, cnt, nil, true
	}

	imgIdList, err := cache.GetLaunchScreenCache(s.ctx, utils.GenerateRedisKeyByStuId(req.StudentId, req.SType))
	if err != nil {
		return nil, -1, err, true
	}
	imgList, cnt, err := db.GetImageByIdList(s.ctx, &imgIdList)
	if err != nil {
		return nil, -1, err, true
	}
	return imgList, cnt, nil, false
}

func (s *LaunchScreenService) SetCache(imgList *[]db.Picture, req *launch_screen.MobileGetImageRequest) error {
	imgIdList := make([]int64, len(*imgList))

	for i, img := range *imgList {
		imgIdList[i] = img.ID
	}
	if err := cache.SetLaunchScreenCache(s.ctx, utils.GenerateRedisKeyByStuId(req.StudentId, req.SType), &imgIdList); err != nil {
		return err
	}

	id, err := db.GetLastImageId(s.ctx)
	if err != nil {
		return err
	}
	if err = cache.SetLastLaunchScreenIdCache(s.ctx, id); err != nil {
		return err
	}
	return nil
}

func (s *LaunchScreenService) AddShowTimes(pictureList *[]db.Picture) error {
	return db.AddImageListShowTime(s.ctx, pictureList)
}
