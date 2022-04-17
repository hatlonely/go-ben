package util

import (
	"encoding/json"

	"github.com/hatlonely/go-kit/refx"
	"github.com/pkg/errors"
)

func MustMerge(v1 interface{}, vs ...interface{}) interface{} {
	v, err := Merge(v1, vs...)
	refx.Must(err)
	return v
}

func Merge(v1 interface{}, vs ...interface{}) (interface{}, error) {
	buf, err := json.Marshal(v1)
	if err != nil {
		return nil, errors.Wrap(err, "json.Marshal failed")
	}
	var v interface{}
	if err := json.Unmarshal(buf, v); err != nil {
		return nil, errors.Wrap(err, "json.Unmarshal failed")
	}

	for _, v2 := range vs {
		if err := refx.Merge(v, v2); err != nil {
			return nil, errors.WithMessage(err, "refx.Merge failed")
		}
	}

	return v, nil
}
