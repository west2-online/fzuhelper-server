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
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// 通过 common rpc 获取最新的开学日期和最新学期，返回开学日期，学期，和 yjs 格式的学期
func (s *CourseService) getLatestStartTerm() (string, string, string, error) {
	resp, err := s.commonClient.GetTermsList(s.ctx, &common.TermListRequest{})
	if err != nil {
		return "", "", "", fmt.Errorf("CourseService.getLatestStartTerm: get term list failed: %w", err)
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return "", "", "", err
	}
	// 防止空指针错误，也许有更好的写法？
	if resp.TermLists == nil || resp.TermLists.Terms == nil || resp.TermLists.Terms[0] == nil {
		return "", "", "", errno.NewErrNo(errno.InternalServiceErrorCode, "CourseService.getLatestStartTerm: term list is nil")
	}
	yjsTerm, err := utils.TransformSemester(*resp.TermLists.CurrentTerm)
	if err != nil {
		return "", "", "", fmt.Errorf("CourseService.getLatestStartTerm: transform semester failed: %w", err)
	}
	return *resp.TermLists.Terms[0].StartDate, *resp.TermLists.CurrentTerm, yjsTerm, nil
}
