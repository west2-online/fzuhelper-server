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
	"strings"

	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func (s *OAService) GetFeedbackByID(id int64, stuID string) (*model.Feedback, error) {
	if id <= 0 {
		return nil, fmt.Errorf("service.GetFeedbackByID: invalid id: %d", id)
	}
	stuID = strings.TrimSpace(stuID)
	if stuID == "" {
		return nil, fmt.Errorf("service.GetFeedbackByID: missing stu_id")
	}

	ok, fb, err := s.db.OA.GetFeedbackById(s.ctx, id, stuID)
	if err != nil {
		return nil, fmt.Errorf("service.GetFeedbackByID: get feedback failed: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("service.GetFeedbackByID: feedback not found")
	}

	fb.Screenshots = utils.EnsureJSONArray(fb.Screenshots)
	fb.VersionHistory = utils.EnsureJSONArray(fb.VersionHistory)
	fb.NetworkTraces = utils.EnsureJSON(fb.NetworkTraces)
	fb.Events = utils.EnsureJSONArray(fb.Events)
	fb.UserSettings = utils.EnsureJSONObject(fb.UserSettings)

	switch fb.NetworkEnv {
	case model.Network2G, model.Network3G, model.Network4G,
		model.Network5G, model.NetworkWifi, model.NetworkUnknown:
	default:
		logger.WithCtx(s.ctx).Warnf("feedback has invalid stored NetworkEnv, coercing to %q (report_id=%d, original=%q)",
			model.NetworkUnknown, fb.ReportId, fb.NetworkEnv)
		fb.NetworkEnv = model.NetworkUnknown
	}

	return fb, nil
}
