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
	"strings"
	"time"

	"github.com/west2-online/jwch/constants"
	"github.com/west2-online/jwch/utils"

	"github.com/antchfx/htmlquery"
)

// GetLectures 获取报名的讲座
func (s *Student) GetLectures() ([]*Lecture, error) {
	resp, err := s.GetWithIdentifier(constants.LectureURL)
	if err != nil {
		return nil, err
	}

	list := htmlquery.Find(htmlquery.FindOne(resp, `//*[@id="ContentPlaceHolder1_DataList_xxk"]/tbody`), "tr")

	// 前三个元素分别为 空行 标题 表头
	list = list[3:]

	res := make([]*Lecture, 0)

	for _, row := range list {
		// 跳过空行
		if strings.TrimSpace(htmlquery.SelectAttr(row, "style")) == "" {
			continue
		}

		cells := htmlquery.Find(row, "./td")

		res = append(res, &Lecture{
			Category:         htmlquery.InnerText(cells[0]),
			IssueNumber:      utils.SafeAtoi(htmlquery.InnerText(cells[1])),
			Title:            htmlquery.InnerText(cells[2]),
			Speaker:          htmlquery.InnerText(cells[3]),
			Timestamp:        parseDateTime(htmlquery.InnerText(cells[4])),
			Location:         htmlquery.InnerText(cells[5]),
			AttendanceStatus: htmlquery.InnerText(cells[6]),
		})
	}

	return res, nil
}

func parseDateTime(dateTime string) int64 {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	t, err := time.ParseInLocation("2006-01-02\u00A0\u00A015：04", dateTime, loc)
	if err != nil {
		return 0
	}
	return t.UnixMilli()
}
