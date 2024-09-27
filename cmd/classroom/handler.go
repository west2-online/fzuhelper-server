package main

import (
	"context"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/pack"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/service"
	classroom "github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"strconv"
	"time"
)

// ClassroomServiceImpl implements the last service interface defined in the IDL.
type ClassroomServiceImpl struct{}

// GetEmptyRoom implements the ClassroomServiceImpl interface.
func (s *ClassroomServiceImpl) GetEmptyRoom(ctx context.Context, req *classroom.EmptyRoomRequest) (resp *classroom.EmptyRoomResponse, err error) {
	// TODO: Your code here...
	resp = classroom.NewEmptyRoomResponse()
	l := service.NewClassroomServiceInDefault(ctx)
	res, err := l.GetEmptyRooms(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Rooms = pack.BuildClassRooms(res, req.Campus)
	utils.LoggerObj.Info("GetEmptyRoom success")
	return
}

// CacheEmptyRooms 缓存所有空教室的数据
// 缓存一周的所有信息，每两天更新一次
func CacheEmptyRooms() {
	ctx := context.Background()
	id, cookies := jwch.NewStudent().WithUser(constants.DefaultAccount, constants.DefaultPassword).GetIdentifierAndCookies()
	l := service.NewClassroomService(ctx, id, cookies)

	var dates []string
	currentTime := time.Now()
	//设定一周时间
	for i := 0; i < 7; i++ {
		date := currentTime.AddDate(0, 0, i).Format("2006-01-02")
		dates = append(dates, date)
	}
	for _, date := range dates {
		for _, campus := range constants.CampusArray {
			for startTime := 1; startTime <= 11; startTime++ {
				for endTime := startTime; endTime <= 11; endTime++ {
					args := &classroom.EmptyRoomRequest{
						Date:      date,
						Campus:    campus,
						StartTime: strconv.Itoa(startTime),
						EndTime:   strconv.Itoa(endTime),
					}
					go l.GetEmptyRooms(args)
					//给3s的时间让服务器处理
					time.Sleep(3 * time.Second)
				}
			}
			utils.LoggerObj.Infof("Complete Data of campus %s", campus)
		}
		utils.LoggerObj.Infof("Complete Data of date %s", date)
	}
	time.Sleep(24 * time.Hour * 2)
}
