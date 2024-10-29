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
