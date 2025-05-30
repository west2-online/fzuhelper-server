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

	"golang.org/x/sync/errgroup"

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *LaunchScreenService) CreateImage(req *launch_screen.CreateImageRequest) (pic *model.Picture, err error) {
	/*
		id, err := s.sf.NextVal()
		if err != nil {
			return nil, fmt.Errorf("LaunchScreen.CreateImage SFCreateIDError:%w", err)
		}
	*/
	suffix, err := utils.GetImageFileType(&req.Image)
	if err != nil {
		return nil, err
	}

	imgUrl, remotePath, err := s.ossClient.GenerateImgName(suffix)
	if err != nil {
		return nil, fmt.Errorf("generate image name failed: %w", err)
	}

	var eg errgroup.Group
	eg.Go(func() error {
		pictureModel := &model.Picture{
			Url:        imgUrl,
			Href:       req.Href,
			Text:       req.Text,
			PicType:    req.PicType,
			ShowTimes:  0,
			PointTimes: 0,
			Duration:   *req.Duration,
			SType:      req.SType,
			Frequency:  req.Frequency,
			StartTime:  req.StartTime,
			EndTime:    req.EndTime,
			Regex:      req.Regex,
			StartAt:    time.Unix(req.StartAt, 0),
			EndAt:      time.Unix(req.EndAt, 0),
		}
		pic, err = s.db.LaunchScreen.CreateImage(s.ctx, pictureModel)
		return err
	})
	eg.Go(func() error {
		/* test stream
		return utils.SaveImageFromBytes(req.Image, "jpg")
		*/
		return s.ossClient.UploadImg(req.Image, remotePath)
	})
	if err = eg.Wait(); err != nil {
		return nil, fmt.Errorf("LaunchScreenService.CreateImage error:%w", err)
	}
	return pic, nil
}
