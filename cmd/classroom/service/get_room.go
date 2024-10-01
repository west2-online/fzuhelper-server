package service

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/west2-online/fzuhelper-server/cmd/classroom/dal/cache"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (s *ClassroomService) GetEmptyRoom(req *classroom.EmptyRoomRequest) ([]string, error) {
	//从redis中获取数据
	key := fmt.Sprintf("%s.%s.%s.%s", req.Date, req.Campus, req.StartTime, req.EndTime)
	if ok := cache.IsExistRoomInfo(s.ctx, key); !ok {
		return nil, errors.Wrap(errno.InternalServiceError, "service.GetEmptyRoom: room info not exist")
	}
	emptyRoomList, err := cache.GetEmptyRoomCache(s.ctx, key)
	if err != nil {
		return nil, errors.WithMessage(err, "service.GetEmptyRoom: Get room info failed")
	}
	return emptyRoomList, nil
}
