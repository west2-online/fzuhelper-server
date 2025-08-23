package main

import (
	"context"
	paper "github.com/west2-online/fzuhelper-server/kitex_gen/paper"
)

// PaperServiceImpl implements the last service interface defined in the IDL.
type PaperServiceImpl struct{}

// ListDirFiles implements the PaperServiceImpl interface.
func (s *PaperServiceImpl) ListDirFiles(ctx context.Context, req *paper.ListDirFilesRequest) (resp *paper.ListDirFilesResponse, err error) {
	// TODO: Your code here...
	return
}

// GetDownloadUrl implements the PaperServiceImpl interface.
func (s *PaperServiceImpl) GetDownloadUrl(ctx context.Context, req *paper.GetDownloadUrlRequest) (resp *paper.GetDownloadUrlResponse, err error) {
	// TODO: Your code here...
	return
}
