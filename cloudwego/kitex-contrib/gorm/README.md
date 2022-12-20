### Features:
1. impl gorm logger, add IkvLogger interface for gorm trace 
```golang
type IkvLogger interface {
	// CtxKVLog  kvs must be kv pairs k,v , k,v ...
	CtxKVLog(ctx context.Context, level int, format string, kvs ...interface{})
}
```

2. add IkvLogger impl DefaultLogger for gorm trace