package config

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gfile"
	"github.com/osgochina/dmicro/logger"
	"path/filepath"
	"strings"
)

type Config struct {
	configFile   string
	entries      *gmap.StrAnyMap
	ProgramGroup *ProcessGroup
}

func NewConfig(configFile string) *Config {
	return &Config{
		configFile: configFile,
		entries:    gmap.NewStrAnyMap(true),
	}
}

func (that *Config) Load() error {
	if !gfile.Exists(that.configFile) {
		return gerror.Newf("要载入的配置文件[%s]不存在", that.configFile)
	}
	that.ProgramGroup = NewProcessGroup()
	res, err := ParseIni(gfile.GetBytes(that.configFile))
	if err != nil {
		return err
	}
	if len(res) > 0 {
		that.parse(res)
	}

	return nil
}

func (that *Config) parse(res map[string]interface{}) {
	for key, val := range res {
		prefix := "program:"
		if strings.HasPrefix(key, prefix) {
			entrys := that.parseProgram(key[len(prefix):], val.(map[string]interface{}))
			if len(entrys) > 0 {
				for _, entry := range entrys {
					that.entries.Set(entry.Name, entry)
				}
			}
		}
	}
	return
}

//解析程序
func (that *Config) parseProgram(programName string, data map[string]interface{}) []*Entry {
	keyValues := gmap.NewStrAnyMap(true)
	keyValues.Sets(data)

	numProcess := keyValues.GetVar("numprocs").Int()
	if numProcess <= 0 {
		numProcess = 1
	}
	originalProcName := programName
	procName := keyValues.GetVar("process_name").String()
	if numProcess > 1 && strings.Index(procName, "%(process_num)") == -1 {
		logger.Errorf("no process_num[%d] in process name[%s]", numProcess, procName)
	}
	if len(procName) > 0 {
		originalProcName = procName
	}

	originalCmd := keyValues.GetVar("command").String()
	var entrys []*Entry
	for i := 1; i <= numProcess; i++ {
		keyValuesNum := keyValues.Clone()
		pne := NewProcessNameExpression(g.Map{
			"program_name": programName,
			"process_num":  fmt.Sprintf("%d", i),
			"group_name":   that.ProgramGroup.GetGroup(programName, programName),
			"here":         that.GetConfigFileDir()},
		)
		//把环境变量加入
		envValue := keyValues.GetVar("environment").String()
		if len(envValue) > 0 {
			for k, v := range *ParseEnv(envValue) {
				pne.Add(fmt.Sprintf("ENV_%s", k), v)
			}
		}
		cmd, err := pne.Eval(originalCmd)
		if err != nil {
			logger.Errorf("program [%s] get envs failed.err:%v", programName, err)
			continue
		}
		keyValuesNum.Set("command", cmd)

		procName, err := pne.Eval(originalProcName)
		if err != nil {
			logger.Errorf("program [%s] get envs failed，err:%v", programName, err)
			continue
		}
		keyValuesNum.Set("process_name", procName)
		keyValuesNum.Set("numprocs_start", fmt.Sprintf("%d", i-1))
		keyValuesNum.Set("process_num", fmt.Sprintf("%d", i))
		entry := NewEntry(procName, keyValuesNum, that.configFile, "program")
		group := that.ProgramGroup.GetGroup(programName, programName)
		entry.Group = group
		entrys = append(entrys, entry)
	}

	return entrys
}

func (that *Config) GetConfigFileDir() string {
	return filepath.Dir(that.configFile)
}

// GetEntries 获取栏目
func (that *Config) GetEntries(filterFunc func(entry *Entry) bool) []*Entry {
	result := make([]*Entry, 0)
	for _, val := range that.entries.Map() {
		entry := val.(*Entry)
		if filterFunc(entry) {
			result = append(result, entry)
		}
	}
	return result
}

// GetPrograms 获取可执行程序
func (that *Config) GetPrograms() []*Entry {
	programs := that.GetEntries(func(entry *Entry) bool {
		return entry.IsProgram()
	})

	return programs
}
