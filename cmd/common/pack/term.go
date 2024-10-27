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

func BuildTermList(termList []jwch.CalTerm) (res []*model.Term) {
	res = make([]*model.Term, 0, len(termList))

	for _, term := range termList {
		res = append(res, &model.Term{
			TermId:     term.TermId,
			SchoolYear: term.SchoolYear,
			Term:       term.Term,
			StartDate:  term.StartDate,
			EndDate:    term.EndDate,
		})
	}

	return res
}

func buildTermEvent(termEvent *jwch.CalTermEvent) *model.TermEvent {
	return &model.TermEvent{
		Name:      termEvent.Name,
		StartDate: termEvent.StartDate,
		EndDate:   termEvent.EndDate,
	}
}

func BuildTermEvents(termEvents []jwch.CalTermEvent) (res []*model.TermEvent) {
	res = make([]*model.TermEvent, 0, len(termEvents))

	for _, termEvent := range termEvents {
		res = append(res, buildTermEvent(&termEvent))
	}

	return res
}
