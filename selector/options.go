package selector

import (
	"context"
	"github.com/osgochina/dmicro/registry"
)

// Options selector的配置参数
type Options struct {
	// 服务注册表
	Registry registry.Registry
	// 节点选择策略引擎
	Strategy Strategy
	// 扩展配置，可以添加自定义选项
	Context context.Context
}

// SelectOptions 节点选择器参数
type SelectOptions struct {
	// 节点过滤器列表
	Filters []Filter
	// 节点选择策略引擎
	Strategy Strategy
	// 扩展配置，可以添加自定义选项
	Context context.Context
}

// Option 根据配置选项初始化 selector
type Option func(*Options)

// SelectOption 调用select 方法的时候传入的配置
type SelectOption func(*SelectOptions)

// Registry 设置selector的注册表对象
func Registry(r registry.Registry) Option {
	return func(o *Options) {
		o.Registry = r
	}
}

// SetStrategy 设置节点策略引擎
func SetStrategy(fn Strategy) Option {
	return func(o *Options) {
		o.Strategy = fn
	}
}

// WithFilter 添加节点过滤规则
func WithFilter(fn ...Filter) SelectOption {
	return func(o *SelectOptions) {
		o.Filters = append(o.Filters, fn...)
	}
}

// WithStrategy 在调用select方法时候传入节点策略引擎
func WithStrategy(fn Strategy) SelectOption {
	return func(o *SelectOptions) {
		o.Strategy = fn
	}
}
