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
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/oss"
)

type logOSSRepo struct {
	category      string
	suffix        string
	uploaded      []byte
	generateError error
	uploadError   error
}

func (r *logOSSRepo) GenerateFileName(category, suffix string) (string, string, error) {
	r.category = category
	r.suffix = suffix
	return "https://log.example.com/feedback/log/1.json.gz", "/feedback/log/1.json.gz", r.generateError
}

func (r *logOSSRepo) Upload(file []byte, _ string) error {
	r.uploaded = append([]byte(nil), file...)
	return r.uploadError
}

func TestUploadFeedbackLog(t *testing.T) {
	validLog := []byte("{\"event\":\"launch\"}\n{\"event\":\"request\",\"status\":200}\n")

	type testCase struct {
		name          string
		file          []byte
		generateError error
		uploadError   error
		expectError   string
	}
	testCases := []testCase{
		{name: "upload log", file: validLog},
		{name: "empty file", expectError: "反馈日志不能为空"},
		{name: "file too large", file: make([]byte, constants.FeedbackLogMaxSize+1), expectError: "不能超过 2MB"},
		{name: "invalid utf8", file: []byte{0xff, 0xfe}, expectError: "UTF-8 JSON Lines"},
		{name: "invalid json line", file: []byte("{\"event\":1}\ninvalid\n"), expectError: "UTF-8 JSON Lines"},
		{name: "generate name error", file: validLog, generateError: errno.InternalServiceError, expectError: "生成反馈日志名称失败"},
		{name: "upload error", file: validLog, uploadError: errno.UpcloudError, expectError: "上传反馈日志失败"},
	}

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			repo := &logOSSRepo{generateError: tc.generateError, uploadError: tc.uploadError}
			service := &OAService{ctx: context.Background(), ossClient: repo}

			url, err := service.UploadFeedbackLog(tc.file)
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
				assert.Empty(t, url)
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, url)
			assert.Equal(t, oss.FeedbackLogCategory, repo.category)
			assert.Equal(t, oss.FeedbackLogFileExtension, repo.suffix)

			reader, err := gzip.NewReader(bytes.NewReader(repo.uploaded))
			require.NoError(t, err)
			decompressed, err := io.ReadAll(reader)
			require.NoError(t, err)
			require.NoError(t, reader.Close())
			assert.Equal(t, tc.file, decompressed)
		})
	}

	mockey.PatchConvey("oss not initialized", t, func() {
		service := &OAService{ctx: context.Background()}
		_, err := service.UploadFeedbackLog(validLog)
		assert.ErrorContains(t, err, "反馈文件存储未初始化")
	})
}

func TestRedactFeedbackLog(t *testing.T) {
	entry := map[string]interface{}{
		"headers": map[string]string{
			"Authorization": "Bearer abc.def.ghi",
			"Cookies":       "session=value",
			"Id":            "052106112",
		},
		"password": "password",
		"message": "phones 13800138000 13900139000, stu_id=102301000 student_id=102301001, " +
			"credential Bearer short-token, jwt eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxIn0.signature, " +
			"token=abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
		"phone":      json.Number("13700137000"),
		"session_id": "session",
	}
	logData, err := json.Marshal(entry)
	require.NoError(t, err)
	logData = append(logData, '\n')

	redacted, err := redactedFeedbackLogBytes(logData)
	require.NoError(t, err)
	assert.NotContains(t, string(redacted), "Bearer abc.def.ghi")
	assert.NotContains(t, string(redacted), "session=value")
	assert.NotContains(t, string(redacted), "052106112")
	assert.NotContains(t, string(redacted), "13800138000")
	assert.NotContains(t, string(redacted), "13900139000")
	assert.NotContains(t, string(redacted), "13700137000")
	assert.NotContains(t, string(redacted), "102301000")
	assert.NotContains(t, string(redacted), "102301001")
	assert.NotContains(t, string(redacted), "short-token")
	assert.NotContains(t, string(redacted), "eyJhbGciOiJIUzI1NiJ9")
	assert.NotContains(t, string(redacted), "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	assert.Contains(t, string(redacted), `"session_id":"session"`)
	assert.Contains(t, string(redacted), feedbackLogRedactedValue)
}

func TestUploadFeedbackLogRedactsBeforeStorage(t *testing.T) {
	repo := &logOSSRepo{}
	service := &OAService{ctx: context.Background(), ossClient: repo}
	logData := []byte(
		`{"timestamp":1788912000,"message":"password=hunter2&token=short-token",` +
			`"headers":{"Set-Cookie":"session=short-value"}}`,
	)
	_, err := service.UploadFeedbackLog(logData)
	require.NoError(t, err)
	reader, err := gzip.NewReader(bytes.NewReader(repo.uploaded))
	require.NoError(t, err)
	defer reader.Close()
	content, err := io.ReadAll(reader)
	require.NoError(t, err)
	require.JSONEq(t, `{"timestamp":1788912000,"message":"password=***&token=***","headers":{"Set-Cookie":"***"}}`, string(content))

	repo = &logOSSRepo{}
	service.ossClient = repo
	_, err = service.UploadFeedbackLog([]byte("{}\ninvalid"))
	require.Error(t, err)
	require.Empty(t, repo.category, "invalid input must not generate an OSS name")
	require.Nil(t, repo.uploaded, "invalid input must not upload partial content")
}

func TestUploadFeedbackLogValidationKeepsOriginalErrors(t *testing.T) {
	for _, tc := range []struct {
		name    string
		file    []byte
		message string
	}{
		{name: "empty", message: "service.UploadFeedbackLog: 反馈日志不能为空"},
		{name: "oversized", file: make([]byte, constants.FeedbackLogMaxSize+1), message: "service.UploadFeedbackLog: 反馈日志不能超过 2MB"},
		{
			name: "blank", file: []byte(" \r\n"),
			message: "service.UploadFeedbackLog: 反馈日志必须是 UTF-8 JSON Lines: feedback log contains no JSON values",
		},
		{
			name: "invalid UTF8", file: []byte{0xff},
			message: "service.UploadFeedbackLog: 反馈日志必须是 UTF-8 JSON Lines: feedback log must be UTF-8 JSON Lines",
		},
		{
			name: "multiple values", file: []byte("{} {}"),
			message: "service.UploadFeedbackLog: 反馈日志必须是 UTF-8 JSON Lines: invalid JSON at line 1",
		},
		{
			name: "bad final line", file: []byte(strings.Repeat("{}\n", 1000) + "invalid"),
			message: "service.UploadFeedbackLog: 反馈日志必须是 UTF-8 JSON Lines: invalid JSON at line 1001",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			repo := &logOSSRepo{}
			s := &OAService{ossClient: repo}
			url, err := s.UploadFeedbackLog(tc.file)
			require.EqualError(t, err, tc.message)
			require.Equal(t, int64(errno.InternalServiceErrorCode), errno.ConvertErr(err).ErrorCode)
			require.Empty(t, url)
			require.Empty(t, repo.category)
			require.Nil(t, repo.uploaded)
		})
	}
}

type feedbackCloseFailWriter struct {
	err error
}

func (w feedbackCloseFailWriter) Write(data []byte) (int, error) {
	// gzip 首次写入的 10 字节头部成功，小文件的压缩数据在 Close 时输出并失败。
	if len(data) == 10 && data[0] == 0x1f && data[1] == 0x8b {
		return len(data), nil
	}
	return 0, w.err
}

func TestCompressFeedbackLogWriterFailures(t *testing.T) {
	want := errors.New("storage writer failed")
	data := []byte(`{"ts":1788912000}`)
	err := compressFeedbackLog(data, feedbackFailWriter{err: want}, gzip.DefaultCompression)
	require.ErrorIs(t, err, want)
	err = compressFeedbackLog(data, feedbackCloseFailWriter{err: want}, gzip.DefaultCompression)
	require.ErrorIs(t, err, want)
	require.ErrorContains(t, err, "service.UploadFeedbackLog: 压缩反馈日志失败")

	mockey.PatchConvey("compression failure must not upload", t, func() {
		mockey.Mock(compressFeedbackLog).Return(want).Build()
		repo := &logOSSRepo{}
		s := &OAService{ossClient: repo}
		url, err := s.UploadFeedbackLog(data)
		require.ErrorIs(t, err, want)
		require.Equal(t, int64(errno.InternalServiceErrorCode), errno.ConvertErr(err).ErrorCode)
		require.Empty(t, url)
		require.Empty(t, repo.category)
		require.Nil(t, repo.uploaded)
	})
}

func TestCompressFeedbackLogMatchesBufferedOutput(t *testing.T) {
	for _, data := range [][]byte{
		[]byte("\r\n{\"token\":\"secret\",\"ts\":1788912000}\r\n{\"duration\":42}"),
		[]byte(`{"message":"` + strings.Repeat("x", 70*1024) + `"}`),
		append(bytes.Repeat([]byte("{}\n"), constants.FeedbackLogMaxSize/3), ' ', ' '),
	} {
		want, err := redactedFeedbackLogBytes(data)
		require.NoError(t, err)
		for _, level := range []int{gzip.DefaultCompression, gzip.BestSpeed} {
			var dst bytes.Buffer
			require.NoError(t, compressFeedbackLog(data, &dst, level))
			reader, err := gzip.NewReader(&dst)
			require.NoError(t, err)
			got, err := io.ReadAll(reader)
			require.NoError(t, err, "read through gzip footer to validate checksum")
			require.NoError(t, reader.Close())
			require.Equal(t, want, got)
		}
	}
}
