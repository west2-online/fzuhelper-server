package yjsy

import (
	"sort"
	"strconv"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/west2-online/yjsy/constants"
)

func (s *Student) GetTerms() (*Term, error) {
	resp, err := s.GetWithIdentifier(constants.TermURL, map[string]string{})
	if err != nil {
		return nil, err
	}
	rows := htmlquery.Find(resp, `//div[@id='divContent']//table//tr[position()>1]`)
	terms := new(Term)
	for _, row := range rows {
		cells := htmlquery.Find(row, `td`)
		term := strings.TrimSpace(htmlquery.InnerText(cells[0]))
		terms.Terms = append(terms.Terms, term)
	}
	// 去重
	termsMap := make(map[string]bool)
	for _, term := range terms.Terms {
		termsMap[term] = true
	}
	terms.Terms = make([]string, 0, len(termsMap))
	for term := range termsMap {
		terms.Terms = append(terms.Terms, term)
	}
	// 降序
	sort.Slice(terms.Terms, func(i, j int) bool {
		return terms.Terms[i] > terms.Terms[j]
	})
	return terms, nil
}

func (s *Student) GetSemesterCourses(term string) ([]*Course, error) {
	allCourses := make([]*Course, 0)
	// 递归解析当前页数据
	currentURL := constants.CourseURL
	courses, nextPageURL, err := s.parseSinglePageByXNXQ(currentURL, term)
	if err != nil {
		return nil, err
	}

	// 追加当前页课程
	allCourses = append(allCourses, courses...)

	// 如果没有下一页，结束循环
	for nextPageURL != "" {
		currentURL = nextPageURL
		courses, nextPageURL, err = s.parseNextPage(currentURL)
		if err != nil {
			return nil, err
		}
		allCourses = append(allCourses, courses...)
	}

	return allCourses, nil
}

func (s *Student) parseSinglePageByXNXQ(url string, term string) ([]*Course, string, error) {
	resp, err := s.GetWithIdentifier(url, map[string]string{
		"strwhere": strings.Join([]string{"XNXQ='", term, "'"}, ""),
	})
	if err != nil {
		return nil, "", err
	}

	// Locate the rows in the course table
	rows := htmlquery.Find(resp, `//div[@id='divContent']//table//tr[position()>1]`)

	courses := make([]*Course, 0)

	for _, row := range rows {
		cells := htmlquery.Find(row, `td`)
		if len(cells) < 8 {
			continue
		}
		// Parse fields
		term := strings.TrimSpace(htmlquery.InnerText(cells[0]))
		name := strings.TrimSpace(htmlquery.InnerText(cells[2]))
		teacher := strings.TrimSpace(htmlquery.InnerText(cells[5]))
		rawScheduleHTML := htmlquery.OutputHTML(cells[6], false) // Extract full HTML for schedule rules to recognize multiple lessons in one week
		remark := strings.TrimSpace(htmlquery.InnerText(cells[8]))
		lessonPlan := ""
		lessonPlanLink := htmlquery.FindOne(cells[7], `.//a[@href]`)
		if lessonPlanLink != nil {
			lessonPlan = htmlquery.SelectAttr(lessonPlanLink, "href")
			lessonPlan = strings.Join([]string{constants.YjsyPrefix, strings.TrimPrefix(lessonPlan, "..")}, "")
		}

		// Parse schedule rules
		scheduleRules := parseScheduleRulesFromHTML(rawScheduleHTML)

		// Append to the result
		courses = append(courses, &Course{
			Name:             name,
			Teacher:          teacher,
			ScheduleRules:    scheduleRules,
			Remark:           remark,
			LessonPlan:       lessonPlan,
			RawScheduleRules: rawScheduleHTML,
			RawAdjust:        "",
			Term:             term,
		})
	}

	nextPage := htmlquery.FindOne(resp, `//div[@id='divPage']//a[contains(text(), '下一页')]`)
	var nextPageURL string
	if nextPage != nil {
		href := htmlquery.SelectAttr(nextPage, "href")
		nextPageURL = strings.Join([]string{constants.CourseURL, href}, "")

	}
	return courses, nextPageURL, nil
}

// Function to parse schedule rules from HTML
func parseScheduleRulesFromHTML(rawScheduleHTML string) []CourseScheduleRule {
	// Replace <br> tags with newlines
	rawScheduleHTML = strings.ReplaceAll(rawScheduleHTML, "<br>", "\n")
	rawScheduleHTML = strings.ReplaceAll(rawScheduleHTML, "<br/>", "\n")
	return parseScheduleRules(rawScheduleHTML)
}

// Existing parseScheduleRules function
func parseScheduleRules(rawScheduleRules string) []CourseScheduleRule {
	lines := strings.Split(rawScheduleRules, "\n")
	var rules []CourseScheduleRule

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Example: "1-8周 星期3:9-11节 东3-109"
		parts := strings.Fields(line)
		if len(parts) < 3 {
			continue
		}

		// Parsing week, day, and location
		weekInfo := strings.Split(parts[0], "-")
		dayInfo := strings.Split(parts[1], ":")
		classInfo := strings.Split(strings.TrimSuffix(dayInfo[1], "节"), "-")

		startWeek, _ := strconv.Atoi(strings.TrimSuffix(weekInfo[0], "周"))
		endWeek, _ := strconv.Atoi(strings.TrimSuffix(weekInfo[1], "周"))
		weekday, _ := strconv.Atoi(strings.TrimPrefix(dayInfo[0], "星期"))
		startClass, _ := strconv.Atoi(classInfo[0])
		endClass, _ := strconv.Atoi(classInfo[1])
		location := parts[2]

		rules = append(rules, CourseScheduleRule{
			Location:   location,
			StartClass: startClass,
			EndClass:   endClass,
			StartWeek:  startWeek,
			EndWeek:    endWeek,
			Weekday:    weekday,
			Single:     true,
			Double:     true,
			Adjust:     false,
		})
	}

	return rules
}

func (s *Student) parseNextPage(url string) ([]*Course, string, error) {
	resp, err := s.GetWithIdentifier(url, map[string]string{})
	if err != nil {
		return nil, "", err
	}

	// Locate the rows in the course table
	rows := htmlquery.Find(resp, `//div[@id='divContent']//table//tr[position()>1]`)

	courses := make([]*Course, 0)

	for _, row := range rows {
		cells := htmlquery.Find(row, `td`)

		// Parse fields
		term := strings.TrimSpace(htmlquery.InnerText(cells[0]))
		name := strings.TrimSpace(htmlquery.InnerText(cells[2]))
		teacher := strings.TrimSpace(htmlquery.InnerText(cells[5]))
		rawScheduleHTML := htmlquery.OutputHTML(cells[6], false) // Extract full HTML for schedule rules to recognize multiple lessons in one week
		remark := strings.TrimSpace(htmlquery.InnerText(cells[8]))
		lessonPlan := ""
		lessonPlanLink := htmlquery.FindOne(cells[7], `.//a[@href]`)
		if lessonPlanLink != nil {
			lessonPlan = htmlquery.SelectAttr(lessonPlanLink, "href")
			lessonPlan = strings.Join([]string{constants.YjsyPrefix, strings.TrimPrefix(lessonPlan, "..")}, "")
		}

		// Parse schedule rules
		scheduleRules := parseScheduleRulesFromHTML(rawScheduleHTML)

		// Append to the result
		courses = append(courses, &Course{
			Name:             name,
			Teacher:          teacher,
			ScheduleRules:    scheduleRules,
			Remark:           remark,
			LessonPlan:       lessonPlan,
			RawScheduleRules: rawScheduleHTML,
			RawAdjust:        "",
			Term:             term,
		})
	}

	nextPage := htmlquery.FindOne(resp, `//div[@id='divPage']//a[contains(text(), '下一页')]`)
	var nextPageURL string
	if nextPage != nil {
		href := htmlquery.SelectAttr(nextPage, "href")
		nextPageURL = strings.Join([]string{constants.CourseURL, href}, "")

	}
	return courses, nextPageURL, nil
}
