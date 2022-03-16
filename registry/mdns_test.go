package registry

import (
	"github.com/gogf/gf/test/gtest"
	"testing"
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
		}

		var opts []Option

		r := newRegistry(opts...)

		for _, service := range testData {
			if err := r.Register(service); err != nil {
				t.Fatal(err)
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
		}
	})

}
