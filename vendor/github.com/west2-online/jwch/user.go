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

package jwch

import (
	"encoding/json"
	"net/http"
	"regexp"

	"github.com/west2-online/jwch/constants"
	"github.com/west2-online/jwch/errno"
	"github.com/west2-online/jwch/utils"

	"github.com/antchfx/htmlquery"
)

// SSO登录返回
type ssoLoginResponse struct {
	Code int    `json:"code"` // 状态码
	Info string `json:"info"` // 返回消息
	// Data interface{} `json:"data"`
}

type verifyCodeResponse struct {
	Message string `json:"message"`
}

// Login 模拟教务处登录/刷新Session
func (s *Student) Login() error {
	// 清除cookie
	s.ClearLoginData()

	code := verifyCodeResponse{}
	loginResp := ssoLoginResponse{}
	passMD5 := utils.Md5Hash(s.Password, 16)
	// 获取验证码图片
	resp, err := s.NewRequest().Get("https://jwcjwxt2.fzu.edu.cn:82/plus/verifycode.asp")
	if err != nil {
		return errno.HTTPQueryError.WithErr(err)
	}

	// 请求西二服务器，自动识别验证码
	resp, err = s.NewRequest().SetFormData(map[string]string{
		"validateCode": utils.Base64EncodeHTTPImage(resp.Body()),
	}).Post(constants.AutoCaptchaVerifyURL)
	if err != nil {
		return errno.HTTPQueryError.WithMessage("automatic code identification failed")
	}

	err = json.Unmarshal(resp.Body(), &code)
	if err != nil {
		return errno.HTTPQueryError.WithErr(err)
	}

	// 登录验证
	_, err = s.NewRequest().SetHeaders(map[string]string{
		"Referer": "https://jwch.fzu.edu.cn",
		"Origin":  "https://jwch.fzu.edu.cn",
	}).SetFormData(map[string]string{
		"Verifycode": code.Message,
		"muser":      s.ID,
		"passwd":     passMD5,
	}).Post("https://jwcjwxt2.fzu.edu.cn:82/logincheck.asp")

	// 由于禁用了302，这里正常情况下会返回一个错误，跳转链接中包含了我们要的全部信息
	if err == nil {
		return errno.LoginCheckFailedError
	}

	// 获取token，第一个是匹配的全部字符，第二个是我们需要的
	token := regexp.MustCompile(`token=(.*?)&`).FindStringSubmatch(err.Error())
	if len(token) < 1 {
		return errno.LoginCheckFailedError
	}

	id := regexp.MustCompile(`id=(.*?)&`).FindStringSubmatch(err.Error())[1]
	num := regexp.MustCompile(`num=(.*?)&`).FindStringSubmatch(err.Error())[1]

	// SSO登录
	resp, err = s.NewRequest().SetHeaders(map[string]string{
		"X-Requested-With": "XMLHttpRequest",
	}).SetFormData(map[string]string{
		"token": token[1],
	}).Post(constants.SSOLoginURL)
	if err != nil {
		return errno.HTTPQueryError.WithErr(err)
	}

	err = json.Unmarshal(resp.Body(), &loginResp)
	if err != nil {
		return errno.HTTPQueryError.WithErr(err)
	}

	// 获取account不存在是400，登录成功是200
	if loginResp.Code != 200 {
		return errno.SSOLoginFailedError
	}

	// 获取cookies
	resp, err = s.NewRequest().SetHeaders(map[string]string{
		"Referer": constants.JwchReferer,
		"Origin":  constants.JwchOrigin,
	}).SetQueryParams(map[string]string{
		"id":       id,
		"num":      num,
		"ssourl":   "https://jwcjwxt2.fzu.edu.cn",
		"hosturl":  "https://jwcjwxt2.fzu.edu.cn:81",
		"ssologin": "",
	}).Get("https://jwcjwxt2.fzu.edu.cn:81/loginchk_xs.aspx")

	// 保存这部分Cookie，这部分Cookie是用来后续鉴权的[ASP.NET_SessionId]
	s.SetCookies(resp.RawResponse.Cookies())

	// 这里是err == nil 因为禁止了重定向，正常登录是会出现异常的
	if err == nil {
		return errno.CookieError
	}

	data := regexp.MustCompile(`id=(.*?)&`).FindStringSubmatch(err.Error())

	if len(data) < 1 {
		return errno.CookieError
	}

	s.SetIdentifier(data[1])

	return nil
}

// GetIdentifierAndCookies 方面服务端进行测试设置的接口
func (s *Student) GetIdentifierAndCookies() (string, []*http.Cookie, error) {
	err := s.CheckSession()
	if err != nil {
		if err := s.Login(); err != nil {
			return "", nil, err
		}
	}
	return s.Identifier, s.client.Cookies, nil
}

// CheckSession returns not nil if SessionExpired or AccountConflict
func (s *Student) CheckSession() error {
	// 逻辑: 如果session没用，我们会返回一个302定向到https://jwcjwxt2.fzu.edu.cn:82/error.asp?id=300，但是我们禁用了重定向，意味着这里HTTP会抛出异常
	// 旧版处理过程： 查询Body中是否含有[当前用户]这四个字

	// 检查过期
	resp, err := s.GetWithIdentifier(constants.UserInfoURL)
	if err != nil {
		return err
	}

	// 检查串号
	res := htmlquery.FindOne(resp, `//*[@id="ContentPlaceHolder1_LB_xh"]`)

	if res == nil {
		return errno.CookieError
	}

	if htmlquery.OutputHTML(res, false) != s.ID {
		return errno.AccountConflictError
	}

	return nil
}

// GetInfo 获取学生个人信息
func (s *Student) GetInfo() (resp *StudentDetail, err error) {
	res, err := s.GetWithIdentifier(constants.UserInfoURL)
	if err != nil {
		return nil, err
	}

	resp = &StudentDetail{
		Name:             safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_xm"]`),
		Birthday:         safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_csrq"]`),
		Sex:              safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_xb"]`),
		Phone:            safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_lxdh"]`),
		Email:            safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_email"]`),
		College:          safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_xymc"]`),
		Grade:            safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_nj"]`),
		StatusChanges:    safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_xjxx"]`),
		Major:            safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_zymc"]`),
		Counselor:        safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_zdy"]`),
		ExamineeCategory: safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_kslb"]`),
		Nationality:      safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_mz"]`),
		Country:          safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_gb"]`),
		PoliticalStatus:  safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_zzmm"]`),
		Source:           safeExtractHTMLFirst(res, `//*[@id="ContentPlaceHolder1_LB_xssy"]`),
	}

	return resp, nil
}
