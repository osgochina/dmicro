package drpc

import (
	"errors"
	"github.com/gogf/gf/os/glog"
)

// PluginContainer 插件容器
type PluginContainer struct {
	*pluginSingleContainer
	left        *pluginSingleContainer
	middle      *pluginSingleContainer
	right       *pluginSingleContainer
	refreshTree func()
}

//创建一个插件容器
func newPluginContainer() *PluginContainer {
	p := &PluginContainer{
		pluginSingleContainer: newPluginSingleContainer(),
		left:                  newPluginSingleContainer(),
		middle:                newPluginSingleContainer(),
		right:                 newPluginSingleContainer(),
	}
	p.refreshTree = func() { p.refresh() }
	return p
}

// 克隆新的服务
func (that *PluginContainer) cloneAndAppendMiddle(plugins ...Plugin) *PluginContainer {
	middle := newPluginSingleContainer()
	middle.plugins = append(that.middle.GetAll(), plugins...)

	//clone 新的插件对象，并刷新它
	newPluginC := newPluginContainer()
	newPluginC.middle = middle
	newPluginC.left = that.left
	newPluginC.right = that.right
	newPluginC.refresh()

	//老的插件容器也要保存，并且把每次clone的新的插件容器对象的刷新方法加入
	//因为老的容器有可能会删除，添加，这样刷新的时候，可以自动更新到它的克隆体中
	oldRefreshTree := that.refreshTree
	that.refreshTree = func() {
		oldRefreshTree()
		newPluginC.refresh()
	}
	return newPluginC
}

// AppendLeft 追加插件到左边
func (that *PluginContainer) AppendLeft(plugins ...Plugin) {
	that.left.appendLeft(plugins...)
	that.refreshTree()
}

// AppendRight 追加插件到右边
func (that *PluginContainer) AppendRight(plugins ...Plugin) {
	that.right.appendRight(plugins...)
	that.refreshTree()
}

// Remove 根据插件名移除插件
func (that *PluginContainer) Remove(pluginName string) error {
	err := that.pluginSingleContainer.remove(pluginName)
	if err != nil {
		return err
	}
	_ = that.left.remove(pluginName)
	_ = that.middle.remove(pluginName)
	_ = that.right.remove(pluginName)
	that.refreshTree()
	return nil
}

//刷新
func (that *PluginContainer) refresh() {
	count := len(that.left.plugins) + len(that.middle.plugins) + len(that.right.plugins)
	allPlugins := make([]Plugin, count)
	copy(allPlugins[0:], that.left.plugins)
	copy(allPlugins[0+len(that.left.plugins):], that.middle.plugins)
	copy(allPlugins[0+len(that.left.plugins)+len(that.middle.plugins):], that.right.plugins)
	m := make(map[string]bool, count)
	for _, plugin := range allPlugins {
		if plugin == nil {
			glog.Fatalf("plugin cannot be nil!")
			return
		}
		if m[plugin.Name()] {
			glog.Fatalf("repeat add plugin: %s", plugin.Name())
			return
		}
		m[plugin.Name()] = true
	}
	that.pluginSingleContainer.plugins = allPlugins
}

// 插件的单一容器
type pluginSingleContainer struct {
	plugins []Plugin
}

func newPluginSingleContainer() *pluginSingleContainer {
	return &pluginSingleContainer{
		plugins: make([]Plugin, 0),
	}
}

// 把新的插件追加到左边
func (that *pluginSingleContainer) appendLeft(plugins ...Plugin) {
	if len(plugins) == 0 {
		return
	}
	that.plugins = append(plugins, that.plugins...)
}

//把新的插件追加到右边
func (that *pluginSingleContainer) appendRight(plugins ...Plugin) {
	if len(plugins) == 0 {
		return
	}
	that.plugins = append(that.plugins, plugins...)
}

// GetByName 通过插件名字获取插件
func (that *pluginSingleContainer) GetByName(pluginName string) Plugin {
	if that.plugins == nil {
		return nil
	}
	for _, plugin := range that.plugins {
		if plugin.Name() == pluginName {
			return plugin
		}
	}
	return nil
}

// GetAll 获取所有插件列表
func (that *pluginSingleContainer) GetAll() []Plugin {
	return that.plugins
}

// 通过插件名字在插件容器中移除插件
func (that *pluginSingleContainer) remove(pluginName string) error {
	if that.plugins == nil {
		return errors.New("no plugins are registered yet")
	}
	if len(pluginName) == 0 {
		//return error: cannot delete an unamed plugin
		return errors.New("plugin with an empty name cannot be removed")
	}
	indexToRemove := -1
	for i, plugin := range that.plugins {
		if plugin.Name() == pluginName {
			indexToRemove = i
			break
		}
	}
	if indexToRemove == -1 {
		return errors.New("cannot remove a plugin which isn't exists")
	}
	that.plugins = append(that.plugins[:indexToRemove], that.plugins[indexToRemove+1:]...)
	return nil
}

func warnInvalidHandlerHooks(plugin []Plugin) {
	for _, p := range plugin {
		switch p.(type) {
		case BeforeNewEndpointPlugin:
			glog.Debugf("invalid BeforeNewEndpointPlugin in router: %s", p.Name())
		case AfterNewEndpointPlugin:
			glog.Debugf("invalid AfterNewEndpointPlugin in router: %s", p.Name())
		case AfterDialPlugin:
			glog.Debugf("invalid AfterDialPlugin in router: %s", p.Name())
		case AfterAcceptPlugin:
			glog.Debugf("invalid AfterAcceptPlugin in router: %s", p.Name())
		case BeforeWriteCallPlugin:
			glog.Debugf("invalid BeforeWriteCallPlugin in router: %s", p.Name())
		case AfterWriteCallPlugin:
			glog.Debugf("invalid AfterWriteCallPlugin in router: %s", p.Name())
		case BeforeWritePushPlugin:
			glog.Debugf("invalid BeforeWritePushPlugin in router: %s", p.Name())
		case AfterWritePushPlugin:
			glog.Debugf("invalid AfterWritePushPlugin in router: %s", p.Name())
		case BeforeReadHeaderPlugin:
			glog.Debugf("invalid BeforeReadHeaderPlugin in router: %s", p.Name())
		case AfterReadCallHeaderPlugin:
			glog.Debugf("invalid AfterReadCallHeaderPlugin in router: %s", p.Name())
		case AfterReadPushHeaderPlugin:
			glog.Debugf("invalid AfterReadPushHeaderPlugin in router: %s", p.Name())
		}
	}
}
