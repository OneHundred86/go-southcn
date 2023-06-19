package context

import (
	"github.com/OneHundred86/go-southcn/session"
	"github.com/gin-gonic/gin"
	"log"
)

func GetSession(c *gin.Context) session.Session {
	var sess session.Session
	if s, e := c.Get("session"); !e {
		log.Panic("session not started")
	} else {
		sess = s.(session.Session)
	}
	return sess
}
