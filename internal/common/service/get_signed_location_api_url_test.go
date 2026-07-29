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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
)

func TestGetSignedApiUrl(t *testing.T) {
	type testCase struct {
		name                string
		location            string
		enabled             bool
		disableMsg          string
		endpoint            string
		mockDoError         error
		mockRespBody        string
		expectRequestBody   string
		expectSignedURL     string
		expectHeaders       map[string]string
		expectErrorContains string
	}

	const (
		location = "119.262647,26.106131"
		endpoint = "http://location-service.test/sign"
	)

	testCases := []testCase{
		{
			name:       "success",
			location:   location,
			enabled:    true,
			disableMsg: "should not be returned",
			endpoint:   endpoint,
			mockRespBody: `{
				"data": {
					"signed_url": "https://restapi.amap.com/v3/place/around?key=xxx&scode=abc",
					"headers": {"User-Agent": "AMAP_Location_SDK_Android"}
				},
				"base": {"code": 10000, "msg": "ok"}
			}`,
			expectRequestBody: `{"location":"119.262647,26.106131"}`,
			expectSignedURL:   "https://restapi.amap.com/v3/place/around?key=xxx&scode=abc",
			expectHeaders:     map[string]string{"User-Agent": "AMAP_Location_SDK_Android"},
		},
		{
			name:                "service_disabled_returns_configured_message",
			location:            location,
			enabled:             false,
			disableMsg:          "location service is under maintenance",
			endpoint:            endpoint,
			expectErrorContains: "location service is under maintenance",
		},
		{
			name:                "empty_location",
			location:            "",
			enabled:             true,
			disableMsg:          "should not be returned",
			endpoint:            endpoint,
			expectErrorContains: "location is empty",
		},
		{
			name:       "business_error",
			location:   location,
			enabled:    true,
			disableMsg: "should not be returned",
			endpoint:   endpoint,
			mockRespBody: `{
				"data": {
					"signed_url": "https://should-not-be-returned.example.com",
					"headers": {}
				},
				"base": {"code": 50001, "msg": "signing failed"}
			}`,
			expectRequestBody:   `{"location":"119.262647,26.106131"}`,
			expectErrorContains: "[50001] signing failed",
		},
		{
			name:                "request_error",
			location:            location,
			enabled:             true,
			disableMsg:          "should not be returned",
			endpoint:            endpoint,
			mockDoError:         assert.AnError,
			expectRequestBody:   `{"location":"119.262647,26.106131"}`,
			expectErrorContains: "request service failed assert.AnError general error for testing",
		},
		{
			name:                "invalid_response_body",
			location:            location,
			enabled:             true,
			disableMsg:          "should not be returned",
			endpoint:            endpoint,
			mockRespBody:        `invalid json`,
			expectRequestBody:   `{"location":"119.262647,26.106131"}`,
			expectErrorContains: "unmarshal response failed",
		},
		{
			name:                "nil_base",
			location:            location,
			enabled:             true,
			disableMsg:          "should not be returned",
			endpoint:            endpoint,
			mockRespBody:        `{"data": null}`,
			expectRequestBody:   `{"location":"119.262647,26.106131"}`,
			expectErrorContains: "response base is nil",
		},
		{
			name:                "nil_data",
			location:            location,
			enabled:             true,
			disableMsg:          "should not be returned",
			endpoint:            endpoint,
			mockRespBody:        `{"data": null, "base": {"code": 10000, "msg": "ok"}}`,
			expectRequestBody:   `{"location":"119.262647,26.106131"}`,
			expectErrorContains: "SignedUrlData is nil",
		},
	}

	httpClient, _ := client.NewClient()
	requestCtx := context.Background()
	_ = config.InitForTest("common")

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			config.SignedLocationApiUrl.Enabled = tc.enabled
			config.SignedLocationApiUrl.DisableMsg = tc.disableMsg
			config.SignedLocationApiUrl.Endpoint = tc.endpoint

			mockey.Mock((*client.Client).Do).To(
				func(c *client.Client, ctx context.Context, req *protocol.Request, resp *protocol.Response) error {
					if tc.mockDoError != nil {
						return tc.mockDoError
					}
					resp.SetBodyString(tc.mockRespBody)
					return nil
				},
			).Build()

			commonService := NewCommonService(requestCtx, &base.ClientSet{HzClient: httpClient}, new(taskqueue.BaseTaskQueue))
			signedURL, headers, err := commonService.GetSignedApiUrl(tc.location)

			if tc.expectErrorContains != "" {
				assert.ErrorContains(t, err, tc.expectErrorContains)
				assert.Empty(t, signedURL)
				assert.Nil(t, headers)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectSignedURL, signedURL)
				assert.Equal(t, tc.expectHeaders, headers)
			}
		})
	}
}
