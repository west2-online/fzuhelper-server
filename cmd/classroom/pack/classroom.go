package pack

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom"
	"github.com/west2-online/fzuhelper-server/pkg/utils"
	"regexp"
	"strings"
)

func BuildClassroom(str string) (res *classroom.Classroom) {
	//旗山东1-103 0(0) 机房
	//晋江A102 150(75) 多媒体
	//铜盘A109 120(60) 多媒体
	//泉港教-110 63(40) 多媒体
	//怡山北301 92(0) 多媒体
	//鼓浪屿多媒体1 0(0) 多媒体
	//集美1-311 287(287) 多媒体
	res = new(classroom.Classroom)
	buildingPattern := regexp.MustCompile(`([^\d]+[\d]*)`)
	locationPattern := regexp.MustCompile(`[\d]+(-[\d]+)?`)
	seatsPattern := regexp.MustCompile(`\d+`)

	// 使用正则表达式提取 building 部分
	building := buildingPattern.FindString(str)

	// 使用正则表达式提取 location 部分
	location := locationPattern.FindString(str)

	// 剩余部分
	remaining := strings.TrimPrefix(str, building+location)

	// 从剩余部分提取人数和教室类型
	parts := strings.Fields(remaining)
	totalSeats := seatsPattern.FindString(parts[0]) // 提取括号外的数字
	classroomType := parts[1]

	// 输出结果
	utils.LoggerObj.Infof("classroom.pack.Buildclassroom: Building: %s, Location: %s, 人数: %s, 教室类型: %s\n", strings.TrimSpace(building), location, totalSeats, classroomType)
	res.Build = strings.TrimSpace(building)
	res.Location = location
	res.Capacity = totalSeats
	res.Type = classroomType

	return res
}

func BuildClassRooms(strs []string) (res []*classroom.Classroom) {
	for _, str := range strs {
		res = append(res, BuildClassroom(str))
	}
	return
}
