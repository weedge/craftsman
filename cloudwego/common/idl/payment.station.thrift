include "base.thrift"
include "payment.base.thrift"
namespace go payment.station

struct BizAssetChangesReq{
    1: required list<BizEventAssetChange> bizAssetChanges,
}

struct BizEventAssetChange{
    1: required string eventId,
    2: required i64 opUserId,
    3: required payment.base.EventType eventType,
    4: required i64 bizId,
    5: required payment.base.BizType bizType,
    6: required string objId,
    7: required UserAssetChangeInfo opUserAssetChange,
    8: required UserAssetChangeInfo toUserAssetChange,
}

struct UserAssetChangeInfo{
    1: required i64 userId,
    2: required payment.base.AssetType assetType,
    3: required i32 incr,
}

struct BizAssetChangesResp{
    1: required list<BizEventAssetChangerRes> bizAssetChangeResList,
    255: required base.BaseResp baseResp,
}

struct BizEventAssetChangerRes{
    1: required string eventId,
    2: required bool changeRes,
    3: required string failMsg,
    4: required payment.base.UserAsset opUserAsset,
}

service PaymentService{
    BizAssetChangesResp ChangeAsset(1: BizAssetChangesReq req)(api.post = '/payment/station/v1/changeassets', api.param = 'true', api.serializer = 'json')
}