package cmd

import (
	"echonic/scanner"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var scanCmd = &cobra.Command{
	Use:     "scan [target]",
	Example: "scan 127.0.0.1",
	Short:   "Scan a target for open ports and services",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		target := args[0]
		timeout := 2 * time.Second
		maxConcurrency := 1000 // 建议降低并发数，避免系统资源耗尽

		openPorts, err := scanner.PortScan(target, timeout, maxConcurrency)
		if err != nil {
			log.Printf("Scan error: %v\n", err)
			return
		}

		if len(openPorts) == 0 {
			fmt.Println("No open ports found")
			return
		}

		for _, port := range openPorts {
			service := scanner.DetectService(port)
			fmt.Printf("Port: %d, Service: %s\n", port, service)
		}
	},
}
