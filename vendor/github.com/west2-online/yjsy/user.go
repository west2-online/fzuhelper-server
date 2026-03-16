package yjsy

import (
	"encoding/base64"
	"net/http"
	"strings"

	"github.com/west2-online/yjsy/constants"
	"github.com/west2-online/yjsy/errno"
)

// Login 模拟教务处登录/刷新Session
func (s *Student) Login() error {
	// 清除cookie
	s.ClearLoginData()

	// 登录验证
	res, err := s.NewRequest().SetHeaders(map[string]string{
		"Referer":    constants.YJSYReferer,
		"Origin":     constants.YJSYOrigin,
		"User-Agent": constants.UserAgent,
	}).SetFormData(map[string]string{
		"muser":  s.ID,
		"passwd": base64.StdEncoding.EncodeToString([]byte(s.Password)),
	}).Post(constants.LoginURL)

	// 由于禁用了302，这里正常情况下会返回一个错误,没有错误就说明登陆失败了
	// 注意，研究生登陆失败超过5次会被封半个小时
	if err == nil {
		return errno.LoginCheckFailedError
	}

	// 解析 Cookie
	cookies := res.Header().Values("Set-Cookie")
	if cookies != nil {
		var parsedCookies []*http.Cookie
		// cookie有两个.ASPXAUTH ,我们只需要第一个 path 为 / 的cookie
		for _, cookie := range cookies {
			cookie = strings.TrimSpace(cookie)

			parts := strings.Split(cookie, ";")

			cookieParts := strings.Split(parts[0], "=")

			var path string
			for _, part := range parts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "path=") {
					path = strings.TrimPrefix(part, "path=")
					break
				}
			}

			// 如果 cookie 解析有效，存储到 parsedCookies 中
			if len(cookieParts) == 2 && path == "/" {
				parsedCookies = append(parsedCookies, &http.Cookie{
					Name:  cookieParts[0],
					Value: cookieParts[1],
					Path:  path,
				})
			}
		}
		s.SetCookies(parsedCookies)
	} else {
		return errno.CookieError
	}

	return nil

}

// GetCookies
func (s *Student) GetCookies() ([]*http.Cookie, error) {
	return s.client.Cookies, nil
}

// CheckSession returns not nil if SessionExpired or AccountConflict
func (s *Student) CheckSession() error {

	// 进行一次成绩查询根据是否正常返回来判断是否过期
	// 如果会话过期，GetWithIdentifier 会直接返回 cookie error
	_, err := s.GetWithIdentifier(constants.MarksQueryURL, nil)
	if err != nil {
		return err
	}

	return nil
}
