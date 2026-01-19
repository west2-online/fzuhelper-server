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
	"fmt"
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	customContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	coursecache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbcourse "github.com/west2-online/fzuhelper-server/pkg/db/course"
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/umeng"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func TestCourseService_GetCourseList(t *testing.T) {
	mockTerm := &jwch.Term{
		Terms:           []string{"202401"},
		ViewState:       "viewstate123",
		EventValidation: "eventvalidation123",
	}

	mockCourses := []*jwch.Course{
		{
			Type:    "Required",
			Name:    "Mathematics",
			Credits: "3.0",
			Teacher: "Prof. John",
			ScheduleRules: []jwch.CourseScheduleRule{
				{
					Location:     "A-202",
					StartClass:   2,
					EndClass:     4,
					StartWeek:    1,
					EndWeek:      16,
					Weekday:      1,
					Single:       false,
					Double:       true,
					Adjust:       false,
					FromFullWeek: false,
				},
			},
		},
		{
			Type:    "Elective",
			Name:    "Physics",
			Credits: "3.0",
			Teacher: "Prof. Smith",
			ScheduleRules: []jwch.CourseScheduleRule{
				{
					Location:     "A-203",
					StartClass:   3,
					EndClass:     4,
					StartWeek:    2,
					EndWeek:      17,
					Weekday:      2,
					Single:       false,
					Double:       true,
					Adjust:       false,
					FromFullWeek: false,
				},
			},
		},
	}

	mockResult := []*model.Course{
		{
			Name:    "Mathematics",
			Teacher: "Prof. John",
			ScheduleRules: []*model.CourseScheduleRule{
				{
					Location:   "A-202",
					StartClass: 2,
					EndClass:   4,
					StartWeek:  1,
					EndWeek:    16,
					Weekday:    1,
					Single:     false,
					Double:     true,
					Adjust:     false,
				},
			},
		},
		{
			Name:    "Physics",
			Teacher: "Prof. Smith",
			ScheduleRules: []*model.CourseScheduleRule{
				{
					Location:   "A-203",
					StartClass: 3,
					EndClass:   4,
					StartWeek:  2,
					EndWeek:    17,
					Weekday:    2,
					Single:     false,
					Double:     true,
					Adjust:     false,
				},
			},
		},
	}
	type testCase struct {
		name                 string
		mockTerms            *jwch.Term
		mockCourses          []*jwch.Course
		expectedResult       []*model.Course
		expectingError       bool
		expectedErrorMsg     string
		mockTermsReturn      *jwch.Term
		mockTermsError       error
		mockCoursesReturn    []*jwch.Course
		mockCoursesError     error
		cacheExist           bool
		cacheTermsGetError   error
		cacheCoursesGetError error
		cacheTermsList       []string
		term                 string
		isRefresh            *bool
	}

	// Test cases
	testCases := []testCase{
		{
			name:              "GetCourseListSuccess",
			mockTerms:         mockTerm,
			mockCourses:       mockCourses,
			expectedResult:    mockResult,
			expectingError:    false,
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: mockCourses,
		},
		{
			name:              "GetCourseListInvalidTerm",
			mockTerms:         mockTerm,
			mockCourses:       nil,
			expectedResult:    nil,
			expectingError:    true,
			expectedErrorMsg:  "Invalid term",
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: nil,
			mockCoursesError:  fmt.Errorf("Invalid term"),
		},
		{
			name:             "GetCourseListGetTermsFailed",
			mockTerms:        nil,
			mockCourses:      nil,
			expectedResult:   nil,
			expectingError:   true,
			expectedErrorMsg: "Get terms failed",
			mockTermsReturn:  nil,
			mockTermsError:   fmt.Errorf("Get terms failed"),
		},
		{
			name:              "GetCourseListGetCoursesFailed",
			mockTerms:         mockTerm,
			mockCourses:       nil,
			expectedResult:    nil,
			expectingError:    true,
			expectedErrorMsg:  "Get semester courses failed",
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: nil,
			mockCoursesError:  fmt.Errorf("Get semester courses failed"),
		},
		{
			name:                 "cache exist success",
			cacheExist:           true, // 缓存里已存在
			cacheTermsGetError:   nil,  // 获取缓存不报错
			cacheCoursesGetError: nil,
			cacheTermsList:       mockTerm.Terms,
			expectedResult:       mockResult,
		},
		{
			name:               "cache exist but GetTermsCache error",
			cacheExist:         true,
			cacheTermsGetError: fmt.Errorf("terms cache error"),
			expectedResult:     nil,
			expectingError:     true,
			expectedErrorMsg:   "service.GetCourseList: Get term fail",
		},
		{
			name:                 "cache exist courses cache error",
			cacheExist:           true,
			cacheTermsList:       []string{"202401"},
			cacheCoursesGetError: fmt.Errorf("courses cache error"),
			expectedResult:       nil,
			expectingError:       true,
			expectedErrorMsg:     "service.GetCourseList: Get courses fail",
		},
		{
			name:              "cache terms ok but term not in top2 => fallback to jwch",
			cacheExist:        true,
			cacheTermsList:    []string{"202402", "202401"},
			term:              "202399", // 不在 top2
			mockTermsReturn:   &jwch.Term{Terms: []string{"202399"}, ViewState: "v", EventValidation: "e"},
			mockCoursesReturn: mockCourses,
			expectedResult:    mockResult,
		},
		{
			name:              "isRefresh=true bypass cache",
			cacheExist:        true,
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: mockCourses,
			expectedResult:    mockResult,
			isRefresh:         func() *bool { b := true; return &b }(),
		},
		{
			name:              "duplicate courses are removed",
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: append([]*jwch.Course{mockCourses[0]}, mockCourses...),
			expectedResult:    mockResult, // 仍为去重后的两门
		},
	}

	mockLoginData := &model.LoginData{
		Id:      "102301517",
		Cookies: "cookie1=value1; cookie2=value2",
	}

	defer mockey.UnPatchAll()

	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockey.Mock((*jwch.Student).GetTerms).Return(tc.mockTermsReturn, tc.mockTermsError).Build()
			mockey.Mock((*jwch.Student).GetSemesterCourses).Return(tc.mockCoursesReturn, tc.mockCoursesError).Build()
			mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool {
				return tc.cacheExist
			}).Build()
			if tc.cacheExist {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).To(
					func(ctx context.Context, key string) ([]string, error) {
						if tc.cacheTermsGetError != nil {
							return nil, tc.cacheTermsGetError
						}
						if tc.cacheTermsList != nil {
							return tc.cacheTermsList, nil
						}
						return mockTerm.Terms, nil
					},
				).Build()
				mockey.Mock((*coursecache.CacheCourse).GetCoursesCache).To(
					func(ctx context.Context, key string) ([]*jwch.Course, error) {
						if tc.cacheCoursesGetError != nil {
							return nil, tc.cacheCoursesGetError
						}
						return mockCourses, nil
					},
				).Build()
			} else {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).To(
					func(ctx context.Context, key string) ([]string, error) {
						return nil, fmt.Errorf("should not be called if cache doesn't exist")
					},
				).Build()
			}
			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			mockClientSet := new(base.ClientSet)
			mockClientSet.SFClient = new(utils.Snowflake)
			mockClientSet.DBClient = new(db.Database)
			mockClientSet.CacheClient = new(cache.Cache)

			ctx := customContext.WithLoginData(context.Background(), mockLoginData)
			courseService := NewCourseService(ctx, mockClientSet, new(taskqueue.BaseTaskQueue))

			// 每个用例可自定义 term 和 isRefresh
			term := "202401"
			if tc.term != "" {
				term = tc.term
			}
			req := &course.CourseListRequest{Term: term}
			if tc.isRefresh != nil {
				req.IsRefresh = tc.isRefresh
			}

			result, err := courseService.GetCourseList(req, &model.LoginData{Id: "123456789", Cookies: "cookie1=value1;cookie2=value2"})

			if tc.expectingError {
				assert.Nil(t, result)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedErrorMsg)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedResult, result)
			}
		})
	}
}

func TestCourseService_RemoveDuplicateCourses(t *testing.T) {
	s := &CourseService{}
	courses := []*model.Course{
		{
			Name:    "Math",
			Teacher: "Alice",
			ScheduleRules: []*model.CourseScheduleRule{
				{StartClass: 1, EndClass: 2, StartWeek: 1, EndWeek: 16},
				{StartClass: 3, EndClass: 4, StartWeek: 1, EndWeek: 16},
			},
		},
		// 重复（规则顺序不同）
		{
			Name:    "Math",
			Teacher: "Alice",
			ScheduleRules: []*model.CourseScheduleRule{
				{StartClass: 3, EndClass: 4, StartWeek: 1, EndWeek: 16},
				{StartClass: 1, EndClass: 2, StartWeek: 1, EndWeek: 16},
			},
		},
		// 不同教师，应保留
		{
			Name:    "Math",
			Teacher: "Bob",
			ScheduleRules: []*model.CourseScheduleRule{
				{StartClass: 1, EndClass: 2, StartWeek: 1, EndWeek: 16},
			},
		},
	}

	got := s.removeDuplicateCourses(courses)
	assert.Equal(t, 2, len(got))
	// 保证包含不同教师那门
	hasBob := false
	for _, c := range got {
		if c.Teacher == "Bob" {
			hasBob = true
		}
	}
	assert.True(t, hasBob)
}

func TestCourseService_GetSemesterCourses(t *testing.T) {
	defer mockey.UnPatchAll()

	// 构造基础依赖
	mockClientSet := new(base.ClientSet)
	mockClientSet.SFClient = new(utils.Snowflake)
	mockClientSet.DBClient = new(db.Database)
	mockClientSet.CacheClient = new(cache.Cache)
	s := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))

	stuID := "102301001"
	term := "202401"

	jwchCourses := []*jwch.Course{
		{Name: "A", Teacher: "T1", ScheduleRules: []jwch.CourseScheduleRule{{StartClass: 1, EndClass: 2, StartWeek: 1, EndWeek: 16}}},
	}

	// 1) 缓存命中成功
	mockey.PatchConvey("cache hit success", t, func() {
		mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool { return true }).Build()
		mockey.Mock((*coursecache.CacheCourse).GetCoursesCache).Return(jwchCourses, nil).Build()
		res, err := s.getSemesterCourses(stuID, term)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, "A", res[0].Name)
	})

	// 2) 缓存读取失败
	mockey.PatchConvey("cache get error", t, func() {
		mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool { return true }).Build()
		mockey.Mock((*coursecache.CacheCourse).GetCoursesCache).Return(nil, fmt.Errorf("cache error")).Build()
		res, err := s.getSemesterCourses(stuID, term)
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "service.GetSemesterCourses: Get courses fail")
	})

	// 3) DB 查询报错
	mockey.PatchConvey("db error", t, func() {
		mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool { return false }).Build()
		mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseByStuIdAndTerm).Return(nil, fmt.Errorf("db error")).Build()
		res, err := s.getSemesterCourses(stuID, term)
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "service.GetSemesterCourses: Get courses fail")
	})

	// 4) DB 返回空
	mockey.PatchConvey("db return nil", t, func() {
		mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool { return false }).Build()
		mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseByStuIdAndTerm).Return((*dbmodel.UserCourse)(nil), nil).Build()
		res, err := s.getSemesterCourses(stuID, term)
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "there is no course in database")
	})

	// 5) Unmarshal 失败
	mockey.PatchConvey("db unmarshal fail", t, func() {
		mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool { return false }).Build()
		mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseByStuIdAndTerm).Return(&dbmodel.UserCourse{TermCourses: "{"}, nil).Build()
		res, err := s.getSemesterCourses(stuID, term)
		assert.Nil(t, res)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unmarshal fail")
	})

	// 6) DB 成功
	mockey.PatchConvey("db success", t, func() {
		list := []*model.Course{{Name: "B", Teacher: "T2"}}
		jsonStr, _ := utils.JSONEncode(list)
		mockey.Mock((*cache.Cache).IsKeyExist).To(func(ctx context.Context, key string) bool { return false }).Build()
		mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseByStuIdAndTerm).Return(&dbmodel.UserCourse{TermCourses: jsonStr}, nil).Build()
		// 避免异步 cache 任务的副作用
		mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()
		res, err := s.getSemesterCourses(stuID, term)
		assert.Nil(t, err)
		assert.Equal(t, 1, len(res))
		assert.Equal(t, "B", res[0].Name)
	})
}

func TestCourseService_PutCourseToDatabase(t *testing.T) {
	defer mockey.UnPatchAll()

	mockClientSet := new(base.ClientSet)
	mockClientSet.SFClient = new(utils.Snowflake)
	mockClientSet.DBClient = new(db.Database)
	mockClientSet.CacheClient = new(cache.Cache)
	tq := new(taskqueue.BaseTaskQueue)
	s := NewCourseService(context.Background(), mockClientSet, tq)

	stuId := "102301517"
	term := "202401"
	courses := []*model.Course{{Name: "C", Teacher: "T3"}}

	// 1) 无旧记录 => Create
	mockey.PatchConvey("no old -> create", t, func() {
		var createCalled int
		mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseSha256ByStuIdAndTerm).Return((*dbmodel.UserCourse)(nil), nil).Build()
		mockey.Mock((*utils.Snowflake).NextVal).Return(int64(123), nil).Build()
		mockey.Mock((*dbcourse.DBCourse).CreateUserTermCourse).To(func(_ *dbcourse.DBCourse, _ context.Context, uc *dbmodel.UserCourse) (*dbmodel.UserCourse, error) {
			createCalled++
			assert.Equal(t, stuId, uc.StuId)
			assert.Equal(t, term, uc.Term)
			return uc, nil
		}).Build()
		err := s.putCourseToDatabase(stuId, term, courses)
		assert.Nil(t, err)
		assert.Equal(t, 1, createCalled)
	})

	// 2) 旧记录相同 SHA => 不更新、不推送
	mockey.PatchConvey("same sha -> no update", t, func() {
		// 计算与实现一致的 sha
		jsonStr, _ := utils.JSONEncode(courses)
		sameSha := utils.SHA256(jsonStr)
		old := &dbmodel.UserCourse{Id: 1, TermCoursesSha256: sameSha}
		var updateCalled, addCalled int
		mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseSha256ByStuIdAndTerm).Return(old, nil).Build()
		mockey.Mock((*dbcourse.DBCourse).UpdateUserTermCourse).To(func(_ *dbcourse.DBCourse, _ context.Context, uc *dbmodel.UserCourse) (*dbmodel.UserCourse, error) {
			updateCalled++
			return uc, nil
		}).Build()
		mockey.Mock((*taskqueue.BaseTaskQueue).Add).To(func(_ *taskqueue.BaseTaskQueue, _ string, _ taskqueue.QueueTask) { addCalled++ }).Build()
		err := s.putCourseToDatabase(stuId, term, courses)
		assert.Nil(t, err)
		assert.Equal(t, 0, updateCalled)
		assert.Equal(t, 0, addCalled)
	})

	// 3) 旧记录不同 SHA => 更新并异步任务
	mockey.PatchConvey("diff sha -> update and task", t, func() {
		old := &dbmodel.UserCourse{Id: 2, TermCoursesSha256: "oldsha"}
		var updateCalled, addCalled int
		mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseSha256ByStuIdAndTerm).Return(old, nil).Build()
		mockey.Mock((*dbcourse.DBCourse).UpdateUserTermCourse).To(func(_ *dbcourse.DBCourse, _ context.Context, uc *dbmodel.UserCourse) (*dbmodel.UserCourse, error) {
			updateCalled++
			return uc, nil
		}).Build()
		mockey.Mock((*taskqueue.BaseTaskQueue).Add).To(func(_ *taskqueue.BaseTaskQueue, _ string, _ taskqueue.QueueTask) { addCalled++ }).Build()
		err := s.putCourseToDatabase(stuId, term, courses)
		assert.Nil(t, err)
		assert.Equal(t, 1, updateCalled)
		assert.Equal(t, 1, addCalled)
	})
}

func TestCourseService_HandleCourseUpdate(t *testing.T) {
	defer mockey.UnPatchAll()
	s := &CourseService{}

	term := "202401"
	// 旧数据：RawAdjust = ""
	oldList := []*model.Course{{
		Name:             "C",
		Teacher:          "T",
		ElectiveType:     "elective",
		RawScheduleRules: "[]",
		RawAdjust:        "",
	}}
	oldJSON, _ := utils.JSONEncode(oldList)
	old := &dbmodel.UserCourse{TermCourses: oldJSON}

	// 新数据：RawAdjust 改变
	newList := []*model.Course{{
		Name:             "C",
		Teacher:          "T",
		ElectiveType:     "elective",
		RawScheduleRules: "[]",
		RawAdjust:        "1",
	}}

	mockey.PatchConvey("adjust changed -> notify once", t, func() {
		var notifyCnt int
		mockey.Mock((*CourseService).sendNotifications).To(func(_ *CourseService, courseName, tag string) error {
			notifyCnt++
			assert.Equal(t, "C", courseName)
			return nil
		}).Build()
		err := s.handleCourseUpdate(term, newList, old)
		assert.Nil(t, err)
		assert.Equal(t, 1, notifyCnt)
	})

	mockey.PatchConvey("no change -> no notify", t, func() {
		var notifyCnt int
		mockey.Mock((*CourseService).sendNotifications).To(func(_ *CourseService, courseName, tag string) error { notifyCnt++; return nil }).Build()
		// 新旧相同
		err := s.handleCourseUpdate(term, oldList, old)
		assert.Nil(t, err)
		assert.Equal(t, 0, notifyCnt)
	})

	mockey.PatchConvey("old json invalid -> error", t, func() {
		badOld := &dbmodel.UserCourse{TermCourses: "{"}
		err := s.handleCourseUpdate(term, oldList, badOld)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Unmarshal old courses failed")
	})
}

func TestCourseService_SendNotifications(t *testing.T) {
	defer mockey.UnPatchAll()
	s := &CourseService{}

	_ = config.InitForTest("api")
	mockey.Mock(time.Sleep).To(func(d time.Duration) {}).Build()

	mockey.PatchConvey("both platform success", t, func() {
		mockey.Mock(umeng.SendAndroidGroupcastWithGoApp).Return(nil).Build()
		mockey.Mock(umeng.SendIOSGroupcast).Return(nil).Build()
		err := s.sendNotifications("CourseX", "tag-1")
		assert.Nil(t, err)
	})

	mockey.PatchConvey("android fail", t, func() {
		mockey.Mock(umeng.SendAndroidGroupcastWithGoApp).Return(fmt.Errorf("android err")).Build()
		err := s.sendNotifications("CourseX", "tag-1")
		assert.Error(t, err)
	})

	mockey.PatchConvey("ios fail after android ok", t, func() {
		mockey.Mock(umeng.SendAndroidGroupcastWithGoApp).Return(nil).Build()
		mockey.Mock(umeng.SendIOSGroupcast).Return(fmt.Errorf("ios err")).Build()
		err := s.sendNotifications("CourseX", "tag-1")
		assert.Error(t, err)
	})
}
