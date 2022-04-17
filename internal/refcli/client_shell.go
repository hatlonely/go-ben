package refcli

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
)

type ShellClientOptions struct {
	Shebang string            `yaml:"shebang,omitempty" dft:"bash"`
	Args    []string          `yaml:"args,omitempty" dft:"['-c']"`
	Envs    map[string]string `yaml:"envs,omitempty"`
}

func NewShellClientWithOptions(options *ShellClientOptions) (*ShellClient, error) {
	var envs []string
	for k, v := range options.Envs {
		envs = append(envs, fmt.Sprintf(`%s=%s`, k, strings.TrimSpace(v)))
	}

	return &ShellClient{
		Shebang: options.Shebang,
		Args:    options.Args,
		Envs:    envs,
	}, nil
}

type ShellClient struct {
	Shebang string
	Args    []string
	Envs    []string
}

type ShellClientReq struct {
	Command string            `json:"command,omitempty"`
	Envs    map[string]string `json:"envs,omitempty"`
	Decoder string            `json:"decoder,omitempty"`
}

type ShellClientRes struct {
	Stdout   string `json:"stdout,omitempty"`
	Stderr   string `json:"stderr,omitempty"`
	ExitCode int    `json:"ExitCode,omitempty"`
}

func (c *ShellClient) Do(reqv interface{}) (string, interface{}, error) {
	buf, err := jsoniter.Marshal(reqv)
	if err != nil {
		return "", nil, errors.Wrap(err, "jsoniter.Marshal failed")
	}
	var req ShellClientReq
	if err := jsoniter.Unmarshal(buf, &req); err != nil {
		return "", nil, errors.Wrap(err, "jsoniter.Unmarshal failed")
	}

	res, err := c.do(&req)
	if err != nil {
		return "", nil, errors.WithMessage(err, "c.do failed")
	}

	buf, err = jsoniter.Marshal(res)
	if err != nil {
		return "", nil, errors.Wrap(err, "jsoniter.Marshal failed")
	}
	var resv interface{}
	if err := jsoniter.Unmarshal(buf, &resv); err != nil {
		return "", nil, errors.Wrap(err, "jsoniter.Unmarshal failed")
	}

	return strings.Fields(req.Command)[0], resv, nil
}

func (c *ShellClient) name(req *ShellClientReq) string {
	return strings.Fields(req.Command)[0]
}

func (c *ShellClient) do(req *ShellClientReq) (*ShellClientRes, error) {
	var envs []string
	for k, v := range req.Envs {
		envs = append(envs, fmt.Sprintf(`%s=%s`, k, strings.TrimSpace(v)))
	}

	cmd := exec.Command(c.Shebang, append(c.Args, req.Command)...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, c.Envs...)
	cmd.Env = append(cmd.Env, envs...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Start(); err != nil {
		return nil, errors.Wrap(err, "cmd.Start failed")
	}

	if err := cmd.Wait(); err != nil {
		exitCode := -1
		if e, ok := err.(*exec.ExitError); ok {
			if status, ok := e.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}

		return &ShellClientRes{
			Stdout:   stdout.String(),
			Stderr:   stderr.String(),
			ExitCode: exitCode,
		}, errors.Wrap(err, "cmd.Wait failed")
	}

	return &ShellClientRes{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: 0,
	}, nil
}
