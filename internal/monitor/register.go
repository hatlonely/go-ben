package monitor

func init() {
	RegisterMonitor("Psutil", NewPsutilMonitorWithOptions)
}
