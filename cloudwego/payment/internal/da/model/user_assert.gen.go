// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package model

import (
	"time"
)

const TableNameUserAssert = "user_assert"

// UserAssert mapped from table <user_assert>
type UserAssert struct {
	ID        int64     `gorm:"column:id;type:bigint unsigned;primaryKey;autoIncrement:true" json:"id"`
	UserID    int64     `gorm:"column:userId;type:bigint unsigned;not null" json:"userId"`
	AssetCn   int64     `gorm:"column:assetCn;type:bigint unsigned;not null" json:"assetCn"`
	AssetType int32     `gorm:"column:assetType;type:tinyint unsigned;not null" json:"assetType"`
	Version   int64     `gorm:"column:version;type:bigint unsigned;not null" json:"version"`
	CreatedAt time.Time `gorm:"column:createdAt;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updatedAt;type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updatedAt"`
}

// TableName UserAssert's table name
func (*UserAssert) TableName() string {
	return TableNameUserAssert
}
