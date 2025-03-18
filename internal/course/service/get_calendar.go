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
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	ics "github.com/arran4/golang-ical"

	"github.com/west2-online/jwch"
)

// source: https://github.com/renbaoshuo/fzu-ics

// 作息时间
var CLASS_TIME = [][2][2]int{
	{{0, 0}, {23, 59}}, // [[起始小时, 起始分钟], [结束小时, 结束分钟]]
	{{8, 20}, {9, 5}},  // 1
	{{9, 15}, {10, 0}},
	{{10, 20}, {11, 5}},
	{{11, 15}, {12, 0}},
	{{14, 0}, {14, 45}},
	{{14, 55}, {15, 40}},
	{{15, 50}, {16, 35}},
	{{16, 45}, {17, 30}},
	{{19, 0}, {19, 45}},
	{{19, 55}, {20, 40}},
	{{20, 50}, {21, 35}}, // 11
}

func (s *CourseService) GetCalendar(stuID string) ([]byte, error) {
	// 初始化
	cstSh, _ := time.LoadLocation("Asia/Shanghai")
	time.Local = cstSh

	// 转换为 ics 格式
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	cal.SetXWRCalName(fmt.Sprintf("福州大学课程表 [%s]", stuID))
	cal.SetTimezoneId("Asia/Shanghai")
	cal.SetXWRTimezone("Asia/Shanghai")

	// 获取学期开始时间
	calendar, err := jwch.NewStudent().GetSchoolCalendar()
	if err != nil {
		return nil, fmt.Errorf("CourseService.GetCalendar: get student calendar failed: %w", err)
	}
	curTermStartDate, err := time.Parse("2006-01-02", calendar.Terms[0].StartDate)
	if err != nil {
		return nil, fmt.Errorf("CourseService.GetCalendar: parse current term start date failed: %w", err)
	}

	// 获取学期课程表
	courses, err := s.getSemesterCourses(stuID, calendar.CurrentTerm)
	if err != nil {
		return nil, fmt.Errorf("CourseService.GetCalendar: get semester courses failed: %w", err)
	}

	for _, course := range courses {
		for _, scheduleRule := range course.ScheduleRules {
			// TODO: 整周课程处理逻辑

			eventIdBase := fmt.Sprintf("%s__%s_%s_%d-%d_%d_%d-%d_%s_%t_%t",
				calendar.CurrentTerm, course.Name, course.Teacher,
				scheduleRule.StartWeek, scheduleRule.EndWeek, scheduleRule.Weekday,
				scheduleRule.StartClass, scheduleRule.EndClass,
				scheduleRule.Location, scheduleRule.Single, scheduleRule.Double)

			startTime, endTime := calcClassTime(scheduleRule.StartWeek, scheduleRule.Weekday, scheduleRule.StartClass, scheduleRule.EndClass, curTermStartDate)
			_, repeatEndTime := calcClassTime(scheduleRule.EndWeek, scheduleRule.Weekday, scheduleRule.StartClass, scheduleRule.EndClass, curTermStartDate)

			event := cal.AddEvent(md5Str(eventIdBase))
			event.SetCreatedTime(curTermStartDate)
			event.SetDtStampTime(time.Now())
			event.SetModifiedAt(time.Now())
			event.SetSummary(course.Name)
			event.SetLocation(scheduleRule.Location)
			event.SetStartAt(startTime)
			event.SetEndAt(endTime)
			if scheduleRule.Single && scheduleRule.Double { // 单双周都有
				// RRULE:FREQ=WEEKLY;UNTIL=20170101T000000Z
				event.AddRrule("FREQ=WEEKLY;UNTIL=" + repeatEndTime.Format("20060102T150405Z"))
			} else {
				// RRULE:FREQ=WEEKLY;UNTIL=20170101T000000Z;INTERVAL=2
				event.AddRrule("FREQ=WEEKLY;UNTIL=" + repeatEndTime.Format("20060102T150405Z") + ";INTERVAL=2")
			}
		}
	}

	calendarContent := cal.Serialize()

	return []byte(calendarContent), nil
}

func calcClassTime(week int64, weekday int64, startClass int64, endClass int64, dateBase time.Time) (time.Time, time.Time) {
	startHour, startMinute := CLASS_TIME[startClass][0][0], CLASS_TIME[startClass][0][1]
	endHour, endMinute := CLASS_TIME[endClass][1][0], CLASS_TIME[endClass][1][1]

	startTime := dateBase.AddDate(0, 0, int((week-1)*7+(weekday-1)))
	startTime = time.Date(startTime.Year(), startTime.Month(), startTime.Day(), startHour, startMinute, 0, 0, time.Local)
	endTime := dateBase.AddDate(0, 0, int((week-1)*7+(weekday-1)))
	endTime = time.Date(endTime.Year(), endTime.Month(), endTime.Day(), endHour, endMinute, 0, 0, time.Local)

	return startTime, endTime
}

func md5Str(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	fullHash := hex.EncodeToString(hasher.Sum(nil)) // 32-bit (full) hash

	return fullHash
}
