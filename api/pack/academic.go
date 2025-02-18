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
	academicModel "github.com/west2-online/fzuhelper-server/api/model/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func BuildScore(res *model.Score) *academicModel.Score {
	return &academicModel.Score{
		Credit:   res.Credit,
		Gpa:      res.Gpa,
		Name:     res.Name,
		Score:    res.Score,
		Teacher:  res.Teacher,
		Term:     res.Term,
		ExamType: res.ExamType,
	}
}

func BuildScoreList(res []*model.Score) []*academicModel.Score {
	list := make([]*academicModel.Score, 0, len(res))
	for _, v := range res {
		list = append(list, BuildScore(v))
	}
	return list
}

func BuildGPA(res *model.GPABean) *academicModel.GPABean {
	gpa := &academicModel.GPABean{Time: res.Time, Data: make([]*academicModel.GPAData, len(res.Data))}

	for i := 0; i < len(res.Data); i++ {
		gpa.Data[i] = &academicModel.GPAData{Type: res.Data[i].Type, Value: res.Data[i].Value}
	}

	return gpa
}

func BuildCredit(res []*model.Credit) []*academicModel.Credit {
	credit := make([]*academicModel.Credit, len(res))
	for i := 0; i < len(res); i++ {
		credit[i] = &academicModel.Credit{Type: res[i].Type, Gain: res[i].Gain, Total: res[i].Total}
	}

	return credit
}

func BuildUnifiedExam(res []*model.UnifiedExam) []*academicModel.UnifiedExam {
	unified := make([]*academicModel.UnifiedExam, len(res))
	for i := 0; i < len(res); i++ {
		unified[i] = &academicModel.UnifiedExam{Name: res[i].Name, Score: res[i].Score, Term: res[i].Term}
	}

	return unified
}
