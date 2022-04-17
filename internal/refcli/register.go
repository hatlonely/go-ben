package refcli

func init() {
	RegisterClient("Shell", NewShellClientWithOptions)
}
