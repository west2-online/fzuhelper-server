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

package service

import (
	"context"
	"net/http"
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

type CreateFeedbackReq struct {
	StuId        string `json:"stuId"`
	Name         string `json:"name"`
	College      string `json:"college"`
	ContactPhone string `json:"contactPhone"`
	ContactQQ    string `json:"contactQQ"`
	ContactEmail string `json:"contactEmail"`

	NetworkEnv   string `json:"networkEnv"`
	IsOnCampus   bool   `json:"isOnCampus"`
	OsName       string `json:"osName"`
	OsVersion    string `json:"osVersion"`
	Manufacturer string `json:"manufacturer"`
	DeviceModel  string `json:"deviceModel"`

	ProblemDesc    string `json:"problemDesc"`
	Screenshots    string `json:"screenshots"`
	AppVersion     string `json:"appVersion"`
	VersionHistory string `json:"versionHistory"`

	NetworkTraces string `json:"networkTraces"`
	Events        string `json:"events"`
	UserSettings  string `json:"userSettings"`
}

type FeedbackListReq struct {
	StuId       string `json:"stuId"`
	Name        string `json:"name"`
	NetworkEnv  string `json:"networkEnv"`
	IsOnCampus  *bool  `json:"isOnCampus"`
	OsName      string `json:"osName"`
	ProblemDesc string `json:"problemDesc"`
	AppVersion  string `json:"appVersion"`

	Limit     int   `json:"limit"`
	PageToken int64 `json:"pageToken"`
	OrderDesc *bool `json:"orderDesc"`

	BeginTime *time.Time `json:"beginTime"`
	EndTime   *time.Time `json:"endTime"`
}

type OAService struct {
	ctx   context.Context
	db    *db.Database
	cache *cache.Cache
	sf    *utils.Snowflake
}

func NewOAService(ctx context.Context, identifier string, cookies []*http.Cookie, clientset *base.ClientSet) *OAService {
	return &OAService{
		ctx:   ctx,
		db:    clientset.DBClient,
		cache: clientset.CacheClient,
		sf:    clientset.SFClient,
	}
}
