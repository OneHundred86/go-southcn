package main

import (
	"github.com/OneHundred86/go-southcn/cache"
	"github.com/OneHundred86/go-southcn/encrypter"
	"github.com/OneHundred86/go-southcn/gin/context"
	"github.com/OneHundred86/go-southcn/gin/middlewares"
	"github.com/OneHundred86/go-southcn/response"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"time"
)

const (
	Sm4Key = "4f17d993e1c602bc7cfa92377e223e6b"
)

func main() {
	r := gin.Default()

	redisCache := cache.NewRedisCache(redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root123",
		DB:       1,
	}), &cache.JsonEncoder{}, "test-")
	tokenSession := middlewares.TokenSession{
		Encrypter: encrypter.NewSm4EcbEncrypter(Sm4Key),
		Cache:     redisCache,
		Header:    "Auth",
		KeyPrefix: "session-",
		TTL:       time.Hour * 2,
	}
	group := r.Group("/test/session", tokenSession.Middleware())
	{
		group.POST("/set", func(c *gin.Context) {
			session := context.GetSession(c)

			if id, ok := c.GetPostForm("id"); ok {
				session.Set("id", id)

				encrypt := encrypter.NewSm4EcbEncrypter(Sm4Key)
				token, err := encrypt.Encrypt(session.GetSessionID())
				if err != nil {
					log.Panic(err)
				}

				c.JSON(http.StatusOK, response.Ok(map[string]any{
					"Auth": token,
				}))
			} else {
				c.JSON(http.StatusOK, response.ErrMsg("id undefined"))
			}
		})

		group.GET("/all", func(c *gin.Context) {
			session := context.GetSession(c)

			c.JSON(http.StatusOK, response.Ok(session.All()))
		})
	}

	r.Run(":8000")
}

/**
curl -XPOST localhost:8000/test/session/set -d"id=123"
curl -H"Auth:xxx" localhost:8000/test/session/all
*/
