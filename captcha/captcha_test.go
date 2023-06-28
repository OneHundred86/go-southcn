package captcha

import (
	"github.com/go-redis/redis"
	"log"
	"testing"
	"time"
)

func TestNewCaptcha(t *testing.T) {
	redis := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root123",
		DB:       0,
	})
	captcha := NewCaptcha(redis, "test", 5*time.Second)
	id, b64s, err := captcha.Generate()
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("base64 captcha: %s", b64s)

	answer := captcha.Store.Get(id, false)
	// time.Sleep(5 * time.Second)

	match := captcha.Verify(id, answer, true)

	log.Printf("answer: %s, match: %v \n", answer, match)
}
