package framework

import (
	"io/ioutil"
	"os"

	"github.com/hatlonely/go-kit/refx"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type Ctx struct {
	Name        string
	Description string
	Var         interface{}
	Ctx         map[string]refx.Options
	Seed        map[string]refx.Options
	Plan        []Plan
}

func (f *Framework) LoadCtx(defaultName string, filepath string) (*Ctx, error) {
	stat, err := os.Stat(filepath)
	if errors.Is(err, os.ErrNotExist) || (err == nil && stat.IsDir()) {
		return &Ctx{
			Name:        defaultName,
			Description: "",
			Var:         nil,
			Ctx:         nil,
			Seed:        nil,
			Plan:        nil,
		}, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "os.Stat [%s] failed", filepath)
	}

	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrapf(err, "ioutil.ReadFile failed")
	}

	var ctx Ctx
	if err := yaml.Unmarshal(buf, &ctx); err != nil {
		return nil, errors.Wrapf(err, "yaml.Unmarshal failed")
	}

	if len(ctx.Name) == 0 {
		ctx.Name = defaultName
	}

	return &ctx, nil
}
