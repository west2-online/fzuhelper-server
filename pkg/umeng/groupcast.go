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
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func getChannelProperties(title, content string) AndroidChannelProperties {
	return AndroidChannelProperties{
		ChannelActivity:         config.Vendors.ChannelActivity,
		VivoCategory:            config.Vendors.VivoCategory,
		OppoChannelID:           config.Vendors.Oppo.ChannelID,
		OppoCategory:            config.Vendors.Oppo.Category,
		OppoNotifyLevel:         config.Vendors.Oppo.NotifyLevel,
		HuaweiChannelImportance: config.Vendors.Huawei.ChannelImportance,
		HuaweiChannelCategory:   config.Vendors.Huawei.ChannelCategory,
		OppoPrivateMsgTemplate: OppoPrivateMsgTemplate{
			PrivateMsgTemplateID: config.Vendors.Oppo.PrivateMsgTemplate.PrivateMsgTemplateID,
			PrivateTitleParameters: OppoPrivateTitleParameters{
				Title: title,
			},
			PrivateContentParameters: OppoPrivateContentParameters{
				Content: content,
			},
		},
		LocalProperties: LocalProperties{
			ChannelID:   config.Vendors.LocalProperties.ChannelID,
			ChannelName: config.Vendors.LocalProperties.ChannelName,
		},
	}
}

// getXiaomiNoticeProperties 按推送类型从配置中选取小米模板，
// 模板参数直接取调用方透传的 keywords[0]，不做字符串反解
func getXiaomiNoticeProperties(pushType string, keywords []string, notice config.XiaomiNotice) (string, *XiaomiExtraProperties) {
	template, ok := notice[pushType]
	if !ok {
		return "", nil
	}
	if template.ChannelID == "" || template.TemplateID == "" || len(keywords) == 0 || keywords[0] == "" {
		return "", nil
	}

	// 小米推送要求 extra.template_param 是 JSON 字符串而非 JSON 对象，
	// 例如 {"keywords1":"数据结构"}，MiPush 会按模板中预配置的 {$keywords1$} 变量拼装消息，
	// 因此这里先用 json.Marshal 把 map 序列化成合法 JSON 字符串再随请求下发，
	templateParam, err := json.Marshal(map[string]string{
		constants.UmengXiaomiTemplateKeyword: keywords[0],
	})
	if err != nil {
		logger.Errorf("umeng.getXiaomiNoticeProperties: failed to marshal xiaomi template_param: %v", err)
		return "", nil
	}

	return template.ChannelID, &XiaomiExtraProperties{
		TemplateID:    template.TemplateID,
		TemplateParam: string(templateParam),
	}
}

// PushByType 按推送类型同时下发安卓、iOS 与鸿蒙推送，任一端失败仅记录日志，不影响业务（尽力而为）。
// text 由调用方按推送类型拼装好传入；keywords 为模板参数（透传复用，不做字符串反解），例如成绩通知传 []string{courseName}
func PushByType(pushType, title, text string, keywords []string, ticker, tag, description, deeplink string) {
	if err := SendAndroidGroupcastWithGoApp(pushType, title, text, ticker, tag, description, deeplink, keywords); err != nil {
		logger.Errorf("umeng.PushByType: %s failed to send Android groupcast: %v", pushType, err)
	}
	if err := SendIOSGroupcast(title, "", text, tag, description, deeplink); err != nil {
		logger.Errorf("umeng.PushByType: %s failed to send IOS groupcast: %v", pushType, err)
	}
	if err := SendHarmonyGroupcast(title, text, tag, description, deeplink); err != nil {
		logger.Errorf("umeng.PushByType: %s failed to send Harmony groupcast: %v", pushType, err)
	}
}

func SendAndroidGroupcastWithGoApp(pushType, title, text, ticker, tag, description, deeplink string, keywords []string) error {
	channelProperties := getChannelProperties(title, text)
	xiaomiChannelID, xiaomiExtraProperties := getXiaomiNoticeProperties(
		pushType,
		keywords,
		config.Vendors.XiaomiNotice,
	)
	if xiaomiChannelID != "" && xiaomiExtraProperties != nil {
		channelProperties.XiaoMiChannelID = xiaomiChannelID
		channelProperties.XiaoMiExtraProperties = xiaomiExtraProperties
	}

	message := AndroidGroupcastMessage{
		AppKey:    config.Umeng.Android.AppKey,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Type:      "groupcast",
		Filter: Filter{
			Where: Where{
				And: []map[string]string{
					{"tag": tag},
				},
			},
		},
		Payload: AndroidPayload{
			DisplayType: "notification",
			Body: AndroidBody{
				Title:       title,
				Text:        text,
				Ticker:      ticker,
				PlaySound:   "true",
				PlayVibrate: "true",
				PlayLights:  "true",
				AfterOpen:   "go_app",
			},
			Extra: map[string]string{
				"deeplink": deeplink,
			},
		},
		Policy: AndroidPolicy{
			ExpireTime:               time.Now().Add(constants.UmengMessageExpireTime).Format("2006-01-02 15:04:05"),
			NotificationClosedFilter: true,
		},
		Description:       description,
		Category:          1, // 系统消息
		ChannelProperties: channelProperties,
	}

	return sendGroupcast(config.Umeng.Android.AppMasterSecret, message)
}

// Android广播函数
func SendAndroidGroupcastWithUrl(title, text, ticker, url, tag, description string) error {
	message := AndroidGroupcastMessage{
		AppKey:    config.Umeng.Android.AppKey,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Type:      "groupcast",
		Filter: Filter{
			Where: Where{
				And: []map[string]string{
					{"tag": tag},
				},
			},
		},
		Payload: AndroidPayload{
			DisplayType: "notification",
			Body: AndroidBody{
				Title:       title,
				Text:        text,
				Ticker:      ticker,
				PlaySound:   "true",
				PlayVibrate: "true",
				PlayLights:  "true",
				AfterOpen:   "go_url",
				URL:         url,
			},
		},
		Policy: AndroidPolicy{
			ExpireTime:               time.Now().Add(constants.UmengMessageExpireTime).Format("2006-01-02 15:04:05"),
			NotificationClosedFilter: true,
		},
		Description:       description,
		Category:          1, // 系统消息
		ChannelProperties: getChannelProperties(title, text),
	}

	return sendGroupcast(config.Umeng.Android.AppMasterSecret, message)
}

// iOS广播函数
func SendIOSGroupcast(title, subtitle, body, tag, description, deeplink string) error {
	message := IOSGroupcastMessage{
		AppKey:    config.Umeng.IOS.AppKey,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Type:      "groupcast",
		Filter: Filter{
			Where: Where{
				And: []map[string]string{
					{"tag": tag},
				},
			},
		},
		Payload: IOSPayload{
			Aps: IOSAps{
				Alert: IOSAlert{
					Title:    title,
					Subtitle: subtitle,
					Body:     body,
				},
				Sound:             "default",
				InterruptionLevel: "active",
			},
			Deeplink: deeplink,
		},
		Policy: IOSPolicy{
			ExpireTime: time.Now().Add(constants.UmengMessageExpireTime).Format("2006-01-02 15:04:05"),
		},
		Description: description,
	}

	return sendGroupcast(config.Umeng.IOS.AppMasterSecret, message)
}

// Harmony广播函数
func SendHarmonyGroupcast(title, text, tag, description, deeplink string) error {
	message := HarmonyGroupcastMessage{
		AppKey:    config.Umeng.Harmony.AppKey,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Type:      "groupcast",
		Filter: Filter{
			Where: Where{
				And: []map[string]string{
					{"tag": tag},
				},
			},
		},
		Payload: HarmonyPayload{
			DisplayType: "notification",
			Body: HarmonyBody{
				Title: title,
				Text:  text,
			},
			Extra: map[string]string{
				"deeplink": deeplink,
			},
		},
		Policy: HarmonyPolicy{
			ExpireTime: time.Now().Add(constants.UmengMessageExpireTime).Format("2006-01-02 15:04:05"),
		},
		Description: description,
		ChannelProperties: HarmonyChannelProperties{
			HarmonyChannelCategory: config.Vendors.Harmony.ChannelCategory,
		},
	}

	return sendGroupcast(config.Umeng.Harmony.AppMasterSecret, message)
}

// 通用广播发送逻辑
func sendGroupcast(appMasterSecret string, message interface{}) error {
	postBody, err := json.Marshal(message)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendGroupcast : failed to marshal JSON: %v", err)
	}

	sign := generateSign("POST", constants.UmengURL, string(postBody), appMasterSecret)

	req, err := http.NewRequest("POST", constants.UmengURL+"?sign="+sign, bytes.NewBuffer(postBody))
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendGroupcast : failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendGroupcast : failed to send request: %v", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logger.Warnf("umeng.sendGroupcast : failed to close response body: %v", err)
		}
	}(resp.Body)

	var response UmengResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendGroupcast : failed to decode response: %v", err)
	}

	if response.Ret != "SUCCESS" {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendGroupcast : Groupcast failed: %s (%s)", response.Data.ErrorMsg, response.Data.ErrorCode)
	}

	return nil
}

// 生成MD5签名
func generateSign(method, url, postBody, appMasterSecret string) string {
	data := fmt.Sprintf("%s%s%s%s", method, url, postBody, appMasterSecret)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
