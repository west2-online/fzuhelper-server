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
	JWTValue = "MTAxNTkwMTg1Mw=="

	MuxConnection           = 1                     // (RPC) 最大连接数
	RPCTimeout              = 3 * time.Second       // (RPC) RPC请求超时时间
	ConnectTimeout          = 50 * time.Millisecond // (RPC) 连接超时时间
	StreamBufferSize        = 1024                  // (RPC) 流请求 Buffer 尺寸
	CheckFileTypeBufferSize = 512                   // 判断文件类型时需读取前512个字节

	TemplateServiceName      = "template"
	ClassroomServiceName     = "classroom"
	CourseServiceName        = "course"
	UserServiceName          = "user"
	ApiServiceName           = "api"
	LaunchScreenServiceName  = "launch_screen"
	PaperServiceName         = "paper"
	URLServiceName           = "url"
	AcademicServiceName      = "academic"
	TemplateServiceTableName = "template"
	UserTableName            = "user"
	LaunchScreenTableName    = "launch_screen"
	CourseTableName          = "course"

	RedisDBEmptyRoom      = 0
	RedisDBLaunchScreen   = 1
	RedisDBPaper          = 2
	ClassroomKeyExpire    = 2 * 24 * time.Hour
	LaunchScreenKeyExpire = 2 * 24 * time.Hour
	LastLaunchScreenIdKey = "last_launch_screen_id"
	RedisSlowQuery        = 10 // ms redis默认的慢查询时间

	SnowflakeWorkerID     = 0
	SnowflakeDatacenterID = 0

	MaxQPS          = 100
	MaxVideoSize    = 300000
	MaxListLength   = 100
	MaxGoroutines   = 10
	MaxOpenConns    = 100
	MaxConnections  = 1000             // (DB) 最大连接数
	MaxIdleConns    = 10               // (DB) 最大空闲连接数
	ConnMaxLifetime = 10 * time.Second // (DB) 最大可复用时间
	ConnMaxIdleTime = 5 * time.Minute  // (DB) 最长保持空闲状态时间
	PageSize        = 10

	NumWorkers = 10 // 最大的并发数量

	FailureRateLimiterBaseDelay = time.Minute
	FailureRateLimiterMaxDelay  = 30 * time.Minute

	ClassroomWorker        = 1              // (class_room) 同时启用的 goroutine 数量
	ClassroomScheduledTime = 24 * time.Hour // (class_room) 空教室非当天同步时间
	ClassroomUpdatedTime   = 6 * time.Hour  // (class_room) 当天空教室更新间隔

	MaxRetries   = 5               // 最大重试次数
	InitialDelay = 1 * time.Second // 初始等待时间

	CACHE_FILEDIR = "UssFileDir"

	KafkaReadMinBytes      = 512 * B
	KafkaReadMaxBytes      = 1 * MB
	KafkaRetries           = 3
	DefaultReaderGroupID   = "r"
	DefaultTimeRetainHours = 6 // 6小时

	// 获取验证码结果的本地python服务url，需要保证 login-verify 和 api 处于同一个 dokcer 网络中
	ValidateCodeURL = "http://login-verify:8081/api/v1/jwch/user/validateCode"

	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB

	DayTime        = 24
	TimeZoneOffset = 8
)

const (
	OneDay = 24 * time.Hour
)

var CampusArray = []string{"旗山校区", "厦门工艺美院", "铜盘校区", "怡山校区", "晋江校区", "泉港校区"}
