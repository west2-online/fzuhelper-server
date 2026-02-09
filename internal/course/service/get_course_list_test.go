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
	"testing"
	"time"

	"github.com/bytedance/mockey"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/internal/course/pack"
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
	"github.com/west2-online/yjsy"
)

func TestGetCourseList(t *testing.T) {
	type testCase struct {
		name                 string
		mockTerms            *jwch.Term
		mockCourses          []*jwch.Course
		expectResult         []*model.Course
		expectError          string
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

	mockTerm := &jwch.Term{
		Terms:           []string{"202401"},
		ViewState:       "viewstate123",
		EventValidation: "eventvalidation123",
	}

	mockCourses := []*jwch.Course{
		{
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

	mockResult := pack.BuildCourse(mockCourses)

	// Test cases
	testCases := []testCase{
		{
			name:              "GetCourseListSuccess",
			mockTerms:         mockTerm,
			mockCourses:       mockCourses,
			expectResult:      mockResult,
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: mockCourses,
		},
		{
			name:           "GetCourseListGetTermsFailed",
			expectError:    "assert.AnError",
			mockTermsError: assert.AnError,
		},
		{
			name:              "GetCourseListInvalidTerm",
			expectError:       "Invalid term",
			mockTermsReturn:   &jwch.Term{Terms: []string{"202402", "202403"}}, // terms中不包含202401
			term:              "202401",                                        // 请求的term不在返回的terms中
			mockCoursesReturn: mockCourses,                                     // 这个不会被调用
		},
		{
			name:             "GetCourseListGetCoursesFailed",
			mockTerms:        mockTerm,
			expectError:      "assert.AnError",
			mockTermsReturn:  mockTerm,
			mockCoursesError: assert.AnError,
		},
		{
			name:           "cache exist success",
			cacheExist:     true, // 缓存里已存在
			cacheTermsList: mockTerm.Terms,
			expectResult:   mockResult,
		},
		{
			name:               "cache exist but GetTermsCache error",
			cacheExist:         true,
			cacheTermsGetError: assert.AnError,
			expectError:        "service.GetCourseList: Get term fail",
		},
		{
			name:                 "cache exist courses cache error",
			cacheExist:           true,
			cacheTermsList:       []string{"202401"},
			cacheCoursesGetError: assert.AnError,
			expectError:          "service.GetCourseList: Get courses fail",
		},
		{
			name:              "cache terms ok but term not in top2 => fallback to jwch",
			cacheExist:        true,
			cacheTermsList:    []string{"202402", "202401"},
			term:              "202399", // 不在 top2
			mockTermsReturn:   &jwch.Term{Terms: []string{"202399"}, ViewState: "v", EventValidation: "e"},
			mockCoursesReturn: mockCourses,
			expectResult:      mockResult,
		},
		{
			name:              "isRefresh=true bypass cache",
			cacheExist:        true,
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: mockCourses,
			expectResult:      mockResult,
			isRefresh:         func() *bool { b := true; return &b }(),
		},
		{
			name:              "duplicate courses are removed",
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: append([]*jwch.Course{mockCourses[0]}, mockCourses...),
			expectResult:      mockResult, // 仍为去重后的两门
		},
	}

	mockLoginData := &model.LoginData{
		Id:      "102301517",
		Cookies: "cookie1=value1; cookie2=value2",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*jwch.Student).GetTerms).Return(tc.mockTermsReturn, tc.mockTermsError).Build()

			mockey.Mock((*jwch.Student).GetSemesterCourses).Return(tc.mockCoursesReturn, tc.mockCoursesError).Build()

			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()
			if tc.cacheExist {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).To(
					func(ctx context.Context, key string) ([]string, error) {
						if tc.cacheTermsList != nil {
							return tc.cacheTermsList, tc.cacheTermsGetError
						}
						return mockTerm.Terms, tc.cacheTermsGetError
					},
				).Build()

				mockey.Mock((*coursecache.CacheCourse).GetCoursesCache).Return(mockCourses, tc.cacheCoursesGetError).Build()
			} else {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).Return(nil, assert.AnError).Build()
			}

			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			// 每个用例可自定义 term 和 isRefresh
			term := "202401"
			if tc.term != "" {
				term = tc.term
			}
			req := &course.CourseListRequest{Term: term}
			if tc.isRefresh != nil {
				req.IsRefresh = tc.isRefresh
			}

			ctx := customContext.WithLoginData(context.Background(), mockLoginData)
			courseService := NewCourseService(ctx, mockClientSet, new(taskqueue.BaseTaskQueue))
			result, err := courseService.GetCourseList(req, mockLoginData)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}

func TestGetCourseListYjsy(t *testing.T) {
	type testCase struct {
		name                 string
		mockTermsReturn      *yjsy.Term
		mockTermsError       error
		mockCoursesReturn    []*yjsy.Course
		mockCoursesError     error
		expectResult         []*model.Course
		expectError          string
		cacheExist           bool
		cacheTermsGetError   error
		cacheCoursesGetError error
		cacheTermsList       []string
		term                 string
		isRefresh            *bool
	}

	mockTerm := &yjsy.Term{
		Terms: []string{"202401"},
	}

	mockCourses := []*yjsy.Course{
		{
			Name:    "Mathematics",
			Teacher: "Prof. John",
			ScheduleRules: []yjsy.CourseScheduleRule{
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
			ScheduleRules: []yjsy.CourseScheduleRule{
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

	mockResult := pack.BuildCourseYjsy(mockCourses)

	// Test cases
	testCases := []testCase{
		{
			name:           "YjsyCacheHitSuccess",
			cacheExist:     true,
			cacheTermsList: mockTerm.Terms,
			expectResult:   mockResult,
		},
		{
			name:               "YjsyCacheTermsError",
			cacheExist:         true,
			cacheTermsGetError: assert.AnError,
			expectError:        "service.GetCourseListYjsy: Get terms fail",
		},
		{
			name:                 "YjsyCacheCoursesError",
			cacheExist:           true,
			cacheTermsList:       []string{"202401"},
			cacheCoursesGetError: assert.AnError,
			expectError:          "service.GetCourseListYjsy: Get courses fail",
		},
		{
			name:           "YjsyGetTermsFailed",
			mockTermsError: assert.AnError,
			expectError:    "assert.AnError",
		},
		{
			name:              "YjsyInvalidTerm",
			mockTermsReturn:   &yjsy.Term{Terms: []string{"202402", "202403"}}, // terms中不包含202401
			term:              "202401",                                        // 请求的term不在返回的terms中
			expectError:       "Invalid term",
			mockCoursesReturn: mockCourses, // 这个不会被调用
		},
		{
			name:             "YjsyGetCoursesFailed",
			mockTermsReturn:  mockTerm,
			mockCoursesError: assert.AnError,
			expectError:      "assert.AnError",
		},
		{
			name:              "YjsyCacheNotTop2Fallback",
			cacheExist:        true,
			cacheTermsList:    []string{"202402", "202401"},
			term:              "202399", // 不在 top2，需要回源
			mockTermsReturn:   &yjsy.Term{Terms: []string{"202399"}},
			mockCoursesReturn: mockCourses,
			expectResult:      mockResult,
		},
		{
			name:              "YjsyIsRefreshBypassCache",
			cacheExist:        true,
			mockTermsReturn:   mockTerm,
			mockCoursesReturn: mockCourses,
			expectResult:      mockResult,
			isRefresh:         func() *bool { b := true; return &b }(),
		},
	}

	mockLoginData := &model.LoginData{
		Id:      "102301517",
		Cookies: "cookie1=value1; cookie2=value2",
	}

	defer mockey.UnPatchAll()
	for _, tc := range testCases {
		mockey.PatchConvey(tc.name, t, func() {
			mockClientSet := &base.ClientSet{
				SFClient:    new(utils.Snowflake),
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
			}

			mockey.Mock((*yjsy.Student).GetTerms).Return(tc.mockTermsReturn, tc.mockTermsError).Build()

			mockey.Mock((*yjsy.Student).GetSemesterCourses).Return(tc.mockCoursesReturn, tc.mockCoursesError).Build()

			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()
			if tc.cacheExist {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).To(
					func(ctx context.Context, key string) ([]string, error) {
						if tc.cacheTermsList != nil {
							return tc.cacheTermsList, tc.cacheTermsGetError
						}
						return mockTerm.Terms, tc.cacheTermsGetError
					},
				).Build()

				mockey.Mock((*coursecache.CacheCourse).GetCoursesCacheYjsy).Return(mockCourses, tc.cacheCoursesGetError).Build()
			} else {
				mockey.Mock((*coursecache.CacheCourse).GetTermsCache).Return(nil, assert.AnError).Build()
			}

			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			term := "202401"
			if tc.term != "" {
				term = tc.term
			}
			req := &course.CourseListRequest{Term: term}
			if tc.isRefresh != nil {
				req.IsRefresh = tc.isRefresh
			}

			ctx := customContext.WithLoginData(context.Background(), mockLoginData)
			courseService := NewCourseService(ctx, mockClientSet, new(taskqueue.BaseTaskQueue))
			result, err := courseService.GetCourseListYjsy(req, mockLoginData)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectResult, result)
			}
		})
	}
}

func TestGetSemesterCourses(t *testing.T) {
	type testCase struct {
		name               string
		cacheExist         bool
		cacheGetError      error
		dbReturnNil        bool
		dbGetError         error
		dbTermCoursesValue string
		expectError        string
		expectResult       []*model.Course
	}

	stuID := "102301001"
	term := "202401"

	jwchCourses := []*jwch.Course{
		{Name: "A", Teacher: "T1", ScheduleRules: []jwch.CourseScheduleRule{{StartClass: 1, EndClass: 2, StartWeek: 1, EndWeek: 16}}},
	}

	courseB := []*model.Course{{Name: "B", Teacher: "T2"}}

	// Test cases
	testCases := []testCase{
		{
			name:       "GetSemesterCoursesCacheHitSuccess",
			cacheExist: true,
		},
		{
			name:          "GetSemesterCoursesCacheGetError",
			cacheExist:    true,
			cacheGetError: assert.AnError,
			expectError:   "service.GetSemesterCourses: Get courses fail",
		},
		{
			name:        "GetSemesterCoursesDbError",
			dbGetError:  assert.AnError,
			expectError: "service.GetSemesterCourses: Get courses fail",
		},
		{
			name:        "GetSemesterCoursesDbReturnNil",
			dbReturnNil: true,
			expectError: "there is no course in database",
		},
		{
			name:               "GetSemesterCoursesDbUnmarshalFail",
			dbTermCoursesValue: "{",
			expectError:        "Unmarshal fail",
		},
		{
			name:               "GetSemesterCoursesDbSuccess",
			dbTermCoursesValue: "",
			expectResult:       courseB,
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

			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.cacheExist).Build()
			if tc.cacheExist {
				mockey.Mock((*coursecache.CacheCourse).GetCoursesCache).Return(jwchCourses, tc.cacheGetError).Build()
			} else {
				mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseByStuIdAndTerm).To(
					func(ctx context.Context, stuIdArg string, termArg string) (*dbmodel.UserCourse, error) {
						if tc.dbReturnNil {
							return nil, tc.dbGetError
						}
						termCoursesValue := tc.dbTermCoursesValue
						if termCoursesValue == "" {
							jsonStr, _ := utils.JSONEncode(courseB)
							termCoursesValue = jsonStr
						}
						return &dbmodel.UserCourse{TermCourses: termCoursesValue}, tc.dbGetError
					},
				).Build()
			}

			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			courseService := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			res, err := courseService.getSemesterCourses(stuID, term)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
				if tc.cacheExist {
					assert.Equal(t, len(jwchCourses), len(res))
					assert.Equal(t, "A", res[0].Name)
				} else {
					assert.Equal(t, tc.expectResult, res)
				}
			}
		})
	}
}

func TestCourseToDatabase(t *testing.T) {
	type testCase struct {
		name            string
		getSha256Return *dbmodel.UserCourse
		getSha256Error  error
		encodeError     error
		nextValReturn   int64
		nextValError    error
		createError     error
		updateError     error
		expectError     string
	}

	stuId := "102301517"
	term := "202401"
	courses := []*model.Course{{Name: "C", Teacher: "T3"}}

	// 计算课程的实际 SHA256
	coursesJSON, _ := utils.JSONEncode(courses)
	newSha := utils.SHA256(coursesJSON)

	// Test cases
	testCases := []testCase{
		{
			name:           "GetSha256Error",
			getSha256Error: assert.AnError,
			expectError:    "assert.AnError",
		},
		{
			name:        "EncodeError",
			encodeError: assert.AnError,
			expectError: "assert.AnError",
		},
		{
			name:          "CreateNewCourseSuccess",
			nextValReturn: int64(123),
		},
		{
			name:         "CreateNewCourseNextValError",
			nextValError: assert.AnError,
			expectError:  "assert.AnError",
		},
		{
			name:          "CreateNewCourseCreateError",
			nextValReturn: int64(123),
			createError:   assert.AnError,
			expectError:   "assert.AnError",
		},
		{
			name:            "UpdateCourseSameShaNoUpdate",
			getSha256Return: &dbmodel.UserCourse{Id: 1, TermCoursesSha256: newSha},
		},
		{
			name:            "UpdateCourseDifferentShaSuccess",
			getSha256Return: &dbmodel.UserCourse{Id: 2, TermCoursesSha256: "oldsha"},
		},
		{
			name:            "UpdateCourseUpdateError",
			getSha256Return: &dbmodel.UserCourse{Id: 2, TermCoursesSha256: "oldsha"},
			updateError:     assert.AnError,
			expectError:     "assert.AnError",
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

			mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseSha256ByStuIdAndTerm).Return(tc.getSha256Return, tc.getSha256Error).Build()

			mockey.Mock(utils.JSONEncode).Return("", tc.encodeError).Build()

			mockey.Mock((*utils.Snowflake).NextVal).Return(tc.nextValReturn, tc.nextValError).Build()

			mockey.Mock((*dbcourse.DBCourse).CreateUserTermCourse).Return(nil, tc.createError).Build()

			mockey.Mock((*dbcourse.DBCourse).UpdateUserTermCourse).Return(nil, tc.updateError).Build()

			mockey.Mock((*taskqueue.BaseTaskQueue).Add).Return().Build()

			courseService := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			err := courseService.putCourseToDatabase(stuId, term, courses)

			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestHandleCourseUpdate(t *testing.T) {
	type testCase struct {
		name              string
		newList           []*model.Course
		old               *dbmodel.UserCourse
		sendError         error
		expectNotifyCount int
		expectError       string
	}
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

	testCases := []testCase{
		{
			name:              "adjust changed -> notify once",
			newList:           newList,
			old:               old,
			expectNotifyCount: 1,
		},
		{
			name:        "sendNotifications error",
			newList:     newList,
			old:         old,
			sendError:   assert.AnError,
			expectError: "Send notifications failed",
		},
		{
			name:              "no change -> no notify",
			newList:           oldList,
			old:               old,
			expectNotifyCount: 0,
		},
		{
			name:        "old json invalid -> error",
			newList:     oldList,
			old:         &dbmodel.UserCourse{TermCourses: "{"},
			expectError: "Unmarshal old courses failed",
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

			mockey.Mock((*CourseService).sendNotifications).Return(tc.sendError).Build()

			courseService := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			err := courseService.handleCourseUpdate(term, tc.newList, tc.old)
			if tc.expectError != "" {
				assert.ErrorContains(t, err, tc.expectError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestSendNotifications(t *testing.T) {
	type testCase struct {
		name         string
		androidError error
		iosError     error
		expectError  bool
	}

	courseName := "CourseX"
	tag := "tag-1"

	testCases := []testCase{
		{
			name: "both platform success",
		},
		{
			name:         "android fail",
			androidError: assert.AnError,
			expectError:  true,
		},
		{
			name:        "ios fail after android ok",
			iosError:    assert.AnError,
			expectError: true,
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

			_ = config.InitForTest("course")

			mockey.Mock(time.Sleep).To(func(d time.Duration) {}).Build()
			mockey.Mock(umeng.SendAndroidGroupcastWithGoApp).Return(tc.androidError).Build()
			mockey.Mock(umeng.SendIOSGroupcast).Return(tc.iosError).Build()

			courseService := NewCourseService(context.Background(), mockClientSet, new(taskqueue.BaseTaskQueue))
			err := courseService.sendNotifications(courseName, tag)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
