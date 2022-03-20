package logger

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

type (
	AnimapuHook struct {
	}
)

var (
	GlobalLog []string
)

func Initialize() {
	GlobalLog = []string{}
}

func (ah *AnimapuHook) Fire(entry *logrus.Entry) error {
	ctx := entry.Context
	var reqID interface{}
	if ctx != nil {
		reqID = ctx.Value("request_id")
	}
	formattedLog := fmt.Sprintf("[%v][%v]: %v", time.Now().UTC(), reqID, entry.Message)
	GlobalLog = append(GlobalLog, formattedLog)
	return nil
}

func (ah *AnimapuHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
