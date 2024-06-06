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

		fmt.Fprint(os.Stdout, "host;uptimeSec;overallStatus;connectionState;inMaintenanceMode;powerState;standbyMode;bootTime;proxyStatus\n")
		for _, hs := range hss {
			if *entityNameFlag != "all" && hs.Summary.Config.Name != *entityNameFlag {
				continue
			}

			fmt.Fprintf(os.Stdout, "%s;%d;%s;%s;%v;%s;%s;%s;%s\n",
				hs.Summary.Config.Name,
				hs.Summary.QuickStats.Uptime,
				hs.Summary.OverallStatus,
				hs.Summary.Runtime.ConnectionState,
				hs.Summary.Runtime.InMaintenanceMode,
				hs.Summary.Runtime.PowerState,
				hs.Summary.Runtime.StandbyMode,
				hs.Summary.Runtime.BootTime.Format("2006-01-02 15:04:05"),
				"OK")
			//
			hostFound = true
		}
		if !hostFound {
			fmt.Fprintf(os.Stdout, "%s;%d;%s;%s;%v;%s;%s;%s;%s\n",
				*entityNameFlag, 0, "NA", "NA", false, "NA", "NA", "NA", "HOST_NOT_FOUND")
			fmt.Fprintf(os.Stderr, "\nHost %s not found\n", *entityNameFlag)
			os.Exit(0)
		}
		return nil
	})
}
