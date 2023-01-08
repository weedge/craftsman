package tracer

import (
	"context"
	"strconv"

	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/weedge/craftsman/cloudwego/common/pkg/constants"
	"go.opentelemetry.io/otel/trace"
)

// MsgTracePayload for msg otel tracing
type MsgTracePayload struct {
	TraceID    string `json:"TraceID"`
	SpanID     string `json:"SpanID"`
	TraceFlags string `json:"TraceFlags"`
	TraceState string `json:"TraceState"`
	Remote     bool   `json:"Remote"`
}

// ContextWithOtelTraceSpanContextFromMsg
func ContextWithOtelTraceSpanContextFromMsg(ctx context.Context, msg *primitive.MessageExt) context.Context {
	spanStr := msg.GetProperty(constants.MqTraceSpanKey)

	if len(spanStr) > 0 {
		spanConf := trace.SpanContextConfig{}
		msgTracePayload := MsgTracePayload{}
		err := sonic.Unmarshal(utils.StringToSliceByte(spanStr), &msgTracePayload)
		if err != nil {
			klog.CtxErrorf(ctx, "sonic.Unmarshal spanStr:%s err:%s", spanStr, err.Error())
		}

		spanConf.TraceID, err = trace.TraceIDFromHex(msgTracePayload.TraceID)
		if err != nil {
			klog.CtxErrorf(ctx, "trace.TraceIDFromHex(%s) err:%s", msgTracePayload.TraceID, err.Error())
		}
		spanConf.SpanID, err = trace.SpanIDFromHex(msgTracePayload.SpanID)
		if err != nil {
			klog.CtxErrorf(ctx, "trace.SpanIDFromHex(%s) err:%s", msgTracePayload.SpanID, err.Error())
		}
		// otherwise base, int bitsize
		traceFlags, _ := strconv.ParseInt(msgTracePayload.TraceFlags, 10, 0)
		spanConf.TraceFlags = trace.TraceFlags(traceFlags)
		spanConf.Remote = msgTracePayload.Remote
		//klog.CtxDebugf(ctx, "spanStr:%s msgTracePayload:%+v spanConf:%+v", spanStr, msgTracePayload, spanConf)
		ctx = trace.ContextWithSpanContext(ctx, trace.NewSpanContext(spanConf))
	}

	return ctx
}
