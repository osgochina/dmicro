package codec

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/util/gconv"
	"net/url"
	"reflect"
	"strconv"
	"time"
)

const (
	FormName = "form"
	FormId   = 'f'
)

func init() {
	Reg(new(FormCodec))
}

// FormCodec url encode codec
type FormCodec struct{}

// Name returns codec name.
func (FormCodec) Name() string {
	return FormName
}

// ID returns codec id.
func (FormCodec) ID() byte {
	return FormId
}

func (FormCodec) Marshal(v interface{}) ([]byte, error) {
	var b []byte
	switch vv := v.(type) {
	case nil:
	case url.Values:
		b = gconv.Bytes(vv.Encode())
	case *url.Values:
		b = gconv.Bytes(vv.Encode())
	case map[string][]string:
		b = gconv.Bytes((url.Values)(vv).Encode())
	case *map[string][]string:
		b = gconv.Bytes((url.Values)(*vv).Encode())
	default:
		vvv := reflect.ValueOf(v)
		for vvv.Kind() == reflect.Ptr {
			vvv = vvv.Elem()
		}
		if vvv.Kind() == reflect.Struct {
			q := make(url.Values)
			setStructToForm(q, vvv)
			return gconv.Bytes(q.Encode()), nil
		}
		return nil, fmt.Errorf("form codec: %T can not be encoded to urlencoded string", v)
	}
	return b, nil
}

func setStructToForm(q url.Values, val reflect.Value) {
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		inputFieldName := typeField.Tag.Get(FormName)
		if inputFieldName == "" {
			if structField.Kind() == reflect.Struct {
				setStructToForm(q, structField)
				continue
			}
			inputFieldName = typeField.Name
		}
		a, ok := q[inputFieldName]
		if !ok {
			a = make([]string, 0, 1)
		}
		if structField.Kind() == reflect.Slice || structField.Kind() == reflect.Array {
			for i := structField.Len() - 1; i >= 0; i-- {
				if s, ok := formatProperType(structField.Index(i)); ok {
					a = append(a, s)
				}
			}
		} else if s, ok := formatProperType(structField); ok {
			a = append(a, s)
		}
		q[inputFieldName] = a
	}
}

func formatProperType(v reflect.Value) (string, bool) {
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.String:
		return v.String(), true
	case reflect.Bool:
		return strconv.FormatBool(v.Bool()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(v.Uint(), 10), true
	case reflect.Float32:
		return strconv.FormatFloat(v.Float(), 'f', -1, 32), true
	case reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'f', -1, 64), true
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return gconv.String(v.Bytes()), true
		} else {
			return "", false
		}
	case reflect.Invalid:
		return "", true
	default:
		return "", false
	}
}

func (FormCodec) Unmarshal(data []byte, v interface{}) error {
	form, err := url.ParseQuery(gconv.String(data))
	if err != nil {
		return fmt.Errorf("form codec: %s", err.Error())
	}
	switch vv := v.(type) {
	case nil:
	case *url.Values:
		*vv = form
	case *map[string][]string:
		*vv = form
	case *interface{}:
		*vv = form
	default:
		vvv := reflect.ValueOf(v)
		for vvv.Kind() == reflect.Ptr {
			vvv = vvv.Elem()
		}
		switch vvv.Kind() {
		case reflect.Interface:
			// *interface{}
			vvv.Set(reflect.ValueOf(form))
		case reflect.Struct:
			return mapFormToStruct(vvv, form)
		}
		return fmt.Errorf("plain codec: []byte can not be converted to %T type", v)
	}
	return nil
}

func mapFormToStruct(val reflect.Value, form map[string][]string) error {
	for val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()
	for i := 0; i < typ.NumField(); i++ {
		typeField := typ.Field(i)
		structField := val.Field(i)
		if !structField.CanSet() {
			continue
		}

		structFieldKind := structField.Kind()
		inputFieldName := typeField.Tag.Get(FormName)
		if inputFieldName == "" {
			// if NAME_FORM tag is nil, we inspect if the field is a struct.
			// this would not make sense for JSON parsing but it does for a form
			// since data is flatten
			if structFieldKind == reflect.Struct {
				err := mapFormToStruct(structField, form)
				if err != nil {
					return err
				}
				continue
			}
			inputFieldName = typeField.Name
		}
		inputValue, exists := form[inputFieldName]
		if !exists {
			continue
		}

		numElems := len(inputValue)
		if structFieldKind == reflect.Array && numElems > 0 {
			for i := 0; i < numElems; i++ {
				arrayOf := structField.Type().Elem().Kind()
				if err := setWithProperType(arrayOf, inputValue[i], structField.Index(i)); err != nil {
					return err
				}
			}
		} else if structFieldKind == reflect.Slice && numElems > 0 {
			sliceOf := structField.Type().Elem().Kind()
			slice := reflect.MakeSlice(structField.Type(), numElems, numElems)
			for i := 0; i < numElems; i++ {
				if err := setWithProperType(sliceOf, inputValue[i], slice.Index(i)); err != nil {
					return err
				}
			}
			val.Field(i).Set(slice)
		} else {
			if _, isTime := structField.Interface().(time.Time); isTime {
				if err := setTimeField(inputValue[0], typeField, structField); err != nil {
					return err
				}
				continue
			}
			if err := setWithProperType(typeField.Type.Kind(), inputValue[0], structField); err != nil {
				return err
			}
		}
	}
	return nil
}

func setWithProperType(valueKind reflect.Kind, val string, structField reflect.Value) error {
	switch valueKind {
	case reflect.Int:
		return setIntField(val, 0, structField)
	case reflect.Int8:
		return setIntField(val, 8, structField)
	case reflect.Int16:
		return setIntField(val, 16, structField)
	case reflect.Int32:
		return setIntField(val, 32, structField)
	case reflect.Int64:
		return setIntField(val, 64, structField)
	case reflect.Uint:
		return setUintField(val, 0, structField)
	case reflect.Uint8:
		return setUintField(val, 8, structField)
	case reflect.Uint16:
		return setUintField(val, 16, structField)
	case reflect.Uint32:
		return setUintField(val, 32, structField)
	case reflect.Uint64:
		return setUintField(val, 64, structField)
	case reflect.Bool:
		return setBoolField(val, structField)
	case reflect.Float32:
		return setFloatField(val, 32, structField)
	case reflect.Float64:
		return setFloatField(val, 64, structField)
	case reflect.String:
		structField.SetString(val)
	default:
		return errors.New("Unknown type ")
	}
	return nil
}

func setIntField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	intVal, err := strconv.ParseInt(val, 10, bitSize)
	if err == nil {
		field.SetInt(intVal)
	}
	return err
}

func setUintField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0"
	}
	uintVal, err := strconv.ParseUint(val, 10, bitSize)
	if err == nil {
		field.SetUint(uintVal)
	}
	return err
}

func setBoolField(val string, field reflect.Value) error {
	if val == "" {
		val = "false"
	}
	boolVal, err := strconv.ParseBool(val)
	if err == nil {
		field.SetBool(boolVal)
	}
	return nil
}

func setFloatField(val string, bitSize int, field reflect.Value) error {
	if val == "" {
		val = "0.0"
	}
	floatVal, err := strconv.ParseFloat(val, bitSize)
	if err == nil {
		field.SetFloat(floatVal)
	}
	return err
}

func setTimeField(val string, structField reflect.StructField, value reflect.Value) error {
	timeFormat := structField.Tag.Get("time_format")
	if timeFormat == "" {
		return errors.New("Blank time format")
	}

	if val == "" {
		value.Set(reflect.ValueOf(time.Time{}))
		return nil
	}

	l := time.Local
	if isUTC, _ := strconv.ParseBool(structField.Tag.Get("time_utc")); isUTC {
		l = time.UTC
	}

	if locTag := structField.Tag.Get("time_location"); locTag != "" {
		loc, err := time.LoadLocation(locTag)
		if err != nil {
			return err
		}
		l = loc
	}

	t, err := time.ParseInLocation(timeFormat, val, l)
	if err != nil {
		return err
	}

	value.Set(reflect.ValueOf(t))
	return nil
}
