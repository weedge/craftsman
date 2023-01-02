package gw

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/weedge/craftsman/cloudwego/payment/internal/gw"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "station",
		Short: "start station server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			server, err := gw.NewServer(ctx)
			if err != nil {
				return err
			}
			return server.Run(ctx)
		},
	}
}
