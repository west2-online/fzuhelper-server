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
	"encoding/base64"
	"errors"
	"io"
	"mime/multipart"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/cloudwego/hertz/pkg/route"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

// 测试数据常量
const (
	testImageBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg=="
	testImageName   = "test.png"
)

// buildCreateImageForm 构建创建图片的 form 数据
func buildCreateImageForm() (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("pic_type", "1")
	_ = w.WriteField("start_at", "1609459200")
	_ = w.WriteField("end_at", "1609545600")
	_ = w.WriteField("s_type", "1")
	_ = w.WriteField("frequency", "1")
	_ = w.WriteField("start_time", "6")
	_ = w.WriteField("end_time", "18")
	_ = w.WriteField("text", "test")
	_ = w.WriteField("regex", ".*")
	// Create a fake image file
	part, _ := w.CreateFormFile("image", testImageName)
	imageData, _ := base64.StdEncoding.DecodeString(testImageBase64)
	_, _ = io.Copy(part, bytes.NewReader(imageData))
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

// buildChangeImageForm 构建修改图片的 form 数据
func buildChangeImageForm() (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("picture_id", "1")
	// Create a fake image file
	part, _ := w.CreateFormFile("image", testImageName)
	imageData, _ := base64.StdEncoding.DecodeString(testImageBase64)
	_, _ = io.Copy(part, bytes.NewReader(imageData))
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

// buildEmptyForm 构建空的 form 数据（用于错误测试）
func buildEmptyForm() (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

// buildChangeImageFormWithoutImage 构建不含图片文件的 form（用于错误测试）
func buildChangeImageFormWithoutImage() (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("picture_id", "1")
	_ = w.Close()
	return &buf, w.FormDataContentType()
}

func TestCreateImage(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *model.Picture
		mockErr        error
		expectContains string
		buildForm      func() (*bytes.Buffer, string)
	}

	testCases := []testCase{
		{
			name: "success",
			url:  "/api/v1/launch_screen/api/image",
			mockResp: &model.Picture{
				Id:      1,
				Type:    1,
				StartAt: 1609459200,
				EndAt:   1609545600,
			},
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
			buildForm:      buildCreateImageForm,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image",
			mockResp:       nil,
			mockErr:        errors.New("service error"),
			expectContains: `{"code":"50001","message":`,
			buildForm:      buildCreateImageForm,
		},
		{
			name:           "bind error - missing required params",
			url:            "/api/v1/launch_screen/api/image",
			mockResp:       nil,
			mockErr:        nil,
			expectContains: `{"code":"20001","message":`,
			buildForm:      buildEmptyForm,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/launch_screen/api/image", CreateImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.CreateImageRPC).To(func(ctx context.Context, req *launch_screen.CreateImageRequest, file [][]byte) (*model.Picture, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			buf, contentType := tc.buildForm()
			result := ut.PerformRequest(router, consts.MethodPost, tc.url,
				&ut.Body{Body: buf, Len: buf.Len()},
				ut.Header{Key: "Content-Type", Value: contentType})
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestGetImage(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *model.Picture
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/launch_screen/api/image?picture_id=1",
			mockResp:       &model.Picture{},
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image?picture_id=1",
			mockResp:       nil,
			mockErr:        errors.New("service error"),
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/launch_screen/api/image",
			mockResp:       nil,
			mockErr:        nil,
			expectContains: `{"code":"20001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/launch_screen/api/image", GetImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetImageRPC).To(func(ctx context.Context, req *launch_screen.GetImageRequest) (*model.Picture, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestChangeImageProperty(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *model.Picture
		mockErr        error
		expectContains string
	}
	//nolint:lll
	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/launch_screen/api/image?picture_id=1&pic_type=1&start_at=1609459200&end_at=1609545600&s_type=1&frequency=1&start_time=6&end_time=18&text=test&regex=",
			mockResp:       &model.Picture{},
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image?picture_id=1&pic_type=1&start_at=1609459200&end_at=1609545600&s_type=1&frequency=1&start_time=6&end_time=18&text=test&regex=",
			mockResp:       nil,
			mockErr:        errors.New("service error"),
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/launch_screen/api/image",
			mockResp:       nil,
			mockErr:        nil,
			expectContains: `{"code":"20001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.PUT("/api/v1/launch_screen/api/image", ChangeImageProperty)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.ChangeImagePropertyRPC).To(func(ctx context.Context, req *launch_screen.ChangeImagePropertyRequest) (*model.Picture, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPut, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestChangeImage(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *model.Picture
		mockErr        error
		expectContains string
		buildForm      func() (*bytes.Buffer, string)
	}

	testCases := []testCase{
		{
			name: "success",
			url:  "/api/v1/launch_screen/api/image/img",
			mockResp: &model.Picture{
				Id:      1,
				Type:    1,
				StartAt: 1609459200,
				EndAt:   1609545600,
			},
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
			buildForm:      buildChangeImageForm,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image/img",
			mockResp:       nil,
			mockErr:        errors.New("service error"),
			expectContains: `{"code":"50001","message":`,
			buildForm:      buildChangeImageForm,
		},
		{
			name:           "bind error - missing image file",
			url:            "/api/v1/launch_screen/api/image/img",
			mockResp:       nil,
			mockErr:        nil,
			expectContains: `{"code":"20001","message":`,
			buildForm:      buildChangeImageFormWithoutImage,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.PUT("/api/v1/launch_screen/api/image/img", ChangeImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.ChangeImageRPC).To(func(ctx context.Context, req *launch_screen.ChangeImageRequest, file [][]byte) (*model.Picture, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			buf, contentType := tc.buildForm()
			result := ut.PerformRequest(router, consts.MethodPut, tc.url,
				&ut.Body{Body: buf, Len: buf.Len()},
				ut.Header{Key: "Content-Type", Value: contentType})
			assert.Equal(t, consts.StatusOK, result.Result().StatusCode())
			assert.Contains(t, string(result.Result().Body()), tc.expectContains)
		})
	}
}

func TestDeleteImage(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/launch_screen/api/image?picture_id=1",
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image?picture_id=1",
			mockErr:        errors.New("service error"),
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/launch_screen/api/image",
			mockErr:        nil,
			expectContains: `{"code":"20001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.DELETE("/api/v1/launch_screen/api/image", DeleteImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.DeleteImageRPC).To(func(ctx context.Context, req *launch_screen.DeleteImageRequest) error {
				return tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodDelete, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestMobileGetImage(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       []*model.Picture
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/launch_screen/api/screen?type=1&student_id=202400001&device=ios",
			mockResp:       []*model.Picture{{}},
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/screen?type=1&student_id=202400001&device=ios",
			mockResp:       nil,
			mockErr:        errors.New("service error"),
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/launch_screen/api/screen",
			mockResp:       nil,
			mockErr:        nil,
			expectContains: `{"code":"20001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/launch_screen/api/screen", MobileGetImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.MobileGetImageRPC).To(func(ctx context.Context, req *launch_screen.MobileGetImageRequest) ([]*model.Picture, *int64, error) {
				count := int64(len(tc.mockResp))
				return tc.mockResp, &count, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestAddImagePointTime(t *testing.T) {
	type testCase struct {
		name           string
		url            string
		mockResp       *model.Picture
		mockErr        error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/launch_screen/api/image/point?picture_id=1",
			mockResp:       &model.Picture{},
			mockErr:        nil,
			expectContains: `{"code":"10000","message":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image/point?picture_id=1",
			mockResp:       nil,
			mockErr:        errors.New("service error"),
			expectContains: `{"code":"50001","message":`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/launch_screen/api/image/point",
			mockResp:       nil,
			mockErr:        nil,
			expectContains: `{"code":"20001","message":`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/launch_screen/api/image/point", AddImagePointTime)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.AddImagePointTimeRPC).To(func(ctx context.Context, req *launch_screen.AddImagePointTimeRequest) (*model.Picture, error) {
				return tc.mockResp, tc.mockErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
