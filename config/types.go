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

type service struct {
	Name     string
	AddrList []string
	LB       bool `mapstructure:"load-balance"`
}

/*
for android
用于构造又拍云的上传参数
*/
type url struct {
	Password     string
	Bucket       string
	Operator     string
	Pass         string
	TokenTimeout int64 `mapstructure:"token-timeout"`
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
* struct upyun 又拍云配置
* @Bucket: 存储桶
* @Opearator: 操作员
* @Password: 密码
* @TokenSecret: 对应又拍云里的SecretAccessKey
* @TokenTimeout: Token过期时间
* @UssDomain: 域名
* @UnCheckedDir: 上传目录
 */
type upyun struct {
	Bucket         string
	Operator       string
	Password       string
	TokenSecret    string `mapstructure:"token-secret"`
	TokenTimeout   int64  `mapstructure:"token-timeout"`
	UssDomain      string `mapstructure:"uss-domain"`
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

type umeng struct {
	Android AndroidUmeng `mapstructure:"android"`
	IOS     IOSUmeng     `mapstructure:"ios"`
}

type vendor struct {
	ExamNotifications        string `mapstructure:"ExamNotifications"`
	ExamResultsNotifications string `mapstructure:"ExamResultsNotifications"`
	JwchNotice               string `mapstructure:"JwchNotice"`
}

type vendors struct {
	Xiaomi vendor `mapstructure:"xiaomi"`
	Huawei vendor `mapstructure:"huawei"`
}

type mcp struct {
	Name    string `mapstructure:"name"`
	Version string `mapstructure:"version"`
}
type friend struct {
	MaxNum int64 `mapstructure:"max-nums"`
}

type config struct {
	Server        server
	MCP           mcp `mapstructure:"mcp"`
	Snowflake     snowflake
	MySQL         mySQL
	Jaeger        jaeger
	Etcd          etcd
	RabbitMQ      rabbitMQ
	Redis         redis
	OSS           oss
	Elasticsearch elasticsearch
	Kafka         kafka
	DefaultUser   defaultUser
	UpYuns        map[string]upyun
	Umeng         umeng
	Url           url
	Vendors       vendors
	Friend        friend
}
