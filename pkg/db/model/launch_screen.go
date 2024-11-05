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

package model

import (
	"time"

	"gorm.io/gorm"
)

type Picture struct {
	ID         int64
	Url        string
	Href       string
	Text       string
	PicType    int64
	ShowTimes  int64
	PointTimes int64
	Duration   int64
	StartAt    time.Time // 开始时间
	EndAt      time.Time // 结束时间
	StartTime  int64     // 开始时段 0~24
	EndTime    int64     // 结束时段 0~24
	SType      int64     // 类型
	Frequency  int64     // 一天展示次数
	Regex      string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `sql:"index"`
}
