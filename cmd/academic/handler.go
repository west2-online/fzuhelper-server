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
