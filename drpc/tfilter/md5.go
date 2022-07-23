package tfilter

import (
	"bytes"
	"crypto/md5"
	"errors"
)

const md5Length = 16

const (
	Md5Id   = '5'
	Md5Name = "md5"
)

var errDataCheck = errors.New("check failed")

// RegMD5 注册md5校验过滤器
func RegMD5() {
	Reg(&md5Hash{
		id:   Md5Id,
		name: Md5Name,
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
