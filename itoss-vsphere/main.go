package main

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"os"
)

func main() {
	Run(func(ctx context.Context, c *vim25.Client) error {
		// Create a view of HostSystem objects
		m := view.NewManager(c)
		v, err := m.CreateContainerView(ctx, c.ServiceContent.RootFolder, []string{"HostSystem"}, true)
		if err != nil {
			return err
		}
		defer v.Destroy(ctx)
		var hss []mo.HostSystem

		err = v.Retrieve(ctx, []string{"HostSystem"}, []string{"summary"}, &hss)

		if err != nil {
			return err
		}
		hostFound := false
		printedTitles := false
		for _, hs := range hss {
			if *hostFlag != "" && hs.Summary.Config.Name != *hostFlag {
				continue
			}
			if !printedTitles {
				fmt.Fprint(os.Stdout, "host;uptimeSec;overallStatus;connectionState;inMaintenanceMode;powerState;standbyMode;bootTime\n")
				printedTitles = true
			}
			fmt.Fprintf(os.Stdout, "%s;%d;%s;%s;%v;%s;%s;%s\n",
				hs.Summary.Config.Name,
				hs.Summary.QuickStats.Uptime,
				hs.Summary.OverallStatus,
				hs.Summary.Runtime.ConnectionState,
				hs.Summary.Runtime.InMaintenanceMode,
				hs.Summary.Runtime.PowerState,
				hs.Summary.Runtime.StandbyMode,
				hs.Summary.Runtime.BootTime.Format("2006-01-02 15:04:05"))
			//
			hostFound = true
		}
		if !hostFound {
			fmt.Fprint(os.Stdout, -5)
			fmt.Fprintf(os.Stderr, "\nHost %s not found\n", *hostFlag)
			os.Exit(0)
		}
		return nil
	})
}
