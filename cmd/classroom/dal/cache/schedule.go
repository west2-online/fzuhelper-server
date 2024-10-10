package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
)

// ScheduledGetClassrooms 定期获取空教室信息
func ScheduledGetClassrooms() error {
	ctx := context.Background()
	// 定义jwch的stu客户端
	stu := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password)
	// 登录，id和cookies会自动保存在client中
	_ = stu.Login()

	var dates []string
	currentTime := time.Now()
	// 设定一周时间
	for i := 0; i < 7; i++ {
		date := currentTime.AddDate(0, 0, i).Format("2006-01-02")
		dates = append(dates, date)
	}
	logger.Infof("ScheduledGetClassrooms: start to get empty room info in the next 7 days: %v", dates)

	var eg errgroup.Group

	// 对每个日期启动一个 goroutine
	for _, date := range dates {
		date := date
		eg.Go(func() error {
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
						// 从jwch获取空教室信息
						switch campus {
						case "旗山校区":
							res, err = stu.GetQiShanEmptyRoom(args)
						default:
							res, err = stu.GetEmptyRoom(args)
						}
						if err != nil {
							return fmt.Errorf("ScheduledGetClassrooms: failed to get empty room info: %w", err)
						}
						// 收集结果并缓存
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
					}
				}

				logger.Infof("ScheduledGetClassrooms: complete all tasks of campus %s in the day %s", campus, date)
			}
			return nil
		})
	}

	// 等待所有 goroutine 完成，并收集错误
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("ScheduledGetClassrooms: failed to refresh empty room info: %w", err)
	}

	return nil
}
