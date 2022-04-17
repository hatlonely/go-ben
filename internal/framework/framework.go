package framework

import (
	"context"
	"encoding/hex"
	"path"
	"strings"
	"time"

	"github.com/hatlonely/go-ben/internal/util"

	"github.com/hatlonely/go-ben/internal/refcli"
	"github.com/hatlonely/go-ben/internal/seeder"
	"github.com/hatlonely/go-ben/internal/stat"
	uuid "github.com/satori/go.uuid"
)

type Options struct {
	TestDirectory string
	PlanDirectory string
	Customize     string
	Reporter      string
	X             string
	JsonStat      string
	Hook          string
	Lang          string
}

func NewFrameworkWithOptions(options *Options) (*Framework, error) {
	return &Framework{
		options: options,
		id:      hex.EncodeToString(uuid.NewV4().Bytes()),
	}, nil
}

type Framework struct {
	options *Options

	id string
}

type Runtime struct {
	clientMap map[string]refcli.Client
	seederMap map[string]seeder.Seeder
	variables interface{}
}

func (f *Framework) Run() *stat.TestStat {
	return f.RunTest(f.options.TestDirectory, &Runtime{
		clientMap: nil,
		seederMap: nil,
		variables: nil,
	})
}

const (
	LoadingFileCtx         = "ctx.yaml"
	LoadingFileVar         = "var.yaml"
	LoadingFileDescription = "README.md"
)

func (f *Framework) RunTest(directory string, runtime *Runtime) *stat.TestStat {
	defaultName := path.Base(directory)
	ctxDesc, err := f.LoadCtx(defaultName, path.Join(directory, LoadingFileCtx))
	if err != nil {
		return stat.NewTestStat(f.id, directory, defaultName, "").SetError(err)
	}
	description, err := f.LoadDescription(path.Join(directory, LoadingFileDescription))
	if err != nil {
		return stat.NewTestStat(f.id, directory, defaultName, "").SetError(err)
	}
	variables, err := f.LoadVar(path.Join(directory, LoadingFileVar))
	if err != nil {
		return stat.NewTestStat(f.id, directory, defaultName, "").SetError(err)
	}
	variables = util.MustMerge(ctxDesc.Var, variables)
	clientMap := map[string]refcli.Client{}
	for k, v := range runtime.clientMap {
		clientMap[k] = v
	}
	for k, v := range ctxDesc.Ctx {
		cli, err := refcli.NewClientWithOptions(&v)
		if err != nil {
			return stat.NewTestStat(f.id, directory, defaultName, "").SetError(err)
		}
		clientMap[k] = cli
	}
	seederMap := map[string]seeder.Seeder{}
	for k, v := range runtime.seederMap {
		seederMap[k] = v
	}
	for k, v := range ctxDesc.Seed {
		s, err := seeder.NewSeederWithOptions(&v)
		if err != nil {
			return stat.NewTestStat(f.id, directory, defaultName, "").SetError(err)
		}
		seederMap[k] = s
	}

	testStat := stat.NewTestStat(f.id, directory, ctxDesc.Name, description)

	if strings.HasPrefix(directory, f.options.PlanDirectory) {
		for _, plan := range ctxDesc.Plan {
			planStat := f.RunPlan(&Runtime{
				clientMap: clientMap,
				seederMap: seederMap,
				variables: variables,
			}, "", plan)

			testStat.AddPlanStat(planStat)
		}
	}

	return testStat
}

func (f *Framework) RunPlan(runtime *Runtime, planID string, plan PlanDesc) *stat.PlanStat {
	planStat := stat.NewPlanStat(planID, plan.Name)

	for idx, groupDesc := range plan.Group {
		unitGroupStat := f.RunUnitGroup(runtime, idx, &groupDesc, plan.Unit)
		planStat.AddUnitGroupStat(unitGroupStat)
	}

	return planStat
}

func (f *Framework) RunUnitGroup(runtime *Runtime, groupIdx int, groupDesc *GroupDesc, unit []UnitDesc) *stat.UnitGroupStat {
	ctx, cancel := context.WithCancel(context.Background())

	unitStatChan := make(chan *stat.UnitStat, len(unit))
	for idx, unitDesc := range unit {
		go func(idx int, unitDesc *UnitDesc) {
			f.RunUnit(ctx, runtime, groupDesc.Parallel[idx], groupDesc, unitDesc, unitStatChan)
		}(idx, &unitDesc)
	}

	time.Sleep(time.Duration(groupDesc.Seconds) * time.Second)
	cancel()

	unitGroupStat := stat.NewUnitGroupStat(groupIdx, groupDesc.Seconds, 0, groupDesc.Quantile)
	for range unit {
		unitGroupStat.AddUnitStat(<-unitStatChan)
	}
	return unitGroupStat
}

func (f *Framework) RunUnit(
	ctx context.Context, runtime *Runtime, parallel int, groupDesc *GroupDesc, unitDesc *UnitDesc,
	unitStatChan chan<- *stat.UnitStat,
) {
	stepStatChan := make(chan *stat.StepStat, parallel*2)
	for i := 0; i < parallel; i++ {
		go func() {
			for {
				select {
				case <-ctx.Done():
					return
				default:
					stepStatChan <- f.RunStep(runtime, unitDesc)
				}
			}
		}()
	}
	go func() {
		unitStat := stat.NewUnitStat(unitDesc.Name, parallel, 0, int64(groupDesc.Seconds), 0, 100, groupDesc.Quantile, groupDesc.MaxStepSize)
		for {
			select {
			case <-ctx.Done():
				unitStat.Summary()
				unitStatChan <- unitStat
				return
			case stepStat := <-stepStatChan:
				unitStat.AddStepStat(stepStat)
			}
		}
	}()
}

func (f *Framework) RunStep(runtime *Runtime, unitDesc *UnitDesc) *stat.StepStat {
	stepStat := stat.NewStepStat()
	var seed map[string]interface{}
	for k, v := range unitDesc.Seed {
		seed[k] = runtime.seederMap[v].Seed()
	}

	for _, step := range unitDesc.Step {
		req := step.Req
		client := runtime.clientMap[step.Ctx]
		now := time.Now()
		name, res, err := client.Do(req)
		if err != nil {
			stepStat.AddErrStat(name, err)
			return stepStat
		}
		stepStat.AddSubStepStat(&stat.SubStepStat{
			Req:     req,
			Res:     res,
			Name:    name,
			Code:    "",
			Success: false,
			Elapse:  time.Now().Sub(now),
		})
	}
	return stepStat
}
