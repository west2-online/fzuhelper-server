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
	"regexp"
	"strings"

	"github.com/west2-online/jwch/utils"

	"github.com/antchfx/htmlquery"

	"github.com/west2-online/jwch/constants"
)

func (s *Student) GetSchoolCalendar() (*SchoolCalendar, error) {
	resp, err := s.GetWithIdentifier(constants.SchoolCalendarURL)
	if err != nil {
		return nil, err
	}

	rawCurTerm := htmlquery.InnerText(htmlquery.FindOne(resp, `//html/body/center/div`))
	rawCurTerm, err = utils.ConvertGB2312ToUTF8([]byte(rawCurTerm))
	if err != nil {
		return nil, err
	}
	curTermRegex := regexp.MustCompile(`当前学期：(\d{6})`)
	curTerm := curTermRegex.FindStringSubmatch(rawCurTerm)[1]

	res := &SchoolCalendar{
		CurrentTerm: curTerm,
	}

	list := htmlquery.Find(resp, `//select[@name="xq"]/option/@value`)

	for _, node := range list {
		// 需要取前16个年份
		if len(res.Terms) >= 16 {
			break
		}
		rawTerm := htmlquery.SelectAttr(node, "value")
		/*
			2024012024082620250117
			[0] 202401
			[1] 20240826
			[2] 20250117
		*/
		schoolYear := rawTerm[0:4]
		term := rawTerm[0:6]
		startDate := rawTerm[6:14]
		endDate := rawTerm[14:22]

		// convert 20240826 to 2024-08-26
		startDate = startDate[0:4] + "-" + startDate[4:6] + "-" + startDate[6:8]
		endDate = endDate[0:4] + "-" + endDate[4:6] + "-" + endDate[6:8]

		res.Terms = append(res.Terms, CalTerm{
			TermId:     rawTerm,
			SchoolYear: schoolYear,
			Term:       term,
			StartDate:  startDate,
			EndDate:    endDate,
		})
	}

	return res, nil
}

func (s *Student) GetTermEvents(termId string) (*CalTermEvents, error) {
	resp, err := s.PostWithIdentifier(constants.SchoolCalendarURL, map[string]string{
		"xq":     termId,
		"submit": "提交",
	})
	if err != nil {
		return nil, err
	}
	res := &CalTermEvents{
		TermId:     termId,
		Term:       termId[0:6],
		SchoolYear: termId[0:4],
	}
	// 远古校历没有任何内容
	table := htmlquery.FindOne(resp, `/html/body/table[2]/tbody/tr`)
	if table == nil {
		return res, nil
	}
	rawTermDetail := htmlquery.InnerText(table)
	rawTermDetail = strings.ReplaceAll(rawTermDetail, " ", " ")
	rawTermDetail, _ = utils.ConvertGB2312ToUTF8([]byte(rawTermDetail))

	termDetail := strings.Split(rawTermDetail, "；")

	for _, event := range termDetail {
		event = strings.TrimSpace(event)
		if event == "" {
			continue
		}

		rawData := strings.Split(event, "为")
		if len(rawData) < 2 {
			// 远古学期的数据格式可能不统一，不做处理
			res.Events = append(res.Events, CalTermEvent{Name: strings.TrimSpace(event)})
			continue
		}

		rawDate := strings.Split(strings.TrimSpace(rawData[0]), "至")
		name := strings.TrimSpace(strings.Join(rawData[1:], "为"))

		if len(rawDate) >= 2 {
			startDate := strings.TrimSpace(rawDate[0])
			endDate := strings.TrimSpace(rawDate[1])

			res.Events = append(res.Events, CalTermEvent{
				Name:      name,
				StartDate: startDate,
				EndDate:   endDate,
			})
		} else {
			// 兼容单日事件格式: 2025-09-07为学生注册
			date := strings.TrimSpace(rawDate[0])
			res.Events = append(res.Events, CalTermEvent{
				Name:      name,
				StartDate: date,
				EndDate:   date,
			})
		}
	}

	return res, nil
}
