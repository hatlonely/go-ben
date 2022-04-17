package refcli

import (
	"fmt"
	"testing"

	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
	. "github.com/smartystreets/goconvey/convey"
)

func TestNewShellClient(t *testing.T) {
	Convey("TestNewShellClient", t, func() {
		client, err := NewShellClientWithOptions(&ShellClientOptions{
			Shebang: "bash",
			Args:    []string{"-c"},
			Envs: map[string]string{
				"KEY1": "hello",
			},
		})

		So(err, ShouldBeNil)

		req := &ShellClientReq{
			Command: "echo -n ${KEY1} ${KEY2}",
			Envs: map[string]string{
				"KEY2": "world",
			},
			Decoder: "text",
		}

		So(client.name(req), ShouldEqual, "echo")

		res, err := client.do(req)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, &ShellClientRes{
			Stdout:   "hello world",
			Stderr:   "",
			ExitCode: 0,
		})
	})
}

func TestShellClient_Do(t *testing.T) {
	Convey("TestShellClient_Do", t, func() {
		client, err := NewClientWithOptions(&refx.TypeOptions{
			Type: "Shell",
			Options: &ShellClientOptions{
				Args:    []string{"-c"},
				Shebang: "bash",
			},
		})

		name, res, err := client.Do(map[string]interface{}{
			"command": "date +%s",
		})
		So(err, ShouldBeNil)
		So(name, ShouldEqual, "date")
		fmt.Println(strx.JsonMarshalIndent(res))
	})
}
