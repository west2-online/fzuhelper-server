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
	"context"
	"errors"
	"testing"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// mockCommonClient implements commonservice.Client for testing getLatestStartTerm.
type mockCommonClient struct {
	termResp *common.TermListResponse
	termErr  error
}

func (m *mockCommonClient) GetTermsList(ctx context.Context, req *common.TermListRequest, opts ...callopt.Option) (*common.TermListResponse, error) {
	return m.termResp, m.termErr
}

// unused methods
func (m *mockCommonClient) GetCSS(context.Context, *common.GetCSSRequest, ...callopt.Option) (*common.GetCSSResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCommonClient) GetHtml(context.Context, *common.GetHtmlRequest, ...callopt.Option) (*common.GetHtmlResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCommonClient) GetUserAgreement(context.Context, *common.GetUserAgreementRequest, ...callopt.Option) (*common.GetUserAgreementResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCommonClient) GetTerm(context.Context, *common.TermRequest, ...callopt.Option) (*common.TermResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCommonClient) GetNotices(context.Context, *common.NoticeRequest, ...callopt.Option) (*common.NoticeResponse, error) {
	return nil, errors.New("not implemented")
}

//nolint:lll
func (m *mockCommonClient) GetContributorInfo(context.Context, *common.GetContributorInfoRequest, ...callopt.Option) (*common.GetContributorInfoResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCommonClient) GetToolboxConfig(context.Context, *common.GetToolboxConfigRequest, ...callopt.Option) (*common.GetToolboxConfigResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockCommonClient) PutToolboxConfig(context.Context, *common.PutToolboxConfigRequest, ...callopt.Option) (*common.PutToolboxConfigResponse, error) {
	return nil, errors.New("not implemented")
}

func TestCourseService_getLatestStartTerm(t *testing.T) {
	successBase := &kitexModel.BaseResp{Code: errno.SuccessCode, Msg: "ok"}

	mkTerm := func(start string) *kitexModel.Term {
		return &kitexModel.Term{StartDate: &start}
	}

	cases := []struct {
		name          string
		resp          *common.TermListResponse
		rpcErr        error
		expectErr     bool
		errContains   string
		expectStart   string
		expectTerm    string
		expectYjsTerm string
	}{
		{
			name:        "rpc error",
			rpcErr:      errors.New("rpc fail"),
			expectErr:   true,
			errContains: "get term list failed",
		},
		{
			name:        "base error",
			resp:        &common.TermListResponse{Base: &kitexModel.BaseResp{Code: errno.BizJwchEvaluationNotFoundCode, Msg: "bad"}},
			expectErr:   true,
			errContains: "bad",
		},
		{
			name:        "term list nil",
			resp:        &common.TermListResponse{Base: successBase},
			expectErr:   true,
			errContains: "term list is nil",
		},
		{
			name: "transform error",
			resp: &common.TermListResponse{
				Base: successBase,
				TermLists: &kitexModel.TermList{
					CurrentTerm: strPtr("123"),
					Terms:       []*kitexModel.Term{mkTerm("2024-09-01")},
				},
			},
			expectErr:   true,
			errContains: "transform semester failed",
		},
		{
			name: "success",
			resp: &common.TermListResponse{
				Base: successBase,
				TermLists: &kitexModel.TermList{
					CurrentTerm: strPtr("202401"),
					Terms:       []*kitexModel.Term{mkTerm("2024-09-01")},
				},
			},
			expectStart:   "2024-09-01",
			expectTerm:    "202401",
			expectYjsTerm: "2024-2025-1",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			svc := &CourseService{ctx: context.Background(), commonClient: &mockCommonClient{termResp: tc.resp, termErr: tc.rpcErr}}

			start, term, yjsTerm, err := svc.getLatestStartTerm()

			if tc.expectErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectStart, start)
			assert.Equal(t, tc.expectTerm, term)
			assert.Equal(t, tc.expectYjsTerm, yjsTerm)
		})
	}
}

func strPtr(s string) *string { return &s }
