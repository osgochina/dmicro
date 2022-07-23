package tfilter

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"sync"
)

const (
	GzipId   = 'z'
	GzipName = "gzip"
)

// RegGzip registers a gzip filter for transfer.
func RegGzip(level int) {
	Reg(newGzip(GzipId, GzipName, level))
}

type Gzip struct {
	id    byte
	name  string
	level int
	wPool sync.Pool
	rPool sync.Pool
}

func newGzip(id byte, name string, level int) *Gzip {
	if level < gzip.HuffmanOnly || level > gzip.BestCompression {
		panic(fmt.Sprintf("gzip: invalid compression level: %d", level))
	}
	g := new(Gzip)
	g.level = level
	g.id = id
	g.name = name
	g.wPool = sync.Pool{
		New: func() interface{} {
			gw, _ := gzip.NewWriterLevel(nil, g.level)
			return gw
		},
	}
	g.rPool = sync.Pool{
		New: func() interface{} {
			return new(gzip.Reader)
		},
	}
	return g
}

// ID returns transfer filter id.
func (that *Gzip) ID() byte {
	return that.id
}

// Name  returns transfer filter name.
func (that *Gzip) Name() string {
	return that.name
}

// OnPack performs filtering on packing.
func (that *Gzip) OnPack(src []byte) ([]byte, error) {
	var bb = new(bytes.Buffer)
	gw := that.wPool.Get().(*gzip.Writer)
	gw.Reset(bb)
	_, _ = gw.Write(src)
	err := gw.Close()
	gw.Reset(nil)
	that.wPool.Put(gw)
	if err != nil {
		return nil, err
	}
	return bb.Bytes(), nil
}

// OnUnpack performs filtering on unpacking.
func (that *Gzip) OnUnpack(src []byte) (dest []byte, err error) {
	if len(src) == 0 {
		return src, nil
	}
	gr := that.rPool.Get().(*gzip.Reader)
	err = gr.Reset(bytes.NewReader(src))
	if err == nil {
		dest, err = ioutil.ReadAll(gr)
	}
	_ = gr.Close()
	that.rPool.Put(gr)
	return dest, err
}
