package repository

import (
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
)

type VisitorLog struct {
	mu           sync.Mutex
	QuickUserLog map[string][]string
}

var visitorLog = VisitorLog{}

func LogVisitor(c *gin.Context) {
	visitorLog.mu.Lock()
	defer visitorLog.mu.Unlock()

	visitorLog.QuickUserLog[c.Request.Header.Get("X-Visitor-Id")] = append(
		[]string{fmt.Sprintf("%v", c.Request.Header.Get("X-From-Path"))},
		visitorLog.QuickUserLog[c.Request.Header.Get("X-Visitor-Id")]...,
	)
	if len(visitorLog.QuickUserLog[c.Request.Header.Get("X-Visitor-Id")]) >= 100 {
		visitorLog.QuickUserLog[c.Request.Header.Get("X-Visitor-Id")] = visitorLog.QuickUserLog[c.Request.Header.Get("X-Visitor-Id")][:100]
	}
}

func GetLogVisitor() map[string][]string {
	return visitorLog.QuickUserLog
}
