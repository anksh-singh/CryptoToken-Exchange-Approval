package utils

import (
	conf "bridge-allowance/config"
	"gopkg.in/yaml.v2"
	"reflect"
	"testing"
)

func TestNewUtils(t *testing.T) {

	configobj := conf.LoadConfig("config", "../config/test")

	u := NewUtils(testlogger, configobj)
	if !reflect.DeepEqual(u.conf, configobj) {
		t.Errorf("NewUtils() got = %v, want %v", u.conf, configobj)
	}

	for _, arg := range expectargs {
		t.Run(arg.name, func(t *testing.T) {

			configobj := conf.LoadConfig(arg.name, arg.path)
			u := NewUtils(testlogger, configobj)
			got := u.GetEVMChains()
			var configuration conf.Config
			data := getConfig(arg.name)
			err := yaml.Unmarshal([]byte(data), &configuration)
			if err != nil {
				t.Errorf("TestNewUtils() got = %v", err)
			}

			if !reflect.DeepEqual(got, arg.want) && !(len(got) == 0 && len(arg.want) == 0) {
				t.Errorf("GetEVMChains() got = %v, want %v", got, arg.want)
			}
		})
	}

}
