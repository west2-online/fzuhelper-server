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
	// auth
	JWTValue = "MTAxNTkwMTg1Mw=="
	StartID  = 10000

	// RPC
	MuxConnection    = 1
	RPCTimeout       = 3 * time.Second
	ConnectTimeout   = 50 * time.Millisecond
	StreamBufferSize = 1024

	// service name
	TemplateServiceName     = "template"
	ClassroomServiceName    = "classroom"
	CourseServiceName       = "course"
	UserServiceName         = "user"
	ApiServiceName          = "api"
	LaunchScreenServiceName = "launch_screen"
	PaperServiceName        = "paper"
	AcademicServiceName     = "academic"

	// db table name
	TemplateServiceTableName = "template"
	UserTableName            = "user"
	LaunchScreenTableName    = "launch_screen"
	CourseTableName          = "course"

	// redis
	RedisDBEmptyRoom      = 0
	RedisDBLaunchScreen   = 1
	RedisDBPaper          = 2
	ClassroomKeyExpire    = 2 * 24 * time.Hour
	LaunchScreenKeyExpire = 2 * 24 * time.Hour
	LastLaunchScreenIdKey = "last_launch_screen_id"
	// snowflake
	SnowflakeWorkerID     = 0
	SnowflakeDatacenterID = 0

	// limit
	MaxConnections  = 1000
	MaxQPS          = 100
	MaxVideoSize    = 300000
	MaxListLength   = 100
	MaxIdleConns    = 10
	MaxGoroutines   = 10
	MaxOpenConns    = 100
	ConnMaxLifetime = 10 * time.Second
	PageSize        = 10

	NumWorkers = 10 // 最大的并发数量

	// timeout
	FailureRateLimiterBaseDelay = time.Minute
	FailureRateLimiterMaxDelay  = 30 * time.Minute

	// 定时任务
	ScheduledTime = 24 * time.Hour
	UpdatedTime   = 6 * time.Hour // 当天空教室更新间隔

	// retry
	MaxRetries   = 5               // 最大重试次数
	InitialDelay = 1 * time.Second // 初始等待时间

	// 又拍云
	CACHE_FILEDIR = "UssFileDir"

	// Kafka
	KafkaReadMinBytes      = 512 * B
	KafkaReadMaxBytes      = 1 * MB
	KafkaRetries           = 3
	DefaultReaderGroupID   = "r"
	DefaultTimeRetainHours = 6 // 6小时

	// byte
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

var CampusArray = []string{"旗山校区", "厦门工艺美院", "铜盘校区", "怡山校区", "晋江校区", "泉港校区"}
