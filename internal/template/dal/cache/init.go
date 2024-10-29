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

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/config"
)

func Init() {
	ctx := context.Background()

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       0,
	})

	err := rdb.Set(ctx, "test", "just for test", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := rdb.Get(ctx, "test").Result()
	if err != nil {
		panic(err)
	}

	klog.Infof("val: %v\n", val)

	val, err = rdb.Get(ctx, "test1").Result()

	if err == redis.Nil {
		klog.Info("Not found test1 key")
	} else if err != nil {
		panic(err)
	}

	klog.Infof("val: %v\n", val)
}
