package main

import (
	"context"
	"fmt"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/pack"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/service"
	classroom "github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"time"
)

// ClassroomServiceImpl implements the last service interface defined in the IDL.
type ClassroomServiceImpl struct{}

// GetEmptyRoom implements the ClassroomServiceImpl interface.
func (s *ClassroomServiceImpl) GetEmptyRoom(ctx context.Context, req *classroom.EmptyRoomRequest) (resp *classroom.EmptyRoomResponse, err error) {
	resp = classroom.NewEmptyRoomResponse()
	//实际上前端会给定一个月内的选择，后端为了完整性，还是要判断一下
	//判断req.date只能从今天开始的一个月内
	//首先判断date的格式是否符合要求
	requestDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		logger.LoggerObj.Errorf("Classroom.GetEmptyRoom: date format error, err: %v", err)
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	// 获取当前日期，不包含时间部分
	now := time.Now().Truncate(24 * time.Hour)
	requestDate = requestDate.Truncate(24 * time.Hour)
	// 计算日期差异
	dateDiff := requestDate.Sub(now).Hours() / 24
	if dateDiff < 0 || dateDiff > 30 {
		err = fmt.Errorf("Classroom.GetEmptyRoom: date out of range, date: %v", req.Date)
		logger.LoggerObj.Errorf("Classroom.GetEmptyRoom: %v", err)
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	l := service.NewClassroomService(ctx)
	res, err := l.GetEmptyRoom(req)
	if err != nil {
		logger.LoggerObj.Errorf("Classroom.GetEmptyRoom: GetEmptyRoom failed, err: %v", err)
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Rooms = pack.BuildClassRooms(res, req.Campus)
	logger.LoggerObj.Info("Classroom.GetEmptyRoom: GetEmptyRoom success")
	return
}
