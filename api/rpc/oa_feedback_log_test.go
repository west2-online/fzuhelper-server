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
	"errors"
	"testing"

	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/require"

	oaimpl "github.com/west2-online/fzuhelper-server/internal/oa"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/oa"
	"github.com/west2-online/fzuhelper-server/kitex_gen/oa/oaservice"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

type feedbackLogRPCClient struct {
	oaservice.Client
	call func(context.Context, *oa.UploadFeedbackLogRequest) (*oa.UploadFeedbackLogResponse, error)
}

func (c feedbackLogRPCClient) UploadFeedbackLog(ctx context.Context, req *oa.UploadFeedbackLogRequest,
	_ ...callopt.Option,
) (*oa.UploadFeedbackLogResponse, error) {
	return c.call(ctx, req)
}

func TestUploadFeedbackLogRPCKeepsOriginalErrorMapping(t *testing.T) {
	original := oaClient
	t.Cleanup(func() { oaClient = original })
	impl := oaimpl.NewOAService(&base.ClientSet{})
	oaClient = feedbackLogRPCClient{call: impl.UploadFeedbackLog}
	for _, tc := range []struct {
		name string
		file []byte
	}{
		{name: "empty"},
		{name: "oversized", file: make([]byte, constants.FeedbackLogMaxSize+1)},
		{name: "invalid JSON", file: []byte("{}\ninvalid")},
		{name: "invalid UTF8", file: []byte{0xff}},
		{name: "storage not initialized", file: []byte("{}")},
	} {
		t.Run(tc.name, func(t *testing.T) {
			url, err := UploadFeedbackLogRPC(context.Background(), &oa.UploadFeedbackLogRequest{File: tc.file})
			require.Error(t, err)
			require.Equal(t, int64(errno.BizErrorCode), errno.ConvertErr(err).ErrorCode)
			require.Empty(t, url)
		})
	}
}

func TestUploadFeedbackLogRPCResponses(t *testing.T) {
	original := oaClient
	t.Cleanup(func() { oaClient = original })
	for _, tc := range []struct {
		name string
		resp *oa.UploadFeedbackLogResponse
		err  error
		code int64
	}{
		{name: "transport error", err: errors.New("RPC unavailable"), code: errno.InternalServiceErrorCode},
		{
			name: "infrastructure error", code: errno.BizErrorCode,
			resp: &oa.UploadFeedbackLogResponse{Base: &model.BaseResp{Code: errno.InternalSFErrorCode, Msg: "snowflake failed"}},
		},
		{
			name: "success", code: errno.SuccessCode,
			resp: &oa.UploadFeedbackLogResponse{Base: base.BuildSuccessResp(), Url: "https://example.com/log.json.gz"},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			oaClient = feedbackLogRPCClient{
				call: func(context.Context, *oa.UploadFeedbackLogRequest) (*oa.UploadFeedbackLogResponse, error) {
					return tc.resp, tc.err
				},
			}
			url, err := UploadFeedbackLogRPC(context.Background(), &oa.UploadFeedbackLogRequest{})
			require.Equal(t, tc.code, errno.ConvertErr(err).ErrorCode)
			if tc.code == errno.SuccessCode {
				require.NoError(t, err)
				require.Equal(t, tc.resp.Url, url)
			} else {
				require.Error(t, err)
				require.Empty(t, url)
			}
		})
	}
}
