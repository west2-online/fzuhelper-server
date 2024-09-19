// Code generated by Kitex v0.7.1. DO NOT EDIT.

package courseservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	course "github.com/west2-online/fzuhelper-server/kitex_gen/course"
)

func serviceInfo() *kitex.ServiceInfo {
	return courseServiceServiceInfo
}

var courseServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "CourseService"
	handlerType := (*course.CourseService)(nil)
	methods := map[string]kitex.MethodInfo{
		"GetCourseList": kitex.NewMethodInfo(getCourseListHandler, newCourseServiceGetCourseListArgs, newCourseServiceGetCourseListResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "course",
		"ServiceFilePath": "../../idl/course.thrift",
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.7.1",
		Extra:           extra,
	}
	return svcInfo
}

func getCourseListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*course.CourseServiceGetCourseListArgs)
	realResult := result.(*course.CourseServiceGetCourseListResult)
	success, err := handler.(course.CourseService).GetCourseList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCourseServiceGetCourseListArgs() interface{} {
	return course.NewCourseServiceGetCourseListArgs()
}

func newCourseServiceGetCourseListResult() interface{} {
	return course.NewCourseServiceGetCourseListResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) GetCourseList(ctx context.Context, req *course.CourseListRequest) (r *course.CourseListResponse, err error) {
	var _args course.CourseServiceGetCourseListArgs
	_args.Req = req
	var _result course.CourseServiceGetCourseListResult
	if err = p.c.Call(ctx, "GetCourseList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
