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

	"github.com/west2-online/fzuhelper-server/internal/academic/pack"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

// GetCreditV2 获取V2版本的学分统计
func (s *AcademicService) GetCreditV2() (*pack.CreditResponse, error) {
	loginData, err := context.GetLoginData(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetCreditV2: Get login data fail %w", err)
	}
	stu := jwch.NewStudent().WithLoginData(loginData.Id, utils.ParseCookies(loginData.Cookies))

	majorCredits, minorCredits, err := stu.GetCreditV2()
	if err = base.HandleJwchError(err); err != nil {
		return nil, fmt.Errorf("service.GetCreditV2: Get credit fail %w", err)
	}

	majorItem := pack.BuildCreditCategory("主修专业", majorCredits)
	var response pack.CreditResponse
	response = append(response, majorItem)

	if len(minorCredits) > 0 {
		minorItem := pack.BuildCreditCategory("辅修专业", minorCredits)
		response = append(response, minorItem)
	}

	return &response, nil
}
