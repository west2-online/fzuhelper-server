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

	"github.com/west2-online/fzuhelper-server/internal/version/pack"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
)

func (s *VersionService) GetReleaseVersion() (*pack.Version, error) {
	return s.getLatestVersion(apkTypeRelease)
}

func (s *VersionService) GetBetaVersion() (*pack.Version, error) {
	return s.getLatestVersion(apkTypeBeta)
}

// getLatestVersion retrieves the latest version of the given type.
// It checks Redis cache first; on cache miss, queries the database and populates the cache.
func (s *VersionService) getLatestVersion(versionType string) (*pack.Version, error) {
	// Try cache first
	vh, err := s.cache.Version.GetLatestVersionCache(s.ctx, versionType)
	if err != nil {
		// Cache miss or error — fall through to DB
		vh, err = s.db.Version.GetLatestVersionByType(s.ctx, versionType)
		if err != nil {
			return nil, fmt.Errorf("getLatestVersion: db query error: %w", err)
		}
		if vh == nil {
			return nil, fmt.Errorf("getLatestVersion: no %s version found in database", versionType)
		}
		// Populate cache (ignore cache write errors — non-critical path)
		_ = s.cache.Version.SetLatestVersionCache(s.ctx, vh)
	}

	return dbModelToPackVersion(vh), nil
}

// dbModelToPackVersion converts a DB model VersionHistory to a pack.Version.
func dbModelToPackVersion(vh *model.VersionHistory) *pack.Version {
	return &pack.Version{
		Version: vh.Version,
		Code:    vh.Code,
		Url:     vh.Url,
		Feature: vh.Feature,
		Force:   vh.Force,
	}
}
