package mysql

import "github.com/google/wire"

var RepositorySet = wire.NewSet(NewUserAssetEventRepository, NewUserAssetRecordRepository, NewUserAssetRepository)
