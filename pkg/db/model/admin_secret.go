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

type AdminSecret struct {
	Id         int64          `gorm:"primaryKey;autoIncrement" json:"id"`
	ModuleName string         `gorm:"column:module_name;type:varchar(255);not null" json:"module_name"`
	SecretKey  string         `gorm:"column:secret_key;type:varchar(255);not null" json:"secret_key"`
	CreatedAt  time.Time      `gorm:"column:created_at;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;type:timestamp;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at,omitempty"`
}
