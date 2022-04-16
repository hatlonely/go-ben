package framework

import (
	"context"
	"time"

	"github.com/hatlonely/go-ben/internal/refcli"
	"github.com/hatlonely/go-ben/internal/seeder"
	"github.com/hatlonely/go-ben/internal/stat"
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
	}, nil
}

type Framework struct {
	options *Options
}

type Runtime struct {
	clientMap map[string]refcli.Client
	seederMap map[string]seeder.Seeder
	variables interface{}
}

func (f *Framework) RunUnit(
	ctx context.Context, runtime Runtime, parallel int, groupDesc *GroupDesc, unitDesc *UnitDesc,
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

func (f *Framework) RunStep(runtime Runtime, unitDesc *UnitDesc) *stat.StepStat {
	stepStat := stat.NewStepStat()
	var seed map[string]interface{}
	for k, v := range unitDesc.Seed {
		seed[k] = runtime.seederMap[v].Seed()
	}

	for _, step := range unitDesc.Step {
		req := step.Req
		client := runtime.clientMap[step.Ctx]
		name := client.Name()
		now := time.Now()
		res, err := client.Do(req)
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