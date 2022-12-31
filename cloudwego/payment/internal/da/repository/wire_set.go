package repository

import (
	"github.com/google/wire"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da/repository/mysql"
)

var ProviderSet = wire.NewSet(
	mysql.NewUserAssetEventRepository,
	mysql.NewUserAssetRecordRepository,
	mysql.NewUserAssetRepository,
)
