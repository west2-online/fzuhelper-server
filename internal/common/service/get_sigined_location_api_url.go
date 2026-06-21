package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/cloudwego/hertz/pkg/app/client"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

const dialTimeoutSeconds = 15 * time.Second

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
	disabled_msg := config.SignedLocationApiUrl.DisableMsg
	endpoint := config.SignedLocationApiUrl.Endpoint

	if !enabled {
		return "", nil, fmt.Errorf("service get signed api url: %s", disabled_msg)
	}

	if location == "" {
		return "", nil, fmt.Errorf("service get signed api url: location is empty")
	}

	c, err := client.NewClient(client.WithDialTimeout(dialTimeoutSeconds))
	if err != nil {
		return "", nil, fmt.Errorf("service get signed api url: create client failed %w", err)
	}

	req := &protocol.Request{}
	resp := &protocol.Response{}
	req.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetBodyString(`{"location":"` + location + `"}`)
	req.SetRequestURI(endpoint)

	if err = c.Do(context.Background(), req, resp); err != nil {
		return "", nil, fmt.Errorf("service get signed api url: request service failed %w", err)
	}

	if resp.StatusCode() >= consts.StatusBadRequest {
		return "", nil, fmt.Errorf("service get signed api url: error response with status %d", resp.StatusCode())
	}

	var rsl signedUrlResp
	if err = json.Unmarshal(resp.Body(), &rsl); err != nil {
		return "", nil, fmt.Errorf("service get signed api url: unmarshal response failed %w", err)
	}

	if rsl.Data == nil {
		return "", nil, fmt.Errorf("service get signed api url: request service failed: unmarshal response failed SignedUrlData is nil")
	}

	signedURL := rsl.Data.SignedURL
	headers := rsl.Data.Headers

	return signedURL, headers, nil
}
