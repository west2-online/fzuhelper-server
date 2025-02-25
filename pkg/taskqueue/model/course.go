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
	"context"
	"strings"

	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

type SetCoursesCacheTask struct {
	ctx     context.Context
	cache   *cache.Cache
	userID  string
	term    string
	courses []*jwch.Course
}

func NewSetCoursesCacheTask(ctx context.Context, cache *cache.Cache, userID, term string, courses []*jwch.Course) *SetCoursesCacheTask {
	return &SetCoursesCacheTask{
		ctx:     ctx,
		cache:   cache,
		userID:  userID,
		term:    term,
		courses: courses,
	}
}

func (t *SetCoursesCacheTask) Execute() error {
	key := strings.Join([]string{t.userID, t.term}, ":")
	return t.cache.Course.SetCoursesCache(t.ctx, key, &t.courses)
}

type SetTermsCacheTask struct {
	ctx    context.Context
	cache  *cache.Cache
	userID string
	terms  []string
}

func NewSetTermsCacheTask(ctx context.Context, cache *cache.Cache, userID string, terms []string) *SetTermsCacheTask {
	return &SetTermsCacheTask{
		ctx:    ctx,
		cache:  cache,
		userID: userID,
		terms:  terms,
	}
}

func (t *SetTermsCacheTask) Execute() error {
	return t.cache.Course.SetTermsCache(t.ctx, t.userID, t.terms)
}

type PutCourseListToDatabaseTask struct {
	ctx     context.Context
	db      *db.Database
	id      string
	sf      *utils.Snowflake
	term    string
	courses []*jwch.Course
}

func NewPutCourseListToDatabaseTask(ctx context.Context, db *db.Database, id string, sf *utils.Snowflake,
	term string, courses []*jwch.Course,
) *PutCourseListToDatabaseTask {
	return &PutCourseListToDatabaseTask{
		ctx:     ctx,
		db:      db,
		id:      id,
		sf:      sf,
		term:    term,
		courses: courses,
	}
}

func (t *PutCourseListToDatabaseTask) Execute() error {
	stuId, err := utils.ParseJwchStuId(t.id)
	if err != nil {
		return err
	}

	old, err := t.db.Course.GetUserTermCourseSha256ByStuIdAndTerm(t.ctx, stuId, t.term)
	if err != nil {
		return err
	}

	json, err := utils.JSONEncode(t.courses)
	if err != nil {
		return err
	}

	newSha256 := utils.SHA256(json)

	if old == nil {
		dbId, err := t.sf.NextVal()
		if err != nil {
			return err
		}

		_, err = t.db.Course.CreateUserTermCourse(t.ctx, &model.UserCourse{
			Id:                dbId,
			StuId:             stuId,
			Term:              t.term,
			TermCourses:       json,
			TermCoursesSha256: newSha256,
		})
		if err != nil {
			return err
		}
	} else if old.TermCoursesSha256 != newSha256 {
		_, err = t.db.Course.UpdateUserTermCourse(t.ctx, &model.UserCourse{
			Id:                old.Id,
			TermCourses:       json,
			TermCoursesSha256: newSha256,
		})
		if err != nil {
			return err
		}
	}

	return nil
}

type SetLocateDateCacheTask struct {
	ctx        context.Context
	cache      *cache.Cache
	locateDate *kitexModel.LocateDate
}

func NewSetLocateDateCacheTask(ctx context.Context, cache *cache.Cache, locateDate *kitexModel.LocateDate) *SetLocateDateCacheTask {
	return &SetLocateDateCacheTask{
		ctx:        ctx,
		cache:      cache,
		locateDate: locateDate,
	}
}

func (t *SetLocateDateCacheTask) Execute() error {
	return t.cache.Course.SetLocateDateCache(t.ctx, constants.LocateDateKey, t.locateDate)
}
