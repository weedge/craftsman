namespace go base

struct BaseResp{
    1: required i64 errCode = 0,
    2: required string errMsg = "",
    3: optional map<string,string> extra,
}
