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
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
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

// buildCreateImageFormWithoutImage 构建缺少图片的 form 数据（用于错误测试）
func buildCreateImageFormWithoutImage() (*bytes.Buffer, string) {
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

// buildChangeImageFormWithoutImage 构建不含图片文件的 form（用于错误测试）
func buildChangeImageFormWithoutImage() (*bytes.Buffer, string) {
	var buf bytes.Buffer
	w := multipart.NewWriter(&buf)
	_ = w.WriteField("picture_id", "1")
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

func TestCreateImage(t *testing.T) {
	type testCase struct {
		name            string
		url             string
		mockResp        *model.Picture
		mockRPCErr      error
		expectContains  string
		buildForm       func() (*bytes.Buffer, string)
		mockTypeInvalid bool
		mockFileErr     error
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
			expectContains: `{"code":"10000","message":`,
			buildForm:      buildCreateImageForm,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
			buildForm:      buildCreateImageForm,
		},
		{
			name:           "bind error - missing required params",
			url:            "/api/v1/launch_screen/api/image",
			expectContains: `{"code":"20001","message":"参数错误,`,
			buildForm:      buildEmptyForm,
		},
		{
			name:           "form file error",
			url:            "/api/v1/launch_screen/api/image",
			expectContains: `{"code":"20001","message":"参数错误,`,
			buildForm:      buildCreateImageFormWithoutImage,
		},
		{
			name:            "invalid image suffix",
			url:             "/api/v1/launch_screen/api/image",
			expectContains:  `文件不可用`,
			buildForm:       buildCreateImageForm,
			mockTypeInvalid: true,
		},
		{
			name:           "file read error",
			url:            "/api/v1/launch_screen/api/image",
			expectContains: `{"code":"40001","message":"请求业务出现问题"}`,
			buildForm:      buildCreateImageForm,
			mockFileErr:    errno.BizError,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.POST("/api/v1/launch_screen/api/image", CreateImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.CreateImageRPC).To(func(ctx context.Context, req *launch_screen.CreateImageRequest, file [][]byte) (*model.Picture, error) {
				return tc.mockResp, tc.mockRPCErr
			}).Build()

			mockey.Mock(utils.CheckImageFileType).To(func(header *multipart.FileHeader) (string, bool) {
				return "", !tc.mockTypeInvalid
			}).Build()

			mockey.Mock(utils.FileToByteArray).To(func(file *multipart.FileHeader) ([][]byte, error) {
				return nil, tc.mockFileErr
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
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/launch_screen/api/image?picture_id=1",
			mockResp:       &model.Picture{},
			expectContains: `{"code":"10000","message":"Success","data":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image?picture_id=1",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/launch_screen/api/image",
			expectContains: `{"code":"20001","message":"参数错误,`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/launch_screen/api/image", GetImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.GetImageRPC).To(func(ctx context.Context, req *launch_screen.GetImageRequest) (*model.Picture, error) {
				return tc.mockResp, tc.mockRPCErr
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
		mockRPCErr     error
		expectContains string
	}
	//nolint:lll
	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/launch_screen/api/image?picture_id=1&pic_type=1&start_at=1609459200&end_at=1609545600&s_type=1&frequency=1&start_time=6&end_time=18&text=test&regex=",
			mockResp:       &model.Picture{},
			expectContains: `{"code":"10000","message":"Success","data":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image?picture_id=1&pic_type=1&start_at=1609459200&end_at=1609545600&s_type=1&frequency=1&start_time=6&end_time=18&text=test&regex=",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/launch_screen/api/image",
			expectContains: `{"code":"20001","message":"参数错误,`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.PUT("/api/v1/launch_screen/api/image", ChangeImageProperty)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.ChangeImagePropertyRPC).To(func(ctx context.Context, req *launch_screen.ChangeImagePropertyRequest) (*model.Picture, error) {
				return tc.mockResp, tc.mockRPCErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodPut, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}

func TestChangeImage(t *testing.T) {
	type testCase struct {
		name            string
		url             string
		mockResp        *model.Picture
		mockRPCErr      error
		expectContains  string
		buildForm       func() (*bytes.Buffer, string)
		mockTypeInvalid bool
		mockFileErr     error
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
			expectContains: `{"code":"10000","message":"Success","data":`,
			buildForm:      buildChangeImageForm,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image/img",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
			buildForm:      buildChangeImageForm,
		},
		{
			name:           "bind error - missing image file",
			url:            "/api/v1/launch_screen/api/image/img",
			expectContains: `{"code":"20001","message":"参数错误,`,
			buildForm:      buildEmptyForm,
		},
		{
			name:           "form file error",
			url:            "/api/v1/launch_screen/api/image/img",
			expectContains: `{"code":"20001","message":"参数错误,`,
			buildForm:      buildChangeImageFormWithoutImage,
		},
		{
			name:            "invalid image suffix",
			url:             "/api/v1/launch_screen/api/image/img",
			expectContains:  `文件不可用`,
			buildForm:       buildChangeImageForm,
			mockTypeInvalid: true,
		},
		{
			name:           "file read error",
			url:            "/api/v1/launch_screen/api/image/img",
			expectContains: `{"code":"40001","message":"请求业务出现问题"}`,
			buildForm:      buildChangeImageForm,
			mockFileErr:    errno.BizError,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.PUT("/api/v1/launch_screen/api/image/img", ChangeImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.ChangeImageRPC).To(func(ctx context.Context, req *launch_screen.ChangeImageRequest, file [][]byte) (*model.Picture, error) {
				return tc.mockResp, tc.mockRPCErr
			}).Build()

			mockey.Mock(utils.CheckImageFileType).To(func(header *multipart.FileHeader) (string, bool) {
				return "", !tc.mockTypeInvalid
			}).Build()

			mockey.Mock(utils.FileToByteArray).To(func(file *multipart.FileHeader) ([][]byte, error) {
				return nil, tc.mockFileErr
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
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/launch_screen/api/image?picture_id=1",
			expectContains: `{"code":"10000","message":"ok"}`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image?picture_id=1",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/launch_screen/api/image",
			expectContains: `{"code":"20001","message":"参数错误,`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.DELETE("/api/v1/launch_screen/api/image", DeleteImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.DeleteImageRPC).To(func(ctx context.Context, req *launch_screen.DeleteImageRequest) error {
				return tc.mockRPCErr
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
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/launch_screen/api/screen?type=1&student_id=202400001&device=ios",
			mockResp:       []*model.Picture{{}},
			expectContains: `{"code":"10000","message":"ok","data":`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/screen?type=1&student_id=202400001&device=ios",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/launch_screen/api/screen",
			expectContains: `{"code":"20001","message":"参数错误,`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/launch_screen/api/screen", MobileGetImage)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.MobileGetImageRPC).To(func(ctx context.Context, req *launch_screen.MobileGetImageRequest) ([]*model.Picture, *int64, error) {
				count := int64(len(tc.mockResp))
				return tc.mockResp, &count, tc.mockRPCErr
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
		mockRPCErr     error
		expectContains string
	}

	testCases := []testCase{
		{
			name:           "success",
			url:            "/api/v1/launch_screen/api/image/point?picture_id=1",
			mockResp:       &model.Picture{},
			expectContains: `{"code":"10000","message":"ok"}`,
		},
		{
			name:           "rpc error",
			url:            "/api/v1/launch_screen/api/image/point?picture_id=1",
			mockRPCErr:     errno.InternalServiceError,
			expectContains: `{"code":"50001","message":"内部服务错误"}`,
		},
		{
			name:           "bind error",
			url:            "/api/v1/launch_screen/api/image/point",
			expectContains: `{"code":"20001","message":"参数错误,`,
		},
	}

	router := route.NewEngine(&config.Options{})
	router.GET("/api/v1/launch_screen/api/image/point", AddImagePointTime)

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock(rpc.AddImagePointTimeRPC).To(func(ctx context.Context, req *launch_screen.AddImagePointTimeRequest) (*model.Picture, error) {
				return tc.mockResp, tc.mockRPCErr
			}).Build()

			res := ut.PerformRequest(router, consts.MethodGet, tc.url, nil)
			assert.Equal(t, consts.StatusOK, res.Result().StatusCode())
			assert.Contains(t, string(res.Result().Body()), tc.expectContains)
		})
	}
}
