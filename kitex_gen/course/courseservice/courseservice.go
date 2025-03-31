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

// Code generated by Kitex v0.12.1. DO NOT EDIT.

package courseservice

import (
	"context"
	"errors"

	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"

	course "github.com/west2-online/fzuhelper-server/kitex_gen/course"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"GetCourseList": kitex.NewMethodInfo(
		getCourseListHandler,
		newCourseServiceGetCourseListArgs,
		newCourseServiceGetCourseListResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetTermList": kitex.NewMethodInfo(
		getTermListHandler,
		newCourseServiceGetTermListArgs,
		newCourseServiceGetTermListResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetCalendar": kitex.NewMethodInfo(
		getCalendarHandler,
		newCourseServiceGetCalendarArgs,
		newCourseServiceGetCalendarResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
	"GetLocateDate": kitex.NewMethodInfo(
		getLocateDateHandler,
		newCourseServiceGetLocateDateArgs,
		newCourseServiceGetLocateDateResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	courseServiceServiceInfo                = NewServiceInfo()
	courseServiceServiceInfoForClient       = NewServiceInfoForClient()
	courseServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return courseServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return courseServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return courseServiceServiceInfoForClient
}

// NewServiceInfo creates a new ServiceInfo containing all methods
func NewServiceInfo() *kitex.ServiceInfo {
	return newServiceInfo(false, true, true)
}

// NewServiceInfo creates a new ServiceInfo containing non-streaming methods
func NewServiceInfoForClient() *kitex.ServiceInfo {
	return newServiceInfo(false, false, true)
}
func NewServiceInfoForStreamClient() *kitex.ServiceInfo {
	return newServiceInfo(true, true, false)
}

func newServiceInfo(hasStreaming bool, keepStreamingMethods bool, keepNonStreamingMethods bool) *kitex.ServiceInfo {
	serviceName := "CourseService"
	handlerType := (*course.CourseService)(nil)
	methods := map[string]kitex.MethodInfo{}
	for name, m := range serviceMethods {
		if m.IsStreaming() && !keepStreamingMethods {
			continue
		}
		if !m.IsStreaming() && !keepNonStreamingMethods {
			continue
		}
		methods[name] = m
	}
	extra := map[string]interface{}{
		"PackageName": "course",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.12.1",
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

func getTermListHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*course.CourseServiceGetTermListArgs)
	realResult := result.(*course.CourseServiceGetTermListResult)
	success, err := handler.(course.CourseService).GetTermList(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCourseServiceGetTermListArgs() interface{} {
	return course.NewCourseServiceGetTermListArgs()
}

func newCourseServiceGetTermListResult() interface{} {
	return course.NewCourseServiceGetTermListResult()
}

func getCalendarHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*course.CourseServiceGetCalendarArgs)
	realResult := result.(*course.CourseServiceGetCalendarResult)
	success, err := handler.(course.CourseService).GetCalendar(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCourseServiceGetCalendarArgs() interface{} {
	return course.NewCourseServiceGetCalendarArgs()
}

func newCourseServiceGetCalendarResult() interface{} {
	return course.NewCourseServiceGetCalendarResult()
}

func getLocateDateHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*course.CourseServiceGetLocateDateArgs)
	realResult := result.(*course.CourseServiceGetLocateDateResult)
	success, err := handler.(course.CourseService).GetLocateDate(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newCourseServiceGetLocateDateArgs() interface{} {
	return course.NewCourseServiceGetLocateDateArgs()
}

func newCourseServiceGetLocateDateResult() interface{} {
	return course.NewCourseServiceGetLocateDateResult()
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

func (p *kClient) GetTermList(ctx context.Context, req *course.TermListRequest) (r *course.TermListResponse, err error) {
	var _args course.CourseServiceGetTermListArgs
	_args.Req = req
	var _result course.CourseServiceGetTermListResult
	if err = p.c.Call(ctx, "GetTermList", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetCalendar(ctx context.Context, req *course.GetCalendarRequest) (r *course.GetCalendarResponse, err error) {
	var _args course.CourseServiceGetCalendarArgs
	_args.Req = req
	var _result course.CourseServiceGetCalendarResult
	if err = p.c.Call(ctx, "GetCalendar", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetLocateDate(ctx context.Context, req *course.GetLocateDateRequest) (r *course.GetLocateDateResponse, err error) {
	var _args course.CourseServiceGetLocateDateArgs
	_args.Req = req
	var _result course.CourseServiceGetLocateDateResult
	if err = p.c.Call(ctx, "GetLocateDate", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
