package da

import (
	"context"

	"github.com/spf13/cobra"
	"github.com/weedge/craftsman/cloudwego/payment/internal/da"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "da",
		Short: "start da server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			server, err := da.NewServer(ctx)
			if err != nil {
				return err
			}
			return server.Run(ctx)
		},
	}
}
