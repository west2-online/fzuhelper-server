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

// Code generated by Kitex v0.11.3. DO NOT EDIT.

package classroomservice

import (
	"context"
	"errors"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	api "github.com/west2-online/fzuhelper-server/kitex_gen/api"
)

var errInvalidMessageType = errors.New("invalid message type for service method handler")

var serviceMethods = map[string]kitex.MethodInfo{
	"GetEmptyClassrooms": kitex.NewMethodInfo(
		getEmptyClassroomsHandler,
		newClassRoomServiceGetEmptyClassroomsArgs,
		newClassRoomServiceGetEmptyClassroomsResult,
		false,
		kitex.WithStreamingMode(kitex.StreamingNone),
	),
}

var (
	classRoomServiceServiceInfo                = NewServiceInfo()
	classRoomServiceServiceInfoForClient       = NewServiceInfoForClient()
	classRoomServiceServiceInfoForStreamClient = NewServiceInfoForStreamClient()
)

// for server
func serviceInfo() *kitex.ServiceInfo {
	return classRoomServiceServiceInfo
}

// for stream client
func serviceInfoForStreamClient() *kitex.ServiceInfo {
	return classRoomServiceServiceInfoForStreamClient
}

// for client
func serviceInfoForClient() *kitex.ServiceInfo {
	return classRoomServiceServiceInfoForClient
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
	serviceName := "ClassRoomService"
	handlerType := (*api.ClassRoomService)(nil)
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
		"PackageName": "api",
	}
	if hasStreaming {
		extra["streaming"] = hasStreaming
	}
	svcInfo := &kitex.ServiceInfo{
		ServiceName:     serviceName,
		HandlerType:     handlerType,
		Methods:         methods,
		PayloadCodec:    kitex.Thrift,
		KiteXGenVersion: "v0.11.3",
		Extra:           extra,
	}
	return svcInfo
}

func getEmptyClassroomsHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*api.ClassRoomServiceGetEmptyClassroomsArgs)
	realResult := result.(*api.ClassRoomServiceGetEmptyClassroomsResult)
	success, err := handler.(api.ClassRoomService).GetEmptyClassrooms(ctx, realArg.Request)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newClassRoomServiceGetEmptyClassroomsArgs() interface{} {
	return api.NewClassRoomServiceGetEmptyClassroomsArgs()
}

func newClassRoomServiceGetEmptyClassroomsResult() interface{} {
	return api.NewClassRoomServiceGetEmptyClassroomsResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) GetEmptyClassrooms(ctx context.Context, request *api.EmptyClassroomRequest) (r *api.EmptyClassroomResponse, err error) {
	var _args api.ClassRoomServiceGetEmptyClassroomsArgs
	_args.Request = request
	var _result api.ClassRoomServiceGetEmptyClassroomsResult
	if err = p.c.Call(ctx, "GetEmptyClassrooms", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}
