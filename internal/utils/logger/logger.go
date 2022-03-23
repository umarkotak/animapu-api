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
		RequestID    string `json:"request_id"`
		ErrorMessage string `json:"error_message"`
	}
)

var (
	GlobalLog []AnimapuLog
)

func Initialize() {
	GlobalLog = []AnimapuLog{}
}

func (ah *AnimapuHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	var reqID interface{}
	if ctx != nil {
		reqID = ctx.Value("request_id")
	}
	formattedLog := fmt.Sprintf("[%v][%v]: %v", time.Now().UTC(), reqID, entry.Message)
	GlobalLog = append(GlobalLog, AnimapuLog{RequestID: fmt.Sprintf("%v", reqID), ErrorMessage: formattedLog})
	return nil
}

func (ah *AnimapuHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
