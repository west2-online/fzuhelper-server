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

package main

import (
	"context"

	"github.com/west2-online/fzuhelper-server/cmd/template/pack"
	"github.com/west2-online/fzuhelper-server/cmd/template/service"
	template "github.com/west2-online/fzuhelper-server/kitex_gen/template"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

// TemplateServiceImpl implements the last service interface defined in the IDL.
type TemplateServiceImpl struct{}

// Ping implements the TemplateServiceImpl interface.
func (s *TemplateServiceImpl) Ping(ctx context.Context, req *template.PingRequest) (resp *template.PingResponse, err error) {
	resp = new(template.PingResponse)

	if req.Text != nil && len(*req.Text) == 0 {
		resp.Base = pack.BuildBaseResp(errno.ParamEmpty)
	}

	text, err := service.NewTemplateService(ctx).Ping(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	resp.Base = pack.BuildBaseResp(errno.Success)
	resp.Pong = text
	return
}
