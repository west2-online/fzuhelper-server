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

func getChannelProperties() AndroidChannelProperties {
	return AndroidChannelProperties{
		ChannelActivity: config.Vendors.ChannelActivity,
		XiaoMiChannelID: config.Vendors.XiaoMiChannelID,
		// VivoCategory:            config.Vendors.VivoCategory,
		// OppoChannelID:           config.Vendors.Oppo.ChannelID,
		// OppoCategory:            config.Vendors.Oppo.Category,
		// OppoNotifyLevel:         config.Vendors.Oppo.NotifyLevel,
		HuaweiChannelImportance: config.Vendors.Huawei.ChannelImportance,
		// HuaweiChannelCategory:   config.Vendors.Huawei.ChannelCategory,
		// HonorChannelImportance:  config.Vendors.Honor.ChannelImportance,
	}
}

func SendAndroidGroupcastWithGoApp(title, text, ticker, tag string) error {
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
				Title:     title,
				Text:      text,
				Ticker:    ticker,
				AfterOpen: "go_app",
			},
		},
		Policy: AndroidPolicy{
			ExpireTime:               time.Now().Add(constants.UmengMessageExpireTime).Format("2006-01-02 15:04:05"),
			NotificationClosedFilter: true,
		},
		Description:       "Android-广播通知",
		Category:          0,
		ChannelProperties: getChannelProperties(),
	}

	return sendGroupcast(config.Umeng.Android.AppMasterSecret, message)
}

// Android广播函数
func SendAndroidGroupcastWithUrl(title, text, ticker, url, tag string) error {
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
				Title:     title,
				Text:      text,
				Ticker:    ticker,
				AfterOpen: "go_url",
				URL:       url,
			},
		},
		Policy: AndroidPolicy{
			ExpireTime:               time.Now().Add(constants.UmengMessageExpireTime).Format("2006-01-02 15:04:05"),
			NotificationClosedFilter: true,
		},
		Description:       "Android-广播通知",
		Category:          0,
		ChannelProperties: getChannelProperties(),
	}

	return sendGroupcast(config.Umeng.Android.AppMasterSecret, message)
}

// iOS广播函数
func SendIOSGroupcast(title, subtitle, body, tag string) error {
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
		},
		Policy: IOSPolicy{
			ExpireTime: time.Now().Add(constants.UmengMessageExpireTime).Format("2006-01-02 15:04:05"),
		},
		Description: "iOS-广播通知",
	}

	return sendGroupcast(config.Umeng.IOS.AppMasterSecret, message)
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

	logger.Infof("Groupcast sent successfully! MsgID: %s\n", response.Data.MsgID)
	return nil
}

// 生成MD5签名
func generateSign(method, url, postBody, appMasterSecret string) string {
	data := fmt.Sprintf("%s%s%s%s", method, url, postBody, appMasterSecret)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
