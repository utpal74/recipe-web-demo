
local key = KEYS[1]

local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local request_id = ARGV[4]

local window_start = now - window

redis.call("ZREMRANGEBYSCORE", key, "-inf", window_start)

local count = redis.call("ZCARD", key)

if count > limit then
   return {0, count}
end

redis.call("ZADD", key, now, request_id)
redis.call("EXPIRE", key, window)

return {1, count+1}
