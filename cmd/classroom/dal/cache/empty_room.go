package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"strings"
)

// SetXiaMenEmptyRoomCache 设置厦门工艺美院的空教室缓存
// 因为前端给的数据只有鼓浪屿校区和集美校区的数据，所以这里需要单独对这两份数据进行处理
func SetXiaMenEmptyRoomCache(ctx context.Context, date, start, end string, emptyRoomList []string) {
	//分别整理两个校区的结果
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
	SetEmptyRoomCache(ctx, guLangYuKey, guLangYuEmptyRooms)
	SetEmptyRoomCache(ctx, jiMeiKey, jiMeiEmptyRooms)
}

func SetEmptyRoomCache(ctx context.Context, key string, emptyRoomList []string) {
	emptyRoomJson, err := json.Marshal(emptyRoomList)
	// 两天过期
	err = RedisClient.Set(ctx, key, emptyRoomJson, constants.ClassroomKeyExpire).Err()
	if err != nil {
		logger.LoggerObj.Fatalf("dal.cache.SetEmptyRoomCache failed, err is %v", err)
	}
}
func GetEmptyRoomCache(ctx context.Context, key string) (emptyRoomList []string) {
	data, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		logger.LoggerObj.Fatalf("dal.cache.GetEmptyRoomCache failed, err is %v", err)
		return nil
	}
	err = json.Unmarshal([]byte(data), &emptyRoomList)
	if err != nil {
		logger.LoggerObj.Fatalf("dal.cache.GetEmptyRoomCache Unmarshal failed, err is %v", err)
		return nil
	}
	return
}
func IsExistRoomInfo(ctx context.Context, key string) bool {
	return RedisClient.Exists(ctx, key).Val() == 1
}
