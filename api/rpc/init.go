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

package rpc

import (
	"github.com/west2-online/fzuhelper-server/kitex_gen/academic/academicservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/classroom/classroomservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/common/commonservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/course/courseservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/launch_screen/launchscreenservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/paper/paperservice"
	"github.com/west2-online/fzuhelper-server/kitex_gen/user/userservice"
)

var (
	classroomClient          classroomservice.Client
	courseClient             courseservice.Client
	userClient               userservice.Client
	launchScreenClient       launchscreenservice.Client
	launchScreenStreamClient launchscreenservice.StreamClient
	paperClient              paperservice.Client
	academicClient           academicservice.Client
	commonClient             commonservice.Client
)

func Init() {
	InitClassroomRPC()
	InitUserRPC()
	InitCourseRPC()
	InitLaunchScreenRPC()
	InitLaunchScreenStreamRPC()
	InitPaperRPC()
	InitAcademicRPC()
	InitCommonRPC()
}
