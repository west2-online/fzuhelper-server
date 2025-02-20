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

package academic

import (
	"context"

	"github.com/west2-online/fzuhelper-server/internal/academic/pack"
	"github.com/west2-online/fzuhelper-server/internal/academic/service"
	"github.com/west2-online/fzuhelper-server/internal/academic/syncer"
	"github.com/west2-online/fzuhelper-server/kitex_gen/academic"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
)

// AcademicServiceImpl implements the last service interface defined in the IDL.
type AcademicServiceImpl struct {
	ClientSet *base.ClientSet
	Syncer    *syncer.AcademicSyncer
}

func NewAcademicService(clientSet *base.ClientSet, syncer *syncer.AcademicSyncer) *AcademicServiceImpl {
	return &AcademicServiceImpl{
		ClientSet: clientSet,
		Syncer:    syncer,
	}
}

// GetScores implements the AcademicServiceImpl interface.
func (s *AcademicServiceImpl) GetScores(ctx context.Context, _ *academic.GetScoresRequest) (resp *academic.GetScoresResponse, err error) {
	resp = academic.NewGetScoresResponse()
	var scores []*jwch.Mark

	scores, err = service.NewAcademicService(ctx, s.ClientSet, s.Syncer).GetScores()
	if err != nil {
		logger.Infof("Academic.GetScores: GetScores failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = base.BuildSuccessResp()
	resp.Scores = pack.BuildScores(scores)
	return resp, nil
}

// GetGPA implements the AcademicServiceImpl interface.
func (s *AcademicServiceImpl) GetGPA(ctx context.Context, _ *academic.GetGPARequest) (resp *academic.GetGPAResponse, err error) {
	resp = academic.NewGetGPAResponse()
	var gpa *jwch.GPABean

	gpa, err = service.NewAcademicService(ctx, s.ClientSet, nil).GetGPA()
	if err != nil {
		logger.Infof("Academic.GetGPA: GetGPA failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Gpa = pack.BuildGPA(gpa)
	return resp, nil
}

// GetCredit implements the AcademicServiceImpl interface.
func (s *AcademicServiceImpl) GetCredit(ctx context.Context, _ *academic.GetCreditRequest) (resp *academic.GetCreditResponse, err error) {
	resp = academic.NewGetCreditResponse()
	var credit []*jwch.CreditStatistics

	credit, err = service.NewAcademicService(ctx, s.ClientSet, nil).GetCredit()
	if err != nil {
		logger.Infof("Academic.GetCredit: GetCredit failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Major = pack.BuildCredit(credit)
	// TODO:辨别本专业和辅修专业
	return resp, nil
}

// GetUnifiedExam implements the AcademicServiceImpl interface.
func (s *AcademicServiceImpl) GetUnifiedExam(ctx context.Context, _ *academic.GetUnifiedExamRequest) (resp *academic.GetUnifiedExamResponse, err error) {
	resp = academic.NewGetUnifiedExamResponse()
	var unifiedExam []*jwch.UnifiedExam

	unifiedExam, err = service.NewAcademicService(ctx, s.ClientSet, nil).GetUnifiedExam()
	if err != nil {
		logger.Infof("Academic.GetUnifiedExam: GetUnifiedExam failed, err: %v", err)
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = base.BuildSuccessResp()
	resp.UnifiedExam = pack.BuildUnifiedExam(unifiedExam)
	return resp, nil
}

// GetPlan implements the AcademicServiceImpl interface.
func (s *AcademicServiceImpl) GetPlan(ctx context.Context, _ *academic.GetPlanRequest) (resp *academic.GetPlanResponse, err error) {
	resp = new(academic.GetPlanResponse)
	url, err := service.NewAcademicService(ctx, s.ClientSet, nil).GetPlan()
	if err != nil {
		resp.Base = base.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = base.BuildSuccessResp()
	resp.Url = url
	return resp, nil
}
