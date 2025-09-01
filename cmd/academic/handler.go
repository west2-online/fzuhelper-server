package main

import (
	"context"
	academic "github.com/west2-online/fzuhelper-server/kitex_gen/academic"
)

// AcademicServiceImpl implements the last service interface defined in the IDL.
type AcademicServiceImpl struct{}

// GetScores implements the AcademicServiceImpl interface.
func (s *AcademicServiceImpl) GetScores(ctx context.Context, req *academic.GetScoresRequest) (resp *academic.GetScoresResponse, err error) {
	// TODO: Your code here...
	return
}

// GetGPA implements the AcademicServiceImpl interface.
func (s *AcademicServiceImpl) GetGPA(ctx context.Context, req *academic.GetGPARequest) (resp *academic.GetGPAResponse, err error) {
	// TODO: Your code here...
	return
}

// GetCredit implements the AcademicServiceImpl interface.
func (s *AcademicServiceImpl) GetCredit(ctx context.Context, req *academic.GetCreditRequest) (resp *academic.GetCreditResponse, err error) {
	// TODO: Your code here...
	return
}

// GetUnifiedExam implements the AcademicServiceImpl interface.
func (s *AcademicServiceImpl) GetUnifiedExam(ctx context.Context, req *academic.GetUnifiedExamRequest) (resp *academic.GetUnifiedExamResponse, err error) {
	// TODO: Your code here...
	return
}

// GetPlan implements the AcademicServiceImpl interface.
func (s *AcademicServiceImpl) GetPlan(ctx context.Context, req *academic.GetPlanRequest) (resp *academic.GetPlanResponse, err error) {
	// TODO: Your code here...
	return
}
