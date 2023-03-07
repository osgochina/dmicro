package tfilter

import "github.com/gogf/gf/v2/crypto/gaes"

const (
	AesId   = 'a'
	AesName = "aes"
)

type aesHash struct {
	id   byte
	name string
	key  []byte
}

func RegAES(key []byte) {
	Reg(&aesHash{
		id:   AesId,
		name: AesName,
		key:  key,
	})
}
func (that *aesHash) ID() byte {
	return that.id
}

func (that *aesHash) Name() string {
	return that.name
}

func (that *aesHash) OnPack(src []byte) ([]byte, error) {
	content, err := gaes.Encrypt(src, that.key)
	if err != nil {
		return nil, err
	}
	return content, nil
}

func (that *aesHash) OnUnpack(src []byte) ([]byte, error) {
	content, err := gaes.Decrypt(src, that.key)
	if err != nil {
		return nil, err
	}
	return content, nil
}
