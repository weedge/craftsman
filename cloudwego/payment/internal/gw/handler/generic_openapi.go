package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/common"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/injectors"
	"go.opentelemetry.io/otel/trace"
)

// OpenApiHandle handle the request with the query path of prefix `/openapi`.
// use kitex generic call (http thrift generic client)
func OpenApiHandle(ctx context.Context, c *app.RequestContext) {
	projectName := c.Param("project")
	svcName := c.Param("svc")
	version := c.Param("version")
	methodName := c.Param("method")
	httpMethodName := utils.SliceByteToString(c.Method())

	svcKey := injectors.GetApiSvcKey(projectName, svcName, version)
	genericeCli, ok := gSvcMap[svcKey]
	if !ok {
		klog.CtxWarnf(ctx, "svcName:%s don't register", svcKey)
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayMethodNotFound),
			ErrMsg:  common.Err_GateWayMethodNotFound.String(),
			Extra:   GetTraceSpanExtra(ctx),
		})
		return
	}

	if genericeOpts, ok := gSvcOptsMap[svcKey]; ok {
		for _, closedMethod := range genericeOpts.ClosedMethods {
			if closedMethod.SvcMethod == methodName &&
				strings.EqualFold(closedMethod.HttpMethod, httpMethodName) {
				klog.CtxWarnf(ctx, "%s %s %s is closed", httpMethodName, svcKey, methodName)
				c.JSON(http.StatusOK, &base.BaseResp{
					ErrCode: int64(common.Err_GateWayMethodNotFound),
					ErrMsg:  common.Err_GateWayMethodNotFound.String(),
					Extra:   GetTraceSpanExtra(ctx),
				})
				return
			}
		}
	}

	req, err := http.NewRequest(httpMethodName, "", bytes.NewBuffer(c.Request.BodyBytes()))
	if err != nil {
		klog.CtxWarnf(ctx, "svcName:%s http.NewRequest err:%s", svcName, err.Error())
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayServerHandlerError),
			ErrMsg:  common.Err_GateWayServerHandlerError.String(),
			Extra:   GetTraceSpanExtra(ctx),
		})
		return
	}

	req.URL.Path = fmt.Sprintf("/%s/%s/%s/%s", projectName, svcName, version, methodName)
	req.URL.RawQuery = utils.SliceByteToString(c.Request.QueryString())
	if genericeOpts, ok := gSvcOptsMap[svcKey]; ok {
		for _, key := range genericeOpts.HeaderKeys {
			val := c.Request.Header.Get(key)
			if len(val) > 0 {
				req.Header.Set(key, val)
			}
		}
	}

	innerReq, err := generic.FromHTTPRequest(req)
	if err != nil {
		klog.CtxWarnf(ctx, "req:%+v generic.FromHTTPRequest err:%s", req, err.Error())
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayServerHandlerError),
			ErrMsg:  common.Err_GateWayServerHandlerError.String(),
			Extra:   GetTraceSpanExtra(ctx),
		})
		return
	}

	resp, err := genericeCli.GenericCall(ctx, "", innerReq)
	if err != nil {
		klog.CtxErrorf(ctx, "req:%+v genericeCli.GenericCall err:%s", innerReq, err.Error())
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayServerInnerCallError),
			ErrMsg:  common.Err_GateWayServerInnerCallError.String(),
			Extra:   GetTraceSpanExtra(ctx),
		})
		return
	}

	realResp, ok := resp.(*generic.HTTPResponse)
	if !ok {
		klog.CtxWarnf(ctx, "resp:%+v not generic.HTTPResponse ", resp)
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayServerHandlerError),
			ErrMsg:  common.Err_GateWayServerHandlerError.String(),
			Extra:   GetTraceSpanExtra(ctx),
		})
		return
	}

	klog.CtxInfof(ctx, "genericeCli.GenericCall ok, req:%+v resp:%+v", req, resp)
	if baseResp, ok := realResp.Body["baseResp"]; ok {
		if bResp, ok := baseResp.(map[string]interface{}); ok {
			realResp.Body["errCode"] = bResp["errCode"]
			realResp.Body["errMsg"] = bResp["errMsg"]
			realResp.Body["extra"] = GetTraceSpanExtra(ctx)
			if _, ok := bResp["extra"]; ok {
				realResp.Body["extra"] = MergeExtra(realResp.Body["extra"], bResp["extra"])
			}
		}
		delete(realResp.Body, "baseResp")
	}

	c.JSON(http.StatusOK, realResp.Body)
}

func GetTraceSpanExtra(ctx context.Context) map[string]string {
	span := trace.SpanFromContext(ctx)
	return map[string]string{
		"traceId": span.SpanContext().TraceID().String(),
		//"spanId":     span.SpanContext().SpanID().String(),
		//"traceFlags": span.SpanContext().TraceFlags().String(),
	}
}

func MergeExtra(src interface{}, dest interface{}) interface{} {
	s, ok := src.(map[string]interface{})
	if !ok {
		return src
	}
	d, ok := dest.(map[string]interface{})
	if !ok {
		return src
	}

	for key, item := range d {
		s[key] = item
	}

	return s
}
