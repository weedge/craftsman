include "base.thrift"
namespace go station.payment

enum EventType{
    none,
    interactGift,
    orderApple,
    orderWX,
    orderAlipay,
    orderDouyin,
}

enum BizType{
    none,
    live,
    recharge
}

enum AssetType{
    none,
    goldCoin,
    diamond
}

struct BizAssetChangesReq{
    1: required list<BizEventAssetChange> bizAssetChanges,
}

struct BizEventAssetChange{
    1: required string eventId,
    2: required i64 opUserId,
    3: required EventType eventType,
    4: required i64 bizId,
    5: required BizType bizType,
    6: required string objId,
    7: required UserAssetChangeInfo opUserAssetChange,
    8: required UserAssetChangeInfo toUserAssetChange,
}

struct UserAssetChangeInfo{
    1: required i64 userId,
    2: required AssetType assetType,
    3: required i32 incr,
}

struct BizAssetChangesResp{
    1: required list<BizEventAssetChangerRes> bizAssetChangeResList,
    2: required i64 bizId,
    255: base.BaseResp baseResp,
}

struct BizEventAssetChangerRes{
    1: required string eventId,
    2: required bool changeRes,
    3: required string failMsg,
    4: required UserAsset opUserAsset,
}

struct UserAsset{
    1: required i64 userId,
    2: required AssetType assetType,
    3: required i32 assetCn,
}
