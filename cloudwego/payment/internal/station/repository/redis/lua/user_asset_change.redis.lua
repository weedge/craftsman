-- author: weedge
-- params: KEYS[1] user asset key
-- params: KEYS[2] event msg key
-- params: ARGV[1] incr asset num eg:1,-1
-- params: ARGV[2] user asset key ttl
-- params: ARGV[3] event msg key ttl
-- return >=0:操作成功,资产数, -3:无操作，-1:缓存资产不存在，-2:资产不足，
-- debug:
--  redis-cli --ldb --eval user_asset_change.redis.lua I.asset.{100} M.asset.{100}.`ksuid` , 100 86400 86400
--  redis-cli -c -p 26383 --ldb --eval user_asset_change.redis.lua I.asset.{100} M.asset.{100}.`ksuid` , 100 86400 86400
--  restart to slot debug

if redis.call("exists", KEYS[1]) ~= 1 then
    return -1
end

if redis.call("type", KEYS[1]).ok == "string" then
    local assetStr = redis.call('get', KEYS[1])
    local assetInfo = cjson.decode(assetStr)
    if assetInfo.assetCn == nil then
        assetInfo.assetCn = 0
    end
    local incr = 0
    if ARGV[1] ~= nil then
        incr = tonumber(ARGV[1])
    end
    if assetInfo.assetCn + incr < 0 then
        return -2
    end
    assetInfo.assetCn = assetInfo.assetCn + incr
    assetStr = cjson.encode(assetInfo)
    if redis.call('set', KEYS[1], assetStr, 'ex', tonumber(ARGV[2]))
        and redis.call('set', KEYS[2], 1, 'ex', tonumber(ARGV[3])) then
        return assetInfo.assetCn
    end
    return -3
end

if redis.call("type", KEYS[1]).ok == "hash" then
    local assetCn = redis.call('hincrby', KEYS[1], 'assetCn', 0)
    local incr = 0
    if ARGV[1] ~= nil then
        incr = tonumber(ARGV[1])
    end
    if assetCn + incr < 0 then
        return -2
    end
    assetCn = assetCn + incr
    if redis.call('hincrby', KEYS[1], 'assetCn', incr)
        and redis.call('expire', KEYS[1], tonumber(ARGV[2]))
        and redis.call('set', KEYS[2], 1, 'ex', tonumber(ARGV[3])) then
        return assetCn
    end
    return -3
end
