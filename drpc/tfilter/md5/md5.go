package md5

import (
	"bytes"
	"crypto/md5"
	"errors"
	"github.com/osgochina/dmicro/drpc/tfilter"
)

const md5Length = 16

var errDataCheck = errors.New("check failed")

// Reg 注册md5校验过滤器
func Reg(id byte, name string) {
	tfilter.Reg(&md5Hash{
		id:   id,
		name: name,
	})
}

type md5Hash struct {
	id   byte
	name string
}

func (that *md5Hash) ID() byte {
	return that.id
}

func (that *md5Hash) Name() string {
	return that.name
}

func (that *md5Hash) OnPack(src []byte) ([]byte, error) {
	content, err := getMd5(src)
	if err != nil {
		return nil, err
	}
	src = append(src, content...)

	return src, nil
}

func (that *md5Hash) OnUnpack(src []byte) ([]byte, error) {
	srcLength := len(src)
	if srcLength < md5Length {
		return nil, errDataCheck
	}
	srcData := src[:srcLength-md5Length]
	content, err := getMd5(srcData)
	if err != nil {
		return nil, err
	}
	// Check
	if !bytes.Equal(content, src[srcLength-md5Length:]) {
		return nil, errDataCheck
	}
	return srcData, nil
}

func getMd5(src []byte) ([]byte, error) {
	newMd5 := md5.New()
	_, err := newMd5.Write(src)
	if err != nil {
		return nil, err
	}

	return newMd5.Sum(nil), nil
}
