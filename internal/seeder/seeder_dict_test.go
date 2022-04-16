package seeder

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDictSeeder(t *testing.T) {
	Convey("TestDictSeeder", t, func() {
		seeder, err := NewDictSeederWithOptions(&DictSeederOptions{
			{"Key1": "val1", "Key2": "val2"},
			{"Key1": "val3", "Key2": "val4"},
			{"Key1": "val5", "Key2": "val6"},
		})
		So(err, ShouldBeNil)

		fmt.Println(seeder.Seed())
	})
}
