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

type ToolboxConfig struct {
	Id        int64          `json:"id"`
	ToolID    int64          `json:"tool_id"`
	Visible   bool           `json:"visible"`
	Name      string         `json:"name"`
	Icon      string         `json:"icon"`
	Type      string         `json:"type"`
	Message   string         `json:"message,omitempty"`
	Extra     string         `json:"extra"`
	StudentID string         `json:"student_id,omitempty"`
	Platform  string         `json:"platform,omitempty"`
	Version   int64          `json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}
