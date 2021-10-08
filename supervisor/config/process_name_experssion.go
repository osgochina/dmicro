package config

import (
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/os/genv"
	"github.com/gogf/gf/util/gconv"
	"os"
	"strings"
)

type ProcessNameExpression struct {
	data *gmap.StrAnyMap
}

// NewProcessNameExpression 创建 变量池子
func NewProcessNameExpression(envs ...map[string]interface{}) *ProcessNameExpression {
	pne := &ProcessNameExpression{
		data: gmap.NewStrAnyMap(true),
	}

	for k, v := range genv.Map() {
		pne.data.Set("ENV_"+k, v)
	}
	if len(envs) > 0 {
		pne.data.Sets(envs[0])
	}
	hostname, err := os.Hostname()
	if err == nil {
		pne.data.Set("host_node_name", hostname)
	}

	return pne
}

// Add 添加变量
func (that *ProcessNameExpression) Add(key string, value string) *ProcessNameExpression {
	that.data.Set(key, value)
	return that
}

// Eval 执行变量替换
func (that *ProcessNameExpression) Eval(s string) (string, error) {
	for {
		//判断格式是否正确，必须是 "%(" 开头
		start := strings.Index(s, "%(")
		if start == -1 {
			return s, nil
		}
		end := start + 1
		n := len(s)

		// 查找变量的结束位置
		for end < n && s[end] != ')' {
			end++
		}

		// 查找变量类型
		typ := end + 1
		for typ < n && !((s[typ] >= 'a' && s[typ] <= 'z') || (s[typ] >= 'A' && s[typ] <= 'Z')) {
			typ++
		}

		if typ < n {
			varName := s[start+2 : end]

			varValue, ok := that.data.Search(varName)

			if !ok {
				return "", fmt.Errorf("fail to find the environment variable %s", varName)
			}
			if s[typ] == 'd' {
				i := gconv.Int(varValue)
				if i == 0 {
					return "", fmt.Errorf("can't convert %s to integer", varValue)
				}
				s = s[0:start] + fmt.Sprintf("%"+s[end+1:typ+1], i) + s[typ+1:]
			} else if s[typ] == 's' {
				s = s[0:start] + gconv.String(varValue) + s[typ+1:]
			} else {
				return "", fmt.Errorf("not implement type:%v", s[typ])
			}
		} else {
			return "", fmt.Errorf("invalid string expression format")
		}
	}
}
