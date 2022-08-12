package memory

import (
	"context"
	"github.com/osgochina/dmicro/registry"
	"time"
)

// memory 特殊的节点结构
type node struct {
	*registry.Node
	TTL      time.Duration
	LastSeen time.Time
}

// 记录器
type record struct {
	Name     string
	Version  string
	Metadata map[string]string
	Nodes    map[string]*node
}

//把service对象转换成record对象
func serviceToRecord(s *registry.Service, ttl time.Duration) *record {
	metadata := make(map[string]string, len(s.Metadata))
	for k, v := range s.Metadata {
		metadata[k] = v
	}
	nodes := make(map[string]*node, len(s.Nodes))

	for _, n := range s.Nodes {
		nodes[n.Id] = &node{
			Node:     n,
			TTL:      ttl,
			LastSeen: time.Now(),
		}
	}
	return &record{
		Name:     s.Name,
		Version:  s.Version,
		Metadata: metadata,
		Nodes:    nodes,
	}
}

// 通过记录转换成service对象
func recordToService(r *record) *registry.Service {
	metadata := make(map[string]string, len(r.Metadata))
	for k, v := range r.Metadata {
		metadata[k] = v
	}
	nodes := make([]*registry.Node, len(r.Nodes))
	i := 0
	for _, n := range r.Nodes {
		md := make(map[string]string, len(n.Metadata))
		for k, v := range n.Metadata {
			md[k] = v
		}

		nodes[i] = &registry.Node{
			Id:       n.Id,
			Address:  n.Address,
			Metadata: metadata,
		}
		i++
	}
	return &registry.Service{
		Name:     r.Name,
		Version:  r.Version,
		Metadata: metadata,
		Nodes:    nodes,
	}
}

type ServicesKey struct{}

// 通过上下文参数,获取服务的记录
func getServiceRecords(ctx context.Context) map[string]map[string]*record {
	memServices, ok := ctx.Value(ServicesKey{}).(map[string][]*registry.Service)
	if !ok {
		return nil
	}
	services := make(map[string]map[string]*record)

	for name, svc := range memServices {
		if _, ok = services[name]; !ok {
			services[name] = make(map[string]*record)
		}
		for _, s := range svc {
			services[s.Name][s.Version] = serviceToRecord(s, 0)
		}
	}
	return services
}
