package yjsy

import (
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/west2-online/yjsy/constants"
)

func (s *Student) GetMarks() (resp []*Mark, err error) {
	res, err := s.GetWithIdentifier(constants.MarksQueryURL, nil)
	if err != nil {
		return nil, err
	}

	resp = make([]*Mark, 0)
	rows := htmlquery.Find(res, "//table[@align='center'][2]//tr[position()>6]")
	for _, row := range rows {
		cells := htmlquery.Find(row, "./td")
		// 只有成绩栏有7列
		if len(cells) < 6 {
			continue
		}

		// 不存在的信息默认留空
		mark := &Mark{
			Name:     strings.TrimSpace(htmlquery.InnerText(cells[1])),
			Type:     strings.TrimSpace(htmlquery.InnerText(cells[2])),
			Semester: strings.TrimSpace(htmlquery.InnerText(cells[3])),
			Credits:  strings.TrimSpace(htmlquery.InnerText(cells[5])),
			Score:    strings.TrimSpace(htmlquery.InnerText(cells[6])),
		}

		resp = append(resp, mark)
	}

	return resp, nil
}
