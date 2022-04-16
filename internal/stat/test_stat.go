package stat

type TestStat struct {
	ID          string
	Directory   string
	Name        string
	Description string
	IsErr       string
	Err         string
	SubTests    *TestStat
}
