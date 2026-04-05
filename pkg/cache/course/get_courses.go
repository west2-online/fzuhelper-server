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

package course

import (
	"context"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (c *CacheCourse) GetCoursesCache(ctx context.Context, key string) (course []*jwch.Course, err error) {
	course = make([]*jwch.Course, 0)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errno.Errorf(errno.InternalRedisErrorCode, "dal.GetCoursesCache: Cache failed: %v", err)
	}
	if err = sonic.Unmarshal(data, &course); err != nil {
		return nil, errno.Errorf(errno.InternalJSONErrorCode, "dal.GetCoursesCache: Unmarshal failed: %v", err)
	}
	return course, nil
}

func (c *CacheCourse) GetCoursesCacheYjsy(ctx context.Context, key string) (course []*yjsy.Course, err error) {
	course = make([]*yjsy.Course, 0)
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		return nil, errno.Errorf(errno.InternalRedisErrorCode, "dal.GetCoursesCacheYjsy: Cache failed: %v", err)
	}
	if err = sonic.Unmarshal(data, &course); err != nil {
		return nil, errno.Errorf(errno.InternalJSONErrorCode, "dal.GetCoursesCacheYjsy: Unmarshal failed: %v", err)
	}
	return course, nil
}
