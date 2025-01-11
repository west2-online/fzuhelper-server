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
	"net/http"

	"github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *AcademicService) GetPlan() (*[]byte, error) {
	userHeader, err := context.GetLoginData(s.ctx)
	if err != nil {
		return nil, err
	}
	stu := jwch.NewStudent().WithLoginData(userHeader.Id, utils.ParseCookies(userHeader.Cookies))
	url, err := stu.GetCultivatePlan()
	if err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "AcademicService.GetPlan error:%v", err)
	}

	urlReq, err := http.NewRequest(constants.GetPlanMethod, url, nil)
	if err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "AcademicService.GetPlan request error:%v", err)
	}
	urlReq.Header.Set("Cookie", userHeader.Cookies)
	htmlSource, err := getHtmlSource(urlReq)
	if err != nil {
		return nil, errno.Errorf(errno.InternalServiceErrorCode, "AcademicService.GetPlan getHtmlSource error:%v", err)
	}

	return htmlSource, nil
}

func getHtmlSource(r *http.Request) (*[]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, err
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}
	htmlSource := buf.Bytes()
	return &htmlSource, nil
}
