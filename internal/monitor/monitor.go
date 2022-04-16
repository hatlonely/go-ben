package monitor

import (
	"reflect"
	"time"

	"github.com/hatlonely/go-kit/refx"
	"github.com/pkg/errors"

	"github.com/hatlonely/go-ben/internal/stat"
)

func RegisterMonitor(key string, constructor interface{}) {
	refx.Register(reflect.TypeOf((*Monitor)(nil)).Elem(), key, constructor)
}

func NewMonitorWithOptions(options *refx.TypeOptions, opts ...refx.Option) (Monitor, error) {
	v, err := refx.New(reflect.TypeOf((*Monitor)(nil)).Elem(), options, opts...)
	if err != nil {
		return nil, errors.WithMessage(err, "refx.New failed")
	}
	return v.(Monitor), nil
}

type Monitor interface {
	Collect()
	Stat(start time.Time, end time.Time) *stat.MonitorStat
}
