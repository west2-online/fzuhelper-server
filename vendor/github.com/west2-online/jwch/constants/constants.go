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

package constants

const (
	ClassroomQueryURL   = "https://jwcjwxt2.fzu.edu.cn:81/kkgl/kbcx/kbcx_kjs.aspx"
	CourseURL           = "https://jwcjwxt2.fzu.edu.cn:81/student/xkjg/wdxk/xkjg_list.aspx"
	MarksQueryURL       = "https://jwcjwxt2.fzu.edu.cn:81/student/xyzk/cjyl/score_sheet.aspx"
	CETQueryURL         = "https://jwcjwxt2.fzu.edu.cn:81/student/glbm/cet/cet_cszt.aspx"
	JSQueryURL          = "https://jwcjwxt2.fzu.edu.cn:81/student/glbm/computer/jsj_cszt.aspx"
	UserInfoURL         = "https://jwcjwxt2.fzu.edu.cn:81/jcxx/xsxx/StudentInformation.aspx"
	SSOLoginURL         = "https://jwcjwxt2.fzu.edu.cn/Sfrz/SSOLogin"
	SchoolCalendarURL   = "https://jwcjwxt2.fzu.edu.cn:82/xl.asp"
	CreditQueryURL      = "https://jwcjwxt2.fzu.edu.cn:81/student/xyzk/xftj/CreditStatistics.aspx"
	GPAQueryURL         = "https://jwcjwxt2.fzu.edu.cn:81/student/xyzk/jdpm/GPA_sheet.aspx"
	VerifyCodeURL       = "https://jwcjwxt2.fzu.edu.cn:82/plus/verifycode.asp"
	ExamRoomQueryURL    = "https://jwcjwxt2.fzu.edu.cn:81/student/xkjg/examination/exam_list.aspx"
	NoticeInfoQueryURL  = "https://jwch.fzu.edu.cn/jxtz.htm"
	JwchNoticeURLPrefix = "https://jwch.fzu.edu.cn/"
	CultivatePlanURL    = "https://jwcjwxt2.fzu.edu.cn:81/pyfa/pyjh/pyjh_list.aspx"
	JwchLocateDateUrl   = "https://jwcjwxt2.fzu.edu.cn:82/week.asp"
	LectureURL          = "https://jwcjwxt2.fzu.edu.cn:81/student/glbm/lecture/jxjt_cszt.aspx"
	JwchPingYiUrl       = "student/jscp/TeaList.aspx"

	UserAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/139.0.0.0 Safari/537.36"

	JwchPrefix  = "https://jwcjwxt2.fzu.edu.cn:81"
	JwchReferer = "https://jwcjwxt2.fzu.edu.cn:82/"
	JwchOrigin  = "https://jwcjwxt2.fzu.edu.cn:82/"

	AutoCaptchaVerifyURL = "https://statistics.fzuhelper.w2fzu.com/api/login/validateCode?validateCode"

	// 青果网络代理相关常量
	QingGuoTunnelURL = "https://longterm.proxy.qg.net/query" // 青果网络隧道地址获取接口
)

var BuildingArray = []string{"公共教学楼东1", "公共教学楼东2", "公共教学楼东3", "公共教学楼文科楼", "公共教学楼西1", "公共教学楼西2", "公共教学楼西3", "公共教学楼中楼"}
