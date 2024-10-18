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

	"golang.org/x/sync/errgroup"

	"github.com/west2-online/fzuhelper-server/pkg/upcloud"

	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
)

func (s *LaunchScreenService) UpdateImagePath(req *launch_screen.ChangeImageRequest) (pic *db.Picture, err error) {
	origin, err := db.GetImageById(s.ctx, req.PictureId)
	if err != nil {
		return nil, fmt.Errorf("LaunchScreenService.UpdateImagePath db.GetImageById error: %v", err)
	}

	delUrl := origin.Url
	imgUrl := upcloud.GenerateImgName(req.PictureId)

	var eg errgroup.Group
	eg.Go(func() error {
		err = upcloud.DeleteImg(delUrl)
		if err != nil {
			return err
		}
		err = upcloud.UploadImg(req.Image, imgUrl)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		origin.Url = imgUrl
		pic, err = db.UpdateImage(s.ctx, origin)
		return err
	})
	if err = eg.Wait(); err != nil {
		return nil, fmt.Errorf("LaunchScreenService.UpdateImagePath error: %v", err)
	}
	return pic, nil
}
