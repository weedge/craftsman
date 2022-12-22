namespace go payment.base

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

struct UserAsset{
    1: required i64 userId,
    2: required AssetType assetType,
    3: required i64 assetCn,
}