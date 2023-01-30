package utils

import (
	"bridge-allowance/config"
	"go.uber.org/zap"
	// "reflect"
)

type UtilConf struct {
	log  *zap.SugaredLogger
	conf *config.Config
}

func NewUtils(log *zap.SugaredLogger, conf *config.Config) *UtilConf {
	return &UtilConf{
		log,
		conf,
	}
}
