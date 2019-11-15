package cmd

import (
	"fmt"
	"regexp"
	"sync"

	"github.com/TylerBrock/saw/blade"
	"github.com/TylerBrock/saw/config"
	"github.com/spf13/cobra"
)

// SawCommand is the main top-level command
var SawCommand = &cobra.Command{
	Use:   "saw <command>",
	Short: "A fast, multipurpose tool for AWS CloudWatch Logs",
	Long:  "Saw is a fast, multipurpose tool for AWS CloudWatch Logs.",
	Example: `  saw groups
  saw streams production
  saw watch production`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.HelpFunc()(cmd, args)
	},
}

var awsConfig config.AWSConfiguration

func init() {
	SawCommand.AddCommand(groupsCommand)
	SawCommand.AddCommand(streamsCommand)
	SawCommand.AddCommand(versionCommand)
	SawCommand.AddCommand(watchCommand)
	SawCommand.AddCommand(getCommand)
	SawCommand.PersistentFlags().StringVar(&awsConfig.Region, "region", "", "override profile AWS region")
	SawCommand.PersistentFlags().StringVar(&awsConfig.Profile, "profile", "", "override default AWS profile")
}

func runMultiGroup(pattern string, fn func(string)) error {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return fmt.Errorf("invalid <log group> pattern: %s", err)
	}

	b := blade.NewBlade(&config.Configuration{}, &awsConfig, nil)

	var wg sync.WaitGroup
	for _, group := range b.GetLogGroups() {
		name := group.LogGroupName
		if name == nil || !re.MatchString(*name) {
			continue
		}

		wg.Add(1)
		go func(group string) {
			defer wg.Done()
			fn(group)
		}(*name)
	}

	wg.Wait()
	return nil
}
