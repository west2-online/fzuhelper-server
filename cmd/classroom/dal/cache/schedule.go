package cache

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
	"strconv"
	"time"
)

// ScheduledGetClassrooms 定期获取空教室信息
func ScheduledGetClassrooms() error {
	ctx := context.Background()
	//定义jwch的stu客户端
	stu := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password)
	//登录，id和cookies会自动保存在client中
	stu.Login()

	var dates []string
	currentTime := time.Now()
	//设定一周时间
	for i := 0; i < 7; i++ {
		date := currentTime.AddDate(0, 0, i).Format("2006-01-02")
		dates = append(dates, date)
	}
	//构建jwch的请求参数
	for _, date := range dates {
		for _, campus := range constants.CampusArray {
			for startTime := 1; startTime <= 11; startTime++ {
				for endTime := startTime; endTime <= 11; endTime++ {
					args := jwch.EmptyRoomReq{
						Campus: campus,
						Time:   date,
						Start:  strconv.Itoa(startTime),
						End:    strconv.Itoa(endTime),
					}
					var res []string
					var err error
					if campus == "旗山校区" {
						res, err = stu.GetQiShanEmptyRoom(args)
					} else {
						res, err = stu.GetEmptyRoom(args)
					}
					if err != nil {
						logger.LoggerObj.Errorf("classroom.cache.GetClassrooms: %v", err)
						return errors.Wrap(err, "classroom.cache.GetClassrooms failed")
					}
					if campus == "厦门工艺美院" {
						go SetXiaMenEmptyRoomCache(ctx, date, args.Start, args.End, res)
					} else {
						key := fmt.Sprintf("%s.%s.%s.%s", args.Time, args.Campus, args.Start, args.End)
						go SetEmptyRoomCache(ctx, key, res)
					}
					logger.LoggerObj.Debugf("classroom.cache.GetClassrooms add task %v", args)
				}
			}
			logger.LoggerObj.Infof("classroom.cache.CGetClassrooms add all tasks of campus %v in the day %v", campus, date)
		}
	}
	return nil
}
