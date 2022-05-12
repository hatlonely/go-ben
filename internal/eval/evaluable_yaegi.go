package eval

import (
	"github.com/pkg/errors"
	"github.com/traefik/yaegi/interp"
	"github.com/traefik/yaegi/stdlib"
)

func NewYaegiEvaluable(expr string) (*YaegiEvaluable, error) {
	i := interp.New(interp.Options{})
	err := i.Use(stdlib.Symbols)
	if err != nil {
		return nil, errors.Wrap(err, "i.Use stdlib.Symbols failed")
	}
	_, err = i.Eval(`import "fmt"`)
	if err != nil {
		return nil, errors.Wrap(err, "i.Eval failed")
	}

	_, err = i.Eval(expr)
	if err != nil {
		return nil, errors.Wrap(err, "i.Eval failed")
	}

	v, err := i.Eval("eval")
	if err != nil {
		return nil, errors.Wrap(err, "i.Eval failed")
	}

	fun := v.Interface().(func(map[string]interface{}) (interface{}, error))

	return &YaegiEvaluable{
		fun: fun,
	}, nil
}

type YaegiEvaluable struct {
	fun func(map[string]interface{}) (interface{}, error)
}

func (e *YaegiEvaluable) Evaluate(v interface{}) (interface{}, error) {
	return e.fun(v.(map[string]interface{}))
}
