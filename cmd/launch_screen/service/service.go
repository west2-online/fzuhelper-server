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
	"context"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

type LaunchScreenService struct {
	ctx context.Context
}

func NewLaunchScreenService(ctx context.Context) *LaunchScreenService {
	return &LaunchScreenService{ctx: ctx}
}

func BuildImageResp(dbP *db.Picture) *model.Picture {
	return &model.Picture{
		Id:         dbP.ID,
		Url:        dbP.Url,
		Href:       dbP.Href,
		Text:       dbP.Text,
		PicType:    dbP.PicType,
		ShowTimes:  &dbP.ShowTimes,
		PointTimes: &dbP.PointTimes,
		Duration:   dbP.Duration,
		SType:      &dbP.SType,
		Frequency:  dbP.Frequency,
	}
}
