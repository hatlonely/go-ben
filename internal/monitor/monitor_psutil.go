package monitor

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/net"

	"github.com/hatlonely/go-ben/internal/stat"
)

type PsutilMonitorOptions struct {
	Interval         time.Duration
	Metrics          []string
	NetworkInterface string
	DiskPath         string
}

func NewPsutilMonitorWithOptions(options *PsutilMonitorOptions) (*PsutilMonitor, error) {
	if len(options.DiskPath) == 0 {
		options.DiskPath = "/"
	}

	unit := map[string]string{
		"CPU":    "percent",
		"Mem":    "byte",
		"Disk":   "byte",
		"IOR":    "times",
		"IOW":    "times",
		"NetIOR": "bit",
		"NetIOW": "bit",
	}
	metrics := stat.NewMonitorStat()
	for _, metric := range options.Metrics {
		if v, ok := unit[metric]; !ok {
			return nil, errors.Errorf("unknown metric [%s]", metric)
		} else {
			metrics.Unit[metric] = v
		}
		metrics.Stat[metric] = nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &PsutilMonitor{
		options: options,
		ctx:     ctx,
		cancel:  cancel,
	}, nil
}

type PsutilMonitor struct {
	options *PsutilMonitorOptions
	metrics *stat.MonitorStat

	ctx    context.Context
	cancel context.CancelFunc
}

func (m *PsutilMonitor) Collect() {
	go func() {
		ticker := time.NewTicker(m.options.Interval)
		now := time.Now()
		defer ticker.Stop()
		for {
			select {
			case <-m.ctx.Done():
				return
			case <-ticker.C:
				if _, ok := m.metrics.Stat["CPU"]; ok {
					vs, _ := cpu.Percent(time.Now().Sub(now), false)
					m.metrics.Stat["CPU"] = append(m.metrics.Stat["CPU"], stat.Measurement{
						Time:  now,
						Value: vs[0],
					})
				}
				if _, ok := m.metrics.Stat["Mem"]; ok {
					vm, _ := mem.VirtualMemory()
					m.metrics.Stat["Mem"] = append(m.metrics.Stat["Mem"], stat.Measurement{
						Time:  now,
						Value: float64(vm.Used),
					})
				}
				if _, ok := m.metrics.Stat["Disk"]; ok {
					du, _ := disk.Usage(m.options.DiskPath)
					m.metrics.Stat["Disk"] = append(m.metrics.Stat["Disk"], stat.Measurement{
						Time:  now,
						Value: float64(du.Used),
					})
				}

				_, ok1 := m.metrics.Stat["NetIOR"]
				_, ok2 := m.metrics.Stat["NetIOW"]
				if ok1 || ok2 {
					nios, _ := net.IOCounters(true)
					for _, nio := range nios {
						if nio.Name != m.options.NetworkInterface {
							continue
						}
						if ok1 {
							m.metrics.Stat["NetIOR"] = append(m.metrics.Stat["NetIOR"], stat.Measurement{
								Time:  now,
								Value: float64(nio.BytesRecv),
							})
						}
						if ok2 {
							m.metrics.Stat["NetIOW"] = append(m.metrics.Stat["NetIOW"], stat.Measurement{
								Time:  now,
								Value: float64(nio.BytesSent),
							})
						}
						break
					}
				}

				now = time.Now()
			}
		}
	}()
}

func (m *PsutilMonitor) Stat(start time.Time, end time.Time) *stat.MonitorStat {
	m.cancel()
	return m.metrics
}
