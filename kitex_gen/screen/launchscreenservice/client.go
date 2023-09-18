// Code generated by Kitex v0.7.1. DO NOT EDIT.

package launchscreenservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	screen "github.com/west2-online/fzuhelper-server/kitex_gen/screen"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	PictureCreate(ctx context.Context, req *screen.CreatePictureRequest, callOptions ...callopt.Option) (r *screen.CreatePictureResponse, err error)
	PictureGet(ctx context.Context, req *screen.GetPictureRequest, callOptions ...callopt.Option) (r *screen.GetPictureResponse, err error)
	PictureUpdate(ctx context.Context, req *screen.PutPictureRequset, callOptions ...callopt.Option) (r *screen.PutPictureResponse, err error)
	PictureImgUpdate(ctx context.Context, req *screen.PutPictureImgRequset, callOptions ...callopt.Option) (r *screen.PutPictureResponse, err error)
	PictureDelte(ctx context.Context, req *screen.DeletePictureRequest, callOptions ...callopt.Option) (r *screen.DeletePictureResponse, err error)
	RetPicture(ctx context.Context, req *screen.RetPictureRequest, callOptions ...callopt.Option) (r *screen.RetPictureResponse, err error)
	AddPoint(ctx context.Context, req *screen.AddPointRequest, callOptions ...callopt.Option) (r *screen.AddPointResponse, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kLaunchScreenServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kLaunchScreenServiceClient struct {
	*kClient
}

func (p *kLaunchScreenServiceClient) PictureCreate(ctx context.Context, req *screen.CreatePictureRequest, callOptions ...callopt.Option) (r *screen.CreatePictureResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PictureCreate(ctx, req)
}

func (p *kLaunchScreenServiceClient) PictureGet(ctx context.Context, req *screen.GetPictureRequest, callOptions ...callopt.Option) (r *screen.GetPictureResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PictureGet(ctx, req)
}

func (p *kLaunchScreenServiceClient) PictureUpdate(ctx context.Context, req *screen.PutPictureRequset, callOptions ...callopt.Option) (r *screen.PutPictureResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PictureUpdate(ctx, req)
}

func (p *kLaunchScreenServiceClient) PictureImgUpdate(ctx context.Context, req *screen.PutPictureImgRequset, callOptions ...callopt.Option) (r *screen.PutPictureResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PictureImgUpdate(ctx, req)
}

func (p *kLaunchScreenServiceClient) PictureDelte(ctx context.Context, req *screen.DeletePictureRequest, callOptions ...callopt.Option) (r *screen.DeletePictureResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.PictureDelte(ctx, req)
}

func (p *kLaunchScreenServiceClient) RetPicture(ctx context.Context, req *screen.RetPictureRequest, callOptions ...callopt.Option) (r *screen.RetPictureResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.RetPicture(ctx, req)
}

func (p *kLaunchScreenServiceClient) AddPoint(ctx context.Context, req *screen.AddPointRequest, callOptions ...callopt.Option) (r *screen.AddPointResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.AddPoint(ctx, req)
}
