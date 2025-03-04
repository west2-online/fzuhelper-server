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
	api "github.com/west2-online/fzuhelper-server/api/model/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
)

func BuildTermList(termList *model.TermList) *api.TermList {
	return &api.TermList{
		CurrentTerm: termList.CurrentTerm,
		Terms:       BuildTerms(termList.Terms),
	}
}

func BuildTerms(termList []*model.Term) []*api.Term {
	return base.BuildTypeList(termList, BuildTerm)
}

func BuildTerm(term *model.Term) *api.Term {
	return &api.Term{
		TermID:     term.TermId,
		SchoolYear: term.SchoolYear,
		Term:       term.Term,
		StartDate:  term.StartDate,
		EndDate:    term.EndDate,
	}
}

func BuildTermInfo(termInfo *model.TermInfo) *api.TermInfo {
	return &api.TermInfo{
		TermID:     termInfo.TermId,
		SchoolYear: termInfo.SchoolYear,
		Term:       termInfo.Term,
		Events:     BuildTermEvents(termInfo.Events),
	}
}

func BuildTermEvent(termEvent *model.TermEvent) *api.TermEvent {
	return &api.TermEvent{
		Name:      termEvent.Name,
		StartDate: termEvent.StartDate,
		EndDate:   termEvent.EndDate,
	}
}

func BuildTermEvents(termEvents []*model.TermEvent) []*api.TermEvent {
	return base.BuildTypeList(termEvents, BuildTermEvent)
}

func BuildContributor(contributor *model.Contributor) *api.Contributor {
	return &api.Contributor{
		Name:          contributor.Name,
		AvatarURL:     contributor.AvatarUrl,
		URL:           contributor.Url,
		Contributions: contributor.Contributions,
	}
}

func BuildContributors(contributors []*model.Contributor) []*api.Contributor {
	return base.BuildTypeList(contributors, BuildContributor)
}
