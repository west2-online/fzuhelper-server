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
	"net/mail"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/db/model"
	"github.com/west2-online/fzuhelper-server/pkg/logger"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

func normalizeFeedbackScreenshots(value string) (string, error) {
	if strings.TrimSpace(value) == "" {
		return "[]", nil
	}

	var screenshots []string
	if err := json.Unmarshal([]byte(value), &screenshots); err != nil {
		return "", fmt.Errorf("screenshots must be a JSON string array: %w", err)
	}
	if len(screenshots) > constants.FeedbackScreenshotMaxCount {
		return "", fmt.Errorf("screenshots cannot contain more than 9 items")
	}
	if screenshots == nil {
		screenshots = []string{}
	}
	for _, screenshot := range screenshots {
		if len(screenshot) > constants.FeedbackScreenshotURLMax {
			return "", fmt.Errorf("screenshot URL cannot exceed %d bytes", constants.FeedbackScreenshotURLMax)
		}
		parsed, err := url.ParseRequestURI(screenshot)
		if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
			return "", fmt.Errorf("screenshot URL must be an absolute HTTP(S) URL")
		}
	}

	normalized, err := json.Marshal(screenshots)
	if err != nil {
		return "", fmt.Errorf("normalize screenshots failed: %w", err)
	}
	return string(normalized), nil
}

func normalizeFeedbackJSONObject(value, field string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "{}", nil
	}
	if len(value) > constants.FeedbackJSONMaxSize {
		return "", fmt.Errorf("%s cannot exceed %d bytes", field, constants.FeedbackJSONMaxSize)
	}

	var object map[string]json.RawMessage
	if err := json.Unmarshal([]byte(value), &object); err != nil || object == nil {
		return "", fmt.Errorf("%s must be a JSON object", field)
	}
	normalized, err := json.Marshal(object)
	if err != nil {
		return "", fmt.Errorf("normalize %s failed: %w", field, err)
	}
	return string(normalized), nil
}

func normalizeAndValidateFeedbackFields(req *CreateFeedbackReq) error {
	fields := []struct {
		name     string
		value    *string
		max      int
		optional bool
	}{
		{name: "stu_id", value: &req.StuId, max: constants.FeedbackStuIDMax},
		{name: "name", value: &req.Name, max: constants.FeedbackNameMax},
		{name: "college", value: &req.College, max: constants.FeedbackCollegeMax},
		{name: "contact_phone", value: &req.ContactPhone, max: constants.FeedbackContactMax, optional: true},
		{name: "contact_qq", value: &req.ContactQQ, max: constants.FeedbackContactMax, optional: true},
		{name: "contact_email", value: &req.ContactEmail, max: constants.FeedbackEmailMax, optional: true},
		{name: "network_env", value: &req.NetworkEnv, max: constants.FeedbackNetworkEnvMax},
		{name: "os_name", value: &req.OsName, max: constants.FeedbackOSNameMax},
		{name: "os_version", value: &req.OsVersion, max: constants.FeedbackOSVersionMax},
		{name: "manufacturer", value: &req.Manufacturer, max: constants.FeedbackDeviceInfoMax},
		{name: "device_model", value: &req.DeviceModel, max: constants.FeedbackDeviceInfoMax},
		{name: "problem_desc", value: &req.ProblemDesc, max: constants.FeedbackProblemDescMax},
		{name: "app_version", value: &req.AppVersion, max: constants.FeedbackAppVersionMax},
	}

	for _, field := range fields {
		*field.value = strings.TrimSpace(*field.value)
		if *field.value == "" && !field.optional {
			return fmt.Errorf("service.CreateFeedback: missing required field %s", field.name)
		}
		if utf8.RuneCountInString(*field.value) > field.max {
			return fmt.Errorf("service.CreateFeedback: field %s cannot exceed %d characters", field.name, field.max)
		}
	}
	if req.ContactPhone == "" && req.ContactQQ == "" && req.ContactEmail == "" {
		return fmt.Errorf("service.CreateFeedback: at least one contact method is required")
	}
	if req.ContactPhone != "" && !regexp.MustCompile(`^[0-9]{11}$`).MatchString(req.ContactPhone) {
		return fmt.Errorf("service.CreateFeedback: contact_phone must contain exactly 11 digits")
	}
	if req.ContactEmail != "" {
		address, err := mail.ParseAddress(req.ContactEmail)
		if err != nil || address.Address != req.ContactEmail {
			return fmt.Errorf("service.CreateFeedback: contact_email must be a valid email address")
		}
	}
	return nil
}

func (s *OAService) CreateFeedback(req *CreateFeedbackReq) (int64, error) {
	if req == nil {
		return 0, fmt.Errorf("service.CreateFeedback: request is nil")
	}
	if err := normalizeAndValidateFeedbackFields(req); err != nil {
		return 0, err
	}

	screenshots, err := normalizeFeedbackScreenshots(req.Screenshots)
	if err != nil {
		return 0, err
	}
	req.Screenshots = screenshots
	if len(req.VersionHistory) > constants.FeedbackJSONMaxSize ||
		len(req.Events) > constants.FeedbackJSONMaxSize || len(req.UserSettings) > constants.FeedbackJSONMaxSize {
		return 0, fmt.Errorf("service.CreateFeedback: JSON field cannot exceed %d bytes", constants.FeedbackJSONMaxSize)
	}
	req.VersionHistory = utils.EnsureJSONArray(req.VersionHistory)
	req.NetworkTraces, err = normalizeFeedbackJSONObject(req.NetworkTraces, "network_traces")
	if err != nil {
		return 0, err
	}
	req.Events = utils.EnsureJSONArray(req.Events)
	req.UserSettings = utils.EnsureJSONObject(req.UserSettings)

	switch req.NetworkEnv {
	case string(model.Network2G), string(model.Network3G), string(model.Network4G),
		string(model.Network5G), string(model.NetworkWifi), string(model.NetworkUnknown):
	default:
		logger.WithCtx(s.ctx).Warnf("invalid NetworkEnv=%q, fallback=%q",
			req.NetworkEnv, model.NetworkUnknown)
		req.NetworkEnv = string(model.NetworkUnknown)
	}

	reportID, err := s.sf.NextVal()
	if err != nil {
		return 0, fmt.Errorf("service.CreateFeedback: generate report_id failed: %w", err)
	}

	fb := &model.Feedback{
		ReportId:       reportID,
		StuId:          req.StuId,
		Name:           req.Name,
		College:        req.College,
		ContactPhone:   req.ContactPhone,
		ContactQQ:      req.ContactQQ,
		ContactEmail:   req.ContactEmail,
		NetworkEnv:     model.NetworkEnv(req.NetworkEnv),
		IsOnCampus:     req.IsOnCampus,
		OsName:         req.OsName,
		OsVersion:      req.OsVersion,
		Manufacturer:   req.Manufacturer,
		DeviceModel:    req.DeviceModel,
		ProblemDesc:    req.ProblemDesc,
		Screenshots:    req.Screenshots,
		AppVersion:     req.AppVersion,
		VersionHistory: req.VersionHistory,
		NetworkTraces:  req.NetworkTraces,
		Events:         req.Events,
		UserSettings:   req.UserSettings,
	}

	if err := s.db.OA.CreateFeedback(s.ctx, fb); err != nil {
		logger.WithCtx(s.ctx).Errorf("service.CreateFeedback dal error: %v", err)
		return 0, fmt.Errorf("service.CreateFeedback: create feedback failed: %w", err)
	}

	return reportID, nil
}
