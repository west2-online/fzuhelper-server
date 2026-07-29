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
	"encoding/json"
	"fmt"

	hertzconfig "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

type Data struct {
	SignedURL string            `json:"signed_url"`
	Headers   map[string]string `json:"headers"`
}

type signedUrlResp struct {
	Data *Data           `json:"data"`
	Base *model.BaseResp `json:"base"`
}

func (s *CommonService) GetSignedApiUrl(location string) (string, map[string]string, error) {
	enabled := config.SignedLocationApiUrl.Enabled
	disabledMsg := config.SignedLocationApiUrl.DisableMsg
	endpoint := config.SignedLocationApiUrl.Endpoint

	if !enabled {
		return "", nil, fmt.Errorf("service get signed api url: %s", disabledMsg)
	}

	if location == "" {
		return "", nil, fmt.Errorf("service get signed api url: location is empty")
	}

	req := protocol.AcquireRequest()
	resp := protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseRequest(req)
		protocol.ReleaseResponse(resp)
	}()

	req.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	body, err := json.Marshal(map[string]string{"location": location})
	if err != nil {
		return "", nil, fmt.Errorf("service get signed api url: marshal request body failed %w", err)
	}
	req.SetBody(body)
	req.SetRequestURI(endpoint)
	req.SetOptions(
		hertzconfig.WithDialTimeout(constants.SignedLocationTimeout),
		hertzconfig.WithReadTimeout(constants.SignedLocationTimeout),
		hertzconfig.WithWriteTimeout(constants.SignedLocationTimeout),
		hertzconfig.WithRequestTimeout(constants.SignedLocationTimeout),
	)

	if err = s.httpClient.Do(s.ctx, req, resp); err != nil {
		return "", nil, fmt.Errorf("service get signed api url: request service failed %w", err)
	}

	if resp.StatusCode() != consts.StatusOK {
		return "", nil, fmt.Errorf("service get signed api url: unexpected status code %d, body: %s",
			resp.StatusCode(), resp.Body())
	}

	var respData signedUrlResp
	if err = json.Unmarshal(resp.Body(), &respData); err != nil {
		return "", nil, fmt.Errorf("service get signed api url: unmarshal response failed %w", err)
	}

	if respData.Base == nil {
		return "", nil, fmt.Errorf("service get signed api url: response base is nil")
	}

	if respData.Base.Code != errno.SuccessCode {
		return "", nil, fmt.Errorf("service get signed api url: location service returned business error: %w",
			errno.NewErrNo(respData.Base.Code, respData.Base.Msg))
	}

	if respData.Data == nil {
		return "", nil, fmt.Errorf("service get signed api url: request service failed: unmarshal response failed SignedUrlData is nil")
	}

	return respData.Data.SignedURL, respData.Data.Headers, nil
}
