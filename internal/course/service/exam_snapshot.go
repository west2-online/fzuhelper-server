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
	"sort"
	"strings"

	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

// CourseExamInfo stores the internal exam snapshot extracted from a course.
type CourseExamInfo struct {
	Name     string `json:"name"`
	Teacher  string `json:"teacher"`
	Credit   string `json:"credit"`
	ExamTime string `json:"exam_time"`
}

func buildCourseExamInfo(courses []*jwch.Course) []CourseExamInfo {
	exams := make([]CourseExamInfo, 0, len(courses))
	for _, course := range courses {
		if course == nil || strings.TrimSpace(course.RawExamTime) == "" {
			continue
		}
		exams = append(exams, CourseExamInfo{
			Name:     course.Name,
			Teacher:  course.Teacher,
			Credit:   course.Credits,
			ExamTime: strings.TrimSpace(course.RawExamTime),
		})
	}

	sort.Slice(exams, func(i, j int) bool {
		left, right := courseExamIdentity(exams[i]), courseExamIdentity(exams[j])
		if left != right {
			return left < right
		}
		return exams[i].ExamTime < exams[j].ExamTime
	})
	return exams
}

func courseExamIdentity(exam CourseExamInfo) string {
	return strings.Join([]string{exam.Name, exam.Teacher, exam.Credit}, "|")
}

func courseExamTag(exam CourseExamInfo) string {
	return utils.MD5(courseExamIdentity(exam))
}

type courseExamChange struct {
	ExamHash string
	Tag      string
	Exam     CourseExamInfo
}

func buildCourseExamChanges(term string, oldExams, newExams []CourseExamInfo) []courseExamChange {
	oldByIdentity := make(map[string]CourseExamInfo, len(oldExams))
	for _, exam := range oldExams {
		oldByIdentity[courseExamIdentity(exam)] = exam
	}

	newByIdentity := make(map[string]CourseExamInfo, len(newExams))
	for _, exam := range newExams {
		newByIdentity[courseExamIdentity(exam)] = exam
	}

	identities := make(map[string]struct{}, len(oldByIdentity)+len(newByIdentity))
	for identity := range oldByIdentity {
		identities[identity] = struct{}{}
	}
	for identity := range newByIdentity {
		identities[identity] = struct{}{}
	}

	changes := make([]courseExamChange, 0)
	for identity := range identities {
		oldExam := oldByIdentity[identity]
		newExam := newByIdentity[identity]
		if oldExam.ExamTime == newExam.ExamTime {
			continue
		}
		if newExam.Name == "" {
			newExam = oldExam
			newExam.ExamTime = ""
		}
		changes = append(changes, newCourseExamChange(term, oldExam, newExam))
	}
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].ExamHash < changes[j].ExamHash
	})
	return changes
}

func newCourseExamChange(term string, oldExam, newExam CourseExamInfo) courseExamChange {
	examHash := utils.SHA256(strings.Join([]string{
		newExam.Name,
		term,
		newExam.Teacher,
		newExam.Credit,
		oldExam.ExamTime,
		newExam.ExamTime,
	}, "|"))
	return courseExamChange{
		ExamHash: examHash,
		Tag:      courseExamTag(newExam),
		Exam:     newExam,
	}
}

func courseExamInfoHash(exams []CourseExamInfo) (string, error) {
	ordered := append([]CourseExamInfo(nil), exams...)
	sort.Slice(ordered, func(i, j int) bool {
		left, right := courseExamIdentity(ordered[i]), courseExamIdentity(ordered[j])
		if left != right {
			return left < right
		}
		return ordered[i].ExamTime < ordered[j].ExamTime
	})

	data, err := utils.JSONEncode(ordered)
	if err != nil {
		return "", err
	}
	return utils.SHA256(data), nil
}
