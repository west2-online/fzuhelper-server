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

	"github.com/west2-online/fzuhelper-server/kitex_gen/academic"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *AcademicService) GetUnifiedExam(req *academic.GetUnifiedExamRequest) ([]*jwch.UnifiedExam, error) {
	stu := jwch.NewStudent().WithLoginData(req.Id, utils.ParseCookies(req.Cookies))
	cet, err := stu.GetCET()
	if err != nil {
		return nil, fmt.Errorf("service.GetUnifiedExam: Get cet info fail %w", err)
	}
	js, err := stu.GetJS()
	if err != nil {
		return nil, fmt.Errorf("service.GetUnifiedExam: Get js info fail %w", err)
	}
	unifiedExam := append(append([]*jwch.UnifiedExam{}, cet...), js...)
	return unifiedExam, nil
}
