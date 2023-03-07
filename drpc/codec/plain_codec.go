package codec

import (
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"reflect"
	"strconv"
)

const (
	PlainName = "plain"
	PlainId   = 's'
)

func init() {
	Reg(new(PlainCodec))
}

// PlainCodec plain text codec
type PlainCodec struct{}

// Name returns codec name.
func (PlainCodec) Name() string {
	return PlainName
}

// ID returns codec id.
func (PlainCodec) ID() byte {
	return PlainId
}

func (PlainCodec) Marshal(v interface{}) ([]byte, error) {
	var b []byte
	switch vv := v.(type) {
	case nil:
	case string:
		b = gconv.Bytes(vv)
	case *string:
		b = gconv.Bytes(*vv)
	case []byte:
		b = vv
	case *[]byte:
		b = *vv
	default:
		s, ok := formatProperType(reflect.ValueOf(v))
		if !ok {
			return nil, fmt.Errorf("plain codec: %T can not be directly converted to []byte type", v)
		}
		b = gconv.Bytes(s)
	}
	return b, nil
}

func (PlainCodec) Unmarshal(data []byte, v interface{}) error {
	switch s := v.(type) {
	case nil:
		return nil
	case *string:
		*s = string(data)
	case []byte:
		copy(s, data)
	case *[]byte:
		if length := len(data); cap(*s) < length {
			*s = make([]byte, length)
		} else {
			*s = (*s)[:length]
		}
		copy(*s, data)
	default:
		if !parseProperType(data, reflect.ValueOf(v)) {
			return fmt.Errorf("plain codec: []byte can not be directly converted to %T type", v)
		}
	}
	return nil
}

func parseProperType(data []byte, v reflect.Value) bool {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.CanSet() {
		return false
	}
	s := gconv.String(data)
	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		bol, err := strconv.ParseBool(s)
		if err != nil {
			return false
		}
		v.SetBool(bol)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		d, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return false
		}
		v.SetInt(d)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		d, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return false
		}
		v.SetUint(d)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return false
		}
		v.SetFloat(f)
	case reflect.Slice:
		if v.Type().Elem().Kind() != reflect.Uint8 {
			return false
		}
		v.SetBytes(data)
	case reflect.Invalid:
		return true
	default:
		return false
	}
	return true
}
