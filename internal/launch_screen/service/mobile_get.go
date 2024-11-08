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
	"errors"
	"fmt"
	"regexp"

	"github.com/bytedance/sonic"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *LaunchScreenService) MobileGetImage(req *launch_screen.MobileGetImageRequest) (respList *[]model.Picture, cntResp int64, err error) {
	getFromMysql, err := s.shouldGetFromMySQL(req.StudentId, req.SType)
	if err != nil {
		return nil, 0, err
	}

	if !getFromMysql {
		// 直接从缓存中获取id
		imgIdList, err := s.cache.LaunchScreen.GetLaunchScreenCache(s.ctx, utils.GenerateRedisKeyByStuId(req.StudentId, req.SType))
		if err != nil {
			return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage cache.GetLaunchScreenCache error:%w", err)
		}
		respList, cntResp, err = s.db.LaunchScreen.GetImageByIdList(s.ctx, &imgIdList)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				getFromMysql = true
			} else {
				return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage db.GetImageByIdList error:%w", err)
			}
		}
	}

	if getFromMysql {
		return s.getImagesFromMySQL(req.StudentId, req.SType)
	}

	// addShowtime for cache
	if err = s.db.LaunchScreen.AddImageListShowTime(s.ctx, respList); err != nil {
		return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage db.AddImageListShowTime error:%w", err)
	}

	return respList, cntResp, nil
}

// shouldGetFromMySQL 确认缓存是否过期或不存在
func (s *LaunchScreenService) shouldGetFromMySQL(studentId string, sType int64) (bool, error) {
	if !s.cache.IsKeyExist(s.ctx, utils.GenerateRedisKeyByStuId(studentId, sType)) {
		return true, nil
	}

	if s.cache.LaunchScreen.IsLastLaunchScreenIdCacheExist(s.ctx) {
		id, err := s.db.LaunchScreen.GetLastImageId(s.ctx)
		if err != nil {
			return true, fmt.Errorf("LaunchScreenService.MobileGetImage db.GetLastImageId error:%w", err)
		}
		cacheId, err := s.cache.LaunchScreen.GetLastLaunchScreenIdCache(s.ctx)
		if err != nil {
			return true, fmt.Errorf("LaunchScreenService.MobileGetImage cache.GetLastLaunchScreenIdCache error:%w", err)
		}
		// 当最新存入图片id与缓存中的不一致时，需要重新获取
		if cacheId != id {
			return true, nil
		}
	} else {
		return true, nil
	}

	return false, nil
}

// getImagesFromMySQL 从数据库中获取图片url
func (s *LaunchScreenService) getImagesFromMySQL(studentId string, sType int64) (*[]model.Picture, int64, error) {
	// 获取符合当前时间的imgList
	imgList, cnt, err := s.db.LaunchScreen.GetImageBySType(s.ctx, sType)
	if err != nil {
		return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage db.GetImageBySType error:%w", err)
	}
	// 没有符合条件的图片
	if cnt == 0 {
		return nil, 0, errno.NoRunningPictureError
	}

	// 从中筛选出JSON含当前学号的图片(regex)
	cntResp := int64(0)
	currentImgList := make([]model.Picture, 0)
	for _, picture := range *imgList {
		// 处理JSON
		m := make(map[string]string)
		if err = sonic.Unmarshal([]byte(picture.Regex), &m); err != nil {
			return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage unmarshal JSON error:%w", err)
		}

		match := false
		for k, v := range m {
			if k == "picture_id" || k == "device" {
				continue
			}
			// 此时key是student_id,value是学号JSON
			match, err = regexp.MatchString(studentId, v)
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
			return s.db.LaunchScreen.AddImageListShowTime(s.ctx, imgList)
		})
		eg.Go(func() error {
			// setIdCache
			imgIdList := make([]int64, len(*imgList))

			for i, img := range *imgList {
				imgIdList[i] = img.ID
			}
			return s.cache.LaunchScreen.SetLaunchScreenCache(s.ctx, utils.GenerateRedisKeyByStuId(studentId, sType), &imgIdList)
		})
		eg.Go(func() error {
			// setExpireCheckCache
			id, err := s.db.LaunchScreen.GetLastImageId(s.ctx)
			if err != nil {
				return err
			}
			if err = s.cache.LaunchScreen.SetLastLaunchScreenIdCache(s.ctx, id); err != nil {
				return err
			}
			return nil
		})
		if err = eg.Wait(); err != nil {
			return nil, -1, fmt.Errorf("LaunchScreenService.MobileGetImage set cache error:%w", err)
		}
	} else {
		return nil, 0, errno.NoRunningPictureError
	}

	return &currentImgList, cntResp, nil
}
