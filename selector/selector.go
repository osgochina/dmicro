package selector

import (
	"errors"
	"github.com/osgochina/dmicro/registry"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrNoneAvailable = errors.New("none available")
)

// Selector 服务选择器接口
type Selector interface {
	Init(opts ...Option) error
	Options() Options
	Select(service string, opts ...SelectOption) (Next, error)
	Mark(service string, node *registry.Node, err error)
	Reset(service string)
	Close() error
	String() string
}

// Next 获取可用的节点
type Next func() (*registry.Node, error)

// Filter 过滤节点
type Filter func([]*registry.Service) []*registry.Service

// Strategy 根据策略选择节点
type Strategy func([]*registry.Service) Next
