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

package pack

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func normalizeCourseLocation(location string) string {
	if location == "旗山物理实验教学中心" || location == "铜盘教学楼" {
		return location
	}

	// 去除 {铜盘,旗山} 前缀
	location = strings.TrimPrefix(location, "铜盘")
	location = strings.TrimPrefix(location, "旗山")

	return location
}

func buildScheduleRule(scheduleRule jwch.CourseScheduleRule) *model.CourseScheduleRule {
	return &model.CourseScheduleRule{
		Location:   normalizeCourseLocation(scheduleRule.Location),
		StartClass: int64(scheduleRule.StartClass),
		EndClass:   int64(scheduleRule.EndClass),
		StartWeek:  int64(scheduleRule.StartWeek),
		EndWeek:    int64(scheduleRule.EndWeek),
		Weekday:    int64(scheduleRule.Weekday),
		Single:     scheduleRule.Single,
		Double:     scheduleRule.Double,
		Adjust:     scheduleRule.Adjust,
	}
}

func buildScheduleRules(scheduleRules []jwch.CourseScheduleRule) []*model.CourseScheduleRule {
	var res []*model.CourseScheduleRule
	for _, scheduleRule := range scheduleRules {
		res = append(res, buildScheduleRule(scheduleRule))
	}
	return res
}

func BuildCourse(courses []*jwch.Course) []*model.Course {
	var courseList []*model.Course
	for _, course := range courses {
		courseList = append(courseList, &model.Course{
			Name:             course.Name,
			Syllabus:         course.Syllabus,
			Lessonplan:       course.LessonPlan,
			Teacher:          course.Teacher,
			ScheduleRules:    buildScheduleRules(course.ScheduleRules),
			RawScheduleRules: course.RawScheduleRules,
			RawAdjust:        course.RawAdjust,
			Remark:           course.Remark,
			ExamType:         course.ExamType,
		})
	}
	return courseList
}

func GetTop2Terms(term *jwch.Term) *jwch.Term {
	t := new(jwch.Term)
	if len(term.Terms) <= constants.CourseCacheMaxNum {
		return term
	}
	t.Terms = term.Terms[:constants.CourseCacheMaxNum]
	return t
}

func buildScheduleRuleYjsy(scheduleRule yjsy.CourseScheduleRule) *model.CourseScheduleRule {
	return &model.CourseScheduleRule{
		Location:   normalizeCourseLocation(scheduleRule.Location),
		StartClass: int64(scheduleRule.StartClass),
		EndClass:   int64(scheduleRule.EndClass),
		StartWeek:  int64(scheduleRule.StartWeek),
		EndWeek:    int64(scheduleRule.EndWeek),
		Weekday:    int64(scheduleRule.Weekday),
		Single:     scheduleRule.Single,
		Double:     scheduleRule.Double,
		Adjust:     scheduleRule.Adjust,
	}
}

func buildScheduleRulesYjsy(scheduleRules []yjsy.CourseScheduleRule) []*model.CourseScheduleRule {
	var res []*model.CourseScheduleRule
	for _, scheduleRule := range scheduleRules {
		res = append(res, buildScheduleRuleYjsy(scheduleRule))
	}
	return res
}

func BuildCourseYjsy(courses []*yjsy.Course) []*model.Course {
	var courseList []*model.Course
	for _, course := range courses {
		courseList = append(courseList, &model.Course{
			Name:             course.Name,
			Syllabus:         course.Syllabus,
			Lessonplan:       course.LessonPlan,
			Teacher:          course.Teacher,
			ScheduleRules:    buildScheduleRulesYjsy(course.ScheduleRules),
			RawScheduleRules: course.RawScheduleRules,
			RawAdjust:        course.RawAdjust,
		})
	}
	return courseList
}

func GetTop2TermsYjsy(term *yjsy.Term) *yjsy.Term {
	t := new(yjsy.Term)
	if len(term.Terms) <= constants.CourseCacheMaxNum {
		return term
	}
	t.Terms = term.Terms[:constants.CourseCacheMaxNum]
	return t
}

// GetTop2TermLists 用于提取字符串类型的Top2Term
func GetTop2TermLists(termList []string) []string {
	if len(termList) <= constants.CourseCacheMaxNum {
		return termList
	}
	t := termList[:constants.CourseCacheMaxNum]
	return t
}

func IsYjsyTerm(term string) bool {
	return len(term) == 11 && term[4] == '-' && term[9] == '-'
}

func IsJwchTerm(term string) bool {
	return len(term) == constants.JwchTermLen
}

func MapJwchTerm(term string) string {
	// 202501 → 2024-2025-1
	year := term[:4]
	semester := term[5:]
	currentYear, _ := strconv.Atoi(year)
	prevYear := currentYear - 1
	return fmt.Sprintf("%d-%d-%s", prevYear, currentYear, semester)
}

func MapYjsyTerm(term string) string {
	// 2024-2025-1 → 202501
	return term[5:9] + "0" + term[10:11]
}
