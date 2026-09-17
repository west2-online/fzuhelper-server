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
	"sync"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/west2-online/fzuhelper-server/kitex_gen/common"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	rpcmodel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	coursecache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbcourse "github.com/west2-online/fzuhelper-server/pkg/db/course"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func TestGetAutoAdjustCourseList(t *testing.T) {
	type testCase struct {
		name                  string
		mockTerm              string
		mockAutoAdjustCourses []*model.AutoAdjustCourse
		mockDBErr             error
		cacheExist            bool
		cacheGetErr           error
		expectResult          []*model.AutoAdjustCourse
		expectError           string
	}

	mockTerm := "202501"

	mockAutoAdjustCourses := []*model.AutoAdjustCourse{
		{
			Year:        "2025",
			FromDate:    "2025-05-01",
			ToDate:      new("2025-05-05"),
			Term:        mockTerm,
			FromWeek:    18,
			ToWeek:      new(int64(19)),
			FromWeekday: 1,
			ToWeekday:   new(int64(3)),
			Enabled:     true,
		},
		{
			Year:        "2025",
			FromDate:    "2025-06-04",
			ToDate:      new("2025-06-07"),
			Term:        mockTerm,
			FromWeek:    20,
			ToWeek:      new(int64(2)),
			FromWeekday: 2,
			ToWeekday:   new(int64(4)),
			Enabled:     true,
		},
	}

	testCases := []testCase{
		{
			name:                  "success",
			mockTerm:              mockTerm,
			mockAutoAdjustCourses: mockAutoAdjustCourses,
			expectResult:          mockAutoAdjustCourses,
		},
		{
			name:                  "cache exist success",
			mockTerm:              mockTerm,
			mockAutoAdjustCourses: mockAutoAdjustCourses,
			cacheExist:            true,
			expectResult:          mockAutoAdjustCourses,
		},
		{
			name:        "cache exist but get cache error",
			mockTerm:    mockTerm,
			cacheExist:  true,
			cacheGetErr: assert.AnError,
			expectError: "service.GetAutoAdjustCourseList: Get cache failed",
		},
		{
			name:        "cache miss get from db error",
			mockTerm:    mockTerm,
			cacheExist:  false,
			mockDBErr:   assert.AnError,
			expectError: "service.GetAutoAdjustCourseList: Get from db failed",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			shouldWait := !tc.cacheExist && tc.mockDBErr == nil
			var wg sync.WaitGroup
			if shouldWait {
				wg.Add(1)
			}

			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()

			if tc.cacheExist {
				mockey.Mock((*coursecache.CacheCourse).GetAutoAdjustCourseListCache).
					Return(tc.mockAutoAdjustCourses, tc.cacheGetErr).
					Build()
			}

			mockey.Mock((*dbcourse.DBCourse).GetAutoAdjustCourseListByTerm).
				Return(tc.mockAutoAdjustCourses, tc.mockDBErr).
				Build()

			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			courseService := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			result, err := courseService.GetAutoAdjustCourseList(tc.mockTerm)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}

func TestUpdateAutoAdjustCourse(t *testing.T) {
	type testCase struct {
		name            string
		req             *course.UpdateAdjustCourseRequest
		mockCheckPwd    bool
		mockOriginal    *model.AutoAdjustCourse
		mockOriginalErr error
		mockUpdateErr   error
		termResp        *common.TermListResponse
		termErr         error
		expectError     string
	}

	boolPtr := func(b bool) *bool { return &b }

	mockID := int64(1)
	mockTerm := "202501"
	mockTermStr := mockTerm

	mockOriginal := &model.AutoAdjustCourse{
		Id:          mockID,
		Term:        mockTerm,
		Year:        "2025",
		FromDate:    "2025-05-01",
		FromWeek:    18,
		FromWeekday: 1,
		Enabled:     true,
	}

	successBase := &rpcmodel.BaseResp{Code: errno.SuccessCode, Msg: "ok"}
	errorBase := &rpcmodel.BaseResp{Code: errno.InternalServiceErrorCode, Msg: "internal error"}

	successTermResp := &common.TermListResponse{
		Base: successBase,
		TermLists: &rpcmodel.TermList{
			CurrentTerm: &mockTermStr,
			Terms: []*rpcmodel.Term{
				{
					Term:      &mockTermStr,
					StartDate: new("2025-02-17"),
					EndDate:   new("2025-06-30"),
				},
			},
		},
	}

	testCases := []testCase{
		{
			name: "invalid secret",
			req: &course.UpdateAdjustCourseRequest{
				Id:     mockID,
				Secret: "wrong",
			},
			mockCheckPwd: false,
			expectError:  "invalid admin secret",
		},
		{
			name: "success update enabled only",
			req: &course.UpdateAdjustCourseRequest{
				Id:      mockID,
				Secret:  "secret",
				Enabled: boolPtr(true),
			},
			mockCheckPwd: true,
			mockOriginal: mockOriginal,
		},
		{
			name: "get original record failed",
			req: &course.UpdateAdjustCourseRequest{
				Id:      mockID,
				Secret:  "secret",
				Enabled: boolPtr(false),
			},
			mockCheckPwd:    true,
			mockOriginalErr: assert.AnError,
			expectError:     "service.UpdateAutoAdjustCourse: Get original record failed",
		},
		{
			name: "update db failed",
			req: &course.UpdateAdjustCourseRequest{
				Id:      mockID,
				Secret:  "secret",
				Enabled: boolPtr(false),
			},
			mockCheckPwd:  true,
			mockOriginal:  mockOriginal,
			mockUpdateErr: assert.AnError,
			expectError:   "service.UpdateAutoAdjustCourse: Update failed",
		},
		{
			name: "get terms list rpc failed",
			req: &course.UpdateAdjustCourseRequest{
				Id:       mockID,
				Secret:   "secret",
				FromDate: new("2025-05-01"),
			},
			mockCheckPwd: true,
			termErr:      assert.AnError,
			expectError:  "service.UpdateAutoAdjustCourse: Get terms list failed",
		},
		{
			name: "terms list base resp error",
			req: &course.UpdateAdjustCourseRequest{
				Id:       mockID,
				Secret:   "secret",
				FromDate: new("2025-05-01"),
			},
			mockCheckPwd: true,
			termResp: &common.TermListResponse{
				Base: errorBase,
			},
			expectError: "service.UpdateAutoAdjustCourse: term list resp error",
		},
		{
			name: "no term found for from_date",
			req: &course.UpdateAdjustCourseRequest{
				Id:       mockID,
				Secret:   "secret",
				FromDate: new("2024-01-01"),
			},
			mockCheckPwd: true,
			termResp:     successTermResp,
			expectError:  "no term found for date",
		},
		{
			name: "success with from_date update",
			req: &course.UpdateAdjustCourseRequest{
				Id:       mockID,
				Secret:   "secret",
				FromDate: new("2025-05-01"),
			},
			mockCheckPwd: true,
			termResp:     successTermResp,
			mockOriginal: mockOriginal,
		},
		{
			name: "success with to_date empty cancellation",
			req: &course.UpdateAdjustCourseRequest{
				Id:     mockID,
				Secret: "secret",
				ToDate: new(""),
			},
			mockCheckPwd: true,
			termResp:     successTermResp,
			mockOriginal: mockOriginal,
		},
		{
			name: "success with to_date set",
			req: &course.UpdateAdjustCourseRequest{
				Id:     mockID,
				Secret: "secret",
				ToDate: new("2025-06-04"),
			},
			mockCheckPwd: true,
			termResp:     successTermResp,
			mockOriginal: mockOriginal,
		},
		{
			name: "no term found for to_date",
			req: &course.UpdateAdjustCourseRequest{
				Id:     mockID,
				Secret: "secret",
				ToDate: new("2024-01-01"),
			},
			mockCheckPwd: true,
			termResp:     successTermResp,
			expectError:  "no term found for to_date",
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock(utils.CheckPwd).Return(tc.mockCheckPwd).Build()
			mockey.Mock((*dbcourse.DBCourse).GetAutoAdjustCourseByID).Return(tc.mockOriginal, tc.mockOriginalErr).Build()
			mockey.Mock((*dbcourse.DBCourse).UpdateAutoAdjustCourse).Return(tc.mockUpdateErr).Build()
			mockey.Mock((*dbcourse.DBCourse).GetAutoAdjustCourseListByTerm).Return(nil, nil).Build()
			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()
			mockey.Mock(utils.GetWeekdayByDate).Return(18, 1, nil).Build()

			courseService := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			courseService.commonClient = &mockCommonClient{
				termResp: tc.termResp,
				termErr:  tc.termErr,
			}

			err := courseService.UpdateAutoAdjustCourse(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestBuildAutoAdjustCourse(t *testing.T) {
	mockTermStr := "202501"
	terms := []*rpcmodel.Term{
		{
			Term:      &mockTermStr,
			StartDate: new("2025-02-17"),
			EndDate:   new("2025-06-30"),
		},
	}

	type testCase struct {
		name        string
		item        *course.CreateAdjustCourseItem
		expectError string
		expect      *model.AutoAdjustCourse
	}

	testCases := []testCase{
		{
			name: "success with to_date set",
			item: &course.CreateAdjustCourseItem{FromDate: "2025-05-01", ToDate: new("2025-05-08")},
			expect: &model.AutoAdjustCourse{
				Year:        "2025",
				FromDate:    "2025-05-01",
				ToDate:      new("2025-05-08"),
				Term:        mockTermStr,
				FromWeek:    11,
				FromWeekday: 4,
				ToWeek:      new(int64(12)),
				ToWeekday:   new(int64(4)),
				Enabled:     false,
			},
		},
		{
			name: "success with to_date empty means canceled",
			item: &course.CreateAdjustCourseItem{FromDate: "2025-05-01", ToDate: nil},
			expect: &model.AutoAdjustCourse{
				Year:        "2025",
				FromDate:    "2025-05-01",
				ToDate:      nil,
				Term:        mockTermStr,
				FromWeek:    11,
				FromWeekday: 4,
				ToWeek:      nil,
				ToWeekday:   nil,
				Enabled:     false,
			},
		},
		{
			name:        "invalid from_date",
			item:        &course.CreateAdjustCourseItem{FromDate: "not-a-date"},
			expectError: "invalid from_date",
		},
		{
			name:        "no term found for from_date",
			item:        &course.CreateAdjustCourseItem{FromDate: "2024-05-01"},
			expectError: "no term found for from_date",
		},
		{
			name:        "invalid to_date",
			item:        &course.CreateAdjustCourseItem{FromDate: "2025-05-01", ToDate: new("not-a-date")},
			expectError: "invalid to_date",
		},
		{
			name:        "no term found for to_date",
			item:        &course.CreateAdjustCourseItem{FromDate: "2025-05-01", ToDate: new("2024-05-08")},
			expectError: "no term found for to_date",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := buildAutoAdjustCourse(tc.item, terms)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
				assert.Nil(t, result)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, result)
		})
	}
}

func TestCreateAutoAdjustCourse(t *testing.T) {
	type testCase struct {
		name           string
		req            *course.CreateAdjustCourseRequest
		termResp       *common.TermListResponse
		termErr        error
		createErrs     []error // 按调用次序返回，超出长度后视为成功
		expectError    string
		expectCount    int64
		expectRefresh  []string
		expectRPCalled bool
	}

	mockTermStr := "202501"
	successBase := &rpcmodel.BaseResp{Code: errno.SuccessCode, Msg: "ok"}
	successTermResp := &common.TermListResponse{
		Base: successBase,
		TermLists: &rpcmodel.TermList{
			CurrentTerm: &mockTermStr,
			Terms: []*rpcmodel.Term{
				{
					Term:      &mockTermStr,
					StartDate: new("2025-02-17"),
					EndDate:   new("2025-06-30"),
				},
			},
		},
	}
	errorTermResp := &common.TermListResponse{
		Base: &rpcmodel.BaseResp{Code: errno.InternalServiceErrorCode, Msg: "internal error"},
	}
	nilTermListResp := &common.TermListResponse{
		Base:      successBase,
		TermLists: nil,
	}

	testCases := []testCase{
		{
			name:        "empty items does not call rpc",
			req:         &course.CreateAdjustCourseRequest{Items: nil},
			expectCount: 0,
		},
		{
			name:        "term list rpc error",
			req:         &course.CreateAdjustCourseRequest{Items: []*course.CreateAdjustCourseItem{{FromDate: "2025-05-01"}}},
			termErr:     assert.AnError,
			expectError: "Get terms list failed",
		},
		{
			name:        "term list resp error",
			req:         &course.CreateAdjustCourseRequest{Items: []*course.CreateAdjustCourseItem{{FromDate: "2025-05-01"}}},
			termResp:    errorTermResp,
			expectError: "term list resp error",
		},
		{
			name:        "empty rpc response",
			req:         &course.CreateAdjustCourseRequest{Items: []*course.CreateAdjustCourseItem{{FromDate: "2025-05-01"}}},
			termResp:    nil,
			expectError: "empty rpc response",
		},
		{
			name:        "nil term list",
			req:         &course.CreateAdjustCourseRequest{Items: []*course.CreateAdjustCourseItem{{FromDate: "2025-05-01"}}},
			termResp:    nilTermListResp,
			expectError: "term list is nil",
		},
		{
			name: "success with two items",
			req: &course.CreateAdjustCourseRequest{Items: []*course.CreateAdjustCourseItem{
				{FromDate: "2025-05-01", ToDate: new("2025-05-08")},
				{FromDate: "2025-05-02"},
			}},
			termResp:       successTermResp,
			expectCount:    2,
			expectRefresh:  []string{"refreshAutoAdjustCourseCache:202501"},
			expectRPCalled: true,
		},
		{
			name: "invalid item is skipped but valid one is created",
			req: &course.CreateAdjustCourseRequest{Items: []*course.CreateAdjustCourseItem{
				{FromDate: "not-a-date"},
				{FromDate: "2024-05-01"},
				{FromDate: "2025-05-01"},
			}},
			termResp:       successTermResp,
			expectCount:    1,
			expectRefresh:  []string{"refreshAutoAdjustCourseCache:202501"},
			expectRPCalled: true,
		},
		{
			name: "duplicated from_date is skipped",
			req: &course.CreateAdjustCourseRequest{Items: []*course.CreateAdjustCourseItem{
				{FromDate: "2025-05-01"},
			}},
			termResp:       successTermResp,
			createErrs:     []error{gorm.ErrDuplicatedKey},
			expectCount:    0,
			expectRPCalled: true,
		},
		{
			name: "create failed",
			req: &course.CreateAdjustCourseRequest{Items: []*course.CreateAdjustCourseItem{
				{FromDate: "2025-05-01"},
			}},
			termResp:       successTermResp,
			createErrs:     []error{assert.AnError},
			expectError:    "create failed",
			expectRPCalled: true,
		},
		{
			// 中途失败时，已写入的记录同样要刷新缓存，否则管理端在 TTL 内看不到它们
			name: "partial failure still refreshes cache for written terms",
			req: &course.CreateAdjustCourseRequest{Items: []*course.CreateAdjustCourseItem{
				{FromDate: "2025-05-01"},
				{FromDate: "2025-05-02"},
			}},
			termResp:       successTermResp,
			createErrs:     []error{nil, assert.AnError},
			expectError:    "create failed",
			expectCount:    1,
			expectRefresh:  []string{"refreshAutoAdjustCourseCache:202501"},
			expectRPCalled: true,
		},
		{
			// 全部条目都被跳过时不应产生缓存刷新
			name: "all items skipped refreshes nothing",
			req: &course.CreateAdjustCourseRequest{Items: []*course.CreateAdjustCourseItem{
				{FromDate: "not-a-date"},
			}},
			termResp:       successTermResp,
			expectCount:    0,
			expectRPCalled: false,
		},
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			var refreshKeys []string
			var createCalled bool
			var createCalls int

			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*dbcourse.DBCourse).CreateAutoAdjustCourse).
				To(func(ctx context.Context, m *model.AutoAdjustCourse) (*model.AutoAdjustCourse, error) {
					createCalled = true
					var err error
					if createCalls < len(tc.createErrs) {
						err = tc.createErrs[createCalls]
					}
					createCalls++
					return m, err
				}).Build()
			mockey.Mock((*dbcourse.DBCourse).GetAutoAdjustCourseListByTerm).Return(nil, nil).Build()
			mockey.Mock((*taskqueue.BaseTaskQueue).Add).
				To(func(key string, task taskqueue.QueueTask) {
					refreshKeys = append(refreshKeys, key)
				}).Build()

			courseService := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			courseService.commonClient = &mockCommonClient{
				termResp: tc.termResp,
				termErr:  tc.termErr,
			}

			count, err := courseService.CreateAutoAdjustCourse(tc.req)

			if tc.expectError != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, tc.expectCount, count)
			assert.Equal(t, tc.expectRefresh, refreshKeys)
			assert.Equal(t, tc.expectRPCalled, createCalled)
		})
	}
}
