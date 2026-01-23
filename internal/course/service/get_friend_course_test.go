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
	"errors"
	"strings"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/stretchr/testify/assert"

	"github.com/west2-online/fzuhelper-server/internal/course/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	customContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	coursecache "github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	dbcourse "github.com/west2-online/fzuhelper-server/pkg/db/course"
	dbmodel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

// mockUserClient implements userservice.Client for tests.
// 仅实现 VerifyFriend，其他方法在表驱动用例中不会被调用。
type mockUserClient struct {
	verifyResp *user.VerifyFriendResponse
	verifyErr  error
}

func (m *mockUserClient) VerifyFriend(ctx context.Context, req *user.VerifyFriendRequest, opts ...callopt.Option) (*user.VerifyFriendResponse, error) {
	return m.verifyResp, m.verifyErr
}

// unused client methods
func (m *mockUserClient) GetLoginData(context.Context, *user.GetLoginDataRequest, ...callopt.Option) (*user.GetLoginDataResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserClient) GetUserInfo(context.Context, *user.GetUserInfoRequest, ...callopt.Option) (*user.GetUserInfoResponse, error) {
	return nil, errors.New("not implemented")
}

//nolint:lll
func (m *mockUserClient) GetGetLoginDataForYJSY(context.Context, *user.GetLoginDataForYJSYRequest, ...callopt.Option) (*user.GetLoginDataForYJSYResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserClient) GetInvitationCode(context.Context, *user.GetInvitationCodeRequest, ...callopt.Option) (*user.GetInvitationCodeResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserClient) BindInvitation(context.Context, *user.BindInvitationRequest, ...callopt.Option) (*user.BindInvitationResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserClient) GetFriendList(context.Context, *user.GetFriendListRequest, ...callopt.Option) (*user.GetFriendListResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserClient) DeleteFriend(context.Context, *user.DeleteFriendRequest, ...callopt.Option) (*user.DeleteFriendResponse, error) {
	return nil, errors.New("not implemented")
}

func (m *mockUserClient) CancelInvite(context.Context, *user.CancelInviteRequest, ...callopt.Option) (*user.CancelInviteResponse, error) {
	return nil, errors.New("not implemented")
}

func TestCourseService_GetFriendCourse(t *testing.T) {
	defer mockey.UnPatchAll()

	baseRespOK := &kitexModel.BaseResp{Code: errno.SuccessCode, Msg: "ok"}
	loginData := &kitexModel.LoginData{Id: "stu-1", Cookies: "ck"}

	type testCase struct {
		name            string
		verifyResp      *user.VerifyFriendResponse
		verifyErr       error
		isKeyExistFn    func(string) bool
		termsCache      []string
		termsCacheErr   error
		termsDB         *dbmodel.UserTerm
		termsDBErr      error
		coursesCache    []*jwch.Course
		coursesCacheErr error
		coursesYjsy     []*yjsy.Course
		coursesYjsyErr  error
		dbCourse        *dbmodel.UserCourse
		dbCourseErr     error
		reqTerm         string
		expectErr       bool
		errContains     string
		expectLen       int
		expectFirstName string
	}

	cases := []testCase{
		{
			name:        "verify friend rpc error",
			verifyErr:   errors.New("rpc fail"),
			reqTerm:     "202401",
			expectErr:   true,
			errContains: "verify friend failed",
		},
		{
			name:        "HandleBaseRespWithCookie error",
			verifyResp:  &user.VerifyFriendResponse{Base: &kitexModel.BaseResp{Code: errno.ParamErrorCode, Msg: "bad param"}},
			reqTerm:     "202401",
			expectErr:   true,
			errContains: "bad param",
		},
		{
			name:        "not friend",
			verifyResp:  &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: false},
			reqTerm:     "202401",
			expectErr:   true,
			errContains: "只能查看好友的课表",
		},
		{
			name:          "terms cache error",
			verifyResp:    &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn:  func(key string) bool { return strings.HasPrefix(key, "terms:") },
			termsCacheErr: errors.New("cache read error"),
			reqTerm:       "202401",
			expectErr:     true,
			errContains:   "Get term fail",
		},
		{
			name:         "terms db error",
			verifyResp:   &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn: func(string) bool { return false },
			termsDBErr:   errors.New("db error"),
			reqTerm:      "202401",
			expectErr:    true,
			errContains:  "Get term from database fail",
		},
		{
			name:         "terms empty",
			verifyResp:   &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn: func(string) bool { return false },
			termsDB:      (*dbmodel.UserTerm)(nil),
			reqTerm:      "202401",
			expectErr:    true,
			errContains:  "Friend termList empty",
		},
		{
			name:         "invalid term - not in terms",
			verifyResp:   &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn: func(string) bool { return false },
			termsDB:      &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401", "202402"})},
			reqTerm:      "202501",
			expectErr:    true,
			errContains:  "Invalid term",
		},
		{
			name:         "invalid term - yjsy format but mapped term not in terms",
			verifyResp:   &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn: func(string) bool { return false },
			termsDB:      &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202301", "202302"})}, // 只有202301和202302
			reqTerm:      "2024-2025-1",                                                                 // yjsy格式，映射到202401，不在terms中
			expectErr:    true,
			errContains:  "Invalid term",
		},
		{
			name:         "invalid term - jwch format but mapped term not in terms",
			verifyResp:   &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn: func(string) bool { return false },
			termsDB:      &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"2023-2024-1", "2023-2024-2"})}, // 只有yjsy格式的terms
			reqTerm:      "202401",                                                                                // jwch格式，映射到2024-2025-1，不在terms中
			expectErr:    true,
			errContains:  "Invalid term",
		},
		{
			name:         "invalid term - default case (unrecognized format)",
			verifyResp:   &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn: func(string) bool { return false },
			termsDB:      &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401"})},
			reqTerm:      "invalid-format", // 既不是yjsy也不是jwch格式，且不在terms中
			expectErr:    true,
			errContains:  "Invalid term",
		},
		{
			name:            "yjsy term format mapping success",
			verifyResp:      &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn:    func(key string) bool { return strings.HasPrefix(key, "terms:") || strings.HasPrefix(key, "course:") },
			termsCache:      []string{"202401", "202402"}, // 包含映射后的"202401"
			coursesCache:    []*jwch.Course{{Name: "Math", Teacher: "T1"}},
			reqTerm:         "2024-2025-1", // yjsy格式，MapYjsyTerm映射为"202401"
			expectLen:       1,
			expectErr:       false,
			expectFirstName: "Math",
		},
		{
			name:            "jwch term format mapping success",
			verifyResp:      &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn:    func(key string) bool { return strings.HasPrefix(key, "terms:") || strings.HasPrefix(key, "course:") },
			termsCache:      []string{"2024-2025-1", "2024-2025-2"}, // 包含映射后的"2024-2025-1"
			coursesCache:    []*jwch.Course{{Name: "Physics", Teacher: "T2"}},
			reqTerm:         "202401", // jwch格式，MapJwchTerm映射为"2024-2025-1"
			expectLen:       1,
			expectErr:       false,
			expectFirstName: "Physics",
		},
		{
			name:            "term cache hit and course cache hit",
			verifyResp:      &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn:    func(key string) bool { return strings.HasPrefix(key, "terms:") || strings.HasPrefix(key, "course:") },
			termsCache:      []string{"202401", "202402"},
			coursesCache:    []*jwch.Course{{Name: "Math", Teacher: "T1"}},
			reqTerm:         "202401",
			expectLen:       1,
			expectFirstName: "Math",
		},
		{
			name:            "course cache error",
			verifyResp:      &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn:    func(key string) bool { return strings.HasPrefix(key, "terms:") || strings.HasPrefix(key, "course:") },
			termsCache:      []string{"202401"},
			coursesCacheErr: errors.New("cache course error"),
			reqTerm:         "202401",
			expectErr:       true,
			errContains:     "Get courses fail",
		},
		{
			name:            "course cache empty -> yjsy cache fallback",
			verifyResp:      &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn:    func(key string) bool { return strings.HasPrefix(key, "terms:") || strings.HasPrefix(key, "course:") },
			termsCache:      []string{"202401"},
			coursesCache:    []*jwch.Course(nil),
			coursesYjsy:     []*yjsy.Course{{Name: "Algo", Teacher: "T2"}},
			reqTerm:         "202401",
			expectLen:       1,
			expectFirstName: "Algo",
		},
		{
			name:           "yjsy cache error",
			verifyResp:     &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn:   func(key string) bool { return strings.HasPrefix(key, "terms:") || strings.HasPrefix(key, "course:") },
			termsCache:     []string{"202401"},
			coursesCache:   []*jwch.Course(nil),
			coursesYjsyErr: errors.New("yjsy cache error"),
			reqTerm:        "202401",
			expectErr:      true,
			errContains:    "Get courses fail",
		},
		{
			name:         "term from db not in recent two",
			verifyResp:   &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn: func(string) bool { return false },
			termsDB:      &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401", "202302", "202201"})},
			reqTerm:      "202201",
			expectErr:    true,
			errContains:  "只能查看好友最近两个学期的课表",
		},
		{
			name:         "db course error",
			verifyResp:   &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn: func(string) bool { return false },
			termsDB:      &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401"})},
			dbCourseErr:  errors.New("db query error"),
			reqTerm:      "202401",
			expectErr:    true,
			errContains:  "Get courses fail",
		},
		{
			name:         "db course missing",
			verifyResp:   &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn: func(string) bool { return false },
			termsDB:      &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401"})},
			dbCourse:     (*dbmodel.UserCourse)(nil),
			reqTerm:      "202401",
			expectErr:    true,
			errContains:  "there is no course in database",
		},
		{
			name:         "db course unmarshal error",
			verifyResp:   &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn: func(string) bool { return false },
			termsDB:      &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401"})},
			dbCourse:     &dbmodel.UserCourse{TermCourses: "{"},
			reqTerm:      "202401",
			expectErr:    true,
			errContains:  "Unmarshal fail",
		},
		{
			name:            "db course success",
			verifyResp:      &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true},
			isKeyExistFn:    func(string) bool { return false },
			termsDB:         &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401"})},
			dbCourse:        &dbmodel.UserCourse{TermCourses: `[{"name":"Database"}]`},
			reqTerm:         "202401",
			expectErr:       false,
			expectLen:       1,
			expectFirstName: "Database",
		},
	}

	for _, tc := range cases {
		mockey.PatchConvey(tc.name, t, func() {
			userCli := &mockUserClient{verifyResp: tc.verifyResp, verifyErr: tc.verifyErr}

			clientSet := &base.ClientSet{
				DBClient:    &db.Database{Course: new(dbcourse.DBCourse)},
				CacheClient: &cache.Cache{Course: new(coursecache.CacheCourse)},
				UserClient:  userCli,
			}

			mockey.Mock((*cache.Cache).IsKeyExist).To(func(_ context.Context, key string) bool {
				if tc.isKeyExistFn != nil {
					return tc.isKeyExistFn(key)
				}
				return false
			}).Build()

			mockey.Mock((*coursecache.CacheCourse).GetTermsCache).To(func(_ context.Context, _ string) ([]string, error) {
				return tc.termsCache, tc.termsCacheErr
			}).Build()
			mockey.Mock((*coursecache.CacheCourse).GetCoursesCache).To(func(_ context.Context, _ string) ([]*jwch.Course, error) {
				return tc.coursesCache, tc.coursesCacheErr
			}).Build()
			mockey.Mock((*coursecache.CacheCourse).GetCoursesCacheYjsy).To(func(_ context.Context, _ string) ([]*yjsy.Course, error) {
				return tc.coursesYjsy, tc.coursesYjsyErr
			}).Build()

			mockey.Mock((*dbcourse.DBCourse).GetUserTermByStuId).To(func(_ *dbcourse.DBCourse, _ context.Context, _ string) (*dbmodel.UserTerm, error) {
				return tc.termsDB, tc.termsDBErr
			}).Build()
			mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseByStuIdAndTerm).To(
				func(_ *dbcourse.DBCourse, _ context.Context, _ string, _ string) (*dbmodel.UserCourse, error) {
					return tc.dbCourse, tc.dbCourseErr
				}).Build()

			ctx := customContext.WithLoginData(context.Background(), loginData)
			svc := NewCourseService(ctx, clientSet, nil)

			res, err := svc.GetFriendCourse(&course.GetFriendCourseRequest{Id: "f1", Term: tc.reqTerm}, loginData)

			if tc.expectErr {
				assert.Error(t, err)
				if tc.errContains != "" {
					assert.Contains(t, err.Error(), tc.errContains)
				}
				assert.Nil(t, res)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, res, tc.expectLen)
			if tc.expectLen > 0 {
				assert.Equal(t, tc.expectFirstName, res[0].Name)
			}
		})
	}
}
