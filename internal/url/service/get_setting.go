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
	"regexp"
	"strings"
	"time"

	"github.com/west2-online/fzuhelper-server/kitex_gen/url"
	"github.com/west2-online/fzuhelper-server/pkg/constants"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
)

// GetCloudSetting todo:rewrite it
func (s *UrlService) GetCloudSetting(req *url.GetSettingRequest) (string, error) {
	data, err := utils.GetJSON(constants.StatisticPath + visitsFileName)
	if err != nil {
		return "", fmt.Errorf("UrlService.GetCloudSetting error:%w", err)
	}
	visitsDict := make(map[string]int64)
	err = json.Unmarshal(data, &visitsDict)
	if err != nil {
		return "", fmt.Errorf("UrlService.GetCloudSetting error:%w", err)
	}

	date := time.Now().Format("2006-01-02") // 获取当前日期
	if count, exists := visitsDict[date]; exists {
		visitsDict[date] = count + 1
	} else {
		visitsDict[date] = 1
	}
	saveData, err := json.Marshal(&visitsDict)
	if err != nil {
		return "", fmt.Errorf("UrlService.GetCloudSetting error:%w", err)
	}
	err = utils.SaveJSON(constants.StatisticPath+visitsFileName, saveData)
	if err != nil {
		return "", fmt.Errorf("UrlService.GetCloudSetting error:%w", err)
	}

	criteria := map[string]string{
		"account":   *req.Account,
		"version":   *req.Version,
		"beta":      *req.Beta,
		"phone":     *req.Phone,
		"isLogin":   *req.IsLogin,
		"loginType": *req.LoginType,
	}

	cloudSettingNoComments, err := loadCloudWithoutComments()
	if err != nil {
		return "", fmt.Errorf("UrlService.GetCloudSetting error:%w", err)
	}
	matchingPlan, err := findMatchingPlan(cloudSettingNoComments, criteria)
	if err != nil {
		return "", fmt.Errorf("UrlService.GetCloudSetting error:%w", err)
	}
	setting, ok := matchingPlan.(string)
	if !ok {
		return "", fmt.Errorf("UrlService.GetCloudSetting error:%w", err)
	}
	return setting, nil
}

// 匹配单个计划
func matches(plan map[string]interface{}, key, value string) bool {
	// 如果 key 不在计划中，则匹配通过
	if _, exists := plan[key]; !exists {
		return true
	}

	// 使用正则表达式匹配值
	pattern, ok := plan[key].(string)
	if !ok {
		return false
	}
	matched, _ := regexp.MatchString(pattern, value)
	return matched
}

// 查找匹配的计划
func findMatchingPlan(plans interface{}, criteria map[string]string) (interface{}, error) {
	plansList, ok := plans.([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid plans format")
	}

	var matchingPlan interface{}
	for _, plan := range plansList {
		planMap, ok := plan.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid plan format")
		}

		// 检查所有 criteria 是否匹配
		allMatch := true
		for key, value := range criteria {
			if !matches(planMap, key, value) {
				allMatch = false
				break
			}
		}

		if allMatch {
			matchingPlan = planMap["plan"]
		}
	}
	return matchingPlan, nil
}

func loadCloudWithoutComments() (map[string]string, error) {
	jsonStr, err := getAllCloudSettingStrFromFile()
	if err != nil {
		return nil, err
	}

	cloudSettingNoComments, err := getJSONWithoutComments(jsonStr)
	if err != nil {
		return nil, err
	}
	return cloudSettingNoComments, nil
}

func getAllCloudSettingStrFromFile() (string, error) {
	settingJson, err := utils.GetJSON(constants.StatisticPath + cloudSettingFileName)
	if err != nil {
		return "", err
	}
	return string(settingJson), nil
}

func getJSONWithoutComments(input string) (map[string]string, error) {
	lines := strings.Split(input, "\n")
	var cleanLines []string

	for _, line := range lines {
		cleanLine := removeComments(line)
		if cleanLine != "" {
			cleanLines = append(cleanLines, cleanLine)
		}
	}

	jsonStr := strings.Join(cleanLines, "\n")
	var result map[string]string

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return result, nil
}

func removeComments(line string) string {
	inString := false
	stringChar := ""

	for i := 0; i < len(line); i++ {
		// 处理字符串边界
		if line[i] == '"' || line[i] == '\'' {
			if !inString {
				inString = true
				stringChar = string(line[i])
			} else if stringChar == string(line[i]) {
				inString = false
			}
		}
		// 检查是否是注释
		if i+1 < len(line) && line[i:i+2] == "//" && !inString {
			return line[:i]
		}
	}

	return line
}
