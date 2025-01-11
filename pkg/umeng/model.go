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

// 公共返回结构
type UmengResponse struct {
	Ret  string `json:"ret"`
	Data struct {
		MsgID     string `json:"msg_id,omitempty"`
		ErrorCode string `json:"error_code,omitempty"`
		ErrorMsg  string `json:"error_msg,omitempty"`
	} `json:"data"`
}

// Android广播消息结构
type AndroidBroadcastMessage struct {
	AppKey            string            `json:"appkey"`
	Timestamp         string            `json:"timestamp"`
	Type              string            `json:"type"`
	Payload           AndroidPayload    `json:"payload"`
	Policy            Policy            `json:"policy"`
	ChannelProperties map[string]string `json:"channel_properties"`
	Description       string            `json:"description"`
}

type AndroidPayload struct {
	DisplayType string      `json:"display_type"`
	Body        AndroidBody `json:"body"`
}

type AndroidBody struct {
	Ticker    string `json:"ticker"`
	Title     string `json:"title"`
	Text      string `json:"text"`
	AfterOpen string `json:"after_open"`
}

// iOS广播消息结构
type IOSBroadcastMessage struct {
	AppKey      string     `json:"appkey"`
	Timestamp   string     `json:"timestamp"`
	Type        string     `json:"type"`
	Payload     IOSPayload `json:"payload"`
	Policy      Policy     `json:"policy"`
	Description string     `json:"description"`
}

type IOSPayload struct {
	Aps IOSAps `json:"aps"`
}

type IOSAps struct {
	Alert IOSAlert `json:"alert"`
}

type IOSAlert struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Body     string `json:"body"`
}

// 公共策略结构
type Policy struct {
	ExpireTime string `json:"expire_time"`
}
