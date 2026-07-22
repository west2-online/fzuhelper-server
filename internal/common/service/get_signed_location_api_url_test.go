package service

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
)

func TestGetSignedApiUrl(t *testing.T) {
	type testCase struct {
		name            string
		location        string
		enabled         bool
		disableMsg      string
		mockDoError     error
		mockRespBody    []byte
		mockStatusCode  int
		expectSignedURL string
		expectHeaders   map[string]string
		expectError     string
	}

	testCases := []testCase{
		{
			name:     "Success",
			location: "119.262647,26.106131",
			enabled:  true,
			mockRespBody: []byte(`{
				"data": {
					"signed_url": "https://restapi.amap.com/v3/place/around?key=xxx&scode=abc",
					"headers": {"User-Agent": "AMAP_Location_SDK_Android"}
				},
				"base": {"code": 0, "msg": "success"}
			}`),
			mockStatusCode:  200,
			expectSignedURL: "https://restapi.amap.com/v3/place/around?key=xxx&scode=abc",
			expectHeaders:   map[string]string{"User-Agent": "AMAP_Location_SDK_Android"},
		},
		{
			name:        "ServiceDisabled",
			location:    "119.262647,26.106131",
			enabled:     false,
			disableMsg:  "Service is unavailable",
			expectError: "Service is unavailable",
		},
		{
			name:        "EmptyLocation",
			location:    "",
			enabled:     true,
			expectError: "location is empty",
		},
		{
			name:           "UpstreamError",
			location:       "119.262647,26.106131",
			enabled:        true,
			mockStatusCode: 500,
			expectError:    "error response with status 500",
		},
		{
			name:           "InvalidResponseBody",
			location:       "119.262647,26.106131",
			enabled:        true,
			mockRespBody:   []byte(`invalid json`),
			mockStatusCode: 200,
			expectError:    "unmarshal response failed",
		},
		{
			name:           "NilData",
			location:       "119.262647,26.106131",
			enabled:        true,
			mockRespBody:   []byte(`{"data": null, "base": {"code": 0, "msg": "success"}}`),
			mockStatusCode: 200,
			expectError:    "SignedUrlData is nil",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {

			// mock c.Do，控制 resp 的状态码和 body
			mockey.Mock((*client.Client).Do).To(func(_ *client.Client, _ context.Context, _ *protocol.Request, resp *protocol.Response) error {
				if tc.mockDoError != nil {
					return tc.mockDoError
				}
				resp.SetStatusCode(tc.mockStatusCode)
				resp.SetBody(tc.mockRespBody)
				return nil
			}).Build()

			mockClientSet := &base.ClientSet{}
			svc := NewCommonService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			signedURL, headers, err := svc.GetSignedApiUrl(tc.location)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
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
