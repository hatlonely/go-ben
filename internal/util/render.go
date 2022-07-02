package util

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/PaesslerAG/gval"
	"github.com/generikvault/gvalstrings"
	"github.com/hatlonely/go-kit/cast"
	"github.com/hatlonely/go-kit/refx"
	"github.com/pkg/errors"
)

var Lang = gval.NewLanguage(
	gval.Arithmetic(),
	gval.Bitmask(),
	gval.Text(),
	gval.PropositionalLogic(),
	gval.JSON(),
	gvalstrings.SingleQuoted(),
	gval.Function("date", func(arguments ...interface{}) (interface{}, error) {
		if len(arguments) != 1 {
			return nil, fmt.Errorf("date() expects exactly one string argument")
		}
		s, ok := arguments[0].(string)
		if !ok {
			return nil, fmt.Errorf("date() expects exactly one string argument")
		}
		for _, format := range [...]string{
			time.ANSIC,
			time.UnixDate,
			time.RubyDate,
			time.Kitchen,
			time.RFC3339,
			time.RFC3339Nano,
			"2006-01-02",                         // RFC 3339
			"2006-01-02 15:04",                   // RFC 3339 with minutes
			"2006-01-02 15:04:05",                // RFC 3339 with seconds
			"2006-01-02 15:04:05-07:00",          // RFC 3339 with seconds and timezone
			"2006-01-02T15Z0700",                 // ISO8601 with hour
			"2006-01-02T15:04Z0700",              // ISO8601 with minutes
			"2006-01-02T15:04:05Z0700",           // ISO8601 with seconds
			"2006-01-02T15:04:05.999999999Z0700", // ISO8601 with nanoseconds
		} {
			ret, err := time.ParseInLocation(format, s, time.Local)
			if err == nil {
				return ret, nil
			}
		}
		return nil, fmt.Errorf("date() could not parse %s", s)
	}),
	gval.Function("len", func(x interface{}) (int, error) {
		return len(x.(string)), nil
	}),
	gval.Function("int", func(x interface{}) (int, error) {
		switch v := x.(type) {
		case string:
			return cast.ToIntE(strings.Fields(v)[0])
		}
		return cast.ToIntE(x)
	}),
)

func Render(v interface{}, args interface{}) (interface{}, error) {
	var w interface{}
	if err := refx.InterfaceTravel(v, func(key string, val interface{}) error {
		if strings.HasPrefix(key, "#") {
			str, ok := val.(string)
			if !ok {
				return errors.Errorf("val [%v] is not string. key [%v]", val, key)
			}
			eval, err := Lang.NewEvaluable(str)
			if err != nil {
				return errors.Wrap(err, "lang.NewEvaluable failed")
			}
			v, err := eval(context.Background(), args)
			if err != nil {
				return errors.Wrap(err, "eval.EvalString failed")
			}
			if err := refx.InterfaceSet(&w, key[1:], v); err != nil {
				return errors.WithMessage(err, "refx.InterfaceSet failed")
			}
		} else {
			if err := refx.InterfaceSet(&w, key, val); err != nil {
				return errors.WithMessage(err, "refx.InterfaceSet failed")
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return w, nil
}
