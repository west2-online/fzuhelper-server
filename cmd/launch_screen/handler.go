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

package main

import (
	"context"
	"encoding/binary"
	"regexp"

	"github.com/bytedance/sonic"

	"github.com/west2-online/fzuhelper-server/pkg/errno"

	"github.com/west2-online/fzuhelper-server/pkg/upcloud"

	"golang.org/x/sync/errgroup"

	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/pack"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/service"
	launch_screen "github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

// LaunchScreenServiceImpl implements the last service interface defined in the IDL.
type LaunchScreenServiceImpl struct{}

// CreateImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) CreateImage(ctx context.Context, req *launch_screen.CreateImageRequest) (resp *launch_screen.CreateImageResponse, err error) {
	resp = new(launch_screen.CreateImageResponse)

	imgUrl := upcloud.GenerateImgName(req.StuId)
	pic := new(db.Picture)

	var eg errgroup.Group
	eg.Go(func() error {
		picture := &model.Picture{
			Url:       imgUrl,
			Href:      *req.Href,
			Text:      req.Text,
			PicType:   req.PicType,
			Duration:  *req.Duration,
			SType:     &req.SType,
			Frequency: req.Frequency,
			StartTime: req.StartTime,
			EndTime:   req.EndTime,
			StartAt:   req.StartAt,
			EndAt:     req.EndAt,
			Regex:     req.Regex,
		}
		pic, err = service.NewLaunchScreenService(ctx).CreateImage(picture)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		err = upcloud.UploadImg(req.Image, imgUrl)
		if err != nil {
			return err
		}
		return nil
	})
	if err = eg.Wait(); err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Picture = service.BuildImageResp(pic)
	return resp, nil
}

// GetImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) GetImage(ctx context.Context, req *launch_screen.GetImageRequest) (resp *launch_screen.GetImageResponse, err error) {
	resp = new(launch_screen.GetImageResponse)

	pic, err := service.NewLaunchScreenService(ctx).GetImageById(req.PictureId)
	resp.Base = pack.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Picture = service.BuildImageResp(pic)
	return resp, nil
}

// GetImagesByUserId implements the LaunchScreenServiceImpl interface.
// 这个接口实际用不到
func (s *LaunchScreenServiceImpl) GetImagesByUserId(ctx context.Context, req *launch_screen.GetImagesByUserIdRequest) (resp *launch_screen.GetImagesByUserIdResponse, err error) {
	resp = new(launch_screen.GetImagesByUserIdResponse)

	pic, cnt, err := service.NewLaunchScreenService(ctx).GetImagesByStuId(req.StuId)
	resp.Base = pack.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Count = &cnt
	resp.PictureList = service.BuildImagesResp(pic)
	return resp, nil
}

// ChangeImageProperty implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) ChangeImageProperty(ctx context.Context, req *launch_screen.ChangeImagePropertyRequest) (resp *launch_screen.ChangeImagePropertyResponse, err error) {
	resp = new(launch_screen.ChangeImagePropertyResponse)

	origin, err := service.NewLaunchScreenService(ctx).GetImageById(req.PictureId)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	pic, err := service.NewLaunchScreenService(ctx).UpdateImageProperty(req, origin)
	resp.Base = pack.BuildBaseResp(err)
	if err != nil {
		return resp, nil
	}
	resp.Picture = service.BuildImageResp(pic)
	return resp, nil
}

// ChangeImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) ChangeImage(ctx context.Context, req *launch_screen.ChangeImageRequest) (resp *launch_screen.ChangeImageResponse, err error) {
	resp = new(launch_screen.ChangeImageResponse)

	origin, err := service.NewLaunchScreenService(ctx).GetImageById(req.PictureId)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	delUrl := origin.Url
	imgUrl := upcloud.GenerateImgName(req.StuId)
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
	pic := new(db.Picture)
	eg.Go(func() error {
		pic, err = service.NewLaunchScreenService(ctx).UpdateImagePath(imgUrl, origin)
		return err
	})
	if err = eg.Wait(); err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Picture = service.BuildImageResp(pic)
	return resp, nil
}

// DeleteImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) DeleteImage(ctx context.Context, req *launch_screen.DeleteImageRequest) (resp *launch_screen.DeleteImageResponse, err error) {
	resp = new(launch_screen.DeleteImageResponse)

	pic, err := service.NewLaunchScreenService(ctx).DeleteImage(req.PictureId)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	if err = upcloud.DeleteImg(pic.Url); err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Picture = service.BuildImageResp(pic)
	return resp, nil
}

// MobileGetImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) MobileGetImage(ctx context.Context, req *launch_screen.MobileGetImageRequest) (resp *launch_screen.MobileGetImageResponse, err error) {
	resp = new(launch_screen.MobileGetImageResponse)
	pictureList, cnt, err, isGetFromMysql := service.NewLaunchScreenService(ctx).MobileGetImage(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		return resp, nil
	}

	if !isGetFromMysql {
		resp.Count = &cnt
		resp.PictureList = service.BuildImagesResp(pictureList)
		return resp, nil
	}

	var respList []db.Picture
	var cntCurrent int64 = 0
	for _, picture := range *pictureList {
		match := false
		m := make(map[string]string)
		if err = sonic.Unmarshal([]byte(picture.Regex), &m); err != nil {
			resp.Base = pack.BuildBaseResp(errno.InternalServiceError)
			return resp, nil
		}
		for k, v := range m {
			if k == "picture_id" || k == "_id" {
				continue
			}
			stuId := make([]byte, 8)
			binary.LittleEndian.PutUint64(stuId, uint64(req.StudentId))
			match, err = regexp.Match(v, stuId)
			if err != nil {
				continue
			}
			if !match {
				continue
			}
		}
		if match {
			cntCurrent++
			respList = append(respList, picture)
		}
	}

	resp.Count = &cntCurrent
	resp.PictureList = service.BuildImagesResp(&respList)
	return resp, nil
}

// AddImagePointTime implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) AddImagePointTime(ctx context.Context, req *launch_screen.AddImagePointTimeRequest) (resp *launch_screen.AddImagePointTimeResponse, err error) {
	resp = new(launch_screen.AddImagePointTimeResponse)
	err = service.NewLaunchScreenService(ctx).AddPointTime(req.PictureId)
	resp.Base = pack.BuildBaseResp(err)
	return resp, nil
}
