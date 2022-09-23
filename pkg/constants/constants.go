package constants

import "time"

const (
	GrpcPort       = "GRPC_PORT"
	HttpPort       = "HTTP_PORT"
	ConfigPath     = "CONFIG_PATH"
	JaegerHostPort = "JAEGER_HOST"
	RedisAddr      = "REDIS_ADDR"
	MongoDbURI     = "MONGO_URI"
	PostgresqlHost = "POSTGRES_HOST"
	PostgresqlPort = "POSTGRES_PORT"
	Yaml           = "yaml"
	Json           = "json"

	GRPC     = "GRPC"
	METHOD   = "METHOD"
	NAME     = "NAME"
	METADATA = "METADATA"
	REQUEST  = "REQUEST"
	REPLY    = "REPLY"
	TIME     = "TIME"

	Topic     = "topic"
	Partition = "partition"
	Message   = "message"
	WorkerID  = "workerID"
	Offset    = "offset"
	Time      = "time"
)

const (
	ErrBadRequestTitle          = "Bad Request"
	ErrConflictTitle            = "Conflict Error"
	ErrNotFoundTitle            = "Not Found"
	ErrUnauthorizedTitle        = "Unauthorized"
	ErrForbiddenTitle           = "Forbidden"
	ErrRequestTimeoutTitle      = "Request Timeout"
	ErrInternalServerErrorTitle = "Internal Server Error"
	ErrDomainTitle              = "Domain Model Error"
	ErrApplicationTitle         = "Application Service Error"
	ErrApiTitle                 = "Api Error"
)

const (
	MaxHeaderBytes       = 1 << 20
	StackSize            = 1 << 10 // 1 KB
	BodyLimit            = "2M"
	ReadTimeout          = 15 * time.Second
	WriteTimeout         = 15 * time.Second
	GzipLevel            = 5
	WaitShotDownDuration = 3 * time.Second
	Dev                  = "development"
	Test                 = "test"
	Production           = "production"
)
