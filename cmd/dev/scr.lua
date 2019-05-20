local res, j, queue
res = redis.call('zrangebyscore', KEYS[1], '-inf', ARGV[2], 'LIMIT', 0, 1)
if #res > 0 then
  j = cjson.decode(res[1])
  redis.call('zrem', KEYS[1], res[1])
  queue = ARGV[1] .. j['name']
  for _,v in pairs(KEYS) do
    if v == queue then
      j['t'] = tonumber(ARGV[2])
      redis.call('lpush', queue, cjson.encode(j))
      return 'ok'
    end
  end
  j['err'] = 'unknown job when requeueing'
  j['failed_at'] = tonumber(ARGV[2])
  redis.call('zadd', KEYS[2], ARGV[2], cjson.encode(j))
  return 'dead' -- put on dead queue
end
return nil