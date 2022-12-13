
# This repo's root import path (under $GOPATH)
KITEX_MODULE := github.com/weedge/craftsman/cloudwego/common

.PHONY: kitex
kitex:
	@rm -rf kitex_gen/
	@kitex -module $(KITEX_MODULE) ./idl/base.thrift
	@kitex -module $(KITEX_MODULE) ./idl/common.thrift
	@kitex -module $(KITEX_MODULE) ./idl/payment.thrift
