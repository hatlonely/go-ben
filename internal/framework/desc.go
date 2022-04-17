package framework

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/hatlonely/go-kit/refx"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

type ResDesc struct {
	GroupBy string
	Success string
}

type StepDesc struct {
	Name string
	Ctx  string
	Req  interface{}
	Res  ResDesc
}

type UnitDesc struct {
	Name string
	Seed map[string]string
	Step []StepDesc
}

type GroupDesc struct {
	Seconds     int
	Parallel    []int
	Quantile    []float64
	MaxStepSize int
}

type PlanDesc struct {
	Name    string
	Group   []GroupDesc
	Unit    []UnitDesc
	Monitor map[string]refx.TypeOptions
}

type CtxDesc struct {
	Name        string
	Description string
	Var         interface{}
	Ctx         map[string]refx.TypeOptions
	Seed        map[string]refx.TypeOptions
	Plan        []PlanDesc
}

func (f *Framework) LoadCtx(defaultName string, filepath string) (*CtxDesc, error) {
	buf, err := ReadFileOrNil(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "ReadFileOrNil failed")
	}
	if len(buf) == 0 {
		return &CtxDesc{Name: defaultName}, nil
	}

	var ctx CtxDesc
	if err := yaml.Unmarshal(buf, &ctx); err != nil {
		return nil, errors.Wrapf(err, "yaml.Unmarshal failed")
	}

	if len(ctx.Name) == 0 {
		ctx.Name = defaultName
	}

	return &ctx, nil
}

func (f *Framework) LoadVar(filepath string) (interface{}, error) {
	buf, err := ReadFileOrNil(filepath)
	if err != nil {
		return nil, errors.Wrap(err, "ReadFileOrNil failed")
	}
	var v interface{}
	if err := yaml.Unmarshal(buf, &v); err != nil {
		return nil, errors.Wrap(err, "yaml.Unmarshal failed")
	}
	return v, nil
}

func (f *Framework) LoadPlan(directory string, filename string) (*PlanDesc, error) {
	buf, err := ReadFileOrNil(path.Join(directory, filename))
	if err != nil {
		return nil, errors.Wrapf(err, "ioutil.ReadFile failed")
	}

	var plan PlanDesc
	if err := yaml.Unmarshal(buf, &plan); err != nil {
		return nil, errors.Wrapf(err, "yaml.Unmarshal failed")
	}

	return &plan, nil
}

func (f *Framework) LoadDescription(filepath string) (string, error) {
	buf, err := ReadFileOrNil(filepath)
	if err != nil {
		return "", errors.Wrap(err, "ReadFileOrNil failed")
	}
	return string(buf), nil
}

func ReadFileOrNil(filepath string) ([]byte, error) {
	stat, err := os.Stat(filepath)
	if errors.Is(err, os.ErrNotExist) || (err == nil && stat.IsDir()) {
		return nil, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "os.Stat [%s] failed", filepath)
	}
	buf, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, errors.Wrapf(err, "ioutil.ReadFile [%s] failed", filepath)
	}
	return buf, nil
}
