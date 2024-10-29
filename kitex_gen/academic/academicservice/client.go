// Code generated by Kitex v0.11.3. DO NOT EDIT.

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
