package refcli

import (
	"testing"

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

		So(client.Name(req), ShouldEqual, "echo")

		res, err := client.Do(req)
		So(err, ShouldBeNil)
		So(res, ShouldResemble, &ShellClientRes{
			Stdout:   "hello world",
			Stderr:   "",
			ExitCode: 0,
		})
	})
}
