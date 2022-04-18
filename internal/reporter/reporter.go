package reporter

import (
	"reflect"

	"github.com/hatlonely/go-ben/internal/stat"

	"github.com/hatlonely/go-kit/refx"
	"github.com/pkg/errors"
)

func RegisterReporter(key string, constructor interface{}) {
	refx.Register(reflect.TypeOf((*Reporter)(nil)).Elem(), key, constructor)
}

func NewReporterWithOptions(options *refx.TypeOptions, opts ...refx.Option) (Reporter, error) {
	v, err := refx.New(reflect.TypeOf((*Reporter)(nil)).Elem(), options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "refx.New failed")
	}
	return v.(Reporter), nil
}

type Reporter interface {
	Report(test *stat.TestStat) string
}
