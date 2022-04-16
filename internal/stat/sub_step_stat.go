package stat

import "time"

type SubStepStat struct {
	Req     interface{}
	Res     interface{}
	Name    string
	Code    string
	Success bool
	Elapse  time.Duration
}
