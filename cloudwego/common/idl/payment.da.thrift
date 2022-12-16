include "base.thrift"
include "payment.base.thrift"
namespace go payment.da

struct GetAssetReq{
    1: list<i64> userIds(api.query = 'userIds'),
}

struct GetAssetResp{
    1: list<payment.base.UserAsset> userAsset,
    255: required base.BaseResp baseResp,
}

service PaymentService{
    GetAssetResp GetAsset(1: GetAssetReq req)(api.get = '/payment/getassets', api.param = 'true') 
}
