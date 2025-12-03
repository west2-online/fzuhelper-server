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
	"github.com/west2-online/fzuhelper-server/kitex_gen/model"
	"github.com/west2-online/jwch"
)

func BuildCredit(data []*jwch.CreditStatistics) []*model.Credit {
	credit := make([]*model.Credit, len(data))
	for i := 0; i < len(data); i++ {
		credit[i] = &model.Credit{Type: data[i].Type, Gain: data[i].Gain, Total: data[i].Total}
	}

	return credit
}

func BuildCreditResponse(data *CreditResponse) model.CreditResponse {
	if data == nil {
		return nil
	}

	// Convert []*jwch.CreditCategory to []*model.CreditCategory
	creditResponse := make([]*model.CreditCategory, len(*data))
	for i, category := range *data {
		// Convert each category
		creditCategory := &model.CreditCategory{
			Type: category.Type,
			Data: make([]*model.CreditDetail, len(category.Data)),
		}

		// Convert each detail in the category
		for j, detail := range category.Data {
			creditCategory.Data[j] = &model.CreditDetail{
				Key:   detail.Key,
				Value: detail.Value,
			}
		}

		creditResponse[i] = creditCategory
	}

	return creditResponse
}
