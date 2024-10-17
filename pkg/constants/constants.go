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
	MuxConnection  = 1
	RPCTimeout     = 3 * time.Second
	ConnectTimeout = 50 * time.Millisecond

	// service name
	TemplateServiceName  = "template"
	ClassroomServiceName = "classroom"
	UserServiceName      = "user"
	ApiServiceName       = "api"
	PaperServiceName     = "paper"

	// db table name
	TemplateServiceTableName = "template"
	CheckServiceTableName    = "check"

	// redis
	RedisDBEmptyRoom   = 0
	RedisDBPaper       = 1
	ClassroomKeyExpire = 24 * time.Hour
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

	NumWorkers = 10 // 最大的并发数量

	// timeout
	FailureRateLimiterBaseDelay = time.Minute
	FailureRateLimiterMaxDelay  = 30 * time.Minute

	// 定时任务
	ScheduledTime = 24 * time.Hour

	//又拍云
	CACHE_FILEDIR = "UssFileDir"
	CacheDst      = ".cache/"
)

var CampusArray = []string{"旗山校区", "厦门工艺美院", "铜盘校区", "怡山校区", "晋江校区", "泉港校区"}
