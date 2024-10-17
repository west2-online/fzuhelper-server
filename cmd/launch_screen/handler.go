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

	launch_screen "github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
)

// LaunchScreenServiceImpl implements the last service interface defined in the IDL.
type LaunchScreenServiceImpl struct{}

func (s *LaunchScreenServiceImpl) CreateImage(stream launch_screen.LaunchScreenService_CreateImageServer) (err error) {
	println("CreateImage called")
	return
}

// GetImage implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) GetImage(ctx context.Context, req *launch_screen.GetImageRequest) (resp *launch_screen.GetImageResponse, err error) {
	// TODO: Your code here...
	return
}

// ChangeImageProperty implements the LaunchScreenServiceImpl interface.
func (s *LaunchScreenServiceImpl) ChangeImageProperty(ctx context.Context, req *launch_screen.ChangeImagePropertyRequest) (resp *launch_screen.ChangeImagePropertyResponse, err error) {
	// TODO: Your code here...
	return
}

func (s *LaunchScreenServiceImpl) ChangeImage(stream launch_screen.LaunchScreenService_ChangeImageServer) (err error) {
	println("ChangeImage called")
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
