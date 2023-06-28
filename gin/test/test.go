package main

import (
	"github.com/OneHundred86/go-southcn/cache"
	myCaptcha "github.com/OneHundred86/go-southcn/captcha"
	"github.com/OneHundred86/go-southcn/encrypter"
	"github.com/OneHundred86/go-southcn/gin/context"
	"github.com/OneHundred86/go-southcn/gin/middlewares"
	"github.com/OneHundred86/go-southcn/response"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	Sm4Key = "4f17d993e1c602bc7cfa92377e223e6b"
)

var (
	Redis *redis.Client
)

func init() {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "root123",
		DB:       1,
	})
}

func main() {
	r := gin.Default()

	redisCache := cache.NewRedisCache(Redis, &cache.JsonEncoder{}, "test-")
	tokenSession := middlewares.TokenSession{
		Encrypter: encrypter.NewSm4EcbEncrypter(Sm4Key),
		Cache:     redisCache,
		Header:    "Auth",
		KeyPrefix: "session",
		TTL:       time.Hour * 2,
	}
	captcha := myCaptcha.NewCaptcha(Redis, "test", time.Minute*10)

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

	captGroup := r.Group("captcha")
	{
		captGroup.GET("code", func(c *gin.Context) {
			id, s, err := captcha.Generate()
			if err != nil {
				c.JSON(http.StatusInternalServerError, response.ErrMsg(err.Error()))
				return
			}

			c.JSON(http.StatusOK, response.Ok(gin.H{
				"id":   id,
				"code": s,
			}))
		})

		captGroup.POST("verify", func(c *gin.Context) {
			var req struct {
				ID   string `form:"id" binding:"required"`
				Code string `form:"code" binding:"required"`
			}

			if err := c.ShouldBind(&req); err != nil {
				c.JSON(http.StatusBadRequest, response.ErrMsg(err.Error()))
				return
			}

			answer := strings.ToUpper(req.Code)

			if captcha.Verify(req.ID, answer, true) {
				c.JSON(http.StatusOK, response.Ok(nil))
			} else {
				c.JSON(http.StatusOK, response.ErrMsg("验证码错误"))
			}
		})
	}

	r.Run(":8000")
}

/**
curl -XPOST localhost:8000/test/session/set -d"id=123"
curl -H"Auth:xxx" localhost:8000/test/session/all

curl localhost:8000/captcha/code
curl -XPOST localhost:8000/captcha/verify -d"id=xoqyjMNEYa0v6k9daQB9&code=ATGG"
*/
