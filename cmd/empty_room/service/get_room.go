package service

import (
	"fmt"

	"github.com/west2-online/fzuhelper-server/cmd/empty_room/dal/cache"
	"github.com/west2-online/fzuhelper-server/kitex_gen/empty_room"
	"github.com/west2-online/jwch"
)

func (s *EmptyRoomService) GetRoom(req *empty_room.EmptyRoomRequest) (empty_room []string, err error) {
	key := fmt.Sprintf("%s.%s.%s.%s", req.Time, req.Start, req.End, req.Building)
	if exist, err := cache.IsExistRoomInfo(s.ctx, key); exist == 1 {
		// 获取缓存
		if err != nil {
			return nil, err
		}
		empty_room, err = cache.GetEmptyRoomCache(s.ctx, key)
		if err != nil {
			return nil, err
		}
	} else {
		// 未命中缓存，登录进行爬取
		student := jwch.NewStudent().WithUser(*req.Account, *req.Password)

		err = student.Login()
		if err != nil {
			return nil, err
		}

		err = student.CheckSession()
		if err != nil {
			return nil, err
		}

		err, empty_room = student.GetEmptyRoom(jwch.EmptyRoomReq{
			Time:     req.Time,
			Start:    req.Start,
			End:      req.End,
			Building: req.Building,
		})
		if err != nil {
			return nil, err
		}

		// 异步写入缓存
		go cache.SetEmptyRoomCache(s.ctx, key, empty_room)
	}
	return empty_room, nil
}
