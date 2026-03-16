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
	"sort"
	"strconv"
	"strings"

	"github.com/west2-online/jwch/constants"
	"github.com/west2-online/jwch/errno"
	"github.com/west2-online/jwch/utils"

	"github.com/antchfx/htmlquery"
)

// 获取我的学期
func (s *Student) GetTerms() (*Term, error) {
	resp, err := s.GetWithIdentifier(constants.CourseURL)
	if err != nil {
		return nil, err
	}

	res := &Term{}

	res.ViewState = htmlquery.SelectAttr(htmlquery.FindOne(resp, `//*[@id="__VIEWSTATE"]`), "value")
	res.EventValidation = htmlquery.SelectAttr(htmlquery.FindOne(resp, `//*[@id="__EVENTVALIDATION"]`), "value")

	// 获取学年学期，例如 202202/202201/202102/202101 需要获取value
	list := htmlquery.Find(resp, `//*[@id="ContentPlaceHolder1_DDL_xnxq"]/option/@value`)

	// 这里考虑过使用 len(list) < 1，但是实际上这没必要，因为小于1那么它必定是0
	if len(list) == 0 {
		return nil, errno.HTMLParseError.WithMessage("empty terms")
	}

	for _, node := range list {
		res.Terms = append(res.Terms, htmlquery.SelectAttr(node, "value"))
	}

	return res, nil
}

// 获取我的选课
func (s *Student) GetSemesterCourses(term, viewState, eventValidation string) ([]*Course, error) {
	resp, err := s.PostWithIdentifier(constants.CourseURL, map[string]string{
		"ctl00$ContentPlaceHolder1$DDL_xnxq":  term,
		"ctl00$ContentPlaceHolder1$BT_submit": "确定",
		"__VIEWSTATE":                         viewState,
		"__EVENTVALIDATION":                   eventValidation,
	})
	if err != nil {
		return nil, err
	}

	list := htmlquery.Find(htmlquery.FindOne(resp, `//*[@id="ContentPlaceHolder1_DataList_xxk"]/tbody`), "tr")

	// 去除第一个元素，第一个元素是标题栏，有个判断文本是“课程名称”
	// TODO: 我们如何确保第一个元素一定是标题栏?
	list = list[2:]

	res := make([]*Course, 0)

	for _, node := range list {
		// 教务处的表格HTML是不规范的，因此XPath解析会出现一些BUG
		if strings.TrimSpace(htmlquery.SelectAttr(node, "style")) == "" {
			continue
		}
		info := htmlquery.Find(node, `td`) // 一行的所有信息

		// 这个表格有12栏
		if len(info) < 12 {
			return nil, errno.HTMLParseError.WithMessage("get course info failed")
		}

		// 解析调课信息
		// 第二个理论上来说是一个不标准的调课信息，但是不知道为什么被录入了教务系统导致炸掉，所以修改了一下解析的正则表达式来做了兼容。 -- @renbaoshuo
		/*
			06周 星期3:5-6节  调至  09周 星期1:7-8节  旗山西1-206
			4 周 星期2:3-4节  调至  05周 星期2:7-8节  旗山东3-101
		*/
		courseInfo11 := strings.Split(utils.InnerTextWithBr(info[11]), "\n")
		// 注意：下面的正则里面有 NO-BREAK SPACE (U+00A0 %C2%A0)
		adjustRegex := regexp.MustCompile(`(\d{1,2})[\s ]*周[\s ]*星期(\d):(\d{1,2})-(\d{1,2})节[\s ]*调至[\s ]*(\d{1,2})[\s ]*周[\s ]*星期(\d):(\d{1,2})-(\d{1,2})节[\s ]*(\S*)`)
		adjustRules := []CourseAdjustRule{}

		for i := 0; i < len(courseInfo11); i++ {
			courseInfo11[i] = strings.TrimSpace(courseInfo11[i])

			if courseInfo11[i] == "" { // 空行
				continue
			}

			adjustMatchArr := adjustRegex.FindStringSubmatch(courseInfo11[i])

			if len(adjustMatchArr) < 10 {
				return nil, errno.HTMLParseError.WithMessage("get course adjust failed")
			}

			adjustRules = append(adjustRules, CourseAdjustRule{
				OldWeek:       utils.SafeAtoi(adjustMatchArr[1]),
				OldWeekday:    utils.SafeAtoi(adjustMatchArr[2]),
				OldStartClass: utils.SafeAtoi(adjustMatchArr[3]),
				OldEndClass:   utils.SafeAtoi(adjustMatchArr[4]),

				NewWeek:       utils.SafeAtoi(adjustMatchArr[5]),
				NewWeekday:    utils.SafeAtoi(adjustMatchArr[6]),
				NewStartClass: utils.SafeAtoi(adjustMatchArr[7]),
				NewEndClass:   utils.SafeAtoi(adjustMatchArr[8]),
				NewLocation:   adjustMatchArr[9],
			})
		}

		// 解析上课时间、地点，融合调课信息
		/*
			05-18 星期1:3-4节 铜盘A110
			05-17 星期3:1-2节 铜盘A110
			05-17 星期5:3-4节 铜盘A110
		*/
		courseInfo8 := strings.Split(utils.InnerTextWithBr(info[8]), "\n")
		scheduleRules := []CourseScheduleRule{}
		fullWeekScheduleRules := []CourseFullWeekScheduleRule{}

		for i := 0; i < len(courseInfo8); i++ {
			courseInfo8[i] = strings.TrimSpace(courseInfo8[i])

			if courseInfo8[i] == "" { // 空行
				continue
			}

			lineData := strings.Fields(courseInfo8[i])

			if len(lineData) < 3 {
				return nil, errno.HTMLParseError.WithMessage("get course info failed")
			}

			if strings.Contains(lineData[0], "周") { // 处理整周的课程，比如军训
				/*
					03周  星期1  -  04周  星期7
					[0] 03周
					[1] 星期1
					[2] -
					[3] 04周
					[4] 星期7
				*/
				startWeek, _ := strconv.Atoi(strings.TrimSuffix(lineData[0], "周"))
				endWeek, _ := strconv.Atoi(strings.TrimSuffix(lineData[3], "周"))
				startWeekday, _ := strconv.Atoi(strings.TrimPrefix(lineData[1], "星期"))
				endWeekday, _ := strconv.Atoi(strings.TrimPrefix(lineData[4], "星期"))

				// 向普通课程的格式转换（应用于课表显示）
				for weekday := 1; weekday <= 7; weekday++ {
					curStartWeek := startWeek
					curEndWeek := endWeek

					if weekday < startWeekday {
						curStartWeek++
					}

					if weekday > endWeekday {
						curEndWeek--
					}

					if curStartWeek > curEndWeek {
						continue
					}

					scheduleRules = append(scheduleRules, CourseScheduleRule{
						Location:     "",
						StartClass:   1,
						EndClass:     8,
						StartWeek:    curStartWeek,
						EndWeek:      curEndWeek,
						Weekday:      weekday,
						Single:       true,
						Double:       true,
						Adjust:       false,
						FromFullWeek: true,
					})
				}

				// 记录整周课程的信息（应用于日历生成）
				fullWeekScheduleRules = append(fullWeekScheduleRules, CourseFullWeekScheduleRule{
					StartWeek:    startWeek,
					StartWeekDay: startWeekday,
					EndWeek:      endWeek,
					EndWeekDay:   endWeekday,
				})
			} else { // 处理周内的正常课程
				/*
					08-16 星期5:7-8节 铜盘A508
					[0] 08-16
					[1] 星期5:7-8节
					[2] 铜盘A508
				*/
				/*
					02-14 星期1:1-2节(双) 旗山西1-206
					[0] 02-14
					[1] 星期1:1-2节(双)
					[2] 旗山西1-206
				*/
				/*
					01-13 星期1:3-4节(单) 旗山西1-206
					[0] 01-13
					[1] 星期1:3-4节(单)
					[2] 旗山西1-206
				*/

				// 是不是用正则表达式更好一点？
				weekInfo := strings.SplitN(lineData[0], "-", 2)    // [8, 16]
				dayInfo := strings.SplitN(lineData[1], ":", 2)     // ["星期5", "7-8节"] or ["星期1", "1-2节(双)"]
				classBasicInfo := strings.Split(dayInfo[1], "节")   // ["7-8", ""] or ["1-2", "(双)"]
				classInfo := strings.Split(classBasicInfo[0], "-") // ["7", "8"]
				location := lineData[2]
				startClass := utils.SafeAtoi(classInfo[0])
				endClass := utils.SafeAtoi(classInfo[1])
				startWeek := utils.SafeAtoi(weekInfo[0])
				endWeek := utils.SafeAtoi(weekInfo[1])
				weekDay := utils.SafeAtoi(strings.TrimPrefix(dayInfo[0], "星期"))
				single := !strings.Contains(classBasicInfo[1], "双")
				double := !strings.Contains(classBasicInfo[1], "单")

				if len(adjustRules) == 0 {
					scheduleRules = append(scheduleRules, CourseScheduleRule{
						Location:     location,
						StartClass:   startClass,
						EndClass:     endClass,
						StartWeek:    startWeek,
						EndWeek:      endWeek,
						Weekday:      weekDay,
						Single:       single,
						Double:       double,
						Adjust:       false,
						FromFullWeek: false,
					})
				} else {
					startWeek := utils.SafeAtoi(weekInfo[0])
					endWeek := utils.SafeAtoi(weekInfo[1])
					startClass := utils.SafeAtoi(classInfo[0])
					endClass := utils.SafeAtoi(classInfo[1])
					removedWeeks := []int{}

					for _, adjustRule := range adjustRules {
						// 匹配是否是对应的调课信息
						if adjustRule.OldWeek < startWeek ||
							adjustRule.OldWeek > endWeek ||
							adjustRule.OldStartClass != startClass ||
							adjustRule.OldEndClass != endClass ||
							adjustRule.OldWeekday != weekDay {
							continue
						}

						// 记录被去掉的周次
						removedWeeks = append(removedWeeks, adjustRule.OldWeek)

						// 添加新的课程信息
						scheduleRules = append(scheduleRules, CourseScheduleRule{
							Location:     adjustRule.NewLocation,
							StartClass:   adjustRule.NewStartClass,
							EndClass:     adjustRule.NewEndClass,
							StartWeek:    adjustRule.NewWeek,
							EndWeek:      adjustRule.NewWeek,
							Weekday:      adjustRule.NewWeekday,
							Single:       true,
							Double:       true,
							Adjust:       true, // 调课
							FromFullWeek: false,
						})
					}

					sort.Ints(removedWeeks)
					// 去掉被调课的周次
					curStartWeek := startWeek

					for _, removedWeek := range removedWeeks {
						if removedWeek == curStartWeek {
							curStartWeek++

							continue
						}

						scheduleRules = append(scheduleRules, CourseScheduleRule{
							Location:     location,
							StartClass:   startClass,
							EndClass:     endClass,
							StartWeek:    curStartWeek,
							EndWeek:      removedWeek - 1,
							Weekday:      weekDay,
							Single:       single,
							Double:       double,
							Adjust:       false,
							FromFullWeek: false,
						})

						curStartWeek = removedWeek + 1
					}

					if curStartWeek <= endWeek {
						scheduleRules = append(scheduleRules, CourseScheduleRule{
							Location:     location,
							StartClass:   startClass,
							EndClass:     endClass,
							StartWeek:    curStartWeek,
							EndWeek:      endWeek,
							Weekday:      weekDay,
							Single:       single,
							Double:       double,
							Adjust:       false,
							FromFullWeek: false,
						})
					}
				}
			}
		}

		// TODO: performance optimization
		res = append(res, &Course{
			Type:       htmlquery.OutputHTML(info[0], false),
			Name:       htmlquery.OutputHTML(info[1], false),
			Syllabus:   constants.JwchPrefix + safeExtractRegex(`javascript:pop1\('(.*?)&`, safeExtractionValue(info[2], "a", "href", 0)),
			LessonPlan: constants.JwchPrefix + safeExtractRegex(`javascript:pop1\('(.*?)&`, safeExtractionValue(info[2], "a", "href", 1)),
			// PaymentStatus: safeExtractionFirst(info[3], "font"),
			Credits:               safeExtractionFirst(info[4], "span"),
			ElectiveType:          utils.GetChineseCharacter(htmlquery.OutputHTML(info[5], false)),
			ExamType:              utils.GetChineseCharacter(htmlquery.OutputHTML(info[6], false)),
			Teacher:               htmlquery.OutputHTML(info[7], false),
			ScheduleRules:         scheduleRules,
			FullWeekScheduleRules: fullWeekScheduleRules,
			RawScheduleRules:      strings.Join(courseInfo8, "\n"),
			RawExamTime:           strings.TrimSpace(htmlquery.InnerText(info[9])),
			RawAdjust:             strings.Join(courseInfo11, "\n"),
			Remark:                htmlquery.OutputHTML(info[10], false),
		})
	}

	return res, nil
}

func (s *Student) GetLocateDate() (*LocateDate, error) {
	resp, err := s.NewRequest().Get(constants.JwchLocateDateUrl)
	if err != nil {
		return nil, err
	}
	data := string(resp.Body())

	// 使用正则表达式解析返回内容
	re := regexp.MustCompile(`var week = "([0-9]+)";\s*//.*\s*var xn = "([0-9]{4})";\s*//.*\s*var xq = "([0-9]{2})";`)
	matches := re.FindStringSubmatch(data)
	if len(matches) < 4 {
		return nil, errno.HTMLParseError.WithMessage("failed to parse response from JWCH_LOCATE_DATE_URL")
	}

	return &LocateDate{Week: matches[1], Year: matches[2], Term: matches[3]}, nil
}
