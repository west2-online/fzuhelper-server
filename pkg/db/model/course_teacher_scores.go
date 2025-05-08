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

type CourseTeacherScores struct {
	Id          int64
	StuIdSha256 string
	CourseName  string
	TeacherName string
	Semester    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `sql:"index"`
}

// CourseTeacherScoreRecord 用于承载一行待写入的数据
type CourseTeacherScoreRecord struct {
	ID          int64
	StuIdSha256 string
	CourseName  string
	TeacherName string
	Semester    string
	Score       float64
}
