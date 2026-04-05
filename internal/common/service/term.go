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
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/jwch"
)

func (s *CommonService) GetTermList() (*jwch.SchoolCalendar, error) {
	if s.cache.IsKeyExist(s.ctx, constants.TermListKey) {
		list, err := s.cache.Common.GetTermListCache(s.ctx, constants.TermListKey)
		if err != nil {
			return nil, errno.ErrNoWithPreMessage(err, "Common.GetTermList: Get term list cache failed")
		}
		return list, nil
	}

	calendar, err := jwch.NewStudent().GetSchoolCalendar()
	if err = base.HandleJwchError(err); err != nil {
		return nil, errno.ErrNoWithPreMessage(err, "Common.GetTermList: Get term list failed")
	}
	s.taskQueue.Add(fmt.Sprintf("setTermListCache:%s", constants.TermListKey), taskqueue.QueueTask{Execute: func() error {
		return s.cache.Common.SetTermListCache(s.ctx, constants.TermListKey, calendar)
	}})

	return calendar, nil
}

func (s *CommonService) GetTerm(req *common.TermRequest) (*jwch.CalTermEvents, error) {
	key := s.cache.Common.TermInfoKey(req.Term)
	if s.cache.IsKeyExist(s.ctx, key) {
		events, err := s.cache.Common.GetTermInfo(s.ctx, key)
		if err != nil {
			return nil, errno.ErrNoWithPreMessage(err, "Common.GetTerm: Get term failed")
		}
		return events, nil
	}

	events, err := jwch.NewStudent().GetTermEvents(req.Term)
	if err = base.HandleJwchError(err); err != nil {
		return nil, errno.ErrNoWithPreMessage(err, "Common.GetTerm: Get term failed")
	}
	s.taskQueue.Add(fmt.Sprintf("setTermInfoCache:%s", key), taskqueue.QueueTask{Execute: func() error {
		return s.cache.Common.SetTermInfo(s.ctx, key, events)
	}})

	return events, nil
}
