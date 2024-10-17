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

package db

import (
	"gorm.io/gorm"
)

type Check struct {
	gorm.Model
	Filename string `gorm:"size:255"`      // 文件名
	Path     string `gorm:"size:1024"`     // 上传路径
	User     string `gorm:"size:16"`       // 上传人
	Uuid     string `gorm:"size:40;index"` // uuid
	Status   int    // 审核状态
}

const (
	uncheck = iota
	pass
	reject
)

func AddCheck(filename, path, user, uuid string) {
	DB.Create(&Check{
		Filename: filename,
		Path:     path,
		User:     user,
		Uuid:     uuid,
		Status:   uncheck,
	})
}
