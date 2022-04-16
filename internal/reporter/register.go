package reporter

func init() {
	RegisterReporter("Json", NewJsonReporterWithOptions)
}
