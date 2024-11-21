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
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *LaunchScreenService) UpdateImagePath(req *launch_screen.ChangeImageRequest) (pic *model.Picture, err error) {
	origin, err := s.db.LaunchScreen.GetImageById(s.ctx, req.PictureId)
	if err != nil {
		return nil, errors.Errorf("LaunchScreenService.UpdateImagePath db.GetImageById error: %v", err)
	}

	delUrl := origin.Url

	suffix, err := utils.GetImageFileType(&req.Image)
	if err != nil {
		return nil, err
	}

	imgUrl := upyun.GenerateImgName(suffix)

	var eg errgroup.Group
	var err2 error
	eg.Go(func() error {
		err = upyun.DeleteImg(delUrl)
		if err != nil {
			return err
		}
		err = upyun.UploadImg(req.Image, imgUrl)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		origin.Url = imgUrl
		pic, err2 = s.db.LaunchScreen.UpdateImage(s.ctx, origin)
		return err2
	})
	if err = eg.Wait(); err != nil {
		return nil, errors.Errorf("LaunchScreenService.UpdateImagePath error: %v", err)
	}
	return pic, nil
}
