package cmd

import (
	"errors"

	"github.com/TylerBrock/saw/blade"
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

var watchConfig config.Configuration

var WatchCommand = &cobra.Command{
	Use:   "watch <log group>",
	Short: "Continously stream log events",
	Long:  "",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("watching streams requires log group argument")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		watchConfig.Group = args[0]
		b := blade.NewBlade(&watchConfig)
		b.StreamEvents()
	},
}

func init() {
	WatchCommand.Flags().StringVar(&watchConfig.Prefix, "prefix", "", "log group prefix filter")
	WatchCommand.Flags().StringVar(&watchConfig.Filter, "filter", "", "event filter pattern")
	WatchCommand.Flags().BoolVar(&watchConfig.Expand, "expand", false, "indent JSON log messages")
	WatchCommand.Flags().BoolVar(&watchConfig.Invert, "invert", false, "invert colors for light terminal themes")
	WatchCommand.Flags().BoolVar(&watchConfig.RawString, "rawString", false, "print JSON strings without escaping")
}
