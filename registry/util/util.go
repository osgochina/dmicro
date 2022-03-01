package util

import "github.com/osgochina/dmicro/registry"

// Copy 复制services列表
func Copy(current []*registry.Service) []*registry.Service {
	services := make([]*registry.Service, len(current))

	for i, service := range current {
		services[i] = CopyService(service)
	}
	return services
}

func CopyService(service *registry.Service) *registry.Service {
	s := new(registry.Service)
	*s = *service

	nodes := make([]*registry.Node, len(service.Nodes))
	for i, node := range service.Nodes {
		n := new(registry.Node)
		*n = *node
		nodes[i] = n
	}
	s.Nodes = nodes

	return s
}
