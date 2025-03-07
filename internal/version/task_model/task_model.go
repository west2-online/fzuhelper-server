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

package task_model

import (
	"context"
	"fmt"
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

type VersionVisitDailySyncTask struct {
	db    *db.Database
	cache *cache.Cache
}

func NewVersionVisitDailySyncTask(db *db.Database, cache *cache.Cache) *VersionVisitDailySyncTask {
	return &VersionVisitDailySyncTask{
		db:    db,
		cache: cache,
	}
}

func (t *VersionVisitDailySyncTask) Execute() error {
	ctx := context.Background()
	now := time.Now().Add(-1 * constants.ONE_DAY).In(constants.ChinaTZ)
	key := now.Format("2006-01-02")
	if !t.cache.IsKeyExist(ctx, key) {
		return fmt.Errorf("version.TaskQueue: get cache error")
	}
	visits, err := t.cache.Version.GetVisit(ctx, key)
	if err != nil {
		return fmt.Errorf("version.TaskQueue: get visits error: %w", err)
	}
	ok, _, err := t.db.Version.GetVersion(ctx, key)
	if err != nil {
		return fmt.Errorf("version.TaskQueue: get version error: %w", err)
	}
	if !ok {
		err = t.db.Version.CreateVersion(ctx, &model.Visit{
			Date:   key,
			Visits: visits,
		})
		if err != nil {
			return fmt.Errorf("version.TaskQueue: create version error: %w", err)
		}
		return nil
	} else {
		err = t.db.Version.UpdateVersion(ctx, &model.Visit{
			Date:   key,
			Visits: visits,
		})
		if err != nil {
			return fmt.Errorf("version.TaskQueue: create version error: %w", err)
		}
		return nil
	}
}

func (t *VersionVisitDailySyncTask) GetScheduleTime() time.Duration {
	now := time.Now().In(constants.ChinaTZ)
	nextRun := time.Date(now.Year(), now.Month(), now.Day(), constants.VersionVisitRefreshHour, constants.VersionVisitRefreshMinute, 0, 0, time.Local)
	if !now.Before(nextRun) {
		nextRun = nextRun.Add(constants.ONE_DAY)
	}
	// 动态处理到下一跳的刷新时间
	return nextRun.Sub(now)
}
