package yjsy

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"strings"

	"github.com/antchfx/htmlquery"
	"github.com/go-resty/resty/v2"
	"github.com/west2-online/yjsy/constants"
	"github.com/west2-online/yjsy/errno"
	"golang.org/x/net/html"
)

func NewStudent() *Student {
	// 从环境变量加载配置
	config := LoadConfigFromEnv()
	// Disable HTTP/2.0
	// Disable Redirect
	transport := &http.Transport{
		TLSNextProto:    make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	// 如果启用了代理，先获取隧道地址再设置代理
	if config.Proxy.Enabled {
		_, err := config.GetTunnelAddress()
		if err == nil && config.Proxy.ProxyServer != "" {
			proxyURL, err := config.GetProxyURL()
			if err == nil {
				transport.Proxy = http.ProxyURL(proxyURL)
			}
		}
	}

	client := resty.New().
		SetTransport(transport).
		SetRedirectPolicy(resty.NoRedirectPolicy())

	return &Student{
		client: client,
	}
}
func (s *Student) WithLoginData(cookies []*http.Cookie) *Student {
	s.cookies = cookies
	s.client.SetCookies(cookies)
	return s
}

// WithUser 携带账号密码
func (s *Student) WithUser(id, password string) *Student {
	s.ID = id
	s.Password = password
	return s
}

func (s *Student) SetCookies(cookies []*http.Cookie) {
	s.cookies = cookies
	s.client.SetCookies(cookies)
}

func (s *Student) ClearLoginData() {
	s.cookies = []*http.Cookie{}
	s.client.Cookies = []*http.Cookie{}
}

func (s *Student) NewRequest() *resty.Request {
	return s.client.R()
}

func (s *Student) GetWithIdentifier(url string, queryParams map[string]string) (*html.Node, error) {
	request := s.NewRequest().SetHeader("Referer", constants.YJSYReferer).SetHeader("User-Agent", constants.UserAgent)
	if queryParams != nil {
		request = request.SetQueryParams(queryParams)
	}
	// 会话过期：会直接重定向，但我们禁用了重定向，所以会有error
	resp, err := request.Get(url)
	if err != nil {
		return nil, errno.CookieError
	}

	// 原文是一个Alert：您没有权限进入本页或当前登录用户已过期！\n请重新登录或与管理员联系！
	if strings.Contains(string(resp.Body()), "当前登录用户已过期") {
		return nil, errno.CookieError
	}

	// 系统发生错误，原文是一个页面，包含以下文字：系统发生错误，该信息已被系统记录，请稍后重试或与管理员联系。
	if strings.Contains(string(resp.Body()), "系统发生错误") {
		return nil, errno.SystemError.WithMessage("教务系统内部错误")
	}

	return htmlquery.Parse(bytes.NewReader(resp.Body()))
}

func (s *Student) PostWithIdentifier(url string, formData map[string]string) (*html.Node, error) {
	resp, err := s.NewRequest().SetHeader("Referer", constants.YJSYReferer).SetHeader("User-Agent", constants.UserAgent).SetFormData(formData).Post(url)

	// 会话过期：会直接重定向，但我们禁用了重定向，所以会有error
	if err != nil {
		return nil, errno.CookieError
	}

	// id 或 cookie 缺失或者解析错误 TODO: 判断条件有点简陋
	if strings.Contains(string(resp.Body()), "当前登录用户已过期") {
		return nil, errno.CookieError
	}

	// 系统发生错误，原文是一个页面，包含以下文字：系统发生错误，该信息已被系统记录，请稍后重试或与管理员联系。
	if strings.Contains(string(resp.Body()), "系统发生错误") {
		return nil, errno.SystemError.WithMessage("教务系统内部错误")
	}

	return htmlquery.Parse(strings.NewReader(strings.TrimSpace(string(resp.Body()))))
}
