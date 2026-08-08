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

package model

import (
	"time"

	"gorm.io/gorm"
)

type UserCourse struct {
	Id                int64
	StuId             string
	Term              string
	TermCourses       string
	TermCoursesSha256 string
	ExamInfo          string
	ExamInfoSHA256    string
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `sql:"index"`
}

type ExamOffering struct {
	ID        int64          `json:"id"`
	ExamHash  string         `json:"exam_hash"`
	Tag       string         `json:"tag"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at,omitempty"`
}

type UserTerm struct {
	Id        int64
	StuId     string
	TermTime  string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `sql:"index"`
}

type AutoAdjustCourse struct {
	Id          int64
	Year        string
	FromDate    string
	ToDate      *string
	Term        string
	FromWeek    int64
	ToWeek      *int64
	FromWeekday int64
	ToWeekday   *int64
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `sql:"index"`
}

// UserCustomCourse 用户自定义课程表
type UserCustomCourse struct {
	ID         int64          `gorm:"column:id;primaryKey;autoIncrement"`
	StuId      string         `gorm:"column:stu_id;size:50;not null"`
	Term       string         `gorm:"column:term;size:20;not null"`
	CourseId   string         `gorm:"column:course_id;size:64;not null"`
	Name       string         `gorm:"column:name;size:100;not null"`
	Teacher    string         `gorm:"column:teacher;size:50"`
	Location   string         `gorm:"column:location;size:100;not null"`
	StartClass int            `gorm:"column:start_class;not null"`
	EndClass   int            `gorm:"column:end_class;not null"`
	StartWeek  int            `gorm:"column:start_week;not null"`
	EndWeek    int            `gorm:"column:end_week;not null"`
	Weekday    int            `gorm:"column:weekday;not null"`
	IsSingle   bool           `gorm:"column:is_single;default:false"`
	IsDouble   bool           `gorm:"column:is_double;default:false"`
	Color      string         `gorm:"column:color;size:20;default:#FF5733"`
	Remark     string         `gorm:"column:remark;size:200"`
	CreatedAt  time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"column:deleted_at"`
}

// TableName 指定表名
func (UserCustomCourse) TableName() string {
	return "user_custom_courses"
}
