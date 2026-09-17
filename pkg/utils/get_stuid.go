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

package utils

import (
	"strings"

	"github.com/cloudwego/hertz/pkg/app"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// 获取jwt写入context中的stuID
func GetStuID(c *app.RequestContext) (string, bool) {
	value, ok := c.Get(constants.StuIDContextKey)
	if !ok {
		return "", false
	}
	stuID, ok := value.(string)
	if !ok {
		return "", false
	}
	stuID = strings.TrimSpace(stuID)
	return stuID, stuID != ""
}
