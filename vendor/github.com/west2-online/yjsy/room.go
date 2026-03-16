package yjsy

import (
	"fmt"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/west2-online/yjsy/constants"
	"github.com/west2-online/yjsy/errno"
	"golang.org/x/net/html"
)

func (s *Student) GetExamRoom(req ExamRoomReq) ([]*ExamRoomInfo, error) {

	res, err := s.GetWithIdentifier(constants.ExamRoomQueryURL, map[string]string{
		// 非常抽象的格式,小心踩坑
		"strwhere": fmt.Sprintf("XNXQ='%s'", req.Term),
	})

	if err != nil {
		return nil, err
	}
	examInfos, err := parseExamRoom(res)
	if err != nil {
		return nil, err
	}
	return examInfos, nil
}

func parseExamRoom(doc *html.Node) ([]*ExamRoomInfo, error) {
	var examInfos []*ExamRoomInfo
	table := htmlquery.FindOne(doc, "//div[@id='divContent']/table")
	if table == nil {
		return nil, errno.HTMLParseError
	}
	rows := htmlquery.Find(table, "//tr[position()>1]")

	for _, row := range rows {
		columns := htmlquery.Find(row, "./td")
		if len(columns) < 6 {
			continue // 如果列数不足 6，则跳过这一行
		}
		var date string
		var time string
		dateTime := strings.Fields(htmlquery.InnerText(columns[4]))
		if len(dateTime) == 0 {
			date = ""
			time = ""
		} else {
			date = strings.Replace(dateTime[0], "/", "-", 2)
			time = dateTime[1]
		}

		examInfo := &ExamRoomInfo{
			CourseName: strings.TrimSpace(htmlquery.InnerText(columns[1])),
			Credit:     "", // 页面没有学分信息，留空
			Teacher:    "", // 页面没有教师信息，留空
			Date:       date,
			Time:       time,
			Location:   strings.TrimSpace(htmlquery.InnerText(columns[5])),
		}

		examInfos = append(examInfos, examInfo)
	}

	return examInfos, nil
}
