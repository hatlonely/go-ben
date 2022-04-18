package i18n

type Title struct {
	Report        string
	Test          string
	Plan          string
	UnitGroup     string
	Unit          string
	Step          string
	Err           string
	Idx           string
	Seconds       string
	Times         string
	Parallel      string
	Limit         string
	Total         string
	Elapse        string
	Rate          string
	ResTime       string
	QPS           string
	Code          string
	Summary       string
	Quantile      string
	QuantileShort string
	Monitor       string
}

type Status struct {
	Fail string
	Succ string
	Skip string
}

type Tooltip struct {
	Save string
	Copy string
}

type I18n struct {
	Title   Title
	Status  Status
	Tooltip Tooltip
}

var defaultI18n = map[string]*I18n{
	"dft": {
		Title: Title{
			Report:        "Report",
			Test:          "Test",
			Plan:          "Plan",
			UnitGroup:     "Unit Group",
			Unit:          "Unit",
			Step:          "Step",
			Err:           "Err",
			Idx:           "Index",
			Seconds:       "Seconds",
			Times:         "Times",
			Parallel:      "Parallel",
			Limit:         "Limit",
			Total:         "Total",
			Elapse:        "Elapse",
			Rate:          "Rate",
			ResTime:       "ResTime",
			QPS:           "QPS",
			Code:          "Code",
			Summary:       "Summary",
			Quantile:      "Quantile",
			QuantileShort: "Q",
			Monitor:       "Monitor",
		},
		Status: Status{
			Fail: "Fail",
			Succ: "Succ",
			Skip: "Skip",
		},
		Tooltip: Tooltip{
			Save: "Save",
			Copy: "Copy",
		},
	},
	"zh": {
		Title: Title{
			Report:        "报告",
			Test:          "测试",
			Plan:          "计划",
			UnitGroup:     "单元组",
			Unit:          "单元",
			Step:          "步骤",
			Err:           "错误",
			Idx:           "序列",
			Seconds:       "测试时长",
			Times:         "测试次数",
			Parallel:      "并发",
			Limit:         "限流",
			Total:         "总共",
			Elapse:        "耗时",
			Rate:          "成功率",
			ResTime:       "响应时间",
			QPS:           "QPS",
			Code:          "错误码",
			Summary:       "汇总",
			Quantile:      "分位数",
			QuantileShort: "Q",
			Monitor:       "监测",
		},
		Status: Status{
			Fail: "失败",
			Succ: "成功",
			Skip: "跳过",
		},
		Tooltip: Tooltip{
			Save: "保存",
			Copy: "复制",
		},
	},
}

func NewI18n(lang string, i18n *I18n) *I18n {
	return defaultI18n["dft"]
}
