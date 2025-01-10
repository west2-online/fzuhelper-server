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
	"fmt"
	"net/http"

	"github.com/west2-online/fzuhelper-server/kitex_gen/academic"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *AcademicService) GetPlan(req *academic.GetPlanRequest) (*[]byte, error) {
	stu := jwch.NewStudent().WithLoginData(req.Id, utils.ParseCookies(req.Cookies))
	url, err := stu.GetCultivatePlan()
	if err != nil {
		return nil, fmt.Errorf("AcademicService.GetPlan error:%w", err)
	}

	urlReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("AcademicService.GetPlan request error:%w", err)
	}
	urlReq.Header.Set("Cookie", req.Cookies)
	urlReq.Header.Set("Identifier", req.Id)
	client := &http.Client{}
	resp, err := client.Do(urlReq)
	if err != nil {
		return nil, fmt.Errorf("AcademicService.GetPlan request error:%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("AcademicService.GetPlan request status code error:%d", resp.StatusCode)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("AcademicService.GetPlan response body error:%w", err)
	}
	htmlSource := buf.Bytes()

	return &htmlSource, nil
}
