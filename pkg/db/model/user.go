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

type Student struct {
	StuId     string `gorm:"primary_key"`
	Name      string
	Sex       string
	Birthday  string
	College   string
	Grade     int64
	Major     string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}
type FollowRelation struct {
	Id         int64  `gorm:"primary_key"`
	FollowerId string // 关注者
	FollowedId string // 被关注者
	OrderSeq   int64  // 排序序号，越大越靠前
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  gorm.DeletedAt `sql:"index"`
	// ActiveFlag 是数据库层的 GENERATED VIRTUAL 列，根据 IF(deleted_at IS NULL, 1, NULL) 求值，只读（已经标了 <-:false
	// active_flag = 1 等价于 deleted_at IS NULL，可利用 uk_active_relation 唯一索引加速查询。
	ActiveFlag *int8 `gorm:"column:active_flag;<-:false"`
}
type UserFriend struct {
	FriendId  string `gorm:"column:followed_id"`
	OrderSeq  int64
	CreatedAt time.Time
}
