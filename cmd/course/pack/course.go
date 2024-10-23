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
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/jwch"
)

func buildScheduleRule(scheduleRule jwch.CourseScheduleRule) *model.CourseScheduleRule {
	return &model.CourseScheduleRule{
		Location:   scheduleRule.Location,
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
		})
	}
	return courseList
}
