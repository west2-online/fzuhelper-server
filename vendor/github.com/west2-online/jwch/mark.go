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

	"golang.org/x/net/html"

	"github.com/west2-online/jwch/constants"
	"github.com/west2-online/jwch/errno"
	"github.com/west2-online/jwch/utils"

	"github.com/antchfx/htmlquery"
)

// 获取成绩，由于教务处缺陷，这里会返回全部的成绩
func (s *Student) GetMarks() (resp []*Mark, err error) {
	res, err := s.GetWithIdentifier(constants.MarksQueryURL)
	if err != nil {
		return nil, err
	}

	// window.alert( '你尚有学费未缴清，暂时不能查询成绩，如有疑问请与计财处联系！');
	htmlStr := htmlquery.OutputHTML(res, false)
	re := regexp.MustCompile(`window\.alert\s*\(\s*'([^']*)'\s*\)`)
	matches := re.FindStringSubmatch(htmlStr)
	if len(matches) == 2 {
		message := matches[1]
		return nil, errno.HTMLParseError.WithMessage(message)
	}

	table := htmlquery.FindOne(res, `//*[@id="ContentPlaceHolder1_DataList_xxk"]/tbody`)
	if table == nil {
		return nil, errno.HTMLParseError.WithMessage("marks table not found")
	}

	list := htmlquery.Find(table, "tr")
	if len(list) < 2 {
		return nil, errno.HTMLParseError.WithMessage("insufficient table rows")
	}

	// 去除第一个元素，第一个元素是标题栏，有个判断文本是“课程名称”
	// TODO: 我们如何确保第一个元素一定是标题栏?
	list = list[2:]

	resp = make([]*Mark, 0)

	for _, node := range list {
		// 教务处的表格HTML是不规范的，因此XPath解析会出现一些BUG
		if strings.TrimSpace(htmlquery.SelectAttr(node, "style")) == "" {
			continue
		}
		info := htmlquery.Find(node, `td`) // 一行的所有信息

		// 这个表格有12栏
		if len(info) < 12 {
			return nil, errno.HTMLParseError.WithMessage("get mark info failed")
		}

		// TODO: performance optimization
		resp = append(resp, &Mark{
			Type:          htmlquery.OutputHTML(info[0], false),
			Semester:      htmlquery.OutputHTML(info[1], false),
			Name:          htmlquery.OutputHTML(info[2], false),
			Credits:       safeExtractionFirst(info[3], "span"),
			Score:         safeExtractionFirst(info[4], "font"),
			GPA:           htmlquery.OutputHTML(info[5], false),
			EarnedCredits: htmlquery.OutputHTML(info[6], false),
			ElectiveType:  utils.GetChineseCharacter(htmlquery.OutputHTML(info[7], false)),
			ExamType:      utils.GetChineseCharacter(htmlquery.OutputHTML(info[8], false)),
			Teacher:       htmlquery.OutputHTML(info[9], false),
			Classroom:     strings.TrimSpace(htmlquery.InnerText(info[10])),
			ExamTime:      strings.TrimSpace(htmlquery.InnerText(info[11])),
		})
	}

	return resp, nil
}

// 获取CET成绩
func (s *Student) GetCET() ([]*UnifiedExam, error) {
	resp, err := s.GetWithIdentifier(constants.CETQueryURL)
	if err != nil {
		return nil, err
	}

	return s.parseUnifiedExam(resp)
}

// 获取省计算机成绩
func (s *Student) GetJS() ([]*UnifiedExam, error) {
	resp, err := s.GetWithIdentifier(constants.JSQueryURL)
	if err != nil {
		return nil, err
	}

	return s.parseUnifiedExam(resp)
}

// 解析统一考试成绩
func (s *Student) parseUnifiedExam(resp *html.Node) ([]*UnifiedExam, error) {
	var exams []*UnifiedExam

	// 查找包含成绩的表格
	table := htmlquery.FindOne(resp, `//*[@id="ContentPlaceHolder1_DataList_xxk"]`)
	if table == nil {
		return nil, fmt.Errorf("failed to find the exam table")
	}

	// 查找所有考试成绩行
	rows := htmlquery.Find(table, `.//tr[@onmouseover]`)
	if len(rows) == 0 {
		return nil, nil // 这里不返回错误，因为有可能没有考试成绩
	}

	// 遍历每一行，提取成绩信息
	for _, row := range rows {
		tds := htmlquery.Find(row, `.//td`)
		if len(tds) < 3 {
			continue // 如果某行的列数不满足要求则跳过
		}

		// 创建一个新的 UnifiedExam 对象
		exam := &UnifiedExam{
			Name:  htmlquery.InnerText(tds[0]),
			Score: htmlquery.InnerText(tds[2]),
			Term:  htmlquery.InnerText(tds[1]),
		}

		exams = append(exams, exam)
	}

	return exams, nil
}
