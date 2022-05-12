package eval

import "github.com/pkg/errors"

func NewEvaluable(typ string, expr string) (Evaluable, error) {
	switch typ {
	case "Tengo":
		return NewTengoEvaluable(expr)
	case "Yaegi":
		return NewYaegiEvaluable(expr)
	case "Govaluate":
		return NewGovaluateEvaluable(expr)
	case "Gvalstrings":
		return NewGvalstringsEvaluable(expr)
	}

	return nil, errors.Errorf("unknown Evaluable type [%s]", typ)
}

type Evaluable interface {
	Evaluate(v interface{}) (interface{}, error)
}
