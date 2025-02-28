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

	"github.com/west2-online/fzuhelper-server/internal/academic/task_model"
	loginmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (s *AcademicService) GetScores(loginData *loginmodel.LoginData) ([]*jwch.Mark, error) {
	key := fmt.Sprintf("scores:%s", context.ExtractIDFromLoginData(loginData))
	loginData, err := context.GetLoginData(s.ctx)
	if err != nil {
		return nil, fmt.Errorf("service.GetScores: Get login data fail %w", err)
	}

	if ok := s.cache.IsKeyExist(s.ctx, key); ok {
		scores, err := s.cache.Academic.GetScoresCache(s.ctx, key)
		if err != nil {
			return nil, fmt.Errorf("service.GetScores: Get scores info from redis error %w", err)
		}
		return scores, nil
	} else {
		stu := jwch.NewStudent().WithLoginData(loginData.Id, utils.ParseCookies(loginData.Cookies))
		scores, err := stu.GetMarks()
		if err = base.HandleJwchError(err); err != nil {
			return nil, fmt.Errorf("service.GetScores: Get scores info fail %w", err)
		}
		s.taskQueue.Add(task_model.NewSetScoresCacheTask(key, scores, s.cache, s.ctx))
		s.taskQueue.Add(task_model.NewPutScoresToDatabaseTask(s.ctx, s.db, context.ExtractIDFromLoginData(loginData), scores))
		return scores, nil
	}
}

func (s *AcademicService) GetScoresYjsy(loginData *loginmodel.LoginData) ([]*yjsy.Mark, error) {
	key := fmt.Sprintf("scores:%s", loginData.Id[len(loginData.Id)-constants.StudentIDLength:])
	if ok := s.cache.IsKeyExist(s.ctx, key); ok {
		scores, err := s.cache.Academic.GetScoresCacheYjsy(s.ctx, key)
		if err != nil {
			return nil, fmt.Errorf("service.GetScoresYjsy: Get scores info from redis error %w", err)
		}
		return scores, nil
	} else {
		stu := yjsy.NewStudent().WithLoginData(utils.ParseCookies(loginData.Cookies))
		scores, err := stu.GetMarks()
		if err = base.HandleYjsyError(err); err != nil {
			return nil, fmt.Errorf("service.GetScoresYjsy: Get scores info fail %w", err)
		}
		s.taskQueue.Add(task_model.NewSetScoresCacheYjsyTask(key, scores, s.cache, s.ctx))
		// 研究生暂时不做成绩推送，也就不做持久化
		return scores, nil
	}
}
