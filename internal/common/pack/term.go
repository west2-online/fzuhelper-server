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

func BuildTermsList(calendar *jwch.SchoolCalendar) *model.TermList {
	return &model.TermList{
		CurrentTerm: &calendar.CurrentTerm,
		Terms:       BuildTerms(calendar.Terms),
	}
}

func BuildTermInfo(term *jwch.CalTermEvents) *model.TermInfo {
	return &model.TermInfo{
		TermId:     &term.TermId,
		Term:       &term.Term,
		SchoolYear: &term.SchoolYear,
		Events:     BuildTermEvents(term.Events),
	}
}

func BuildTerm(term jwch.CalTerm) *model.Term {
	return &model.Term{
		TermId:     &term.TermId,
		SchoolYear: &term.SchoolYear,
		Term:       &term.Term,
		StartDate:  &term.StartDate,
		EndDate:    &term.EndDate,
	}
}

func BuildTerms(terms []jwch.CalTerm) []*model.Term {
	if len(terms) == 0 {
		return nil
	}
	termList := make([]*model.Term, len(terms))
	for i, term := range terms {
		termList[i] = BuildTerm(term)
	}
	return termList
}

func BuildTermEvent(term jwch.CalTermEvent) *model.TermEvent {
	return &model.TermEvent{
		Name:      &term.Name,
		StartDate: &term.StartDate,
		EndDate:   &term.EndDate,
	}
}

func BuildTermEvents(events []jwch.CalTermEvent) []*model.TermEvent {
	if len(events) == 0 {
		return nil
	}

	termEvents := make([]*model.TermEvent, len(events))
	for i, event := range events {
		termEvents[i] = BuildTermEvent(event)
	}
	return termEvents
}
