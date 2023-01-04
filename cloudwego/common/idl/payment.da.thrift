include "base.thrift"
include "payment.base.thrift"
namespace go payment.da

struct GetAssetsReq{
    1: list<i64> userIds(api.query = 'userIds'),
}

struct GetAssetsResp{
    1: list<payment.base.UserAsset> userAssets,
    255: required base.BaseResp baseResp,
}

service PaymentService{
    GetAssetsResp GetAssets(1: GetAssetsReq req)(api.get = '/payment/da/v1/getassets', api.param = 'true')
}
