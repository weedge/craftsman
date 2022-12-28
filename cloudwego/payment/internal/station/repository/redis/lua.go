package redis

var (
	assetStringChangeLua = `
if redis.call("exists", KEYS[1]) ~= 1 then
    return -1
end

if redis.call("type", KEYS[1]).ok == "string" then
    local assetStr = redis.call('get', KEYS[1]);
    local assetInfo = cjson.decode(assetStr);
    if assetInfo.assetCn == nil then
        assetInfo.assetCn = 0
    end
    local incr = 0
    if ARGV[1] ~= nil then
        incr = tonumber(ARGV[1]);
    end
    if assetInfo.assetCn + incr < 0 then
        return -2;
    end
    assetInfo.assetCn = assetInfo.assetCn + incr;
    assetStr = cjson.encode(assetInfo);
    if redis.call('set', KEYS[1], assetStr, 'ex', tonumber(ARGV[2]))
        and redis.call('set', KEYS[2], 1, 'ex', tonumber(ARGV[3])) then
        return 1;
    end
    return 0;
end

if redis.call("type", KEYS[1]).ok == "hash" then
    local assetCn = redis.call('hincrby', KEYS[1], 'assetCn', 0);
    local incr = 0
    if ARGV[1] ~= nil then
        incr = tonumber(ARGV[1]);
    end
    if assetCn + incr < 0 then
        return -2;
    end
    assetCn = assetCn + incr;
    if redis.call('hincrby', KEYS[1], 'assetCn', incr)
        and redis.call('expire', KEYS[1], tonumber(ARGV[2]))
        and redis.call('set', KEYS[2], 1, 'ex', tonumber(ARGV[3])) then
        return 1;
    end
    return 0;
end
`
)
