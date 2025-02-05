package cmd

import (
	"echonic/network_mapper"
	"fmt"
	"github.com/spf13/cobra"
)

var mapCmd = &cobra.Command{
	Use:   "map [target]",
	Short: "Map the network and discover connected devices",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		fmt.Println("Mapping network for target:", target)
		network_mapper.ScanNetwork(target)
	},
}
