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
	"strings"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	dbModel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func normalizeCourseLocation(location string, isGraduate bool) string {
	if isGraduate {
		return location
	}

	if location == "旗山物理实验教学中心" || location == "铜盘教学楼" || strings.HasPrefix(location, "晋江校区") {
		return location
	}

	// 非研究生 去除 {铜盘,旗山,晋江} 前缀
	// 先去晋江的前缀 以防晋江楼被误删
	location = strings.TrimPrefix(location, "晋江")
	location = strings.TrimPrefix(location, "铜盘")
	location = strings.TrimPrefix(location, "旗山")

	return location
}

func buildScheduleRule(scheduleRule jwch.CourseScheduleRule) *model.CourseScheduleRule {
	return &model.CourseScheduleRule{
		Location:   normalizeCourseLocation(scheduleRule.Location, false),
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

func buildAdjustRule(adjustRule jwch.CourseAdjustRule) *model.CourseAdjustRule {
	return &model.CourseAdjustRule{
		OldWeek:        int64(adjustRule.OldWeek),
		OldDay:         int64(adjustRule.OldWeekday),
		OldStartClass:  int64(adjustRule.OldStartClass),
		OldEndClass:    int64(adjustRule.OldEndClass),
		Canceled:       adjustRule.Canceled,
		NewWeek_:       int64(adjustRule.NewWeek),
		NewDay_:        int64(adjustRule.NewWeekday),
		NewStartClass_: int64(adjustRule.NewStartClass),
		NewEndClass_:   int64(adjustRule.NewEndClass),
		NewLocation_:   normalizeCourseLocation(adjustRule.NewLocation, false),
	}
}

func buildAdjustRules(adjustRules []jwch.CourseAdjustRule) []*model.CourseAdjustRule {
	var res []*model.CourseAdjustRule
	for _, adjustRule := range adjustRules {
		res = append(res, buildAdjustRule(adjustRule))
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
			AdjustRules:      buildAdjustRules(course.AdjustRules),
			RawScheduleRules: course.RawScheduleRules,
			RawAdjust:        course.RawAdjust,
			Remark:           course.Remark,
			ExamType:         course.ExamType,
			ElectiveType:     course.ElectiveType,
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
		Location:   normalizeCourseLocation(scheduleRule.Location, true),
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
			ElectiveType:     course.ElectiveType,
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

// BuildTermOnDB 用于转换成存储在db中的termList
func BuildTermOnDB(termList []string) string {
	return strings.Join(termList, "|")
}

// ParseTerm 用于db中的termList转换为string数组
func ParseTerm(termList string) []string {
	if termList == "" {
		return nil
	}
	return strings.Split(termList, "|")
}

func ToJwchScheduleRules(rules []*model.CourseScheduleRule) []jwch.CourseScheduleRule {
	res := make([]jwch.CourseScheduleRule, 0, len(rules))
	for _, r := range rules {
		res = append(res, jwch.CourseScheduleRule{
			Location:   r.Location,
			StartClass: int(r.StartClass),
			EndClass:   int(r.EndClass),
			StartWeek:  int(r.StartWeek),
			EndWeek:    int(r.EndWeek),
			Weekday:    int(r.Weekday),
			Single:     r.Single,
			Double:     r.Double,
			Adjust:     r.Adjust,
		})
	}
	return res
}

func FromJwchScheduleRules(rules []jwch.CourseScheduleRule) []*model.CourseScheduleRule {
	return buildScheduleRules(rules)
}

func BuildAdjustCourse(c *dbModel.AutoAdjustCourse) *model.AdjustCourse {
	toDate := ""
	if c.ToDate != nil {
		toDate = *c.ToDate
	}

	var toWeek int64
	if c.ToWeek != nil {
		toWeek = *c.ToWeek
	}

	var toWeekday int64
	if c.ToWeekday != nil {
		toWeekday = *c.ToWeekday
	}

	return &model.AdjustCourse{
		Id:          c.Id,
		Enabled:     c.Enabled,
		Year:        c.Year,
		Term:        c.Term,
		FromDate:    c.FromDate,
		FromWeek:    c.FromWeek,
		FromWeekday: c.FromWeekday,
		ToDate:      toDate,
		ToWeek:      toWeek,
		ToWeekday:   toWeekday,
	}
}

func BuildAdjustCourseList(list []*dbModel.AutoAdjustCourse) []*model.AdjustCourse {
	res := make([]*model.AdjustCourse, 0, len(list))
	for _, c := range list {
		res = append(res, BuildAdjustCourse(c))
	}
	return res
}
