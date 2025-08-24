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
	"fmt"
	"strings"

	"github.com/west2-online/fzuhelper-server/pkg/base/environment"
)

func (c *CacheClassroom) SetXiaMenEmptyRoomCache(ctx context.Context, date, start, end string, emptyRoomList []string) (err error) {
	if environment.IsTestEnvironment() {
		return nil
	}
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
	err = c.SetEmptyRoomCache(ctx, guLangYuKey, guLangYuEmptyRooms)
	if err != nil {
		return fmt.Errorf("dal.SetXiaMenEmptyRoomCache: Set guLangYu rooms info failed: %w", err)
	}
	err = c.SetEmptyRoomCache(ctx, jiMeiKey, jiMeiEmptyRooms)
	if err != nil {
		return fmt.Errorf("dal.SetXiaMenEmptyRoomCache: Set jiMei rooms info failed: %w", err)
	}
	return nil
}
