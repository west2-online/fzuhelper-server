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

type NetworkEnv string

const (
	Network2G      NetworkEnv = "2G"
	Network3G      NetworkEnv = "3G"
	Network4G      NetworkEnv = "4G"
	Network5G      NetworkEnv = "5G"
	NetworkWifi    NetworkEnv = "wifi"
	NetworkUnknown NetworkEnv = "unknown"
)

type Feedback struct {
	ReportId     int64
	StuId        string
	Name         string
	College      string
	ContactPhone string
	ContactQQ    string
	ContactEmail string

	NetworkEnv   NetworkEnv
	IsOnCampus   bool
	OsName       string
	OsVersion    string
	Manufacturer string
	DeviceModel  string

	ProblemDesc    string
	Screenshots    string `json:"scores_info"`
	AppVersion     string
	VersionHistory string `json:"version_history"`
	NetworkTraces  string `json:"network_traces"`
	Events         string `json:"events"`
	UserSettings   string `json:"user_settings"`

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

type FeedbackListReq struct {
	StuId       string
	Name        string
	NetworkEnv  NetworkEnv
	IsOnCampus  *bool
	OsName      string
	ProblemDesc string
	AppVersion  string

	Limit     int
	PageToken int64
	OrderDesc *bool

	BeginTime *time.Time
	EndTime   *time.Time
}

type FeedbackListItem struct {
	ReportId    int64
	Name        string
	NetworkEnv  NetworkEnv
	ProblemDesc string
	AppVersion  string

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
