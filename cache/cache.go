package cache

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

type Cache interface {
	Get(key string, value interface{}) bool
	Set(key string, value interface{}, expiration time.Duration) bool
	Forever(key string, value interface{}) bool
	Delete(key string) bool
	Pull(key string, value interface{}) bool
}

type RedisCache struct {
	Cache

	redis     *redis.Client
	encoder   Encoder
	keyPrefix string
}

func NewRedisCache(client *redis.Client, encoder Encoder, keyPrefix string) Cache {
	return &RedisCache{
		redis:     client,
		encoder:   encoder,
		keyPrefix: keyPrefix,
	}
}

func makeCacheKey(prefix, key string) string {
	return fmt.Sprintf("%s%s", prefix, key)
}

func (cache *RedisCache) Get(key string, value interface{}) bool {
	key = makeCacheKey(cache.keyPrefix, key)
	result, err := cache.redis.Get(key).Result()
	if err != nil {
		// log todo
		return false
	}

	err = cache.encoder.Decode(result, value)
	if err != nil {
		// log todo
		return false
	}

	return true
}

func (cache *RedisCache) Set(key string, value interface{}, expiration time.Duration) bool {
	v, err := cache.encoder.Encode(value)
	if err != nil {
		log.Panic(err)
		return false
	}

	key = makeCacheKey(cache.keyPrefix, key)
	result, err := cache.redis.Set(key, v, expiration).Result()
	if err != nil {
		log.Panic(err)
		return false
	}

	return result == "OK"
}

func (cache *RedisCache) Forever(key string, value interface{}) bool {
	key = makeCacheKey(cache.keyPrefix, key)
	return cache.Set(key, value, 0)
}

func (cache *RedisCache) Delete(key string) bool {
	key = makeCacheKey(cache.keyPrefix, key)
	result, err := cache.redis.Del(key).Result()
	if err != nil {
		// log todo
		return false
	}

	return result == 1
}

func (cache *RedisCache) Pull(key string, value interface{}) bool {
	key = makeCacheKey(cache.keyPrefix, key)
	result, err := cache.redis.Do("GETDEL", key).Result()
	if err != nil {
		// log todo
		return false
	}

	str, ok := result.(string)
	if !ok {
		panic(fmt.Sprintf("key[%s] not string", key))
	}

	err = cache.encoder.Decode(str, value)
	if err != nil {
		// log todo
		return false
	}
	return true
}
