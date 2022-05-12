package eval

import (
	"github.com/Knetic/govaluate"
	"github.com/pkg/errors"
)

func NewGovaluateEvaluable(expr string) (*GovaluateEvaluable, error) {
	e, err := govaluate.NewEvaluableExpression(expr)
	if err != nil {
		return nil, errors.Wrap(err, "govaluate.NewEvaluableExpression failed")
	}

	return &GovaluateEvaluable{
		expr: e,
	}, nil
}

type GovaluateEvaluable struct {
	expr *govaluate.EvaluableExpression
}

func (e *GovaluateEvaluable) Evaluate(v interface{}) (interface{}, error) {
	return e.expr.Evaluate(v.(map[string]interface{}))
}
