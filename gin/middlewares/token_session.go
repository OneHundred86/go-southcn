package middlewares

import (
	"github.com/OneHundred86/go-southcn/cache"
	"github.com/OneHundred86/go-southcn/encrypter"
	"github.com/OneHundred86/go-southcn/session"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
	"time"
)

type TokenSession struct {
	Encrypter encrypter.Encrypter
	Cache     cache.Cache
	Header    string
	KeyPrefix string
	TTL       time.Duration
}

func (s *TokenSession) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		cacheSession := session.NewCacheSession(s.Cache, s.KeyPrefix, s.TTL)

		tokenCipher := c.GetHeader(s.Header)
		if tokenCipher != "" {
			token, err := s.Encrypter.Decrypt(tokenCipher)
			if err == nil && token != "" {
				cacheSession.SetSessionID(token)
				cacheSession.Load()
			} else {
				token = s.newToken()
				cacheSession.SetSessionID(token)
			}
		} else {
			token := s.newToken()
			cacheSession.SetSessionID(token)
		}

		// 放入gin.Context
		c.Set("session", cacheSession)

		c.Next()

		cacheSession.Save()
	}
}

func (s *TokenSession) newToken() string {
	token := uuid.NewString()
	return strings.ReplaceAll(token, "-", "")
}
