package main

import (
	"context"
	"fmt"
	"time"

	"github.com/west2-online/fzuhelper-server/cmd/classroom/pack"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/service"
	classroom "github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// ClassroomServiceImpl implements the last service interface defined in the IDL.
type ClassroomServiceImpl struct{}

// GetEmptyRoom implements the ClassroomServiceImpl interface.
func (s *ClassroomServiceImpl) GetEmptyRoom(ctx context.Context, req *classroom.EmptyRoomRequest) (resp *classroom.EmptyRoomResponse, err error) {
	resp = classroom.NewEmptyRoomResponse()
	// 实际上前端会给定一个月内的选择，后端为了完整性，还是要判断一下
	// 判断req.date只能从今天开始的一个月内，在当前日期前或超过 30 天则报错
	// 首先判断date的格式是否符合要求
	requestDate, err := utils.TimeParse(req.Date)
	if err != nil {
		logger.Errorf("Classroom.GetEmptyRoom: date format error, err: %v", err)
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	now := time.Now().Truncate(24 * time.Hour)
	requestDate = requestDate.Truncate(24 * time.Hour)
	dateDiff := requestDate.Sub(now).Hours() / 24
	if dateDiff < 0 || dateDiff > 30 {
		err = fmt.Errorf("date out of range, date: %v", req.Date)
		logger.Infof("Classroom.GetEmptyRoom: %v", err)
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	l := service.NewClassroomService(ctx)
	res, err := l.GetEmptyRoom(req)
	if err != nil {
		logger.Infof("Classroom.GetEmptyRoom: GetEmptyRoom failed, err: %v", err)
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Rooms = pack.BuildClassRooms(res, req.Campus)
	logger.Info("Classroom.GetEmptyRoom: GetEmptyRoom success")
	return resp, nil
}
