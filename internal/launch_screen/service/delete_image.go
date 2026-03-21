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

	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *LaunchScreenService) DeleteImage(req *launch_screen.DeleteImageRequest) error {
	if !utils.CheckPwd(req.Password) {
		return buildAuthFailedError()
	}
	pic, err := s.db.LaunchScreen.DeleteImage(s.ctx, req.PictureId)
	if err != nil {
		return fmt.Errorf("LaunchScreenService.DeleteImage error:%w", err)
	}
	remotePath := s.ossClient.GetRemotePathFromUrl(pic.Url)
	if err = s.ossClient.DeleteImg(remotePath); err != nil {
		return fmt.Errorf("LaunchScreen.DeleteImage error: %w", err)
	}
	return nil
}
