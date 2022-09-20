package memory

import (
	"fmt"
	"github.com/gogf/gf/test/gtest"
	"github.com/osgochina/dmicro/registry"
	"testing"
	"time"
)

var (
	testData = map[string][]*registry.Service{
		"foo": {
			{
				Name:    "foo",
				Version: "1.0.0",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.0-123",
						Address: "localhost:9999",
					},
					{
						Id:      "foo-1.0.0-321",
						Address: "localhost:9999",
					},
				},
			},
			{
				Name:    "foo",
				Version: "1.0.1",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.1-321",
						Address: "localhost:6666",
					},
				},
			},
			{
				Name:    "foo",
				Version: "1.0.3",
				Nodes: []*registry.Node{
					{
						Id:      "foo-1.0.3-345",
						Address: "localhost:8888",
					},
				},
			},
		},
	}
)

// 测试注册中心
func TestMemoryRegistry(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := NewRegistry()

		// 服务注册测试
		for _, v := range testData {
			serviceCount := 0
			for _, service := range v {
				err := m.Register(service)
				t.Assert(err, nil)
				serviceCount++
				services, err := m.GetService(service.Name)
				t.Assert(err, nil)
				t.Assert(len(services), serviceCount)
			}
		}

		// 获取service 与测试数据对比
		fn := func(serviceName string, services []*registry.Service) {
			scvs, err := m.GetService(serviceName)
			t.Assert(err, nil)
			t.Assert(len(scvs), len(services))

			for _, ss := range services {
				seen := false
				for _, s := range scvs {
					if s.Version == ss.Version {
						seen = true
						break
					}
				}
				t.Assert(seen, true)
			}
		}
		for k, v := range testData {
			fn(k, v)
		}

		// 对比所有已注册的服务,与测试数据对比
		services, err := m.ListServices()
		t.Assert(err, nil)
		totalServiceCount := 0
		for _, testSvc := range testData {
			for range testSvc {
				totalServiceCount++
			}
		}
		t.Assert(len(services), totalServiceCount)

		// 测试注销
		for _, v := range testData {
			for _, service := range v {
				err = m.Deregister(service)
				t.Assert(err, nil)
			}
		}
		// 服务已注销,返回的应该是错误
		for _, v := range testData {
			for _, service := range v {
				services, err = m.GetService(service.Name)
				t.Assert(err, registry.ErrNotFound)
				t.Assert(len(services), 0)
			}
		}
	})
}

func TestMemoryRegistryTTL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m := NewRegistry()
		for _, v := range testData {
			for _, service := range v {
				err := m.Register(service, registry.OptRegisterTTL(time.Millisecond))
				t.Assert(err, nil)
			}
		}
		time.Sleep(ttlPruneTime * 2)
		for name := range testData {
			svcs, err := m.GetService(name)
			t.Assert(err, nil)

			for _, svc := range svcs {
				t.Assert(len(svc.Nodes), 0)
			}
		}
	})
}

func TestMemoryRegistryTTLConcurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		concurrency := 1000
		waitTime := ttlPruneTime * 2
		m := NewRegistry()
		// 先注册服务
		for _, v := range testData {
			for _, service := range v {
				err := m.Register(service, registry.OptRegisterTTL(waitTime/2))
				t.Assert(err, nil)
			}
		}
		errChan := make(chan error, concurrency)
		syncChan := make(chan struct{})
		//并发读取
		for i := 0; i < concurrency; i++ {
			go func() {
				<-syncChan
				for name := range testData {
					svcs, err := m.GetService(name)
					if err != nil {
						errChan <- err
						return
					}
					for _, svc := range svcs {
						if len(svc.Nodes) > 0 {
							errChan <- fmt.Errorf("Service %q still has nodes registered ", name)
							return
						}
					}
				}
				errChan <- nil
			}()
		}
		time.Sleep(waitTime)
		close(syncChan)

		for i := 0; i < concurrency; i++ {
			err := <-errChan
			t.Assert(err, nil)
		}
	})
}
