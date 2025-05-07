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

package academic

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

const (
	scoreExcellent = 89.90
	scoreGood      = 79.90
	scoreMedium    = 69.90
	scorePass      = 59.90

	scoreInvalid    = -2.00
	scoreFailed     = -1.00
	scoreParamCount = 6
)

var numRegex = regexp.MustCompile(`^[0-9]+(\.[0-9]+)?$`)

type ScoreInfo struct {
	Name        string `json:"name"`
	TeacherList string `json:"teacher"`
	Semester    string `json:"semester"`
	Score       string `json:"score"`
}

// CourseTeacherScoreRecord 用于承载一行待写入的数据
type CourseTeacherScoreRecord struct {
	ID          int64
	StuIdSha256 string
	CourseName  string
	TeacherName string
	Semester    string
	Score       float64
}

// mapScore 对应 SQL 中的 CASE ... END
func mapScore(s string) float64 {
	switch {
	case numRegex.MatchString(s):
		if v, err := strconv.ParseFloat(s, 64); err == nil {
			return v
		}
	case s == "优秀" || s == "优":
		return scoreExcellent
	case s == "良好" || s == "良":
		return scoreGood
	case s == "中等":
		return scoreMedium
	case s == "及格" || s == "合格":
		return scorePass
	case s == "不及格" || s == "不合格":
		return scoreFailed
	}
	return scoreInvalid
}

// generateSHA256 对应 SQL 中的 SHA2(stu_id,256)
func generateSHA256(str string) string {
	h := sha256.Sum256([]byte(str))
	return hex.EncodeToString(h[:])
}

func (c *DBAcademic) UpdateCourseTeacherScores(ctx context.Context, offset, limit int) error {
	rows, err := c.client.WithContext(ctx).
		Table(constants.ScoreTableName).
		Select("stu_id", "scores_info").
		Offset(offset).
		Limit(limit).
		Rows()
	if err != nil {
		return errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.GetStuScoresInfoRows error: %v", err))
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			// 这里不需要返回错误，直接打印日志即可
			logger.Errorf("dal.GetStuScoresInfoRows rows.Close error: %v", err)
		}
	}(rows)

	var recs []CourseTeacherScoreRecord
	for rows.Next() {
		var stuID string
		var raw []byte
		if err := rows.Scan(&stuID, &raw); err != nil {
			return err
		}
		// 第一次 JSON_ARRAY 展开
		var infos []ScoreInfo
		if err := json.Unmarshal(raw, &infos); err != nil {
			return err
		}
		for _, info := range infos {
			// 第二次 teacher_list 拆分
			teachers := strings.Split(info.TeacherList, ",")
			for _, t := range teachers {
				teacher := strings.TrimSpace(t)
				idVal, err := c.sf.NextVal()
				if err != nil {
					return err
				}
				recs = append(recs, CourseTeacherScoreRecord{
					ID:          idVal,
					StuIdSha256: generateSHA256(stuID),
					CourseName:  info.Name,
					TeacherName: teacher,
					Semester:    info.Semester,
					Score:       mapScore(info.Score),
				})
			}
		}
	}
	if len(recs) == 0 {
		return nil
	}
	// 批量插入或更新
	if err := c.BulkUpsertCourseTeacherScores(ctx, recs); err != nil {
		return err
	}
	logger.Infof("Processed page offset=%d, count=%d", offset, len(recs))
	return nil
}

// BulkUpsertCourseTeacherScores 批量插入或更新 course_teacher_scores 表
func (c *DBAcademic) BulkUpsertCourseTeacherScores(ctx context.Context, recs []CourseTeacherScoreRecord) error {
	if len(recs) == 0 {
		return nil
	}

	// VALUES 部分
	var (
		sb   strings.Builder
		args = make([]interface{}, 0, len(recs)*scoreParamCount) // 每条 6 个占位符
	)
	sb.WriteString("INSERT INTO ")
	sb.WriteString(constants.CourseTeacherScoresTableName)
	sb.WriteString(" (id, stu_id_sha256, course_name, teacher_name, semester, score) VALUES ")

	for i, r := range recs {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString("(?,?,?,?,?,?)")
		args = append(args,
			r.ID,
			r.StuIdSha256,
			r.CourseName,
			r.TeacherName,
			r.Semester,
			r.Score,
		)
	}

	// ON DUPLICATE KEY UPDATE 部分
	sb.WriteString(" ON DUPLICATE KEY UPDATE ")
	sb.WriteString("score = IF(VALUES(score) <> ")
	sb.WriteString(constants.CourseTeacherScoresTableName)
	sb.WriteString(".score, VALUES(score), ")
	sb.WriteString(constants.CourseTeacherScoresTableName)
	sb.WriteString(".score)")

	// 执行
	if err := c.client.
		WithContext(ctx).
		Exec(sb.String(), args...).
		Error; err != nil {
		return errno.NewErrNo(
			errno.InternalDatabaseErrorCode,
			fmt.Sprintf("BulkUpsertCourseTeacherScores failed: %v", err),
		)
	}
	return nil
}
