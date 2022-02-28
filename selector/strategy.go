package selector

import (
	"github.com/osgochina/dmicro/registry"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Random 随机选择节点
func Random(services []*registry.Service) Next {
	nodes := make([]*registry.Node, 0, len(services))

	for _, service := range services {
		nodes = append(nodes, service.Nodes...)
	}

	return func() (*registry.Node, error) {
		if len(nodes) == 0 {
			return nil, ErrNoneAvailable
		}
		i := rand.Int() % len(nodes)
		return nodes[i], nil
	}
}
