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
	"net/http"
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
)

func SendAndroidGroupcastWithGoApp(appKey, appMasterSecret, ticker, title, text, tag string) error {
	message := AndroidGroupcastMessage{
		AppKey:    appKey,
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
				Ticker:    ticker,
				Title:     title,
				Text:      text,
				AfterOpen: "go_app",
			},
		},
		Policy: Policy{
			ExpireTime: time.Now().Add(constants.UmengMessageExpireTime).Format("2006-01-02 15:04:05"),
		},
		Description: "Android-广播通知",
	}

	return sendGroupcast(appMasterSecret, message)
}

// Android广播函数
func SendAndroidGroupcastWithUrl(appKey, appMasterSecret, ticker, title, text, tag, url string, channelProperties map[string]string) error {
	message := AndroidGroupcastMessage{
		AppKey:    appKey,
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
				Ticker:    ticker,
				Title:     title,
				Text:      text,
				AfterOpen: "go_url",
				URL:       url,
			},
		},
		Policy: Policy{
			ExpireTime: time.Now().Add(constants.UmengMessageExpireTime).Format("2006-01-02 15:04:05"),
		},
		ChannelProperties: channelProperties,
		Description:       "Android-广播通知",
	}

	return sendGroupcast(appMasterSecret, message)
}

// iOS广播函数
func SendIOSGroupcast(appKey, appMasterSecret, title, subtitle, body, tag string) error {
	message := IOSGroupcastMessage{
		AppKey:    appKey,
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
			},
		},
		Policy: Policy{
			ExpireTime: time.Now().Add(constants.UmengMessageExpireTime).Format("2006-01-02 15:04:05"),
		},
		Description: "iOS-广播通知",
	}

	return sendGroupcast(appMasterSecret, message)
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
	defer resp.Body.Close()

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
