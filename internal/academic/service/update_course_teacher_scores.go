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
	"regexp"
	"strconv"
	"strings"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

var numericScoreRegexp = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)

const (
	InvalidScoreValue   = -10.00 // 无法识别的成绩
	AbsentScoreValue    = -1.00  // 缺考
	CheatingScoreValue  = -2.00  // 作弊
	ExcellentScoreValue = 90.00  // 优秀
	GoodScoreValue      = 80.00  // 良好
	MediumScoreValue    = 68.00  // 中等
	PassScoreValue      = 60.00  // 及格/合格
	FailScoreValue      = 0.00   // 不及格/不合格
)

func convertScore(raw string) float64 {
	if numericScoreRegexp.MatchString(raw) {
		v, err := strconv.ParseFloat(raw, 64)
		if err == nil {
			return v
		}
	}
	switch raw {
	case "优秀", "优":
		return ExcellentScoreValue
	case "良好", "良":
		return GoodScoreValue
	case "中等":
		return MediumScoreValue
	case "及格", "合格":
		return PassScoreValue
	case "不及格", "不合格":
		return FailScoreValue
	case "缺考":
		return AbsentScoreValue
	case "作弊":
		return CheatingScoreValue
	default:
		return InvalidScoreValue
	}
}

func (s *AcademicService) UpdateCourseTeacherScores() error {
	lastStuId := ""
	batchSize := constants.CourseTeacherScoresBatchReadSize

	for {
		scores, err := s.db.Academic.GetScoresBatchByStuId(s.ctx, lastStuId, batchSize)
		if err != nil {
			return fmt.Errorf("UpdateCourseTeacherScores: get scores batch error: %w", err)
		}
		if len(scores) == 0 {
			logger.Infof("UpdateCourseTeacherScores: all records have been processed")
			break
		}

		logger.Infof("UpdateCourseTeacherScores: processing batch, stu_id > %s, count=%d", lastStuId, len(scores))

		var records []*model.CourseTeacherScore
		for _, score := range scores {
			stuIdSHA256 := utils.SHA256(score.StuID)

			var marks []*jwch.Mark
			if err = sonic.UnmarshalString(score.ScoresInfo, &marks); err != nil {
				logger.Errorf("UpdateCourseTeacherScores: unmarshal scores_info for stu_id=%s error: %v", score.StuID, err)
				continue
			}

			for _, mark := range marks {
				teachers := splitTeachers(mark.Teacher)
				numericScore := convertScore(mark.Score)

				for _, teacher := range teachers {
					id, err := s.sf.NextVal()
					if err != nil {
						return fmt.Errorf("UpdateCourseTeacherScores: generate snowflake id error: %w", err)
					}
					records = append(records, &model.CourseTeacherScore{
						ID:           id,
						StuIdSHA256:  stuIdSHA256,
						CourseName:   mark.Name,
						ElectiveType: mark.ElectiveType,
						TeacherName:  teacher,
						Semester:     mark.Semester,
						Score:        numericScore,
					})
				}
			}
		}

		if err = s.db.Academic.UpsertCourseTeacherScores(s.ctx, records); err != nil {
			return fmt.Errorf("UpdateCourseTeacherScores: upsert batch error: %w", err)
		}

		lastStuId = scores[len(scores)-1].StuID
	}

	return nil
}

func splitTeachers(teacherList string) []string {
	if strings.TrimSpace(teacherList) == "" {
		return []string{""}
	}
	parts := strings.Split(teacherList, ",")
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		result = append(result, t)
	}
	return result
}
