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

	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitClassroomRPC() {
	c, err := client.InitClassroomRPC()
	if err != nil {
		logger.Fatalf("api.rpc.classroom InitClassroomRPC failed, err  %v", err)
	}
	classroomClient = *c
}

func GetEmptyRoomRPC(ctx context.Context, req *classroom.EmptyRoomRequest) (emptyRooms []*model.Classroom, err error) {
	resp, err := classroomClient.GetEmptyRoom(ctx, req)
	if err != nil {
		logger.Errorf("GetEmptyRoomRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.BizError.WithMessage(resp.Base.Msg)
	}
	return resp.Rooms, nil
}

func GetExamRoomInfoRPC(ctx context.Context, req *classroom.ExamRoomInfoRequest) (roomInfo []*model.ExamRoomInfo, err error) {
	resp, err := classroomClient.GetExamRoomInfo(ctx, req)
	if err != nil {
		return nil, errno.Errorf(errno.InternalGRPCErrorCode, "GetExamRoomInfoRPC: RPC called failed: %v", err.Error())
	}
	if !utils.IsSuccess(resp.Base) {
		return nil, errno.NewErrNo(resp.Base.Code, resp.Base.Msg) // 由于 rpc 的错误是通过 base 返回的，所以 api 部分只需要再次简单包装一下继续上抛，保持住原始错误
	}
	return resp.Rooms, nil
}
