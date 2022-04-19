package reporter

import "github.com/hatlonely/go-ben/internal/stat"

type NoneReporter struct{}

func (r *NoneReporter) Report(test *stat.TestStat) string {
	return ""
}
