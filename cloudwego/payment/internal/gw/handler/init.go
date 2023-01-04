package handler

import (
	"github.com/cloudwego/kitex/client/genericclient"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/injectors"
)

// if watch to reload diff config, use sync.RWMutex+Map
var gSvcMap = map[string]genericclient.Client{}
var gSvcOptsMap = map[string]*injectors.GenericEndpointsOpts{}

func InitSvcGenericClientMap(svcMap map[string]genericclient.Client, optsMap map[string]*injectors.GenericEndpointsOpts) {
	gSvcMap = svcMap
	gSvcOptsMap = optsMap
}
