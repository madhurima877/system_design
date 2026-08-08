# kafkaWithGolang
main
 ↓
Processor.Start()
 ↓
worker pool
 ↓
jobs channel
 ↓
RateLimiter.Allow()
 ↓
processFn(msg)
 ↓
retry up to 3 times
 ↓
Store.Save()
 ↓
processed / failed counters
 ↓
Shutdown()
 ↓
Stats()