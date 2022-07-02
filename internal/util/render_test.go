package util

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRender(t *testing.T) {
	Convey("TestRender", t, func() {
		v, err := Render(map[string]interface{}{
			"#key1": "strkey",
			"#key2": "intkey",
		}, map[string]interface{}{
			"strkey": "strval",
			"intkey": 123,
		})
		So(err, ShouldBeNil)
		So(v, ShouldResemble, map[string]interface{}{
			"key1": "strval",
			"key2": 123,
		})
	})
}
