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

// Config
const (
	MaxConnections  = 1000            // (DB) 最大连接数
	MaxIdleConns    = 10              // (DB) 最大空闲连接数
	ConnMaxLifetime = 10 * ONE_SECOND // (DB) 最大可复用时间
	ConnMaxIdleTime = 5 * ONE_MINUTE  // (DB) 最长保持空闲状态时间

)

// Table Name
const (
	UserTableName            = "student"
	UserRelationTableName    = " follow_relation"
	CourseTableName          = "course"
	LaunchScreenTableName    = "launch_screen"
	NoticeTableName          = "notice"
	ScoreTableName           = "scores"
	VisitTableName           = "visit"
	CourseOfferingsTableName = "course_offerings"
	ToolboxConfigTableName   = "toolbox_config"
	AdminSecretTableName     = "admin_secrets"
	FeedbackTableName        = "feedback"
)

// Biz
const (
	StuInfoExpireTime = ONE_WEEK // 存储在db的学生信息最大刷新时间
)
