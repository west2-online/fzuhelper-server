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

package pack

import (
	apimodel "github.com/west2-online/fzuhelper-server/api/model/model"
	model "github.com/west2-online/fzuhelper-server/kitex_gen/model"
)

func BuildFeedbackList(dbItems []*model.FeedbackListItem) []*apimodel.FeedbackListItem {
	if len(dbItems) == 0 {
		return nil
	}
	out := make([]*apimodel.FeedbackListItem, len(dbItems))
	for i := range dbItems {
		it := dbItems[i]
		out[i] = &apimodel.FeedbackListItem{
			ReportID:    it.ReportId,
			Name:        it.Name,
			NetworkEnv:  it.NetworkEnv,
			ProblemDesc: it.ProblemDesc,
			AppVersion:  it.AppVersion,
		}
	}
	return out
}
