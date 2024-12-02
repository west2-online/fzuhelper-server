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

	"github.com/west2-online/fzuhelper-server/kitex_gen/course"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/client"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func InitCourseRPC() {
	c, err := client.InitCourseRPC()
	if err != nil {
		logger.Fatalf("api.rpc.course InitCourseRPC failed, err  %v", err)
	}
	courseClient = *c
}

func GetCourseListRPC(ctx context.Context, req *course.CourseListRequest) (courses []*model.Course, err error) {
	resp, err := courseClient.GetCourseList(ctx, req)
	if err != nil {
		logger.Errorf("GetCourseListRPC: RPC called failed: %v", err.Error())
		return nil, errno.InternalServiceError.WithMessage(err.Error())
	}

	if !utils.IsSuccess(resp.Base) {
		return nil, errno.BizError.WithMessage(resp.Base.Msg)
	}

	return resp.Data, nil
}
