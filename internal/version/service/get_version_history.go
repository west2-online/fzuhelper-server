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

	kitexModel "github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/fzuhelper-server/kitex_gen/version"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	dbModel "github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// GetVersionHistoryList retrieves version history records from the database,
// converts them from DB model to RPC type, and returns them with the next page token.
// Returns an empty slice (not nil) when no versions have been uploaded.
func (s *VersionService) GetVersionHistoryList(req *version.GetVersionHistoryListRequest) ([]*kitexModel.VersionHistory, int64, error) {
	if req == nil || !utils.CheckPwd(req.Password) {
		return nil, 0, buildAuthFailedError()
	}

	limit := req.GetLimit()
	if limit <= 0 || limit > constants.VersionHistoryMaxPageSize {
		limit = constants.VersionHistoryDefaultPageSize
	}

	records, nextPageToken, err := s.db.Version.GetVersionHistoryList(s.ctx, int(limit), req.GetPageToken())
	if err != nil {
		return nil, 0, fmt.Errorf("GetVersionHistoryList: get version history list error: %w", err)
	}

	result := make([]*kitexModel.VersionHistory, 0, len(records))
	for _, r := range records {
		result = append(result, buildVersionHistory(r))
	}
	return result, nextPageToken, nil
}

// buildVersionHistory converts a DB model VersionHistory to a kitex RPC VersionHistory.
// The CreatedAt time is formatted as "2006-01-02 15:04:05" for consistent JSON output.
func buildVersionHistory(r *dbModel.VersionHistory) *kitexModel.VersionHistory {
	return &kitexModel.VersionHistory{
		Id:        r.Id,
		Version:   r.Version,
		Code:      r.Code,
		Url:       r.Url,
		Feature:   r.Feature,
		Force:     r.Force,
		Type:      r.Type,
		CreatedAt: r.CreatedAt.Format("2006-01-02 15:04:05"),
	}
}
