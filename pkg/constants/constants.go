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

	// db table name
	TemplateServiceTableName = "template"

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
)
