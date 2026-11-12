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

package classroom

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/environment"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
)

// DayAvailability 保存某校区一天的完整教室占用快照。
type DayAvailability struct {
	Date  string             `json:"date"`
	Term  string             `json:"term"`
	Rooms []*model.Classroom `json:"rooms"`
}

const availabilityFailureMessage = "教室数据拉取失败，请稍后重试"

func AvailabilityKey(date, campus string) string {
	return fmt.Sprintf("classroom:availability:%s:%s", date, campus)
}

func availabilityFailureKey(date, campus string) string {
	return fmt.Sprintf("classroom:availability:failure:%s:%s", date, campus)
}

func (c *CacheClassroom) GetAvailability(ctx context.Context, date, campus string) (*DayAvailability, error) {
	failure, err := c.client.Get(ctx, availabilityFailureKey(date, campus)).Result()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("get classroom availability failure status: %w", err)
	}
	if err == nil && failure != "" {
		return nil, errors.New(failure)
	}
	data, err := c.client.Get(ctx, AvailabilityKey(date, campus)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("get classroom availability: %w", err)
	}
	var day DayAvailability
	if err := json.Unmarshal(data, &day); err != nil {
		return nil, fmt.Errorf("decode classroom availability: %w", err)
	}
	// Rooms 为 nil 说明快照缺失有效内容，不能当作"当天无教室"处理。
	if day.Rooms == nil {
		return nil, fmt.Errorf("invalid classroom availability snapshot")
	}
	return &day, nil
}

// MarkAvailabilityFailure 标记重试后仍失败的查询，阻止客户端把旧快照当作本次刷新结果。
func (c *CacheClassroom) MarkAvailabilityFailure(ctx context.Context, date, campus string, failedBuildings ...string) error {
	if environment.IsTestEnvironment() {
		return nil
	}
	message := availabilityFailureText(failedBuildings)
	// 教务处将两个美院校区合并查询，失败状态也必须覆盖两个客户端校区。
	_, err := c.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, targetCampus := range availabilityCampuses(campus) {
			pipe.Set(ctx, availabilityFailureKey(date, targetCampus), message, constants.ClassroomKeyExpire)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("mark classroom availability failure: %w", err)
	}
	return nil
}

func availabilityFailureText(failedBuildings []string) string {
	message := availabilityFailureMessage
	if buildings := uniqueBuildingNames(failedBuildings); len(buildings) > 0 {
		message += "（失败教学楼：" + strings.Join(buildings, "、") + "）"
	}
	return message
}

func uniqueBuildingNames(buildings []string) []string {
	seen := make(map[string]struct{}, len(buildings))
	result := make([]string, 0, len(buildings))
	for _, building := range buildings {
		if building == "" {
			continue
		}
		if _, ok := seen[building]; ok {
			continue
		}
		seen[building] = struct{}{}
		result = append(result, building)
	}
	return result
}

func (c *CacheClassroom) SetAvailability(ctx context.Context, campus string, day *DayAvailability, failedBuildings ...string) error {
	if environment.IsTestEnvironment() {
		return nil
	}
	days := splitAvailability(campus, day)
	encoded := make(map[string][]byte, len(days))
	for targetCampus, snapshot := range days {
		data, err := json.Marshal(snapshot)
		if err != nil {
			return fmt.Errorf("encode classroom availability: %w", err)
		}
		encoded[targetCampus] = data
	}
	_, err := c.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		for targetCampus, data := range encoded {
			pipe.Set(ctx, AvailabilityKey(day.Date, targetCampus), data, constants.ClassroomKeyExpire)
			if len(failedBuildings) == 0 {
				pipe.Del(ctx, availabilityFailureKey(day.Date, targetCampus))
			} else {
				pipe.Set(ctx, availabilityFailureKey(day.Date, targetCampus), availabilityFailureText(failedBuildings), constants.ClassroomKeyExpire)
			}
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("set classroom availability: %w", err)
	}
	return nil
}

func availabilityCampuses(campus string) []string {
	if campus == "厦门工艺美院" {
		return []string{"鼓浪屿校区", "集美校区"}
	}
	return []string{campus}
}

func splitAvailability(campus string, day *DayAvailability) map[string]*DayAvailability {
	if campus != "厦门工艺美院" {
		return map[string]*DayAvailability{campus: day}
	}
	guLangYu := &DayAvailability{Date: day.Date, Term: day.Term, Rooms: make([]*model.Classroom, 0)}
	jiMei := &DayAvailability{Date: day.Date, Term: day.Term, Rooms: make([]*model.Classroom, 0)}
	for _, room := range day.Rooms {
		if strings.Contains(room.Location, "鼓浪屿") {
			copyRoom := *room
			copyRoom.Build = "鼓浪屿校区"
			guLangYu.Rooms = append(guLangYu.Rooms, &copyRoom)
		} else if strings.Contains(room.Location, "集美") {
			copyRoom := *room
			copyRoom.Build = "集美校区"
			jiMei.Rooms = append(jiMei.Rooms, &copyRoom)
		}
	}
	return map[string]*DayAvailability{"鼓浪屿校区": guLangYu, "集美校区": jiMei}
}
