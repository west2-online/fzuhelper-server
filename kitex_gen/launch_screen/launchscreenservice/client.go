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

// Code generated by Kitex v0.9.1. DO NOT EDIT.

package launchscreenservice

import (
	"context"

	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"

	launch_screen "github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	CreateImage(ctx context.Context, req *launch_screen.CreateImageRequest, callOptions ...callopt.Option) (r *launch_screen.CreateImageResponse, err error)
	GetImage(ctx context.Context, req *launch_screen.GetImageRequest, callOptions ...callopt.Option) (r *launch_screen.GetImageResponse, err error)
	GetImagesByUserId(ctx context.Context, req *launch_screen.GetImagesByUserIdRequest, callOptions ...callopt.Option) (r *launch_screen.GetImagesByUserIdResponse, err error)
	ChangeImageProperty(ctx context.Context, req *launch_screen.ChangeImagePropertyRequest, callOptions ...callopt.Option) (r *launch_screen.ChangeImagePropertyResponse, err error)
	ChangeImage(ctx context.Context, req *launch_screen.ChangeImageRequest, callOptions ...callopt.Option) (r *launch_screen.ChangeImageResponse, err error)
	DeleteImage(ctx context.Context, req *launch_screen.DeleteImageRequest, callOptions ...callopt.Option) (r *launch_screen.DeleteImageResponse, err error)
	MobileGetImage(ctx context.Context, req *launch_screen.MobileGetImageRequest, callOptions ...callopt.Option) (r *launch_screen.MobileGetImageResponse, err error)
	AddImagePointTime(ctx context.Context, req *launch_screen.AddImagePointTimeRequest, callOptions ...callopt.Option) (r *launch_screen.AddImagePointTimeResponse, err error)
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
	return &kLaunchScreenServiceClient{
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

type kLaunchScreenServiceClient struct {
	*kClient
}

func (p *kLaunchScreenServiceClient) CreateImage(ctx context.Context, req *launch_screen.CreateImageRequest, callOptions ...callopt.Option) (r *launch_screen.CreateImageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.CreateImage(ctx, req)
}

func (p *kLaunchScreenServiceClient) GetImage(ctx context.Context, req *launch_screen.GetImageRequest, callOptions ...callopt.Option) (r *launch_screen.GetImageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetImage(ctx, req)
}

func (p *kLaunchScreenServiceClient) GetImagesByUserId(ctx context.Context, req *launch_screen.GetImagesByUserIdRequest, callOptions ...callopt.Option) (r *launch_screen.GetImagesByUserIdResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetImagesByUserId(ctx, req)
}

func (p *kLaunchScreenServiceClient) ChangeImageProperty(ctx context.Context, req *launch_screen.ChangeImagePropertyRequest, callOptions ...callopt.Option) (r *launch_screen.ChangeImagePropertyResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ChangeImageProperty(ctx, req)
}

func (p *kLaunchScreenServiceClient) ChangeImage(ctx context.Context, req *launch_screen.ChangeImageRequest, callOptions ...callopt.Option) (r *launch_screen.ChangeImageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ChangeImage(ctx, req)
}

func (p *kLaunchScreenServiceClient) DeleteImage(ctx context.Context, req *launch_screen.DeleteImageRequest, callOptions ...callopt.Option) (r *launch_screen.DeleteImageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DeleteImage(ctx, req)
}

func (p *kLaunchScreenServiceClient) MobileGetImage(ctx context.Context, req *launch_screen.MobileGetImageRequest, callOptions ...callopt.Option) (r *launch_screen.MobileGetImageResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.MobileGetImage(ctx, req)
}

func (p *kLaunchScreenServiceClient) AddImagePointTime(ctx context.Context, req *launch_screen.AddImagePointTimeRequest, callOptions ...callopt.Option) (r *launch_screen.AddImagePointTimeResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.AddImagePointTime(ctx, req)
}
