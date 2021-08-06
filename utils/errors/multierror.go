package errors

import (
	"bytes"
	"strconv"
)

type multiError struct {
	errs []error
	text string
}

var mergePrefix = []byte("MultiError:\n")

func (that *multiError) Error() string {
	if len(that.text) > 0 {
		return that.text
	}
	var bText = make([]byte, len(mergePrefix), 56)
	copy(bText, mergePrefix)
	for i, err := range that.errs {
		bText = append(bText, strconv.Itoa(i+1)...)
		bText = append(bText, ". "...)
		bText = append(bText, bytes.Trim([]byte(err.Error()), "\n")...)
		bText = append(bText, '\n')
	}
	that.text = string(bText)
	return that.text
}

func Merge(errs ...error) error {
	return Append(nil, errs...)
}

func Append(err error, errs ...error) error {
	count := len(errs)
	if count == 0 {
		return err
	}
	var merged []error
	if err != nil {
		if e, ok := err.(*multiError); ok {
			_count := len(e.errs)
			merged = make([]error, _count, count+_count)
			copy(merged, e.errs)
		} else {
			merged = make([]error, 1, count+1)
			merged[0] = err
		}
	}
	for _, err := range errs {
		switch e := err.(type) {
		case nil:
			continue
		case *multiError:
			merged = append(merged, e.errs...)
		default:
			merged = append(merged, e)
		}
	}
	if len(merged) == 0 {
		return nil
	}
	return &multiError{
		errs: merged,
	}
}
