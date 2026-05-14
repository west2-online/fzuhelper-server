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

package singleflight

import "fmt"

func ScoresKey(stuID string, isGraduate bool) string {
	// 本科和研究生成绩来自不同上游，返回结构也不同，需要按身份隔离。
	return fmt.Sprintf("scores:%s:%t", stuID, isGraduate)
}

func ExamRoomsKey(stuID, term string, isGraduate bool) string {
	// 考场结果按学期和身份区分，同一学生不同学期或不同身份不能复用结果。
	return fmt.Sprintf("exam_rooms:%s:%s:%t", stuID, term, isGraduate)
}

func CourseListKey(stuID, term string, isGraduate, isRefresh bool) string {
	// 将刷新标记纳入 key，避免强刷请求复用普通请求的 singleflight 结果。
	return fmt.Sprintf("courses:%s:%s:%t:%t", stuID, term, isGraduate, isRefresh)
}

func CourseTermsKey(stuID string, isGraduate bool) string {
	// 本科和研究生学期来源不同，同一学号也要按身份隔离。
	return fmt.Sprintf("terms:%s:%t", stuID, isGraduate)
}

func TermKey(term string) string {
	return fmt.Sprintf("term:%s", term)
}

func NoticeKey(pageNum int64) string {
	return fmt.Sprintf("notice:%d", pageNum)
}

func PaperDirKey(path string) string {
	return fmt.Sprintf("dir:%s", path)
}

func UserInfoKey(stuID string, isGraduate bool) string {
	// 本科和研究生用户信息来自不同上游，按身份隔离避免复用到错误来源的数据。
	return fmt.Sprintf("user_info:%s:%t", stuID, isGraduate)
}

func FriendListKey(stuID string) string {
	return fmt.Sprintf("friend_list:%s", stuID)
}
