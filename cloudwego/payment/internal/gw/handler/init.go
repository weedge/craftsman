package handler

import (
	"github.com/cloudwego/kitex/client/genericclient"
)

var gSvcMap = map[string]genericclient.Client{}

func InitSvcGenericClientMap(svcMap map[string]genericclient.Client) error {
	gSvcMap = svcMap
	return nil
}
