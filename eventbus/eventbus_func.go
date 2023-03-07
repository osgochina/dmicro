package eventbus

import (
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"strings"
)

// 检查事件名是否符合规范
func checkName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", gerror.New("event: 事件名不能为空")
	}

	if !gregex.IsMatchString(`^[a-zA-Z][\w-.*]*$`, name) {
		return "", gerror.New(`event: 事件名格式不正确,请匹配'^[a-zA-Z][\w-.]*$'`)
	}

	return name, nil
}
