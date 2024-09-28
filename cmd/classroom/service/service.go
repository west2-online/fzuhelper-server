package service

import (
	"context"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"github.com/west2-online/jwch"
	"net/http"
	"strconv"
	"time"
)

type ClassroomService struct {
	ctx        context.Context
	Identifier string
	cookies    []*http.Cookie
}

func NewClassroomServiceInDefault(ctx context.Context) *ClassroomService {

	id, cookies := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password).GetIdentifierAndCookies()
	return &ClassroomService{
		ctx:        ctx,
		Identifier: id,
		cookies:    cookies,
	}
}

func NewClassroomService(ctx context.Context, identifier string, cookies []*http.Cookie) *ClassroomService {
	return &ClassroomService{
		ctx:        ctx,
		Identifier: identifier,
		cookies:    cookies,
	}
}

// CacheEmptyRooms 缓存所有空教室的数据
// 缓存一周的所有信息，每两天更新一次
func CacheEmptyRooms() {
	for {
		ctx := context.Background()
		id, cookies := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password).GetIdentifierAndCookies()
		l := NewClassroomService(ctx, id, cookies)
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
		time.Sleep(constants.ClassroomKeyExpire)
	}
}
