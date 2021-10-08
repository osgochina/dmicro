package config

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gvar"
)

type Entry struct {
	Name       string
	Group      string
	entryType  string
	configFile string
	keyValues  *gmap.StrAnyMap
}

// NewEntry 创建配置文件条目
func NewEntry(name string, data *gmap.StrAnyMap, configFile string, entryType string) *Entry {
	return &Entry{
		configFile: configFile,
		Name:       name,
		keyValues:  data,
		entryType:  entryType,
	}
}

func (that *Entry) IsProgram() bool {
	return that.entryType == "program"
}

func (that *Entry) String() string {
	return that.keyValues.String()
}

func (that *Entry) Get(key string) *gvar.Var {
	return that.keyValues.GetVar(key)
}

func (that *Entry) Set(key string, val interface{}) {
	that.keyValues.Set(key, val)
}

func (that *Entry) Map() *gmap.StrAnyMap {
	return that.keyValues
}
