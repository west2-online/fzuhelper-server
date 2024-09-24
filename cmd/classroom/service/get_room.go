package service

import (
	"github.com/pkg/errors"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
)

func (s *ClassroomService) GetEmptyRooms(req *classroom.EmptyRoomRequest) (rooms []string, err error) {
	//TODO: using redis cache
	if req.Campus == "旗山校区" {
		rooms, err = s.getQiShanEmptyRooms(req)
		if err != nil {
			return nil, errors.WithMessage(err, "service.GetQiShanEmptyRooms failed")
		}
	} else {
		rooms, err = s.getOtherEmptyRooms(req)
		if err != nil {
			return nil, errors.WithMessage(err, "service.GetEmptyRooms failed")
		}
	}
	//key := fmt.Sprintf("%s.%s.%s.%s", req.Time, req.Start, req.End, req.Building)
	//if exist, err := cache.IsExistRoomInfo(s.ctx, key); exist == 1 {
	//	// 获取缓存
	//	if err != nil {
	//		return nil, err
	//	}
	//	empty_room, err = cache.GetEmptyRoomCache(s.ctx, key)
	//	if err != nil {
	//		return nil, err
	//	}
	//} else {
	//	// 未命中缓存，登录进行爬取
	//	student := jwch.NewStudent().WithUser(*req.Account, *req.Password)
	//
	//	err = student.Login()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	err = student.CheckSession()
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	err, empty_room = student.GetEmptyRoom(jwch.EmptyRoomReq{
	//		Time:     req.Time,
	//		Start:    req.Start,
	//		End:      req.End,
	//		Building: req.Building,
	//	})
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//	// 异步写入缓存
	//	go cache.SetEmptyRoomCache(s.ctx, key, empty_room)
	//}
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
