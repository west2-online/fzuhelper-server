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
	"strconv"

	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/cache"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
)

func (s *LaunchScreenService) MobileGetImage(req *launch_screen.MobileGetImageRequest) (*[]db.Picture, int64, error, bool) {
	getFromMysql := false
	if !cache.IsLaunchScreenCacheExist(s.ctx, strconv.FormatInt(req.StudentId, 10)) {
		// no cache
		getFromMysql = true
	}

	if cache.IsLastLaunchScreenIdCacheExist(s.ctx) && !getFromMysql {
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
	if getFromMysql {
		imageModel := &db.Picture{
			SType: req.SType,
		}
		imgList, cnt, err := db.GetImageBySType(s.ctx, imageModel)
		if err != nil {
			return nil, -1, err, true
		}

		imgIdList := make([]int64, len(*imgList))
		for i, img := range *imgList {
			imgIdList[i] = img.ID
		}
		if err = cache.SetLaunchScreenCache(s.ctx, strconv.FormatInt(req.StudentId, 10), &imgIdList); err != nil {
			return nil, -1, err, true
		}

		id, err := db.GetLastImageId(s.ctx)
		if err != nil {
			return nil, -1, err, true
		}
		if err = cache.SetLastLaunchScreenIdCache(s.ctx, id); err != nil {
			return nil, -1, err, true
		}

		return imgList, cnt, nil, true
	}
	imgIdList, err := cache.GetLaunchScreenCache(s.ctx, strconv.FormatInt(req.StudentId, 10))
	if err != nil {
		return nil, -1, err, true
	}
	imgList, cnt, err := db.GetImageByIdList(s.ctx, imgIdList)
	if err != nil {
		return nil, -1, err, true
	}
	return imgList, cnt, nil, false
}
