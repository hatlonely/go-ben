package eval

import (
	"context"

	"github.com/d5/tengo/v2"
)

func NewTengoEvaluable(expr string) (*TengoEvaluable, error) {
	return &TengoEvaluable{
		expr: expr,
	}, nil
}

type TengoEvaluable struct {
	expr string
}

func (e *TengoEvaluable) Evaluate(v interface{}) (interface{}, error) {
	return tengo.Eval(context.Background(), e.expr, v.(map[string]interface{}))
}
