FROM: https://github.com/kitex-contrib/obs-opentelemetry/tree/main/logging/zap


CHANGE:
1. add zap logger for zap fields, log kv pairs
```golang
func (l *Logger) CtxKVLog(ctx context.Context, level klog.Level, format string, kvs ...interface{})
