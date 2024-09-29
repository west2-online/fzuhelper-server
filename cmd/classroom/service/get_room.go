package service

import (
	"errors"
	"fmt"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/dal/cache"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
)

func (s *ClassroomService) GetEmptyRoom(req *classroom.EmptyRoomRequest) ([]string, error) {
	//从redis中获取数据
	key := fmt.Sprintf("%s.%s.%s.%s", req.Date, req.Campus, req.StartTime, req.EndTime)
	if ok := cache.IsExistRoomInfo(s.ctx, key); !ok {
		return nil, errors.New("service.GetEmptyRoom IsExistRoomInfo failed")
	}
	emptyRoomList := cache.GetEmptyRoomCache(s.ctx, key)
	return emptyRoomList, nil
}
