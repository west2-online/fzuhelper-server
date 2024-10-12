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
	api "github.com/west2-online/fzuhelper-server/cmd/api/biz/model/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func BuildLaunchScreen(image *model.Picture) *api.Picture {
	return &api.Picture{
		ID:         image.Id,
		UserID:     image.UserId,
		URL:        image.Url,
		PicType:    image.PicType,
		Duration:   image.Duration,
		Href:       image.Href,
		ShowTimes:  image.ShowTimes,
		SType:      image.SType,
		Frequency:  image.Frequency,
		Text:       image.Text,
		StartAt:    image.StartAt,
		EndAt:      image.EndAt,
		StartTime:  image.StartTime,
		EndTime:    image.EndTime,
		StudentID:  image.StudentId,
		DeviceType: image.DeviceType,
	}
}

func BuildLaunchScreenList(kitexPictures []*model.Picture) []*api.Picture {
	imagesResp := make([]*api.Picture, 0)
	for _, v := range kitexPictures {
		imagesResp = append(imagesResp, BuildLaunchScreen(v))
	}
	return imagesResp
}
