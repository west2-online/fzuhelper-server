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

package api

import (
	"context"
	"testing"

	"github.com/bytedance/mockey"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/hertz/pkg/common/ut"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/api/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func TestClassroomVersionedRoutes(t *testing.T) {
	mockey.PatchConvey("V1 and V2 use distinct RPC methods", t, func() {
		mockey.Mock(rpc.GetEmptyRoomRPC).To(func(_ context.Context, req *classroom.EmptyRoomRequest) ([]*model.Classroom, error) {
			require.Equal(t, "3", req.StartTime)
			require.Equal(t, "5", req.EndTime)
			return []*model.Classroom{{Location: "free room"}}, nil
		}).Build()
		mockey.Mock(rpc.GetEmptyRoomV2RPC).To(func(_ context.Context, req *classroom.EmptyRoomV2Request) ([]*model.Classroom, error) {
			require.Equal(t, "2026-10-11", req.Date)
			require.Equal(t, "旗山校区", req.Campus)
			return []*model.Classroom{{Location: "occupied room", Occupancies: []*model.ClassroomOccupancy{{Start: 4, End: 4, Event: "补课"}}}}, nil
		}).Build()
		router := server.New()
		Register(router)
		v1 := ut.PerformRequest(router.Engine, consts.MethodGet, "/api/v1/common/classroom/empty?date=2026-10-11&campus=旗山校区&startTime=3&endTime=5", nil)
		require.Equal(t, consts.StatusOK, v1.Code)
		require.Contains(t, string(v1.Result().Body()), `"location":"free room"`)
		require.NotContains(t, string(v1.Result().Body()), "occupancies")
		v2 := ut.PerformRequest(router.Engine, consts.MethodGet, "/api/v2/common/classroom/empty?date=2026-10-11&campus=旗山校区", nil)
		require.Equal(t, consts.StatusOK, v2.Code)
		require.Contains(t, string(v2.Result().Body()), `"location":"occupied room"`)
		require.Contains(t, string(v2.Result().Body()), `"event":"补课"`)
	})
}
