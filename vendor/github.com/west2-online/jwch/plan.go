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

package jwch

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/antchfx/htmlquery"

	"github.com/west2-online/jwch/constants"
)

func (s *Student) GetCultivatePlan() (string, error) {
	info, err := s.GetInfo()
	if err != nil {
		return "", err
	}

	// 获取初始页面状态
	viewStateMap, err := s.getState(constants.CultivatePlanURL)
	if err != nil {
		return "", err
	}

	// 尝试精确匹配学院和专业代码
	url, err := s.getCultivatePlanWithPreciseMatch(info, viewStateMap)
	if err == nil {
		return url, nil
	}

	// 如果精确匹配失败，使用fallback逻辑
	return s.getCultivatePlanWithFallback(info, viewStateMap)
}

// 精确匹配学院和专业代码获取培养方案
func (s *Student) getCultivatePlanWithPreciseMatch(info *StudentDetail, viewStateMap map[string]string) (string, error) {
	// 获取学院选择页面
	initialDoc, err := s.GetWithIdentifier(constants.CultivatePlanURL)
	if err != nil {
		return "", err
	}

	// 查找学院代码
	collegeSelect := htmlquery.FindOne(initialDoc, `//select[@id="xymcdpl"]`)
	if collegeSelect == nil {
		return "", fmt.Errorf("college select not found")
	}

	collegeCode := ""
	collegeOptions := htmlquery.Find(collegeSelect, ".//option")
	for _, option := range collegeOptions {
		optionText := htmlquery.InnerText(option)
		optionValue := htmlquery.SelectAttr(option, "value")

		// 直接匹配学院名称
		if optionText == info.College {
			collegeCode = optionValue
			break
		}

		// 处理学院改名的情况
		if (strings.Contains(optionText, "计算机与大数据") || strings.Contains(optionText, "数学与统计")) &&
			strings.Contains(info.College, "数学与计算机") {
			collegeCode = optionValue
			break
		}
	}

	if collegeCode == "" {
		return "", fmt.Errorf("college code not found for %s", info.College)
	}

	viewStateGenerator := htmlquery.SelectAttr(htmlquery.FindOne(initialDoc, `//*[@id="__VIEWSTATEGENERATOR"]`), "value")

	// 选择年级和学院后获取专业列表
	majorListResp, err := s.PostWithIdentifier(constants.CultivatePlanURL, map[string]string{
		"__VIEWSTATE":                         viewStateMap["VIEWSTATE"],
		"__EVENTVALIDATION":                   viewStateMap["EVENTVALIDATION"],
		"__EVENTTARGET":                       "ctl00$njdpl",
		"__EVENTARGUMENT":                     "",
		"__VIEWSTATEGENERATOR":                viewStateGenerator,
		"ctl00$njdpl":                         info.Grade,  // 年级
		"ctl00$xymcdpl":                       collegeCode, // 学院名称
		"ctl00$dldpl":                         "<-全部->",    // 大类
		"ctl00$zymcdpl":                       "<-全部->",    // 专业代码
		"ctl00$zylbdpl":                       "本专业",       // 修读类别：本专业/辅修
		"ctl00$ContentPlaceHolder1$DDL_syxw":  "<-全部->",    // 授予学位
		"ctl00$ContentPlaceHolder1$BT_submit": "确定",
	})
	if err != nil {
		return "", err
	}

	// 查找专业代码
	majorSelect := htmlquery.FindOne(majorListResp, `//select[@id="zymcdpl"]`)
	if majorSelect == nil {
		return "", fmt.Errorf("major select not found")
	}

	majorCode := ""
	majorOptions := htmlquery.Find(majorSelect, ".//option")
	for _, option := range majorOptions {
		optionText := htmlquery.InnerText(option)
		if optionText == info.Major {
			majorCode = htmlquery.SelectAttr(option, "value")
			break
		}
	}

	if majorCode == "" {
		return "", fmt.Errorf("major code not found for %s", info.Major)
	}

	// 构造最终URL
	finalURL := fmt.Sprintf("/pyfa/pyjh/pyfa_bzy.aspx?nj=%s&xyh=%s&zyh=%s&zylb=本专业&id=%s",
		info.Grade, collegeCode, majorCode, s.Identifier)

	return constants.JwchPrefix + finalURL, nil
}

// fallback逻辑：当精确匹配失败时使用
func (s *Student) getCultivatePlanWithFallback(info *StudentDetail, viewStateMap map[string]string) (string, error) {
	// 获取初始页面状态
	initialDoc, err := s.GetWithIdentifier(constants.CultivatePlanURL)
	if err != nil {
		return "", err
	}

	viewStateGenerator := htmlquery.SelectAttr(htmlquery.FindOne(initialDoc, `//*[@id="__VIEWSTATEGENERATOR"]`), "value")

	// 只选择年级，提交查询
	res, err := s.PostWithIdentifier(constants.CultivatePlanURL,
		map[string]string{
			"__VIEWSTATE":                         viewStateMap["VIEWSTATE"],
			"__EVENTVALIDATION":                   viewStateMap["EVENTVALIDATION"],
			"__EVENTTARGET":                       "",
			"__EVENTARGUMENT":                     "",
			"__VIEWSTATEGENERATOR":                viewStateGenerator,
			"ctl00$njdpl":                         info.Grade,
			"ctl00$dldpl":                         "<-全部->",
			"ctl00$zymcdpl":                       "<-全部->",
			"ctl00$zylbdpl":                       "本专业",
			"ctl00$ContentPlaceHolder1$DDL_syxw":  "<-全部->",
			"ctl00$ContentPlaceHolder1$BT_submit": "确定",
		})
	if err != nil {
		return "", err
	}
	xpathExpr := fmt.Sprintf("//tr[td[matches(string(.), '^（.*?）%s$')]]/td/a[contains(@href, 'pyfa')]/@href", regexp.QuoteMeta(info.Major))
	node := htmlquery.FindOne(res, xpathExpr)
	if node == nil {
		return "", fmt.Errorf("cultivate plan not found for major: %s", info.Major)
	}

	url := htmlquery.SelectAttr(node, "href")
	formatUrl := constants.JwchPrefix + "/pyfa/pyjh/" + strings.TrimPrefix(strings.TrimSuffix(url, "')"), "javascript:pop1('")
	return formatUrl, nil
}
