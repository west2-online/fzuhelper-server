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
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/bytedance/sonic"
	hertzconfig "github.com/cloudwego/hertz/pkg/common/config"
	"github.com/cloudwego/hertz/pkg/protocol"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

const (
	jobFairSourceURL      = "http://fjrclh.fzu.edu.cn/CmsInterface/getDateZPHKeynoteList_month"
	jobFairDetailURL      = "http://fjrclh.fzu.edu.cn/cms/zphdetail.html"
	lectureDetailURL      = "http://fjrclh.fzu.edu.cn/cms/xjhdetail.html"
	jobFairRequestTimeout = 10 * time.Second
)

var chinaStandardTime = time.FixedZone("Asia/Shanghai", int(8*time.Hour/time.Second))

type jobFairSourceItem struct {
	ID        string  `json:"id"`
	Title     *string `json:"title"`
	Place     *string `json:"place"`
	StartTime string  `json:"start_time"`
	Time      string  `json:"time"`
	ZPHVal    string  `json:"zphval"`
}

type jobFairSourceResponse struct {
	Success bool                `json:"success"`
	Events  []jobFairSourceItem `json:"zhaopinhui_keynoteList"`
}

func (s *CommonService) GetJobFair(month string) ([]*model.JobFairEvent, error) {
	parsedMonth, err := time.Parse("2006-01", month)
	if err != nil || parsedMonth.Format("2006-01") != month {
		return nil, errno.ParamError.WithMessage(fmt.Sprintf("invalid month %q, expected YYYY-MM", month))
	}

	cacheKey := ""
	if s.cache != nil && s.cache.Common != nil {
		cacheKey = s.cache.Common.JobFairKey(month)
		if s.cache.IsKeyExist(s.ctx, cacheKey) {
			events, cacheErr := s.cache.Common.GetJobFair(s.ctx, cacheKey)
			if cacheErr != nil {
				return nil, fmt.Errorf("service get job fair: read cache failed: %w", cacheErr)
			}
			return events, nil
		}
	}

	req := protocol.AcquireRequest()
	resp := protocol.AcquireResponse()
	defer func() {
		protocol.ReleaseRequest(req)
		protocol.ReleaseResponse(resp)
	}()

	form := url.Values{"dateday": {strings.Replace(month, "-", "/", 1)}}
	req.SetMethod(consts.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/x-www-form-urlencoded; charset=UTF-8"))
	req.SetBodyString(form.Encode())
	req.SetRequestURI(jobFairSourceURL)
	req.SetOptions(
		hertzconfig.WithDialTimeout(jobFairRequestTimeout),
		hertzconfig.WithReadTimeout(jobFairRequestTimeout),
		hertzconfig.WithWriteTimeout(jobFairRequestTimeout),
		hertzconfig.WithRequestTimeout(jobFairRequestTimeout),
	)

	if err = s.httpClient.Do(s.ctx, req, resp); err != nil {
		return nil, fmt.Errorf("service get job fair: request source failed: %w", err)
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return nil, fmt.Errorf("service get job fair: unexpected status code %d", resp.StatusCode())
	}

	var sourceResponse jobFairSourceResponse
	if err = sonic.Unmarshal(resp.Body(), &sourceResponse); err != nil {
		return nil, fmt.Errorf("service get job fair: unmarshal response failed: %w", err)
	}
	if !sourceResponse.Success {
		return nil, fmt.Errorf("service get job fair: source response was unsuccessful")
	}

	events := make([]*model.JobFairEvent, 0, len(sourceResponse.Events))
	for _, sourceEvent := range sourceResponse.Events {
		startsAt, parseErr := time.ParseInLocation(
			"20060102 15:04",
			sourceEvent.StartTime+" "+sourceEvent.Time,
			chinaStandardTime,
		)
		if parseErr != nil {
			return nil, fmt.Errorf("service get job fair: parse event time failed for %q: %w", sourceEvent.ID, parseErr)
		}

		detailURL := jobFairDetailURL
		if sourceEvent.ZPHVal == "3" {
			detailURL = lectureDetailURL
		}

		events = append(events, &model.JobFairEvent{
			Id:        sourceEvent.ID,
			Title:     normalizeJobFairText(sourceEvent.Title),
			Place:     normalizeJobFairText(sourceEvent.Place),
			Time:      sourceEvent.Time,
			StartsAt:  startsAt.Unix(),
			DateKey:   startsAt.Format(time.DateOnly),
			DetailUrl: detailURL + "?id=" + url.QueryEscape(sourceEvent.ID),
		})
	}
	if cacheKey != "" {
		if cacheErr := s.cache.Common.SetJobFair(s.ctx, cacheKey, events); cacheErr != nil {
			logger.Errorf("service get job fair: write cache failed: %v", cacheErr)
		}
	}

	return events, nil
}

func normalizeJobFairText(value *string) string {
	if value == nil {
		return ""
	}
	return strings.TrimSpace(html.UnescapeString(*value))
}
