package eval

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEvaluable_Evaluate(t *testing.T) {
	Convey("TestEvaluable_Evaluate", t, func() {
		params := map[string]interface{}{
			"x": "hello",
			"y": "world",
		}

		for _, unit := range []struct {
			Type string
			Expr string
		}{
			{"Tengo", "x + y"},
			{"Yaegi", `
func eval(v map[string]interface{}) (interface{}, error) {
	return v["x"].(string) + v["y"].(string), nil
}
`},
			{"Govaluate", `x + y`},
			{"Gval", `x + y`},
		} {
			eval, err := NewEvaluable(unit.Type, unit.Expr)
			So(err, ShouldBeNil)
			val, err := eval.Evaluate(params)
			So(err, ShouldBeNil)
			So(val, ShouldEqual, "helloworld")
		}
	})
}
