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

// Android广播函数
func SendAndroidBroadcast(appKey, appMasterSecret, ticker, title, text, expireTime string) error {
	message := AndroidBroadcastMessage{
		AppKey:    appKey,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Type:      "broadcast",
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
			ExpireTime: expireTime,
		},
		ChannelProperties: map[string]string{
			"channel_activity": "xxx",
		},
		Description: "测试广播通知-Android",
	}

	return sendBroadcast(appMasterSecret, message)
}

// iOS广播函数
func SendIOSBroadcast(appKey, appMasterSecret, title, subtitle, body, expireTime string) error {
	message := IOSBroadcastMessage{
		AppKey:    appKey,
		Timestamp: fmt.Sprintf("%d", time.Now().Unix()),
		Type:      "broadcast",
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
			ExpireTime: expireTime,
		},
		Description: "测试广播通知-iOS",
	}

	return sendBroadcast(appMasterSecret, message)
}

// 通用广播发送逻辑
func sendBroadcast(appMasterSecret string, message interface{}) error {
	postBody, err := json.Marshal(message)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendBroadcast : failed to marshal JSON: %v", err)
	}

	sign := generateSign("POST", constants.UmengURL, string(postBody), appMasterSecret)

	req, err := http.NewRequest("POST", constants.UmengURL+"?sign="+sign, bytes.NewBuffer(postBody))
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendBroadcast : failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendBroadcast : failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendBroadcast : unexpected response code: %v", resp.StatusCode)
	}

	var response UmengResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendBroadcast : failed to decode response: %v", err)
	}

	if response.Ret != "SUCCESS" {
		return errno.Errorf(errno.InternalServiceErrorCode, "umeng.sendBroadcast : broadcast failed: %s (%s)", response.Data.ErrorMsg, response.Data.ErrorCode)
	}

	logger.Infof("Broadcast sent successfully! MsgID: %s\n", response.Data.MsgID)
	return nil
}

// 生成MD5签名
func generateSign(method, url, postBody, appMasterSecret string) string {
	data := fmt.Sprintf("%s%s%s%s", method, url, postBody, appMasterSecret)
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
