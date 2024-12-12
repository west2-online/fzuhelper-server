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

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/jwch"
)

func (s *TermService) GetTermList() (*jwch.SchoolCalendar, error) {
	calendar, err := jwch.NewStudent().GetSchoolCalendar()
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetTermList: Get term list failed %w", err)
	}
	return calendar, nil
}

func (s *TermService) GetTerm(req *common.TermRequest) (*jwch.CalTermEvents, error) {
	events, err := jwch.NewStudent().GetTermEvents(req.Term)
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetTerm: Get term  failed %w", err)
	}
	return events, nil
}
