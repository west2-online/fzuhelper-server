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
	"testing"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func TestGetXiaomiNoticeProperties(t *testing.T) {
	notice := config.XiaomiNotice{
		Score: config.XiaomiNoticeTemplate{
			ChannelID:  "score-channel",
			TemplateID: "P12395",
		},
		Exam: config.XiaomiNoticeTemplate{
			ChannelID:  "exam-channel",
			TemplateID: "P12394",
		},
		Teaching: config.XiaomiNoticeTemplate{
			ChannelID:  "teaching-channel",
			TemplateID: "P12325",
		},
	}

	tests := []struct {
		name          string
		text          string
		deeplink      string
		wantChannelID string
		wantTemplate  string
		wantKeyword   string
	}{
		{
			name:          "score",
			text:          "数据结构" + constants.UmengGradeNotificationBodySuffix,
			deeplink:      constants.UmengGradeDeeplink,
			wantChannelID: "score-channel",
			wantTemplate:  "P12395",
			wantKeyword:   "数据结构",
		},
		{
			name:          "exam",
			text:          "数据结构" + constants.UmengExamNotificationBodySuffix,
			deeplink:      constants.UmengExamRoomDeeplink,
			wantChannelID: "exam-channel",
			wantTemplate:  "P12394",
			wantKeyword:   "数据结构",
		},
		{
			name:          "teaching",
			text:          "关于补考安排的通知",
			deeplink:      constants.UmengJwchNoticeDeeplink + "?url=https%3A%2F%2Fexample.com%2Fnotice",
			wantChannelID: "teaching-channel",
			wantTemplate:  "P12325",
			wantKeyword:   "关于补考安排的通知",
		},
		{
			name:     "unknown deeplink",
			deeplink: "fzuhelper://unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChannelID, gotProperties := getXiaomiNoticeProperties(tt.text, tt.deeplink, notice)
			if gotChannelID != tt.wantChannelID {
				t.Fatalf("channel ID = %q, want %q", gotChannelID, tt.wantChannelID)
			}
			if tt.wantTemplate == "" {
				if gotProperties != nil {
					t.Fatalf("properties = %+v, want nil", gotProperties)
				}
				return
			}
			if gotProperties == nil {
				t.Fatal("properties = nil, want configured Xiaomi properties")
			}
			if gotProperties.TemplateID != tt.wantTemplate {
				t.Errorf("template ID = %q, want %q", gotProperties.TemplateID, tt.wantTemplate)
			}
			wantParam := `{"keywords1":"` + tt.wantKeyword + `"}`
			if gotProperties.TemplateParam != wantParam {
				t.Errorf("template param = %q, want %q", gotProperties.TemplateParam, wantParam)
			}
		})
	}
}
