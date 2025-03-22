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
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// 通过 common rpc 获取最新的开学日期和最新学期
func (s *CourseService) getLatestStartTerm() (string, string, error) {
	// 获取最新 term
	latestTerm, err := s.GetLocateDate()
	if err != nil {
		return "", "", fmt.Errorf("CourseService.GetCalendar: get locate date failed: %w", err)
	}
	term, err := s.commonClient.GetTerm(s.ctx, &common.TermRequest{Term: latestTerm.Year + latestTerm.Term})
	if err != nil {
		return "", "", fmt.Errorf("CourseService.GetCalendar: get term failed: %w", err)
	}
	if err = utils.HandleBaseRespWithCookie(term.Base); err != nil {
		return "", "", err
	}
	// 防止空指针错误，也许有更好的写法？
	if term.TermInfo == nil || term.TermInfo.Events == nil || term.TermInfo.Events[0] == nil || term.TermInfo.Events[0].StartDate == nil {
		return "", "", fmt.Errorf("CourseService.GetCalendar: get term info failed: term is nil")
	}
	return *term.TermInfo.Events[0].StartDate, *term.TermInfo.Term, nil
}
