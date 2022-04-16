package refcli

type Client interface {
	Name() string
	Do(req interface{}) (interface{}, error)
}
