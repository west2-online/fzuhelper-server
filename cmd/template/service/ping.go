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

import "github.com/west2-online/fzuhelper-server/kitex_gen/template"

func (s *TemplateService) Ping(req *template.PingRequest) (string, error) {
	return *req.Text, nil

	// 这里负责处理业务请求，如果有需要可以继续和 dal（data access layout）做交互
}
