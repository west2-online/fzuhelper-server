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
	"strings"

	"github.com/antchfx/htmlquery"

	"github.com/west2-online/jwch/constants"
)

func (s *Student) GetCredit() (creditStatistics []*CreditStatistics, err error) {
	resp, err := s.GetWithIdentifier(constants.CreditQueryURL)
	if err != nil {
		return nil, err
	}

	spanNode := htmlquery.FindOne(resp, `//*[@id="ContentPlaceHolder1_LB_kb"]`)
	if spanNode == nil {
		return nil, fmt.Errorf("failed to find the statistics span element")
	}

	tables := htmlquery.Find(spanNode, "//table")
	if len(tables) == 0 {
		return nil, fmt.Errorf("failed to find tables within the span element")
	}
	tables = tables[:len(tables)-1] // 去掉最后一个表格

	creditStatistics = make([]*CreditStatistics, 0)

	for _, table := range tables {
		rows := htmlquery.Find(table, "//tr")

		// 临时存储三列数据
		temp := [3][]string{
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
		}
		// 遍历每行，提取单元格数据
		for index, row := range rows {
			cells := htmlquery.Find(row, "./td")
			// 提取单元格数据
			for _, cell := range cells {
				text := htmlquery.InnerText(cell)
				if text != "查" { // 因为表格的第三行多了一个单元格，去掉一个无用的格子它使得表格规整
					temp[index] = append(temp[index], text)
				}
			}
		}
		// 构建 CreditStatistics
		for i := 0; i < len(temp[0]); i++ {
			// 去掉个人信息的列（这列第一个单元格式空的）和“修习情况”这个无效的列
			if strings.TrimSpace(temp[0][i]) != "" && !strings.Contains(temp[0][i], "情况") {
				bean := &CreditStatistics{
					Type:  temp[0][i],
					Gain:  temp[2][i],
					Total: temp[1][i],
				}
				creditStatistics = append(creditStatistics, bean)
			}
		}
	}

	return creditStatistics, nil
}

func (s *Student) GetGPA() (gpa *GPABean, err error) {
	gpa = &GPABean{}
	resp, err := s.GetWithIdentifier(constants.GPAQueryURL)
	if err != nil {
		return gpa, err
	}

	document := htmlquery.FindOne(resp, `//*[@id="ContentPlaceHolder1_Label1"]`)
	if document == nil {
		return gpa, fmt.Errorf("failed to find the time element")
	}

	timeText := htmlquery.InnerText(document)
	gpa.Time = strings.TrimSpace(timeText)

	table := htmlquery.FindOne(resp, `//*[@id="ContentPlaceHolder1_DataList_xxk"]`)
	if table == nil {
		return gpa, fmt.Errorf("failed to find the GPA table")
	}

	// 获取表头标题
	titleRow := htmlquery.FindOne(table, `//tr[@style="height:30px; background:#efefef; border-bottom:1px solid gray; border-left:1px solid gray; vertical-align:middle;"]`)
	if titleRow == nil {
		return gpa, fmt.Errorf("failed to find the title row in GPA table")
	}

	// 获取每个表头标题的单元格
	tdsTitle := htmlquery.Find(titleRow, `./td[@align="center"]`)
	width := len(tdsTitle)

	// 获取表格中的所有数据
	tdsFull := htmlquery.Find(table, `.//td[@align="center"]`)
	if len(tdsFull) == 0 {
		return gpa, fmt.Errorf("failed to find GPA data cells")
	}

	height := len(tdsFull)/width - 1

	var data []GPAData
	for h := 1; h <= height; h++ {
		for w := 0; w < width; w++ {
			data = append(data, GPAData{
				Type:  htmlquery.InnerText(tdsTitle[w]),
				Value: htmlquery.InnerText(tdsFull[width*h+w]),
			})
		}
	}
	gpa.Data = data

	return gpa, nil
}

// GetCreditV2 用于获取原始的学分统计
func (s *Student) GetCreditV2() (majorCredits, minorCredits []*CreditStatistics, err error) {
	resp, err := s.GetWithIdentifier(constants.CreditQueryURL)
	if err != nil {
		return nil, nil, err
	}

	spanNode := htmlquery.FindOne(resp, `//*[@id="ContentPlaceHolder1_LB_kb"]`)
	if spanNode == nil {
		return nil, nil, fmt.Errorf("failed to find the statistics span element")
	}

	tables := htmlquery.Find(spanNode, "//table")
	if len(tables) == 0 {
		return nil, nil, fmt.Errorf("failed to find tables within the span element")
	}
	tables = tables[:len(tables)-1] // 去掉最后一个表格

	// 处理主修专业和辅修专业
	for tableIndex, table := range tables {
		rows := htmlquery.Find(table, "//tr")

		// 临时存储三列数据
		temp := [3][]string{
			make([]string, 0),
			make([]string, 0),
			make([]string, 0),
		}
		// 遍历每行，提取单元格数据
		for index, row := range rows {
			cells := htmlquery.Find(row, "./td")
			// 提取单元格数据
			for _, cell := range cells {
				text := htmlquery.InnerText(cell)
				if text != "查" { // 因为表格的第三行多了一个单元格，去掉一个无用的格子它使得表格规整
					temp[index] = append(temp[index], text)
				}
			}
		}
		// 构建 CreditStatistics
		for i := 0; i < len(temp[0]); i++ {
			// 去掉个人信息的列（这列第一个单元格式空的）和"修习情况"这个无效的列
			if strings.TrimSpace(temp[0][i]) != "" && !strings.Contains(temp[0][i], "情况") {
				bean := &CreditStatistics{
					Type:  temp[0][i],
					Gain:  temp[2][i],
					Total: temp[1][i],
				}
				// 第一个表格是主修专业，第二个表格是辅修专业
				if tableIndex == 0 {
					majorCredits = append(majorCredits, bean)
				} else {
					minorCredits = append(minorCredits, bean)
				}
			}
		}
	}

	return majorCredits, minorCredits, err
}
