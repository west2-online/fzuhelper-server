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
	"time"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/jwch"
)

func (s *CourseService) GetLocateDate() (*model.LocateDate, error) {
	// 获取当前日期和星期几
	currentDate := time.Now()
	formattedCurrentDate := currentDate.Format("2006-01-02 15:04:05")
	currentDay := int(currentDate.Weekday())
	if currentDay == 0 {
		currentDay = 7
	}

	var result *model.LocateDate
	if ok := s.cache.IsKeyExist(s.ctx, constants.LocateDateKey); ok {
		// 解析出缓存的日期和星期几
		cachedLocateDate, err := s.cache.Course.GetDateCache(s.ctx, constants.LocateDateKey)
		if err != nil {
			return nil, err
		}
		// 这里需要指定时区
		cachedDate, err := time.ParseInLocation(time.DateTime, cachedLocateDate.Date, time.Local)
		if err != nil {
			return nil, fmt.Errorf("failed to parse cached date: %w", err)
		}
		cachedDay := int(cachedDate.Weekday())
		if cachedDay == 0 {
			currentDay = 7
		}
		// 判断是否跨周,没跨周直接返回缓存数据
		// 日期相差少于7天，并且星期n（currentDay）是增长的，也就是没有从星期日跨到星期一，就没有跨周
		if currentDate.Sub(cachedDate) < 7*24*time.Hour && currentDay >= cachedDay {
			result = &model.LocateDate{
				Year: cachedLocateDate.Year,
				Week: cachedLocateDate.Week,
				Term: cachedLocateDate.Term,
				Date: formattedCurrentDate,
			}
			return result, nil
		}
	}
	// 缓存不存在或者跨周,重新获取数据
	locateDate, err := jwch.NewStudent().GetLocateDate()
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetLocateDate: Get locate date fail %w", err)
	}
	result = &model.LocateDate{
		Year: locateDate.Year,
		Week: locateDate.Week,
		Term: locateDate.Term,
		Date: formattedCurrentDate,
	}
	s.taskQueue.Add(constants.LocateDateTaskKey, taskqueue.QueueTask{Execute: func() error {
		return cache.SetStructCache(s.cache, s.ctx, constants.LocateDateKey, result, constants.KeyNeverExpire, "Common.SetLocateDate")
	}})
	return result, nil
}
