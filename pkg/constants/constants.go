package constants

import "time"

const (
	// auth
	JWTValue = "MTAxNTkwMTg1Mw=="
	StartID  = 10000

	// RPC
	MuxConnection  = 1
	RPCTimeout     = 3 * time.Second
	ConnectTimeout = 50 * time.Millisecond

	// service name
	TemplateServiceName = "template"
	ClassroomService    = "classroom"

	// db table name
	TemplateServiceTableName = "template"

	// redis
	RedisDBEmptyRoom = 1

	// snowflake
	SnowflakeWorkerID     = 0
	SnowflakeDatacenterID = 0

	// limit
	MaxConnections  = 1000
	MaxQPS          = 100
	MaxVideoSize    = 300000
	MaxListLength   = 100
	MaxIdleConns    = 10
	MaxGoroutines   = 10
	MaxOpenConns    = 100
	ConnMaxLifetime = 10 * time.Second

	// default account and password
	DefaultAccount  = "082100170"
	DefaultPassword = "Zhuyinfan815"
)

var CampusArray = []string{"旗山校区", "鼓浪屿校区", "集美校区", "铜盘校区", "怡山校区", "晋江校区", "泉港校区"}
