package registry

import (
	"errors"
	"time"
)

const (
	DefaultRegisterInterval = time.Second * 30
	DefaultRegisterTTL      = time.Second * 90
)

var (
	DefaultRegistry   = NewRegistry()
	ErrNotFound       = errors.New("service not found")
	ErrWatcherStopped = errors.New("watcher stopped")
)

// Registry 服务注册接口
type Registry interface {
	Init(...Option) error
	Options() Options
	Register(*Service, ...RegisterOption) error
	Deregister(*Service, ...DeregisterOption) error
	GetService(string, ...GetOption) ([]*Service, error)
	ListServices(...ListOption) ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	String() string
}

// Service 服务
type Service struct {
	Name     string            `json:"name"`
	Version  string            `json:"version"`
	Metadata map[string]string `json:"metadata"`
	Nodes    []*Node           `json:"nodes"`
}

// Node 节点
type Node struct {
	Id       string            `json:"id"`
	Address  string            `json:"address"`
	Paths    []string          `json:"paths"`
	Metadata map[string]string `json:"metadata"`
}

type Option func(*Options)

type RegisterOption func(*RegisterOptions)

type WatchOption func(*WatchOptions)

type DeregisterOption func(*DeregisterOptions)

type GetOption func(*GetOptions)

type ListOption func(*ListOptions)

// Register 注册节点到注册表
func Register(s *Service, opts ...RegisterOption) error {
	return DefaultRegistry.Register(s, opts...)
}

// Deregister 从注册表这中注销节点
func Deregister(s *Service, opts ...DeregisterOption) error {
	return DefaultRegistry.Deregister(s, opts...)
}

// GetService 通过服务名获取服务列表,为什么返回的是一个数组？因为服务是有版本区别的，多个版本的服务是可以共存的
func GetService(name string) ([]*Service, error) {
	return DefaultRegistry.GetService(name)
}

// ListServices 获取所有服务的列表
func ListServices() ([]*Service, error) {
	return DefaultRegistry.ListServices()
}

// Watch 监听服务变化
func Watch(opts ...WatchOption) (Watcher, error) {
	return DefaultRegistry.Watch(opts...)
}

func String() string {
	return DefaultRegistry.String()
}
