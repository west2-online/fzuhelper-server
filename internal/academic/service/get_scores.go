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

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *AcademicService) GetScores() ([]*jwch.Mark, error) {
	loginData, err := context.GetLoginData(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetScores: Get login data fail %w", err)
	}
	stu := jwch.NewStudent().WithLoginData(loginData.Id, utils.ParseCookies(loginData.Cookies))
	scores, err := stu.GetMarks()
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetScores: Get scores info fail %w", err)
	}

	return scores, nil
}
