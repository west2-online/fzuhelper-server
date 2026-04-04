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

package user

import (
	"context"

	"github.com/west2-online/fzuhelper-server/internal/user/pack"
	"github.com/west2-online/fzuhelper-server/internal/user/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/base"
	metainfoContext "github.com/west2-online/fzuhelper-server/pkg/base/context"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/taskqueue"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	ClientSet *base.ClientSet
	taskQueue taskqueue.TaskQueue
}

func NewUserService(clientSet *base.ClientSet, taskQueue taskqueue.TaskQueue) *UserServiceImpl {
	return &UserServiceImpl{
		ClientSet: clientSet,
		taskQueue: taskQueue,
	}
}

// GetLoginData implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetLoginData(ctx context.Context, req *user.GetLoginDataRequest) (resp *user.GetLoginDataResponse, err error) {
	resp = new(user.GetLoginDataResponse)
	id, cookies, err := service.NewUserService(ctx, s.ClientSet, nil).GetLoginData(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Id = id
	resp.Cookies = cookies
	return resp, nil
}

// GetUserInfo implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetUserInfo(ctx context.Context, request *user.GetUserInfoRequest) (resp *user.GetUserInfoResponse, err error) {
	resp = new(user.GetUserInfoResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		resp.Base = base.BuildBaseResp(errno.Errorf(errno.AuthErrorCode, "User.GetUserInfo: Get login data fail %v", err))
		return resp, nil
	}
	if utils.IsGraduate(loginData.Id) {
		info, err := service.NewUserService(ctx, s.ClientSet, s.taskQueue).GetUserInfoYjsy(loginData)
		resp.Base = base.BuildBaseResp(err)
		if err != nil {
			return resp, nil
		}
		resp.Data = pack.BuildInfoResp(info)
		return resp, nil
	} else {
		info, err := service.NewUserService(ctx, s.ClientSet, s.taskQueue).GetUserInfo(loginData)
		resp.Base = base.BuildBaseResp(err)
		if err != nil {
			return resp, nil
		}
		resp.Data = pack.BuildInfoResp(info)
		return resp, nil
	}
}

// GetGetLoginDataForYJSY implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetGetLoginDataForYJSY(ctx context.Context, req *user.GetLoginDataForYJSYRequest) (
	resp *user.GetLoginDataForYJSYResponse, err error,
) {
	resp = new(user.GetLoginDataForYJSYResponse)
	cookies, err := service.NewUserService(ctx, s.ClientSet, nil).GetLoginDataForYJSY(req)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Id = utils.MarkGraduate(req.Id) // yjsy的访问不需要id，5个前导0+学号表示研究生标识
	resp.Cookies = cookies
	return resp, nil
}

// GetInvitationCode implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetInvitationCode(ctx context.Context, request *user.GetInvitationCodeRequest) (
	resp *user.GetInvitationCodeResponse, err error,
) {
	resp = new(user.GetInvitationCodeResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		resp.Base = base.BuildBaseResp(errno.Errorf(errno.AuthErrorCode, "User.GetInvitationCode: Get login data fail %v", err))
		return resp, nil
	}
	code, expireAt, err := service.NewUserService(ctx, s.ClientSet, s.taskQueue).GetInvitationCode(loginData, request.GetIsRefresh())
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.InvitationCode = code
	resp.ExpireAt = expireAt
	return resp, nil
}

// BindInvitation implements the UserServiceImpl interface.
func (s *UserServiceImpl) BindInvitation(ctx context.Context, request *user.BindInvitationRequest) (
	resp *user.BindInvitationResponse, err error,
) {
	resp = new(user.BindInvitationResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		resp.Base = base.BuildBaseResp(errno.Errorf(errno.AuthErrorCode, "User.BindInvitation: Get login data fail %v", err))
		return resp, nil
	}
	err = service.NewUserService(ctx, s.ClientSet, s.taskQueue).BindInvitation(loginData, request.InvitationCode)
	resp.Base = base.BuildBaseResp(err)
	return resp, nil
}

// GetFriendList implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFriendList(ctx context.Context, request *user.GetFriendListRequest) (
	resp *user.GetFriendListResponse, err error,
) {
	resp = new(user.GetFriendListResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		resp.Base = base.BuildBaseResp(errno.Errorf(errno.AuthErrorCode, "User.GetFriendList: Get login data fail %v", err))
		return resp, nil
	}
	data, err := service.NewUserService(ctx, s.ClientSet, s.taskQueue).GetFriendList(loginData)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Data = data
	return resp, nil
}

// DeleteFriend implements the UserServiceImpl interface.
func (s *UserServiceImpl) DeleteFriend(ctx context.Context, request *user.DeleteFriendRequest) (
	resp *user.DeleteFriendResponse, err error,
) {
	resp = new(user.DeleteFriendResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		resp.Base = base.BuildBaseResp(errno.Errorf(errno.AuthErrorCode, "User.DeleteFriend: Get login data fail %v", err))
		return resp, nil
	}
	err = service.NewUserService(ctx, s.ClientSet, s.taskQueue).DeleteUserFriend(loginData, request.Id)
	resp.Base = base.BuildBaseResp(err)
	return resp, nil
}

func (s *UserServiceImpl) VerifyFriend(ctx context.Context, request *user.VerifyFriendRequest) (
	resp *user.VerifyFriendResponse, err error,
) {
	resp = new(user.VerifyFriendResponse)
	res, err := service.NewUserService(ctx, s.ClientSet, s.taskQueue).VerifyUserFriend(request.Id, request.FriendId)
	resp.Base = base.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.FriendExist = res
	return resp, nil
}

func (s *UserServiceImpl) CancelInvite(ctx context.Context, request *user.CancelInviteRequest) (
	resp *user.CancelInviteResponse, err error,
) {
	resp = new(user.CancelInviteResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		resp.Base = base.BuildBaseResp(errno.Errorf(errno.AuthErrorCode, "User.CancelInvite: Get login data fail %v", err))
		return resp, nil
	}
	err = service.NewUserService(ctx, s.ClientSet, s.taskQueue).CancelInvitationCode(loginData)
	resp.Base = base.BuildBaseResp(err)
	return resp, nil
}

// GetFriendMaxNum implements the UserServiceImpl interface.
func (s *UserServiceImpl) GetFriendMaxNum(ctx context.Context, request *user.GetFriendMaxNumRequest) (
	resp *user.GetFriendMaxNumResponse, err error,
) {
	resp = new(user.GetFriendMaxNumResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		resp.Base = base.BuildBaseResp(errno.Errorf(errno.AuthErrorCode, "User.GetFriendMaxNum: Get login data fail %v", err))
		return resp, nil
	}
	maxNum := service.NewUserService(ctx, s.ClientSet, s.taskQueue).GetFriendMaxNum(loginData)
	resp.Base = base.BuildSuccessResp()
	resp.Data = &model.FriendMaxNumInfo{MaxNum: maxNum}
	return resp, nil
}

// ReorderFriendList implements the UserServiceImpl interface.
func (s *UserServiceImpl) ReorderFriendList(ctx context.Context, request *user.ReorderFriendListRequest) (
	resp *user.ReorderFriendListResponse, err error,
) {
	resp = new(user.ReorderFriendListResponse)
	loginData, err := metainfoContext.GetLoginData(ctx)
	if err != nil {
		resp.Base = base.BuildBaseResp(errno.Errorf(errno.AuthErrorCode, "User.ReorderFriendList: Get login data fail %v", err))
		return resp, nil
	}
	err = service.NewUserService(ctx, s.ClientSet, s.taskQueue).ReorderFriendList(loginData, request.FriendIds)
	resp.Base = base.BuildBaseResp(err)
	return resp, nil
}
