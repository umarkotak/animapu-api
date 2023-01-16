package repository

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	quickUserLog = map[string][]string{}
)

func LogVisitor(c *gin.Context) {
	quickUserLog[c.Request.Header.Get("X-Visitor-Id")] = append(
		[]string{fmt.Sprintf("%v", c.Request.Header.Get("X-From-Path"))},
		quickUserLog[c.Request.Header.Get("X-Visitor-Id")]...,
	)
	if len(quickUserLog[c.Request.Header.Get("X-Visitor-Id")]) >= 100 {
		quickUserLog[c.Request.Header.Get("X-Visitor-Id")] = quickUserLog[c.Request.Header.Get("X-Visitor-Id")][:100]
	}
}

func GetLogVisitor() map[string][]string {
	return quickUserLog
}
