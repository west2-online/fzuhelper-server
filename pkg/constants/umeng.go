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
