package main

import (
	"os"

	"github.com/hatlonely/go-ben/internal/framework"
	"github.com/hatlonely/go-kit/flag"
	"github.com/hatlonely/go-kit/refx"
	"github.com/hatlonely/go-kit/strx"
)

var Version = "Unknown"

type Options struct {
	Help    bool `flag:"-h; usage: show help info"`
	Version bool `flag:"-v; usage: show version"`

	framework.Options
}

const (
	ECOK = 0
	ECErr
)

func main() {
	var options Options
	refx.Must(flag.Struct(&options, refx.WithKebabName()))
	refx.Must(flag.Parse(flag.WithJsonVal()))
	if options.Help {
		strx.Trac(flag.Usage())
		strx.Trac(`
  ben -t ops/example
`)
		os.Exit(ECOK)
	}
	if options.Version {
		strx.Trac(Version)
		os.Exit(ECOK)
	}

	fw, err := framework.NewFrameworkWithOptions(&options.Options)
	refx.Must(err)

	if len(options.JsonStat) != 0 {
		fw.Format()
		os.Exit(ECOK)
	} else if fw.Run() {
		os.Exit(ECOK)
	}

	os.Exit(ECErr)
}
