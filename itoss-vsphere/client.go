package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/vmware/govmomi/session/cache"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
	"os"
)

var urlFlag = flag.String("url", "", "-url <https://username:password@host/sdk>")
var insecureFlag = flag.Bool("insecure", false, "-insecure")
var hostFlag = flag.String("host", "", "-host <hostname> (if not specified list all hosts)")
var statusFlag = flag.Bool("status", false, "-status")

// NewClient creates a vim25.Client for use in the examples
func NewClient(ctx context.Context) (*vim25.Client, error) {
	// Parse URL from string

	u, err := soap.ParseURL(*urlFlag)

	if err != nil {
		return nil, err
	}

	s := &cache.Session{
		URL:      u,
		Insecure: *insecureFlag,
	}

	c := new(vim25.Client)
	err = s.Login(ctx, c, nil)
	if err != nil {
		return nil, err
	}

	return c, nil
}

// Run calls f with Client create from the -url flag if provided,
// otherwise runs the example against vcsim.
func Run(f func(context.Context, *vim25.Client) error) {

	flag.Parse()
	var err error
	var c *vim25.Client

	//if *urlFlag == "" {
	//	err = simulator.VPX().Run(f)
	//	os.Exit(0)
	//} else {
	if *urlFlag == "" {
		fmt.Fprint(os.Stdout, "You must specify an url.\n")
		flag.Usage()
		os.Exit(0)
	}
	ctx := context.Background()
	c, err = NewClient(ctx)
	if err == nil {
		err = f(ctx, c)
	}
	//}

	if *statusFlag == false {
		fmt.Fprint(os.Stdout, "Option not implemented, set status flag.\n")
		flag.Usage()
		os.Exit(0)
	}
	if err != nil {
		fmt.Fprint(os.Stdout, -4)
		fmt.Fprintf(os.Stderr, "\nError: %s\n", err)
		os.Exit(0)
	}
}
