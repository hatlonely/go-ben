package refcli

import (
	"reflect"

	"github.com/hatlonely/go-kit/refx"
	"github.com/pkg/errors"
)

func RegisterClient(key string, constructor interface{}) {
	refx.Register(reflect.TypeOf((*Client)(nil)).Elem(), key, constructor)
}

func NewClientWithOptions(options *refx.TypeOptions, opts ...refx.Option) (Client, error) {
	v, err := refx.New(reflect.TypeOf((*Client)(nil)).Elem(), options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "refx.New failed")
	}
	return v.(Client), nil
}

type Client interface {
	Name() string
	Do(req interface{}) (interface{}, error)
}
