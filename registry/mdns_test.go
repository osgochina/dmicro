package registry

import (
	"github.com/gogf/gf/test/gtest"
	"testing"
	"time"
)

func TestMDNS(t *testing.T) {

	gtest.C(t, func(t *gtest.T) {
		testData := []*Service{
			{
				Name:    "test1",
				Version: "1.0.1",
				Nodes: []*Node{
					{
						Id:      "test1_1",
						Address: "10.10.0.1:10001",
						Metadata: map[string]string{
							"foo": "boo",
						},
					},
				},
			},
			{
				Name:    "test2",
				Version: "1.0.2",
				Nodes: []*Node{
					{
						Id:      "test2-1",
						Address: "10.0.0.2:10002",
						Metadata: map[string]string{
							"foo2": "bar2",
						},
					},
				},
			},
			{
				Name:    "test3",
				Version: "1.0.3",
				Nodes: []*Node{
					{
						Id:      "test3-1",
						Address: "10.0.0.3:10003",
						Metadata: map[string]string{
							"foo3": "bar3",
						},
					},
				},
			},
			{
				Name:    "test4",
				Version: "1.0.4",
				Nodes: []*Node{
					{
						Id:      "test4-1",
						Address: "[::]:10004",
						Metadata: map[string]string{
							"foo4": "bar4",
						},
					},
				},
			},
		}

		var opts []Option

		r := newRegistry(opts...)

		for _, service := range testData {
			if err := r.Register(service); err != nil {
				t.Fatal(err)
			}
			s, err := r.GetService(service.Name)
			if err != nil {
				t.Fatal(err)
			}
			if len(s) != 1 {
				t.Fatalf("Expected one result for %s got %d", service.Name, len(s))
			}

			if s[0].Name != service.Name {
				t.Fatalf("Expected name %s got %s", service.Name, s[0].Name)
			}
			if s[0].Version != service.Version {
				t.Fatalf("Expected version %s got %s", service.Version, s[0].Version)
			}

			if len(s[0].Nodes) != 1 {
				t.Fatalf("Expected 1 node, got %d", len(s[0].Nodes))
			}
			node := s[0].Nodes[0]
			if node.Id != service.Nodes[0].Id {
				t.Fatalf("Expected node id %s got %s", service.Nodes[0].Id, node.Id)
			}
			if node.Address != service.Nodes[0].Address {
				t.Fatalf("Expected node address %s got %s", service.Nodes[0].Address, node.Address)
			}
		}

		//time.Sleep(1 * time.Second)

		services, err := r.ListServices()

		if err != nil {
			t.Fatal(err)
		}

		for _, service := range testData {
			var seen bool
			for _, s := range services {
				if s.Name == service.Name {
					seen = true
					break
				}

			}

			if !seen {
				t.Fatalf("没有获取到预期的服务 %s", service.Name)
			}

			// 注销服务
			if err := r.Deregister(service); err != nil {
				t.Fatal(err)
			}

			time.Sleep(time.Millisecond * 5)

			// 再次检查服务是否存在
			s, _ := r.GetService(service.Name)
			if len(s) > 0 {
				t.Fatalf("Expected nothing got %+v", s[0])
			}
		}
	})

}
