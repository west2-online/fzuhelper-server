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
	courseModel "github.com/west2-online/fzuhelper-server/api/model/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func BuildCourseScheduleRule(res *model.CourseScheduleRule) *courseModel.CourseScheduleRule {
	return &courseModel.CourseScheduleRule{
		Location:   res.Location,
		StartClass: res.StartClass,
		EndClass:   res.EndClass,
		StartWeek:  res.StartWeek,
		EndWeek:    res.EndWeek,
		Weekday:    res.Weekday,
		Single:     res.Single,
		Double:     res.Double,
		Adjust:     res.Adjust,
	}
}

func BuildCourseScheduleRuleList(res []*model.CourseScheduleRule) []*courseModel.CourseScheduleRule {
	list := make([]*courseModel.CourseScheduleRule, 0, len(res))
	for _, v := range res {
		list = append(list, BuildCourseScheduleRule(v))
	}
	return list
}

func BuildCourse(res *model.Course) *courseModel.Course {
	return &courseModel.Course{
		Name:             res.Name,
		Teacher:          res.Teacher,
		ScheduleRules:    BuildCourseScheduleRuleList(res.ScheduleRules),
		Remark:           res.Remark,
		Lessonplan:       res.Lessonplan,
		Syllabus:         res.Syllabus,
		RawScheduleRules: res.RawScheduleRules,
		RawAdjust:        res.RawAdjust,
	}
}

func BuildCourseList(res []*model.Course) []*courseModel.Course {
	list := make([]*courseModel.Course, 0, len(res))
	for _, v := range res {
		list = append(list, BuildCourse(v))
	}
	return list
}

func BuildLocateDate(date *model.LocateDate) *courseModel.LocateDate {
	return &courseModel.LocateDate{
		Week: date.Week,
		Year: date.Year,
		Term: date.Term,
		Date: date.Date,
	}
}
