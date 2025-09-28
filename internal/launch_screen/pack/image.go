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
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	db "github.com/west2-online/fzuhelper-server/pkg/db/model"
)

func BuildImageResp(dbP *db.Picture) *model.Picture {
	return &model.Picture{
		Id:  dbP.ID,
		Url: dbP.Url,
		// Url:        strings.Join([]string{config.UpYun.DownloadDomain, strings.TrimPrefix(dbP.Url, config.UpYun.UssDomain)}, ""),
		Href:       dbP.Href,
		Text:       dbP.Text,
		Type:       dbP.PicType,
		ShowTimes:  &dbP.ShowTimes,
		PointTimes: &dbP.PointTimes,
		Duration:   dbP.Duration,
		SType:      &dbP.SType,
		Frequency:  dbP.Frequency,
		StartAt:    dbP.StartAt.Unix(),
		EndAt:      dbP.EndAt.Unix(),
		StartTime:  dbP.StartTime,
		EndTime:    dbP.EndTime,
		Regex:      dbP.Regex,
	}
}

func BuildImagesResp(dbPictures *[]db.Picture) []*model.Picture {
	var pictureList []*model.Picture
	for _, msg := range *dbPictures {
		pictureList = append(pictureList, BuildImageResp(&msg))
	}
	return pictureList
}
