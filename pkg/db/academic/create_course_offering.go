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

package academic

import (
	"context"
	"fmt"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

func (c *DBAcademic) CreateCourseOffering(ctx context.Context, course *model.CourseOffering) (*model.CourseOffering, error) {
	if err := c.client.WithContext(ctx).Table(constants.CourseOfferingsTableName).Create(course).Error; err != nil {
		return nil, errno.NewErrNo(errno.InternalDatabaseErrorCode, fmt.Sprintf("dal.CreateCourseOffering error: %v", err))
	}
	return course, nil
}
