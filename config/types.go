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

package config

type server struct {
	Secret   string `mapstructure:"private-key"`
	Version  string
	Name     string
	LogLevel string `mapstructure:"log-level"`
}

type snowflake struct {
	WorkerID      int64 `mapstructure:"worker-id"`
	DatancenterID int64 `mapstructure:"datancenter-id"`
}

type admin struct {
	Secret string `mapstructure:"secret"`
}

type service struct {
	Name     string
	AddrList []string
	LB       bool `mapstructure:"load-balance"`
}

/*
for android
用于构造 COS PostObject 的表单上传参数(服务端签名后下发给客户端)
*/
type cosUpload struct {
	Bucket       string
	Region       string
	SecretID     string `mapstructure:"secret-id"`
	SecretKey    string `mapstructure:"secret-key"`
	TokenTimeout int64  `mapstructure:"token-timeout"`
	Path         string
}

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type jaeger struct {
	Addr string
}

type otel struct {
	Endpoint string
}

type etcd struct {
	Addr string
}

type rabbitMQ struct {
	Addr     string
	Username string
	Password string
}

type redis struct {
	Addr     string
	Password string
}

type oss struct {
	Endpoint        string
	AccessKeyID     string `mapstructure:"accessKey-id"`
	AccessKeySecret string `mapstructure:"accessKey-secret"`
	BucketName      string
	MainDirectory   string `mapstructure:"main-directory"`
}

type elasticsearch struct {
	Addr string
	Host string
}

type kafka struct {
	Address  string
	Network  string
	User     string
	Password string
}

type defaultUser struct {
	Account  string `mapstructure:"account"`
	Password string `mapstructure:"password"`
}

/*
* struct cos 腾讯云 COS 存储配置
* @Bucket: 存储桶名(需带 appid 后缀,如 fzuhelper-paper-cos-125000000)
* @Region: 桶所在地域,如 ap-guangzhou
* @SecretID: 子账号密钥 ID(每个桶独立子账号,禁止使用主账号密钥)
* @SecretKey: 子账号密钥 Key
* @TokenSecret: EO Token 鉴权 PrivateKey(历史卷防盗链,留空则下载链接不签名)
* @TokenTimeout: 签名过期时间(秒)
* @DownloadDomain: 用户侧下载域名(EO 加速域名)
* @Path: 对象存储根路径
* @AvatarPath: 贡献者头像存储路径
 */
type cos struct {
	Bucket         string
	Region         string
	SecretID       string `mapstructure:"secret-id"`
	SecretKey      string `mapstructure:"secret-key"`
	TokenSecret    string `mapstructure:"token-secret"`
	TokenTimeout   int64  `mapstructure:"token-timeout"`
	DownloadDomain string `mapstructure:"download-domain"`
	Path           string
	AvatarPath     string `mapstructure:"avatar-path"`
}

type AndroidUmeng struct {
	AppKey          string `mapstructure:"app_key"`
	AppMasterSecret string `mapstructure:"app_master_secret"`
}

type IOSUmeng struct {
	AppKey          string `mapstructure:"app_key"`
	AppMasterSecret string `mapstructure:"app_master_secret"`
}

type HarmonyUmeng struct {
	AppKey          string `mapstructure:"app_key"`
	AppMasterSecret string `mapstructure:"app_master_secret"`
}

type umeng struct {
	Android AndroidUmeng `mapstructure:"android"`
	IOS     IOSUmeng     `mapstructure:"ios"`
	Harmony HarmonyUmeng `mapstructure:"harmony"`
}

type oppo struct {
	ChannelID          string `mapstructure:"channel_id"`
	Category           string `mapstructure:"category"`
	NotifyLevel        string `mapstructure:"notify_level"`
	PrivateMsgTemplate struct {
		PrivateMsgTemplateID string `mapstructure:"private_msg_template_id"`
	} `mapstructure:"private_msg_template"`
}

type huawei struct {
	ChannelImportance string `mapstructure:"channel_importance"`
	ChannelCategory   string `mapstructure:"channel_category"`
}

type harmonyVendor struct {
	ChannelCategory string `mapstructure:"channel_category"`
}

type localProperties struct {
	ChannelID   string `mapstructure:"channel_id"`
	ChannelName string `mapstructure:"channel_name"`
}

type XiaomiNoticeTemplate struct {
	ChannelID  string `mapstructure:"channel_id"`
	TemplateID string `mapstructure:"template_id"`
}

// XiaomiNotice 小米推送模板配置，key 为推送类型（score/exam/teaching），对应 pkg/constants 中的 UmengPushType*
type XiaomiNotice map[string]XiaomiNoticeTemplate

type vendors struct {
	ChannelActivity string          `mapstructure:"channel_activity"`
	XiaomiNotice    XiaomiNotice    `mapstructure:"xiaomi_notice"`
	VivoCategory    string          `mapstructure:"vivo_category"`
	Oppo            oppo            `mapstructure:"oppo"`
	Huawei          huawei          `mapstructure:"huawei"`
	Harmony         harmonyVendor   `mapstructure:"harmony"`
	LocalProperties localProperties `mapstructure:"local_properties"`
}

type mcp struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}

type ai struct {
	Key      string `mapstructure:"key"`
	Endpoint string `mapstructure:"endpoint"`
}

type friend struct {
	MaxNum int64 `mapstructure:"max-nums"`
}

type signedLocationApiUrl struct {
	Endpoint   string `mapstructure:"endpoint"`
	Enabled    bool   `mapstructure:"enabled"`
	DisableMsg string `mapstructure:"disable_msg"`
}

type apiMonitorConfig struct {
	Enabled              bool     `mapstructure:"enabled"`
	WindowSeconds        int64    `mapstructure:"window-seconds"`
	CheckIntervalSeconds int64    `mapstructure:"check-interval-seconds"`
	ErrorRateThreshold   float64  `mapstructure:"error-rate-threshold"`
	MinRequests          int64    `mapstructure:"min-requests"`
	AlertCooldownSeconds int64    `mapstructure:"alert-cooldown-seconds"`
	RouteBlacklist       []string `mapstructure:"route-blacklist"`
	CodeBlacklist        []int64  `mapstructure:"code-blacklist"`
}

type config struct {
	Server               server
	MCP                  mcp `mapstructure:"mcp"`
	Admin                admin
	AI                   ai
	Snowflake            snowflake
	MySQL                mySQL
	Jaeger               jaeger
	Otel                 otel
	Etcd                 etcd
	RabbitMQ             rabbitMQ
	Redis                redis
	OSS                  oss
	Elasticsearch        elasticsearch
	Kafka                kafka
	DefaultUser          defaultUser
	Coss                 map[string]cos
	CosUpload            cosUpload `mapstructure:"cos-upload"`
	Umeng                umeng
	Vendors              vendors
	Friend               friend
	SignedLocationApiUrl signedLocationApiUrl `mapstructure:"signed_location_api_url"`
	APIMonitor           apiMonitorConfig     `mapstructure:"api-monitor"`
}
