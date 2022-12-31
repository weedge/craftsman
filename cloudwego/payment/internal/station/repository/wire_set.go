package repository

import (
	"github.com/google/wire"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/repository/redis"
	"github.com/weedge/craftsman/cloudwego/payment/internal/station/repository/rmq"
	"golang.org/x/sync/singleflight"
)

var ProviderSet = wire.NewSet(
	wire.Value(&singleflight.Group{}),
	redis.NewUserAssetCallBack,

	redis.NewUserAssetEventRepository,
	redis.NewUserAssetRepository,

	rmq.NewUserAssetEventMsg,
	rmq.NewUserAssetEventMsgLister,
)
