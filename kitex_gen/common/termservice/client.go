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

// Code generated by Kitex v0.11.3. DO NOT EDIT.

package termservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	common "github.com/west2-online/fzuhelper-server/kitex_gen/common"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	GetTermsList(ctx context.Context, req *common.TermListRequest, callOptions ...callopt.Option) (r *common.TermListResponse, err error)
	GetTerm(ctx context.Context, req *common.TermRequest, callOptions ...callopt.Option) (r *common.TermResponse, err error)
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
	return &kTermServiceClient{
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

type kTermServiceClient struct {
	*kClient
}

func (p *kTermServiceClient) GetTermsList(ctx context.Context, req *common.TermListRequest, callOptions ...callopt.Option) (r *common.TermListResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetTermsList(ctx, req)
}

func (p *kTermServiceClient) GetTerm(ctx context.Context, req *common.TermRequest, callOptions ...callopt.Option) (r *common.TermResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetTerm(ctx, req)
}
