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
	"regexp"
	"strconv"

	"github.com/bytedance/sonic"
	"golang.org/x/sync/errgroup"

	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/cache"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *LaunchScreenService) MobileGetImage(req *launch_screen.MobileGetImageRequest) (respList *[]db.Picture, cntResp int64, err error) {
	getFromMysql := false
	if !cache.IsLaunchScreenCacheExist(s.ctx, utils.GenerateRedisKeyByStuId(req.StudentId, req.SType)) {
		// no cache exist
		getFromMysql = true
	}

	if !getFromMysql {
		if cache.IsLastLaunchScreenIdCacheExist(s.ctx) {
			id, err := db.GetLastImageId(s.ctx)
			if err != nil {
				return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage db.GetLastImageId error:%v", err.Error())
			}
			cacheId, err := cache.GetLastLaunchScreenIdCache(s.ctx)
			if err != nil {
				return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage cache.GetLastLaunchScreenIdCache error:%v", err.Error())
			}
			// 当最新存入图片id与缓存中的不一致时，需要重新获取
			if cacheId != id {
				getFromMysql = true
			}
		}
	}

	if getFromMysql {
		// 获取符合当前时间的imgList
		imgList, cnt, err := db.GetImageBySType(s.ctx, req.SType)
		if err != nil {
			return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage db.GetImageBySType error:%v", err.Error())
		}
		// 没有符合条件的图片
		if cnt == 0 {
			return nil, 0, errno.NoRunningPictureError
		}

		// 从中筛选出JSON含当前学号的图片(regex)
		cntResp = 0
		currentImgList := make([]db.Picture, 0)
		for _, picture := range *imgList {
			// 处理JSON
			m := make(map[string]string)
			if err = sonic.Unmarshal([]byte(picture.Regex), &m); err != nil {
				return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage unmarshal JSON error:%v", err.Error())
			}

			match := false
			for k, v := range m {
				if k == "picture_id" || k == "device" {
					continue
				}
				// 此时key是student_id,value是学号JSON
				match, err = regexp.MatchString(strconv.Itoa(int(req.StudentId)), v)
				if err != nil {
					continue
				}
				if !match {
					continue
				}
			}
			if match {
				cntResp++
				currentImgList = append(currentImgList, picture)
			}
		}

		if cntResp != 0 {
			var eg errgroup.Group
			eg.Go(func() error {
				// addShowTime
				return db.AddImageListShowTime(s.ctx, imgList)
			})
			eg.Go(func() error {
				// setIdCache
				imgIdList := make([]int64, len(*imgList))

				for i, img := range *imgList {
					imgIdList[i] = img.ID
				}
				return cache.SetLaunchScreenCache(s.ctx, utils.GenerateRedisKeyByStuId(req.StudentId, req.SType), &imgIdList)
			})
			eg.Go(func() error {
				// setExpireCheckCache
				id, err := db.GetLastImageId(s.ctx)
				if err != nil {
					return err
				}
				if err = cache.SetLastLaunchScreenIdCache(s.ctx, id); err != nil {
					return err
				}
				return nil
			})
			if err = eg.Wait(); err != nil {
				return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage set cache error:%v", err.Error())
			}
		} else {
			return nil, 0, errno.NoRunningPictureError
		}

		return &currentImgList, cntResp, nil
	}

	// 直接从缓存中获取id
	imgIdList, err := cache.GetLaunchScreenCache(s.ctx, utils.GenerateRedisKeyByStuId(req.StudentId, req.SType))
	if err != nil {
		return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage cache.GetLaunchScreenCache error:%v", err.Error())
	}
	respList, cntResp, err = db.GetImageByIdList(s.ctx, &imgIdList)
	if err != nil {
		return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage db.GetImageByIdList error:%v", err.Error())
	}
	// addShowtime
	if err = db.AddImageListShowTime(s.ctx, respList); err != nil {
		return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage db.AddImageListShowTime error:%v", err.Error())
	}

	return respList, cntResp, nil
}
