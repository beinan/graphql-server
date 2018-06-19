package downstreams

import "github.com/go-redis/redis"

var RedisClient = redis.NewClient(&redis.Options{
	Addr:     "redis:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})
