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
	"encoding/json"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

func TestPushByType(t *testing.T) {
	tests := []struct {
		name       string
		pushType   string
		androidErr error
		iosErr     error
		harmonyErr error
	}{
		{
			name:     "send to android, ios and harmony",
			pushType: constants.UmengPushTypeScore,
		},
		{
			name:       "ignore android, ios and harmony failures",
			pushType:   constants.UmengPushTypeExam,
			androidErr: assert.AnError,
			iosErr:     assert.AnError,
			harmonyErr: assert.AnError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var androidPushType string
			androidPatch := mockey.Mock(SendAndroidGroupcastWithGoApp).To(
				func(pushType, title, text, ticker, tag, description, deeplink string) error {
					androidPushType = pushType
					return tt.androidErr
				},
			).Build()
			defer androidPatch.UnPatch()

			iosCalled := false
			iosPatch := mockey.Mock(SendIOSGroupcast).To(
				func(title, subtitle, body, tag, description, deeplink string) error {
					iosCalled = true
					return tt.iosErr
				},
			).Build()
			defer iosPatch.UnPatch()

			harmonyCalled := false
			harmonyPatch := mockey.Mock(SendHarmonyGroupcast).To(
				func(title, text, tag, description, deeplink string) error {
					harmonyCalled = true
					return tt.harmonyErr
				},
			).Build()
			defer harmonyPatch.UnPatch()

			// PushByType 为尽力而为，三端失败仅记录日志，不应 panic
			PushByType(tt.pushType, "title", "text", "ticker", "tag", "description", "deeplink")

			assert.Equal(t, tt.pushType, androidPushType)
			assert.True(t, iosCalled)
			assert.True(t, harmonyCalled)
		})
	}
}

func TestHarmonyGroupcastMessageJSON(t *testing.T) {
	tests := []struct {
		name                string
		msg                 HarmonyGroupcastMessage
		wantTitle           string
		wantText            string
		wantDeeplink        string
		wantExpireTime      string
		wantChannelCategory string
		wantCategoryExists  bool
	}{
		{
			name: "embed all fields",
			msg: HarmonyGroupcastMessage{
				AppKey:    "harmony-appkey",
				Timestamp: "1750000000",
				Type:      "groupcast",
				Filter: Filter{
					Where: Where{
						And: []map[string]string{{"tag": "jwch-notice"}},
					},
				},
				Payload: HarmonyPayload{
					DisplayType: "notification",
					Body: HarmonyBody{
						Title: "教务处通知",
						Text:  "关于补考安排的通知",
					},
					Extra: map[string]string{"deeplink": "fzuhelper://office_notice"},
				},
				Policy:      HarmonyPolicy{ExpireTime: "2026-08-10 15:04:05"},
				Description: "教务处通知",
				ChannelProperties: HarmonyChannelProperties{
					HarmonyChannelCategory: "SERVICE",
				},
			},
			wantTitle:           "教务处通知",
			wantText:            "关于补考安排的通知",
			wantDeeplink:        "fzuhelper://office_notice",
			wantExpireTime:      "2026-08-10 15:04:05",
			wantChannelCategory: "SERVICE",
			wantCategoryExists:  true,
		},
		{
			name: "omit empty channel category",
			msg: HarmonyGroupcastMessage{
				AppKey:            "harmony-appkey",
				Timestamp:         "1750000000",
				Type:              "groupcast",
				Filter:            Filter{},
				Payload:           HarmonyPayload{DisplayType: "notification", Body: HarmonyBody{Title: "t", Text: "x"}},
				Policy:            HarmonyPolicy{},
				ChannelProperties: HarmonyChannelProperties{},
			},
			wantTitle: "t",
			wantText:  "x",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.msg)
			assert.NoError(t, err)

			var got map[string]any
			assert.NoError(t, json.Unmarshal(data, &got))

			assert.Equal(t, "harmony-appkey", got["appkey"])
			assert.Equal(t, "groupcast", got["type"])

			payload, ok := got["payload"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, "notification", payload["display_type"])
			body, ok := payload["body"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, tt.wantTitle, body["title"])
			assert.Equal(t, tt.wantText, body["text"])
			if tt.wantDeeplink != "" {
				extra, ok := payload["extra"].(map[string]any)
				assert.True(t, ok)
				assert.Equal(t, tt.wantDeeplink, extra["deeplink"])
			}

			policy, ok := got["policy"].(map[string]any)
			assert.True(t, ok)
			assert.Equal(t, tt.wantExpireTime, policy["expire_time"])

			channelProperties, ok := got["channel_properties"].(map[string]any)
			assert.True(t, ok)
			if tt.wantCategoryExists {
				assert.Equal(t, tt.wantChannelCategory, channelProperties["harmony_channel_category"])
			} else {
				_, exists := channelProperties["harmony_channel_category"]
				assert.False(t, exists)
			}
		})
	}
}

func TestGetXiaomiNoticeProperties(t *testing.T) {
	notice := config.XiaomiNotice{
		constants.UmengPushTypeScore: {
			ChannelID:  "score-channel",
			TemplateID: "P12395",
		},
		constants.UmengPushTypeExam: {
			ChannelID:  "exam-channel",
			TemplateID: "P12394",
		},
		constants.UmengPushTypeTeaching: {
			ChannelID:  "teaching-channel",
			TemplateID: "P12325",
		},
	}

	tests := []struct {
		name          string
		text          string
		pushType      string
		wantChannelID string
		wantTemplate  string
		wantKeyword   string
	}{
		{
			name:          "score",
			text:          "数据结构" + constants.UmengGradeNotificationBodySuffix,
			pushType:      constants.UmengPushTypeScore,
			wantChannelID: "score-channel",
			wantTemplate:  "P12395",
			wantKeyword:   "数据结构",
		},
		{
			name:          "exam",
			text:          "数据结构" + constants.UmengExamNotificationBodySuffix,
			pushType:      constants.UmengPushTypeExam,
			wantChannelID: "exam-channel",
			wantTemplate:  "P12394",
			wantKeyword:   "数据结构",
		},
		{
			name:          "teaching",
			text:          "关于补考安排的通知",
			pushType:      constants.UmengPushTypeTeaching,
			wantChannelID: "teaching-channel",
			wantTemplate:  "P12325",
			wantKeyword:   "关于补考安排的通知",
		},
		{
			name:     "unknown push type",
			pushType: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChannelID, gotProperties := getXiaomiNoticeProperties(tt.text, tt.pushType, notice)
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
