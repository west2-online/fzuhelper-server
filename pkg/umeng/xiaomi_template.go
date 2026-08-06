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

// XiaomiTemplateKind 描述一种推送类型对应的小米模板 keyword 提取规则
type XiaomiTemplateKind struct {
	// KeywordSuffix 从 text 中提取 keywords1 时去掉的正文后缀，空串表示取完整 text
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
