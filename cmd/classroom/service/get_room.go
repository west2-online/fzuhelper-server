package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/dal/cache"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"strings"
)

func (s *ClassroomService) GetEmptyRooms(req *classroom.EmptyRoomRequest) (rooms []string, err error) {
	key := fmt.Sprintf("%s.%s.%s.%s", req.Date, req.Campus, req.StartTime, req.EndTime)
	if cache.IsExistRoomInfo(s.ctx, key) {
		return cache.GetEmptyRoomCache(s.ctx, key), nil
	}
	switch req.Campus {
	case "旗山校区":
		rooms, err = s.getQiShanEmptyRooms(req)
		if err != nil {
			return nil, errors.WithMessage(err, "service.GetQiShanEmptyRooms failed")
		}
		go cache.SetEmptyRoomCache(s.ctx, key, rooms)
	case "鼓浪屿校区":
		rooms, err = s.getGuLangYuEmptyRooms(req)
		if err != nil {
			return nil, errors.WithMessage(err, "service.GetGuLangYuEmptyRooms failed")
		}
	case "集美校区":
		rooms, err = s.getJiMeiEmptyRooms(req)
		if err != nil {
			return nil, errors.WithMessage(err, "service.GetJiMeiEmptyRooms failed")
		}
	default:
		rooms, err = s.getOtherEmptyRooms(req)
		if err != nil {
			return nil, errors.WithMessage(err, "service.GetEmptyRooms failed")
		}
		go cache.SetEmptyRoomCache(s.ctx, key, rooms)
	}

	return rooms, nil
}

func (s *ClassroomService) getQiShanEmptyRooms(req *classroom.EmptyRoomRequest) (rooms []string, err error) {
	rooms, err = jwch.NewStudent().WithLoginData(s.Identifier, s.cookies).GetQiShanEmptyRoom(jwch.EmptyRoomReq{
		Campus: req.Campus,
		Time:   req.Date,
		Start:  req.StartTime,
		End:    req.EndTime,
	})
	if err != nil {
		utils.LoggerObj.Errorf("service.getQiShanEmptyRooms: %v", err)
		return nil, errors.Wrap(err, "service.getQiShanEmptyRooms failed")
	}
	return rooms, nil
}

func (s *ClassroomService) getOtherEmptyRooms(req *classroom.EmptyRoomRequest) (rooms []string, err error) {
	rooms, err = jwch.NewStudent().WithLoginData(s.Identifier, s.cookies).GetEmptyRoom(jwch.EmptyRoomReq{
		Campus: req.Campus,
		Time:   req.Date,
		Start:  req.StartTime,
		End:    req.EndTime,
	})
	if err != nil {
		utils.LoggerObj.Errorf("service.getOtherEmptyRooms: %v", err)
		return nil, errors.Wrap(err, "service.getOtherEmptyRooms failed")
	}

	return rooms, nil
}

func (s *ClassroomService) getGuLangYuEmptyRooms(req *classroom.EmptyRoomRequest) (res []string, err error) {
	key := fmt.Sprintf("%s.%s.%s.%s", req.Date, "鼓浪屿", req.StartTime, req.EndTime)
	if cache.IsExistRoomInfo(s.ctx, key) {
		return cache.GetEmptyRoomCache(s.ctx, key), nil
	}
	temp := req.Campus
	req.Campus = "厦门工艺美院"
	rooms, err := s.getOtherEmptyRooms(req)
	if err != nil {
		return nil, errors.WithMessage(err, "service.GetEmptyRooms.getGuLangYuEmptyRooms failed")
	}
	req.Campus = temp
	for _, room := range rooms {
		if strings.Contains(room, "鼓浪屿") {
			res = append(res, room)
		}
	}
	go cache.SetEmptyRoomCache(s.ctx, key, res)
	return res, nil
}

func (s *ClassroomService) getJiMeiEmptyRooms(req *classroom.EmptyRoomRequest) (res []string, err error) {
	key := fmt.Sprintf("%s,%s.%s.%s", req.Date, "集美", req.StartTime, req.EndTime)
	if cache.IsExistRoomInfo(s.ctx, key) {
		return cache.GetEmptyRoomCache(s.ctx, key), nil
	}
	temp := req.Campus
	req.Campus = "厦门工艺美院"
	rooms, err := s.getOtherEmptyRooms(req)
	if err != nil {
		return nil, errors.WithMessage(err, "service.GetEmptyRooms.getGuLangYuEmptyRooms failed")
	}
	req.Campus = temp
	for _, room := range rooms {
		if strings.Contains(room, "集美") {
			res = append(res, room)
		}
	}
	go cache.SetEmptyRoomCache(s.ctx, key, res)
	return res, nil
}
