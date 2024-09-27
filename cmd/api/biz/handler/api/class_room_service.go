// Code generated by hertz generator.

package api

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/model/api"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/pack"
	"github.com/west2-online/fzuhelper-server/cmd/api/biz/rpc"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// GetEmptyClassrooms .
// @router /api/v1/common/classroom/empty [GET]
// 获取空教室统一不需要id和cookies
func GetEmptyClassrooms(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.EmptyClassroomRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.LoggerObj.Error("api.GetEmptyClassrooms: BindAndValidate", err)
		pack.RespError(c, errno.ParamEmpty)
		return
	}
	res, err := rpc.GetEmptyRoomRPC(ctx, &classroom.EmptyRoomRequest{
		Date:      req.Date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Campus:    req.Campus,
	})
	if err != nil {
		utils.LoggerObj.Error("api.GetEmptyClassrooms: GetEmptyRoomRPC", err)
		pack.RespError(c, errno.InternalServiceError)
		return
	}
	resp := new(api.EmptyClassroomResponse)
	resp.Classrooms = pack.BuildClassroomList(res)
	pack.RespList(c, resp.Classrooms)
}
