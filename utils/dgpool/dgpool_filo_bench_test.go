package dgpool_test

import (
	"context"
	"github.com/osgochina/dmicro/utils/dgpool"
	"sync"
	"testing"
)

func BenchmarkGoPool_MustGo(b *testing.B) {
	gp := dgpool.NewFILOPool(10000000, 0)
	wg := new(sync.WaitGroup)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		gp.MustGo(func() {})
	}
	wg.Wait()
}

func BenchmarkGoPool_MustGo_Background(b *testing.B) {
	gp := dgpool.NewFILOPool(10000000, 0)
	wg := new(sync.WaitGroup)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		gp.MustGo(func() {
			wg.Done()
		}, context.Background())
	}
	wg.Wait()
}

func BenchmarkGoPool_go(b *testing.B) {
	wg := new(sync.WaitGroup)
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			wg.Done()
		}()
	}
	wg.Wait()
}
