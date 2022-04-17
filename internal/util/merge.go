package util

import (
	"encoding/json"

	"github.com/hatlonely/go-kit/refx"
	"github.com/pkg/errors"
)

func MustMerge(vs ...interface{}) interface{} {
	v, err := Merge(vs...)
	refx.Must(err)
	return v
}

func Merge(vs ...interface{}) (interface{}, error) {
	if len(vs) == 0 {
		return nil, nil
	}
	if vs[0] == nil {
		return Merge(vs[1:]...)
	}
	buf, err := json.Marshal(vs[0])
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal failed")
	}
	var v interface{}
	if err := json.Unmarshal(buf, v); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal failed")
	}

	for _, v2 := range vs[1:] {
		if err := refx.Merge(v, v2); err != nil {
			return nil, errors.WithMessage(err, "refx.Merge failed")
		}
	}

	return v, nil
}
