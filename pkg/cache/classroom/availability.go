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

const availabilityFailureValue = "failure"

var ErrAvailabilityFailure = errors.New("classroom availability refresh failed")

func AvailabilityKey(date, campus string) string {
	return fmt.Sprintf("classroom:availability:%s:%s", date, campus)
}

func (c *CacheClassroom) GetAvailability(ctx context.Context, date, campus string) (*DayAvailability, error) {
	data, err := c.client.Get(ctx, AvailabilityKey(date, campus)).Bytes()
	if err != nil {
		return nil, fmt.Errorf("get classroom availability: %w", err)
	}
	if string(data) == availabilityFailureValue {
		return nil, ErrAvailabilityFailure
	}
	var day DayAvailability
	if err := json.Unmarshal(data, &day); err != nil {
		return nil, fmt.Errorf("decode classroom availability: %w", err)
	}
	return &day, nil
}

func (c *CacheClassroom) SetAvailabilityFailure(ctx context.Context, date, campus string) error {
	if environment.IsTestEnvironment() {
		return nil
	}
	_, err := c.client.Set(ctx, AvailabilityKey(date, campus), availabilityFailureValue, constants.ClassroomKeyExpire).Result()
	if err != nil {
		return fmt.Errorf("set classroom availability failure: %w", err)
	}
	return nil
}

func (c *CacheClassroom) SetAvailability(ctx context.Context, campus string, day *DayAvailability) error {
	if environment.IsTestEnvironment() {
		return nil
	}
	data, err := json.Marshal(day)
	if err != nil {
		return fmt.Errorf("encode classroom availability: %w", err)
	}
	_, err = c.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Set(ctx, AvailabilityKey(day.Date, campus), data, constants.ClassroomKeyExpire)
		return nil
	})
	if err != nil {
		return fmt.Errorf("set classroom availability: %w", err)
	}
	return nil
}
