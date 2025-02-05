package cmd

import (
	"github.com/spf13/cobra"
)

// 根命令
var rootCmd = &cobra.Command{
	Use:   "echonic",
	Short: "Echonic is a penetration testing tool",
	Long:  `Echonic is a tool for port scanning, service detection, and network mapping.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 初始化子命令
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(bruteCmd)
	rootCmd.AddCommand(mapCmd)
}
