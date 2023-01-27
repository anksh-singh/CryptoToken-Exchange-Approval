package utils

import (
	"go.uber.org/zap"
	"runtime/debug"
)

func (u *UtilConf) CleanUp(log *zap.SugaredLogger) {
	if r := recover(); r != nil {
		u.log.Errorf("Recovered from a panic %v  \n  stack trace from panic \n %v  ", r, string(debug.Stack()))
	}
}
