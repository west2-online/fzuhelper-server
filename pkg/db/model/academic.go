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

type Score struct {
	StuID            string         `json:"stu_id"`
	ScoresInfo       string         `json:"scores_info"`
	ScoresInfoSHA256 string         `json:"scores_info_sha256"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type CourseOffering struct {
	ID           int64          `json:"id"`
	Name         string         `json:"name"`
	Term         string         `json:"term"`
	Teacher      string         `json:"teacher"`
	ElectiveType string         `json:"elective_type"`
	CourseHash   string         `json:"course_hash"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at,omitempty"`
}
