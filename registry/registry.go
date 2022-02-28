package registry

import (
	"errors"
	"time"
)

const (
	DefaultRegisterInterval = time.Second * 30
	DefaultRegisterTTL      = time.Second * 90
)

var ErrNotFound = errors.New("service not found")

// Registry 服务注册接口
type Registry interface {
	Init(...Option) error
	Options() Options
	Register(*Service, ...RegisterOption) error
	Deregister(*Service, ...DeregisterOption) error
	GetService(string, ...GetOption) ([]*Service, error)
	ListServices(...ListOption) ([]*Service, error)
	Watch(...WatchOption) (Watcher, error)
	Stop()
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
