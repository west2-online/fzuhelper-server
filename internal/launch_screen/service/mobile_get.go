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
	"regexp"

	"github.com/bytedance/sonic"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *LaunchScreenService) MobileGetImage(req *launch_screen.MobileGetImageRequest) (respList *[]model.Picture, cntResp int64, err error) {
	getFromMysql, err := s.shouldGetFromMySQL(req.StudentId, req.SType, req.Device)
	if err != nil {
		return nil, 0, err
	}

	if !getFromMysql {
		// 直接从缓存中获取id
		imgIdList, err := s.cache.LaunchScreen.GetLaunchScreenCache(s.ctx, utils.GenerateRedisKeyByStuId(req.StudentId, req.SType, req.Device))
		if err != nil {
			return nil, -1, errors.Errorf("LaunchScreenService.MobileGetImage cache.GetLaunchScreenCache error:%v", err)
		}
		respList, cntResp, err = s.db.LaunchScreen.GetImageByIdList(s.ctx, &imgIdList)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				getFromMysql = true
			} else {
				return nil, -1, errors.Errorf("LaunchScreenService.MobileGetImage db.GetImageByIdList error:%v", err)
			}
		}
	}

	if getFromMysql {
		return s.getImagesFromMySQL(req.StudentId, req.SType, req.Device)
	}

	// addShowtime for cache
	if cntResp != 0 {
		if err = s.db.LaunchScreen.AddImageListShowTime(s.ctx, respList); err != nil {
			return nil, -1, errors.Errorf("LaunchScreenService.MobileGetImage db.AddImageListShowTime error:%v", err)
		}
	} else {
		return nil, 0, errno.NoRunningPictureError
	}

	return respList, cntResp, nil
}

// shouldGetFromMySQL 确认缓存是否过期或不存在
func (s *LaunchScreenService) shouldGetFromMySQL(studentId string, sType int64, device string) (bool, error) {
	if !s.cache.IsKeyExist(s.ctx, utils.GenerateRedisKeyByStuId(studentId, sType, device)) {
		return true, nil
	}

	if s.cache.LaunchScreen.IsLastLaunchScreenIdCacheExist(s.ctx, device) {
		id, err := s.db.LaunchScreen.GetLastImageId(s.ctx)
		if err != nil {
			return true, errors.Errorf("LaunchScreenService.MobileGetImage cache.GetLaunchScreenCache error: %v", err)
		}
		cacheId, err := s.cache.LaunchScreen.GetLastLaunchScreenIdCache(s.ctx, device)
		if err != nil {
			return true, errors.Errorf("LaunchScreenService.MobileGetImage db.GetImageByIdList error: %v", err)
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
func (s *LaunchScreenService) getImagesFromMySQL(studentId string, sType int64, device string) (*[]model.Picture, int64, error) {
	// 获取符合当前时间的imgList
	imgList, cnt, err := s.db.LaunchScreen.GetImageBySType(s.ctx, sType)
	if err != nil {
		return nil, -1, errors.Errorf("LaunchScreenService.MobileGetImage db.GetImageBySType error:%v", err)
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
			return nil, -1, errors.Errorf("LaunchScreenService.MobileGetImage unmarshal JSON error:%v", err)
		}

		matchId := false
		matchDevice := false
		for k, v := range m {
			switch k {
			case "student_id":
				if v == "" {
					matchId = true
					continue
				}
				matchId, err = regexp.MatchString(studentId, v)
				if err != nil {
					continue
				}
				if !matchId {
					break
				}
			case "device":
				matchDevice, err = regexp.MatchString(device, v)
				if err != nil {
					continue
				}
				if !matchDevice {
					break
				}
			default:
				continue
			}
		}
		if matchId && matchDevice {
			cntResp++
			currentImgList = append(currentImgList, picture)
		}
	}

	if cntResp != 0 {
		var eg errgroup.Group
		eg.Go(func() error {
			// addShowTime
			return s.db.LaunchScreen.AddImageListShowTime(s.ctx, &currentImgList)
		})
		eg.Go(func() error {
			// setIdCache
			imgIdList := make([]int64, len(currentImgList))

			for i, img := range currentImgList {
				imgIdList[i] = img.ID
			}
			return s.cache.LaunchScreen.SetLaunchScreenCache(s.ctx, utils.GenerateRedisKeyByStuId(studentId, sType, device), &imgIdList)
		})
		eg.Go(func() error {
			// setExpireCheckCache
			id, err := s.db.LaunchScreen.GetLastImageId(s.ctx)
			if err != nil {
				return err
			}
			if err = s.cache.LaunchScreen.SetLastLaunchScreenIdCache(s.ctx, id, device); err != nil {
				return err
			}
			return nil
		})
		if err = eg.Wait(); err != nil {
			return nil, -1, errors.Errorf("LaunchScreenService.MobileGetImage set cache error:%v", err)
		}
	} else {
		return nil, 0, errno.NoRunningPictureError
	}

	return &currentImgList, cntResp, nil
}
