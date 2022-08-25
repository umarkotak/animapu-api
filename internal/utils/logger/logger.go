package logger

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	AnimapuHook struct {
	}

	AnimapuLog struct {
		RequestID    string      `json:"request_id"`
		ErrorMessage string      `json:"error_message"`
		UtcTime      time.Time   `json:"utc_time"`
		Data         interface{} `json:"data"`
	}
)

var (
	GlobalLog []AnimapuLog
)

func Initialize() {
	GlobalLog = []AnimapuLog{}
}

func (ah *AnimapuHook) Fire(entry *logrus.Entry) error {
	if entry.Level != logrus.ErrorLevel {
		return nil
	}

	ctx := entry.Context
	var reqID interface{}
	if ctx != nil {
		reqID = ctx.Value("request_id")
	}

	if reqID == nil {
		reqID = ""
	}

	GlobalLog = append(GlobalLog, AnimapuLog{
		RequestID:    fmt.Sprintf("%v", reqID),
		UtcTime:      time.Now().UTC(),
		ErrorMessage: entry.Message,
		Data:         entry.Data,
	})

	return nil
}

func (ah *AnimapuHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
