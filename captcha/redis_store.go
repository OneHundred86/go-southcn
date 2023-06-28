package captcha

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/mojocn/base64Captcha"
	"time"
)

func NewRedisStore(redis *redis.Client, prefix string, expiration time.Duration) base64Captcha.Store {
	return &RedisStore{
		redis:      redis,
		prefix:     prefix,
		expiration: expiration,
	}
}

type RedisStore struct {
	base64Captcha.Store
	redis      *redis.Client
	prefix     string
	expiration time.Duration
}

func (s *RedisStore) genKey(id string) string {
	return fmt.Sprintf("%s@captcha@%s", s.prefix, id)
}

func (s *RedisStore) Set(id string, value string) error {
	key := s.genKey(id)
	cmd := s.redis.Set(key, value, s.expiration)
	return cmd.Err()
}

func (s *RedisStore) Get(id string, clear bool) string {
	key := s.genKey(id)
	cmd := s.redis.Get(key)

	if clear {
		s.redis.Del(key)
	}

	return cmd.Val()
}

func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	code := s.Get(id, clear)
	return code == answer
}
