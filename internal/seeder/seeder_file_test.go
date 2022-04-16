package seeder

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFileSeeder(t *testing.T) {
	Convey("TestFileSeeder", t, func() {
		So(ioutil.WriteFile("test.json", []byte(`
{"Key1": "Val1", "Key2": "Val2"}
{"Key1": "Val1", "Key2": "Val2"}
{"Key1": "Val1", "Key2": "Val2"}
`), 0644), ShouldBeNil)
		defer os.RemoveAll("test.json")

		seeder, err := NewFileSeederWithOptions(&FileSeederOptions{
			Name:             "test.json",
			IgnoreParseError: false,
		})

		So(err, ShouldBeNil)

		fmt.Println(seeder.Seed())
	})
}
