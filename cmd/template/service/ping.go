package service

import "github.com/west2-online/fzuhelper-server/kitex_gen/template"

func (s *TemplateService) Ping(req *template.PingRequest) (string, error) {
	return *req.Text, nil

	// 这里负责处理业务请求，如果有需要可以继续和 dal（data access layout）做交互
}
