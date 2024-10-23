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
	"fmt"
	"time"

	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *LaunchScreenService) UpdateImageProperty(req *launch_screen.ChangeImagePropertyRequest, origin *db.Picture) (*db.Picture, error) {
	Loc := utils.LoadCNLocation()
	origin.PicType = req.PicType
	origin.SType = req.SType
	origin.Duration = *req.Duration
	origin.Href = *req.Href
	origin.Frequency = req.Frequency
	origin.Text = req.Text
	origin.StartAt = time.Unix(req.StartAt, 0).In(Loc)
	origin.EndAt = time.Unix(req.EndAt, 0).In(Loc)
	origin.StartTime = req.StartTime
	origin.EndTime = req.EndTime
	origin.Regex = req.Regex
	pic, err := db.UpdateImage(s.ctx, origin)
	if err != nil {
		return nil, fmt.Errorf("LaunchScreenService.UpdateImageProperty error: %v", err)
	}
	return pic, nil
}
