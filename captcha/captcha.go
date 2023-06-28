package captcha

import (
	"github.com/go-redis/redis"
	"github.com/mojocn/base64Captcha"
	"time"
)

func NewCaptcha(redis *redis.Client, prefix string, expiration time.Duration) *base64Captcha.Captcha {
	store := NewRedisStore(redis, prefix, expiration)
	source := "ABCDEFGHJKMNQRSTUVXYZ23456789"
	driver := base64Captcha.NewDriverString(100, 300, 0, base64Captcha.OptionShowHollowLine|base64Captcha.OptionShowSlimeLine, 4, source, nil, nil, nil)
	return base64Captcha.NewCaptcha(driver, store)
}
