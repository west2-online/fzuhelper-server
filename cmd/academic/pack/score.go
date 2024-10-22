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

// 将 semester 分成 year 和 term
func parseSemester(semester string) (string, string) {
	year := semester[:4]
	term := semester[4:]
	return year, term
}

func BuildScores(data []*jwch.Mark) []*model.Score {
	scores := make([]*model.Score, len(data))

	for i := 0; i < len(data); i++ {
		year, term := parseSemester(data[i].Semester)
		scores[i] = &model.Score{
			Credit:  data[i].Credits,
			Gpa:     data[i].GPA,
			Name:    data[i].Name,
			Score:   data[i].Score,
			Teacher: data[i].Teacher,
			Term:    term,
			Year:    year,
		}
	}

	return scores
}
