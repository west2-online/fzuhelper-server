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
	"github.com/west2-online/yjsy"
)

func BuildScores(data []*jwch.Mark) []*model.Score {
	scores := make([]*model.Score, len(data))
	var score string
	for i := 0; i < len(data); i++ {
		score = data[i].Score
		if score == "成绩尚未录入" {
			score = "暂无"
		} else if score == "成绩只录一遍" {
			data[i].Score = "录入中"
		}
		scores[i] = &model.Score{
			Credit:       data[i].Credits,
			Gpa:          data[i].GPA,
			Name:         data[i].Name,
			Score:        score,
			Teacher:      data[i].Teacher,
			Term:         data[i].Semester,
			ExamType:     data[i].ExamType,
			ElectiveType: data[i].ElectiveType,
		}
	}

	return scores
}

func BuildScoresYjsy(data []*yjsy.Mark) []*model.Score {
	scores := make([]*model.Score, len(data))
	var score string
	for i := 0; i < len(data); i++ {
		score = data[i].Score
		if score == "成绩尚未录入" {
			score = "暂无"
		} else if score == "成绩只录一遍" {
			data[i].Score = "录入中"
		}
		scores[i] = &model.Score{
			Credit:       data[i].Credits,
			Gpa:          data[i].GPA,
			Name:         data[i].Name,
			Score:        score,
			Teacher:      data[i].Teacher,
			Term:         data[i].Semester,
			ExamType:     data[i].ExamType,
			ElectiveType: data[i].ElectiveType,
		}
	}

	return scores
}
