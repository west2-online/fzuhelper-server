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
	// 获取当前日期和星期几，使用中国时区
	currentDate := time.Now().In(constants.ChinaTZ)
	formattedCurrentDate := currentDate.Format(time.DateTime)

	var result *model.LocateDate
	if ok := s.cache.IsKeyExist(s.ctx, constants.LocateDateKey); ok {
		// 解析出缓存的日期和星期几
		cachedLocateDate, err := s.cache.Course.GetDateCache(s.ctx, constants.LocateDateKey)
		if err != nil {
			// 缓存获取失败时，降级到获取新数据
			return s.fetchAndCacheNewDate(formattedCurrentDate)
		}

		// 使用中国时区解析缓存的日期
		cachedDate, err := time.ParseInLocation(time.DateTime, cachedLocateDate.Date, constants.ChinaTZ)
		if err != nil {
			// 日期解析失败时，降级到获取新数据
			return s.fetchAndCacheNewDate(formattedCurrentDate)
		}

		// 判断是否跨周
		currentYear, currentWeek := currentDate.ISOWeek()
		cachedYear, cachedWeek := cachedDate.ISOWeek()

		if currentYear == cachedYear && currentWeek == cachedWeek {
			result = &model.LocateDate{
				Year: cachedLocateDate.Year,
				Week: cachedLocateDate.Week,
				Term: cachedLocateDate.Term,
				Date: formattedCurrentDate,
			}
			return result, nil
		}
	}

	return s.fetchAndCacheNewDate(formattedCurrentDate)
}

func (s *CourseService) fetchAndCacheNewDate(formattedCurrentDate string) (*model.LocateDate, error) {
	// 缓存不存在或者跨周,重新获取数据
	locateDate, err := jwch.NewStudent().GetLocateDate()
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetLocateDate: Get locate date fail %w", err)
	}
	result := &model.LocateDate{
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
