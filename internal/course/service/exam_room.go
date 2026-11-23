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
	"strings"

	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

// examTimePattern 匹配选课页「考试时间地点」列，实际格式为日期、时间、地点直接相连。
// 例如：2026年06月26日14:00-18:00旗山东3-106。
var examTimePattern = regexp.MustCompile(`^(\d{4}年\d{1,2}月\d{1,2}日)\s*(\d{1,2}:\d{2}\s*-\s*\d{1,2}:\d{2})\s*(.*)$`)

// noExamRoomData 用于未安排考场的课程。
const noExamRoomData = "暂无考场数据"

// buildExamRoomInfo 从「学期选课」页面解析出考场信息。
func buildExamRoomInfo(courses []*jwch.Course) []*kitexModel.ExamRoomInfo {
	rooms := make([]*kitexModel.ExamRoomInfo, 0, len(courses))
	seen := make(map[string]struct{}, len(courses))
	for _, course := range courses {
		if course == nil {
			continue
		}
		room := newExamRoomInfo(course.Name, course.Credits, course.Teacher, course.RawExamTime)
		// 选课页同一门课可能重复出现，考场信息按内容去重。
		identity := examRoomIdentity(room)
		if _, ok := seen[identity]; ok {
			continue
		}
		seen[identity] = struct{}{}
		rooms = append(rooms, room)
	}
	return rooms
}

// buildExamRoomInfoYjsy 与 buildExamRoomInfo 相同，只是数据源为研究生「学期选课」页面。
func buildExamRoomInfoYjsy(courses []*yjsy.Course) []*kitexModel.ExamRoomInfo {
	rooms := make([]*kitexModel.ExamRoomInfo, 0, len(courses))
	seen := make(map[string]struct{}, len(courses))
	for _, course := range courses {
		if course == nil {
			continue
		}
		room := newExamRoomInfo(course.Name, course.Credits, course.Teacher, course.RawExamTime)
		identity := examRoomIdentity(room)
		if _, ok := seen[identity]; ok {
			continue
		}
		seen[identity] = struct{}{}
		rooms = append(rooms, room)
	}
	return rooms
}

func newExamRoomInfo(name, credit, teacher, rawExamTime string) *kitexModel.ExamRoomInfo {
	date, examTime, location, ok := parseExamTimeAndLocation(rawExamTime)
	if !ok {
		// 没有考场信息（或格式无法解析）时，仍保留课程并返回占位值。
		return &kitexModel.ExamRoomInfo{
			Name:     name,
			Credit:   credit,
			Teacher:  teacher,
			Location: noExamRoomData,
		}
	}
	return &kitexModel.ExamRoomInfo{
		Name:     name,
		Credit:   credit,
		Teacher:  teacher,
		Date:     date,
		Time:     examTime,
		Location: location,
	}
}

// parseExamTimeAndLocation 拆分选课页「考试时间地点」列。
// 未安排考试时该列为空字符串。
func parseExamTimeAndLocation(raw string) (date, examTime, location string, ok bool) {
	matches := examTimePattern.FindStringSubmatch(strings.TrimSpace(raw))
	if matches == nil {
		return "", "", "", false
	}
	// 时间段之间即便出现空白也统一去掉，输出给客户端的保持紧凑格式。
	return matches[1], strings.ReplaceAll(matches[2], " ", ""), strings.TrimSpace(matches[3]), true
}

func examRoomIdentity(room *kitexModel.ExamRoomInfo) string {
	return strings.Join([]string{
		room.Name, room.Teacher, room.Credit, room.Date, room.Time, room.Location,
	}, "|")
}
