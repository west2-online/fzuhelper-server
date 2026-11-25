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

package cos

import (
	"hash/crc64"
	"io"
	"net/http"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	tencentyun "github.com/tencentyun/cos-go-sdk-v5"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

type feedbackRoundTripFunc func(*http.Request) (*http.Response, error)

func (f feedbackRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestFeedbackCOSUpload(t *testing.T) {
	for _, tc := range []struct {
		name       string
		remotePath string
		body       []byte
		status     int
	}{
		{"screenshot", "/feedback/img/1.png", []byte("image"), http.StatusOK},
		{"log", "/feedback/log/2.json.gz", []byte{0x1f, 0x8b, 0x08}, http.StatusOK},
		{"key without slash", "feedback/img/3.jpg", []byte("image"), http.StatusOK},
		{"access denied", "/feedback/img/4.png", []byte("image"), http.StatusForbidden},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			httpClient := &http.Client{Transport: feedbackRoundTripFunc(func(req *http.Request) (*http.Response, error) {
				calls++
				require.Equal(t, http.MethodPut, req.Method)
				require.Equal(t, "/"+strings.TrimPrefix(tc.remotePath, "/"), req.URL.Path)
				body, err := io.ReadAll(req.Body)
				require.NoError(t, err)
				require.Equal(t, tc.body, body)
				header := make(http.Header)
				checksum := crc64.Checksum(body, crc64.MakeTable(crc64.ECMA))
				header.Set("x-cos-hash-crc64ecma", strconv.FormatUint(checksum, 10))
				responseBody := ""
				if tc.status != http.StatusOK {
					responseBody = "<Error><Code>AccessDenied</Code><Message>Access denied</Message></Error>"
				}
				return &http.Response{
					StatusCode: tc.status,
					Header:     header,
					Body:       io.NopCloser(strings.NewReader(responseBody)),
					Request:    req,
				}, nil
			})}
			bucketURL, err := tencentyun.NewBucketURL("feedback-test-1250000000", "ap-guangzhou", true)
			require.NoError(t, err)
			client := tencentyun.NewClient(&tencentyun.BaseURL{BucketURL: bucketURL}, httpClient)
			repo := NewFeedbackCOSCli(client, "/feedback/", "https://feedback.example.com", nil)
			err = repo.Upload(tc.body, tc.remotePath)
			if tc.status == http.StatusOK {
				require.NoError(t, err)
			} else {
				require.ErrorIs(t, err, errno.UpcloudError)
			}
			require.Equal(t, 1, calls)
		})
	}
}

func TestFeedbackCOSGenerateFileName(t *testing.T) {
	sf, err := utils.NewSnowflake(0, 0)
	require.NoError(t, err)
	for _, rootPath := range []string{"/feedback/", "feedback", "/feedback"} {
		t.Run(rootPath, func(t *testing.T) {
			repo := NewFeedbackCOSCli(nil, rootPath, "https://feedback.example.com/", sf)
			imageURL, imagePath, err := repo.GenerateFileName(FeedbackImageCategory, "png")
			require.NoError(t, err)
			require.Regexp(t, "^/feedback/img/[0-9]+[.]png$", imagePath)
			require.Equal(t, "https://feedback.example.com"+imagePath, imageURL)

			logURL, logPath, err := repo.GenerateFileName(FeedbackLogCategory, FeedbackLogFileExtension)
			require.NoError(t, err)
			require.Regexp(t, "^/feedback/log/[0-9]+[.]json[.]gz$", logPath)
			require.Equal(t, "https://feedback.example.com"+logPath, logURL)

			_, nextPath, err := repo.GenerateFileName(FeedbackImageCategory, "png")
			require.NoError(t, err)
			require.NotEqual(t, imagePath, nextPath)
		})
	}
}
