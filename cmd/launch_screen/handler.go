package main

import (
	"context"
	"github.com/cloudwego/kitex/client"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/dal/db"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/pack"
	"github.com/west2-online/fzuhelper-server/cmd/launch_screen/service"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen/launchscreenservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"golang.org/x/sync/errgroup"
)

// LaunchScreenServiceImpl implements the last service interface defined in the IDL.
type LaunchScreenServiceImpl struct {
	launchScreenCli launchscreenservice.Client
}

func NewLaunchScreenClient(addr string) (launchscreenservice.Client, error) {
	return launchscreenservice.NewClient(constants.LaunchScreenServiceName, client.WithHostPorts(addr))
}

// CreateImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) CreateImage(ctx context.Context, req *launch_screen.CreateImageRequest) (resp *launch_screen.CreateImageResponse, err error) {
	resp = new(launch_screen.CreateImageResponse)
	claim, err := utils.CheckToken(req.Token)
	if err != nil {
		resp.Base = utils.BuildBaseResp(err)
		return resp, nil
	}
	uid := claim.UserId
	imgUrl := pack.GenerateImgName(uid)
	pic := new(db.Picture)

	var eg errgroup.Group
	eg.Go(func() error {
		picture := &model.Picture{
			UserId:    uid,
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
		}
		pic, err = service.NewLaunchScreenService(ctx).PutImage(picture)
		if err != nil {
			return err
		}
		return nil
	})
	eg.Go(func() error {
		err = pack.UploadImg(req.Image, imgUrl)
		if err != nil {
			return err
		}
		return nil
	})
	if err = eg.Wait(); err != nil {
		resp.Base = utils.BuildBaseResp(err)
		return resp, nil
	}
	resp.Base = utils.BuildBaseResp(nil)
	resp.Picture = service.BuildImageResp(pic)
	return resp, nil
}

// GetImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) GetImage(ctx context.Context, req *launch_screen.GetImageRequest) (resp *launch_screen.GetImagesByUserIdResponse, err error) {
	// TODO: Your code here...
	return
}

// GetImagesByUserId implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) GetImagesByUserId(ctx context.Context, req *launch_screen.GetImagesByUserIdRequest) (resp *launch_screen.GetImagesByUserIdResponse, err error) {
	// TODO: Your code here...
	return
}

// ChangeImageProperty implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) ChangeImageProperty(ctx context.Context, req *launch_screen.ChangeImagePropertyRequest) (resp *launch_screen.ChangeImagePropertyResponse, err error) {
	// TODO: Your code here...
	return
}

// ChangeImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) ChangeImage(ctx context.Context, req *launch_screen.ChangeImageRequest) (resp *launch_screen.ChangeImageResponse, err error) {
	// TODO: Your code here...
	return
}

// DeleteImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) DeleteImage(ctx context.Context, req *launch_screen.DeleteImageRequest) (resp *launch_screen.DeleteImageResponse, err error) {
	// TODO: Your code here...
	return
}

// MobileGetImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) MobileGetImage(ctx context.Context, req *launch_screen.MobileGetImageRequest) (resp *launch_screen.MobileGetImageResponse, err error) {
	// TODO: Your code here...
	return
}

// AddImagePointTime implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) AddImagePointTime(ctx context.Context, req *launch_screen.AddImagePointTimeRequest) (resp *launch_screen.AddImagePointTimeResponse, err error) {
	// TODO: Your code here...
	return
}
