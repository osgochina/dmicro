package config

import (
	"bytes"
	"github.com/gogf/gf/container/gmap"
	"strings"
)

type ProcessGroup struct {
	processGroup *gmap.StrStrMap
}

func NewProcessGroup() *ProcessGroup {
	return &ProcessGroup{
		processGroup: gmap.NewStrStrMap(true),
	}
}

func (that *ProcessGroup) Clone() *ProcessGroup {
	newPg := NewProcessGroup()
	newPg.processGroup = that.processGroup.Clone()
	return newPg
}

func (that *ProcessGroup) Add(group string, procName string) {
	that.processGroup.Set(procName, group)
}

func (that *ProcessGroup) Remove(procName string) {
	that.processGroup.Remove(procName)
}

func (that *ProcessGroup) GetAllGroup() []string {

	groups := make(map[string]bool)
	for _, group := range that.processGroup.Map() {
		groups[group] = true
	}

	result := make([]string, 0)
	for group := range groups {
		result = append(result, group)
	}
	return result
}

func (that *ProcessGroup) GetAllProcess(group string) []string {
	result := make([]string, 0)
	for procName, groupName := range that.processGroup.Map() {
		if group == groupName {
			result = append(result, procName)
		}
	}
	return result
}

func (that *ProcessGroup) InGroup(procName string, group string) bool {
	groupName, ok := that.processGroup.Search(procName)
	if ok && group == groupName {
		return true
	}
	return false
}

func (that *ProcessGroup) GetGroup(procName string, defGroup string) string {
	group, ok := that.processGroup.Search(procName)

	if ok {
		return group
	}
	that.processGroup.Set(procName, defGroup)
	return defGroup
}

func (that *ProcessGroup) ForEachProcess(procFunc func(group string, procName string)) {
	that.processGroup.Iterator(func(procName string, groupName string) bool {
		procFunc(groupName, procName)
		return true
	})
}

func (that *ProcessGroup) String() string {
	buf := bytes.NewBuffer(make([]byte, 0))
	for _, group := range that.processGroup.Map() {
		buf.WriteString(group)
		buf.WriteString(":")
		buf.WriteString(strings.Join(that.GetAllProcess(group), ","))
		buf.WriteString(";")
	}
	return buf.String()
}
