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

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/jwch"
	"github.com/west2-online/yjsy"
)

func (c *CacheCourse) SetCoursesCache(ctx context.Context, key string, course *[]*jwch.Course) error {
	coursesJson, err := sonic.Marshal(course)
	if err != nil {
		logger.Errorf("dal.SetCoursesCache: Marshal info failed: %v", err)
		return err
	}
	if err = c.client.Set(ctx, key, coursesJson, constants.CourseTermsKeyExpire).Err(); err != nil {
		logger.Errorf("dal.SetCoursesCache: Set info failed: %v", err)
		return err
	}
	return nil
}

func (c *CacheCourse) SetCoursesCacheYjsy(ctx context.Context, key string, course *[]*yjsy.Course) error {
	coursesJson, err := sonic.Marshal(course)
	if err != nil {
		logger.Errorf("dal.SetCoursesCacheYjsy: Marshal info failed: %v", err)
		return err
	}
	if err = c.client.Set(ctx, key, coursesJson, constants.CourseTermsKeyExpire).Err(); err != nil {
		logger.Errorf("dal.SetCoursesCacheYjsy: Set info failed: %v", err)
		return err
	}
	return nil
}
