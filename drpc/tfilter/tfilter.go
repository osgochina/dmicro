package tfilter

import (
	"errors"
	"fmt"
	"math"
)

// TransferFilter 传输过滤器接口
type TransferFilter interface {
	// ID 过滤器id
	ID() byte
	// Name 过滤器名字
	Name() string
	// OnPack 过滤器打包方法
	OnPack([]byte) ([]byte, error)
	// OnUnpack 过滤器解包方法
	OnUnpack([]byte) ([]byte, error)
}

var transferFilterMap = struct {
	idMap   map[byte]TransferFilter
	nameMap map[string]TransferFilter
}{
	idMap:   make(map[byte]TransferFilter),
	nameMap: make(map[string]TransferFilter),
}

var ErrTransferFilterTooLong = errors.New("The length of transfer pipe cannot be bigger than 255 ")

// Reg 注册过滤器
func Reg(tFilter TransferFilter) {
	id := tFilter.ID()
	name := tFilter.Name()
	if _, ok := transferFilterMap.idMap[id]; ok {
		panic(fmt.Sprintf("multi-register transfer filter id: %d", tFilter.ID()))
	}
	if _, ok := transferFilterMap.nameMap[name]; ok {
		panic("multi-register transfer filter name: " + tFilter.Name())
	}
	transferFilterMap.idMap[id] = tFilter
	transferFilterMap.nameMap[name] = tFilter
}

// Get 通过id获取过滤器对象
func Get(id byte) (TransferFilter, error) {
	tFilter, ok := transferFilterMap.idMap[id]
	if !ok {
		return nil, fmt.Errorf("unsupported transfer filter id: %d", id)
	}
	return tFilter, nil
}

// GetByName 通过过滤器名称返回过滤器对象
func GetByName(name string) (TransferFilter, error) {
	tFilter, ok := transferFilterMap.nameMap[name]
	if !ok {
		return nil, fmt.Errorf("unsupported transfer filter name: %s", name)
	}
	return tFilter, nil
}

// PipeTFilter 传输过滤器管道
type PipeTFilter struct {
	filters []TransferFilter
}

// NewPipeTFilter 创建传输过滤器管道
func NewPipeTFilter() *PipeTFilter {
	return new(PipeTFilter)
}

// Reset 清除传输过滤器管道
func (that *PipeTFilter) Reset() {
	that.filters = that.filters[:0]
}

// Append 追加传输过滤器
func (that *PipeTFilter) Append(filterID ...byte) error {
	for _, id := range filterID {
		filter, err := Get(id)
		if err != nil {
			return err
		}
		that.filters = append(that.filters, filter)
	}
	return that.check()
}

// AppendFrom 从指定传输过滤器管道追加过滤器
func (that *PipeTFilter) AppendFrom(src *PipeTFilter) {
	for _, filter := range src.filters {
		that.filters = append(that.filters, filter)
	}
}

//检查传输过滤器管道是否大于256
func (that *PipeTFilter) check() error {
	if that.Len() > math.MaxUint8 {
		return ErrTransferFilterTooLong
	}
	return nil
}

// Len 当前传输过滤器管道的长度
func (that *PipeTFilter) Len() int {
	if that == nil {
		return 0
	}
	return len(that.filters)
}

// IDs 获取当前传输过滤器管道中的id
func (that *PipeTFilter) IDs() []byte {
	var ids = make([]byte, that.Len())
	if that.Len() == 0 {
		return ids
	}
	for i, filter := range that.filters {
		ids[i] = filter.ID()
	}
	return ids
}

//获取当前传输过滤器管道中的名字
func (that *PipeTFilter) Names() []string {
	var names = make([]string, that.Len())
	if that.Len() == 0 {
		return names
	}
	for i, filter := range that.filters {
		names[i] = filter.Name()
	}
	return names
}

// Iterator 迭代
func (that *PipeTFilter) Iterator(callback func(idx int, filter TransferFilter) bool) {
	for idx, filter := range that.filters {
		if !callback(idx, filter) {
			break
		}
	}
}

// OnPack 打包，从最内层到最外层
func (that *PipeTFilter) OnPack(data []byte) ([]byte, error) {
	var err error
	for i := that.Len() - 1; i >= 0; i-- {
		if data, err = that.filters[i].OnPack(data); err != nil {
			return data, err
		}
	}
	return data, err
}

// OnUnpack 解包，从最外层到最内层
func (that *PipeTFilter) OnUnpack(data []byte) ([]byte, error) {
	var err error
	var count = that.Len()
	for i := 0; i < count; i++ {
		if data, err = that.filters[i].OnUnpack(data); err != nil {
			return data, err
		}
	}
	return data, err
}
