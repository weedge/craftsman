package handler

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/generic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/utils"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/base"
	"github.com/weedge/craftsman/cloudwego/common/kitex_gen/common"
)

type requestBodyParams struct {
	Method    string `form:"method,required" json:"method"`
	BizParams string `form:"bizParams,required" json:"bizParams"`
}

// OpenApiHandle handle the request with the query path of prefix `/openapi`.
// use kitex generic call (http thrift generic client)
func OpenApiHandle(ctx context.Context, c *app.RequestContext) {
	svcName := c.Param("svc")
	genericeCli, ok := gSvcMap[svcName]
	if !ok {
		klog.CtxWarnf(ctx, "svcName:%s don't register", svcName)
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayMethodNotFound),
			ErrMsg:  common.Err_GateWayMethodNotFound.String(),
			Extra:   nil,
		})
		return
	}

	params := &requestBodyParams{}
	if err := c.BindAndValidate(&params); err != nil {
		klog.CtxWarnf(ctx, "svcName:%s BindAndValidate Params err:%s", svcName, err.Error())
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayBadRequest),
			ErrMsg:  common.Err_GateWayBadRequest.String(),
			Extra:   nil,
		})
		return
	}

	req, err := http.NewRequest(utils.SliceByteToString(c.Method()), "", bytes.NewBuffer([]byte(params.BizParams)))
	if err != nil {
		klog.CtxWarnf(ctx, "svcName:%s http.NewRequest err:%s", svcName, err.Error())
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayServerHandlerError),
			ErrMsg:  common.Err_GateWayServerHandlerError.String(),
			Extra:   nil,
		})
		return
	}

	req.URL.Path = fmt.Sprintf("/%s/%s", svcName, params.Method)
	innerReq, err := generic.FromHTTPRequest(req)
	if err != nil {
		klog.CtxWarnf(ctx, "req:%+v generic.FromHTTPRequest err:%s", req, err.Error())
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayServerHandlerError),
			ErrMsg:  common.Err_GateWayServerHandlerError.String(),
			Extra:   nil,
		})
		return
	}

	resp, err := genericeCli.GenericCall(ctx, "", innerReq)
	if err != nil {
		klog.CtxErrorf(ctx, "req:%+v genericeCli.GenericCall err:%s", innerReq, err.Error())
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayServerInnerCallError),
			ErrMsg:  common.Err_GateWayServerInnerCallError.String(),
			Extra:   nil,
		})
		return
	}

	realResp, ok := resp.(*generic.HTTPResponse)
	if !ok {
		klog.CtxWarnf(ctx, "resp:%+v not generic.HTTPResponse ", resp)
		c.JSON(http.StatusOK, &base.BaseResp{
			ErrCode: int64(common.Err_GateWayServerHandlerError),
			ErrMsg:  common.Err_GateWayServerHandlerError.String(),
			Extra:   nil,
		})
		return
	}

	c.JSON(http.StatusOK, realResp.Body)
}
