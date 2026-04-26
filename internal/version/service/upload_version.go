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
	"encoding/json"
	"fmt"

	"github.com/west2-online/fzuhelper-server/internal/version/pack"
	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
	"github.com/west2-online/fzuhelper-server/pkg/upyun"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *VersionService) UploadVersion(req *version.UploadRequest) error {
	if !utils.CheckPwd(req.Password) {
		return buildAuthFailedError()
	}
	v := &pack.Version{
		Version: req.Version,
		Code:    req.Code,
		Url:     req.Url,
		Feature: req.Feature,
		Force:   req.Force,
	}
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("VersionService.UploadVersion json marshal err: %w", err)
	}

	switch req.Type {
	case apkTypeRelease:
		err = upyun.URlUploadFile(jsonBytes, upyun.JoinFileName(releaseVersionFileName))
		if err != nil {
			return fmt.Errorf("VersionService.UploadVersion upload release err: %w", err)
		}
	case apkTypeBeta:
		err = upyun.URlUploadFile(jsonBytes, upyun.JoinFileName(betaVersionFileName))
		if err != nil {
			return fmt.Errorf("VersionService.UploadVersion upload beta err: %w", err)
		}
	default:
		return errno.ParamError
	}

	// Persist to version_history table for the admin list endpoint.
	// If DB write fails, the error message indicates that the UPYUN upload already succeeded.
	vh := &model.VersionHistory{
		Version: req.Version,
		Code:    req.Code,
		Url:     req.Url,
		Feature: req.Feature,
		Force:   req.Force,
		Type:    req.Type,
	}
	if dbErr := s.db.Version.CreateVersionHistory(s.ctx, vh); dbErr != nil {
		return fmt.Errorf("VersionService.UploadVersion save history err: %w (upload succeeded)", dbErr)
	}

	// Refresh the latest version cache for this type so GetReleaseVersion/GetBetaVersion
	// return the newly uploaded version immediately. Cache write failure is non-critical.
	_ = s.cache.Version.SetLatestVersionCache(s.ctx, vh)

	return nil
}
