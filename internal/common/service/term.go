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
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
)

func (s *CommonService) GetTermList() (*jwch.SchoolCalendar, error) {
	if s.cache.IsKeyExist(s.ctx, constants.TermListKey) {
		list, err := s.cache.Common.GetTermListCache(s.ctx, constants.TermListKey)
		if err != nil {
			return nil, fmt.Errorf("service.GetTermList: Get term list cache failed %w", err)
		}
		return list, nil
	}

	calendar, err := jwch.NewStudent().GetSchoolCalendar()
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetTermList: Get term list failed %w", err)
	}
	go func() {
		err = s.cache.Common.SetTermListCache(s.ctx, constants.TermListKey, calendar)
		if err != nil {
			logger.Errorf("service.GetTermList: set term list cache failed %v", err)
		}
	}()

	return calendar, nil
}

func (s *CommonService) GetTerm(req *common.TermRequest) (bool, *jwch.CalTermEvents, error) {
	var err error
	var events *jwch.CalTermEvents

	key := s.cache.Common.TermInfoKey(req.Term)
	if ok := s.cache.IsKeyExist(s.ctx, key); ok {
		events, err = s.cache.Common.GetTermInfo(s.ctx, key)
		if err != nil {
			return false, nil, fmt.Errorf("service.GetTerm: Get term  failed %w", err)
		}
		return true, events, nil
	}

	events, err = jwch.NewStudent().GetTermEvents(req.Term)
	if err = base.HandleJwchError(err); err != nil {
		return false, nil, fmt.Errorf("service.GetTerm: Get term  failed %w", err)
	}
	if err = s.cache.Common.SetTermInfo(s.ctx, key, events); err != nil {
		return true, nil, fmt.Errorf("service.GetTerm set term info cache failed %w", err)
	}
	return true, events, err
}
