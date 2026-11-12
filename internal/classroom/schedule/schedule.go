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

// Package schedule 计算空教室刷新的触发时刻和日期范围。
package schedule

import (
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// NextDelay 返回距离下一个刷新时刻的时长；当天时刻都已过去时顺延到次日第一个时刻。
func NextDelay(now time.Time) time.Duration {
	for _, hour := range constants.ClassroomRefreshHours {
		next := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, now.Location())
		if next.After(now) {
			return next.Sub(now)
		}
	}
	next := time.Date(now.Year(), now.Month(), now.Day()+1, constants.ClassroomRefreshHours[0], 0, 0, 0, now.Location())
	return next.Sub(now)
}

// Dates 返回从当天到下周日的所有日期（含两端）。
func Dates(now time.Time) []time.Time {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	// time.Weekday 中周日为 0，所以本周剩余到周日的天数为 (7-今天) % 7。
	daysUntilSunday := (7 - int(today.Weekday())) % 7
	end := today.AddDate(0, 0, daysUntilSunday+7) // 下周日
	dates := make([]time.Time, 0, daysUntilSunday+8)
	for day := today; !day.After(end); day = day.AddDate(0, 0, 1) {
		dates = append(dates, day)
	}
	return dates
}
