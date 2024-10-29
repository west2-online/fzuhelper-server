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
	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/jwch"
)

func (s *TermService) GetTermList() (*jwch.SchoolCalendar, error) {
	return jwch.NewStudent().GetSchoolCalendar()
}

func (s *TermService) GetTerm(req *common.TermRequest) (*jwch.CalTermEvents, error) {
	// TODO: 验证 req.Term 是否合法
	// 考虑有没有必要额外增加一次请求
	return jwch.NewStudent().GetTermEvents(req.Term)
}
