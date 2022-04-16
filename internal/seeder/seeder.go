package seeder

import (
	"reflect"

	"github.com/hatlonely/go-kit/refx"
	"github.com/pkg/errors"
)

func RegisterSeeder(key string, constructor interface{}) {
	refx.Register(reflect.TypeOf((*Seeder)(nil)).Elem(), key, constructor)
}

func NewSeederWithOptions(options *refx.TypeOptions, opts ...refx.Option) (Seeder, error) {
	v, err := refx.New(reflect.TypeOf((*Seeder)(nil)).Elem(), options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "refx.New failed")
	}
	return v.(Seeder), nil
}

type Seeder interface {
	Seed() interface{}
}
