package cache

import (
	"context"
	"fmt"
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
	_ = stu.Login()

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
					//从jwch获取空教室信息
					//分为从旗山校区爬取和其他校区爬取
					switch campus {
					case "旗山校区":
						res, err = stu.GetQiShanEmptyRoom(args)
					default:
						res, err = stu.GetEmptyRoom(args)
					}
					if err != nil {
						return fmt.Errorf("ScheduledGetClassrooms: failed to get empty room info: %w", err)
					}
					//收集结果，如果是厦门工艺美院，分为集美校区和鼓浪屿校区，需要单独分开处理
					switch campus {
					case "厦门工艺美院":
						err = SetXiaMenEmptyRoomCache(ctx, date, args.Start, args.End, res)
						if err != nil {
							return fmt.Errorf("ScheduledGetClassrooms: failed to set xiamen empty room cache: %w", err)
						}
					default:
						key := fmt.Sprintf("%s.%s.%s.%s", args.Time, args.Campus, args.Start, args.End)
						err = SetEmptyRoomCache(ctx, key, res)
						if err != nil {
							return fmt.Errorf("ScheduledGetClassrooms: failed to set empty room cache: %w", err)
						}
					}
					logger.LoggerObj.Debugf("ScheduledGetClassrooms: complete the task %v", args)
				}
			}
			logger.LoggerObj.Infof("classroom.cache.GetClassrooms complete all tasks of campus %v in the day %v", campus, date)
		}
	}
	return nil
}
