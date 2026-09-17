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

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"mime/multipart"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	oa "github.com/west2-online/fzuhelper-server/kitex_gen/oa"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

const feedbackStuID = "102301000"

func feedbackTestRouter(method, path string, handler app.HandlerFunc) *route.Engine {
	router := route.NewEngine(&config.Options{})
	setStuID := func(ctx context.Context, c *app.RequestContext) {
		c.Set(constants.StuIDContextKey, feedbackStuID)
		c.Next(ctx)
	}
	if method == consts.MethodGet {
		router.GET(path, setStuID, handler)
	} else {
		router.POST(path, setStuID, handler)
	}
	return router
}

func feedbackMultipartBody(t *testing.T, data []byte) (*ut.Body, ut.Header) {
	t.Helper()
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	part, err := writer.CreateFormFile("file", "feedback")
	require.NoError(t, err)
	_, err = part.Write(data)
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	return &ut.Body{Body: bytes.NewReader(buffer.Bytes()), Len: buffer.Len()},
		ut.Header{Key: "Content-Type", Value: writer.FormDataContentType()}
}

const createFeedbackBody = `{
"name":"张三",
"college":"计算机学院",
"contact_phone":"13800000000",
"contact_qq":"10001",
"contact_email":"a@b.com",
"network_env":"wifi",
"is_on_campus":true,"os_name":
"Android","os_version":"14",
"manufacturer":"Xiaomi","device_model":
"Mi 14",
"problem_desc":"登录白屏",
"app_version":"1.2.3"
}`

func TestCreateFeedback(t *testing.T) {
	defer mockey.UnPatchAll()
	var rpcStuID string
	mockey.Mock(rpc.CreateFeedbackRPC).To(func(_ context.Context, req *oa.CreateFeedbackRequest) (int64, error) {
		rpcStuID = req.StuId
		return 123, nil
	}).Build()

	router := feedbackTestRouter(consts.MethodPost, "/api/v1/feedback/create", CreateFeedback)
	res := ut.PerformRequest(router, consts.MethodPost, "/api/v1/feedback/create",
		&ut.Body{Body: bytes.NewBufferString(createFeedbackBody), Len: len(createFeedbackBody)},
		ut.Header{Key: "Content-Type", Value: "application/json"})

	assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
	assert.Contains(t, string(res.Result().Body()), `"report_id":123`)
	assert.Equal(t, feedbackStuID, rpcStuID)
}

func TestCreateFeedbackSingleContact(t *testing.T) {
	testCases := []struct {
		field, value string
	}{
		{field: "contact_phone", value: "13800000000"},
		{field: "contact_qq", value: "10001"},
		{field: "contact_email", value: "a@b.com"},
	}
	for _, tc := range testCases {
		mockey.PatchConvey(tc.field, t, func() {
			var payload map[string]interface{}
			require.NoError(t, json.Unmarshal([]byte(createFeedbackBody), &payload))
			delete(payload, "contact_phone")
			delete(payload, "contact_qq")
			delete(payload, "contact_email")
			payload[tc.field] = tc.value
			body, err := json.Marshal(payload)
			require.NoError(t, err)

			called := false
			mockey.Mock(rpc.CreateFeedbackRPC).To(func(_ context.Context, req *oa.CreateFeedbackRequest) (int64, error) {
				called = true
				contacts := map[string]string{
					"contact_phone": req.GetContactPhone(),
					"contact_qq":    req.GetContactQq(),
					"contact_email": req.GetContactEmail(),
				}
				assert.Equal(t, tc.value, contacts[tc.field])
				delete(contacts, tc.field)
				for _, value := range contacts {
					assert.Empty(t, value)
				}
				assert.Equal(t, feedbackStuID, req.StuId)
				return 123, nil
			}).Build()

			router := feedbackTestRouter(consts.MethodPost, "/api/v1/feedback/create", CreateFeedback)
			res := ut.PerformRequest(router, consts.MethodPost, "/api/v1/feedback/create",
				&ut.Body{Body: bytes.NewReader(body), Len: len(body)},
				ut.Header{Key: "Content-Type", Value: "application/json"})
			assert.True(t, called)
			assert.Contains(t, string(res.Result().Body()), `"report_id":123`)
		})
	}
}

func TestCreateFeedbackRPCError(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.Mock(rpc.CreateFeedbackRPC).Return(int64(0), errno.InternalServiceError).Build()

	router := feedbackTestRouter(consts.MethodPost, "/api/v1/feedback/create", CreateFeedback)
	res := ut.PerformRequest(router, consts.MethodPost, "/api/v1/feedback/create",
		&ut.Body{Body: bytes.NewBufferString(createFeedbackBody), Len: len(createFeedbackBody)},
		ut.Header{Key: "Content-Type", Value: "application/json"})

	assert.Contains(t, string(res.Result().Body()), `"code":"50001"`)
}

func TestGetFeedbackByID(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.Mock(rpc.GetFeedbackByIDRPC).To(func(_ context.Context, req *oa.GetFeedbackByIDRequest) (*model.Feedback, error) {
		assert.Equal(t, feedbackStuID, req.StuId)
		assert.Equal(t, int64(123), req.ReportId)
		return &model.Feedback{ReportId: 123, StuId: feedbackStuID, ProblemDesc: "登录白屏"}, nil
	}).Build()

	router := feedbackTestRouter(consts.MethodGet, "/api/v1/feedbacks/detail", GetFeedbackByID)
	res := ut.PerformRequest(router, consts.MethodGet, "/api/v1/feedbacks/detail?report_id=123", nil)

	assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
	assert.Contains(t, string(res.Result().Body()), `"report_id":123`)
}

func TestListFeedback(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.Mock(rpc.ListFeedbackRPC).To(func(_ context.Context, req *oa.GetListFeedbackRequest) ([]*model.FeedbackListItem, int64, error) {
		assert.Equal(t, feedbackStuID, req.StuId)
		return []*model.FeedbackListItem{{ReportId: 123, ProblemDesc: "登录白屏"}}, 100, nil
	}).Build()

	router := feedbackTestRouter(consts.MethodGet, "/api/v1/feedbacks/list", ListFeedback)
	res := ut.PerformRequest(router, consts.MethodGet, "/api/v1/feedbacks/list?limit=20", nil)

	assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
	assert.Contains(t, string(res.Result().Body()), `"page_token":100`)
}

func TestUploadFeedbackScreenshot(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.Mock(rpc.UploadFeedbackScreenshotRPC).Return("https://example.com/feedback/img/1.jpg", nil).Build()

	jpeg := []byte{0xff, 0xd8, 0xff, 0xe0, 0x00, 0x10, 'J', 'F', 'I', 'F', 0x00}
	body, header := feedbackMultipartBody(t, jpeg)
	router := feedbackTestRouter(consts.MethodPost, "/api/v1/feedback/upload", UploadFeedbackScreenshot)
	res := ut.PerformRequest(router, consts.MethodPost, "/api/v1/feedback/upload", body, header)

	assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
	assert.Contains(t, string(res.Result().Body()), `"url":"https://example.com/feedback/img/1.jpg"`)
}

func TestUploadFeedbackScreenshotRejectsInvalidFile(t *testing.T) {
	body, header := feedbackMultipartBody(t, []byte("not an image"))
	router := feedbackTestRouter(consts.MethodPost, "/api/v1/feedback/upload", UploadFeedbackScreenshot)
	res := ut.PerformRequest(router, consts.MethodPost, "/api/v1/feedback/upload", body, header)

	assert.Contains(t, string(res.Result().Body()), "反馈截图仅支持 JPEG 和 PNG")
}

func TestUploadFeedbackLog(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.Mock(rpc.UploadFeedbackLogRPC).Return("https://example.com/feedback/log/1.json.gz", nil).Build()

	body, header := feedbackMultipartBody(t, []byte("{\"event\":\"launch\"}\n"))
	router := feedbackTestRouter(consts.MethodPost, "/api/v1/feedback/upload-log", UploadFeedbackLog)
	res := ut.PerformRequest(router, consts.MethodPost, "/api/v1/feedback/upload-log", body, header)

	assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
	assert.Contains(t, string(res.Result().Body()), `"url":"https://example.com/feedback/log/1.json.gz"`)
}

func TestUploadFeedbackLogRejectsInvalidJSONLines(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.Mock(rpc.UploadFeedbackLogRPC).To(func(_ context.Context, _ *oa.UploadFeedbackLogRequest) (string, error) {
		t.Error("invalid JSON Lines must be rejected before RPC")
		return "", nil
	}).Build()
	body, header := feedbackMultipartBody(t, []byte("invalid json lines"))
	router := feedbackTestRouter(consts.MethodPost, "/api/v1/feedback/upload-log", UploadFeedbackLog)
	res := ut.PerformRequest(router, consts.MethodPost, "/api/v1/feedback/upload-log", body, header)

	assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
	assert.Contains(t, string(res.Result().Body()), `"code":"20009"`)
	assert.Contains(t, string(res.Result().Body()), "反馈日志必须是 UTF-8 JSON Lines")
}

func TestUploadFeedbackLogRejectsSizeBeforeRPC(t *testing.T) {
	defer mockey.UnPatchAll()
	mockey.Mock(rpc.UploadFeedbackLogRPC).To(func(_ context.Context, _ *oa.UploadFeedbackLogRequest) (string, error) {
		t.Error("empty or oversized log must not reach RPC")
		return "", nil
	}).Build()
	for _, tc := range []struct {
		name string
		file []byte
		code string
	}{
		{name: "empty", code: "20012"},
		{name: "oversized", file: make([]byte, constants.FeedbackLogMaxSize+1), code: "20010"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body, header := feedbackMultipartBody(t, tc.file)
			router := feedbackTestRouter(consts.MethodPost, "/api/v1/feedback/upload-log", UploadFeedbackLog)
			res := ut.PerformRequest(router, consts.MethodPost, "/api/v1/feedback/upload-log", body, header)
			assert.Contains(t, string(res.Result().Body()), `"code":"`+tc.code+`"`)
		})
	}
}
