package main

import (
	"flag"
	"fmt"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/spf13/cobra"
	"github.com/weedge/craftsman/cloudwego/payment/cmd/da"
	"github.com/weedge/craftsman/cloudwego/payment/cmd/station"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/configparser"
	"github.com/weedge/craftsman/cloudwego/payment/pkg/version"
)

var (
	moduleName = version.Get().Module
	rootCmd    = &cobra.Command{
		Use:   moduleName,
		Short: fmt.Sprintf("%s module", moduleName),
	}
)

func main() {
	rootCmd.PersistentFlags().AddGoFlagSet(flag.CommandLine)
	rootCmd.AddCommand(
		da.NewCommand(),
		station.NewCommand(),
	)

	configparser.Flags(rootCmd.PersistentFlags())
	if err := rootCmd.Execute(); err != nil {
		klog.Fatal(err.Error())
	}
}
