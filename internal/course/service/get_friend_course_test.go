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

func TestGetFriendCourse(t *testing.T) {
	type testCase struct {
		name            string
		verifyResp      *user.VerifyFriendResponse
		verifyErr       error
		isKeyExistFn    bool
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
		expectErr       string
		expectLen       int
		expectFirstName string
	}

	baseRespOK := &kitexModel.BaseResp{Code: errno.SuccessCode, Msg: "ok"}
	baseRespErr := &kitexModel.BaseResp{Code: errno.ParamErrorCode, Msg: "bad param"}
	baseResponseOK := &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: true}
	baseResponseNotFriend := &user.VerifyFriendResponse{Base: baseRespOK, FriendExist: false}
	loginData := &kitexModel.LoginData{Id: "stu-1", Cookies: "ck"}

	cases := []testCase{
		{
			name:      "verify friend rpc error",
			verifyErr: assert.AnError,
			reqTerm:   "202401",
			expectErr: "verify friend failed",
		},
		{
			name:       "HandleBaseRespWithCookie error",
			verifyResp: &user.VerifyFriendResponse{Base: baseRespErr},
			reqTerm:    "202401",
			expectErr:  "bad param",
		},
		{
			name:       "not friend",
			verifyResp: baseResponseNotFriend,
			reqTerm:    "202401",
			expectErr:  "只能查看好友的课表",
		},
		{
			name:          "terms cache error",
			verifyResp:    baseResponseOK,
			isKeyExistFn:  true,
			termsCacheErr: assert.AnError,
			reqTerm:       "202401",
			expectErr:     "Get term fail",
		},
		{
			name:       "terms db error",
			verifyResp: baseResponseOK,
			termsDBErr: assert.AnError,
			reqTerm:    "202401",
			expectErr:  "Get term from database fail",
		},
		{
			name:       "terms empty",
			verifyResp: baseResponseOK,
			termsDB:    (*dbmodel.UserTerm)(nil),
			reqTerm:    "202401",
			expectErr:  "Friend termList empty",
		},
		{
			name:       "invalid term - not in terms",
			verifyResp: baseResponseOK,
			termsDB:    &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401", "202402"})},
			reqTerm:    "202501",
			expectErr:  "Invalid term",
		},
		{
			name:       "invalid term - yjsy format but mapped term not in terms",
			verifyResp: baseResponseOK,
			termsDB:    &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202301", "202302"})}, // 只有202301和202302
			reqTerm:    "2024-2025-1",                                                                 // yjsy格式，映射到202401，不在terms中
			expectErr:  "Invalid term",
		},
		{
			name:       "invalid term - jwch format but mapped term not in terms",
			verifyResp: baseResponseOK,
			termsDB:    &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"2023-2024-1", "2023-2024-2"})}, // 只有yjsy格式的terms
			reqTerm:    "202401",                                                                                // jwch格式，映射到2024-2025-1，不在terms中
			expectErr:  "Invalid term",
		},
		{
			name:       "invalid term - default case (unrecognized format)",
			verifyResp: baseResponseOK,
			termsDB:    &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401"})},
			reqTerm:    "invalid-format", // 既不是yjsy也不是jwch格式，且不在terms中
			expectErr:  "Invalid term",
		},
		{
			name:            "yjsy term format mapping success",
			verifyResp:      baseResponseOK,
			isKeyExistFn:    true,
			termsCache:      []string{"202401", "202402"}, // 包含映射后的"202401"
			coursesCache:    []*jwch.Course{{Name: "Math", Teacher: "T1"}},
			reqTerm:         "2024-2025-1", // yjsy格式，MapYjsyTerm映射为"202401"
			expectLen:       1,
			expectFirstName: "Math",
		},
		{
			name:            "jwch term format mapping success",
			verifyResp:      baseResponseOK,
			isKeyExistFn:    true,
			termsCache:      []string{"2024-2025-1", "2024-2025-2"}, // 包含映射后的"2024-2025-1"
			coursesCache:    []*jwch.Course{{Name: "Physics", Teacher: "T2"}},
			reqTerm:         "202401", // jwch格式，MapJwchTerm映射为"2024-2025-1"
			expectLen:       1,
			expectFirstName: "Physics",
		},
		{
			name:            "term cache hit and course cache hit",
			verifyResp:      baseResponseOK,
			isKeyExistFn:    true,
			termsCache:      []string{"202401", "202402"},
			coursesCache:    []*jwch.Course{{Name: "Math", Teacher: "T1"}},
			reqTerm:         "202401",
			expectLen:       1,
			expectFirstName: "Math",
		},
		{
			name:            "course cache error",
			verifyResp:      baseResponseOK,
			isKeyExistFn:    true,
			termsCache:      []string{"202401"},
			coursesCacheErr: assert.AnError,
			reqTerm:         "202401",
			expectErr:       "Get courses fail",
		},
		{
			name:            "course cache empty -> yjsy cache fallback",
			verifyResp:      baseResponseOK,
			isKeyExistFn:    true,
			termsCache:      []string{"202401"},
			coursesCache:    []*jwch.Course(nil),
			coursesYjsy:     []*yjsy.Course{{Name: "Algo", Teacher: "T2"}},
			reqTerm:         "202401",
			expectLen:       1,
			expectFirstName: "Algo",
		},
		{
			name:           "yjsy cache error",
			verifyResp:     baseResponseOK,
			isKeyExistFn:   true,
			termsCache:     []string{"202401"},
			coursesCache:   []*jwch.Course(nil),
			coursesYjsyErr: assert.AnError,
			reqTerm:        "202401",
			expectErr:      "Get courses fail",
		},
		{
			name:       "term from db not in recent two",
			verifyResp: baseResponseOK,
			termsDB:    &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401", "202302", "202201"})},
			reqTerm:    "202201",
			expectErr:  "只能查看好友最近两个学期的课表",
		},
		{
			name:        "db course error",
			verifyResp:  baseResponseOK,
			termsDB:     &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401"})},
			dbCourseErr: assert.AnError,
			reqTerm:     "202401",
			expectErr:   "Get courses fail",
		},
		{
			name:       "db course missing",
			verifyResp: baseResponseOK,
			termsDB:    &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401"})},
			dbCourse:   (*dbmodel.UserCourse)(nil),
			reqTerm:    "202401",
			expectErr:  "there is no course in database",
		},
		{
			name:       "db course unmarshal error",
			verifyResp: baseResponseOK,
			termsDB:    &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401"})},
			dbCourse:   &dbmodel.UserCourse{TermCourses: "{"},
			reqTerm:    "202401",
			expectErr:  "Unmarshal fail",
		},
		{
			name:            "db course success",
			verifyResp:      baseResponseOK,
			termsDB:         &dbmodel.UserTerm{TermTime: pack.BuildTermOnDB([]string{"202401"})},
			dbCourse:        &dbmodel.UserCourse{TermCourses: `[{"name":"Database"}]`},
			reqTerm:         "202401",
			expectLen:       1,
			expectFirstName: "Database",
		},
	}

	defer mockey.UnPatchAll()
	for _, tc := range cases {
		mockey.PatchConvey(tc.name, t, func() {
			userCli := &mockUserClient{verifyResp: tc.verifyResp, verifyErr: tc.verifyErr}
			clientSet := &base.ClientSet{
				DBClient:    new(db.Database),
				CacheClient: new(cache.Cache),
				UserClient:  userCli,
			}

			mockey.Mock((*cache.Cache).IsKeyExist).Return(tc.isKeyExistFn).Build()
			mockey.Mock((*coursecache.CacheCourse).GetTermsCache).Return(tc.termsCache, tc.termsCacheErr).Build()
			mockey.Mock((*coursecache.CacheCourse).GetCoursesCache).Return(tc.coursesCache, tc.coursesCacheErr).Build()
			mockey.Mock((*coursecache.CacheCourse).GetCoursesCacheYjsy).Return(tc.coursesYjsy, tc.coursesYjsyErr).Build()
			mockey.Mock((*dbcourse.DBCourse).GetUserTermByStuId).Return(tc.termsDB, tc.termsDBErr).Build()
			mockey.Mock((*dbcourse.DBCourse).GetUserTermCourseByStuIdAndTerm).Return(tc.dbCourse, tc.dbCourseErr).Build()

			ctx := customContext.WithLoginData(context.Background(), loginData)
			svc := NewCourseService(ctx, clientSet, nil)
			res, err := svc.GetFriendCourse(&course.GetFriendCourseRequest{Id: "f1", Term: tc.reqTerm}, loginData)
			if tc.expectErr != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tc.expectErr)
				assert.Nil(t, res)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, res)
				assert.Len(t, res, tc.expectLen)
				if tc.expectLen > 0 {
					assert.Equal(t, tc.expectFirstName, res[0].Name)
				}
			}
		})
	}
}
