package repository

import (
	"fmt"
	"sync"
)

type VisitorLog struct {
	mu           sync.Mutex
	QuickUserLog map[string][]string
}

var visitorLog = VisitorLog{
	QuickUserLog: map[string][]string{},
}

func LogVisitor(visitorID, fromPath string) {
	visitorLog.mu.Lock()
	defer visitorLog.mu.Unlock()

	visitorLog.QuickUserLog[visitorID] = append(
		[]string{fmt.Sprintf("%v", fromPath)},
		visitorLog.QuickUserLog[visitorID]...,
	)
	if len(visitorLog.QuickUserLog[visitorID]) >= 100 {
		visitorLog.QuickUserLog[visitorID] = visitorLog.QuickUserLog[visitorID][:100]
	}
}

func GetLogVisitor() map[string][]string {
	return visitorLog.QuickUserLog
}
