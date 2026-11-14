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

// availabilityRedisHook 模拟 Redis 命令执行，验证快照读写，无需连接真实 Redis。
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

func TestAvailabilityFailureValue(t *testing.T) {
	mockey.PatchConvey("失败值和快照共用同一个 key", t, func() {
		mockey.Mock(environment.IsTestEnvironment).Return(false).Build()
		hook := &availabilityRedisHook{values: make(map[string]string)}
		client := redis.NewClient(&redis.Options{})
		client.AddHook(hook)
		t.Cleanup(func() { require.NoError(t, client.Close()) })
		cache := NewCacheClassroom(client)
		ctx := context.Background()
		day := &DayAvailability{Date: "2026-10-12", Term: "202601", Rooms: []*model.Classroom{}}
		require.NoError(t, cache.SetAvailability(ctx, "旗山校区", day))
		result, err := cache.GetAvailability(ctx, day.Date, "旗山校区")
		require.NoError(t, err)
		assert.Equal(t, day.Date, result.Date)
		assert.NotNil(t, result.Rooms)
		require.NoError(t, cache.SetAvailabilityFailure(ctx, day.Date, "旗山校区"))
		_, err = cache.GetAvailability(ctx, day.Date, "旗山校区")
		require.ErrorIs(t, err, ErrAvailabilityFailure)
		require.NoError(t, cache.SetAvailability(ctx, "旗山校区", day))
		result, err = cache.GetAvailability(ctx, day.Date, "旗山校区")
		require.NoError(t, err)
		assert.Equal(t, day.Date, result.Date)
		_, err = cache.GetAvailability(ctx, day.Date, "没有数据的校区")
		require.ErrorIs(t, err, redis.Nil)
	})
}

func TestAvailabilityRedisFailure(t *testing.T) {
	mockey.PatchConvey("Redis 异常必须传递", t, func() {
		mockey.Mock(environment.IsTestEnvironment).Return(false).Build()
		redisErr := errors.New("redis unavailable")
		client := redis.NewClient(&redis.Options{})
		client.AddHook(&availabilityRedisHook{err: redisErr})
		t.Cleanup(func() { require.NoError(t, client.Close()) })
		cache := NewCacheClassroom(client)
		ctx := context.Background()
		_, err := cache.GetAvailability(ctx, "2026-10-12", "旗山校区")
		require.ErrorIs(t, err, redisErr)
		require.ErrorIs(t, cache.SetAvailabilityFailure(ctx, "2026-10-12", "旗山校区"), redisErr)
		require.ErrorIs(t, cache.SetAvailability(ctx, "旗山校区", &DayAvailability{Date: "2026-10-12", Rooms: []*model.Classroom{}}), redisErr)
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
