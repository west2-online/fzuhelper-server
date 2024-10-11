/*
Copyright 2024 The west2-online Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/west2-online/fzuhelper-server/pkg/utils"

	"golang.org/x/sync/errgroup"

	"github.com/west2-online/fzuhelper-server/config"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
)

// ScheduledGetClassrooms 定期获取空教室信息
func ScheduledGetClassrooms() error {
	ctx := context.Background()
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
		// 将 date 作为参数传递给 goroutine
		eg.Go(func(date string) func() error {
			return func() error {
				// 创建新的 stu 对象，每个 goroutine 都有自己的 client
				// 定义 jwch 的 stu 客户端
				stu := jwch.NewStudent().WithUser(config.DefaultUser.Account, config.DefaultUser.Password)
				// 登录，id 和 cookies 会自动保存在 client 中
				// 如果登录失败，重试
				err := utils.RetryLogin(stu)
				if err != nil {
					return fmt.Errorf("ScheduledGetClassrooms: failed to login: %w", err)
				}
				for _, campus := range constants.CampusArray {
					for startTime := 1; startTime <= 11; startTime++ {
						for endTime := startTime; endTime <= 11; endTime++ {
							args := jwch.EmptyRoomReq{
								Campus: campus,
								Time:   date, // 使用传递进来的 date 参数
								Start:  strconv.Itoa(startTime),
								End:    strconv.Itoa(endTime),
							}
							var res []string
							var err error
							// 从 jwch 获取空教室信息
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
			}
		}(date))
	}

	// 等待所有 goroutine 完成，并收集错误
	if err := eg.Wait(); err != nil {
		return fmt.Errorf("ScheduledGetClassrooms: failed to refresh empty room info: %w", err)
	}

	return nil
}
