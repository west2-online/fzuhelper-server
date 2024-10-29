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
	"strings"
)

// SetXiaMenEmptyRoomCache 设置厦门工艺美院的空教室缓存
// 因为前端给的数据只有鼓浪屿校区和集美校区的数据，所以这里需要单独对这两份数据进行处理
func SetXiaMenEmptyRoomCache(ctx context.Context, date, start, end string, emptyRoomList []string) (err error) {
	// 分别整理两个校区的结果
	guLangYuEmptyRooms := make([]string, 0)
	jiMeiEmptyRooms := make([]string, 0)
	for _, room := range emptyRoomList {
		if strings.Contains(room, "鼓浪屿") {
			guLangYuEmptyRooms = append(guLangYuEmptyRooms, room)
		} else if strings.Contains(room, "集美") {
			jiMeiEmptyRooms = append(jiMeiEmptyRooms, room)
		}
	}
	guLangYuKey := fmt.Sprintf("%s.%s.%s.%s", date, "鼓浪屿校区", start, end)
	jiMeiKey := fmt.Sprintf("%s.%s.%s.%s", date, "集美校区", start, end)
	err = SetEmptyRoomCache(ctx, guLangYuKey, guLangYuEmptyRooms)
	if err != nil {
		return fmt.Errorf("dal.SetXiaMenEmptyRoomCache: Set guLangYu rooms info failed: %w", err)
	}
	err = SetEmptyRoomCache(ctx, jiMeiKey, jiMeiEmptyRooms)
	if err != nil {
		return fmt.Errorf("dal.SetXiaMenEmptyRoomCache: Set jiMei rooms info failed: %w", err)
	}
	return nil
}
