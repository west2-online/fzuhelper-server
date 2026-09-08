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

package pack

import (
	"strconv"

	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	dbModel "github.com/west2-online/fzuhelper-server/pkg/db/model"
)

// BuildCustomCourseItems 将数据库模型转换为 Thrift 类型
func BuildCustomCourseItems(courses []*dbModel.UserCustomCourse) []*kitexModel.CustomCourse {
	result := make([]*kitexModel.CustomCourse, 0, len(courses))
	for _, c := range courses {
		id := strconv.FormatInt(c.Id, 10)
		item := &kitexModel.CustomCourse{
			Id:         &id,
			Name:       c.Name,
			Location:   c.Location,
			StartClass: int32(c.StartClass),
			EndClass:   int32(c.EndClass),
			StartWeek:  int32(c.StartWeek),
			EndWeek:    int32(c.EndWeek),
			Weekday:    int32(c.Weekday),
		}
		if c.Teacher != "" {
			item.Teacher = &c.Teacher
		}
		item.Single = c.IsSingle
		item.Double_ = c.IsDouble
		if c.Color != "" {
			item.Color = &c.Color
		}
		if c.Remark != "" {
			item.Remark = &c.Remark
		}
		result = append(result, item)
	}
	return result
}
