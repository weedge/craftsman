package logutils

import (
	"os"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	kitexZap "github.com/weedge/craftsman/cloudwego/kitex-contrib/obs-opentelemetry/logging/zap"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func ZapEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.CallerKey = "file"
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.999999"))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)

	return encoder
}

func ZapWriteSyncer() (ws zapcore.WriteSyncer) {
	ws = zapcore.AddSync(os.Stdout)

	return
}

func ZapLogLevel(level Level) (lvl zap.AtomicLevel) {
	lvl = zap.NewAtomicLevelAt(level.ZapLogLevel())

	return
}

func ZapOptions(meta map[string]interface{}) (opts []zap.Option) {
	var fields []zap.Field
	for key, val := range meta {
		fields = append(fields, zap.Any(key, val))
	}
	opts = append(opts, zap.Fields(fields...))
	opts = append(opts, zap.AddCaller())
	opts = append(opts, zap.Development())

	return
}

func NewkitexZapLogger(level Level, meta map[string]interface{}) (logger klog.FullLogger) {
	logger = kitexZap.NewLogger(
		kitexZap.WithCoreEnc(ZapEncoder()),
		kitexZap.WithCoreWs(ZapWriteSyncer()),
		kitexZap.WithCoreLevel(ZapLogLevel(level)),
		kitexZap.WithZapOptions(ZapOptions(meta)...),
	)

	return
}
