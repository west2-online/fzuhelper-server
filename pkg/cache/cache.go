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

	"github.com/redis/go-redis/v9"

	"github.com/west2-online/fzuhelper-server/pkg/cache/academic"
	"github.com/west2-online/fzuhelper-server/pkg/cache/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/cache/common"
	"github.com/west2-online/fzuhelper-server/pkg/cache/course"
	"github.com/west2-online/fzuhelper-server/pkg/cache/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/cache/paper"
	"github.com/west2-online/fzuhelper-server/pkg/cache/user"
	"github.com/west2-online/fzuhelper-server/pkg/cache/version"
)

type Cache struct {
	client       *redis.Client
	Classroom    *classroom.CacheClassroom
	Paper        *paper.CachePaper
	LaunchScreen *launch_screen.CacheLaunchScreen
	Academic     *academic.CacheAcademic
	User         *user.CacheUser
	Common       *common.CacheCommon
	Course       *course.CacheCourse
	Version      *version.CacheVersion
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client:       client,
		Classroom:    classroom.NewCacheClassroom(client),
		LaunchScreen: launch_screen.NewCacheLaunchScreen(client),
		Paper:        paper.NewCachePaper(client),
		Academic:     academic.NewCacheAcademic(client),
		User:         user.NewCacheUser(client),
		Common:       common.NewCacheCommon(client),
		Course:       course.NewCacheCourse(client),
		Version:      version.NewCacheVersion(client),
	}
}

// IsKeyExist will check if key exist
func (c *Cache) IsKeyExist(ctx context.Context, key string) bool {
	return c.client.Exists(ctx, key).Val() == 1
}
