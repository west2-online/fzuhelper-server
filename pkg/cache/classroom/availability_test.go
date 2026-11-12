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
	"testing"

	"github.com/bytedance/mockey"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/base/environment"
)

// availabilityRedisHook 模拟 Redis 命令执行，验证失败标记与快照读写，无需连接真实 Redis。
type availabilityRedisHook struct {
	values map[string]string
	err    error
}

func (h *availabilityRedisHook) DialHook(next redis.DialHook) redis.DialHook {
	return next
}

func (h *availabilityRedisHook) ProcessHook(_ redis.ProcessHook) redis.ProcessHook {
	return func(_ context.Context, cmd redis.Cmder) error {
		if h.err != nil {
			return h.err
		}
		value, exists := h.values[fmt.Sprint(cmd.Args()[1])]
		if !exists {
			cmd.SetErr(redis.Nil)
			return redis.Nil
		}
		result, ok := cmd.(*redis.StringCmd)
		if !ok {
			return fmt.Errorf("unexpected Redis command %s", cmd.Name())
		}
		result.SetVal(value)
		return nil
	}
}

func (h *availabilityRedisHook) ProcessPipelineHook(_ redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(_ context.Context, commands []redis.Cmder) error {
		if h.err != nil {
			return h.err
		}
		for _, cmd := range commands {
			switch cmd.Name() {
			case "set":
				value := fmt.Sprint(cmd.Args()[2])
				if data, ok := cmd.Args()[2].([]byte); ok {
					value = string(data)
				}
				h.values[fmt.Sprint(cmd.Args()[1])] = value
			case "del":
				delete(h.values, fmt.Sprint(cmd.Args()[1]))
			}
		}
		return nil
	}
}

func TestAvailabilityFailureAndRecovery(t *testing.T) {
	mockey.PatchConvey("失败标记优先于旧快照，成功刷新清除两个美院校区的标记", t, func() {
		mockey.Mock(environment.IsTestEnvironment).Return(false).Build()
		hook := &availabilityRedisHook{values: make(map[string]string)}
		client := redis.NewClient(&redis.Options{})
		client.AddHook(hook)
		t.Cleanup(func() { require.NoError(t, client.Close()) })
		cache := NewCacheClassroom(client)
		ctx := context.Background()
		day := &DayAvailability{Date: "2026-10-12", Term: "202601", Rooms: []*model.Classroom{}}
		require.NoError(t, cache.SetAvailability(ctx, "厦门工艺美院", day))
		require.NoError(t, cache.MarkAvailabilityFailure(ctx, day.Date, "厦门工艺美院"))
		for _, campus := range []string{"鼓浪屿校区", "集美校区"} {
			_, err := cache.GetAvailability(ctx, day.Date, campus)
			require.EqualError(t, err, availabilityFailureMessage)
		}
		require.NoError(t, cache.SetAvailability(ctx, "厦门工艺美院", day))
		for _, campus := range []string{"鼓浪屿校区", "集美校区"} {
			result, err := cache.GetAvailability(ctx, day.Date, campus)
			require.NoError(t, err)
			assert.Equal(t, day.Date, result.Date)
			assert.NotNil(t, result.Rooms)
		}
		require.NoError(t, cache.MarkAvailabilityFailure(ctx, day.Date, "旗山校区"))
		_, err := cache.GetAvailability(ctx, day.Date, "旗山校区")
		require.EqualError(t, err, availabilityFailureMessage)
		require.NoError(t, cache.MarkAvailabilityFailure(ctx, day.Date, "旗山校区", "楼A", "楼B", "楼A"))
		_, err = cache.GetAvailability(ctx, day.Date, "旗山校区")
		require.EqualError(t, err, availabilityFailureMessage+"（失败教学楼：楼A、楼B）")
		// 部分快照与失败楼名同批写入，客户端不能将其读作完整成功结果。
		require.NoError(t, cache.SetAvailability(ctx, "旗山校区", day, "楼B"))
		_, err = cache.GetAvailability(ctx, day.Date, "旗山校区")
		require.EqualError(t, err, availabilityFailureMessage+"（失败教学楼：楼B）")
		var partial DayAvailability
		require.NoError(t, json.Unmarshal([]byte(hook.values[AvailabilityKey(day.Date, "旗山校区")]), &partial))
		require.Equal(t, day.Date, partial.Date)
		require.NoError(t, cache.SetAvailability(ctx, "旗山校区", day))
		_, err = cache.GetAvailability(ctx, day.Date, "旗山校区")
		require.NoError(t, err)
	})
}

func TestAvailabilityRedisFailure(t *testing.T) {
	mockey.PatchConvey("Redis 异常必须传递，不能当作没有失败标记", t, func() {
		mockey.Mock(environment.IsTestEnvironment).Return(false).Build()
		redisErr := errors.New("redis unavailable")
		client := redis.NewClient(&redis.Options{})
		client.AddHook(&availabilityRedisHook{err: redisErr})
		t.Cleanup(func() { require.NoError(t, client.Close()) })
		cache := NewCacheClassroom(client)
		ctx := context.Background()
		_, err := cache.GetAvailability(ctx, "2026-10-12", "旗山校区")
		require.ErrorIs(t, err, redisErr)
		require.ErrorIs(t, cache.MarkAvailabilityFailure(ctx, "2026-10-12", "旗山校区"), redisErr)
		require.ErrorIs(t, cache.SetAvailability(ctx, "旗山校区", &DayAvailability{Date: "2026-10-12"}), redisErr)
	})
}

func TestAvailabilityOldSnapshotCompatibility(t *testing.T) {
	// 已有 Redis 快照仍能读取，旧状态字段自动忽略，空闲区间和占用记录保留。
	const date = "2026-10-12"
	const campus = "旗山校区"
	client := redis.NewClient(&redis.Options{})
	client.AddHook(&availabilityRedisHook{values: map[string]string{
		AvailabilityKey(date, campus): `{"date":"2026-10-12","term":"202601","rooms":[
			{"build":"东2","location":"旗山东2-101","capacity":"100","type":"多媒体",
			"is_free":false,
			"free_ranges":[{"start":1,"end":3},{"start":5,"end":11}],
			"occupancies":[{"start":4,"end":4,"event":"补课","occupant":"","raw_text":""}]},
			{"build":"东2","location":"旗山东2-102","capacity":"100","type":"多媒体",
			"free_ranges":[{"start":1,"end":11}],"occupancies":[]}
		]}`,
	}})
	t.Cleanup(func() { require.NoError(t, client.Close()) })
	day, err := NewCacheClassroom(client).GetAvailability(context.Background(), date, campus)
	require.NoError(t, err)
	require.Len(t, day.Rooms, 2)
	require.Len(t, day.Rooms[0].Occupancies, 1)
	assert.EqualValues(t, 4, day.Rooms[0].Occupancies[0].Start)
	assert.Equal(t, "补课", day.Rooms[0].Occupancies[0].Event)
	assert.Zero(t, day.Rooms[0].Occupancies[0].StartWeek)
	assert.Zero(t, day.Rooms[0].Occupancies[0].EndWeek)
	assert.Empty(t, day.Rooms[1].Occupancies)
	encoded, err := json.Marshal(day)
	require.NoError(t, err)
	assert.Contains(t, string(encoded), `"free_ranges"`)
	assert.NotContains(t, string(encoded), `"is_free"`)
}

func TestSplitAvailability(t *testing.T) {
	day := &DayAvailability{Date: "2026-10-11", Term: "202601", Rooms: []*model.Classroom{
		{Build: "厦门工艺美院", Location: "鼓浪屿多媒体1"},
		{Build: "厦门工艺美院", Location: "集美1-311", Occupancies: []*model.ClassroomOccupancy{{Start: 1, End: 4, Event: "讲座"}}},
	}}
	days := splitAvailability("厦门工艺美院", day)
	require.Len(t, days["鼓浪屿校区"].Rooms, 1)
	require.Len(t, days["集美校区"].Rooms, 1)
	assert.Equal(t, "鼓浪屿校区", days["鼓浪屿校区"].Rooms[0].Build)
	assert.Equal(t, "集美校区", days["集美校区"].Rooms[0].Build)
	assert.Equal(t, "讲座", days["集美校区"].Rooms[0].Occupancies[0].Event)
	assert.Equal(t, "202601", days["集美校区"].Term)
	assert.Equal(t, "厦门工艺美院", day.Rooms[0].Build)
	empty := splitAvailability("厦门工艺美院", &DayAvailability{Date: day.Date, Term: day.Term, Rooms: []*model.Classroom{}})
	assert.NotNil(t, empty["集美校区"].Rooms)
	assert.NotNil(t, empty["鼓浪屿校区"].Rooms)
}
