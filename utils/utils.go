package utils

import (
	"bridge-allowance/config"
	"go.uber.org/zap"
	"reflect"
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

// FetchAttributes Common util to fetch interface into structured attributes based on reflect types
func (u *UtilConf) FetchAttributes(m interface{}) map[string]reflect.Type {
	typ := reflect.TypeOf(m)
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	attrs := make(map[string]reflect.Type)
	if typ.Kind() != reflect.Struct {
		u.log.Infof("%v type can't have attributes inspected\n", typ.Kind())
		return attrs
	}
	for i := 0; i < typ.NumField(); i++ {
		p := typ.Field(i)
		if !p.Anonymous {
			attrs[p.Name] = p.Type
		}
	}
	return attrs
}
