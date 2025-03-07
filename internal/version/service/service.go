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
	"context"
	"net/http"

	"github.com/west2-online/fzuhelper-server/pkg/base"
	"github.com/west2-online/fzuhelper-server/pkg/cache"
	"github.com/west2-online/fzuhelper-server/pkg/db"
	"github.com/west2-online/fzuhelper-server/pkg/errno"
)

const (
	cloudSettingFileName   = "cloud_setting.json"
	visitsFileName         = "visits.json"
	releaseVersionFileName = "version.json"
	betaVersionFileName    = "versionbeta.json"

	apkTypeRelease = "release"
	apkTypeBeta    = "beta"

	authFailedErrorMsg = "authorization failed"
)

type VersionService struct {
	ctx   context.Context
	db    *db.Database
	cache *cache.Cache
}

func NewVersionService(ctx context.Context, clientset *base.ClientSet) *VersionService {
	return &VersionService{
		ctx:   ctx,
		db:    clientset.DBClient,
		cache: clientset.CacheClient,
	}
}

// buildAuthFailedError customize the client
func buildAuthFailedError() errno.ErrNo {
	return errno.NewErrNo(http.StatusUnauthorized, authFailedErrorMsg)
}
