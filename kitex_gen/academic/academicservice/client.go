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

// Code generated by Kitex v0.12.1. DO NOT EDIT.

package academicservice

import (
	"context"

	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"

	academic "github.com/west2-online/fzuhelper-server/kitex_gen/academic"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	GetScores(ctx context.Context, req *academic.GetScoresRequest, callOptions ...callopt.Option) (r *academic.GetScoresResponse, err error)
	GetGPA(ctx context.Context, req *academic.GetGPARequest, callOptions ...callopt.Option) (r *academic.GetGPAResponse, err error)
	GetCredit(ctx context.Context, req *academic.GetCreditRequest, callOptions ...callopt.Option) (r *academic.GetCreditResponse, err error)
	GetUnifiedExam(ctx context.Context, req *academic.GetUnifiedExamRequest, callOptions ...callopt.Option) (r *academic.GetUnifiedExamResponse, err error)
	GetPlan(ctx context.Context, req *academic.GetPlanRequest, callOptions ...callopt.Option) (r *academic.GetPlanResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfoForClient(), options...)
	if err != nil {
		return nil, err
	}
	return &kAcademicServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kAcademicServiceClient struct {
	*kClient
}

func (p *kAcademicServiceClient) GetScores(ctx context.Context, req *academic.GetScoresRequest, callOptions ...callopt.Option) (r *academic.GetScoresResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetScores(ctx, req)
}

func (p *kAcademicServiceClient) GetGPA(ctx context.Context, req *academic.GetGPARequest, callOptions ...callopt.Option) (r *academic.GetGPAResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetGPA(ctx, req)
}

func (p *kAcademicServiceClient) GetCredit(ctx context.Context, req *academic.GetCreditRequest, callOptions ...callopt.Option) (r *academic.GetCreditResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetCredit(ctx, req)
}

func (p *kAcademicServiceClient) GetUnifiedExam(ctx context.Context, req *academic.GetUnifiedExamRequest, callOptions ...callopt.Option) (r *academic.GetUnifiedExamResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUnifiedExam(ctx, req)
}

func (p *kAcademicServiceClient) GetPlan(ctx context.Context, req *academic.GetPlanRequest, callOptions ...callopt.Option) (r *academic.GetPlanResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetPlan(ctx, req)
}
