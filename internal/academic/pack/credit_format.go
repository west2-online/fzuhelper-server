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
	"fmt"
	"strconv"
	"strings"

	"github.com/west2-online/jwch"
)

// CreditDetail 用于表示学分的详细数据项
type CreditDetail struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// CreditCategory 用于表示学分分类
type CreditCategory struct {
	Type string          `json:"type"`
	Data []*CreditDetail `json:"data"`
}

// CreditResponse 用于学分响应
type CreditResponse []*CreditCategory

// BuildCreditCategory 格式化处理学分统计
func BuildCreditCategory(categoryType string, credits []*jwch.CreditStatistics) *CreditCategory {
	category := &CreditCategory{
		Type: categoryType,
		Data: make([]*CreditDetail, 0, len(credits)),
	}

	// 用于计算"总计"栏的还需学分总和
	var totalNeed float64
	specialKeys := []string{"奖励", "其它", "重修", "正在修习", "CET"}

	for _, credit := range credits {
		isSpecial := false
		for _, key := range specialKeys {
			if strings.Contains(credit.Type, key) {
				isSpecial = true
				break
			}
		}

		creditType := credit.Type
		switch {
		case creditType == "总计":
			// 处理"总计"行
			gain, _ := strconv.ParseFloat(credit.Gain, 64)
			total, _ := strconv.ParseFloat(credit.Total, 64)
			value := fmt.Sprintf("%g / %g (还需 %g 分)", gain, total, totalNeed)

			category.Data = append(category.Data, &CreditDetail{
				Key:   credit.Type,
				Value: value,
			})
		case isSpecial:
			category.Data = append(category.Data, &CreditDetail{
				Key:   credit.Type,
				Value: credit.Gain,
			})
		default:
			gain, _ := strconv.ParseFloat(credit.Gain, 64)
			total, _ := strconv.ParseFloat(credit.Total, 64)
			var value string

			if gain >= total {
				value = fmt.Sprintf("%g / %g", gain, total)
			} else {
				need := total - gain
				value = fmt.Sprintf("%g / %g (还需 %g 分)", gain, total, need)
				// 未修满时需要的学分计入总计
				totalNeed += need
			}

			category.Data = append(category.Data, &CreditDetail{
				Key:   credit.Type,
				Value: value,
			})
		}
	}
	return category
}
