package cmd

import (
	"echonic/bruteforce"
	"fmt"
	"github.com/spf13/cobra"
)

var bruteCmd = &cobra.Command{
	Use:   "brute [target] [ports]",
	Short: "Perform a brute force attack on specified ports",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		ports := args[1]
		fmt.Printf("Brute forcing target: %s on ports: %s\n", target, ports)
		bruteforce.Attack(target) // Example ports
	},
}
