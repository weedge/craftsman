package domain

import "errors"

var (
	ErrorNoEnoughAsset             = errors.New("asset not enough")
	ErrorSetAssetFromCBTimeOut     = errors.New("set asset from db exec timeout")
	ErrorSetAssetDoneOut           = errors.New("set asset from db exec done out")
	ErrorWatchAssetCasLoopMaxRetry = errors.New("watch user asset CAS loop reached maximum number of retries")
)
