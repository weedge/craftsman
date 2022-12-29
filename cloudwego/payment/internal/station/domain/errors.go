package domain

import "errors"

var (
	ErrorMsgDecode = errors.New("event msg decode error")
	ErrorEventDone = errors.New("event had been done")

	ErrorInternalRedis = errors.New("internal redis error")

	ErrorNoEnoughAsset             = errors.New("asset not enough")
	ErrorSetAssetFromCBTimeOut     = errors.New("set asset from db exec timeout")
	ErrorSetAssetDoneOut           = errors.New("set asset from db exec done out")
	ErrorWatchAssetCasLoopMaxRetry = errors.New("watch user asset CAS loop reached maximum number of retries")

	ErrorSendAssetChangeEventMsg = errors.New("send user asset change event msg fail")
)
