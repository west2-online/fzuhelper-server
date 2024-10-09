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

package service

import (
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"time"
)

func (s *LaunchScreenService) PutImage(picture *model.Picture) (*db.Picture, error) {
	Loc, _ := time.LoadLocation("Asia/Shanghai")
	pictureModel := &db.Picture{
		ID:         picture.Id,
		Uid:        picture.UserId,
		Url:        picture.Url,
		Href:       picture.Href,
		Text:       picture.Text,
		PicType:    picture.PicType,
		ShowTimes:  0,
		PointTimes: 0,
		Duration:   picture.Duration,
		SType:      *picture.SType,
		Frequency:  picture.Frequency,
		StartTime:  picture.StartTime,
		EndTime:    picture.EndTime,
		StartAt:    time.Unix(picture.StartAt, 0).In(Loc),
		EndAt:      time.Unix(picture.EndAt, 0).In(Loc),
	}
	return db.CreateImage(s.ctx, pictureModel)
}
