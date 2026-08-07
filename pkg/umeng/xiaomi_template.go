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

package umeng

import (
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// XiaomiTemplateKind 描述一种推送类型对应的小米模板拼装规则
type XiaomiTemplateKind struct {
	// KeywordSuffix 拼装推送正文时追加在 keywords[0] 后的固定后缀，空串表示正文即 keywords[0]
	KeywordSuffix string
}

// xiaomiTemplateKinds 推送类型注册表，新增推送类型只需在配置中增加模板并在这里注册一项
var xiaomiTemplateKinds = map[string]XiaomiTemplateKind{
	constants.UmengPushTypeScore: {
		KeywordSuffix: constants.UmengGradeNotificationBodySuffix,
	},
	constants.UmengPushTypeExam: {
		KeywordSuffix: constants.UmengExamNotificationBodySuffix,
	},
	constants.UmengPushTypeTeaching: {},
}

// buildPushText 按推送类型用透传的模板参数拼出通知正文，
// 例如成绩/考试通知为 keywords[0]+固定后缀，教务处通知正文即 keywords[0]
func buildPushText(pushType string, keywords []string) string {
	kind, ok := xiaomiTemplateKinds[pushType]
	if !ok || len(keywords) == 0 {
		return ""
	}
	return keywords[0] + kind.KeywordSuffix
}
