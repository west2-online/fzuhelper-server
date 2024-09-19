// Code generated by Kitex v0.7.1. DO NOT EDIT.

package userservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	user "github.com/west2-online/fzuhelper-server/kitex_gen/user"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	GetUserInfo(ctx context.Context, req *user.UserInfoRequest, callOptions ...callopt.Option) (r *user.UserInfoResp, err error)
	ValidateCode(ctx context.Context, req *user.ValidateCodeRequest, callOptions ...callopt.Option) (r *user.ValidateCodeResp, err error)
	ChangePassword(ctx context.Context, req *user.ChangePasswordRequest, callOptions ...callopt.Option) (r *user.ChangePasswordResp, err error)
	GetSchoolCalendar(ctx context.Context, req *user.SchoolCalendarRequest, callOptions ...callopt.Option) (r *user.SchoolCalendarResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kUserServiceClient{
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

type kUserServiceClient struct {
	*kClient
}

func (p *kUserServiceClient) GetUserInfo(ctx context.Context, req *user.UserInfoRequest, callOptions ...callopt.Option) (r *user.UserInfoResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetUserInfo(ctx, req)
}

func (p *kUserServiceClient) ValidateCode(ctx context.Context, req *user.ValidateCodeRequest, callOptions ...callopt.Option) (r *user.ValidateCodeResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ValidateCode(ctx, req)
}

func (p *kUserServiceClient) ChangePassword(ctx context.Context, req *user.ChangePasswordRequest, callOptions ...callopt.Option) (r *user.ChangePasswordResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ChangePassword(ctx, req)
}

func (p *kUserServiceClient) GetSchoolCalendar(ctx context.Context, req *user.SchoolCalendarRequest, callOptions ...callopt.Option) (r *user.SchoolCalendarResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetSchoolCalendar(ctx, req)
}