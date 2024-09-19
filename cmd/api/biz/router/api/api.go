// Code generated by hertz generator. DO NOT EDIT.

package api

import (
	"github.com/cloudwego/hertz/pkg/app/server"
	api "github.com/west2-online/fzuhelper-server/biz/handler/api"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_academic := root.Group("/academic", _academicMw()...)
		_academic.GET("/credit", append(_getcreditMw(), api.GetCredit)...)
		_academic.GET("/gpa", append(_getgpaMw(), api.GetGPA)...)
		_academic.GET("/plan", append(_getplanMw(), api.GetPlan)...)
		_academic.GET("/scores", append(_getscoresMw(), api.GetScores)...)
		_academic.GET("/unifiedExam", append(_getunifiedexamMw(), api.GetUnifiedExam)...)
	}
	{
		_classroom := root.Group("/classroom", _classroomMw()...)
		_classroom.GET("/empty", append(_getemptyroomMw(), api.GetEmptyRoom)...)
		_classroom.GET("/exam", append(_getexamMw(), api.GetExam)...)
	}
	{
		_course := root.Group("/course", _courseMw()...)
		_course.GET("/list", append(_getcourselistMw(), api.GetCourseList)...)
	}
	{
		_user := root.Group("/user", _userMw()...)
		_user.GET("/info", append(_getuserinfoMw(), api.GetUserInfo)...)
		_user.PUT("/info", append(_changepasswordMw(), api.ChangePassword)...)
		_user.GET("/schoolCalendar", append(_getschoolcalendarMw(), api.GetSchoolCalendar)...)
		_user.POST("/validateCode", append(_validatecodeMw(), api.ValidateCode)...)
	}
}
