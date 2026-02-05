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

// 具体定义查看友盟官方文档：https://developer.umeng.com/docs/67966/detail/68343
// UmengResponse 公共返回结构
type UmengResponse struct {
	Ret  string `json:"ret"`
	Data struct {
		MsgID     string `json:"msg_id,omitempty"`
		TaskID    string `json:"task_id,omitempty"`
		ErrorCode string `json:"error_code,omitempty"`
		ErrorMsg  string `json:"error_msg,omitempty"`
	} `json:"data"`
}

// AndroidGroupcastMessage Android广播消息结构
type AndroidGroupcastMessage struct {
	AppKey            string                   `json:"appkey"`
	Timestamp         string                   `json:"timestamp"`
	Type              string                   `json:"type"`
	Filter            Filter                   `json:"filter"`
	Payload           AndroidPayload           `json:"payload"`
	Policy            AndroidPolicy            `json:"policy"`
	Description       string                   `json:"description"`
	Category          int                      `json:"category"`
	ChannelProperties AndroidChannelProperties `json:"channel_properties"`
}

type AndroidPayload struct {
	DisplayType string      `json:"display_type"`
	Body        AndroidBody `json:"body"`
}

type AndroidBody struct {
	Title       string `json:"title"`
	Text        string `json:"text"`
	Ticker      string `json:"ticker"`
	PlaySound   string `json:"play_sound"`
	PlayVibrate string `json:"play_vibrate"`
	PlayLights  string `json:"play_lights"`
	AfterOpen   string `json:"after_open"`
	URL         string `json:"url"`
}

type AndroidChannelProperties struct {
	ChannelActivity string `json:"channel_activity"`
	XiaoMiChannelID string `json:"xiaomi_channel_id"`
	// VivoCategory            string `json:"vivo_category"`
	// OppoChannelID           string `json:"oppo_channel_id"`
	// OppoCategory            string `json:"oppo_category"`
	// OppoNotifyLevel         string `json:"oppo_notify_id"`
	HuaweiChannelImportance string `json:"huawei_channel_importance"`
	HuaweiChannelCategory   string `json:"huawei_channel_category"`
	// HonorChannelImportance  string `json:"honor_channel_importance"`
	// OppoPrivateMsgTemplate  OppoPrivateMsgTemplate `json:"oppo_private_msg_template"`
}

type OppoPrivateMsgTemplate struct {
	PrivateMsgTemplateID     string `json:"private_msg_template_id"`
	PrivateTitleParameters   string `json:"private_title_parameters"`
	PrivateContentParameters string `json:"private_content_parameters"`
}

// IOSGroupcastMessage iOS广播消息结构
type IOSGroupcastMessage struct {
	AppKey      string     `json:"appkey"`
	Timestamp   string     `json:"timestamp"`
	Type        string     `json:"type"`
	Filter      Filter     `json:"filter"`
	Payload     IOSPayload `json:"payload"`
	Policy      IOSPolicy  `json:"policy"`
	Description string     `json:"description"`
}

type IOSPayload struct {
	Aps IOSAps `json:"aps"`
}

type IOSAps struct {
	Alert             IOSAlert `json:"alert"`
	Sound             string   `json:"sound"`
	InterruptionLevel string   `json:"interruption-level"`
}

type IOSAlert struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Body     string `json:"body"`
}

// Policy 公共策略结构
type AndroidPolicy struct {
	ExpireTime               string `json:"expire_time"`
	NotificationClosedFilter bool   `json:"notification_closed_filter"`
}

type IOSPolicy struct {
	ExpireTime string `json:"expire_time"`
}

// Filter 过滤条件结构
type Filter struct {
	Where Where `json:"where"`
}

// Where 结构体，表示 where 条件
type Where struct {
	And []map[string]string `json:"and"` // 多个条件
}
