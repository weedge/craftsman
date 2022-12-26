package constants

import (
	"fmt"
	"time"
)

const (
	RedisClientTypeReplica = 1
	RedisClientTypeCluster = 2
)
const (
	RedisKeyPrefixUserAssetInfo     = "I.asset."
	RedisKeyPrefixGiftInfo          = "I.gift."
	RedisKeyPrefixUserInfo          = "I.user."
	RedisKeyPrefixRoomInfo          = "I.room."
	RedisKeyPrefixUserAssetInfoLock = "L.asset."
	RedisKeyPrefixAssetEventMsg     = "M.asset."
)
const (
	RedisKeyExpireUserAssetInfo     = 86400 * time.Second
	RedisKeyExpireGiftInfo          = 7 * 86400 * time.Second
	RedisKeyExpireUserInfo          = 86400 * time.Second
	RedisKeyExpireRoomInfo          = 86400 * time.Second
	RedisKeyExpireUserAssetInfoLock = 60 * time.Second
	RedisKeyExpireAssetEventMsg     = 86400 * time.Second
)
const (
	DisLockerBlockRetryCn = 100
	// MaxRetries Redis transactions use optimistic locking.
	MaxRetries = 1000
)
const (
	TimeOutSetAssetFromCB = 3 * time.Second
)

func GetUserAssetInfoKey(userId int64) string {
	return fmt.Sprintf("%s{%d}", RedisKeyPrefixUserAssetInfo, userId)
}
func GetGiftInfoKey(giftId int64) string {
	return fmt.Sprintf("%s{%d}", RedisKeyPrefixGiftInfo, giftId)
}
func GetUserInfoKey(userId int64) string {
	return fmt.Sprintf("%s{%d}", RedisKeyPrefixUserInfo, userId)
}
func GetRoomInfoKey(roomId int64) string {
	return fmt.Sprintf("%s{%d}", RedisKeyPrefixRoomInfo, roomId)
}
func GetUserAssetInfoLockKey(userId int64, tag string) string {
	return fmt.Sprintf("%s{%d}.%s", RedisKeyPrefixUserAssetInfoLock, userId, tag)
}
func GetUserAssetEventMsgKey(opUserId int64, eventId string) string {
	return fmt.Sprintf("%s{%d}.%s", RedisKeyPrefixAssetEventMsg, opUserId, eventId)
}
