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

package constants

import "time"

const (
	UmengURL               = "https://msgapi.umeng.com/api/send" // 推送 API
	UmengMessageExpireTime = 3 * ONE_DAY                         // 推送消息过期时间
	UmengRateLimitDelay    = 5 * time.Second                     // 用于在成绩 diff 循环中等待，防止被友盟限流
)

// Tag
const (
	UmengJwchNoticeTag = "jwch-notice" // 教务处通知的tag
)
