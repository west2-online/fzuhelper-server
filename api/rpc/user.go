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

package rpc

import (
	"context"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/tracing"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitUserRPC() {
	c, err := client.InitUserRPC()
	if err != nil {
		logger.Fatalf("api.rpc.user InitUserRPC failed, err is %v", err)
	}
	userClient = *c
}

func GetLoginDataRPC(ctx context.Context, req *user.GetLoginDataRequest) (string, string, error) {
	resp, err := userClient.GetLoginData(ctx, req)
	if err != nil {
		tracing.RecordErrorAndSetStatus(ctx, err)
		logger.WithCtx(ctx).Errorf("GetLoginDataRPC: RPC called failed: %v", err.Error())
		return "", "", errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return "", "", errno.BizError.WithMessage("教务处登录失败: " + resp.Base.Msg)
	}
	return resp.Id, resp.Cookies, nil
}

func GetUserInfoRPC(ctx context.Context, req *user.GetUserInfoRequest) (*model.UserInfo, error) {
	resp, err := userClient.GetUserInfo(ctx, req)
	if err != nil {
		tracing.RecordErrorAndSetStatus(ctx, err)
		logger.WithCtx(ctx).Errorf("GetUserInfoRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if err = utils.HandleBaseRespWithCookie(resp.Base); err != nil {
		return nil, err
	}
	return resp.Data, nil
}

func GetLoginDataForYJSYRPC(ctx context.Context, req *user.GetLoginDataForYJSYRequest) (string, string, error) {
	resp, err := userClient.GetGetLoginDataForYJSY(ctx, req)
	if err != nil {
		tracing.RecordErrorAndSetStatus(ctx, err)
		logger.WithCtx(ctx).Errorf("GetLoginDataRPC: RPC called failed: %v", err.Error())
		return "", "", errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return "", "", errno.BizError.WithMessage("研究生管理系统登录失败: " + resp.Base.Msg)
	}
	return resp.Id, resp.Cookies, nil
}

func GetInvitationCodeRPC(ctx context.Context, req *user.GetInvitationCodeRequest) (string, int64, error) {
	resp, err := userClient.GetInvitationCode(ctx, req)
	if err != nil {
		tracing.RecordErrorAndSetStatus(ctx, err)
		logger.WithCtx(ctx).Errorf("GetInvitationCodeRPC: RPC called failed: %v", err.Error())
		return "", -1, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return "", -1, errno.BizError.WithMessage("尝试生成邀请码失败: " + resp.Base.Msg)
	}
	return resp.InvitationCode, resp.ExpireAt, nil
}

func BindInvitationRPC(ctx context.Context, req *user.BindInvitationRequest) error {
	resp, err := userClient.BindInvitation(ctx, req)
	if err != nil {
		tracing.RecordErrorAndSetStatus(ctx, err)
		logger.WithCtx(ctx).Errorf("BindInvitationRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.BizError.WithMessage("验证邀请码失败: " + resp.Base.Msg)
	}
	return nil
}

func GetFriendListRPC(ctx context.Context, req *user.GetFriendListRequest) ([]*model.UserFriendInfo, error) {
	resp, err := userClient.GetFriendList(ctx, req)
	if err != nil {
		tracing.RecordErrorAndSetStatus(ctx, err)
		logger.WithCtx(ctx).Errorf("GetFriendListRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.BizError.WithMessage("查看好友列表失败: " + resp.Base.Msg)
	}
	return resp.Data, nil
}

func DeleteFriendRPC(ctx context.Context, req *user.DeleteFriendRequest) error {
	resp, err := userClient.DeleteFriend(ctx, req)
	if err != nil {
		tracing.RecordErrorAndSetStatus(ctx, err)
		logger.WithCtx(ctx).Errorf("DeleteFriendRpc: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.BizError.WithMessage("删除好友失败: " + resp.Base.Msg)
	}
	return nil
}

func CancelInviteRPC(ctx context.Context, req *user.CancelInviteRequest) error {
	resp, err := userClient.CancelInvite(ctx, req)
	if err != nil {
		tracing.RecordErrorAndSetStatus(ctx, err)
		logger.WithCtx(ctx).Errorf("CancelInviteRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.BizError.WithMessage("作废邀请码失败: " + resp.Base.Msg)
	}
	return nil
}

func GetFriendMaxNumRPC(ctx context.Context, req *user.GetFriendMaxNumRequest) (*model.FriendMaxNumInfo, error) {
	resp, err := userClient.GetFriendMaxNum(ctx, req)
	if err != nil {
		tracing.RecordErrorAndSetStatus(ctx, err)
		logger.WithCtx(ctx).Errorf("GetFriendMaxNumRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.BizError.WithMessage("获取好友数量上限失败: " + resp.Base.Msg)
	}
	return resp.Data, nil
}

func ReorderFriendListRPC(ctx context.Context, req *user.ReorderFriendListRequest) error {
	resp, err := userClient.ReorderFriendList(ctx, req)
	if err != nil {
		tracing.RecordErrorAndSetStatus(ctx, err)
		logger.WithCtx(ctx).Errorf("ReorderFriendListRPC: RPC called failed: %v", err.Error())
		return errno.InternalServiceError.WithError(err)
	}
	if !utils.IsSuccess(resp.Base) {
		return errno.BizError.WithMessage("好友列表排序失败: " + resp.Base.Msg)
	}
	return nil
}
