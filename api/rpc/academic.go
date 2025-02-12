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

package rpc

import (
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/academic"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitAcademicRPC() {
	c, err := client.InitAcademicRPC()
	if err != nil {
		logger.Fatalf("api.rpc.academic InitAcademicRPC failed, err  %v", err)
	}
	academicClient = *c
}

func GetScoresRPC(ctx context.Context, req *academic.GetScoresRequest) (scores []*model.Score, err error) {
	resp, err := academicClient.GetScores(ctx, req)
	if err != nil {
		logger.Errorf("GetScoresRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}
	return resp.Scores, nil
}

func GetGPARPC(ctx context.Context, req *academic.GetGPARequest) (gpa *model.GPABean, err error) {
	resp, err := academicClient.GetGPA(ctx, req)
	if err != nil {
		logger.Errorf("GetGPARPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}
	return resp.Gpa, nil
}

func GetCreditRPC(ctx context.Context, req *academic.GetCreditRequest) (credit []*model.Credit, err error) {
	resp, err := academicClient.GetCredit(ctx, req)
	if err != nil {
		logger.Errorf("GetCreditRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}

	return resp.Major, nil
}

func GetUnifiedExamRPC(ctx context.Context, req *academic.GetUnifiedExamRequest) (unifiedExam []*model.UnifiedExam, err error) {
	resp, err := academicClient.GetUnifiedExam(ctx, req)
	if err != nil {
		logger.Errorf("GetUnifiedExamRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}

	return resp.UnifiedExam, nil
}

func GetCultivatePlanRPC(ctx context.Context, req *academic.GetPlanRequest) (string, error) {
	resp, err := academicClient.GetPlan(ctx, req)
	if err != nil {
		logger.Errorf("GetCultivatePlanRPC: RPC called failed: %v", err.Error())
		return "", errno.InternalServiceError.WithMessage(err.Error())
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return "", err
	}

	return resp.Url, nil
}
