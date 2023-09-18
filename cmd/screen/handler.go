package main

import (
	"bytes"
	"context"

	"github.com/west2-online/fzuhelper-server/cmd/screen/pack"
	"github.com/west2-online/fzuhelper-server/cmd/screen/service"
	screen "github.com/west2-online/fzuhelper-server/kitex_gen/screen"
)

// LaunchScreenServiceImpl implements the last service interface defined in the IDL.
type LaunchScreenServiceImpl struct{}

// PictureCreate implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) PictureCreate(ctx context.Context, req *screen.CreatePictureRequest) (resp *screen.CreatePictureResponse, err error) {
	// TODO: Your code here...
	// 校验token
	// _, err = utils.CheckToken(req.Token)
	// if err != nil {
	// 	resp.Base = pack.BuildBaseResp(err)
	// 	resp.Picture = nil
	// 	return resp, err
	// }
	img := bytes.NewReader(req.Imgfile)
	picture, err := service.NewScreenService(ctx).CreatePicture(req, img)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		resp.Picture = nil
		return resp, err
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Picture = pack.ConvertPicture(picture)
	return resp, nil
}

// PictureGet implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) PictureGet(ctx context.Context, req *screen.GetPictureRequest) (resp *screen.GetPictureResponse, err error) {
	// TODO: Your code here...
	// _, err = utils.CheckToken(req.Token)
	// if err != nil {
	// 	resp.Base = pack.BuildBaseResp(err)
	// 	resp.Picture = nil
	// 	return resp, err
	// }

	pictures, err := service.NewScreenService(ctx).GetPicture(req.PictureId)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		resp.Picture = nil
		return resp, err
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Picture = pack.BuildPicturesResp(pictures)
	resp.Total = int64(len(resp.Picture))
	return
}

// PictureUpdate implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) PictureUpdate(ctx context.Context, req *screen.PutPictureRequset) (resp *screen.PutPictureResponse, err error) {
	// TODO: Your code here...
	// _, err = utils.CheckToken(req.Token)
	// if err != nil {
	// 	resp.Base = pack.BuildBaseResp(err)
	// 	resp.Picture = nil
	// 	return
	// }
	picture, err := service.NewScreenService(ctx).UpdatePicture(req)

	resp.Base = pack.BuildBaseResp(nil)
	resp.Picture = pack.ConvertPicture(picture)
	return
}

// PictureImgUpdate implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) PictureImgUpdate(ctx context.Context, req *screen.PutPictureImgRequset) (resp *screen.PutPictureResponse, err error) {
	// TODO: Your code here...

	// _, err = utils.CheckToken(req.Token)
	// if err != nil {
	// 	resp.Base = pack.BuildBaseResp(err)
	// 	resp.Picture = nil
	// 	return
	// }
	img := bytes.NewReader(req.Imgfile)
	picture, err := service.NewScreenService(ctx).UpdatePictureImg(req, img)
	resp.Base = pack.BuildBaseResp(nil)
	resp.Picture = pack.ConvertPicture(picture)
	return
}

// PictureDelte implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) PictureDelte(ctx context.Context, req *screen.DeletePictureRequest) (resp *screen.DeletePictureResponse, err error) {
	// TODO: Your code here...
	// _, err = utils.CheckToken(req.Token)
	// if err != nil {
	// 	resp.Base = pack.BuildBaseResp(err)
	// 	resp.Picture = nil
	// 	return
	// }
	picture, err := service.NewScreenService(ctx).DeletePicture(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		resp.Picture = nil
		return
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Picture = pack.ConvertPicture(picture)
	return
}

// RetPicture implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) RetPicture(ctx context.Context, req *screen.RetPictureRequest) (resp *screen.RetPictureResponse, err error) {
	// TODO: Your code here...
	// _, err = utils.CheckToken(req.Token)
	// if err != nil {
	// 	resp.Base = pack.BuildBaseResp(err)
	// 	resp.Picture = nil
	// 	return
	// }
	imgs, err := service.NewScreenService(ctx).RetPicture(req)
	if err != nil {
		resp.Base = pack.BuildBaseResp(err)
		resp.Picture = nil
		return
	}
	resp.Base = pack.BuildBaseResp(nil)
	resp.Picture = pack.BuildPicturesResp(imgs)
	resp.Total = int64(len(resp.Picture))
	return
}

// AddPoint implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) AddPoint(ctx context.Context, req *screen.AddPointRequest) (resp *screen.AddPointResponse, err error) {
	// TODO: Your code here...
	// _, err = utils.CheckToken(req.Token)
	// if err != nil {
	// 	resp.Base = pack.BuildBaseResp(err)
	// 	return
	// }

	service.NewScreenService(ctx).AddPoint(req.PictureId)
	return
}
