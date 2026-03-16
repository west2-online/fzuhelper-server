// go
package yjsy

import (
	"github.com/west2-online/yjsy/constants"
)

// GetStudentInfo 请求学生信息页面并解析为 HTML 文档后提取学生信息
func (s *Student) GetStudentInfo() (*StudentDetail, error) {
	resp, err := s.GetWithIdentifier(constants.UserInfoURL, nil)
	if err != nil {
		return &StudentDetail{}, err
	}

	// 以下 xpath 表达式请根据实际响应 HTML 结构调整
	xpaths := map[string]string{
		"stu_id":   "//*[@id='xxTable']/table/tbody/tr[1]/td[2]",
		"name":     "//*[@id='xxTable']/table/tbody/tr[1]/td[4]",
		"birthday": "//*[@id='xxTable']/table/tbody/tr[3]/td[4]",
		"sex":      "//*[@id='xxTable']/table/tbody/tr[2]/td[4]",
		"college":  "//*[@id='xxTable']/table/tbody/tr[15]/td[2]",
		"grade":    "//*[@id='xxTable']/table/tbody/tr[14]/td[4]",
		"major":    "//*[@id='xxTable']/table/tbody/tr[15]/td[4]",
	}

	result := StudentDetail{}
	for key, xp := range xpaths {
		value := safeExtractHTMLFirst(resp, xp)
		switch key {
		case "stu_id":
			result.StuID = value
		case "name":
			result.Name = value
		case "birthday":
			result.Birthday = value
		case "sex":
			result.Sex = value
		case "college":
			result.College = value
		case "grade":
			result.Grade = value
		case "major":
			result.Major = value
		}
	}
	return &result, nil
}
