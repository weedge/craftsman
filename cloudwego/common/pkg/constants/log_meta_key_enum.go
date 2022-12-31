package constants

// biz meta key
const (
	ProjectNameKey = "projectName"
	AppNameKey     = "appName"
	EnvKey         = "env"
	LocalIpKey     = "localIp"
	OutBoundIpKey  = "outBoundIp"
)

// http header key
const (
	SessionId            = "Session-Id"
	RequestSessionId     = "Request-Session-Id"
	LastRequestSessionId = "Last-Request-Session-Id"
	XRequestId           = "X-Request-ID"
	CurrentUserId        = "Current-User-Id"

	XB3ZipKinOpenTraceId = "x-b3-traceid" // open tracing id for B3 ZipKin https://github.com/openzipkin/b3-propagation
)

// mq property key
const (
	MqTraceIdKey   = "mq_trace_id"
	MqSpanIdKey    = "mq_sapn_id"
	MqTraceSpanKey = "mq_trace_span"
)

// log type
const (
	LogTypeKey        = "LOGTYPE"
	LogTypeServer     = "server"   // biz server logic log
	LogTypeAccess     = "access"   // access log
	LogTypeMqConsumer = "consumer" // mq consumer
	LogTypeModule     = "module"   // dependent resource module (db,cache)
	LogTypeHttpReq    = "httpreq"  // dependent http resource access lay request/response
	LogTypeRpc        = "rpc"      // dependent rpc resource access lay
	LogTypeJob        = "job"      // cron job
)

// log type -> key  related lib to define
const ()

// access resource proto
const (
	ProtoKey          = "PROTO"
	ProtoKeyHttp      = "http"
	ProtoKeyGRPC      = "grpc"
	ProtoKeyThriftRPC = "thrift"
	ProtoKeyRedis     = "redis"
	ProtoKeyMongo     = "mongo"
	ProtoKeyPgSql     = "pgsql"
	ProtoKeyMySql     = "mysql"
)

// cost
const (
	CostKey             = "cost"
	RequestStartTimeKey = "requestStartTime"
	RequestEndTimeKey   = "requestEndTime"
)
