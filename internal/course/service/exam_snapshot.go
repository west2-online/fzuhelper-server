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
	// 从选课接口结果中提取考试快照；没有考试时间的课程不参与考试通知。
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

	// 统一考试快照的顺序，避免教务处返回顺序变化导致 hash 变化。
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
	// 按课程身份比较新旧考试时间，支持考试新增、时间变更和考试信息清除。
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
	// map 遍历顺序不固定，排序后保证多条考试变化的处理顺序稳定。
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].ExamHash < changes[j].ExamHash
	})
	return changes
}

func newCourseExamChange(term string, oldExam, newExam CourseExamInfo) courseExamChange {
	// Hash 描述一次具体的考试状态迁移，用于跨用户全局去重。
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
	// 快照 hash 只反映考试内容，不依赖教务处返回的课程排列顺序。
	ordered := append([]CourseExamInfo(nil), exams...)
	// 对副本排序，保证相同考试内容始终生成相同的快照 hash。
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
