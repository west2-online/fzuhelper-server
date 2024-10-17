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

// Code generated by Kitex v0.7.1. DO NOT EDIT.

package paperservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	kitex "github.com/cloudwego/kitex/pkg/serviceinfo"
	paper "github.com/west2-online/fzuhelper-server/kitex_gen/paper"
)

func serviceInfo() *kitex.ServiceInfo {
	return paperServiceServiceInfo
}

var paperServiceServiceInfo = NewServiceInfo()

func NewServiceInfo() *kitex.ServiceInfo {
	serviceName := "PaperService"
	handlerType := (*paper.PaperService)(nil)
	methods := map[string]kitex.MethodInfo{
		"UploadFile":     kitex.NewMethodInfo(uploadFileHandler, newPaperServiceUploadFileArgs, newPaperServiceUploadFileResult, false),
		"ListDirFiles":   kitex.NewMethodInfo(listDirFilesHandler, newPaperServiceListDirFilesArgs, newPaperServiceListDirFilesResult, false),
		"GetDownloadUrl": kitex.NewMethodInfo(getDownloadUrlHandler, newPaperServiceGetDownloadUrlArgs, newPaperServiceGetDownloadUrlResult, false),
	}
	extra := map[string]interface{}{
		"PackageName":     "paper",
		"ServiceFilePath": "idl/paper.thrift",
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

func uploadFileHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*paper.PaperServiceUploadFileArgs)
	realResult := result.(*paper.PaperServiceUploadFileResult)
	success, err := handler.(paper.PaperService).UploadFile(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newPaperServiceUploadFileArgs() interface{} {
	return paper.NewPaperServiceUploadFileArgs()
}

func newPaperServiceUploadFileResult() interface{} {
	return paper.NewPaperServiceUploadFileResult()
}

func listDirFilesHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*paper.PaperServiceListDirFilesArgs)
	realResult := result.(*paper.PaperServiceListDirFilesResult)
	success, err := handler.(paper.PaperService).ListDirFiles(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newPaperServiceListDirFilesArgs() interface{} {
	return paper.NewPaperServiceListDirFilesArgs()
}

func newPaperServiceListDirFilesResult() interface{} {
	return paper.NewPaperServiceListDirFilesResult()
}

func getDownloadUrlHandler(ctx context.Context, handler interface{}, arg, result interface{}) error {
	realArg := arg.(*paper.PaperServiceGetDownloadUrlArgs)
	realResult := result.(*paper.PaperServiceGetDownloadUrlResult)
	success, err := handler.(paper.PaperService).GetDownloadUrl(ctx, realArg.Req)
	if err != nil {
		return err
	}
	realResult.Success = success
	return nil
}
func newPaperServiceGetDownloadUrlArgs() interface{} {
	return paper.NewPaperServiceGetDownloadUrlArgs()
}

func newPaperServiceGetDownloadUrlResult() interface{} {
	return paper.NewPaperServiceGetDownloadUrlResult()
}

type kClient struct {
	c client.Client
}

func newServiceClient(c client.Client) *kClient {
	return &kClient{
		c: c,
	}
}

func (p *kClient) UploadFile(ctx context.Context, req *paper.UploadFileRequest) (r *paper.UploadFileResponse, err error) {
	var _args paper.PaperServiceUploadFileArgs
	_args.Req = req
	var _result paper.PaperServiceUploadFileResult
	if err = p.c.Call(ctx, "UploadFile", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) ListDirFiles(ctx context.Context, req *paper.ListDirFilesRequest) (r *paper.ListDirFilesResponse, err error) {
	var _args paper.PaperServiceListDirFilesArgs
	_args.Req = req
	var _result paper.PaperServiceListDirFilesResult
	if err = p.c.Call(ctx, "ListDirFiles", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}

func (p *kClient) GetDownloadUrl(ctx context.Context, req *paper.GetDownloadUrlRequest) (r *paper.GetDownloadUrlResponse, err error) {
	var _args paper.PaperServiceGetDownloadUrlArgs
	_args.Req = req
	var _result paper.PaperServiceGetDownloadUrlResult
	if err = p.c.Call(ctx, "GetDownloadUrl", &_args, &_result); err != nil {
		return
	}
	return _result.GetSuccess(), nil
}